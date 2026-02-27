package container

import (
	_ "embed"
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed gradient.kage
var gradientShaderRaw []byte
var gradientShader *ebiten.Shader

type Gradient[T widget.Containerer] struct {
	*widget.Container
	outer T
	opts  *ebiten.DrawRectShaderOptions
	cache *ebiten.Image
}

func NewGradient[T widget.Containerer](container *widget.Container, outer T) *Gradient[T] {
	return &Gradient[T]{
		Container: container,
		outer:     outer,
	}
}

func (g *Gradient[T]) SetVerticalGradientBackground(top, bottom color.Color) T {
	r1, g1, b1, a1 := top.RGBA()
	r2, g2, b2, a2 := bottom.RGBA()

	g.opts = &ebiten.DrawRectShaderOptions{}
	g.opts.Uniforms = map[string]any{
		"TopColor": [4]float32{
			float32(r1) / 0xffff, float32(g1) / 0xffff, float32(b1) / 0xffff, float32(a1) / 0xffff,
		},
		"BottomColor": [4]float32{
			float32(r2) / 0xffff, float32(g2) / 0xffff, float32(b2) / 0xffff, float32(a2) / 0xffff,
		},
	}

	return g.outer
}

func (g *Gradient[T]) Render(*ebiten.Image) {
	if g.opts == nil {
		return
	}

	if gradientShader == nil {
		var err error
		gradientShader, err = ebiten.NewShader(gradientShaderRaw)
		if err != nil {
			panic("failed to create gradient shader: " + err.Error())
		}
	}

	rect := g.GetWidget().Rect
	if g.cache != nil && g.cache.Bounds().Dy() == rect.Dy() && g.cache.Bounds().Dx() == rect.Dx() {
		return
	}

	g.cache = ebiten.NewImage(rect.Dx(), rect.Dy())
	g.cache.DrawRectShader(rect.Dx(), rect.Dy(), gradientShader, g.opts)
	g.SetBackgroundImage(image.NewNineSlice(g.cache, [3]int{0, rect.Dx(), 0}, [3]int{0, rect.Dy(), 0}))
}
