package widgets

import (
	"image"
	"image/color"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Icon struct {
	*widget.Widget
	icon  *ebiten.Image
	opts  *ebiten.DrawImageOptions
	size  int
	color color.Color
}

func NewIcon(icon theme.Icon, size theme.IconSize) *Icon {
	i := &Icon{
		Widget: widget.NewWidget(),
		size:   theme.Current.IconsTheme.Sizes[size],
		color:  theme.Current.Colors[theme.IconColorDefault],
	}
	i.SetIcon(icon)
	return i
}

func (i *Icon) SetSize(size theme.IconSize) {
	i.size = theme.Current.IconsTheme.Sizes[size]
	i.scale()
}

func (i *Icon) SetColor(c theme.Color) {
	i.color = theme.Current.Colors[c]
	// pas besoin de re-scale pour la couleur, mais ok si tu veux
}

func (i *Icon) SetIcon(icon theme.Icon) {
	i.icon = theme.Current.IconsTheme.Icons[icon]
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
	opts.ColorScale.ScaleWithColor(i.color)

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
