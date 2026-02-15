package async

// Expose provides a limited interface to an Async, intended for use by external components
type Expose[C any, E any] struct {
	as *Async[C, E]
}

func NewExpose[C any, E any](async *Async[C, E]) *Expose[C, E] {
	return &Expose[C, E]{
		as: async,
	}
}

// Wait blocks until the Async has completed its work and is no longer running
func (e *Expose[C, E]) Wait() {
	e.as.Wait()
}

// TryCommand attempts to send a command to the Async
func (e *Expose[C, E]) TryCommand(cmd C) bool {
	return e.as.TryCommand(cmd)
}

// Events returns a read-only channel of events emitted by the Async
func (e *Expose[C, E]) Events() <-chan E {
	return e.as.Events()
}
