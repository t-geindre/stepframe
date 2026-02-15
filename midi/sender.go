package midi

import (
	"context"
	"stepframe/seq"
	"sync/atomic"
	"time"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

type message struct {
	port int
	msg  midi.Message
}

type Sender struct {
	ports    map[int]*Port
	done     chan struct{}
	running  atomic.Bool
	queue    chan message
	commands chan Command
	stagger  time.Duration
	lastSend map[int]time.Time
}

func NewSender(stagger time.Duration) *Sender {
	return &Sender{
		ports:    make(map[int]*Port),
		queue:    make(chan message, 1024),
		commands: make(chan Command, 16),
		stagger:  stagger,
		lastSend: make(map[int]time.Time),
	}
}

func (s *Sender) Send(e seq.Event) {
	if e.Type == seq.EvPanic {
		for p, _ := range s.ports {
			for ch := uint8(0); ch < 16; ch++ {
				s.SendRaw(p, midi.ControlChange(ch, 120, 0)) // All Sound Off
				s.SendRaw(p, midi.ControlChange(ch, 123, 0)) // All Notes Off
			}
		}
		return
	}

	if e.Type == seq.EvNoteOff {
		s.SendRaw(e.Port, midi.NoteOff(e.Channel, e.Note))
	}

	switch e.Type {
	case seq.EvNoteOn:
		s.SendRaw(e.Port, midi.NoteOn(e.Channel, e.Note, e.Vel))
	case seq.EvNoteOff:
		s.SendRaw(e.Port, midi.NoteOff(e.Channel, e.Note))
	case seq.EvClock:
		s.SendRaw(e.Port, midi.TimingClock())
	default:
		// TODO don't panic, log error
		panic(ErrUnknownEventType)
	}
}

func (s *Sender) SendRaw(port int, msg midi.Message) {
	if port == 0 {
		return
	}
	s.queue <- message{
		port: port,
		msg:  msg,
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
	// ensure all remaining events are sent
	time.Sleep(time.Millisecond * 5)

	s.drainMessages(nil)
	// close all ports
	for port := range s.ports {
		s.closePort(port)
	}
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
			return
		case m := <-s.queue:
			s.drainMessages(&m)
		case cmd := <-s.commands:
			s.drainCommands(cmd)
		}
	}
}

func (s *Sender) drainMessages(m *message) {
	if m != nil {
		s.handleMessage(*m)
	}
	for i := 0; i < 1024; i++ { // prevent starvation
		select {
		case m := <-s.queue:
			s.handleMessage(m)
		default:
			return
		}
	}
}

func (s *Sender) handleMessage(m message) {
	port, ok := s.ports[m.port]
	if !ok {
		// TODO don't panic, log error
		panic(ErrPortNotFound)
	}

	// --- STAGGER: ensure we don't flood the same port ---
	if s.stagger > 0 {
		if last, ok := s.lastSend[m.port]; ok {
			elapsed := time.Since(last)
			if elapsed < s.stagger {
				time.Sleep(s.stagger - elapsed)
			}
		}
	}

	if err := port.send(m.msg); err != nil {
		// TODO don't panic, log error
		panic("MIDI Send error: " + err.Error())
	}

	s.lastSend[m.port] = time.Now()
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
		s.closePort(c.Port)
	default:
		// TODO don't panic, log error
		panic(ErrUnknownCommand)
	}
}

func (s *Sender) closePort(port int) {
	if p, ok := s.ports[port]; ok {
		_ = p.close()
		delete(s.ports, port)
	}
}
