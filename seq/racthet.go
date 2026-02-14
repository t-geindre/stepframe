package seq

import (
	"sort"
	"stepframe/clock"
)

type Ratchet struct {
	Clock           clock.Clock
	NotesPerQuarter int64

	// If >0, enforce a fixed number of sub-notes per note (count is applied per original Step).
	// When Count > 0, the transformer will generate exactly Count sub-notes (unless clamped by nextAt/loop),
	// and GateTick is used only for clamping (not for determining Count).
	Count int

	// Intervals is a semitone offset pattern applied per generated sub-note.
	// Example: []int{0, 7, 12, 19} => root, fifth, octave, octave+fifth.
	// If empty or nil, defaults to []int{0} (classic ratchet: repeats the same note).
	Intervals []int

	GateRatio float64

	// MinGate avoids too short gates
	MinGate int64
}

func NewRatchet(c clock.Clock, notesPerQuarter int64) *Ratchet {
	return &Ratchet{
		Clock:           c,
		NotesPerQuarter: notesPerQuarter,
		Intervals:       []int{0},
		GateRatio:       0.9,
		MinGate:         1,
	}
}

func (r *Ratchet) Transform(t *Track) {
	if t == nil || len(t.steps) == 0 || r.Clock == nil {
		return
	}
	tpq := r.Clock.GetTicksPerQuarter()
	if tpq <= 0 {
		return
	}

	if r.NotesPerQuarter <= 0 {
		r.NotesPerQuarter = 4 // default: 1/16
	}
	interval := tpq / r.NotesPerQuarter
	if interval <= 0 {
		interval = 1
	}

	if len(r.Intervals) == 0 {
		r.Intervals = []int{0}
	}

	if r.GateRatio <= 0 {
		r.GateRatio = 1.0
	}
	if r.MinGate <= 0 {
		r.MinGate = 1
	}

	orig := append([]Step(nil), t.steps...)
	sort.Slice(orig, func(i, j int) bool { return orig[i].AtTick < orig[j].AtTick })

	out := make([]Step, 0, len(orig)*2)

	for i := 0; i < len(orig); i++ {
		s := orig[i]

		nextAt := int64(1 << 62)
		if i+1 < len(orig) {
			nextAt = orig[i+1].AtTick
		}
		if t.loop && t.loopTick > 0 && nextAt > t.loopTick {
			nextAt = t.loopTick
		}

		if s.GateTick <= 0 {
			out = append(out, s)
			continue
		}

		// Clamp the available duration so sub-notes never overlap the next event (or loop end).
		maxDur := s.GateTick
		if s.AtTick+maxDur > nextAt {
			maxDur = nextAt - s.AtTick
		}
		if maxDur <= 0 {
			continue
		}

		count := r.Count
		if count <= 0 {
			// Auto-count: generate as many sub-notes as can fit in maxDur.
			count = int(maxDur / interval)
			if maxDur%interval != 0 {
				count++
			}
			if count < 1 {
				count = 1
			}
		}

		stepGate := int64(float64(interval) * r.GateRatio)
		if stepGate < r.MinGate {
			stepGate = r.MinGate
		}

		remain := maxDur
		for k := 0; k < count && remain > 0; k++ {
			ns := s
			ns.AtTick = s.AtTick + int64(k)*interval

			// Apply semitone offset pattern.
			semi := r.Intervals[k%len(r.Intervals)]
			n := int(ns.Note) + semi
			if n < 0 || n > 127 {
				// Skip out-of-range MIDI notes, but still advance time.
				remain -= interval
				continue
			}
			ns.Note = uint8(n)

			g := stepGate
			if g > remain {
				g = remain
			}
			if g < r.MinGate {
				g = r.MinGate
			}
			ns.GateTick = g

			out = append(out, ns)
			remain -= interval
		}
	}

	t.steps = out
}
