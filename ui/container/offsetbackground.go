package container

import (
	"image/color"
	"stepframe/ui/theme"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
)

type OffsetBackground[T widget.Containerer] struct {
	*widget.Container
	outer      T
	background *theme.OffsetImage
	colorize   color.Color

	pulseBackground *theme.OffsetImage
	pulseColorize   color.Color
	pulseDuration   time.Duration
	lastPulse       time.Time
}

func NewOffsetBackground[T widget.Containerer](container *widget.Container, outer T) *OffsetBackground[T] {
	return &OffsetBackground[T]{
		Container: container,
		outer:     outer,
	}
}

func (b *OffsetBackground[T]) SetBackgroundOffsetImage(i theme.Image) T {
	img := theme.Current.Images[i]
	b.background = img

	if img == nil || img.Offset == nil {
		return b.outer
	}

	imw, imh := img.Image.MinSize()
	wmw, wmh := b.GetWidget().MinHeight, b.GetWidget().MinHeight

	if wmw == 0 {
		b.GetWidget().MinWidth = imw - b.background.Offset.Dx()
	}

	if wmh == 0 {
		b.GetWidget().MinHeight = imh - b.background.Offset.Dy()
	}

	return b.outer
}

func (b *OffsetBackground[T]) SetPulseBackgroundOffsetImage(i theme.Image) T {
	b.pulseBackground = theme.Current.Images[i]
	return b.outer
}

func (b *OffsetBackground[T]) SetPulseDuration(d time.Duration) T {
	b.pulseDuration = d
	return b.outer
}

func (b *OffsetBackground[T]) SetPulseColorizedBackgroundOffsetImage(col theme.Color) T {
	b.pulseColorize = theme.Current.Colors[col]
	return b.outer
}

func (b *OffsetBackground[T]) SetColorizedBackgroundOffsetImage(col theme.Color) T {
	b.colorize = theme.Current.Colors[col]
	return b.outer
}

func (b *OffsetBackground[T]) Render(screen *ebiten.Image) {
	b.renderBackground(screen)
	b.renderPulse(screen)
}

func (b *OffsetBackground[T]) Pulse() {
	b.lastPulse = time.Now()
}

func (b *OffsetBackground[T]) renderBackground(screen *ebiten.Image) {
	if b.background == nil {
		return
	}

	rect := b.GetWidget().Rect
	if b.background.Offset != nil {
		rect = b.background.Offset.Apply(rect)
	}

	b.background.Image.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		if b.colorize != nil {
			opts.ColorScale.ScaleWithColor(b.colorize)
		}
	})
}

func (b *OffsetBackground[T]) renderPulse(screen *ebiten.Image) {
	if b.pulseBackground == nil || b.lastPulse.IsZero() {
		return
	}

	elapsed := time.Since(b.lastPulse)
	if elapsed > b.pulseDuration {
		return
	}

	t := 1.0 - float64(elapsed)/float64(b.pulseDuration)
	t = smootherstep(t)
	t = t * t

	rect := b.GetWidget().Rect
	if b.pulseBackground.Offset != nil {
		rect = b.pulseBackground.Offset.Apply(rect)
	}

	b.pulseBackground.Image.Draw(screen, rect.Dx(), rect.Dy(), func(opts *ebiten.DrawImageOptions) {
		opts.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
		if b.pulseColorize != nil {
			opts.ColorScale.ScaleWithColor(b.pulseColorize)
		}
		opts.ColorScale.ScaleAlpha(float32(t))
	})
}

// UTILS

func smootherstep(t float64) float64 {
	if t <= 0 {
		return 0
	}
	if t >= 1 {
		return 1
	}
	// 6t^5 - 15t^4 + 10t^3
	return t * t * t * (t*(t*6-15) + 10)
}
