package main

import (
	"context"
	"os"
	"os/signal"
	"stepframe/clock"
	"stepframe/engine"
	"stepframe/midi"
	"stepframe/seq"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	midi.DebugPorts() // todo remove me
	mdiOut := midi.NewOut(1)

	sched := engine.NewScheduler()
	nm := engine.NewNoteManager(sched)

	clk := clock.NewInternalClock(120.0, 256)
	clk.Start(ctx)

	// cleanup on exit
	defer func() {
		stop()
		mdiOut.PanicAll()
		clk.Stop()
	}()

	tracks := []*seq.Track{getBillieJeanBassTrack()}
	//tracks := []*seq.Track{}

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
			for _, tr := range tracks {
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
