package seq

type EventType uint8

const (
	EvNoteOn EventType = iota
	EvNoteOff
	EvCC
	EvPanic // optional high-level "all notes off" request
	EvClock
)

type Event struct {
	AtTick int64
	Type   EventType

	Channel uint8
	Port    int
	Note    uint8
	Vel     uint8 // for NoteOn
	CC      uint8 // for CC
	Value   uint8 // for CC
}

type NoteEvent struct {
	AtTick   int64
	Channel  uint8
	Port     int
	Note     uint8
	Velocity uint8
	Duration int64
}
