package widgets

import (
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

var currentUi *ebitenui.UI

func NewUi(root widget.Containerer) *ebitenui.UI {
	if theme.Current == nil {
		panic("theme.Current must be set before creating UI")
	}

	ui := &ebitenui.UI{
		Container:    root,
		PrimaryTheme: theme.Current.Theme,
	}

	currentUi = ui

	return ui
}
