package layout

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

func NewMain(th *theme.Theme) *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewStackedLayout(
			widget.StackedLayoutOpts.Padding(th.PanelTheme.Padding),
		)),
		widget.ContainerOpts.BackgroundImage(th.PanelTheme.BackgroundImage),
	)
}
