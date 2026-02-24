package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

type TopBar struct {
	widget.Containerer
	beatLed *widgets.Led

	playIcon  *widgets.Icon
	playLabel *widget.Text
	playing   bool
}

func NewTopBar(sequencer *seq.Sequencer) *TopBar {
	t := &TopBar{Containerer: container.NewGrid().
		SetColumns(3).
		SetStretch([]bool{false, true, false}, []bool{false}).
		SetSpacing(theme.Current.PanelTheme.Spacing, theme.Current.PanelTheme.Spacing).
		SetPadding(theme.Current.PanelTheme.Padding),
	}

	// BPM LED
	t.beatLed = widgets.NewLed(true, theme.IconColorIdle)
	t.beatLed.SetPulseColor(theme.IconColorDefault)

	// Add track
	addBtn := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdNewTrack})
	})
	addBtn.AddChild(
		widgets.NewIcon(theme.IconPlus, theme.IconSizeMedium),
		widgets.NewLabel("Track"),
	)

	// Play button
	t.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	t.playLabel = widgets.NewLabel("Play")
	playButton := widgets.NewButton(func() {
		if t.playing {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPause})
		} else {
			sequencer.TryCommand(seq.Command{Id: seq.CmdPlay})
		}
	})
	playButton.AddChild(t.playIcon, t.playLabel, t.beatLed)

	// Stop button
	stopIcon := widgets.NewIcon(theme.IconStop, theme.IconSizeMedium)
	stopLabel := widgets.NewLabel("Stop")
	stopButton := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdStop})
	})
	stopButton.AddChild(stopIcon, stopLabel)

	// Place widgets
	left, center, right := t.getBox(), t.getBox(), t.getBox()
	left.AddChild(addBtn)
	center.AddChild(playButton, stopButton)
	right.AddChild(widgets.NewLabel("BPM 120"), widgets.NewLabel("BPB 4"))
	t.AddChild(left, center, right)

	return t
}

func (t *TopBar) HandleEvent(event seq.Event) {
	if event.TrackId != nil {
		return // TrackCommands command, ignore
	}
	switch event.Id {
	case seq.EvBeat:
		t.beatLed.Pulse()
	case seq.EvPaused, seq.EvStopped:
		if event.Id == seq.EvPaused {
			t.beatLed.SetColor(theme.IconColorArmed)
		} else {
			t.beatLed.SetColor(theme.IconColorIdle)
		}
		t.playing = false
		t.playIcon.SetIcon(theme.IconPlay)
		t.playLabel.Label = "Play"
		t.Containerer.RequestRelayout()
	case seq.EvPlaying:
		t.beatLed.SetColor(theme.IconColorOn)
		t.playing = true
		t.playIcon.SetIcon(theme.IconPause)
		t.playLabel.Label = "Pause"
		t.Containerer.RequestRelayout()
	}
}

func (t *TopBar) getBox() widget.Containerer {
	return container.NewRow().
		SetBackgroundImage(theme.Current.PanelTheme.ForegroundImage).
		SetPadding(theme.Current.PanelTheme.Padding).
		SetSpacing(theme.Current.PanelTheme.Spacing).
		SetContentStretch(true)
}
