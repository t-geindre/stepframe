package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"
)

type TrackCommands struct {
	*container.Row
	Id                 int
	playLed, recordLed *widgets.Icon
	playIcon           *widgets.Icon
	state              *TrackState
}

func NewTrackCommands(state *TrackState) *TrackCommands {
	tc := &TrackCommands{state: state}

	playButton := widgets.NewButton(state.TogglePlay)
	tc.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	tc.playLed = widgets.NewIcon(theme.IconLed, theme.IconSizeSmall)
	playButton.AddChild(tc.playIcon, tc.playLed)

	recordButton := widgets.NewButton(state.ToggleRecord)
	recordIcon := widgets.NewIcon(theme.IconRecord, theme.IconSizeMedium)
	tc.recordLed = widgets.NewIcon(theme.IconLed, theme.IconSizeSmall)
	recordButton.AddChild(recordIcon, tc.recordLed)

	clearButton := widgets.NewButton(func() {
		// sequencer.TryCommand(seq.Command{Id: seq.CmdClear, TrackId: &id}) // TODO
	})
	clearIcon := widgets.NewIcon(theme.IconClear, theme.IconSizeMedium)
	clearButton.AddChild(clearIcon)

	deleteButton := widgets.NewButton(state.Delete)
	deleteIcon := widgets.NewIcon(theme.IconDelete, theme.IconSizeMedium)
	deleteButton.AddChild(deleteIcon)

	optionsWin := widgets.NewWindow(300, 200).WithTitleBar(theme.IconGear, "TrackCommands Options")
	optionsButton := widgets.NewButton(func() {
		optionsWin.Open()
	})
	optionsWin.AttachedTo(optionsButton)
	optionsWin.AddChild(widgets.NewLabel("test"), widgets.NewIcon(theme.IconPlus, theme.IconSizeMedium))
	optionsIcon := widgets.NewIcon(theme.IconGear, theme.IconSizeMedium)
	optionsButton.AddChild(optionsIcon)

	tc.Row = container.NewRow().SetSpacing(theme.Current.PanelTheme.Spacing)
	tc.Row.AddChild(playButton, recordButton)
	tc.Row.AddChild(clearButton, deleteButton, optionsButton)

	tc.applyVisualState()
	return tc
}

func (tc *TrackCommands) HandleEvent(e seq.Event) {
	if e.Id == seq.EvBeat {
		tc.playLed.Pulse()
		tc.recordLed.Pulse()
		return
	}

	if e.TrackId != nil && *e.TrackId == tc.Id {
		tc.applyVisualState()
		return
	}
}

func (tc *TrackCommands) applyVisualState() {
	playing := tc.state.mode == ModePlaying || tc.state.mode == ModeRecording
	recording := tc.state.mode == ModeRecording

	if playing || tc.state.armed == ArmStop {
		tc.playIcon.SetIcon(theme.IconStop)
	} else {
		tc.playIcon.SetIcon(theme.IconPlay)
	}

	tc.playLed.SetPulseColor(theme.IconColorDefault)
	switch {
	case tc.state.armed == ArmStop:
		tc.playLed.SetColor(theme.IconColorArmed)
	case playing:
		tc.playLed.SetColor(theme.IconColorOn)
	case tc.state.armed == ArmPlay || tc.state.armed == ArmRecord:
		tc.playLed.SetColor(theme.IconColorArmed)
	default:
		tc.playLed.SetColor(theme.IconColorIdle)
		tc.playLed.SetPulseColor(theme.IconColorNone)
	}

	tc.recordLed.SetPulseColor(theme.IconColorDefault)
	switch {
	case tc.state.armed == ArmRecord:
		tc.recordLed.SetColor(theme.IconColorArmed)
	case tc.state.armed == ArmStop && recording:
		tc.recordLed.SetColor(theme.IconColorArmed)
	case recording:
		tc.recordLed.SetColor(theme.IconColorOn)
	default:
		tc.recordLed.SetColor(theme.IconColorIdle)
		tc.recordLed.SetPulseColor(theme.IconColorNone)
	}
}
