package widgets

import (
	"image/color"
	"stepframe/ui/theme"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Led struct {
	*widget.Widget
	color, pulseColor color.Color
	isOn              bool
	lastPulse         time.Time
}

func NewLed(isOn bool, color theme.Color) *Led {
	l := &Led{
		Widget: widget.NewWidget(),
		color:  theme.Current.Colors[color],
		isOn:   isOn,
	}
	return l
}

func (l *Led) GetWidget() *widget.Widget {
	return l.Widget
}

func (l *Led) PreferredSize() (int, int) {
	return theme.Current.LedTheme.Width, theme.Current.LedTheme.Height
}

func (l *Led) Validate() {
}

func (l *Led) Render(screen *ebiten.Image) {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(
		float64(l.Widget.Rect.Min.X+theme.Current.LedTheme.OffsetX),
		float64(l.Widget.Rect.Min.Y+theme.Current.LedTheme.OffsetY),
	)
	opt.ColorScale.ScaleWithColor(l.color)

	img := theme.Current.LedTheme.OnImage
	if !l.isOn {
		img = theme.Current.LedTheme.OffImage
	}
	screen.DrawImage(img, opt)

	if l.pulseColor == nil || l.lastPulse.IsZero() {
		return
	}

	elapsed := time.Since(l.lastPulse)
	if elapsed > theme.Current.LedTheme.PulseDuration {
		return
	}

	t := 1.0 - float64(elapsed)/float64(theme.Current.LedTheme.PulseDuration)
	t = smootherstep(t)
	t = t * t

	pulseOpt := &ebiten.DrawImageOptions{}
	pulseOpt.GeoM.Translate(
		float64(l.Widget.Rect.Min.X+theme.Current.LedTheme.OffsetX),
		float64(l.Widget.Rect.Min.Y+theme.Current.LedTheme.OffsetY),
	)
	pulseOpt.ColorScale.ScaleWithColor(l.pulseColor)
	pulseOpt.ColorScale.ScaleAlpha(float32(t))

	screen.DrawImage(theme.Current.LedTheme.OnImage, pulseOpt)
}

func (l *Led) SetPulseColor(color theme.Color) {
	l.pulseColor = theme.Current.Colors[color]
}

func (l *Led) SetColor(color theme.Color) {
	l.color = theme.Current.Colors[color]
}

func (l *Led) Pulse() {
	l.lastPulse = time.Now()
}
