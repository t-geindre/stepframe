package clock

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"
)

// todo runtime.GOMAXPROCS(1) may reduce scheduler jitter (may affect GC performance, so test)

type Internal struct {
	ppqn     int64
	bpm      atomic.Uint64 // lock-free (float64 bits)
	ticks    chan Tick
	stopOnce sync.Once
	stopCh   chan struct{} // close to stop
	done     chan struct{} // clock is done
	running  atomic.Bool
}

func NewInternalClock(ppqn int64, bpm float64, buffer int) *Internal {
	c := &Internal{
		ppqn:   ppqn,
		ticks:  make(chan Tick, buffer),
		stopCh: make(chan struct{}),
		done:   make(chan struct{}),
	}
	c.SetBPM(bpm)
	return c
}

func (c *Internal) GetTicksPerQuarter() int64 {
	return c.ppqn
}

func (c *Internal) Ticks() <-chan Tick { return c.ticks }

func (c *Internal) Run(ctx context.Context) {
	if !c.running.CompareAndSwap(false, true) {
		return
	}

	go c.run(ctx)
}

// SetBPM thread-safe
func (c *Internal) SetBPM(bpm float64) {
	c.bpm.Store(math.Float64bits(bpm))
}

// BPM tread-safe
func (c *Internal) BPM() float64 {
	return math.Float64frombits(c.bpm.Load())
}

func (c *Internal) Wait() {
	<-c.done
}

func (c *Internal) tickDuration() time.Duration {
	bpm := c.BPM()
	tps := (bpm / 60.0) * float64(c.ppqn)
	secPerTick := 1.0 / tps
	return time.Duration(secPerTick * float64(time.Second))
}

func (c *Internal) run(ctx context.Context) {
	defer close(c.done)

	start := time.Now()
	var tick int64 = 0
	next := start

	for {
		select {
		case <-ctx.Done():
			close(c.ticks)
			return
		case <-c.stopCh:
			close(c.ticks)
			return
		default:
		}

		tickDur := c.tickDuration()
		if tick == 0 {
			next = start
		} else {
			next = next.Add(tickDur)
		}

		sleepUntil(next)

		now := time.Now()

		select {
		case c.ticks <- Tick{N: tick, When: now.UnixNano()}:
		default:
			// no consumer ready, skip tick
		}

		tick++

		// Drift catch-up
		ideal := int64(time.Since(start) / tickDur)
		if ideal-tick > 0 {
			tick = ideal
			next = start.Add(time.Duration(tick) * tickDur)
		}
	}
}
