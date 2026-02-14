package ui

import (
	"stepframe/clock"
	"stepframe/seq"
)

func getBillieJeanBassTrack() *seq.Track {
	t := seq.NewTrack("Bass")
	t.Append(
		seq.Step{AtTick: 0, Note: 64 - 12, Velocity: 110, GateTick: 24},   // F#
		seq.Step{AtTick: 48, Note: 59 - 12, Velocity: 110, GateTick: 24},  // C#
		seq.Step{AtTick: 96, Note: 62 - 12, Velocity: 110, GateTick: 24},  // E
		seq.Step{AtTick: 144, Note: 64 - 12, Velocity: 110, GateTick: 24}, // F#
		seq.Step{AtTick: 192, Note: 62 - 12, Velocity: 110, GateTick: 24}, // E
		seq.Step{AtTick: 240, Note: 59 - 12, Velocity: 110, GateTick: 24}, // C#
		seq.Step{AtTick: 288, Note: 57 - 12, Velocity: 110, GateTick: 24}, // B
		seq.Step{AtTick: 336, Note: 59 - 12, Velocity: 110, GateTick: 24}, // C#
	)
	t.SetLoop(true, 384)
	t.SetChannel(0)
	t.SetPort(1)
	t.Finalize()
	return t
}

func getBillieJeanLeadTrack() *seq.Track {
	t := seq.NewTrack("Lead")
	t.Append(
		seq.Step{AtTick: 0, Note: 64, Velocity: 110, GateTick: 144},   // E
		seq.Step{AtTick: 144, Note: 66, Velocity: 110, GateTick: 240}, // F#
		seq.Step{AtTick: 384, Note: 67, Velocity: 110, GateTick: 144}, // G
		seq.Step{AtTick: 528, Note: 66, Velocity: 110, GateTick: 240}, // F#
	)
	t.SetLoop(true, 384*2)
	t.SetChannel(1)
	t.SetPort(1)
	t.Finalize()
	return t
}

func getBillieJeanLeadTrackWithRatchet(clk clock.Clock) *seq.Track {
	t := getBillieJeanLeadTrack()
	rat := seq.NewRatchet(clk, 1)
	rat.Intervals = []int{0, 0, 7}
	t.Finalize(rat)
	return t
}

func getBillieJeanLeadTrackWithRatchetDouble(clk clock.Clock) *seq.Track {
	t := getBillieJeanLeadTrack()
	rat := seq.NewRatchet(clk, 2)
	rat.Intervals = []int{0}
	t.Finalize(rat)
	return t
}
