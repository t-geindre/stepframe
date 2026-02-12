package main

import (
	"stepframe/seq/engine"
)

func getBillieJeanBassTrack() *engine.Track {
	t := engine.NewTrack("Bass (Billie Jean - played version)")
	t.Append(
		engine.Step{AtTick: 0, Note: 64 - 12, Velocity: 110, GateTick: 36},   // F#
		engine.Step{AtTick: 48, Note: 59 - 12, Velocity: 110, GateTick: 36},  // C#
		engine.Step{AtTick: 96, Note: 62 - 12, Velocity: 110, GateTick: 36},  // E
		engine.Step{AtTick: 144, Note: 64 - 12, Velocity: 110, GateTick: 36}, // F#
		engine.Step{AtTick: 192, Note: 62 - 12, Velocity: 110, GateTick: 36}, // E
		engine.Step{AtTick: 240, Note: 59 - 12, Velocity: 110, GateTick: 36}, // C#
		engine.Step{AtTick: 288, Note: 57 - 12, Velocity: 110, GateTick: 36}, // B
		engine.Step{AtTick: 336, Note: 59 - 12, Velocity: 110, GateTick: 36}, // C#
	)
	t.SetLoop(true, 384)
	t.Sort()
	return t
}

func getBillieJeanLeadTrack() *engine.Track {
	t := engine.NewTrack("Lead (Billie Jean - played version)")
	t.Append(
		engine.Step{AtTick: 0, Note: 64, Velocity: 110, GateTick: 144},   // E
		engine.Step{AtTick: 144, Note: 66, Velocity: 110, GateTick: 240}, // F#
		engine.Step{AtTick: 384, Note: 67, Velocity: 110, GateTick: 144}, // G
		engine.Step{AtTick: 528, Note: 66, Velocity: 110, GateTick: 240}, // F#
	)
	t.SetLoop(true, 384*2)
	t.Sort()
	return t
}
