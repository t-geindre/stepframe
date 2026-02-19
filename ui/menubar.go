package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

type TopBar struct {
	*container.Bar
	beatLed   *widgets.Icon
	playIcon  *widgets.Icon
	playLabel *widget.Text
	playing   bool
}

func NewTopBar(sequencer *seq.Sequencer) *TopBar {
	t := &TopBar{
		Bar: container.NewBar(),
	}

	// BPM LED
	t.beatLed = widgets.NewIcon(theme.IconLed, theme.IconSizeMedium)
	t.beatLed.SetColor(theme.IconColorIdle)
	t.beatLed.SetPulseColor(theme.IconColorOn)

	// Add track
	addBtn := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdNewTrack})
	})
	addBtn.AddChild(
		widgets.NewIcon(theme.IconPlus, theme.IconSizeMedium),
		widgets.NewText("Track"),
	)

	// Play button
	t.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	t.playLabel = widgets.NewText("Play")
	playButton := widgets.NewButton(func() {
		if t.playing {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPause})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPlay})
		}
	})
	playButton.AddChild(t.playIcon, t.playLabel)

	// Stop button
	stopIcon := widgets.NewIcon(theme.IconStop, theme.IconSizeMedium)
	stopLabel := widgets.NewText("Stop")
	stopButton := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdStop})
	})
	stopButton.AddChild(stopIcon, stopLabel)

	// Place widgets
	t.Bar.Left.AddChild(addBtn)

	t.Bar.Center.AddChild(t.beatLed)
	t.Bar.Center.AddChild(playButton)
	t.Bar.Center.AddChild(stopButton)

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
		t.playLabel.Label = "Play"
		t.Bar.RequestRelayout()
	case seq.EvPlaying:
		t.playing = true
		t.playIcon.SetIcon(theme.IconPause)
		t.playLabel.Label = "Pause"
		t.Bar.RequestRelayout()
	}
}
