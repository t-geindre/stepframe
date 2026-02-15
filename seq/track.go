package seq

import (
	"sort"

	"gitlab.com/gomidi/midi/v2"
)

type Transformer interface {
	Transform(track *Track)
}

type TrackId int64

type Step struct {
	AtTick   int64 // LOCAL to track (0-based)
	Note     uint8
	Velocity uint8
	GateTick int64
	midi.Message
}

type Track struct {
	id       TrackId
	name     string
	channel  uint8
	port     int
	steps    []Step
	loop     bool
	loopTick int64
}

func NewTrack(name string) *Track {
	return &Track{
		name: name,
	}
}

func (t *Track) Id() TrackId { return t.id }

func (t *Track) SetChannel(ch uint8) {
	t.channel = ch
}

func (t *Track) SetPort(port int) {
	t.port = port
}

func (t *Track) SetId(id TrackId) {
	t.id = id
}

func (t *Track) Append(steps ...Step) {
	t.steps = append(t.steps, steps...)
}

// Finalize Always sort before committing changes
func (t *Track) Finalize(transformers ...Transformer) {
	for _, tr := range transformers {
		tr.Transform(t)
	}

	sort.Slice(t.steps, func(i, j int) bool {
		return t.steps[i].AtTick < t.steps[j].AtTick
	})
}

func (t *Track) SetLoop(loop bool, loopLen int64) {
	t.loop = loop
	t.loopTick = loopLen
}

func (t *Track) Clone() *Track {
	if t == nil {
		return nil
	}
	out := *t
	out.steps = append([]Step(nil), t.steps...) // deep copy
	return &out
}
