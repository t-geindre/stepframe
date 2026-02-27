package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Row struct {
	*Container[*Row]
	layout    *widget.RowLayout
	cPosition widget.RowLayoutPosition
	cStretch  bool
}

func NewRow() *Row {
	row := &Row{}
	row.layout = widget.NewRowLayout()
	row.Container = NewContainer[*Row](row.layout, row)
	row.cPosition = widget.RowLayoutPositionCenter

	return row
}

func (r *Row) SetDirection(d widget.Direction) *Row {
	widget.RowLayoutOpts.Direction(d)(r.layout)
	return r
}

func (r *Row) SetSpacing(spacing int) *Row {
	widget.RowLayoutOpts.Spacing(spacing)(r.layout)
	return r
}

func (r *Row) SetPadding(padding *widget.Insets) *Row {
	widget.RowLayoutOpts.Padding(padding)(r.layout)
	return r
}

func (r *Row) SetContentPosition(position widget.RowLayoutPosition) *Row {
	r.cPosition = position
	return r
}

func (r *Row) SetContentStretch(stretch bool) *Row {
	r.cStretch = stretch
	return r
}

func (r *Row) SetTheme(th *theme.PanelTheme) *Row {
	r.SetPadding(th.Padding).SetSpacing(th.Spacing)
	r.Container.SetTheme(th)
	return r
}

func (r *Row) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	for _, c := range children {
		c.GetWidget().LayoutData = widget.RowLayoutData{
			Position: r.cPosition,
			Stretch:  r.cStretch,
		}
	}

	return r.Container.AddChild(children...)
}
