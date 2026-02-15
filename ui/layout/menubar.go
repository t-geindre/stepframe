package layout

import "github.com/ebitenui/ebitenui/widget"

func NewMenubar() *widget.Container {
	c := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{false, true, false}, []bool{true}),
		)),
	)

	left := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Stretch([]bool{false}, []bool{true}),
		)),
	)
	c.AddChild(left)

	return c
}
