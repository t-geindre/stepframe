package container

import (
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type Container[T widget.Containerer] struct {
	*widget.Container
	outer T
}

func NewContainer[T widget.Containerer](layout widget.Layouter, outer T) *Container[T] {
	return &Container[T]{
		Container: widget.NewContainer(
			widget.ContainerOpts.Layout(layout),
		),
		outer: outer}
}

func (c *Container[T]) SetBackgroundImage(img *image.NineSlice) T {
	if c.Container.IsValidated() {
		c.Container.SetBackgroundImage(img)
	} else {
		widget.ContainerOpts.BackgroundImage(img)(c.Container)
	}
	return c.outer
}

func (c *Container[T]) SetMinSize(width, height int) T {
	widget.WidgetOpts.MinSize(width, height)(c.Container.GetWidget())
	return c.outer
}
