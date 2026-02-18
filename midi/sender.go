package midi

import (
	"context"
	"stepframe/async"
	"time"

	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
)

type Sender struct {
	async *async.Async[Command, Event]
	*async.Expose[Command, Event]

	port     int
	close    func() error
	send     func(midi.Message) error
	stagger  time.Duration
	lastSend time.Time
	logger   zerolog.Logger
}

func NewSender(logger zerolog.Logger, stagger time.Duration) *Sender {
	logger = logger.With().Str("component", "midi sender").Logger()
	as := async.NewAsync[Command, Event](logger, 16)

	return &Sender{
		async:   as,
		Expose:  async.NewExpose(as),
		stagger: stagger,
		logger:  logger,
		port:    -1,
	}
}

func (s *Sender) Run(ctx context.Context) {
	if ok := s.async.CanRun(); !ok {
		panic(ErrAlreadyRunning) // intentional panic
	}

	go s.run(ctx)
}

func (s *Sender) run(ctx context.Context) {
	defer func() {
		s.closePort()

		s.async.Done()
		s.logger.Info().Msg("stopped")
	}()

	s.logger.Info().Msg("running")

	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-s.async.Commands():
			s.handleCommand(cmd)
		}
	}
}

func (s *Sender) handleCommand(c Command) {
	switch c.Id {
	case CmdOpenPort:
		s.openPort(c.Port)
	case CmdClosePort:
		s.closePort()
	case CmdMessage:
		s.sendMessage(c.Msg)
	case CmdPanic:
		s.panic()

	default:
		s.logger.Err(ErrUnknownCommand).Msg("unknown command")
	}
}

func (s *Sender) openPort(id int) {
	s.closePort()
	logger := s.logger.With().Int("port", id).Logger()

	out, err := midi.OutPort(id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to open port")
		return
	}

	send, err := midi.SendTo(out)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create sender")
		_ = out.Close() // best-effort cleanup
		return
	}

	s.port = id
	s.close = out.Close
	s.send = send

	s.logger.Info().Int("port", s.port).Msg("port opened")
	s.async.TryDispatch(Event{Id: EvPortOpened, Port: id})
}

func (s *Sender) closePort() {
	if s.close == nil {
		return
	}

	logger := s.logger.With().Int("port", s.port).Logger()

	logger.Info().Msg("all notes off")
	s.panic()

	err := s.close()
	if err != nil {
		logger.Err(err).Msg("failed to close port")
		return
	}

	s.close = nil
	s.send = nil

	logger.Info().Msg("port closed")
	s.async.TryDispatch(Event{Id: EvPortClosed, Port: s.port})

	s.port = -1
}

func (s *Sender) sendMessage(m midi.Message) {
	if s.send == nil {
		s.logger.Err(ErrPortNotFound).Msg("no port open, cannot send message")
		return
	}

	if s.stagger > 0 && !s.lastSend.IsZero() {
		// stager enabled: flood avoidance
		elapsed := time.Since(s.lastSend)
		if elapsed < s.stagger {
			time.Sleep(s.stagger - elapsed)
		}
	}

	if s.logger.GetLevel() <= zerolog.DebugLevel {
		// avoid unnecessary string allocation
		s.logger.Debug().Str("msg", m.String()).Msg("send midi message")
	}

	err := s.send(m)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to send message")
		return
	}

	s.lastSend = time.Now()
}

func (s *Sender) panic() {
	for i := uint8(0); i < 16; i++ {
		s.sendMessage(midi.ControlChange(i, 123, 0))
		s.sendMessage(midi.ControlChange(i, 120, 0))
	}
}
