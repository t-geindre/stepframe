package theme

import (
	"bytes"
	_ "embed"
	"image/color"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// SRC: https://feathericons.com/
//
//go:embed icons.png
var iconsImage []byte

type Icons map[Icon]*widget.GraphicImage
type Icon int
type IconAccent int

const (
	IconAdd Icon = iota
	IconCamera
	IconCenter
	IconChunk
	IconDelete
	IconFile
	IconFullscreen
	IconNoise
	IconOpen
	IconSave
	IconStats
	IconZoom

	IconCellSize float64 = 24

	IconAccentNone IconAccent = iota
	IconAccentCamera
	IconAccentFile
	IconAccentNoise
)

var iconsMap = map[Icon]struct {
	x, y float64
	acc  IconAccent
}{
	IconAdd:        {0, 0, IconAccentNoise},
	IconCamera:     {1, 0, IconAccentCamera},
	IconCenter:     {2, 0, IconAccentCamera},
	IconChunk:      {3, 0, IconAccentNoise},
	IconDelete:     {4, 0, IconAccentNoise},
	IconFile:       {5, 0, IconAccentFile},
	IconFullscreen: {6, 0, IconAccentCamera},
	IconNoise:      {0, 1, IconAccentNoise},
	IconOpen:       {1, 1, IconAccentFile},
	IconSave:       {2, 1, IconAccentFile},
	IconStats:      {3, 1, IconAccentFile},
	IconZoom:       {4, 1, IconAccentCamera},
}

type IconsBuilder struct {
	// accents[accents][scale] = sheet
	accents map[IconAccent]map[float64]*Sheet
}

func NewIconsBuilder(accNone, accCam, accFile, accNoise color.Color) *IconsBuilder {
	ib := &IconsBuilder{
		accents: make(map[IconAccent]map[float64]*Sheet),
	}

	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(iconsImage))
	if err != nil {
		panic(err)
	}

	sheet := NewSheet(img, IconCellSize)

	for _, acc := range []struct {
		acc IconAccent
		col color.Color
	}{
		{IconAccentNone, accNone},
		{IconAccentCamera, accCam},
		{IconAccentFile, accFile},
		{IconAccentNoise, accNoise},
	} {
		ib.accents[acc.acc] = make(map[float64]*Sheet)
		ib.accents[acc.acc][1] = sheet.Colorize(acc.col)
	}

	return ib
}

func (ib *IconsBuilder) GetIcons(withIdleAcc, withHoverAcc bool, scale float64) Icons {
	icons := make(Icons)
	for icon, info := range iconsMap {
		idleAcc := IconAccentNone
		if withIdleAcc {
			idleAcc = info.acc
		}

		hoverAcc := IconAccentNone
		if withHoverAcc {
			hoverAcc = info.acc
		}

		icons[icon] = &widget.GraphicImage{
			Idle:  ib.getSheet(idleAcc, scale).Get(info.x, info.y),
			Hover: ib.getSheet(hoverAcc, scale).Get(info.x, info.y),
		}
	}

	return icons
}

func (ib *IconsBuilder) getSheet(acc IconAccent, scale float64) *Sheet {
	if _, ok := ib.accents[acc]; !ok {
		panic("invalid accent")
	}

	if _, ok := ib.accents[acc][scale]; !ok {
		ib.accents[acc][scale] = ib.accents[acc][1].Scale(scale)
	}

	return ib.accents[acc][scale]
}
