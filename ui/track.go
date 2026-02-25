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

	idLabel := widgets.NewLabel(fmt.Sprintf(" %02d", id+1))
	idLabel.SetFont(theme.Current.TrackTheme.Id.Font)
	
	commands := NewTrackCommands(state)
	bars := NewTrackBars()

	cont := container.NewGrid().
		SetColumns(3).
		SetStretch([]bool{false, false, true}, []bool{true}).
		SetSpacing(theme.Current.PanelTheme.Spacing, theme.Current.PanelTheme.Spacing).
		SetPadding(theme.Current.PanelTheme.Padding).
		SetBackgroundImage(theme.Current.PanelTheme.BackgroundImage)

	cont.AddChild(idLabel, commands, bars)

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
