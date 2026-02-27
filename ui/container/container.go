package container

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Container[T widget.Containerer] struct {
	*widget.Container
	*Gradient[T]
	*OffsetBackground[T]
	outer T
}

func NewContainer[T widget.Containerer](layout widget.Layouter, outer T) *Container[T] {
	container := widget.NewContainer(widget.ContainerOpts.Layout(layout))

	return &Container[T]{
		Container:        container,
		Gradient:         NewGradient(container, outer),
		OffsetBackground: NewOffsetBackground(container, outer),
		outer:            outer,
	}
}

func (c *Container[T]) SetBackgroundImage(img *image.NineSlice) T {
	if c.Container.IsValidated() {
		c.Container.SetBackgroundImage(img)
	} else {
		widget.ContainerOpts.BackgroundImage(img)(c.Container)
	}

	imw, imh := img.MinSize()
	wmw, wmh := c.GetWidget().MinHeight, c.GetWidget().MinHeight

	if wmw == 0 {
		c.GetWidget().MinWidth = imw
	}

	if wmh == 0 {
		c.GetWidget().MinHeight = imh
	}

	return c.outer
}

func (c *Container[T]) SetMinSize(width, height int) T {
	widget.WidgetOpts.MinSize(width, height)(c.Container.GetWidget())
	return c.outer
}

func (c *Container[T]) Render(screen *ebiten.Image) {
	c.Gradient.Render(screen)
	c.OffsetBackground.Render(screen)
	c.Container.Render(screen)
}
