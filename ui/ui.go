package ui

import (
	"stepframe/seq"
	"stepframe/ui/layout"
	"stepframe/ui/theme"
	"strconv"

	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/widget"
)

func New(sqr *seq.Sequencer) *ebitenui.UI {
	th := theme.NewDefaultTheme()

	playing := make(map[seq.TrackId]bool)

	grid := layout.NewPlayGrid(th)
	for i := 0; i < 16; i++ {
		grid.AddChild(widget.NewButton(
			widget.ButtonOpts.TextLabel("Play "+strconv.Itoa(i)),
			widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
				trId := seq.TrackId(i)
				if playing[trId] {
					sqr.Commands() <- seq.Command{
						Id:      seq.CmdStop,
						TrackId: seq.TrackId(i),
						At:      seq.CmdAtNextBar,
					}
					args.Button.SetText("Play " + strconv.Itoa(i))
					playing[trId] = false
					return
				}
				sqr.Commands() <- seq.Command{
					Id:      seq.CmdPlay,
					TrackId: seq.TrackId(i),
					At:      seq.CmdAtNextBar,
				}
				args.Button.SetText("Stop " + strconv.Itoa(i))
				playing[trId] = true
			}),
		))
	}

	main := layout.NewMain(th)
	main.AddChild(grid)

	return &ebitenui.UI{
		Container:    main,
		PrimaryTheme: th.Theme,
	}
}
