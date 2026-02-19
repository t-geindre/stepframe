package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Row struct {
	*widget.Container
	layout  *widget.RowLayout
	align   widget.RowLayoutPosition
	stretch bool
}

func newRow(d widget.Direction) *Row {
	l := widget.NewRowLayout(
		widget.RowLayoutOpts.Spacing(theme.Current.PanelTheme.Spacing),
		widget.RowLayoutOpts.Direction(d),
	)
	return &Row{
		Container: widget.NewContainer(
			widget.ContainerOpts.Layout(l),
		),
		layout: l,
		align:  widget.RowLayoutPositionCenter,
	}
}

func NewHorizontalRow() *Row {
	return newRow(widget.DirectionHorizontal)
}

func NewVerticalRow() *Row {
	return newRow(widget.DirectionVertical)
}

func (r *Row) AlignContent(position widget.RowLayoutPosition) *Row {
	r.align = position
	return r
}

func (r *Row) StretchContent() *Row {
	r.stretch = true
	return r
}

func (r *Row) WitSpacing() *Row {
	widget.RowLayoutOpts.Spacing(theme.Current.PanelTheme.Spacing)(r.layout)
	return r
}

func (r *Row) WithPadding() *Row {
	widget.RowLayoutOpts.Padding(theme.Current.PanelTheme.Padding)(r.layout)
	return r
}

func (r *Row) WithBackground() *Row {
	widget.ContainerOpts.BackgroundImage(theme.Current.PanelTheme.BackgroundImage)(r.Container)
	return r
}

func (r *Row) WithForeground() *Row {
	widget.ContainerOpts.BackgroundImage(theme.Current.PanelTheme.ForegroundImage)(r.Container)
	return r
}

func (r *Row) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, child := range children {
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: r.align,
			Stretch:  r.stretch,
		})(child.GetWidget())
	}

	return r.Container.AddChild(children...)
}
