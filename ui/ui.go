package ui

import (
	"fmt"
	"stepframe/clock"
	"stepframe/midi"
	"stepframe/seq"
	"stepframe/ui/layout"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

type Ui struct {
	seq      *seq.Sequencer
	grid     *Grid
	main     *widget.Container
	clock    clock.Clock
	sender   *midi.Sender
	receiver *midi.Receiver
}

func New(clk clock.Clock, sqr *seq.Sequencer, sd *midi.Sender, rc *midi.Receiver) *ebitenui.UI {
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
	// Display in/out ports
	fmt.Println("MIDI OUT PORTS:")
	for _, p := range midi.AllOutPorts() {
		fmt.Printf(" - PORT[%d]: %s\n", p.Id, p.Name)
	}
	fmt.Println("MIDI IN PORTS:")
	for _, p := range midi.AllInPorts() {
		fmt.Printf(" _ PORT[%d]: %s\n", p.Id, p.Name)
	}
	// open ports
	sd.Commands() <- midi.Command{Id: midi.CmdOpenPort, Port: 1} // out
	rc.Commands() <- midi.Command{Id: midi.CmdOpenPort, Port: 2} // in
	// forward
	rc.Commands() <- midi.Command{Id: midi.CmdForward, Port: 2, PortOut: 1}
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
