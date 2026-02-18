package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"
	"stepframe/ui/widgets"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/rs/zerolog"
)

type Tracks struct {
	*container.Row
	addBtn    *widgets.Button
	tacks     map[int]*Track
	sequencer *seq.Sequencer
	logger    zerolog.Logger
}

func NewTracks(logger zerolog.Logger, sequencer *seq.Sequencer) *Tracks {
	t := &Tracks{
		tacks:     make(map[int]*Track),
		sequencer: sequencer,
		logger:    logger.With().Str("component", "ui_tracks").Logger(),
	}
	t.Row = container.NewRow(widget.DirectionVertical)

	t.addBtn = widgets.NewButton(func() {
		sequencer.TryCommand(seq.Command{Id: seq.CmdNewTrack})
	})
	t.addBtn.AddChild(widgets.NewIcon(theme.IconPlus, theme.IconSizeMedium))
	t.Row.AddChild(t.addBtn)

	return t
}

func (t *Tracks) HandleEvent(e seq.Event) {
	if e.TrackId == nil {
		return // No error, ignore command
	}

	switch e.Id {
	case seq.EvTrackAdded:
		t.addTrack(*e.TrackId)
	case seq.EvTrackRemoved:
		t.removeTrack(*e.TrackId)
	default:
		if track, ok := t.tacks[*e.TrackId]; ok {
			track.HandleEvent(e)
			return
		}
		t.logger.Warn().
			Int("track_id", *e.TrackId).
			Int("event_id", int(e.Id)).
			Msg("received event for non-existent track")
	}
}

func (t *Tracks) addTrack(id int) {
	if _, ok := t.tacks[id]; ok {
		t.logger.Warn().Int("track_id", id).Msg("track already exists")
		return
	}
	tr := NewTrack(id, t.sequencer)
	t.Row.AddChild(tr)
	t.tacks[id] = tr
}

func (t *Tracks) removeTrack(id int) {
	if _, ok := t.tacks[id]; !ok {
		t.logger.Warn().Int("track_id", id).Msg("track does not exist")
		return
	}
	t.Row.RemoveChild(t.tacks[id])
	delete(t.tacks, id)
}
