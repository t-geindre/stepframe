package main

import (
	"context"
	"os"
	"os/signal"
	"stepframe/clock"
	"stepframe/game"
	"stepframe/midi"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	mdi "gitlab.com/gomidi/midi/v2"
)

func main() {
	// LOGGER todo: make level configurable
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	defer logger.Info().Msg("exiting")

	// SIGNAL AWARE CONTEXT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// TIMING
	const MidiPPQN = 24 // Pulse per quarter note (PPQN)
	const InternalPPQN = 96

	clk := clock.NewInternalClock(logger, InternalPPQN, 120, 256, false)
	clk.Run(ctx)
	defer clk.Wait()

	// MIDI
	defer midi.CloseDriver(logger)

	sender := midi.NewSender(logger, 800*time.Microsecond)
	sender.Run(ctx)
	defer sender.Wait()

	receiver := midi.NewReceiver(logger)
	receiver.Run(ctx)
	defer receiver.Wait()

	// TODO TEST
	sender.TryCommand(midi.Command{Id: midi.CmdOpenPort, Port: 1})
	receiver.TryCommand(midi.Command{Id: midi.CmdOpenPort, Port: 2})

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case ev := <-receiver.Events():
				if ev.Id == midi.EvMessage {
					sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: ev.Msg})
				}
			}
		}
	}()

	defer func() {
		sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: mdi.ControlChange(1, 123, 0)})
		sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: mdi.ControlChange(1, 120, 0)})
	}()

	// SEQUENCER
	//sqr := seq.NewSequencer(clk, InternalPPQN/MidiPPQN)
	//sqr.Run(ctx, sender.Send)
	//defer sqr.Wait()

	// GUI
	//gui := ui.New(clk, sqr, sender, receiver)
	//update := game.NewUpdateFunc(func() error {
	//	gui.Update()
	//	if ctx.Err() != nil {
	//		return ebiten.Termination
	//	}
	//	return nil
	//})

	// RUN

	game.RunGame(logger, ctx /*update, gui*/)

	stop()
}
