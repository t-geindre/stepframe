package seq

import "sort"

type Step struct {
	AtTick   int64 // LOCAL to track (0-based)
	Note     uint8
	Velocity uint8
	GateTick int64
}

type Track struct {
	name    string
	channel uint8
	steps   []Step

	baseTick int64
	cursor   int

	loop       bool
	lengthTick int64

	curCycle int64
}

func NewTrack(name string) *Track {
	return &Track{name: name}
}

func (t *Track) Reset(startTick int64) {
	t.baseTick = startTick
	t.cursor = 0
	t.curCycle = 0
}

func (t *Track) SetChannel(ch uint8) {
	t.channel = ch
}

func (t *Track) Append(steps ...Step) {
	t.steps = append(t.steps, steps...)
}

func (t *Track) Sort() {
	sort.Slice(t.steps, func(i, j int) bool {
		return t.steps[i].AtTick < t.steps[j].AtTick
	})
}

func (t *Track) SetLoop(loop bool, loopLen int64) {
	t.loop = loop
	t.lengthTick = loopLen
}

func (t *Track) ProcessTick(tick int64) []NoteEvent {
	// Rien à jouer ou pas encore démarré
	if len(t.steps) == 0 || tick < t.baseTick {
		return nil
	}

	// Si on a déjà consommé tous les steps et qu'on ne loop pas
	if !t.loop && t.cursor >= len(t.steps) {
		return nil
	}

	local := tick - t.baseTick
	baseCycleTick := t.baseTick

	if t.loop {
		if t.lengthTick <= 0 {
			return nil
		}

		cycle := local / t.lengthTick
		local = local % t.lengthTick

		// Nouveau cycle → on repart du début de la liste de steps
		if cycle != t.curCycle {
			t.curCycle = cycle
			t.cursor = 0
		}

		baseCycleTick = t.baseTick + cycle*t.lengthTick
	}

	if t.cursor >= len(t.steps) {
		return nil
	}

	var out []NoteEvent
	for t.cursor < len(t.steps) {
		step := t.steps[t.cursor]
		if step.AtTick > local {
			break
		}

		out = append(out, NoteEvent{
			AtTick:   baseCycleTick + step.AtTick,
			Channel:  t.channel,
			Note:     step.Note,
			Velocity: step.Velocity,
			Duration: step.GateTick,
		})

		t.cursor++
	}

	return out
}
