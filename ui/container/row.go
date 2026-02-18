package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Row struct {
	*widget.Container
	layout *widget.RowLayout
}

func NewRow(d widget.Direction) *Row {
	l := widget.NewRowLayout(
		widget.RowLayoutOpts.Spacing(theme.Current.PanelTheme.Spacing),
		widget.RowLayoutOpts.Direction(d),
	)
	return &Row{
		Container: widget.NewContainer(
			widget.ContainerOpts.Layout(l),
		),
		layout: l,
	}
}

func (r *Row) Themed() *Row {
	widget.RowLayoutOpts.Spacing(theme.Current.PanelTheme.Spacing)(r.layout)
	widget.RowLayoutOpts.Padding(theme.Current.PanelTheme.Padding)(r.layout)

	return r
}

func (r *Row) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, child := range children {
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
		})(child.GetWidget())
	}

	return r.Container.AddChild(children...)
}
