package theme

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// SRC: https://feathericons.com/
//
//go:embed icons.png
var iconsImage []byte

type Icons map[Icon]*ebiten.Image
type Icon int

const (
	IconClear Icon = iota
	IconDelete
	IconGear
	IconPause
	IconPlay
	IconPlus
	IconRecord
	IconStop
	IconLed

	IconCellSize float64 = 64
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

var iconsMap = map[Icon]struct{ x, y float64 }{
	IconClear:  {0, 0},
	IconDelete: {1, 0},
	IconGear:   {2, 0},
	IconPause:  {3, 0},
	IconPlay:   {4, 0},
	IconPlus:   {5, 0},
	IconRecord: {6, 0},
	IconStop:   {7, 0},
	IconLed:    {0, 1},
}

type IconsBuilder struct {
	sheet *Sheet
}

func NewIconsBuilder() *IconsBuilder {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(iconsImage))
	if err != nil {
		panic(err)
	}

	return &IconsBuilder{
		sheet: NewSheet(img, IconCellSize),
	}
}

func (ib *IconsBuilder) GetIcons(scale float64, col *color.Color) Icons {
	s := ib.sheet
	if scale != 1 {
		s = s.Scale(scale)
	}
	if col != nil {
		s = s.Colorize(*col)
	}

	icons := make(Icons)

	for i, c := range iconsMap {
		icons[i] = ib.sheet.Get(c.x, c.y)
	}

	return icons
}
