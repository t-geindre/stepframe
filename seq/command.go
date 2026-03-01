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
	CmdToggleBar
)

type Command struct {
	Id      CommandId
	TrackId *int
	Val     int
}

type EventId int

const (
	EvBeat EventId = iota
	EvPlaying
	EvStopped
	EvReset
	EvPaused
	EvRecording
	EvArmedPlaying
	EvArmedStopped
	EvArmedRecording
	EvTrackAdded
	EvTrackRemoved
	EvBarActivated
	EvBarDeactivated
)

type Event struct {
	Id      EventId
	TrackId *int
	Val     int
}
