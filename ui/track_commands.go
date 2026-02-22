package ui

import (
	"fmt"
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"
)

type Mode uint8

const (
	ModeStopped Mode = iota
	ModePlaying
	ModeRecording
)

type Armed uint8

const (
	ArmNone Armed = iota
	ArmPlay
	ArmRecord
	ArmStop
)

type TrackCommands struct {
	*container.Row
	Id int

	playLed, recordLed *widgets.Icon
	playIcon           *widgets.Icon

	mode  Mode
	armed Armed
}

func NewTrackCommands(id int, sequencer *seq.Sequencer) *TrackCommands {
	tc := &TrackCommands{Id: id, mode: ModeStopped, armed: ArmNone}

	playButton := widgets.NewButton(func() {
		if tc.mode != ModeStopped || tc.armed == ArmPlay || tc.armed == ArmRecord {
			sequencer.TryCommand(seq.Command{Id: seq.CmdStop, TrackId: &id})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPlay, TrackId: &id})
		}
	})
	tc.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	tc.playLed = widgets.NewIcon(theme.IconLed, theme.IconSizeSmall)
	playButton.AddChild(tc.playIcon, tc.playLed)

	recordButton := widgets.NewButton(func() {
		if tc.mode == ModeRecording || tc.armed == ArmRecord {
			sequencer.TryCommand(seq.Command{Id: seq.CmdStopRecord, TrackId: &id})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdRecord, TrackId: &id})
		}
	})
	recordIcon := widgets.NewIcon(theme.IconRecord, theme.IconSizeMedium)
	tc.recordLed = widgets.NewIcon(theme.IconLed, theme.IconSizeSmall)
	recordButton.AddChild(recordIcon, tc.recordLed)

	clearButton := widgets.NewButton(func() {
		// sequencer.TryCommand(seq.Command{Id: seq.CmdClear, TrackId: &id}) // TODO
	})
	clearIcon := widgets.NewIcon(theme.IconClear, theme.IconSizeMedium)
	clearButton.AddChild(clearIcon)

	deleteButton := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdRemoveTrack, TrackId: &id})
	})
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

	tc.Row = container.NewRow().
		SetSpacing(theme.Current.PanelTheme.Spacing).
		SetPadding(theme.Current.PanelTheme.Padding)
	tc.Row.AddChild(widgets.NewLabel(fmt.Sprintf(" %02d", id+1)))
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

	if e.TrackId == nil || *e.TrackId != tc.Id {
		return
	}

	switch e.Id {
	case seq.EvPlaying:
		tc.mode = ModePlaying
		tc.armed = ArmNone
	case seq.EvRecording:
		tc.mode = ModeRecording
		tc.armed = ArmNone
	case seq.EvStopped:
		tc.mode = ModeStopped
		tc.armed = ArmNone
	case seq.EvArmedPlaying:
		tc.armed = ArmPlay
	case seq.EvArmedRecording:
		tc.armed = ArmRecord
	case seq.EvArmedStopped:
		tc.armed = ArmStop
	}

	tc.applyVisualState()
}

func (tc *TrackCommands) applyVisualState() {
	playing := tc.mode == ModePlaying || tc.mode == ModeRecording
	recording := tc.mode == ModeRecording

	if playing || tc.armed == ArmStop {
		tc.playIcon.SetIcon(theme.IconStop)
	} else {
		tc.playIcon.SetIcon(theme.IconPlay)
	}

	tc.playLed.SetPulseColor(theme.IconColorDefault)
	switch {
	case tc.armed == ArmStop:
		tc.playLed.SetColor(theme.IconColorArmed)
	case playing:
		tc.playLed.SetColor(theme.IconColorOn)
	case tc.armed == ArmPlay || tc.armed == ArmRecord:
		tc.playLed.SetColor(theme.IconColorArmed)
	default:
		tc.playLed.SetColor(theme.IconColorIdle)
		tc.playLed.SetPulseColor(theme.IconColorNone)
	}

	tc.recordLed.SetPulseColor(theme.IconColorDefault)
	switch {
	case tc.armed == ArmRecord:
		tc.recordLed.SetColor(theme.IconColorArmed)
	case tc.armed == ArmStop && recording:
		tc.recordLed.SetColor(theme.IconColorArmed)
	case recording:
		tc.recordLed.SetColor(theme.IconColorOn)
	default:
		tc.recordLed.SetColor(theme.IconColorIdle)
		tc.recordLed.SetPulseColor(theme.IconColorNone)
	}
}
