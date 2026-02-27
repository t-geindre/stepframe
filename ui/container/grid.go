package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Grid struct {
	*Container[*Grid]
	layout *widget.GridLayout
}

func NewGrid() *Grid {
	grid := &Grid{}
	grid.layout = widget.NewGridLayout(widget.GridLayoutOpts.Columns(1))
	grid.Container = NewContainer[*Grid](grid.layout, grid)

	return grid
}

func (g *Grid) SetColumns(columns int) *Grid {
	widget.GridLayoutOpts.Columns(columns)(g.layout)
	return g
}

func (g *Grid) SetSpacing(vertical, horizontal int) *Grid {
	widget.GridLayoutOpts.Spacing(vertical, horizontal)(g.layout)
	return g
}

func (g *Grid) SetPadding(padding *widget.Insets) *Grid {
	widget.GridLayoutOpts.Padding(padding)(g.layout)
	return g
}

func (g *Grid) SetStretch(columns, rows []bool) *Grid {
	widget.GridLayoutOpts.Stretch(columns, rows)(g.layout)
	return g
}

func (g *Grid) SetDefaultStretch(column, row bool) *Grid {
	widget.GridLayoutOpts.DefaultStretch(column, row)(g.layout)
	return g
}

func (g *Grid) SetTheme(th *theme.PanelTheme) *Grid {
	g.SetPadding(th.Padding).SetSpacing(th.Spacing, th.Spacing)
	g.Container.SetTheme(th)
	return g
}

func (g *Grid) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, c := range children {
		c.GetWidget().LayoutData = widget.GridLayoutData{
			HorizontalPosition: widget.GridLayoutPositionCenter,
			VerticalPosition:   widget.GridLayoutPositionCenter,
		}
	}

	return g.Container.AddChild(children...)
}
