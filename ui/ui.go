package ui

import (
	"stepframe/clock"
	"stepframe/seq"
	"stepframe/ui/layout"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type Ui struct {
	seq   *seq.Sequencer
	grid  *Grid
	main  *widget.Container
	clock clock.Clock
}

func New(clk clock.Clock, sqr *seq.Sequencer) *ebitenui.UI {
	var ui *Ui
	th := theme.NewDefaultTheme()
	ui = &Ui{
		seq:   sqr,
		clock: clk,
		grid:  NewGrid(th, sqr),
		main: layout.NewMain(th, widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.OnUpdate(func(w widget.HasWidget) {
				ui.Update()
			}),
		)),
	}

	// TODO TESTS REMOVE ME
	// Add some tracks
	id := seq.TrackId(0)
	for _, track := range []*seq.Track{
		getBillieJeanBassTrack(),
		getBillieJeanLeadTrack(),
		getBillieJeanLeadTrackWithRatchet(clk),
		getBillieJeanLeadTrackWithRatchetDouble(clk),
	} {
		track.SetId(id)
		id++

		sqr.Commands() <- seq.Command{
			Id:    seq.CmdAdd,
			Track: track,
		}
	}
	// TODO END ---

	// default view
	ui.main.AddChild(ui.grid)

	return &ebitenui.UI{
		Container:    ui.main,
		PrimaryTheme: th.Theme,
	}
}

func (u *Ui) Update() {
	u.drainSequencer()
}

func (u *Ui) drainSequencer() {
	for i := 0; i < 1024; i++ { // prevent starvation
		select {
		case e := <-u.seq.Events():
			u.handleEvent(e)
		default:
			return
		}
	}
}

func (u *Ui) handleEvent(e seq.Event) {
	u.grid.HandleEvent(e)
}
