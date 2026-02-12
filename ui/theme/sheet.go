package theme

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sheet struct {
	cs  float64 // rectangular cells
	img *ebiten.Image
}

func NewSheet(img *ebiten.Image, cs float64) *Sheet {
	return &Sheet{
		img: img,
		cs:  cs,
	}
}

func (s *Sheet) Get(x, y float64) *ebiten.Image {
	x0, yo, x1, y1 := math.Round(x*s.cs), math.Round(y*s.cs), math.Round((x+1)*s.cs), math.Round((y+1)*s.cs)
	rect := image.Rect(int(x0), int(yo), int(x1), int(y1))
	return s.img.SubImage(rect).(*ebiten.Image)
}

func (s *Sheet) Colorize(col color.Color) *Sheet {
	img := ebiten.NewImage(s.img.Bounds().Dx(), s.img.Bounds().Dy())
	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleWithColor(col)
	img.DrawImage(s.img, opts)

	return &Sheet{
		img: img,
		cs:  s.cs,
	}
}

func (s *Sheet) Scale(factor float64) *Sheet {
	w, h := s.img.Bounds().Dx(), s.img.Bounds().Dy()
	nw, nh := int(float64(w)*factor), int(float64(h)*factor)

	img := ebiten.NewImage(nw, nh)
	opts := &ebiten.DrawImageOptions{
		Filter: ebiten.FilterLinear,
	}
	opts.GeoM.Scale(factor, factor)
	img.DrawImage(s.img, opts)

	return &Sheet{
		img: img,
		cs:  s.cs * factor,
	}
}
