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

	receiver *midi.Receiver
	sender   *midi.Sender

	clock clock.Clock
	time  *TimeShifter

	forwardInOut bool

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

		receiver: in,
		sender:   out,

		clock: clock,
		time:  NewTimeShifter(),

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
			s.onCommand(cmd)
		case ev := <-s.receiver.Events():
			s.onMidiEvent(ev)
		case tick := <-s.clock.Ticks():
			s.onTick(tick)
		}
	}
}

func (s *Sequencer) onMidiEvent(ev midi.Event) {
	if s.forwardInOut && ev.Id == midi.EvMessage {
		s.sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: ev.Msg})
	}

	if s.time.State() != TsPlaying {
		return
	}

	s.track.AddEvent(s.time.Now(), ev.Msg)
}

func (s *Sequencer) onCommand(cmd Command) {
	switch cmd.Id {
	case CmdPlay:
		if s.time.State() == TsStopped {
			s.track.Reset()
		}
		s.time.Play()
		s.async.TryDispatch(Event{Id: EvPlaying})
	case CmdPause:
		s.time.Pause()
		s.async.TryDispatch(Event{Id: EvPaused})
		s.sender.TryCommand(midi.Command{Id: midi.CmdPanic})
	case CmdStop:
		s.time.Stop(true)
		s.async.TryDispatch(Event{Id: EvStopped})
		s.sender.TryCommand(midi.Command{Id: midi.CmdPanic})
	}
}

func (s *Sequencer) onTick(tick clock.Tick) {
	s.time.Tick(tick.N)

	if s.time.State() != TsPlaying {
		return
	}

	now := s.time.Now()

	for _, m := range s.track.PollDue(now) {
		s.sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: m})
	}

	if s.time.IsBeat(s.clock.GetTicksPerQuarter()) {
		s.async.TryDispatch(Event{Id: EvBeat})
	}
}
