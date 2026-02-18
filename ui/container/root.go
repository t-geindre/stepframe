package container

import (
	"github.com/ebitenui/ebitenui/widget"
)

type Root struct {
	*widget.Container
}

func NewRoot() *Root {
	return &Root{
		Container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true}),
			)),
		),
	}
}

func (r *Root) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, c := range children {
		c.GetWidget().LayoutData = widget.GridLayoutData{
			HorizontalPosition: widget.GridLayoutPositionCenter,
			VerticalPosition:   widget.GridLayoutPositionCenter,
		}
	}

	return r.Container.AddChild(children...)
}
