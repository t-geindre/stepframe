package ui

import (
	"stepframe/ui/layout"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui"
)

func New() *ebitenui.UI {
	th := theme.NewDefaultTheme()

	main := layout.NewMain(th)

	return &ebitenui.UI{
		Container:    main,
		PrimaryTheme: th.Theme,
	}
}
