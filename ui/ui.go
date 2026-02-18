package ui

import (
	"stepframe/clock"
	"stepframe/midi"
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/rs/zerolog"
)

type Ui struct {
	// Core
	clock     *clock.Clock
	sequencer *seq.Sequencer
	sender    *midi.Sender
	receiver  *midi.Receiver
	// Logging
	logger zerolog.Logger
	// UI
	root   widget.Containerer
	menu   *TopBar
	tracks *Tracks
}

func New(
	clock clock.Clock,
	sequencer *seq.Sequencer,
	sender *midi.Sender,
	receiver *midi.Receiver,
	logger zerolog.Logger,
) *ebitenui.UI {
	theme.SetDefaultTheme()

	ui := &Ui{
		root:      container.NewRoot(),
		clock:     &clock,
		sequencer: sequencer,
		sender:    sender,
		receiver:  receiver,
		logger:    logger.With().Str("component", "ui").Logger(),
		menu:      NewTopBar(sequencer),
		tracks:    NewTracks(logger, sequencer),
	}

	ui.root.GetWidget().OnUpdate = func(w widget.HasWidget) {
		ui.drainSequencer()
	}

	ui.root.AddChild(ui.menu)
	ui.root.AddChild(ui.tracks)

	return &ebitenui.UI{
		Container:    ui.root,
		PrimaryTheme: theme.Current.Theme,
	}
}

func (u *Ui) drainSequencer() {
	for i := 0; i < 1024; i++ { // prevent starvation
		select {
		case e := <-u.sequencer.Events():
			u.handleEvent(e)
		default:
			return
		}
	}
}

func (u *Ui) handleEvent(e seq.Event) {
	u.menu.HandleEvent(e)
	u.tracks.HandleEvent(e)
}
