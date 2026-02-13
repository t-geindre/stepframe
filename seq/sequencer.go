package seq

import (
	"context"
	"stepframe/clock"
	"sync/atomic"
)

type Sequencer struct {
	tracks      []*TrackState
	clock       clock.Clock
	scheduler   *Scheduler
	manager     *NoteManager
	done        chan struct{} // Sequencer is done
	running     atomic.Bool
	clockEvRate int64 // how many clock ticks per clock event
	commands    chan Command
}

func NewSequencer(clock clock.Clock, clockEvRate int64) *Sequencer {
	sched := NewScheduler()

	return &Sequencer{
		tracks:      make([]*TrackState, 0),
		clock:       clock,
		scheduler:   sched,
		manager:     NewNoteManager(sched),
		done:        make(chan struct{}),
		clockEvRate: clockEvRate,
		commands:    make(chan Command, 16),
	}
}

func (s *Sequencer) Wait() {
	<-s.done
}

func (s *Sequencer) Run(ctx context.Context, send func(e Event)) {
	if !s.running.CompareAndSwap(false, true) {
		return
	}

	done := make(chan struct{})
	s.done = done

	for _, tr := range s.tracks {
		tr.Reset(0)
	}

	go s.run(ctx, send, done)
}

func (s *Sequencer) Commands() chan<- Command {
	return s.commands
}

func (s *Sequencer) run(ctx context.Context, send func(e Event), done chan struct{}) {
	defer func() {
		send(Event{Type: EvPanic}) // all note off before exit
		close(done)                // close local done
		s.running.Store(false)
	}()

	evBuf := make([]NoteEvent, 0, 16)

	for {
		select {
		case <-ctx.Done():
			return

		case tk, ok := <-s.clock.Ticks():
			if !ok {
				return
			}
			now := tk.N
			s.clockPulse(now, send)
			s.processTracks(now, send, evBuf)
			s.processCommands(now)
		}
	}
}

func (s *Sequencer) clockPulse(now int64, send func(e Event)) {
	return // todo temp disable clock events, overflows volca buffers
	if now%s.clockEvRate == 0 {
		send(Event{Type: EvClock})
	}
}

func (s *Sequencer) processTracks(now int64, send func(e Event), evBuf []NoteEvent) {
	for _, tr := range s.tracks {
		if tr.playAt > 0 && now >= tr.playAt {
			s.playTrackNow(now, tr)
		}

		if tr.stopAt > 0 && now >= tr.stopAt {
			tr.stopAt = 0
			tr.play = false
			continue
		}

		if !tr.play {
			continue
		}

		evBuf := make([]NoteEvent, 0, len(evBuf)) // todo temp fix
		noteEvents := tr.ProcessTick(now, evBuf)

		if tr.muted {
			continue
		}

		for _, nev := range noteEvents {
			s.manager.HandleNote(nev)
		}
	}

	for _, ev := range s.scheduler.PopDue(now) {
		send(ev)
		s.manager.OnEventSent(ev)
	}
}

func (s *Sequencer) processCommands(now int64) {
	for {
		select {
		case c := <-s.commands:
			switch c.Id {
			case CmdPlay:
				s.playTrackAt(now, c.TrackId, c.At)
			case CmdStop:
				s.stopTrackAt(now, c.TrackId, c.At)
			case CmdAdd:
				s.addTrack(c.Track)
			case CmdRemove:
			case CmdSwap:
			default:
				panic("invalid command")
			}
		default:
			return
		}
	}
}

func (s *Sequencer) addTrack(tr *Track) {
	s.tracks = append(s.tracks, NewTrackState(tr))
}

func (s *Sequencer) findTrack(id TrackId) *TrackState {
	for _, tr := range s.tracks {
		if tr.track.Id() == id {
			return tr
		}
	}

	return nil
}

func (s *Sequencer) playTrackAt(now int64, id TrackId, at CommandAt) {
	tr := s.findTrack(id)

	if tr == nil || tr.play {
		return
	}

	atTk := s.getAtTick(now, at)
	if atTk <= now {
		s.playTrackNow(now, tr)
		return
	}
	tr.playAt = atTk
}

func (s *Sequencer) playTrackNow(now int64, tr *TrackState) {
	tr.playAt = 0
	tr.play = true
	tr.Reset(now)
}

func (s *Sequencer) stopTrackAt(now int64, id TrackId, at CommandAt) {
	tr := s.findTrack(id)

	if tr == nil || !tr.play {
		return
	}

	atTk := s.getAtTick(now, at)
	if atTk <= now {
		tr.stopAt = 0
		tr.play = false
		return
	}
	tr.stopAt = atTk
}

func (s *Sequencer) getAtTick(now int64, at CommandAt) int64 {
	switch at {
	case CmdAtNextBar:
		ppqn := s.clock.GetTicksPerQuarter() * 4 // todo assume 4/4 time signature for now

		rem := now % ppqn
		if rem == 0 {
			return now
		}

		return now + (ppqn - rem)

	case CmdAtNextBeat:
		ppqn := s.clock.GetTicksPerQuarter()

		rem := now % ppqn
		if rem == 0 {
			return now
		}

		return now + (ppqn - rem)

	case CmdAtNow:
		return now

	default:
		panic("invalid at command")
	}
}
