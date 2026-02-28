package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
)

const (
	BarsCount = 10
)

type Bar struct {
	*container.Row
	Active bool
}

type TrackBars struct {
	*container.Grid
	state *TrackState
	bars  []*Bar
	index int
	beats int64
}

func NewTrackBars(state *TrackState) *TrackBars {
	t := &TrackBars{
		Grid: container.NewGrid().
			SetTheme(theme.Current.BarContainerPanelTheme).
			SetColumns(BarsCount).
			SetDefaultStretch(true, true),
		bars:  make([]*Bar, 0),
		state: state,
		beats: -1, // todo time signature
	}

	t.AddBars()
	t.Stop()

	return t
}

func (t *TrackBars) HandleEvent(e seq.Event) {
	switch e.Id {
	case seq.EvStopped, seq.EvReset:
		if e.TrackId != nil && *e.TrackId == t.state.id {
			t.Stop()
		}
	case seq.EvBeat:
		if t.state.mode == ModePlaying || t.state.mode == ModeRecording {
			t.beats++
			if t.beats > 0 && (t.beats)%4 == 0 { // todo time signature
				t.SetIndex(t.index + 1)
			}
			t.bars[t.index].Pulse()
		}
	}
}

func (t *TrackBars) AddBars() {
	for i := 0; i < BarsCount; i++ {
		bar := &Bar{
			Row: container.NewRow().SetOnClick(func() {
				t.DisableBar(i)
			}),
			Active: len(t.bars) == 0, // first bar is active by default
		}

		t.bars = append(t.bars, bar)
		t.Grid.AddChild(bar)
	}
}

func (t *TrackBars) DisableBar(index int) {
	defer t.ApplyBarVisual()

	if !t.bars[index].Active {
		t.bars[index].Active = true
		return
	}

	// at least one bar should be active
	oneActive := false
	for i, b := range t.bars {
		if b.Active && i != index {
			oneActive = true
			break
		}
	}

	if !oneActive {
		return
	}

	t.bars[index].Active = false
}

func (t *TrackBars) ApplyBarVisual() {
	for i, b := range t.bars {
		if i == t.index {
			b.SetTheme(theme.Current.BarActivePanelTheme)
			continue
		}

		if b.Active {
			b.SetTheme(theme.Current.BarOnPanelTheme)
			continue
		}

		b.SetTheme(theme.Current.BarOffPanelTheme)
	}
}

func (t *TrackBars) SetIndex(index int) {
	i := 0
	for i < 1000 {
		if index >= len(t.bars) {
			index = 0
		}
		if t.bars[index].Active {
			break
		}
		index++
		i++
	}
	t.index = index
	t.ApplyBarVisual()
}

func (t *TrackBars) Stop() {
	t.beats = -1 // do not skip first bar
	t.SetIndex(0)
}
