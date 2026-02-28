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
	PanelTheme             *PanelTheme
	AddTrackPanelTheme     *PanelTheme
	BarContainerPanelTheme *PanelTheme
	BarPanelTheme          *PanelTheme
	BarPanelOffTheme       *PanelTheme
	LedPanelTheme          *PanelTheme
	IconsTheme             *IconsTheme
	MainContainerTheme     *MainContainerTheme
	TrackTheme             *TrackTheme
	Colors                 ColorTheme
}

type MainContainerTheme struct {
	GradientTop    color.Color
	GradientBottom color.Color
	Padding        *widget.Insets
	Spacing        int
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
	ColorOn Color = iota
	ColorIdle
	ColorArmed
	IconColorDefault
	ColorNone
)
