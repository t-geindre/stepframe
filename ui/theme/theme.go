package theme

import (
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var Current *Theme

type Theme struct {
	*widget.Theme
	MainPanelTheme         *PanelTheme
	PanelTheme             *PanelTheme
	AddTrackPanelTheme     *PanelTheme
	BarContainerPanelTheme *PanelTheme
	BarOnPanelTheme        *PanelTheme
	BarOffPanelTheme       *PanelTheme
	BarActivePanelTheme    *PanelTheme
	LedPanelTheme          *PanelTheme
	IconsTheme             *IconsTheme
	Images                 map[Image]*OffsetImage
	TrackTheme             *TrackTheme
	Colors                 ColorTheme
}

type TrackTheme struct {
	Id *TrackIdTheme
}

type TrackIdTheme struct {
	Font *text.Face
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
	ColorNone Color = iota
	ColorLedOn
	ColorLedIdle
	ColorLedArmed
	ColorLedPulse
	ColorIconDefault
	ColorBarOn
	ColorBarOff
	ColorBarActive
	ColorBarPulse
	ColorMainGradientTop
	ColorMainGradientBottom
)

type Image int

const (
	ImageNone Image = iota
	ImagePanel
	ImagePanelHover
	ImagePanelPressed
	ImageAddTrackPanel
	ImageBarContainerPanel
	ImageBarOffPanel
	ImageBarOnPanel
	ImageBarActivePanel
	ImageBarActivePanelPulse
	ImageLedPanel
	ImageLedPanelPulse
)
