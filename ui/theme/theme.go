package theme

import (
	img "image"
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
	PanelTheme         *PanelTheme
	VirtualPanelTheme  *PanelTheme
	DropPanelTheme     *PanelTheme
	ShinyPanelTheme    *PanelTheme
	IconsTheme         *IconsTheme
	LedTheme           *LedTheme
	MainContainerTheme *MainContainerTheme
	TrackTheme         *TrackTheme
	Colors             ColorTheme
}

type MainContainerTheme struct {
	GradientTop    color.Color
	GradientBottom color.Color
	Padding        *widget.Insets
	Spacing        int
}

type PanelTheme struct {
	BackgroundImage *OffsetImage
	Padding         *widget.Insets
	Spacing         int
}

type LedTheme struct {
	OnImage       *ebiten.Image
	OffImage      *ebiten.Image
	Width, Height int
	PulseDuration time.Duration
}

type TrackTheme struct {
	Id *TrackIdTheme
}

type TrackIdTheme struct {
	Font *text.Face
}

type OffsetImage struct {
	Image  *image.NineSlice
	Offset img.Point
}

// ICONS

type Icon int
type IconSize int
type IconsTheme struct {
	Icons map[Icon]*ebiten.Image
	Sizes map[IconSize]int
}

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
	IconNone
)

const (
	IconSizeSmall IconSize = iota
	IconSizeMedium
	IconSizeLarge
)

// COLORS SET

type Color int
type ColorTheme map[Color]color.Color

const (
	ColorOn Color = iota
	ColorOff
	ColorIdle
	ColorArmed
	IconColorDefault
	ColorNone
)
