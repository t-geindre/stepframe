package clock

import (
	"context"
)

type Tick struct {
	N    int64 // tick id
	When int64 // monotonic time in nanoseconds
}

type Clock interface {
	Ticks() <-chan Tick
	Run(ctx context.Context)
	SetBPM(bpm float64)
	BPM() float64
	Wait()
	GetTicksPerQuarter() int64
}
