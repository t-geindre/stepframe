package seq

type CommandId int

const (
	CmdPlay CommandId = iota
	CmdStop
	CmdPause
)

type Command struct {
	Id CommandId
}

type EventId int

const (
	EvBeat EventId = iota
	EvPlaying
	EvStopped
	EvPaused
)

type Event struct {
	Id EventId
}
