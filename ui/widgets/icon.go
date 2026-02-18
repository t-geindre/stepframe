package widgets

import (
	"image"
	"image/color"
	"stepframe/ui/theme"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Icon struct {
	*widget.Widget
	icon *ebiten.Image
	opts *ebiten.DrawImageOptions
	size int

	color      color.Color
	pulseColor color.Color
	lastPulse  time.Time
}

func NewIcon(icon theme.Icon, size theme.IconSize) *Icon {
	i := &Icon{
		Widget: widget.NewWidget(),
		size:   theme.Current.IconSizes[size],
		color:  theme.Current.IconColors[theme.IconColorDefault],
	}
	i.SetIcon(icon)
	return i
}

func (i *Icon) Pulse() { i.lastPulse = time.Now() }

func (i *Icon) SetSize(size theme.IconSize) {
	i.size = theme.Current.IconSizes[size]
	i.scale()
}

func (i *Icon) SetColor(c theme.IconColor) {
	i.color = theme.Current.IconColors[c]
	// pas besoin de re-scale pour la couleur, mais ok si tu veux
}

func (i *Icon) SetPulseColor(c theme.IconColor) {
	i.pulseColor = theme.Current.IconColors[c]
}

func (i *Icon) SetIcon(icon theme.Icon) {
	i.icon = theme.Current.Icons[icon]
	i.scale()
}

func (i *Icon) GetWidget() *widget.Widget { return i.Widget }

func (i *Icon) PreferredSize() (int, int) { return i.size, i.size }

func (i *Icon) Render(screen *ebiten.Image) {
	if i.icon == nil || i.opts == nil {
		return
	}

	opts := *i.opts
	opts.ColorScale.Reset()
	opts.ColorScale.ScaleWithColor(i.currentColor())

	screen.DrawImage(i.icon, &opts)
}

func (i *Icon) SetLocation(rect image.Rectangle) {
	i.Widget.SetLocation(rect)
	i.scale()
}

func (i *Icon) Validate() {}

func (i *Icon) Update(updObj *widget.UpdateObject) {
	i.Widget.Update(updObj)
}

func (i *Icon) scale() {
	if i.icon == nil {
		i.opts = nil
		return
	}

	w := float64(i.Widget.Rect.Dx())
	h := float64(i.Widget.Rect.Dy())
	if w <= 0 || h <= 0 {
		i.opts = nil
		return
	}

	iw := float64(i.icon.Bounds().Dx())
	ih := float64(i.icon.Bounds().Dy())

	scale := min(w/iw, h/ih)

	offX := (w - iw*scale) / 2
	offY := (h - ih*scale) / 2

	i.opts = &ebiten.DrawImageOptions{Filter: ebiten.FilterLinear}
	i.opts.GeoM.Scale(scale, scale)
	i.opts.GeoM.Translate(
		float64(i.Widget.Rect.Min.X)+offX,
		float64(i.Widget.Rect.Min.Y)+offY,
	)
}

func (i *Icon) currentColor() color.NRGBA {
	if i.pulseColor == nil {
		return toNRGBA(i.color)
	}

	base := toNRGBA(i.color)
	pulse := toNRGBA(i.pulseColor)
	if pulse.A == 0 { // pulseColor pas set
		pulse = base
	}

	const PulseDuration = 300 * time.Millisecond

	if i.lastPulse.IsZero() {
		return base
	}
	elapsed := time.Since(i.lastPulse)
	if elapsed <= 0 {
		return pulse
	}
	if elapsed >= PulseDuration {
		return base
	}

	t := 1.0 - float64(elapsed)/float64(PulseDuration)
	t = smootherstep(t)
	t = t * t

	return lerpNRGBA(base, pulse, t)
}

func toNRGBA(c color.Color) color.NRGBA {
	if c == nil {
		return color.NRGBA{255, 255, 255, 255}
	}
	r, g, b, a := c.RGBA() // 0..65535
	return color.NRGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func lerpNRGBA(a, b color.NRGBA, t float64) color.NRGBA {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	lerp := func(x, y uint8) uint8 {
		return uint8(float64(x) + (float64(y)-float64(x))*t)
	}
	return color.NRGBA{
		R: lerp(a.R, b.R),
		G: lerp(a.G, b.G),
		B: lerp(a.B, b.B),
		A: lerp(a.A, b.A),
	}
}

func smootherstep(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	// 6t^5 - 15t^4 + 10t^3
	return t * t * t * (t*(t*6-15) + 10)
}
