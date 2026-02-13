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

	mdi := midi.NewSender()
	mdi.Run(ctx)
	defer mdi.Wait()

	clk := clock.NewInternalClock(InternalPPQN, 120, 256)
	clk.Run(ctx)
	defer clk.Wait()

	sqr := seq.NewSequencer(clk, InternalPPQN/MidiPPQN)
	sqr.Run(ctx, mdi.Send)
	defer sqr.Wait()

	// TODO TESTS REMOVE ME
	// open and list all Midi ports
	for _, p := range midi.AllPorts() {
		fmt.Printf("PORT[%d]: %s\n", p.Id, p.Name)
		mdi.Commands() <- midi.Command{
			Id:   midi.CmdOpenPort,
			Port: p.Id,
		}
	}
	// add some tracks
	sqr.Commands() <- seq.Command{
		Id:    seq.CmdAdd,
		Track: getBillieJeanLeadTrack(),
	}
	sqr.Commands() <- seq.Command{
		Id:    seq.CmdAdd,
		Track: getBillieJeanBassTrack(),
	}
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

	fmt.Println("GRACEFUL EXIT")
}
