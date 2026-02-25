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
	beatLeds   []*widgets.Led
	ledCurrent int
	firstBeat  bool

	playIcon  *widgets.Icon
	playLabel *widgets.Label
	playing   bool
}

func NewTopBar(sequencer *seq.Sequencer) *TopBar {
	t := &TopBar{Containerer: container.NewGrid().
		SetColumns(3).
		SetStretch([]bool{false, true, false}, []bool{false}).
		SetSpacing(theme.Current.PanelTheme.Spacing, theme.Current.PanelTheme.Spacing).
		SetPadding(theme.Current.PanelTheme.Padding),
	}

	// BPM LED todo BPB
	t.beatLeds = make([]*widgets.Led, 4)
	for i := 0; i < 4; i++ {
		t.beatLeds[i] = widgets.NewLed(true, theme.ColorIdle)
		t.beatLeds[i].SetPulseColor(theme.ColorOn)
	}

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
	playButton.AddChild(t.playIcon, t.playLabel)

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
	for _, led := range t.beatLeds {
		center.AddChild(led)
	}
	right.AddChild(widgets.NewLabel("BPM 120"), widgets.NewLabel("BPB 4"))
	t.AddChild(left, center, right)

	// Initial state
	t.SetStopped(false)

	return t
}

func (t *TopBar) HandleEvent(event seq.Event) {
	if event.TrackId != nil {
		return // TrackCommands command, ignore
	}
	switch event.Id {
	case seq.EvBeat:
		t.OnBeat()
	case seq.EvPaused:
		t.SetStopped(true)
	case seq.EvStopped:
		t.SetStopped(false)
	case seq.EvPlaying:
		t.SetPlaying()
	}
}

func (t *TopBar) getBox() widget.Containerer {
	return container.NewRow().
		SetBackgroundImage(theme.Current.PanelTheme.BackgroundImage).
		SetPadding(theme.Current.PanelTheme.Padding).
		SetSpacing(theme.Current.PanelTheme.Spacing).
		SetContentStretch(true)
}

func (t *TopBar) SetStopped(pause bool) {
	if !pause {
		for _, led := range t.beatLeds {
			led.SetColor(theme.ColorIdle)
		}
		t.ledCurrent = -1
	}

	t.playing = false
	t.playIcon.SetIcon(theme.IconPlay)
	t.playLabel.Label = "Play"
	t.Containerer.RequestRelayout()
}

func (t *TopBar) SetPlaying() {
	t.playing = true
	t.playIcon.SetIcon(theme.IconPause)
	t.playLabel.Label = "Pause"
	t.Containerer.RequestRelayout()
}

func (t *TopBar) OnBeat() {
	if t.ledCurrent >= 0 {
		t.beatLeds[t.ledCurrent].SetColor(theme.ColorIdle)
	}

	t.ledCurrent++
	if t.ledCurrent >= len(t.beatLeds) {
		t.ledCurrent = 0
	}

	t.beatLeds[t.ledCurrent].SetColor(theme.ColorOn)
	t.beatLeds[t.ledCurrent].Pulse()

}
