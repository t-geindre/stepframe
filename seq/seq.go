package seq

import (
	"context"
	"stepframe/clock"
)

type Sequencer struct {
	tracks      []*Track
	clock       clock.Clock
	scheduler   *Scheduler
	manager     *NoteManager
	done        chan struct{} // Sequencer is done
	clockEvRate int64         // how many clock ticks per clock event
}

func NewSequencer(clock clock.Clock, clockEvRate int64) *Sequencer {
	sched := NewScheduler()

	return &Sequencer{
		tracks:      make([]*Track, 0),
		clock:       clock,
		scheduler:   sched,
		manager:     NewNoteManager(sched),
		done:        make(chan struct{}),
		clockEvRate: clockEvRate,
	}
}

func (s *Sequencer) AddTrack(tr *Track) {
	s.tracks = append(s.tracks, tr)
}

func (s *Sequencer) Wait() {
	<-s.done
}

func (s *Sequencer) Run(ctx context.Context, send func(e Event)) {
	for _, tr := range s.tracks {
		tr.Reset(0)
	}

	go s.run(ctx, send)
}

func (s *Sequencer) run(ctx context.Context, send func(e Event)) {
	defer func() {
		send(Event{Type: EvPanic}) // all note off before exit
		close(s.done)
	}()

	for {
		select {
		case <-ctx.Done():
			return

		case tk, ok := <-s.clock.Ticks():
			if !ok {
				return
			}
			now := tk.N

			// clock pulse
			if tk.N%s.clockEvRate == 0 {
				send(Event{Type: EvClock})
			}

			// poll tracks
			for _, tr := range s.tracks {
				noteEvents := tr.ProcessTick(now)
				for _, nev := range noteEvents {
					s.manager.HandleNote(nev)
				}
			}

			// send due
			for _, ev := range s.scheduler.PopDue(now) {
				send(ev)
				s.manager.OnEventSent(ev)
			}
		}
	}
}
