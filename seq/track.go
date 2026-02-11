package seq

type Step struct {
	On       bool
	Note     uint8
	Velocity uint8
	GateTick int64
}

type Track struct {
	Name         string
	Channel      uint8
	DivisionTick int64
	Steps        []Step

	stepIdx int
	nextAt  int64
}

func (t *Track) Reset(startTick int64) {
	t.stepIdx = 0
	t.nextAt = startTick
}

func (t *Track) ProcessTick(tick int64) []NoteEvent {
	if len(t.Steps) == 0 || t.DivisionTick <= 0 || tick < t.nextAt {
		return nil
	}
	var out []NoteEvent
	for tick >= t.nextAt {
		step := t.Steps[t.stepIdx%len(t.Steps)]
		if step.On {
			out = append(out, NoteEvent{
				AtTick:   t.nextAt,
				Channel:  t.Channel,
				Note:     step.Note,
				Velocity: step.Velocity,
				Duration: step.GateTick,
			})
		}
		t.stepIdx++
		t.nextAt += t.DivisionTick
	}
	return out
}
