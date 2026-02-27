package theme

import (
	img "image"

	"github.com/ebitenui/ebitenui/image"
)

type OffsetImage struct {
	Image  *image.NineSlice
	Offset *Outsets
}

func NewSimpleOffsetImage(img *image.NineSlice) *OffsetImage {
	return &OffsetImage{
		Image: img,
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
