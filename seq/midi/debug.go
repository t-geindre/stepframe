package midi

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
)

func DebugPorts() {
	outs := midi.GetOutPorts()
	for i, p := range outs {
		fmt.Println("Out", i, ":", p.String())
	}

}
