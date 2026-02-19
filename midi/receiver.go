package midi

import (
	"context"
	"stepframe/async"

	"github.com/rs/zerolog"
	"gitlab.com/gomidi/midi/v2"
)

type Receiver struct {
	async *async.Async[Command, Event]
	*async.Expose[Command, Event]

	close  func()
	port   int
	logger zerolog.Logger
}

func NewReceiver(logger zerolog.Logger) *Receiver {
	logger = logger.With().Str("component", "midi_receiver").Logger()
	as := async.NewAsync[Command, Event](logger, 16)

	return &Receiver{
		async:  as,
		Expose: async.NewExpose(as),
		logger: logger,
	}
}

func (r *Receiver) Run(ctx context.Context) {
	if ok := r.async.CanRun(); !ok {
		panic(ErrAlreadyRunning) // intentional panic
	}

	go r.run(ctx)
}

func (r *Receiver) run(ctx context.Context) {
	defer func() {
		r.closePort()
		r.async.Done()
		r.logger.Info().Msg("stopped")
	}()

	r.logger.Info().Msg("running")

	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-r.async.Commands():
			r.handleCommand(cmd)
		}
	}
}

func (r *Receiver) onMessage(msg midi.Message, ts int32) {
	if r.logger.GetLevel() <= zerolog.DebugLevel {
		// avoid unnecessary string allocation
		r.logger.Debug().Str("msg", msg.String()).Int32("timestamp", ts).Msg("received midi message")
	}

	r.async.TryDispatch(Event{Id: EvMessage, Msg: msg})
}

func (r *Receiver) handleCommand(c Command) {
	switch c.Id {
	case CmdOpenPort:
		r.openPort(c.Port)
	case CmdClosePort:
		r.closePort()
	default:
		r.logger.Error().Err(ErrUnknownCommand).Msg("unknown command")
	}
}

func (r *Receiver) openPort(id int) {
	r.closePort()
	logger := r.logger.With().Int("port", id).Logger()

	ports := midi.GetInPorts()
	if id < 0 || id >= len(ports) || ports[id] == nil {
		logger.Error().Err(ErrPortNotFound).Msg("port not found")
		return
	}

	var err error
	r.close, err = midi.ListenTo(ports[id], r.onMessage)

	if err != nil {
		logger.Error().Err(err).Msg("failed to open port")
		return
	}

	r.port = id
	r.async.TryDispatch(Event{Id: EvPortOpened, Port: id})

	logger.Info().Msg("port opened")
}

func (r *Receiver) closePort() {
	if r.close == nil {
		return
	}

	r.close()
	r.close = nil
	r.async.TryDispatch(Event{Id: EvPortClosed, Port: r.port})

	r.logger.Info().Int("port", r.port).Msg("port closed")
}
