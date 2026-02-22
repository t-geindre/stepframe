package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/widget"
)

type Track struct {
	widget.Containerer

	id        int
	sequencer *seq.Sequencer

	commands *TrackCommands
	bars     *TrackBars
}

func NewTrack(id int, sequencer *seq.Sequencer) *Track {
	t := &Track{
		Containerer: container.NewRow().
			SetDirection(widget.DirectionVertical).
			SetSpacing(theme.Current.PanelTheme.Spacing).
			SetPadding(theme.Current.PanelTheme.Padding).
			SetBackgroundImage(theme.Current.PanelTheme.ForegroundImage).
			SetContentPosition(widget.RowLayoutPositionStart).
			SetContentStretch(true),
		id:        id,
		sequencer: sequencer,
		commands:  NewTrackCommands(id, sequencer),
		bars:      NewTrackBars(),
	}

	t.AddChild(t.commands, t.bars)

	return t
}

func (t *Track) HandleEvent(e seq.Event) {
	t.commands.HandleEvent(e)
	t.bars.HandleEvent(e)
}
