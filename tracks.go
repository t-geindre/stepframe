package main

import "stepframe/seq"

func getBillieJeanBassTrack() *seq.Track {
	return &seq.Track{
		Name:         "Bass (Billie Jean - played version)",
		Channel:      0,
		DivisionTick: 24, // 1/16 todo
		Steps: []seq.Step{
			{On: true, Note: 64 - 12, Velocity: 110, GateTick: 32}, // E
			{On: false},
			{On: true, Note: 59 - 12, Velocity: 110, GateTick: 32}, // B (short)
			{On: false},

			{On: true, Note: 62 - 12, Velocity: 110, GateTick: 32}, // D (short)
			{On: false},

			{On: true, Note: 64 - 12, Velocity: 110, GateTick: 32}, // E
			{On: false},

			{On: true, Note: 62 - 12, Velocity: 110, GateTick: 32}, // D
			{On: false},

			{On: true, Note: 59 - 12, Velocity: 110, GateTick: 32}, // B
			{On: false},

			{On: true, Note: 57 - 12, Velocity: 110, GateTick: 32}, // A
			{On: false},

			{On: true, Note: 59 - 12, Velocity: 110, GateTick: 32}, // B
			{On: false},
		},
	}
}
