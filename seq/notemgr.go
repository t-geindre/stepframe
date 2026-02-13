package seq

type NoteManager struct {
	sched *Scheduler
	held  [16][16][128]bool // port, channel, note
}

func NewNoteManager(s *Scheduler) *NoteManager { return &NoteManager{sched: s} }

// HandleNote Policy: retrigger steals previous voice (simple + safe)
func (nm *NoteManager) HandleNote(ev NoteEvent) {
	ch, note := ev.Channel, ev.Note

	// If already held, schedule an immediate off before the on (steal)
	if nm.held[ev.Port][ch][note] {
		nm.sched.Push(Event{
			AtTick:  ev.AtTick,
			Type:    EvNoteOff,
			Channel: ch,
			Port:    ev.Port,
			Note:    note,
		})
		nm.held[ev.Port][ch][note] = false
	}

	// NoteOn now
	nm.sched.Push(Event{
		AtTick:  ev.AtTick,
		Type:    EvNoteOn,
		Channel: ch,
		Port:    ev.Port,
		Note:    note,
		Vel:     ev.Velocity,
	})
	nm.held[ev.Port][ch][note] = true

	// NoteOff later
	nm.sched.Push(Event{
		AtTick:  ev.AtTick + ev.Duration,
		Type:    EvNoteOff,
		Channel: ch,
		Port:    ev.Port,
		Note:    note,
	})
}

func (nm *NoteManager) OnEventSent(e Event) {
	if e.Type == EvNoteOff {
		nm.held[e.Port][e.Channel][e.Note] = false
	}
}
