package widgets

import (
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Label struct {
	*widget.Text
}

func NewLabel(label string) *Label {
	return &Label{
		Text: widget.NewText(
			widget.TextOpts.TextLabel(label),
		),
	}
}

func (l *Label) SetFont(ft *text.Face) {
	widget.TextOpts.TextFace(ft)(l.Text)
}
