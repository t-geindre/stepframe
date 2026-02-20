package seq

import "gitlab.com/gomidi/midi/v2"

type NoteTracker struct {
	notesOn [16][128]bool
	empty   bool
}

func NewNoteTracker() *NoteTracker {
	return &NoteTracker{
		empty: true,
	}
}

func (nt *NoteTracker) Track(msg midi.Message) {
	key, ch, vel := uint8(0), uint8(0), uint8(0)

	switch {
	case msg.GetNoteOn(&ch, &key, &vel):
		nt.notesOn[ch][key] = vel > 0
		nt.empty = false
	case msg.GetNoteOff(&ch, &key, &vel):
		nt.notesOn[ch][key] = false
	}
}

func (nt *NoteTracker) Flush() []midi.Message {
	if nt.empty {
		return nil
	}

	var msgs []midi.Message

	for ch := 0; ch < 16; ch++ {
		for key := 0; key < 128; key++ {
			if nt.notesOn[ch][key] {
				msgs = append(msgs, midi.NoteOff(uint8(ch), uint8(key)))
				nt.notesOn[ch][key] = false
			}
		}
	}

	return msgs
}
