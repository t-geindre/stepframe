package seq

import (
	"stepframe/clock"
	"stepframe/track"

	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
)

type TrackState int

const (
	TrackStateStopped TrackState = iota
	TrackStateRecording
	TrackStatePlaying
	TrackStateNone
)

type Track struct {
	track.Track
	clock           clock.Clock
	quantizer       *track.Quantizer
	state           TrackState
	scheduledState  TrackState
	scheduledAt     int64
	dispatch        func(Event) bool
	logger          zerolog.Logger
	id              int
	recordRealStart int64
}

func NewTrack(logger zerolog.Logger, id int, clock clock.Clock, dispatch func(Event) bool) *Track {
	tr := track.NewTrack(clock.GetTicksPerQuarter() * 4) // todo one bar should be configurable
	quantizer := track.NewQuantizer(
		tr,
		clock.GetTicksPerQuarter()/4, // todo time signature should be configurable
		track.QuantizeNearest,
	)

	return &Track{
		Track:          quantizer,
		quantizer:      quantizer,
		clock:          clock,
		state:          TrackStateStopped,
		scheduledState: TrackStateNone,
		logger:         logger.With().Str("component", "track").Logger(),
		dispatch:       dispatch,
		id:             id,
	}
}
func (t *Track) PollDue(nowLocal int64) []midi.Message {
	t.setDueState(nowLocal)

	if t.state != TrackStatePlaying && t.state != TrackStateRecording {
		return nil
	}

	// Auto stop recording
	if t.state == TrackStateRecording && t.recordRealStart >= 0 {
		if nowLocal-t.recordRealStart > t.Track.GetLengthTick() {
			t.setState(TrackStatePlaying)
		}
	}

	return t.Track.PollDue(nowLocal)
}

func (t *Track) AddEvent(atLocalTick int64, msg midi.Message) {
	t.setDueState(atLocalTick)

	if t.state == TrackStateRecording {
		if t.recordRealStart < 0 {
			t.recordRealStart = atLocalTick
		}
		t.Track.AddEvent(atLocalTick, msg)
	}
}

func (t *Track) HandleCommand(nowLocal int64, cmd Command) {
	switch cmd.Id {
	case CmdPlay:
		t.scheduleState(TrackStatePlaying, t.getNextBarTick(nowLocal))
	case CmdStop:
		if t.state == TrackStateRecording {
			t.setState(TrackStatePlaying)
		}
		t.scheduleState(TrackStateStopped, t.getNextBarTick(nowLocal))
	case CmdRecord:
		t.scheduleState(TrackStateRecording, t.getNextBarTick(nowLocal))
	case CmdStopRecord:
		if t.state == TrackStateRecording {
			t.setState(TrackStatePlaying)
		} else if t.scheduledState == TrackStateRecording {
			t.scheduleState(TrackStatePlaying, t.getNextBarTick(nowLocal))
		}
	default:
		t.logger.Warn().Int("cmdId", int(cmd.Id)).Msg("unknown command")
		return
	}
}

func (t *Track) getNextBarTick(nowLocal int64) int64 {
	ticksPerBar := t.clock.GetTicksPerQuarter() * 4 // todo configurable
	if ticksPerBar <= 0 {
		return nowLocal
	}
	// Now
	if nowLocal%ticksPerBar == 0 {
		return nowLocal
	}
	// Next bar
	return ((nowLocal / ticksPerBar) + 1) * ticksPerBar
}

func (t *Track) setState(state TrackState) {
	t.state = state
	t.scheduledState = TrackStateNone
	t.scheduledAt = 0

	switch state {
	case TrackStatePlaying:
		t.Track.Reset()
		t.dispatch(Event{Id: EvPlaying, TrackId: &t.id})
	case TrackStateStopped:
		t.Track.Reset()
		t.dispatch(Event{Id: EvStopped, TrackId: &t.id})
	case TrackStateRecording:
		t.recordRealStart = -1
		t.Track.Reset()
		t.dispatch(Event{Id: EvRecording, TrackId: &t.id})
	default:
		t.logger.Warn().Int("state", int(state)).Msg("unknown state")
		return
	}
}

func (t *Track) scheduleState(state TrackState, atLocal int64) {
	t.scheduledState = state
	t.scheduledAt = atLocal

	switch state {
	case TrackStatePlaying:
		t.dispatch(Event{Id: EvArmedPlaying, TrackId: &t.id})
	case TrackStateStopped:
		if t.state != TrackStateRecording && t.state != TrackStatePlaying {
			t.setState(TrackStateStopped)
		} else {
			t.dispatch(Event{Id: EvArmedStopped, TrackId: &t.id})
		}
	case TrackStateRecording:
		t.dispatch(Event{Id: EvArmedRecording, TrackId: &t.id})
	default:
		t.logger.Warn().Int("state", int(state)).Msg("unknown state")
		return
	}

	t.logger.Debug().Int("track", t.id).Int64("at", atLocal).Msg("schedule")
}

func (t *Track) setDueState(nowLocal int64) {
	if t.scheduledState != TrackStateNone && nowLocal >= t.scheduledAt {
		t.setState(t.scheduledState)
	}
}
