package seq

import (
	"sort"
)

type Scheduler struct { // todo implement as a min-heap
	q []Event
}

func NewScheduler() *Scheduler { return &Scheduler{} }

func (s *Scheduler) Push(events ...Event) {
	s.q = append(s.q, events...)
	sort.SliceStable(s.q, func(i, j int) bool {
		if s.q[i].AtTick == s.q[j].AtTick {
			// Same tick: prioritize by event type
			return s.priority(s.q[i]) < s.priority(s.q[j])
		}

		return s.q[i].AtTick < s.q[j].AtTick
	})
}

// PopDue returns and removes all events with AtTick <= tick.
func (s *Scheduler) PopDue(tick int64) []Event {
	i := 0
	for i < len(s.q) && s.q[i].AtTick <= tick {
		i++
	}
	if i == 0 {
		return nil
	}
	due := append([]Event(nil), s.q[:i]...)
	s.q = s.q[i:]
	return due
}

func (s *Scheduler) Clear() { s.q = s.q[:0] }

func (s *Scheduler) priority(e Event) int {
	switch e.Type {
	case EvPanic:
		return 0
	case EvNoteOff:
		return 1
	case EvNoteOn:
		return 2
	case EvCC:
		return 3
	case EvClock:
		return 4
	default:
		panic("Unknown event type")
	}
}
