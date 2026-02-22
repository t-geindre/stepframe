package track

import "gitlab.com/gomidi/midi/v2"

type Velocity struct {
	Track
	min, max uint8
}

func NewVelocity(track Track) *Velocity {
	return &Velocity{Track: track, min: 0, max: 127}
}

func (v *Velocity) SetMinMax(min, max uint8) {
	v.min, v.max = min, max
}

func (v *Velocity) PollDue(nowLocal int64) []midi.Message {
	var ch, key, vel uint8
	msgs := v.Track.PollDue(nowLocal)

	for i, msg := range msgs {
		switch {
		case msg.GetNoteOn(&ch, &key, &vel):
			if vel == 0 {
				// Note on with velocity 0 is a note off
				continue
			}
			if vel < v.min {
				vel = v.min
				msgs[i] = midi.NoteOn(ch, key, vel)
			} else if vel > v.max {
				vel = v.max
				msgs[i] = midi.NoteOn(ch, key, vel)
			}
		}
	}

	return msgs
}
