package midi

import (
	"stepframe/seq/engine"

	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv" // autoregister driver
)

type Out struct {
	send func(midi.Message) error
}

func NewOut(port int) *Out {
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

func (o *Out) SendEvent(e engine.Event) error {
	switch e.Type {
	case engine.EvNoteOn:
		return o.send(midi.NoteOn(e.Channel, e.Note, e.Vel))
	case engine.EvNoteOff:
		return o.send(midi.NoteOff(e.Channel, e.Note))
	case engine.EvCC:
		return o.send(midi.ControlChange(e.Channel, e.CC, e.Value))
	default:
		return nil
	}
}

func (o *Out) SendClockPulse() error {
	return o.send(midi.TimingClock())
}

func (o *Out) SendStart() error {
	return o.send(midi.Start())
}

func (o *Out) SendStop() error {
	return o.send(midi.Stop())
}

func (o *Out) PanicAll() {
	for ch := uint8(0); ch < 16; ch++ {
		_ = o.send(midi.ControlChange(ch, 120, 0)) // All Sound Off
		_ = o.send(midi.ControlChange(ch, 123, 0)) // All Notes Off
	}
}
