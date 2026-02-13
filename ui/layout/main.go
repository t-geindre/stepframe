package layout

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

func NewMain(th *theme.Theme, ops ...widget.ContainerOpt) *widget.Container {
	return widget.NewContainer(
		append([]widget.ContainerOpt{
			widget.ContainerOpts.Layout(widget.NewStackedLayout(
				widget.StackedLayoutOpts.Padding(th.PanelTheme.Padding),
			)),
			widget.ContainerOpts.BackgroundImage(th.PanelTheme.BackgroundImage),
		}, ops...)...,
	)
}
