package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
)

const (
	BarsCount = 10
)

type TrackBars struct {
	*container.Grid
	state *TrackState
	bars  []*container.Row
	index int
	beats int64
}

func NewTrackBars(state *TrackState) *TrackBars {
	t := &TrackBars{
		Grid: container.NewGrid().
			SetTheme(theme.Current.BarContainerPanelTheme).
			SetColumns(BarsCount).
			SetDefaultStretch(true, true),
		bars:  make([]*container.Row, 0),
		state: state,
	}

	for i := 0; i < BarsCount; i++ {
		t.AddBar()
	}
	t.ApplyBarVisual()

	return t
}

func (t *TrackBars) HandleEvent(e seq.Event) {
	switch e.Id {
	case seq.EvStopped, seq.EvReset:
		if e.TrackId != nil && *e.TrackId == t.state.id {
			t.beats = -1 // todo time signature
			t.index = 0
			t.ApplyBarVisual()
		}
	case seq.EvBeat:
		if t.state.mode == ModePlaying || t.state.mode == ModeRecording {
			t.beats++
			if t.beats > 0 && (t.beats)%4 == 0 { // todo time signature
				t.index++
				if t.index >= len(t.bars) {
					t.index = 0
				}
			}
			t.ApplyBarVisual()
			t.bars[t.index].Pulse()
		}
	}
}

func (t *TrackBars) AddBar() {
	bar := container.NewRow().SetTheme(theme.Current.BarPanelTheme)

	t.bars = append(t.bars, bar)
	t.Grid.AddChild(bar)
}

func (t *TrackBars) ApplyBarVisual() {
	for _, b := range t.bars {
		b.SetColorizedBackgroundOffsetImage(theme.ColorIdle)
	}

	t.bars[t.index].SetColorizedBackgroundOffsetImage(theme.ColorOn)
}
