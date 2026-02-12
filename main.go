package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"stepframe/clock"
	"stepframe/game"
	"stepframe/midi"
	"stepframe/seq"
	"stepframe/ui"
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Pulse per quarter note (PPQN)
	const MidiPPQN = 24
	const InternalPPQN = 96

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	mdi := midi.NewOut(1)

	clk := clock.NewInternalClock(InternalPPQN, 120, 256)
	clk.Run(ctx)

	sqr := seq.NewSequencer(clk, InternalPPQN/MidiPPQN)
	sqr.AddTrack(getBillieJeanBassTrack())
	sqr.Run(ctx, mdi.SendEvent)

	gui := ui.New()
	update := game.NewUpdateFunc(func() error {
		if ctx.Err() != nil {
			return ebiten.Termination
		}
		return nil
	})

	err := ebiten.RunGame(game.NewGame(update, gui))
	if err != nil && !errors.Is(err, ebiten.Termination) {
		panic(err)
	}

	stop()

	clk.Wait()
	sqr.Wait()

	fmt.Println("GRACEFUL EXIT")
}
