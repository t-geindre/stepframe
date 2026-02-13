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
	"time"

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
	sqr.Run(ctx, mdi.SendEvent)

	// TODO TESTS REMOVE ME
	time.AfterFunc(time.Millisecond*100, func() {
		sqr.Commands() <- seq.Command{
			Id:    seq.CmdAdd,
			Track: getBillieJeanLeadTrack(),
		}
		sqr.Commands() <- seq.Command{
			Id:    seq.CmdAdd,
			Track: getBillieJeanBassTrack(),
		}
	})
	// TODO END ---

	gui := ui.New(sqr)
	update := game.NewUpdateFunc(func() error {
		gui.Update()
		if ctx.Err() != nil {
			return ebiten.Termination
		}
		return nil
	})

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(game.NewGame(update, gui))
	if err != nil && !errors.Is(err, ebiten.Termination) {
		panic(err)
	}

	stop()

	clk.Wait()
	sqr.Wait()

	fmt.Println("GRACEFUL EXIT")
}
