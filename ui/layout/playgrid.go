package layout

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

func NewPlayGrid(th *theme.Theme) *widget.Container {
	return widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(4),
				widget.GridLayoutOpts.Stretch(
					[]bool{true, true, true, true},
					[]bool{true, true, true, true},
				),
				widget.GridLayoutOpts.Spacing(th.PanelTheme.Spacing, th.PanelTheme.Spacing),
			),
		),
	)
}
