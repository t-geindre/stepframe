package seq

import (
	"context"
	"stepframe/seq/clock"
	"stepframe/seq/engine"
	"stepframe/seq/midi"
)

type Sequencer struct {
	tracks []*engine.Track
}

func NewSequencer() *Sequencer {
	return &Sequencer{
		tracks: make([]*engine.Track, 0),
	}
}

func (s *Sequencer) AddTrack(tr *engine.Track) {
	s.tracks = append(s.tracks, tr)
}

func (s *Sequencer) Run(ctx context.Context) {
	mdiOut := midi.NewOut(1)

	sched := engine.NewScheduler()
	nm := engine.NewNoteManager(sched)

	clk := clock.NewInternalClock(120.0, 256)
	clk.Start(ctx)

	// cleanup on exit
	defer func() {
		mdiOut.PanicAll()
		clk.Stop()
	}()

	for _, tr := range s.tracks {
		tr.Reset(0)
	}

	for {
		select {
		case <-ctx.Done():
			return

		case tk, ok := <-clk.Ticks():
			if !ok {
				return
			}
			now := tk.N

			// send clock pulse
			if tk.N%4 == 0 {
				// clock runs 4 times faster than MIDI PPQN
				_ = mdiOut.SendClockPulse()
			}

			// poll tracks
			for _, tr := range s.tracks {
				noteEvents := tr.ProcessTick(now)
				for _, nev := range noteEvents {
					nm.HandleNote(nev)
				}
			}

			// send due
			for _, ev := range sched.PopDue(now) {
				_ = mdiOut.SendEvent(ev)
				nm.OnEventSent(ev)
			}
		}
	}
}
