package midi

import (
	"context"
	"sync/atomic"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

const ForwardNone = -1

type Receiver struct {
	ports    map[int]func()
	done     chan struct{}
	running  atomic.Bool
	commands chan Command
	sender   *Sender
	forwards map[int]int // msg forwarding: in port -> out port
}

func NewReceiver(sd *Sender) *Receiver {
	return &Receiver{
		ports:    make(map[int]func()),
		commands: make(chan Command, 16),
		forwards: make(map[int]int),
		sender:   sd,
	}
}

func (r *Receiver) Run(ctx context.Context) {
	if !r.running.CompareAndSwap(false, true) {
		return
	}

	done := make(chan struct{})
	r.done = done

	go r.run(ctx, done)
}

func (r *Receiver) Wait() {
	done := r.done
	if done == nil {
		return
	}

	<-done

	for _, closePort := range r.ports {
		closePort()
	}
}

func (r *Receiver) Commands() chan<- Command {
	return r.commands
}

func (r *Receiver) run(ctx context.Context, done chan struct{}) {
	defer func() {
		close(done) // close local done
		r.running.Store(false)
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-r.commands:
			r.drainCommands(cmd)
		}
	}
}

func (r *Receiver) drainCommands(c Command) {
	r.handleCommand(c)
	for i := 0; i < 1024; i++ { // prevent starvation
		select {
		case c := <-r.commands:
			r.handleCommand(c)
		default:
			return
		}
	}
}

func (r *Receiver) handleCommand(c Command) {
	switch c.Id {
	case CmdOpenPort:
		if r.ports[c.Port] != nil {
			return // already open
		}

		p := midi.GetInPorts()
		if p[c.Port] == nil {
			// TODO don't panic, log error
			panic("Failed to open MIDI port: " + ErrPortNotFound.Error())
		}

		stop, err := midi.ListenTo(p[c.Port], func(msg midi.Message, ts int32) {
			r.onMessage(c.Port, msg, ts)
		})

		if err != nil {
			// TODO don't panic, log error
			panic("Failed to open MIDI port: " + err.Error())
		}

		r.ports[c.Port] = stop

	case CmdClosePort:
		if r.ports[c.Port] == nil {
			// TODO don't panic, log error
			panic("Failed to close MIDI port: " + ErrPortNotFound.Error())
		}

		r.ports[c.Port]()

	case CmdForward:
		if c.PortOut == ForwardNone {
			delete(r.forwards, c.Port)
			break
		}
		r.forwards[c.Port] = c.PortOut

	default:
		// TODO don't panic, log error
		panic(ErrUnknownCommand)
	}
}

func (r *Receiver) onMessage(p int, msg midi.Message, ts int32) {
	if fwdPort, ok := r.forwards[p]; ok {
		r.sender.SendRaw(fwdPort, msg)
	}
}
