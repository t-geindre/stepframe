package track

import "gitlab.com/gomidi/midi/v2"

type QuantizeMode int

const (
	QuantizeNearest QuantizeMode = iota
	QuantizeFloor
	QuantizeCeil
)

type Quantizer struct {
	Track
	gridTick int64
	Mode     QuantizeMode
}

func NewQuantizer(inner Track, gridTick int64, mode QuantizeMode) *Quantizer {
	if inner == nil {
		panic("inner is nil")
	}
	if gridTick <= 0 {
		panic("gridTick must be > 0")
	}
	return &Quantizer{Track: inner, gridTick: gridTick, Mode: mode}
}

func (q *Quantizer) AddEvent(atLocalTick int64, msg midi.Message) {
	base := q.Track.GetBaseLocalTick()
	rel := atLocalTick - base
	relQ := quantizeTick(rel, q.gridTick, q.Mode)
	localQ := base + relQ
	q.Track.AddEvent(localQ, msg)
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
