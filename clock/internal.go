package clock

import (
	"context"
	"math"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

type Internal struct {
	done    chan struct{} // clock is done
	running atomic.Bool

	logger zerolog.Logger

	ppqn     int64
	bpm      atomic.Uint64 // lock-free (float64 bits)
	ticks    chan Tick
	stopOnce sync.Once

	dropTicks bool
}

func NewInternalClock(logger zerolog.Logger, ppqn int64, bpm float64, buffer int, dropTicks bool) *Internal {
	c := &Internal{
		logger:    logger.With().Str("component", "internal_clock").Logger(),
		ppqn:      ppqn,
		ticks:     make(chan Tick, buffer),
		done:      make(chan struct{}),
		dropTicks: dropTicks,
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

func (c *Internal) tickDurationFromBits(bpmBits uint64) time.Duration {
	bpm := math.Float64frombits(bpmBits)
	if !(bpm > 0) || math.IsNaN(bpm) || math.IsInf(bpm, 0) {
		bpm = 120
	}
	if c.ppqn <= 0 {
		return time.Second
	}
	tps := (bpm / 60.0) * float64(c.ppqn)
	if !(tps > 0) || math.IsNaN(tps) || math.IsInf(tps, 0) {
		return time.Second
	}
	secPerTick := 1.0 / tps
	d := time.Duration(secPerTick * float64(time.Second))
	if d <= 0 {
		return time.Nanosecond
	}
	return d
}

func sleepUntilCtx(ctx context.Context, deadline time.Time) bool {
	const spinThreshold = 300 * time.Microsecond

	for {
		now := time.Now()
		remain := deadline.Sub(now)
		if remain <= 0 {
			return true
		}
		if remain > spinThreshold {
			t := time.NewTimer(remain - spinThreshold)
			select {
			case <-ctx.Done():
				t.Stop()
				return false
			case <-t.C:
			}
			continue
		}
		for time.Now().Before(deadline) {
			select {
			case <-ctx.Done():
				return false
			default:
				runtime.Gosched()
			}
		}
		return true
	}
}

func (c *Internal) sendTick(ctx context.Context, t Tick) bool {
	if c.dropTicks {
		select {
		case c.ticks <- t:
		default:
		}
		return true
	}

	select {
	case <-ctx.Done():
		return false
	case c.ticks <- t:
		return true
	}
}

func (c *Internal) run(ctx context.Context) {
	defer func() {
		close(c.done)
		c.logger.Info().Msg("stopped")
	}()

	c.logger.Info().Msg("running")

	start := time.Now()
	baseTime := start
	var sentTick int64 = -1

	lastBpmBits := c.bpm.Load()
	tickDur := c.tickDurationFromBits(lastBpmBits)

	for {
		select {
		case <-ctx.Done():
			close(c.ticks)
			return
		default:
		}

		now := time.Now()
		elapsed := now.Sub(baseTime)
		if elapsed < 0 {
			elapsed = 0
		}

		curTick := int64(elapsed / tickDur)

		bpmBits := c.bpm.Load()
		if bpmBits != lastBpmBits {
			lastBpmBits = bpmBits
			tickDur = c.tickDurationFromBits(bpmBits)
			baseTime = now.Add(-time.Duration(curTick) * tickDur)
		}

		nextTick := curTick + 1
		next := baseTime.Add(time.Duration(nextTick) * tickDur)

		if !sleepUntilCtx(ctx, next) {
			close(c.ticks)
			return
		}

		emitNow := time.Now()
		emitTick := int64(emitNow.Sub(baseTime) / tickDur)
		if emitTick < nextTick {
			emitTick = nextTick
		}

		for sentTick < emitTick {
			sentTick++
			if !c.sendTick(ctx, Tick{N: sentTick, When: emitNow.UnixNano()}) {
				close(c.ticks)
				return
			}
		}

		// Drift catch-up
	}
}
