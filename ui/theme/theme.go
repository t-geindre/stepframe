package theme

import (
	"image/color"
	"time"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var Current *Theme

type Theme struct {
	*widget.Theme
	PanelTheme      *PanelTheme
	Icons           Icons
	IconSizes       IconSizes
	IconColors      IconColors
	MainMenuTheme   *MainMenuTheme
	PlayTheme       *PlayTheme
	LedTheme        *LedTheme
	BackgroundTheme *BackgroundTheme
}

type BackgroundTheme struct {
	GradientTop    color.Color
	GradientBottom color.Color
	Padding        *widget.Insets
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

type LedTheme struct {
	OnImage          *ebiten.Image
	OffImage         *ebiten.Image
	OffsetX, OffsetY int
	Width, Height    int
	PulseDuration    time.Duration
}

type Icons map[Icon]*ebiten.Image
type Icon int

const (
	IconBarDelete Icon = iota
	IconDelete
	IconGear
	IconMinus
	IconPause
	IconPlay
	IconBarAdd
	IconPlus
	IconRecord
	IconStop
	IconButton
	IconLed
	IconNone
)

type IconSize int
type IconSizes map[IconSize]int

const (
	IconSizeSmall IconSize = iota
	IconSizeMedium
	IconSizeLarge
)

type IconColor int
type IconColors map[IconColor]color.Color

const (
	IconColorOn IconColor = iota
	IconColorOff
	IconColorIdle
	IconColorArmed
	IconColorDefault
	IconColorNone
)
