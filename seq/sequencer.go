package seq

import (
	"context"
	"stepframe/async"
	"stepframe/clock"
	"stepframe/midi"

	"github.com/rs/zerolog"
)

type Sequencer struct {
	async *async.Async[Command, Event]
	*async.Expose[Command, Event]
	logger zerolog.Logger

	receiver     *midi.Receiver
	sender       *midi.Sender
	clock        clock.Clock
	time         *TimeShifter
	forwardInOut bool
	tracks       map[int]*Track
	trackIds     int

	playRequested bool
}

func NewSequencer(
	logger zerolog.Logger,
	clock clock.Clock,
	receiver *midi.Receiver, sender *midi.Sender,
) *Sequencer {
	logger = logger.With().Str("component", "sequencer").Logger()
	as := async.NewAsync[Command, Event](logger, 16)

	return &Sequencer{
		async:        as,
		Expose:       async.NewExpose(as),
		logger:       logger,
		receiver:     receiver,
		sender:       sender,
		clock:        clock,
		time:         NewTimeShifter(),
		forwardInOut: true,
		tracks:       make(map[int]*Track),
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
		case tick := <-s.clock.Ticks():
			s.onTick(tick)
		}
	}
}

func (s *Sequencer) drainMidiEvents() []midi.Event {
	evs := make([]midi.Event, 0, 16)
	for i := 0; i < 16; i++ { // drain up to 16 events at once
		select {
		case ev := <-s.receiver.Events():
			evs = append(evs, ev)
		default:
			return evs
		}
	}
	return evs
}

func (s *Sequencer) onTick(tick clock.Tick) {
	s.time.Tick(tick.N)

	midiEvents := s.drainMidiEvents()
	isPlaying := s.playIfRequired(len(midiEvents) > 0)

	if !isPlaying {
		return
	}

	for _, t := range s.tracks {
		for _, m := range t.PollDue(s.time.Now()) {
			s.sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: m})
		}
	}

	if s.time.IsBeat(s.clock.GetTicksPerQuarter()) {
		s.async.TryDispatch(Event{Id: EvBeat})
	}

	for _, ev := range midiEvents {
		if s.forwardInOut && ev.Id == midi.EvMessage {
			s.sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: ev.Msg})
		}
		for _, t := range s.tracks {
			t.AddEvent(s.time.Now(), ev.Msg)
		}
	}
}

func (s *Sequencer) onCommand(cmd Command) {
	if cmd.TrackId != nil {
		s.onTrackCommand(cmd)
		return
	}

	switch cmd.Id {
	case CmdPlay:
		s.requestPlay()
	case CmdPause:
		s.pause()
	case CmdStop:
		s.stop()
	case CmdNewTrack:
		s.newTrack()
	default:
		s.logger.Warn().Int("cmdId", int(cmd.Id)).Msg("unknown command")
		return
	}
}

func (s *Sequencer) onTrackCommand(cmd Command) {
	id := *cmd.TrackId

	if cmd.Id == CmdRemoveTrack {
		s.removeTrack(id)
		return
	}

	s.tracks[id].HandleCommand(s.time.Now(), cmd)
	return
}

func (s *Sequencer) pause() {
	s.time.Pause()
	s.async.TryDispatch(Event{Id: EvPaused})
	s.sender.TryCommand(midi.Command{Id: midi.CmdPanic})
	s.logger.Info().Msg("play state: paused")
}

func (s *Sequencer) stop() {
	s.time.Stop(true)
	s.async.TryDispatch(Event{Id: EvStopped})
	s.sender.TryCommand(midi.Command{Id: midi.CmdPanic})
	s.logger.Info().Msg("play state: stopped")
}

func (s *Sequencer) newTrack() {
	id := s.trackIds
	s.trackIds++

	t := NewTrack(s.logger, id, s.clock, s.async.TryDispatch)
	s.tracks[id] = t
	s.async.TryDispatch(Event{Id: EvTrackAdded, TrackId: &id})
}

func (s *Sequencer) removeTrack(id int) {
	if _, ok := s.tracks[id]; !ok {
		s.logger.Warn().Int("trackId", id).Msg("invalid track id in remove command")
		return
	}

	for _, msg := range s.tracks[id].flush() {
		s.sender.TryCommand(midi.Command{Id: midi.CmdMessage, Msg: msg})
	}

	delete(s.tracks, id)
	s.async.TryDispatch(Event{Id: EvTrackRemoved, TrackId: &id})

	if len(s.tracks) == 0 {
		s.stop()
	}
}

func (s *Sequencer) requestPlay() {
	s.playRequested = true
}

func (s *Sequencer) play() {
	if s.time.State() == TsStopped {
		for _, t := range s.tracks {
			t.Reset()
		}
	}
	s.time.Play()
	s.async.TryDispatch(Event{Id: EvPlaying})
	s.logger.Info().Msg("play state: playing")
}

func (s *Sequencer) playIfRequired(hasReceivedEvent bool) bool {
	if s.time.State() == TsPlaying {
		return true
	}

	if s.playRequested {
		s.playRequested = false
		s.play()
		return true
	}

	if hasReceivedEvent {
		for _, t := range s.tracks {
			if t.scheduledState == TrackStateRecording {
				s.play()
				return true
			}
		}
	}

	return false
}
