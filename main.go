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

	sender := midi.NewSender(800 * time.Microsecond)
	sender.Run(ctx)
	defer sender.Wait()

	receiver := midi.NewReceiver(sender)
	receiver.Run(ctx)
	defer receiver.Wait()

	clk := clock.NewInternalClock(InternalPPQN, 120, 256)
	clk.Run(ctx)
	defer clk.Wait()

	sqr := seq.NewSequencer(clk, InternalPPQN/MidiPPQN)
	sqr.Run(ctx, sender.Send)
	defer sqr.Wait()

	gui := ui.New(clk, sqr, sender, receiver)
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

	fmt.Println("GRACEFUL EXIT")
}
