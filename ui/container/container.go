package container

import (
	_ "embed"
	"image/color"
	"stepframe/ui/theme"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed gradient.kage
var gradientShaderRaw []byte
var gradientShader *ebiten.Shader

type Container[T widget.Containerer] struct {
	*widget.Container
	outer T

	gradientOn               bool
	gradientFrom, gradientTo [4]float32
	gradientCache            *ebiten.Image

	offsetBackground *theme.OffsetImage
}

func NewContainer[T widget.Containerer](layout widget.Layouter, outer T) *Container[T] {
	return &Container[T]{
		Container: widget.NewContainer(
			widget.ContainerOpts.Layout(layout),
		),
		outer: outer,
	}
}

func (c *Container[T]) SetBackgroundOffsetImage(img *theme.OffsetImage) T {
	if img.Offset.X != 0 || img.Offset.Y != 0 {
		c.offsetBackground = img
		return c.outer
	}

	c.SetBackgroundImage(img.Image)

	return c.outer
}

func (c *Container[T]) SetBackgroundImage(img *image.NineSlice) T {
	if c.Container.IsValidated() {
		c.Container.SetBackgroundImage(img)
	} else {
		widget.ContainerOpts.BackgroundImage(img)(c.Container)
	}
	return c.outer
}

func (c *Container[T]) SetMinSize(width, height int) T {
	widget.WidgetOpts.MinSize(width, height)(c.Container.GetWidget())
	return c.outer
}

func (c *Container[T]) SetVerticalGradientBackground(top, bottom color.Color) T {
	c.gradientOn = true

	r1, g1, b1, a1 := top.RGBA()
	r2, g2, b2, a2 := bottom.RGBA()

	c.gradientFrom = [4]float32{float32(r1) / 0xffff, float32(g1) / 0xffff, float32(b1) / 0xffff, float32(a1) / 0xffff}
	c.gradientTo = [4]float32{float32(r2) / 0xffff, float32(g2) / 0xffff, float32(b2) / 0xffff, float32(a2) / 0xffff}

	return c.outer
}

func (c *Container[T]) Render(screen *ebiten.Image) {
	c.computeGradient()
	c.Container.Render(screen)
	c.renderOffsetBackground(screen)
}

func (c *Container[T]) renderOffsetBackground(screen *ebiten.Image) {
	if c.offsetBackground == nil {
		return
	}

	rect := c.GetWidget().Rect
	rect.Min = rect.Min.Sub(c.offsetBackground.Offset)
	rect.Max = rect.Max.Add(c.offsetBackground.Offset)

	c.offsetBackground.Image.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	})
}

func (c *Container[T]) computeGradient() {
	if !c.gradientOn {
		return
	}

	if gradientShader == nil {
		var err error
		gradientShader, err = ebiten.NewShader(gradientShaderRaw)
		if err != nil {
			panic("failed to create gradient shader: " + err.Error())
		}
	}

	rect := c.GetWidget().Rect
	if c.gradientCache != nil && c.gradientCache.Bounds().Dy() == rect.Dy() && c.gradientCache.Bounds().Dx() == rect.Dx() {
		return
	}

	c.gradientCache = ebiten.NewImage(rect.Dx(), rect.Dy())

	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"TopColor":    c.gradientFrom,
		"BottomColor": c.gradientTo,
	}

	c.gradientCache.DrawRectShader(rect.Dx(), rect.Dy(), gradientShader, opts)
	c.SetBackgroundImage(image.NewNineSlice(c.gradientCache, [3]int{0, rect.Dx(), 0}, [3]int{0, rect.Dy(), 0}))
}
