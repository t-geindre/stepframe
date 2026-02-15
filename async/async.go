package async

import (
	"sync/atomic"

	"github.com/rs/zerolog"
)

// Async is a helper struct to manage asynchronous Commands and events
type Async[C any, E any] struct {
	commands chan C
	events   chan E
	isDone   chan struct{}
	running  atomic.Bool
	closed   atomic.Bool
	logger   zerolog.Logger
}

func NewAsync[C any, E any](logger zerolog.Logger, buffer int) *Async[C, E] {
	return &Async[C, E]{
		commands: make(chan C, buffer),
		events:   make(chan E, buffer),
		isDone:   make(chan struct{}),
		logger:   logger,
	}
}

// Wait blocks until the Async is Done
func (a *Async[C, E]) Wait() {
	<-a.isDone
}

// TryCommand tries to send a command to the Async, drops if full
func (a *Async[C, E]) TryCommand(cmd C) bool {
	select {
	case a.commands <- cmd:
		return true
	default:
		a.logger.Warn().Msg("command dropped")
		return false
	}
}

// Events returns a channel that can be used to receive events from the Async
func (a *Async[C, E]) Events() <-chan E {
	return a.events
}

// TryDispatch tries to send an event to the Async listener
func (a *Async[C, E]) TryDispatch(event E) bool {
	select {
	case a.events <- event:
		return true
	default:
		a.logger.Warn().Msg("event dropped")
		return false
	}
}

// CanRun checks if the Async is not already running and sets it to running
func (a *Async[C, E]) CanRun() bool {
	return a.running.CompareAndSwap(false, true)
}

// Done marks the Async as Done and closes the isDone channel
func (a *Async[C, E]) Done() {
	if a.closed.CompareAndSwap(false, true) {
		close(a.isDone)
	}
}

// Commands returns the channel for receiving commands
func (a *Async[C, E]) Commands() <-chan C {
	return a.commands
}
