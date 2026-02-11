package engine

import "stepframe/seq"

type NoteManager struct {
	sched *Scheduler
	held  [16][128]bool
}

func NewNoteManager(s *Scheduler) *NoteManager { return &NoteManager{sched: s} }

// HandleNote Policy: retrigger steals previous voice (simple + safe)
func (nm *NoteManager) HandleNote(ev seq.NoteEvent) {
	ch, note := ev.Channel, ev.Note

	// If already held, schedule an immediate off before the on (steal)
	if nm.held[ch][note] {
		nm.sched.Push(seq.Event{
			AtTick:  ev.AtTick,
			Type:    seq.EvNoteOff,
			Channel: ch,
			Note:    note,
		})
		nm.held[ch][note] = false
	}

	// NoteOn now
	nm.sched.Push(seq.Event{
		AtTick:  ev.AtTick,
		Type:    seq.EvNoteOn,
		Channel: ch,
		Note:    note,
		Vel:     ev.Velocity,
	})
	nm.held[ch][note] = true

	// NoteOff later
	nm.sched.Push(seq.Event{
		AtTick:  ev.AtTick + ev.Duration,
		Type:    seq.EvNoteOff,
		Channel: ch,
		Note:    note,
	})
}

func (nm *NoteManager) OnEventSent(e seq.Event) {
	if e.Type == seq.EvNoteOff {
		nm.held[e.Channel][e.Note] = false
	}
}
