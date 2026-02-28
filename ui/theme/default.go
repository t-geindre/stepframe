package theme

import (
	img "image"
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
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
		MainPanelTheme: basePanel.
			WithGradient(ColorMainGradientTop, ColorMainGradientBottom).
			WithSimplePadding(30),
		PanelTheme: basePanel.
			WithIdleBackgroundImage(ImagePanel),
		AddTrackPanelTheme: basePanel.
			WithIdleBackgroundImage(ImageAddTrackPanel),
		BarContainerPanelTheme: basePanel.
			WithIdleBackgroundImage(ImageBarContainerPanel),
		BarOffPanelTheme: basePanel.
			WithIdleBackgroundImage(ImageBarOffPanel).
			WithIdleBackgroundColorize(ColorBarOff).
			WithHoverBackgroundColorize(ColorBarOn).
			WithCursor(input.CURSOR_POINTER),
		BarOnPanelTheme: basePanel.
			WithIdleBackgroundImage(ImageBarOnPanel).
			WithIdleBackgroundColorize(ColorBarOn).
			WithHoverBackgroundColorize(ColorBarOff).
			WithCursor(input.CURSOR_POINTER),
		BarActivePanelTheme: basePanel.
			WithIdleBackgroundImage(ImageBarActivePanel).
			WithIdleBackgroundColorize(ColorBarActive).
			WithPulseBackgroundImage(ImageBarActivePanelPulse).
			WithPulseDuration(time.Millisecond * 300).
			WithPulseBackgroundColorize(ColorBarPulse).
			WithHoverBackgroundColorize(ColorBarOff).
			WithCursor(input.CURSOR_POINTER),
		LedPanelTheme: basePanel.
			WithIdleBackgroundImage(ImageLedPanel).
			WithPulseBackgroundImage(ImageLedPanelPulse).
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
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TilePanel),
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
			ColorIconDefault:        colornames.White,
			ColorLedIdle:            colornames.Red,
			ColorLedOn:              colornames.Lime,
			ColorLedArmed:           colornames.Yellow,
			ColorLedPulse:           colornames.White,
			ColorBarOff:             colornames.Dimgray,
			ColorBarOn:              color.RGBA{R: 0xff, G: 0x92, B: 0xff, A: 0xff},
			ColorBarActive:          colornames.Lime,
			ColorBarPulse:           colornames.White,
			ColorMainGradientTop:    color.RGBA{R: 0x22, G: 0x25, B: 0x2b, A: 255},
			ColorMainGradientBottom: color.RGBA{R: 0x0f, G: 0x11, B: 0x15, A: 255},
		},
		Images: map[Image]*OffsetImage{
			ImagePanel: NewSimpleOffsetImage(
				tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TilePanel),
				10, 27, 10,
				10, 27, 10,
			),
			ImagePanelHover: NewSimpleOffsetImage(
				tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TilePanel),
				10, 27, 10,
				10, 27, 10,
			),
			ImagePanelPressed: NewSimpleOffsetImage(
				tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButtonPressed),
				10, 27, 10,
				10, 27, 10,
			),
			ImageAddTrackPanel: NewSimpleOffsetImage(
				tileSheet.GetTile(TileVirtualPanel),
				10, 44, 10,
				10, 44, 10,
			),
			ImageBarContainerPanel: NewSimpleOffsetImage(
				tileSheet.GetTile(TileBarContainerPanel),
				10, 44, 10,
				10, 44, 10,
			),
			ImageBarOffPanel: NewOffsetImage(
				tileSheet.GetTile(TileBarOffPanel),
				21, 22, 21,
				21, 22, 21,
				16,
			),
			ImageBarOnPanel: NewOffsetImage(
				tileSheet.GetTile(TileBarPanel), 21, 22, 21, 21, 22, 21, 16,
			),
			ImageBarActivePanel: NewOffsetImage(
				tileSheet.GetTile(TileBarPanel), 21, 22, 21, 21, 22, 21, 16,
			),
			ImageBarActivePanelPulse: NewOffsetImage(
				tileSheet.GetTile(TileBarPanel), 21, 22, 21, 21, 22, 21, 16,
			),
			ImageLedPanel: NewOffsetImage(
				tileSheet.GetTile(TileLedOn), 0, 64, 0, 0, 64, 0, 23,
			),
			ImageLedPanelPulse: NewOffsetImage(
				tileSheet.GetTile(TileLedOn), 0, 64, 0, 0, 64, 0, 23,
			),
		},
		TrackTheme: &TrackTheme{
			Id: &TrackIdTheme{
				Font: getFontFace(18, true),
			},
		},
	}

	Current = theme
}
