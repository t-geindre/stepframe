package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
)

type TrackCommands struct {
	widget.Containerer
	playPulse, recordPulse *widgets.Led
	playIcon               *widgets.Icon
	state                  *TrackState
}

func NewTrackCommands(state *TrackState) *TrackCommands {
	tc := &TrackCommands{state: state}

	// Play button
	playBtn := widgets.NewButton(state.TogglePlay)
	tc.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	tc.playPulse = widgets.NewLed(true, theme.IconColorNone)
	playBtn.AddChild(tc.playIcon, tc.playPulse)

	// Record button
	recordBtn := widgets.NewButton(state.ToggleRecord)
	recordIcon := widgets.NewIcon(theme.IconRecord, theme.IconSizeMedium)
	tc.recordPulse = widgets.NewLed(true, theme.IconColorNone)
	recordBtn.AddChild(recordIcon, tc.recordPulse)

	// Delete button
	deleteBtn := widgets.NewButton(state.Delete)
	deleteBtn.AddChild(widgets.NewIcon(theme.IconDelete, theme.IconSizeMedium))

	// Delete bar button
	barDeleteBtn := widgets.NewButton(nil) // Todo state.DeleteBar
	barDeleteBtn.AddChild(widgets.NewIcon(theme.IconBarDelete, theme.IconSizeMedium))

	// Add bar button
	barAddBtn := widgets.NewButton(state.AddBar)
	barAddBtn.AddChild(widgets.NewIcon(theme.IconBarAdd, theme.IconSizeMedium))

	// Container
	const columns = 3
	tc.Containerer = container.NewGrid().
		// Todo apply separator instead of more spacing
		SetSpacing(theme.Current.PanelTheme.Spacing*2, theme.Current.PanelTheme.Spacing).
		SetColumns(columns)

	// Columns
	cols := make([]widget.Containerer, columns)
	for i := 0; i < columns; i++ {
		cols[i] = container.NewRow().SetSpacing(theme.Current.PanelTheme.Spacing)
		tc.Containerer.AddChild(cols[i])
	}

	// Placement
	cols[0].AddChild(playBtn, recordBtn)
	cols[1].AddChild(deleteBtn)
	cols[2].AddChild(barAddBtn, barDeleteBtn)

	tc.applyVisualState()

	return tc
}

func (tc *TrackCommands) HandleEvent(e seq.Event) {
	if e.Id == seq.EvBeat {
		tc.playPulse.Pulse()
		tc.recordPulse.Pulse()
		return
	}

	if e.TrackId != nil && *e.TrackId == tc.state.id {
		tc.applyVisualState()
		return
	}
}

func (tc *TrackCommands) applyVisualState() {
	playing := tc.state.mode == ModePlaying || tc.state.mode == ModeRecording
	recording := tc.state.mode == ModeRecording

	if playing || tc.state.armed == ArmStop {
		tc.playIcon.SetIcon(theme.IconStop)
	} else {
		tc.playIcon.SetIcon(theme.IconPlay)
	}

	tc.playPulse.SetPulseColor(theme.IconColorDefault)
	switch {
	case tc.state.armed == ArmStop:
		tc.playPulse.SetColor(theme.IconColorArmed)
	case playing:
		tc.playPulse.SetColor(theme.IconColorOn)
	case tc.state.armed == ArmPlay || tc.state.armed == ArmRecord:
		tc.playPulse.SetColor(theme.IconColorArmed)
	default:
		tc.playPulse.SetColor(theme.IconColorIdle)
		tc.playPulse.SetPulseColor(theme.IconColorNone)
	}

	tc.recordPulse.SetPulseColor(theme.IconColorDefault)
	switch {
	case tc.state.armed == ArmRecord:
		tc.recordPulse.SetColor(theme.IconColorArmed)
	case tc.state.armed == ArmStop && recording:
		tc.recordPulse.SetColor(theme.IconColorArmed)
	case recording:
		tc.recordPulse.SetColor(theme.IconColorOn)
	default:
		tc.recordPulse.SetColor(theme.IconColorIdle)
		tc.recordPulse.SetPulseColor(theme.IconColorNone)
	}
}
