package ui

import (
	"stepframe/seq"
	"stepframe/ui/layout"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

type Grid struct {
	*widget.Container
	seq  *seq.Sequencer
	btns []*widgets.Play
}

func NewGrid(th *theme.Theme, sqr *seq.Sequencer) *Grid {
	g := &Grid{
		Container: layout.NewPlayGrid(th),
		seq:       sqr,
		btns:      make([]*widgets.Play, 0),
	}

	for i := 0; i < 16; i++ {
		btn := widgets.NewPlay(th)
		btn.GetWidget().MouseButtonClickedEvent.AddHandler(func(args interface{}) {
			g.Click(btn)
		})
		g.Container.AddChild(btn)
		g.btns = append(g.btns, btn)
	}

	return g
}

func (g *Grid) Click(btn *widgets.Play) {
	switch btn.GetState() {
	case widgets.PlayStateStopped:
		g.Play(btn.Id())
	case widgets.PlayStatePlaying:
		g.Stop(btn.Id())
	case widgets.PlayStateArmed:
		g.Stop(btn.Id())
	}
}

func (g *Grid) Play(id seq.TrackId) {
	g.seq.Commands() <- seq.Command{
		Id:      seq.CmdPlay,
		TrackId: id,
		At:      seq.CmdAtNextBar,
	}
}

func (g *Grid) Stop(id seq.TrackId) {
	g.seq.Commands() <- seq.Command{
		Id:      seq.CmdStop,
		TrackId: id,
		At:      seq.CmdAtNextBar,
	}
}

func (g *Grid) HandleEvent(e seq.Event) {
	switch e.Type {
	case seq.EvTrackArmed:
		g.SetBtnState(e.TrackId, widgets.PlayStateArmed)
	case seq.EvTrackPlay:
		g.SetBtnState(e.TrackId, widgets.PlayStatePlaying)
	case seq.EvTrackStop:
		g.SetBtnState(e.TrackId, widgets.PlayStateStopped)
	case seq.EvBeat:
		for _, b := range g.btns {
			b.Pulse()
		}
	case seq.EvTrackAdded:
		g.AddTrack(e.TrackId)
	default:
		// ignore
	}
}

func (g *Grid) SetBtnState(id seq.TrackId, state widgets.PlayState) {
	for _, b := range g.btns {
		if b.Id() == id {
			b.SetState(state)
			return
		}
	}
}

func (g *Grid) AddTrack(id seq.TrackId) {
	for _, b := range g.btns {
		if b.GetState() == widgets.PlayStateNone {
			b.SetId(id)
			b.SetState(widgets.PlayStateStopped)
			return
		}
	}
}
