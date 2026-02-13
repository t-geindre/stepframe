package main

import (
	"stepframe/seq"
)

// TODO stagger + prio OFF

func getBillieJeanBassTrack() *seq.Track {
	t := seq.NewTrack(0, "Bass (Billie Jean - played version)")
	t.Append(
		seq.Step{AtTick: 1 + 0, Note: 64 - 12, Velocity: 110, GateTick: 24},   // F#
		seq.Step{AtTick: 1 + 48, Note: 59 - 12, Velocity: 110, GateTick: 24},  // C#
		seq.Step{AtTick: 1 + 96, Note: 62 - 12, Velocity: 110, GateTick: 24},  // E
		seq.Step{AtTick: 1 + 144, Note: 64 - 12, Velocity: 110, GateTick: 24}, // F#
		seq.Step{AtTick: 1 + 192, Note: 62 - 12, Velocity: 110, GateTick: 24}, // E
		seq.Step{AtTick: 1 + 240, Note: 59 - 12, Velocity: 110, GateTick: 24}, // C#
		seq.Step{AtTick: 1 + 288, Note: 57 - 12, Velocity: 110, GateTick: 24}, // B
		seq.Step{AtTick: 1 + 336, Note: 59 - 12, Velocity: 110, GateTick: 24}, // C#
	)
	t.SetLoop(true, 384)
	t.SetChannel(0)
	t.Sort()
	return t
}

func getBillieJeanLeadTrack() *seq.Track {
	t := seq.NewTrack(1, "Lead (Billie Jean - played version)")
	t.Append(
		seq.Step{AtTick: 0, Note: 64 - 12, Velocity: 110, GateTick: 144},   // E
		seq.Step{AtTick: 144, Note: 66 - 12, Velocity: 110, GateTick: 240}, // F#
		seq.Step{AtTick: 384, Note: 67 - 12, Velocity: 110, GateTick: 144}, // G
		seq.Step{AtTick: 528, Note: 66 - 12, Velocity: 110, GateTick: 240}, // F#
	)
	t.SetLoop(true, 384*2)
	t.SetChannel(1)
	t.Sort()
	return t
}
