package theme

import (
	"bytes"
	_ "embed"
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// SRC: https://feathericons.com/
//
//go:embed icons.png
var sheet []byte

type Tile int

const (
	TileBarDelete Tile = iota
	TileDelete
	TileGear
	TileMinus
	TilePause
	TilePlay
	TileBarAdd
	TilePlus
	TileRecord
	TileStop
	TileButtonHover
	TileButton
	TileButtonPressed
	TileLed
)

type TileMap map[Tile]struct{ x, y float64 }

var tileMap = TileMap{
	TileBarDelete:     {0, 0},
	TileDelete:        {1, 0},
	TileGear:          {2, 0},
	TileMinus:         {3, 0},
	TilePause:         {4, 0},
	TilePlay:          {5, 0},
	TileBarAdd:        {6, 0},
	TilePlus:          {7, 0},
	TileRecord:        {0, 1},
	TileStop:          {1, 1},
	TileButtonHover:   {2, 1},
	TileButton:        {3, 1},
	TileButtonPressed: {4, 1},
	TileLed:           {5, 1},
}

type Sheet struct {
	img *ebiten.Image
	cs  float64 // rectangular cells
	tm  TileMap
	cr  image.Rectangle
}

func NewInternalSheet() *Sheet {
	img, _, err := ebitenutil.NewImageFromReader(bytes.NewReader(sheet))
	if err != nil {
		panic(err)
	}
	return NewSheet(img, 64, tileMap)
}

func NewSheet(img *ebiten.Image, cs float64, tm TileMap) *Sheet {
	return &Sheet{
		img: img,
		cs:  cs,
		tm:  tm,
	}
}

func (s *Sheet) GetTile(t Tile) *ebiten.Image {
	tile, ok := s.tm[t]
	if !ok {
		panic("tile not found")
	}

	x := math.Round(tile.x * s.cs)
	y := math.Round(tile.y * s.cs)

	rect := image.Rect(int(x), int(y), int(x+s.cs), int(y+s.cs))
	if !s.cr.Empty() {
		rect.Min.X = rect.Min.X + s.cr.Min.X
		rect.Min.Y = rect.Min.Y + s.cr.Min.Y
		rect.Max.X = rect.Min.X + s.cr.Dx()
		rect.Max.Y = rect.Min.Y + s.cr.Dy()
	}

	return s.img.SubImage(rect).(*ebiten.Image)
}

func (s *Sheet) Crop(rect image.Rectangle) *Sheet {
	ns := NewSheet(s.img, s.cs, s.tm)
	ns.cr = rect
	return ns
}
