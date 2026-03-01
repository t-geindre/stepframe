package track

import (
	"sort"

	"gitlab.com/gomidi/midi/v2"
)

const MaxBars = 10

// Track works in sequencer-local ticks (nowLocal).
type Track interface {
	Reset()
	SetBaseLocalTick(baseLocal int64)
	LocalTick(nowLocal int64) int64
	AddEvent(atLocalTick int64, msg midi.Message)
	PollDue(nowLocal int64) []midi.Message
	GetBaseLocalTick() int64
	GetLengthTick() int64

	ActivateBar(index int)
	DeactivateBar(index int)
	IsBarActive(index int) bool
	GetMaxBars() int
}

type event struct {
	AtTick int64 // in [0..barLenTick)
	Msg    midi.Message
}

type bar struct {
	active bool
	events []event
	cursor int
}

type track struct {
	barLenTick int64
	baseTick   int64

	bars [MaxBars]bar

	// playhead (stateful)
	started    bool
	lastNow    int64
	virtualPos int64 // compressed ticks since base (monotonic, not mod)

	curBar int   // physical bar index [0..MaxBars)
	barPos int64 // tick in current bar [0..barLenTick)
}

func NewTrack(barLenTick int64) Track {
	if barLenTick <= 0 {
		panic("barLenTick must be > 0")
	}

	t := &track{
		barLenTick: barLenTick,
		baseTick:   0,
		curBar:     0,
		barPos:     0,
	}

	for i := 0; i < MaxBars; i++ {
		t.bars[i] = bar{active: false, events: nil, cursor: 0}
	}
	t.bars[0].active = true

	return t
}

func (t *track) GetMaxBars() int { return MaxBars }

func (t *track) GetBaseLocalTick() int64 { return t.baseTick }

func (t *track) SetBaseLocalTick(baseLocal int64) {
	t.baseTick = baseLocal
	t.started = false
	t.lastNow = 0
	t.virtualPos = 0
	t.curBar = 0
	t.barPos = 0
	for i := 0; i < MaxBars; i++ {
		t.bars[i].cursor = 0
	}
}

func (t *track) Reset() {
	t.started = false
	t.lastNow = 0
	t.virtualPos = 0
	t.curBar = t.firstActiveBarOrZero()
	t.barPos = 0
	for i := 0; i < MaxBars; i++ {
		t.bars[i].cursor = 0
	}
}

func (t *track) IsBarActive(index int) bool {
	if index < 0 || index >= MaxBars {
		return false
	}
	return t.bars[index].active
}

func (t *track) ActivateBar(index int) {
	if index < 0 || index >= MaxBars {
		return
	}
	t.bars[index].active = true
}

func (t *track) DeactivateBar(index int) {
	if index < 0 || index >= MaxBars {
		return
	}
	if !t.bars[index].active {
		return
	}

	activeCount := 0
	for i := 0; i < MaxBars; i++ {
		if t.bars[i].active {
			activeCount++
		}
	}
	if activeCount <= 1 {
		return
	}

	t.bars[index].active = false
}

func (t *track) GetLengthTick() int64 {
	count := 0
	for i := 0; i < MaxBars; i++ {
		if t.bars[i].active {
			count++
		}
	}
	if count <= 0 {
		count = 1
	}
	return int64(count) * t.barLenTick
}

// LocalTick is the compressed playhead position in [0..lengthTick).
// IMPORTANT: this is stateful: it advances the playhead to nowLocal.
func (t *track) LocalTick(nowLocal int64) int64 {
	t.advanceTo(nowLocal)
	l := t.GetLengthTick()
	if l <= 0 {
		return 0
	}
	return t.virtualPos % l
}

func (t *track) AddEvent(atLocalTick int64, msg midi.Message) {
	t.advanceTo(atLocalTick)

	b := &t.bars[t.curBar]
	ev := event{AtTick: t.barPos, Msg: msg}

	i := sort.Search(len(b.events), func(i int) bool {
		return b.events[i].AtTick > ev.AtTick
	})

	b.events = append(b.events, event{})
	copy(b.events[i+1:], b.events[i:])
	b.events[i] = ev

	if i <= b.cursor {
		b.cursor++
	}
}

func (t *track) PollDue(nowLocal int64) []midi.Message {
	t.advanceTo(nowLocal)

	b := &t.bars[t.curBar]
	if len(b.events) == 0 {
		return nil
	}

	local := t.barPos

	if b.cursor < len(b.events) && b.events[b.cursor].AtTick < local {
		b.cursor = sort.Search(len(b.events), func(i int) bool {
			return b.events[i].AtTick >= local
		})
	}

	start := b.cursor
	for b.cursor < len(b.events) && b.events[b.cursor].AtTick == local {
		b.cursor++
	}
	if b.cursor == start {
		return nil
	}

	out := make([]midi.Message, 0, b.cursor-start)
	for i := start; i < b.cursor; i++ {
		out = append(out, b.events[i].Msg)
	}
	return out
}

func (t *track) advanceTo(nowLocal int64) {
	if !t.started {
		t.started = true
		t.lastNow = t.baseTick
		t.virtualPos = 0
		t.curBar = t.firstActiveBarOrZero()
		t.barPos = 0
		t.bars[t.curBar].cursor = 0
	}

	if nowLocal <= t.lastNow {
		return
	}

	delta := nowLocal - t.lastNow

	for delta > 0 {
		remain := t.barLenTick - t.barPos
		step := delta
		if step > remain {
			step = remain
		}

		t.barPos += step
		t.virtualPos += step
		t.lastNow += step
		delta -= step

		if t.barPos == t.barLenTick {
			t.barPos = 0
			t.curBar = t.nextActiveBar(t.curBar)
			t.bars[t.curBar].cursor = 0
		}
	}
}

func (t *track) firstActiveBarOrZero() int {
	for i := 0; i < MaxBars; i++ {
		if t.bars[i].active {
			return i
		}
	}
	t.bars[0].active = true
	return 0
}

func (t *track) nextActiveBar(from int) int {
	for step := 1; step <= MaxBars; step++ {
		i := (from + step) % MaxBars
		if t.bars[i].active {
			return i
		}
	}
	t.bars[0].active = true
	return 0
}
