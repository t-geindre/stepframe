package track

import (
	"sort"

	"gitlab.com/gomidi/midi/v2"
)

// Track works in sequencer-local ticks.
type Track interface {
	Reset()
	// SetBaseLocalTick sets the base local tick of the track (usually 0).
	SetBaseLocalTick(baseLocal int64)
	// LocalTick returns the local tick in [0..lengthTick) for a given nowLocal.
	LocalTick(nowLocal int64) int64
	// AddEvent receives a tick expressed in sequencer-local ticks.
	AddEvent(atLocalTick int64, msg midi.Message)
	// PollDue receives nowLocal (sequencer-local) and returns messages to fire now.
	PollDue(nowLocal int64) []midi.Message
	GetBaseLocalTick() int64
}

type event struct {
	AtTick int64 // local in [0..lengthTick)
	Msg    midi.Message
}

type track struct {
	lengthTick int64
	baseTick   int64 // base local tick (sequencer-local) where track local tick == 0

	events []event // sorted by AtTick
	cursor int     // index of first event NOT YET passed in current cycle
}

func NewTrack(lengthTick int64) Track {
	if lengthTick <= 0 {
		panic("lengthTick must be > 0")
	}
	return &track{
		lengthTick: lengthTick,
		baseTick:   0,
		events:     nil,
		cursor:     0,
	}
}

func (t *track) GetBaseLocalTick() int64 {
	return t.baseTick
}

func (t *track) SetBaseLocalTick(baseLocal int64) {
	t.baseTick = baseLocal
}

func (t *track) Reset() {
	t.cursor = 0
}

func (t *track) LocalTick(nowLocal int64) int64 {
	local := nowLocal - t.baseTick
	local %= t.lengthTick
	if local < 0 {
		local += t.lengthTick
	}
	return local
}

func (t *track) AddEvent(atLocalTick int64, msg midi.Message) {
	atTick := t.LocalTick(atLocalTick)

	ev := event{AtTick: atTick, Msg: msg}

	i := sort.Search(len(t.events), func(i int) bool {
		return t.events[i].AtTick > ev.AtTick
	})

	t.events = append(t.events, event{})
	copy(t.events[i+1:], t.events[i:])
	t.events[i] = ev

	if i <= t.cursor {
		t.cursor++
	}
}

func (t *track) PollDue(nowLocal int64) []midi.Message {
	if len(t.events) == 0 {
		return nil
	}

	local := t.LocalTick(nowLocal)

	if local == 0 {
		t.cursor = 0
	}

	if t.cursor < len(t.events) && t.events[t.cursor].AtTick < local {
		t.cursor = sort.Search(len(t.events), func(i int) bool {
			return t.events[i].AtTick >= local
		})
	}

	start := t.cursor
	for t.cursor < len(t.events) && t.events[t.cursor].AtTick == local {
		t.cursor++
	}

	if t.cursor == start {
		return nil
	}

	out := make([]midi.Message, 0, t.cursor-start)
	for i := start; i < t.cursor; i++ {
		out = append(out, t.events[i].Msg)
	}
	return out
}
