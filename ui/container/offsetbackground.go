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

	return b.outer
}

func (b *OffsetBackground[T]) Render(screen *ebiten.Image) {
	if b.background == nil {
		return
	}

	rect := b.GetWidget().Rect
	rect.Min = rect.Min.Sub(b.background.Offset)
	rect.Max = rect.Max.Add(b.background.Offset)

	b.background.Image.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	})
}
