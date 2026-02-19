package container

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Bar struct {
	*widget.Container
	Left, Center, Right *Row
	locked              bool
}

func NewBar() *Bar {
	m := &Bar{
		Container: widget.NewContainer(
			widget.ContainerOpts.Layout(widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Stretch([]bool{false, true, false}, []bool{false}),
				widget.GridLayoutOpts.Padding(theme.Current.PanelTheme.Padding),
				widget.GridLayoutOpts.Spacing(theme.Current.PanelTheme.Spacing, theme.Current.PanelTheme.Spacing),
			)),
		),
		Left:   NewHorizontalRow().WithForeground().WithPadding().WitSpacing(),
		Center: NewHorizontalRow().WithForeground().WithPadding().WitSpacing(),
		Right:  NewHorizontalRow().WithForeground().WithPadding().WitSpacing(),
	}

	m.AddChild(m.Left)
	m.AddChild(m.Center)
	m.AddChild(m.Right)

	m.locked = true

	return m
}

func (b *Bar) AddChild(children ...widget.PreferredSizeLocateableWidget) widget.RemoveChildFunc {
	if b.locked {
		panic("cannot add child to bar directly, add to Left, Center or Right instead")
	}

	for _, c := range children {
		w := c.GetWidget()
		w.LayoutData = widget.GridLayoutData{
			HorizontalPosition: widget.GridLayoutPositionCenter,
			VerticalPosition:   widget.GridLayoutPositionCenter,
		}
		w.MinWidth, w.MinHeight = 1, 1
	}

	return b.Container.AddChild(children...)
}
