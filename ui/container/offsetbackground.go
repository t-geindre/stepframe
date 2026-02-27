package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type OffsetBackground[T widget.Containerer] struct {
	*widget.Container
	outer      T
	background *theme.OffsetImage
}

func NewOffsetBackground[T widget.Containerer](container *widget.Container, outer T) *OffsetBackground[T] {
	return &OffsetBackground[T]{
		Container: container,
		outer:     outer,
	}
}

func (b *OffsetBackground[T]) SetBackgroundOffsetImage(img *theme.OffsetImage) T {
	b.background = img

	if img.Offset == nil {
		return b.outer
	}

	imw, imh := img.Image.MinSize()
	wmw, wmh := b.GetWidget().MinHeight, b.GetWidget().MinHeight

	if wmw == 0 {
		b.GetWidget().MinWidth = imw - b.background.Offset.Dx()
	}

	if wmh == 0 {
		b.GetWidget().MinHeight = imh - b.background.Offset.Dy()
	}

	return b.outer
}

func (b *OffsetBackground[T]) Render(screen *ebiten.Image) {
	if b.background == nil {
		return
	}

	rect := b.GetWidget().Rect
	if b.background.Offset != nil {
		rect = b.background.Offset.Apply(rect)
	}

	b.background.Image.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	})
}
