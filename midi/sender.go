package midi

import (
	"context"
	"errors"
	"stepframe/seq"
	"sync/atomic"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

var ErrUnknownEventType = errors.New("Unknown MIDI event type")

type Sender struct {
	ports    map[int]*Port
	done     chan struct{}
	running  atomic.Bool
	queue    chan seq.Event
	commands chan Command
	stagger  time.Duration
	lastSend map[int]time.Time
}

func NewSender(stagger time.Duration) *Sender {
	return &Sender{
		ports:    make(map[int]*Port),
		queue:    make(chan seq.Event, 1024),
		commands: make(chan Command, 16),
		stagger:  stagger,
		lastSend: make(map[int]time.Time),
	}
}

func (s *Sender) Send(e seq.Event) {
	if e.Type == seq.EvNoteOff || e.Type == seq.EvPanic {
		s.queue <- e // never drop NoteOff or Panic events
		return
	}

	select {
	case s.queue <- e:
	default:
		// queue full, drop event
	}
}

func (s *Sender) Run(ctx context.Context) {
	if !s.running.CompareAndSwap(false, true) {
		return
	}

	done := make(chan struct{})
	s.done = done

	go s.run(ctx, done)
}

func (s *Sender) Wait() {
	done := s.done
	if done == nil {
		return
	}
	<-done
}

func (s *Sender) Commands() chan<- Command {
	return s.commands
}

func (s *Sender) run(ctx context.Context, done chan struct{}) {
	defer func() {
		close(done) // close local done
		s.running.Store(false)
	}()

	for {
		select {
		case <-ctx.Done():
			s.drainEvents(nil)
			return
		case e := <-s.queue:
			s.drainEvents(&e)
		case cmd := <-s.commands:
			s.drainCommands(cmd)
		}
	}
}

func (s *Sender) drainEvents(e *seq.Event) {
	if e != nil {
		s.handleEvent(*e)
	}
	for i := 0; i < 1024; i++ { // prevent starvation
		select {
		case e := <-s.queue:
			s.handleEvent(e)
		default:
			return
		}
	}
}

func (s *Sender) handleEvent(e seq.Event) {
	if e.Type == seq.EvPanic {
		for _, p := range s.ports {
			for ch := uint8(0); ch < 16; ch++ {
				_ = p.send(midi.ControlChange(ch, 120, 0)) // All Sound Off
				_ = p.send(midi.ControlChange(ch, 123, 0)) // All Notes Off
			}
		}

		return
	}

	// --- STAGGER: ensure we don't flood the same port ---
	if s.stagger > 0 {
		if last, ok := s.lastSend[e.Port]; ok {
			elapsed := time.Since(last)
			if elapsed < s.stagger {
				time.Sleep(s.stagger - elapsed)
			}
		}
	}

	p, ok := s.ports[e.Port]
	if !ok {
		return
	}

	var err error

	switch e.Type {
	case seq.EvNoteOn:
		err = p.send(midi.NoteOn(e.Channel, e.Note, e.Vel))
	case seq.EvNoteOff:
		err = p.send(midi.NoteOff(e.Channel, e.Note))
	case seq.EvCC:
		err = p.send(midi.ControlChange(e.Channel, e.CC, e.Value))
	case seq.EvClock:
		err = p.send(midi.TimingClock())
	default:
		err = ErrUnknownEventType
	}

	s.lastSend[e.Port] = time.Now()

	if err != nil {
		// TODO don't panic, log error
		panic("MIDI Send error: " + err.Error())
	}
}

func (s *Sender) drainCommands(c Command) {
	s.handleCommand(c)
	for i := 0; i < 1024; i++ { // prevent starvation
		select {
		case c := <-s.commands:
			s.handleCommand(c)
		default:
			return
		}
	}
}

func (s *Sender) handleCommand(c Command) {
	switch c.Id {
	case CmdOpenPort:
		out, err := midi.OutPort(c.Port)
		if err != nil {
			// TODO don't panic, log error
			panic("MIDI OpenPort error: " + err.Error())
		}
		send, err := midi.SendTo(out)
		if err != nil {
			// TODO don't panic, log error
			panic("MIDI OpenPort error: " + err.Error())
		}
		s.ports[c.Port] = &Port{
			send:  send,
			close: out.Close,
		}

	case CmdClosePort:
		if p, ok := s.ports[c.Port]; ok {
			_ = p.close()
			delete(s.ports, c.Port)
		}
	}
}
