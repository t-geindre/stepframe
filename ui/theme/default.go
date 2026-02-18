package theme

import (
	img "image"
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/utilities/constantutil"
	"github.com/ebitenui/ebitenui/widget"
	"golang.org/x/image/colornames"
)

func SetDefaultTheme() {
	const borderSize = 1

	// COLORS

	// Text
	cText := color.White
	cTextDisabled := color.NRGBA{122, 122, 122, 255}

	// Surfaces
	cSurfacePanelBg := color.RGBA{R: 0x21, G: 0x21, B: 0x21, A: 0xff}
	cSurfacePanelFg := color.RGBA{R: 0x2b, G: 0x2d, B: 0x30, A: 0xff}
	cSurfaceTabBg := color.NRGBA{32, 32, 32, 255}
	cSurfaceListSelected := color.NRGBA{40, 40, 40, 255}
	cSurfaceListSelectedFocused := color.NRGBA{50, 50, 50, 255}
	cSurfaceListFocus := color.NRGBA{99, 99, 99, 255}

	// Buttons / inputs (fills & borders)
	cBtnIdleFill := color.NRGBA{51, 51, 51, 255}
	cBtnIdleBorder := color.NRGBA{81, 81, 81, 255}
	cBtnHoverFill := color.NRGBA{77, 77, 77, 255}
	cBtnPressedFill := color.NRGBA{119, 119, 119, 255}

	cInputIdleFill := color.NRGBA{77, 77, 77, 255}
	cInputDisabledFill := color.NRGBA{47, 47, 47, 255}
	cInputBorder := color.NRGBA{177, 177, 177, 255}

	// Slider handle
	cHandleIdleFill := color.White
	cHandleHoverFill := color.NRGBA{235, 235, 235, 255}
	cHandlePressedFill := color.NRGBA{210, 210, 210, 255}
	cHandleBorder := cInputBorder

	// Scroll container
	cScrollIdle := cInputIdleFill
	cScrollDisabled := cInputDisabledFill

	// List slider handle
	cListHandleIdle := color.NRGBA{110, 110, 110, 255}
	cListHandleHover := color.NRGBA{120, 120, 120, 255}

	// Combo icon tint
	cComboIcon := color.NRGBA{220, 220, 220, 255}

	// Widget background (panel)
	cWidgetPanelBg := colornames.Black

	// Play state colors
	cPlayPlaying := color.RGBA{R: 0x22, G: 0xC5, B: 0x5E, A: 0xFF}
	cPlayStopped := color.RGBA{R: 0xEF, G: 0x44, B: 0x44, A: 0xFF}
	cPlayArmed := color.RGBA{R: 0xF5, G: 0x9E, B: 0x0B, A: 0xFF}
	cPlayNone := colornames.Black
	cPulse := colornames.White

	// ICONS
	iconsBuilder := NewIconsBuilder()

	// FONTS
	face := getFontFace(18)
	menuFace := getFontFace(18)

	theme := &Theme{
		PanelTheme: &PanelTheme{
			ForegroundImage: image.NewBorderedNineSliceColor(cSurfacePanelFg, cSurfacePanelBg, 1),
			BackgroundImage: image.NewNineSliceColor(cSurfacePanelBg),
			Padding:         &widget.Insets{Left: 5, Right: 5, Top: 5, Bottom: 5},
			Spacing:         10,
		},
		Theme: &widget.Theme{
			DefaultFace:      face,
			DefaultTextColor: cText,
			ButtonTheme: &widget.ButtonParams{
				TextColor: &widget.ButtonTextColor{Idle: cText},
				TextFace:  face,
				Image: &widget.ButtonImage{
					Idle:    image.NewBorderedNineSliceColor(cBtnIdleFill, cBtnIdleBorder, borderSize),
					Hover:   image.NewBorderedNineSliceColor(cBtnHoverFill, cBtnIdleFill, borderSize),
					Pressed: image.NewBorderedNineSliceColor(cBtnPressedFill, cBtnHoverFill, borderSize),
				},
				TextPadding: &widget.Insets{Left: 10, Right: 10, Top: 5, Bottom: 5},
				TextPosition: &widget.TextPositioning{
					VTextPosition: widget.TextPositionCenter,
					HTextPosition: widget.TextPositionCenter,
				},
			},
			PanelTheme: &widget.PanelParams{
				BackgroundImage: image.NewNineSliceColor(cWidgetPanelBg),
			},
			LabelTheme: &widget.LabelParams{
				Face:  face,
				Color: &widget.LabelColor{Idle: cText},
			},
			TextTheme: &widget.TextParams{
				Face:  face,
				Color: cText,
				Position: &widget.TextPositioning{
					VTextPosition: widget.TextPositionCenter,
					HTextPosition: widget.TextPositionCenter,
				},
			},
			TabbookTheme: &widget.TabBookParams{
				TabButton: &widget.ButtonParams{
					TextColor: &widget.ButtonTextColor{Idle: cText},
					TextFace:  face,
					Image: &widget.ButtonImage{
						Idle:    image.NewNineSliceColor(cBtnIdleFill),
						Hover:   image.NewNineSliceColor(cBtnHoverFill),
						Pressed: image.NewNineSliceColor(cBtnPressedFill),
					},
					TextPadding: widget.NewInsetsSimple(5),
					MinSize:     &img.Point{98, 40},
				},
				TabSpacing: constantutil.ConstantToPointer(1),
			},
			TabTheme: &widget.TabParams{
				BackgroundImage: image.NewNineSliceColor(cSurfaceTabBg),
			},
			TextInputTheme: &widget.TextInputParams{
				Face: face,
				Image: &widget.TextInputImage{
					Idle:     image.NewBorderedNineSliceColor(cInputIdleFill, cInputBorder, borderSize),
					Disabled: image.NewBorderedNineSliceColor(cInputDisabledFill, cInputBorder, borderSize),
				},
				Color:   &widget.TextInputColor{Idle: cText, Caret: cText},
				Padding: widget.NewInsetsSimple(5),
			},
			SliderTheme: &widget.SliderParams{
				TrackPadding:    widget.NewInsetsSimple(0),
				FixedHandleSize: constantutil.ConstantToPointer(6),
				TrackOffset:     constantutil.ConstantToPointer(0),
				PageSizeFunc:    func() int { return 1 },
				TrackImage: &widget.SliderTrackImage{
					Idle:     image.NewBorderedNineSliceColor(cInputIdleFill, cInputBorder, borderSize),
					Disabled: image.NewBorderedNineSliceColor(cInputDisabledFill, cInputBorder, borderSize),
				},
				HandleImage: &widget.ButtonImage{
					Idle:         image.NewBorderedNineSliceColor(cHandleIdleFill, cHandleBorder, 1),
					Hover:        image.NewBorderedNineSliceColor(cHandleHoverFill, cHandleBorder, borderSize),
					Pressed:      image.NewBorderedNineSliceColor(cHandlePressedFill, cHandleBorder, borderSize),
					PressedHover: image.NewBorderedNineSliceColor(cHandlePressedFill, cHandleBorder, borderSize),
				},
			},
			ListComboButtonTheme: &widget.ListComboButtonParams{
				List: &widget.ListParams{
					EntryFace:                   face,
					EntryTextPadding:            widget.NewInsetsSimple(5),
					EntryTextHorizontalPosition: constantutil.ConstantToPointer(widget.TextPositionStart),
					EntryTextVerticalPosition:   constantutil.ConstantToPointer(widget.TextPositionCenter),
					EntryColor: &widget.ListEntryColor{
						Unselected:                 cText,
						Selected:                   cText,
						SelectedBackground:         cSurfaceListSelected,
						SelectedFocusedBackground:  cSurfaceListSelectedFocused,
						SelectingBackground:        cSurfaceListFocus,
						FocusedBackground:          cSurfaceListFocus,
						SelectingFocusedBackground: cSurfaceListFocus,
					},
					ScrollContainerPadding: widget.NewInsetsSimple(4),
					ScrollContainerImage: &widget.ScrollContainerImage{
						Idle:     image.NewNineSliceColor(cScrollIdle),
						Disabled: image.NewNineSliceColor(cScrollDisabled),
						Mask:     image.NewNineSliceColor(cScrollIdle),
					},
					Slider: &widget.SliderParams{
						TrackImage: &widget.SliderTrackImage{
							Idle:     image.NewNineSliceColor(cScrollIdle),
							Disabled: image.NewNineSliceColor(cScrollDisabled),
						},
						HandleImage: &widget.ButtonImage{
							Idle:    image.NewNineSliceColor(cListHandleIdle),
							Hover:   image.NewNineSliceColor(cListHandleHover),
							Pressed: image.NewBorderedNineSliceColor(cListHandleHover, cListHandleIdle, borderSize),
						},
						TrackPadding: &widget.Insets{Top: 4, Left: 4, Right: 4, Bottom: 4},
					},
				},
				Button: &widget.ButtonParams{
					TextColor: &widget.ButtonTextColor{
						Idle:    cText,
						Hover:   cText,
						Pressed: cText,
					},
					TextFace: face,
					Image: &widget.ButtonImage{
						Idle:    getComboListButtonImage(cBtnIdleFill, cBtnIdleBorder, cComboIcon),
						Hover:   getComboListButtonImage(cBtnHoverFill, cBtnIdleFill, cComboIcon),
						Pressed: getComboListButtonImage(cBtnPressedFill, cBtnHoverFill, cComboIcon),
					},
					TextPadding: &widget.Insets{Left: 15, Right: 25, Top: 5, Bottom: 5},
					TextPosition: &widget.TextPositioning{
						VTextPosition: widget.TextPositionCenter,
						HTextPosition: widget.TextPositionCenter,
					},
				},
				MaxContentHeight: constantutil.ConstantToPointer(200),
			},
			CheckboxTheme: &widget.CheckboxParams{
				Label: &widget.LabelParams{
					Face: face,
					Color: &widget.LabelColor{
						Idle:     cText,
						Disabled: cTextDisabled,
					},
				},
				Image: getCheckboxImage(),
			},
		},
		Icons: iconsBuilder.GetIcons(1, nil),
		IconSizes: IconSizes{
			IconSizeSmall:  24,
			IconSizeMedium: 32,
			IconSizeLarge:  48,
		},
		IconColors: IconColors{
			IconColorLedOn:   colornames.Lime,
			IconColorLedOff:  colornames.Red,
			IconColorDefault: colornames.White,
		},
		MainMenuTheme: &MainMenuTheme{
			ButtonImage: &widget.ButtonImage{
				Idle:  image.NewNineSliceColor(cBtnIdleFill),
				Hover: image.NewNineSliceColor(cBtnHoverFill),
			},
			ButtonPadding: &widget.Insets{Left: 10, Right: 10, Top: 5, Bottom: 5},
			IconSpacing:   5,
			Font:          menuFace,
			TextColor:     cText,
		},
		PlayTheme: &PlayTheme{
			Playing:       NewNineSliceRounded(cPlayPlaying, 10),
			Stopped:       NewNineSliceRounded(cPlayStopped, 10),
			Armed:         NewNineSliceRounded(cPlayArmed, 10),
			None:          NewNineSliceRounded(cPlayNone, 10),
			Pulse:         NewNineSliceRounded(cPulse, 10),
			PulseStrength: 0.45,
			PulseDuration: 150 * time.Millisecond,
		},
	}

	Current = theme
}
