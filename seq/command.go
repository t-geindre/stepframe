package seq

type CommandId int

const (
	CmdPlay CommandId = iota
	CmdStop
	CmdPause
	CmdRecord
	CmdStopRecord
	CmdNewTrack
	CmdRemoveTrack
)

type Command struct {
	Id      CommandId
	TrackId *int
}

type EventId int

const (
	EvBeat EventId = iota
	EvPlaying
	EvStopped
	EvPaused
	EvRecording
	EvArmedPlaying
	EvArmedStopped
	EvArmedPaused
	EvArmedRecording
	EvStopRecording
	EvTrackAdded
	EvTrackRemoved
)

type Event struct {
	Id      EventId
	TrackId *int
}
