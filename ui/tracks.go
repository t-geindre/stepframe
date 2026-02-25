package ui

import (
	"stepframe/seq"
	"stepframe/ui/container"
	"stepframe/ui/theme"

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
		Row: container.NewRow().
			SetDirection(widget.DirectionVertical).
			SetContentPosition(widget.RowLayoutPositionStart).
			SetSpacing(theme.Current.PanelTheme.Spacing).
			SetContentStretch(true),
		tacks:     make(map[int]*Track),
		sequencer: sequencer,
		logger:    logger.With().Str("component", "ui_tracks").Logger(),
	}
}

func (t *Tracks) HandleEvent(e seq.Event) {
	if e.TrackId != nil {
		switch e.Id {
		case seq.EvTrackAdded:
			t.addTrack(*e.TrackId)
			return
		case seq.EvTrackRemoved:
			t.removeTrack(*e.TrackId)
			return
		}
	}

	for _, track := range t.tacks {
		track.HandleEvent(e)
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
