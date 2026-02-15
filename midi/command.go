package midi

type CommandId int

const (
	CmdOpenPort CommandId = iota
	CmdClosePort
	CmdForward
)

type Command struct {
	Id      CommandId
	Port    int // int
	PortOut int
}

type ResultId int

type Result struct {
	Id      ResultId
	Message string
}
