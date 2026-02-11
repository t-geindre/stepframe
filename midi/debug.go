package midi

import "gitlab.com/gomidi/midi/v2"

func DebugPorts() {
	outs := midi.GetOutPorts()
	for i, p := range outs {
		println("Out", i, ":", p.String())
	}

}
