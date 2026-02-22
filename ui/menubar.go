package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"
	"strconv"

	"github.com/ebitenui/ebitenui/widget"
)

type TopBar struct {
	*container.Bar
	beatLed *widgets.Icon

	playIcon  *widgets.Icon
	playLabel *widget.Text
	playing   bool

	beatPerBar      int
	beatPerBarLabel *widget.Text
}

func NewTopBar(sequencer *seq.Sequencer) *TopBar {
	t := &TopBar{
		Bar:        container.NewBar(),
		beatPerBar: 4,
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

	// Beat per bar
	bpbPlusBtn := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdAddBeatPerBar})
	})
	bpbPlusBtn.AddChild(widgets.NewIcon(theme.IconPlus, theme.IconSizeSmall))
	bpbMinusBtn := widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdSubBeatPerBar})
	})
	bpbMinusBtn.AddChild(widgets.NewIcon(theme.IconMinus, theme.IconSizeSmall))
	t.beatPerBarLabel = widgets.NewLabel(strconv.Itoa(t.beatPerBar))

	// Place widgets
	t.Bar.Left.AddChild(addBtn)
	t.Bar.Center.AddChild(t.beatLed, playButton, stopButton)
	t.Bar.Right.AddChild(widgets.NewLabel("BPM 120"))
	t.Bar.Right.AddChild(
		widgets.NewLabel("BPB"), bpbMinusBtn, t.beatPerBarLabel, bpbPlusBtn,
	)

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
	case seq.EvAddBeatPerBar:
		t.SetBeatPerBar(t.beatPerBar + 1)
	case seq.EvSubBeatPerBar:
		t.SetBeatPerBar(t.beatPerBar - 1)
	}
}
func (t *TopBar) SetBeatPerBar(bpb int) {
	if bpb < 1 {
		return
	}

	t.beatPerBar = bpb
	t.beatPerBarLabel.Label = strconv.Itoa(bpb)
	t.Bar.RequestRelayout()
}
