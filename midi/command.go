package midi

type CommandId int

const (
	CmdOpenPort CommandId = iota
	CmdClosePort
)

type Command struct {
	Id   CommandId
	Port int
}

type ResultId int

type Result struct {
	Id      ResultId
	Message string
}
