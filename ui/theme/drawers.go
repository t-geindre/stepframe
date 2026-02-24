package theme

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func getCheckboxImage() *widget.CheckboxImage {
	const size = 16
	const borderSize = 1
	const crossPadding = 4
	const crossWidth = 2

	border := image.NewBorderedNineSliceColor(color.NRGBA{119, 119, 119, 255}, color.White, borderSize)
	idle := ebiten.NewImage(size, size)
	border.Draw(idle, size, size, nil)
	idle9s := image.NewFixedNineSlice(idle)

	// Create checked image
	checked := ebiten.NewImage(size, size)
	border.Draw(checked, size, size, nil)
	vector.StrokeLine(checked, crossPadding, crossPadding, size-crossPadding, size-crossPadding, crossWidth, color.White, true)
	vector.StrokeLine(checked, size-crossPadding, crossPadding, crossPadding, size-crossPadding, crossWidth, color.White, true)

	checked9s := image.NewFixedNineSlice(checked)

	// Create greyed image
	greyed := ebiten.NewImage(size, size)
	border.Draw(greyed, size, size, nil)
	vector.StrokeLine(greyed, 5, 16, 27, 16, 3, color.White, true)
	greyed9s := image.NewFixedNineSlice(greyed)

	return &widget.CheckboxImage{
		Unchecked:         idle9s,
		Checked:           checked9s,
		Greyed:            greyed9s,
		UncheckedHovered:  idle9s,
		CheckedHovered:    checked9s,
		GreyedHovered:     greyed9s,
		UncheckedDisabled: idle9s,
		CheckedDisabled:   checked9s,
		GreyedDisabled:    greyed9s,
	}
}

func getComboListButtonImage(bgCol, borderCol, arrowCol color.Color) *image.NineSlice {
	const (
		border = 1
		w, h   = 44, 24

		leftCap   = 4
		rightCap  = 30
		topCap    = 4
		bottomCap = 4
	)

	i := ebiten.NewImage(w, h)
	i.Fill(bgCol)

	vector.StrokeRect(i, 0, 0, float32(w-border), float32(h-border), float32(border), borderCol, false)

	cx := float32(w - rightCap/2)
	cy := float32(h / 2)

	size := float32(4)

	var p vector.Path
	p.MoveTo(cx-size, cy-size)
	p.LineTo(cx, cy+size)
	p.LineTo(cx+size, cy-size)
	p.Close()

	dpOpt := vector.DrawPathOptions{AntiAlias: true}
	dpOpt.ColorScale.ScaleWithColor(arrowCol)

	flOpt := vector.FillOptions{
		FillRule: vector.FillRuleEvenOdd,
	}

	vector.FillPath(i, &p, &flOpt, &dpOpt)

	ws := [3]int{
		leftCap,
		w - leftCap - rightCap,
		rightCap,
	}
	hs := [3]int{
		topCap,
		h - topCap - bottomCap,
		bottomCap,
	}

	return image.NewNineSlice(i, ws, hs)
}
func NewNineSliceRounded(col color.Color, radius int) *image.NineSlice {
	return NewNineSliceRoundedBordered(col, col, radius, 0)
}

func NewNineSliceRoundedBordered(fillCol, borderCol color.Color, radius, borderWidth int) *image.NineSlice {
	size := 2*radius + 1
	img := ebiten.NewImage(size, size)

	if borderWidth <= 0 {
		drawRoundedRect(img, 0, 0, float32(size), float32(size), float32(radius), fillCol)
		return image.NewNineSlice(img, [3]int{radius, 1, radius}, [3]int{radius, 1, radius})
	}

	// Outer rounded rect = border
	drawRoundedRect(img, 0, 0, float32(size), float32(size), float32(radius), borderCol)

	// Inner rounded rect = fill, inset by borderWidth
	in := float32(borderWidth)
	innerW := float32(size) - 2*in
	innerH := float32(size) - 2*in
	if innerW > 0 && innerH > 0 {
		innerR := float32(radius - borderWidth)
		if innerR < 0 {
			innerR = 0
		}
		drawRoundedRect(img, in, in, innerW, innerH, innerR, fillCol)
	}

	return image.NewNineSlice(img, [3]int{radius, 1, radius}, [3]int{radius, 1, radius})
}

func drawRoundedRect(img *ebiten.Image, x, y, w, h float32, r float32, col color.Color) {
	if w <= 0 || h <= 0 {
		return
	}
	if r <= 0 {
		vector.FillRect(img, x, y, w, h, col, true)
		return
	}

	// Clamp radius so it never exceeds half the size.
	maxR := w / 2
	if h/2 < maxR {
		maxR = h / 2
	}
	if r > maxR {
		r = maxR
	}

	// Cross (two rectangles)
	vector.FillRect(img, x+r, y, w-2*r, h, col, true)
	vector.FillRect(img, x, y+r, w, h-2*r, col, true)

	// Four corners (circles)
	vector.FillCircle(img, x+r, y+r, r, col, true)     // top-left
	vector.FillCircle(img, x+w-r, y+r, r, col, true)   // top-right
	vector.FillCircle(img, x+r, y+h-r, r, col, true)   // bottom-left
	vector.FillCircle(img, x+w-r, y+h-r, r, col, true) // bottom-right
}
