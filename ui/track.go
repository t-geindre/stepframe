package ui

import (
	"fmt"
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

type Track struct {
	widget.Containerer
	state    *TrackState
	commands *TrackCommands
	bars     *TrackBars
}

func NewTrack(id int, sequencer *seq.Sequencer) *Track {
	state := NewTrackState(id, sequencer)

	left := widgets.NewLabel(fmt.Sprintf(" %02d", id+1))

	right := container.NewRow().
		SetDirection(widget.DirectionVertical).
		SetSpacing(theme.Current.PanelTheme.Spacing).
		SetContentStretch(true)

	commands := NewTrackCommands(state)
	bars := NewTrackBars()
	right.AddChild(commands, bars)

	cont := container.NewGrid().
		SetColumns(2).
		SetStretch([]bool{false, true}, []bool{true}).
		SetSpacing(theme.Current.PanelTheme.Spacing, theme.Current.PanelTheme.Spacing).
		SetPadding(theme.Current.PanelTheme.Padding).
		SetBackgroundImage(theme.Current.PanelTheme.ForegroundImage)

	cont.AddChild(left, right)

	t := &Track{
		Containerer: cont,
		state:       state,
		commands:    commands,
		bars:        bars,
	}

	return t
}

func (t *Track) HandleEvent(e seq.Event) {
	t.state.HandleEvent(e)
	t.commands.HandleEvent(e)
	t.bars.HandleEvent(e)
}
