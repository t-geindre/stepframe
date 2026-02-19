package main

import (
	"context"
	"os"
	"os/signal"
	"stepframe/clock"
	"stepframe/game"
	"stepframe/midi"
	"stepframe/seq"
	"stepframe/ui"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	// LOGGER todo: make level configurable
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	logger := zerolog.New(output).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	defer logger.Info().Msg("done")

	// SIGNAL AWARE CONTEXT
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	// TIMING
	const InternalPPQN = 96

	clk := clock.NewInternalClock(logger, InternalPPQN, 120, 256, false)
	clk.Run(ctx)
	defer clk.Wait()

	// MIDI
	midi.Open(logger)
	defer midi.Close(logger)

	sender := midi.NewSender(logger, 800*time.Microsecond)
	sender.Run(ctx)
	defer sender.Wait()

	receiver := midi.NewReceiver(logger)
	receiver.Run(ctx)
	defer receiver.Wait()

	// TODO TEST
	sender.TryCommand(midi.Command{Id: midi.CmdOpenPort, Port: 1})
	receiver.TryCommand(midi.Command{Id: midi.CmdOpenPort, Port: 2})

	// SEQUENCER
	sqr := seq.NewSequencer(logger, clk, receiver, sender)
	sqr.Run(ctx)
	defer sqr.Wait()

	// GUI
	gui := ui.New(clk, sqr, sender, receiver, logger)

	// RUN
	game.RunGame(logger, ctx, gui)
	stop()
}
