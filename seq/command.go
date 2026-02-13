package seq

type CommandId int

const (
	CmdPlay CommandId = iota
	CmdStop
	CmdAdd
	CmdRemove
	CmdSwap
)

type CommandAt int

const (
	CmdAtNow = iota
	CmdAtNextBar
	CmdAtNextBeat
)

type Command struct {
	Id      CommandId
	At      CommandAt
	TrackId TrackId
	Track   *Track
}
