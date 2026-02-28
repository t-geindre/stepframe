package theme

import (
	img "image"
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func SetDefaultTheme() {
	cText := color.White
	tileSheet := NewInternalSheet()
	fontFace := getFontFace(18, false)
	basePanel := NewPanelTheme().
		WithSpacing(10).
		WithPadding(&widget.Insets{Left: 10, Right: 10, Top: 8, Bottom: 8})

	theme := &Theme{
		MainContainerTheme: &MainContainerTheme{
			GradientTop:    color.RGBA{R: 0x22, G: 0x25, B: 0x2b, A: 255},
			GradientBottom: color.RGBA{R: 0x0f, G: 0x11, B: 0x15, A: 255},
			Padding:        widget.NewInsetsSimple(30),
			Spacing:        10,
		},
		PanelTheme: basePanel.
			WithBackgroundImage(
				NewSimpleOffsetImage(
					tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButton),
					10, 27, 10, 10, 27, 10,
				),
			),
		AddTrackPanelTheme: basePanel.
			WithBackgroundImage(
				NewSimpleOffsetImage(tileSheet.GetTile(TileVirtualPanel), 10, 44, 10, 10, 44, 10),
			),
		BarContainerPanelTheme: basePanel.
			WithBackgroundImage(
				NewSimpleOffsetImage(tileSheet.GetTile(TileBarContainerPanel), 10, 44, 10, 10, 44, 10),
			),
		BarOffPanelTheme: basePanel.
			WithBackgroundImage(NewOffsetImage(
				tileSheet.GetTile(TileBarOffPanel), 21, 22, 21, 21, 22, 21, 16,
			)).WithBackgroundColorize(ColorBarOff),
		BarOnPanelTheme: basePanel.
			WithBackgroundImage(NewOffsetImage(
				tileSheet.GetTile(TileBarPanel), 21, 22, 21, 21, 22, 21, 16,
			)).
			WithBackgroundColorize(ColorBarOn),
		BarActivePanelTheme: basePanel.
			WithBackgroundImage(NewOffsetImage(
				tileSheet.GetTile(TileBarPanel), 21, 22, 21, 21, 22, 21, 16,
			)).
			WithBackgroundColorize(ColorBarActive).
			WithPulseBackgroundImage(NewOffsetImage(
				tileSheet.GetTile(TileBarPanel), 21, 22, 21, 21, 22, 21, 16,
			)).
			WithPulseDuration(time.Millisecond * 300).
			WithPulseBackgroundColorize(ColorBarPulse),
		LedPanelTheme: basePanel.
			WithBackgroundImage(NewOffsetImage(
				tileSheet.GetTile(TileLedOn), 0, 64, 0, 0, 64, 0, 23,
			)).
			WithPulseBackgroundImage(NewOffsetImage(
				tileSheet.GetTile(TileLedOn), 0, 64, 0, 0, 64, 0, 23,
			)).
			WithPulseBackgroundColorize(ColorLedPulse).
			WithPulseDuration(time.Millisecond * 300),
		Theme: &widget.Theme{ // todo use a panel instead
			DefaultFace:      fontFace,
			DefaultTextColor: cText,
			ButtonTheme: &widget.ButtonParams{
				TextColor: &widget.ButtonTextColor{Idle: cText},
				TextFace:  fontFace,
				Image: &widget.ButtonImage{
					Idle: image.NewNineSlice(
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButton),
						[3]int{10, 27, 10},
						[3]int{10, 27, 10},
					),
					Hover: image.NewNineSlice(
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButtonHover),
						[3]int{10, 27, 10},
						[3]int{10, 27, 10},
					),
					Pressed: image.NewNineSlice(
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButtonPressed),
						[3]int{10, 27, 10},
						[3]int{10, 27, 10},
					),
				},
				TextPadding: &widget.Insets{Left: 20, Right: 20, Top: 10, Bottom: 10},
				TextPosition: &widget.TextPositioning{
					VTextPosition: widget.TextPositionCenter,
					HTextPosition: widget.TextPositionCenter,
				},
			},
			TextTheme: &widget.TextParams{
				Face:  fontFace,
				Color: cText,
				Position: &widget.TextPositioning{
					VTextPosition: widget.TextPositionCenter,
					HTextPosition: widget.TextPositionCenter,
				},
			},
		},
		IconsTheme: &IconsTheme{
			Icons: map[Icon]*ebiten.Image{
				IconBarDelete: tileSheet.GetTile(TileBarDelete),
				IconDelete:    tileSheet.GetTile(TileDelete),
				IconGear:      tileSheet.GetTile(TileGear),
				IconMinus:     tileSheet.GetTile(TileMinus),
				IconPause:     tileSheet.GetTile(TilePause),
				IconPlay:      tileSheet.GetTile(TilePlay),
				IconBarAdd:    tileSheet.GetTile(TileBarAdd),
				IconPlus:      tileSheet.GetTile(TilePlus),
				IconRecord:    tileSheet.GetTile(TileRecord),
				IconStop:      tileSheet.GetTile(TileStop),
			},
			Sizes: map[IconSize]int{
				IconSizeSmall:  16,
				IconSizeMedium: 24,
				IconSizeLarge:  32,
			},
		},
		Colors: ColorTheme{
			ColorIconDefault: colornames.White,
			ColorLedIdle:     colornames.Red,
			ColorLedOn:       colornames.Lime,
			ColorLedArmed:    colornames.Yellow,
			ColorLedPulse:    colornames.White,
			ColorBarOff:      colornames.Dimgray,
			ColorBarOn:       color.RGBA{R: 0xff, G: 0x92, B: 0xff, A: 0xff},
			ColorBarActive:   colornames.Lime,
			ColorBarPulse:    colornames.White,
		},
		TrackTheme: &TrackTheme{
			Id: &TrackIdTheme{
				Font: getFontFace(18, true),
			},
		},
	}

	Current = theme
}
