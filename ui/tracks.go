package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/rs/zerolog"
)

type Tracks struct {
	*container.Row
	tacks     map[int]*Track
	sequencer *seq.Sequencer
	logger    zerolog.Logger
}

func NewTracks(logger zerolog.Logger, sequencer *seq.Sequencer) *Tracks {
	return &Tracks{
		Row: container.NewVerticalRow().AlignContent(widget.RowLayoutPositionStart).
			WithPadding().StretchContent(),
		tacks:     make(map[int]*Track),
		sequencer: sequencer,
		logger:    logger.With().Str("component", "ui_tracks").Logger(),
	}
}

func (t *Tracks) HandleEvent(e seq.Event) {
	for _, track := range t.tacks {
		track.HandleEvent(e)
	}

	if e.TrackId == nil {
		return // No error, ignore command
	}

	switch e.Id {
	case seq.EvTrackAdded:
		t.addTrack(*e.TrackId)
	case seq.EvTrackRemoved:
		t.removeTrack(*e.TrackId)
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
