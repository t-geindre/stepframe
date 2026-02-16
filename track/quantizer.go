package track

import "gitlab.com/gomidi/midi/v2"

type QuantizeMode int

const (
	QuantizeNearest QuantizeMode = iota
	QuantizeFloor
	QuantizeCeil
)

type Quantizer struct {
	Inner    Track
	GridTick int64
	Mode     QuantizeMode
}

func NewQuantizer(inner Track, gridTick int64, mode QuantizeMode) Track {
	if inner == nil {
		panic("inner is nil")
	}
	if gridTick <= 0 {
		panic("gridTick must be > 0")
	}
	return &Quantizer{Inner: inner, GridTick: gridTick, Mode: mode}
}

func (q *Quantizer) Reset(globalTick0 int64)          { q.Inner.Reset(globalTick0) }
func (q *Quantizer) GetBaseTick() int64               { return q.Inner.GetBaseTick() }
func (q *Quantizer) LocalTick(globalTick int64) int64 { return q.Inner.LocalTick(globalTick) }
func (q *Quantizer) PollDue(globalTick int64) []midi.Message {
	return q.Inner.PollDue(globalTick)
}

// IMPORTANT: atTick is GLOBAL here (by your new convention)
func (q *Quantizer) AddEvent(globalTick int64, msg midi.Message) {
	if msg == nil {
		return
	}

	base := q.Inner.GetBaseTick()

	// tick relative to track start (can be negative)
	rel := globalTick - base

	// quantize relative tick on the grid
	relQ := quantizeTick(rel, q.GridTick, q.Mode)

	// convert back to a global tick aligned on base
	globalQ := base + relQ

	// delegate: inner will modulo to its loop length
	q.Inner.AddEvent(globalQ, msg)
}

func quantizeTick(tick, grid int64, mode QuantizeMode) int64 {
	r := tick % grid
	if r < 0 {
		r += grid
	}
	base := tick - r

	switch mode {
	case QuantizeFloor:
		return base
	case QuantizeCeil:
		if r == 0 {
			return tick
		}
		return base + grid
	default: // QuantizeNearest
		if r*2 < grid {
			return base
		}
		return base + grid
	}
}
