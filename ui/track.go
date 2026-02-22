package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"

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
		Containerer: container.NewVerticalRow().WitSpacing().StretchContent(),
		id:          id,
		sequencer:   sequencer,
		commands:    NewTrackCommands(id, sequencer),
		bars:        NewTrackBars(),
	}

	t.AddChild(t.commands, t.bars)

	return t
}

func (t *Track) HandleEvent(e seq.Event) {
	t.commands.HandleEvent(e)
	t.bars.HandleEvent(e)
}
