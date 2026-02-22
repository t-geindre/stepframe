package ui

import (
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

const (
	BarWidth  = 80 // todo move to theme
	BarHeight = 40
)

type TrackBars struct {
	widget.Containerer
	bars     []*container.Row
	add      *widgets.Button
	barWidth int
}

func NewTrackBars() *TrackBars {
	t := &TrackBars{
		Containerer: container.NewRow().
			SetPadding(theme.Current.PanelTheme.Padding).
			SetSpacing(theme.Current.PanelTheme.Spacing),
		bars:     make([]*container.Row, 0),
		barWidth: BarWidth,
	}

	t.add = widgets.NewButton(t.AddBar)
	t.add.AddChild(widgets.NewIcon(theme.IconPlus, theme.IconSizeMedium))

	t.AddBar()
	t.AddChild(t.add)

	t.Containerer.GetWidget().OnUpdate = func(w widget.HasWidget) {
		t.Resize()
	}

	return t
}

func (t *TrackBars) HandleEvent(e any) {
}

func (t *TrackBars) AddBar() {
	bar := container.NewRow().SetMinSize(t.barWidth, BarHeight).SetBackgroundImage(theme.Current.PanelTheme.BackgroundImage)
	t.bars = append(t.bars, bar)
	t.Containerer.AddChild(bar)
	t.Resize()
}

func (t *TrackBars) Resize() {
	cRect := t.Containerer.GetWidget().Rect
	addRect := t.add.GetWidget().Rect
	avWidth := cRect.Dx() - addRect.Dx() - theme.Current.PanelTheme.Padding.Left - theme.Current.PanelTheme.Padding.Right - theme.Current.PanelTheme.Spacing

	barWidth := t.barWidth
	if avWidth < barWidth*len(t.bars) {
		barWidth = avWidth / len(t.bars)
	}

}
