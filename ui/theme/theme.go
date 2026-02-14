package theme

import (
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Theme struct {
	*widget.Theme
	PanelTheme    *PanelTheme
	Icons         Icons
	MainMenuTheme *MainMenuTheme
	PlayTheme     *PlayTheme
}

type MainMenuTheme struct {
	ButtonImage   *widget.ButtonImage
	ButtonPadding *widget.Insets
	IconSpacing   int
	Font          *text.Face
	TextColor     color.Color
}

type PanelTheme struct {
	BackgroundImage *image.NineSlice
	ForegroundImage *image.NineSlice
	Padding         *widget.Insets
	Spacing         int
}

type PlayTheme struct {
	Playing       *image.NineSlice
	Stopped       *image.NineSlice
	Armed         *image.NineSlice
	None          *image.NineSlice
	Pulse         *image.NineSlice
	PulseStrength float64
	PulseDuration time.Duration
}
