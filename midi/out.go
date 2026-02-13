package midi

import (
	"errors"
	"stepframe/seq"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

var ErrUnknownEventType = errors.New("Unknown event type")

type Out struct {
	send func(midi.Message) error
}

func NewOut(port int) *Out {
	DebugPorts() // TODO REMOVE ME

	out, err := midi.OutPort(port)
	if err != nil {
		panic("OutPort error: " + err.Error())
	}

	send, err := midi.SendTo(out)
	if err != nil {
		panic("Sender error: " + err.Error())
	}

	return &Out{send: send}
}

func (o *Out) SendEvent(e seq.Event) {
	var err error

	switch e.Type {
	case seq.EvNoteOn:
		err = o.send(midi.NoteOn(e.Channel, e.Note, e.Vel))
		if e.Channel == 1 {
		}
	case seq.EvNoteOff:
		err = o.send(midi.NoteOff(e.Channel, e.Note))
	case seq.EvCC:
		err = o.send(midi.ControlChange(e.Channel, e.CC, e.Value))
	case seq.EvClock:
		err = o.send(midi.TimingClock())
	case seq.EvPanic:
		o.PanicAll()
	default:
		err = ErrUnknownEventType
	}

	if err != nil {
		panic("MIDI Send error: " + err.Error())
	}
}

func (o *Out) PanicAll() {
	for ch := uint8(0); ch < 16; ch++ {
		_ = o.send(midi.ControlChange(ch, 120, 0)) // All Sound Off
		_ = o.send(midi.ControlChange(ch, 123, 0)) // All Notes Off
	}
}
