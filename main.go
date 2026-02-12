package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"stepframe/game"
	"stepframe/seq"
	midi2 "stepframe/seq/midi"
	"stepframe/ui"
	"sync"
	"syscall"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	midi2.DebugPorts() // todo remove me

	wg := sync.WaitGroup{}
	wg.Add(1)

	sqr := seq.NewSequencer()
	sqr.AddTrack(getBillieJeanBassTrack())
	go func() {
		defer wg.Done()
		sqr.Run(ctx)
	}()

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
	wg.Wait()

	fmt.Println("GRACEFUL EXIT")
}
