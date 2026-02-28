package theme

import (
	img "image"

	"github.com/ebitenui/ebitenui/image"
	"github.com/hajimehoshi/ebiten/v2"
)

type OffsetImage struct {
	Image  *image.NineSlice
	Offset *Outsets
}

func NewSimpleOffsetImage(img *ebiten.Image, v1, v2, v3, h1, h2, h3 int) *OffsetImage {
	return &OffsetImage{
		Image: image.NewNineSlice(img, [3]int{v1, v2, v3}, [3]int{h1, h2, h3}),
	}
}

func NewOffsetImage(img *ebiten.Image, v1, v2, v3, h1, h2, h3 int, offset int) *OffsetImage {
	return &OffsetImage{
		Image:  image.NewNineSlice(img, [3]int{v1, v2, v3}, [3]int{h1, h2, h3}),
		Offset: NewOutsetsSimple(offset),
	}
}

type Outsets struct {
	Top    int
	Left   int
	Right  int
	Bottom int
}

func NewOutsetsSimple(widthHeight int) *Outsets {
	return &Outsets{
		Top:    widthHeight,
		Left:   widthHeight,
		Right:  widthHeight,
		Bottom: widthHeight,
	}
}

func (i Outsets) Apply(rect img.Rectangle) img.Rectangle {
	rect.Min = rect.Min.Sub(img.Point{X: i.Left, Y: i.Top})
	rect.Max = rect.Max.Add(img.Point{X: i.Right, Y: i.Bottom})
	return rect
}

func (i Outsets) Dx() int {
	return i.Left + i.Right
}

func (i Outsets) Dy() int {
	return i.Top + i.Bottom
}
