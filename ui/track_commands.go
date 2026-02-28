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
	playPulse, recordPulse *container.Row
	playIcon               *widgets.Icon
	state                  *TrackState
}

func NewTrackCommands(state *TrackState) *TrackCommands {
	tc := &TrackCommands{state: state}

	// Play button
	playBtn := widgets.NewButton(state.TogglePlay)
	tc.playIcon = widgets.NewIcon(theme.IconPlay, theme.IconSizeMedium)
	tc.playPulse = container.NewRow().SetTheme(theme.Current.LedPanelTheme)
	playBtn.AddChild(tc.playIcon, tc.playPulse)

	// Record button
	recordBtn := widgets.NewButton(state.ToggleRecord)
	recordIcon := widgets.NewIcon(theme.IconRecord, theme.IconSizeMedium)
	tc.recordPulse = container.NewRow().SetTheme(theme.Current.LedPanelTheme)
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
		SetStretch([]bool{false, false, true}, []bool{true}).
		SetColumns(columns)

	// Columns
	cols := make([]widget.Containerer, columns)
	for i := 0; i < columns; i++ {
		cols[i] = container.NewRow().SetContentPosition(widget.RowLayoutPositionCenter).SetSpacing(theme.Current.PanelTheme.Spacing)
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
		playing := tc.state.mode == ModePlaying || tc.state.mode == ModeRecording || tc.state.armed == ArmPlay || tc.state.armed == ArmRecord
		recording := tc.state.mode == ModeRecording || tc.state.armed == ArmRecord

		if playing {
			tc.playPulse.Pulse()
		}

		if recording {
			tc.recordPulse.Pulse()
		}

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

	switch {
	case tc.state.armed == ArmStop:
		tc.playPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedArmed)
	case playing:
		tc.playPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedOn)
	case tc.state.armed == ArmPlay || tc.state.armed == ArmRecord:
		tc.playPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedArmed)
	default:
		tc.playPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedIdle)
	}

	switch {
	case tc.state.armed == ArmRecord:
		tc.recordPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedArmed)
	case tc.state.armed == ArmStop && recording:
		tc.recordPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedArmed)
	case recording:
		tc.recordPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedOn)
	default:
		tc.recordPulse.SetColorizedBackgroundOffsetImage(theme.ColorLedIdle)
	}
}
