package track

import "gitlab.com/gomidi/midi/v2"

type QuantizeMode int

const (
	QuantizeNearest QuantizeMode = iota
	QuantizeFloor
	QuantizeCeil
)

type Quantizer struct {
	inner    Track
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
	return &Quantizer{inner: inner, gridTick: gridTick, Mode: mode}
}

func (q *Quantizer) AddEvent(atLocalTick int64, msg midi.Message) {
	base := q.inner.GetBaseLocalTick()
	rel := atLocalTick - base
	relQ := quantizeTick(rel, q.gridTick, q.Mode)
	localQ := base + relQ
	q.inner.AddEvent(localQ, msg)
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

// inner forwards

func (q *Quantizer) Reset()                                { q.inner.Reset() }
func (q *Quantizer) SetBaseLocalTick(baseLocal int64)      { q.inner.SetBaseLocalTick(baseLocal) }
func (q *Quantizer) GetBaseLocalTick() int64               { return q.inner.GetBaseLocalTick() }
func (q *Quantizer) LocalTick(nowLocal int64) int64        { return q.inner.LocalTick(nowLocal) }
func (q *Quantizer) PollDue(nowLocal int64) []midi.Message { return q.inner.PollDue(nowLocal) }
func (q *Quantizer) GetLengthTick() int64                  { return q.inner.GetLengthTick() }
