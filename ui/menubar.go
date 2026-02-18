package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"
)

type TopBar struct {
	*container.Bar
	beatLed  *widgets.Icon
	playIcon *widgets.Icon
	playing  bool
}

func NewTopBar(sequencer *seq.Sequencer) *TopBar {
	t := &TopBar{
		Bar: container.NewBar(),
	}

	// BPM LED
	t.beatLed = widgets.NewIcon(theme.IconLed, theme.IconSizeMedium)
	t.beatLed.SetColor(theme.IconColorLedOff)
	t.beatLed.SetPulseColor(theme.IconColorLedOn)

	// Play button
	t.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	playButton := widgets.NewButton(func() {
		if t.playing {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPause})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPlay})
		}
	})
	playButton.AddChild(t.playIcon)

	// Stop button
	stopButton := widgets.NewIconButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdStop})
	}, theme.IconStop, theme.IconSizeMedium)

	// Place widgets
	t.Bar.Left.AddChild(t.beatLed)
	t.Bar.Left.AddChild(playButton)
	t.Bar.Left.AddChild(stopButton)

	t.Bar.Right.AddChild(widgets.NewText("BPM 120"))
	t.Bar.Right.AddChild(widgets.NewText("BPB 4"))

	return t
}

func (t *TopBar) HandleEvent(event seq.Event) {
	if event.TrackId != nil {
		return // Track command, ignore
	}
	switch event.Id {
	case seq.EvBeat:
		t.beatLed.Pulse()
	case seq.EvPaused, seq.EvStopped:
		t.playing = false
		t.playIcon.SetIcon(theme.IconPlay)
	case seq.EvPlaying:
		t.playing = true
		t.playIcon.SetIcon(theme.IconPause)
	}
}
