package seq

import "sort"

type TrackState struct {
	track    *Track
	baseTick int64
	cursor   int
	curCycle int64

	play  bool
	muted bool

	playAt int64
	stopAt int64
}

func NewTrackState(tr *Track) *TrackState {
	return &TrackState{track: tr}
}

func (t *TrackState) Swap(tr *Track, now int64) {
	local := now - t.baseTick

	if t.track.loop && t.track.loopTick > 0 {
		local = local % t.track.loopTick
	}

	t.track = tr

	t.cursor = sort.Search(len(t.track.steps), func(i int) bool {
		return t.track.steps[i].AtTick > local
	})
}

func (t *TrackState) Reset(startTick int64) {
	t.baseTick = startTick
	t.cursor = 0
	t.curCycle = 0
}

func (t *TrackState) ProcessTick(tick int64, out []NoteEvent) []NoteEvent {
	tr := t.track

	if len(tr.steps) == 0 || tick < t.baseTick {
		return nil
	}

	if !tr.loop && t.cursor >= len(tr.steps) {
		return nil
	}

	local := tick - t.baseTick
	baseCycleTick := t.baseTick

	if tr.loop {
		if tr.loopTick <= 0 {
			return nil
		}

		cycle := local / tr.loopTick
		local = local % tr.loopTick

		if cycle != t.curCycle {
			t.curCycle = cycle
			t.cursor = 0
		}

		baseCycleTick = t.baseTick + cycle*tr.loopTick
	}

	if t.cursor >= len(tr.steps) {
		return out
	}

	for t.cursor < len(tr.steps) {
		step := tr.steps[t.cursor]
		if step.AtTick > local {
			break
		}

		out = append(out, NoteEvent{
			AtTick:   baseCycleTick + step.AtTick,
			Channel:  tr.channel,
			Port:     tr.port,
			Note:     step.Note,
			Velocity: step.Velocity,
			Duration: step.GateTick,
		})

		t.cursor++
	}

	return out
}
