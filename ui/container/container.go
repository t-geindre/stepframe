package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type Container[T widget.Containerer] struct {
	*widget.Container
	*Gradient[T]
	*OffsetBackground[T]
	*MouseState[T]
	outer T
}

func NewContainer[T widget.Containerer](layout widget.Layouter, outer T) *Container[T] {
	container := widget.NewContainer(widget.ContainerOpts.Layout(layout))
	gradient := NewGradient(container, outer)
	background := NewOffsetBackground(container, outer)
	mouse := NewMouseState(background, outer)

	return &Container[T]{
		Container:        container,
		Gradient:         gradient,
		OffsetBackground: background,
		MouseState:       mouse,
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
	if !c.GetWidget().IsVisible() {
		return
	}

	c.Gradient.Render(screen)
	c.OffsetBackground.Render(screen)
	c.Container.Render(screen)
}

func (c *Container[T]) SetTheme(th *theme.PanelTheme) T {
	c.OffsetBackground.SetBackgroundOffsetImage(th.IdleBackgroundImage)
	c.OffsetBackground.SetColorizedBackgroundOffsetImage(th.IdleBackgroundColorize)
	c.OffsetBackground.SetPulseBackgroundOffsetImage(th.PulseBackgroundImage)
	c.OffsetBackground.SetPulseColorizedBackgroundOffsetImage(th.PulseBackgroundColorize)
	c.OffsetBackground.SetPulseDuration(th.PulseDuration)
	c.Gradient.SetVerticalGradientBackground(th.GradientTop, th.GradientBottom)
	c.MouseState.SetIdleBackgroundOffsetImage(th.IdleBackgroundImage)
	c.MouseState.SetIdleBackgroundOffsetColorize(th.IdleBackgroundColorize)
	c.MouseState.SetHoverBackgroundOffsetImage(th.HoverBackgroundImage)
	c.MouseState.SetHoverBackgroundOffsetColorize(th.HoverBackgroundColorize)
	c.MouseState.SetPressedBackgroundOffsetImage(th.PressedBackgroundImage)
	c.MouseState.SetPressedBackgroundOffsetColorize(th.PressedBackgroundColorize)
	c.MouseState.SetHoverCursor(th.Cursor)
	return c.outer
}
