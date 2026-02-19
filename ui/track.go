package ui

import (
	"fmt"
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
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

type Track struct {
	*container.Row
	Id int

	playLed, recordLed *widgets.Icon
	playIcon           *widgets.Icon

	mode  Mode
	armed Armed
}

func NewTrack(id int, sequencer *seq.Sequencer) *Track {
	t := &Track{Id: id, mode: ModeStopped, armed: ArmNone}

	playButton := widgets.NewButton(func() {
		if t.mode != ModeStopped || t.armed == ArmPlay {
			sequencer.TryCommand(seq.Command{Id: seq.CmdStop, TrackId: &id})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPlay, TrackId: &id})
		}
	})
	t.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	t.playLed = widgets.NewIcon(theme.IconLed, theme.IconSizeSmall)
	playButton.AddChild(t.playIcon, t.playLed)

	recordButton := widgets.NewButton(func() {
		if t.mode == ModeRecording || t.armed == ArmRecord {
			sequencer.TryCommand(seq.Command{Id: seq.CmdStopRecord, TrackId: &id})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdRecord, TrackId: &id})
		}
	})
	recordIcon := widgets.NewIcon(theme.IconRecord, theme.IconSizeMedium)
	t.recordLed = widgets.NewIcon(theme.IconLed, theme.IconSizeSmall)
	recordButton.AddChild(recordIcon, t.recordLed)

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

	optionsButton := widgets.NewButton(func() {
		// sequencer.TryCommand(seq.Command{Id: seq.CmdOpenTrackOptions, TrackId: &id}) // TODO
	})
	optionsIcon := widgets.NewIcon(theme.IconGear, theme.IconSizeMedium)
	optionsButton.AddChild(optionsIcon)

	t.Row = container.NewRow(widget.DirectionHorizontal)
	t.Row.AddChild(widgets.NewText(fmt.Sprintf(" %02d", id+1)))
	t.Row.AddChild(playButton, recordButton)
	t.Row.AddChild(clearButton, deleteButton, optionsButton)

	t.applyVisualState()
	return t
}

func (t *Track) HandleEvent(e seq.Event) {
	switch e.Id {
	case seq.EvPlaying:
		t.mode = ModePlaying
		t.armed = ArmNone
	case seq.EvRecording:
		t.mode = ModeRecording
		t.armed = ArmNone
	case seq.EvStopped:
		t.mode = ModeStopped
		t.armed = ArmNone
	case seq.EvArmedPlaying:
		t.armed = ArmPlay
	case seq.EvArmedRecording:
		t.armed = ArmRecord
	case seq.EvArmedStopped:
		t.armed = ArmStop
	}

	t.applyVisualState()
}

func (t *Track) applyVisualState() {
	playing := t.mode == ModePlaying || t.mode == ModeRecording
	recording := t.mode == ModeRecording

	if playing || t.armed == ArmStop {
		t.playIcon.SetIcon(theme.IconStop)
	} else {
		t.playIcon.SetIcon(theme.IconPlay)
	}

	switch {
	case t.armed == ArmStop:
		t.playLed.SetColor(theme.IconColorArmed)
	case playing:
		t.playLed.SetColor(theme.IconColorLedOn)
	case t.armed == ArmPlay || t.armed == ArmRecord:
		t.playLed.SetColor(theme.IconColorArmed)
	default:
		t.playLed.SetColor(theme.IconColorLedOff)
	}

	switch {
	case t.armed == ArmRecord:
		t.recordLed.SetColor(theme.IconColorArmed)
	case t.armed == ArmStop && recording:
		t.recordLed.SetColor(theme.IconColorArmed)
	case recording:
		t.recordLed.SetColor(theme.IconColorLedOn)
	default:
		t.recordLed.SetColor(theme.IconColorLedOff)
	}
}
