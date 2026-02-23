package ui

import (
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

// TODO: move to theme
const (
	BarHeight        = 40
	BarsDefaultCount = 4
)

var BarsBreakpoints = []struct{ from, count int }{
	{from: 0, count: 4},
	{from: 500, count: 6},
	{from: 1000, count: 8},
	{from: 1500, count: 10},
	{from: 2000, count: 12},
	{from: 2500, count: 14},
}

// TODO ---

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
			SetSpacing(theme.Current.PanelTheme.Spacing, theme.Current.PanelTheme.Spacing).
			SetColumns(BarsDefaultCount),
		bars:      make([]*container.Row, 0),
		barsCount: BarsDefaultCount,
	}

	t.add = widgets.NewButton(t.AddBar)
	t.add.AddChild(widgets.NewIcon(theme.IconPlus, theme.IconSizeMedium))

	t.AddBar()
	t.AddChild(t.add)
	t.GetWidget().OnUpdate = func(w widget.HasWidget) { t.resizeBars(false) }

	return t
}

func (t *TrackBars) HandleEvent(e any) {
}

func (t *TrackBars) AddBar() {
	bar := container.NewRow().SetBackgroundImage(theme.Current.PanelTheme.BackgroundImage)
	t.bars = append(t.bars, bar)
	t.Grid.AddChild(bar)
	t.resizeBars(true)
}

func (t *TrackBars) resizeBars(force bool) {
	availableWidth := t.GetWidget().Rect.Dx()

	barsCount := 4
	for _, bp := range BarsBreakpoints {
		if availableWidth >= bp.from {
			barsCount = bp.count
			continue
		}
		break
	}

	availableWidth -= (barsCount - 1) * theme.Current.PanelTheme.Spacing
	barWidth := availableWidth / barsCount

	if barsCount != t.barsCount {
		t.barsCount = barsCount
		t.Grid.SetColumns(barsCount)
	}

	if !force && barWidth == t.barWidth {
		return
	}

	t.barWidth = barWidth
	for _, b := range t.bars {
		b.SetMinSize(t.barWidth, BarHeight)
	}

	t.RequestRelayout()
}
