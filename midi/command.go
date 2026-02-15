package midi

import "gitlab.com/gomidi/midi/v2"

type CommandId int

const (
	CmdOpenPort CommandId = iota
	CmdClosePort
	CmdMessage
)

type Command struct {
	Id   CommandId
	Port int // int
	Msg  midi.Message
}

type EventId int

const (
	EvMessage EventId = iota
	EvPortOpened
	EvPortClosed
)

type Event struct {
	Id   EventId
	Msg  midi.Message // for EvMessage
	Port int          // for EvPortOpened and EvPortClosed
}
