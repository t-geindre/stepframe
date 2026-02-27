package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Anchor struct {
	*Container[*Anchor]
	layout     *widget.AnchorLayout
	layoutData widget.AnchorLayoutData
}

func NewAnchor() *Anchor {
	a := &Anchor{}
	a.layout = widget.NewAnchorLayout()
	a.Container = NewContainer[*Anchor](a.layout, a)
	a.layoutData = widget.AnchorLayoutData{
		HorizontalPosition: widget.AnchorLayoutPositionCenter,
		VerticalPosition:   widget.AnchorLayoutPositionCenter,
	}

	return a
}

func (a *Anchor) SetPadding(padding *widget.Insets) *Anchor {
	widget.AnchorLayoutOpts.Padding(padding)(a.layout)
	return a
}

func (a *Anchor) SetContentVerticalPosition(verticalPosition widget.AnchorLayoutPosition) *Anchor {
	a.layoutData.VerticalPosition = verticalPosition
	return a
}

func (a *Anchor) SetContentHorizontalPosition(horizontalPosition widget.AnchorLayoutPosition) *Anchor {
	a.layoutData.HorizontalPosition = horizontalPosition
	return a
}

func (a *Anchor) SetContentPosition(horizontalPosition, verticalPosition widget.AnchorLayoutPosition) *Anchor {
	a.layoutData.HorizontalPosition = horizontalPosition
	a.layoutData.VerticalPosition = verticalPosition
	return a
}

func (a *Anchor) SetTheme(th *theme.PanelTheme) *Anchor {
	a.SetPadding(th.Padding)
	a.Container.SetTheme(th)
	return a
}

func (a *Anchor) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, c := range children {
		c.GetWidget().LayoutData = a.layoutData
	}

	return a.Container.AddChild(children...)
}
