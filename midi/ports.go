package midi

import (
	"gitlab.com/gomidi/midi/v2"
)

type Port struct {
	Id   int
	Name string
}

func AllOutPorts() []Port {
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

func AllInPorts() []Port {
	ports := make([]Port, 0)
	ins := midi.GetInPorts()
	for i, p := range ins {
		ports = append(ports, Port{
			Id:   i,
			Name: p.String(),
		})
	}

	return ports
}
