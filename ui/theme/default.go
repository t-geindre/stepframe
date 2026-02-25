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
	face := getFontFace(18, false)

	theme := &Theme{
		BackgroundTheme: &BackgroundTheme{
			GradientTop:    color.RGBA{R: 0x22, G: 0x25, B: 0x2b, A: 255},
			GradientBottom: color.RGBA{R: 0x0f, G: 0x11, B: 0x15, A: 255},
			Padding:        widget.NewInsetsSimple(30),
		},
		PanelTheme: &PanelTheme{
			BackgroundImage: image.NewNineSlice(
				tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButton),
				[3]int{23, 1, 23},
				[3]int{23, 1, 23},
			),
			Padding: &widget.Insets{Left: 10, Right: 10, Top: 8, Bottom: 8},
			Spacing: 10,
		},
		Theme: &widget.Theme{
			DefaultFace:      face,
			DefaultTextColor: cText,
			ButtonTheme: &widget.ButtonParams{
				TextColor: &widget.ButtonTextColor{Idle: cText},
				TextFace:  face,
				Image: &widget.ButtonImage{
					Idle: image.NewNineSlice(
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButton),
						[3]int{23, 1, 23},
						[3]int{23, 1, 23},
					),
					Hover: image.NewNineSlice(
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButtonHover),
						[3]int{23, 1, 23},
						[3]int{23, 1, 23},
					),
					Pressed: image.NewNineSlice(
						tileSheet.Crop(img.Rect(9, 9, 56, 56)).GetTile(TileButtonPressed),
						[3]int{23, 1, 23},
						[3]int{23, 1, 23},
					),
				},
				TextPadding: &widget.Insets{Left: 15, Right: 15, Top: 5, Bottom: 5},
				TextPosition: &widget.TextPositioning{
					VTextPosition: widget.TextPositionCenter,
					HTextPosition: widget.TextPositionCenter,
				},
			},
			TextTheme: &widget.TextParams{
				Face:  face,
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
			IconColorDefault: colornames.White,
			ColorIdle:        colornames.Red,
			ColorOn:          colornames.Lime,
			ColorOff:         colornames.Black,
			ColorArmed:       colornames.Yellow,
		},
		LedTheme: &LedTheme{
			OnImage:       tileSheet.GetTile(TileLedOn),
			OffImage:      tileSheet.GetTile(TileLedOff),
			Width:         16,
			Height:        16,
			PulseDuration: 350 * time.Millisecond,
		},
		TrackTheme: &TrackTheme{
			Id: &TrackIdTheme{
				Font: getFontFace(18, true),
			},
		},
	}

	Current = theme
}
