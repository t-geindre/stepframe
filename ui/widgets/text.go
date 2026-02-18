package widgets

import (
	"github.com/ebitenui/ebitenui/widget"
)

func NewText(label string) *widget.Text {
	return widget.NewText(
		widget.TextOpts.TextLabel(label),
	)
}
