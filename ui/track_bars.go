package ui

import (
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"
)

const (
	BarsCount = 10
)

type TrackBars struct {
	*container.Grid
	bars []*container.Row
	add  *widgets.Button

	barWidth  int
	barsCount int
}

func NewTrackBars() *TrackBars {
	t := &TrackBars{
		Grid: container.NewGrid().
			SetBackgroundOffsetImage(theme.Current.BarContainerPanelTheme.BackgroundImage).
			SetPadding(theme.Current.BarContainerPanelTheme.Padding).
			SetSpacing(theme.Current.BarContainerPanelTheme.Spacing, theme.Current.BarContainerPanelTheme.Spacing).
			SetColumns(BarsCount).
			SetDefaultStretch(true, true),
		bars:      make([]*container.Row, 0),
		barsCount: BarsCount,
	}

	for i := 0; i < BarsCount; i++ {
		t.AddBar()
	}

	return t
}

func (t *TrackBars) HandleEvent(e any) {
}

func (t *TrackBars) AddBar() {
	bar := container.NewRow().SetBackgroundOffsetImage(theme.Current.BarPanelTheme.BackgroundImage)
	t.bars = append(t.bars, bar)
	t.Grid.AddChild(bar)
}
