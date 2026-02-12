package main

import (
	"stepframe/seq"
)

func getBillieJeanBassTrack() *seq.Track {
	t := seq.NewTrack("Bass (Billie Jean - played version)")
	t.Append(
		seq.Step{AtTick: 0, Note: 64 - 12, Velocity: 110, GateTick: 36},   // F#
		seq.Step{AtTick: 48, Note: 59 - 12, Velocity: 110, GateTick: 36},  // C#
		seq.Step{AtTick: 96, Note: 62 - 12, Velocity: 110, GateTick: 36},  // E
		seq.Step{AtTick: 144, Note: 64 - 12, Velocity: 110, GateTick: 36}, // F#
		seq.Step{AtTick: 192, Note: 62 - 12, Velocity: 110, GateTick: 36}, // E
		seq.Step{AtTick: 240, Note: 59 - 12, Velocity: 110, GateTick: 36}, // C#
		seq.Step{AtTick: 288, Note: 57 - 12, Velocity: 110, GateTick: 36}, // B
		seq.Step{AtTick: 336, Note: 59 - 12, Velocity: 110, GateTick: 36}, // C#
	)
	t.SetLoop(true, 384)
	t.Sort()
	return t
}

func getBillieJeanLeadTrack() *seq.Track {
	t := seq.NewTrack("Lead (Billie Jean - played version)")
	t.Append(
		seq.Step{AtTick: 0, Note: 64, Velocity: 110, GateTick: 144},   // E
		seq.Step{AtTick: 144, Note: 66, Velocity: 110, GateTick: 240}, // F#
		seq.Step{AtTick: 384, Note: 67, Velocity: 110, GateTick: 144}, // G
		seq.Step{AtTick: 528, Note: 66, Velocity: 110, GateTick: 240}, // F#
	)
	t.SetLoop(true, 384*2)
	t.Sort()
	return t
}
