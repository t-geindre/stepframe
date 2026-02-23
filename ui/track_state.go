package ui

import "stepframe/seq"

type Mode uint8

const (
	ModeStopped Mode = iota
	ModePlaying
	ModeRecording
)

type Armed uint8

const (
	ArmNone Armed = iota
	ArmPlay
	ArmRecord
	ArmStop
)

type TrackState struct {
	id        int
	mode      Mode
	armed     Armed
	sequencer *seq.Sequencer
}

func NewTrackState(id int, sequencer *seq.Sequencer) *TrackState {
	return &TrackState{
		mode:      ModeStopped,
		armed:     ArmNone,
		sequencer: sequencer,
	}
}

func (s *TrackState) TogglePlay() {
	if s.mode != ModeStopped || s.armed == ArmPlay || s.armed == ArmRecord {
		s.sequencer.TryCommand(seq.Command{Id: seq.CmdStop, TrackId: &s.id})
	} else {
		s.sequencer.TryCommand(seq.Command{Id: seq.CmdPlay, TrackId: &s.id})
	}
}

func (s *TrackState) ToggleRecord() {
	if s.mode == ModeRecording || s.armed == ArmRecord {
		s.sequencer.TryCommand(seq.Command{Id: seq.CmdStopRecord, TrackId: &s.id})
	} else {
		s.sequencer.TryCommand(seq.Command{Id: seq.CmdRecord, TrackId: &s.id})
	}
}

func (s *TrackState) Delete() {
	s.sequencer.TryCommand(seq.Command{Id: seq.CmdRemoveTrack, TrackId: &s.id})
}

func (s *TrackState) HandleEvent(e seq.Event) {
	if e.TrackId != nil && *e.TrackId != s.id {
		return
	}

	switch e.Id {
	case seq.EvPlaying:
		s.mode = ModePlaying
		s.armed = ArmNone
	case seq.EvRecording:
		s.mode = ModeRecording
		s.armed = ArmNone
	case seq.EvStopped:
		s.mode = ModeStopped
		s.armed = ArmNone
	case seq.EvArmedPlaying:
		s.armed = ArmPlay
	case seq.EvArmedRecording:
		s.armed = ArmRecord
	case seq.EvArmedStopped:
		s.armed = ArmStop
	}
}
