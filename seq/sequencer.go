package seq

import (
	"context"
	"stepframe/async"
	"stepframe/clock"
	"stepframe/midi"
	"stepframe/track"

	"github.com/rs/zerolog"
)

type Sequencer struct {
	async *async.Async[Command, Event]
	*async.Expose[Command, Event]

	logger zerolog.Logger

	in    *midi.Receiver
	out   *midi.Sender
	clock clock.Clock

	forwardInOut bool
	now          int64

	track track.Track
}

func NewSequencer(
	logger zerolog.Logger,
	clock clock.Clock,
	in *midi.Receiver, out *midi.Sender,
) *Sequencer {
	logger = logger.With().Str("component", "sequencer").Logger()
	as := async.NewAsync[Command, Event](logger, 16)

	return &Sequencer{
		async:  as,
		Expose: async.NewExpose(as),
		logger: logger,

		clock: clock,
		in:    in,
		out:   out,

		forwardInOut: true,

		track: track.NewQuantizer(
			track.NewTrack(384),
			clock.GetTicksPerQuarter()/4,
			track.QuantizeNearest,
		),
	}
}

func (s *Sequencer) Run(ctx context.Context) {
	if ok := s.async.CanRun(); !ok {
		panic(ErrAlreadyRunning) // intentional panic
	}

	go s.run(ctx)
}

func (s *Sequencer) run(ctx context.Context) {
	defer func() {
		s.async.Done()
		s.logger.Info().Msg("stopped")
	}()

	s.logger.Info().Msg("running")

	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-s.async.Commands():
			_ = cmd
			s.logger.Warn().Msg("received command")
		case ev := <-s.in.Events():
			s.onInEvent(ev)
		case tick := <-s.clock.Ticks():
			s.onTick(tick)
		}
	}
}

func (s *Sequencer) onInEvent(ev midi.Event) {
	if s.forwardInOut && ev.Id == midi.EvMessage {
		s.out.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: ev.Msg})
	}

	s.track.AddEvent(s.now, ev.Msg)
}

func (s *Sequencer) onTick(tick clock.Tick) {
	s.now = tick.N
	for _, m := range s.track.PollDue(s.now) {
		s.out.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: m})
	}

}
