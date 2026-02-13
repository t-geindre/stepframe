package midi

import (
	"gitlab.com/gomidi/midi/v2"
)

type Port struct {
	Id    int
	Name  string
	send  func(midi.Message) error
	close func() error
}

func AllPorts() []Port {
	ports := make([]Port, 0)
	outs := midi.GetOutPorts()
	for i, p := range outs {
		ports = append(ports, Port{
			Id:   i,
			Name: p.String(),
		})
	}

	return ports
}
