package theme

import (
	"time"

	"github.com/ebitenui/ebitenui/widget"
)

type PanelTheme struct {
	GradientTop             Color
	GradientBottom          Color
	BackgroundImage         Image
	BackgroundColorize      Color
	PulseBackgroundImage    Image
	PulseBackgroundColorize Color
	PulseDuration           time.Duration
	Padding                 *widget.Insets
	Spacing                 int
}

func NewPanelTheme() *PanelTheme {
	return &PanelTheme{
		BackgroundColorize:      ColorNone,
		PulseBackgroundColorize: ColorNone,
	}
}

func (p *PanelTheme) WithBackgroundImage(img Image) *PanelTheme {
	c := p.clone()
	c.BackgroundImage = img
	return c
}

func (p *PanelTheme) WithBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.BackgroundColorize = color
	return c
}

func (p *PanelTheme) WithPulseBackgroundImage(img Image) *PanelTheme {
	c := p.clone()
	c.PulseBackgroundImage = img
	return c
}

func (p *PanelTheme) WithPulseBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.PulseBackgroundColorize = color
	return c
}

func (p *PanelTheme) WithPulseDuration(duration time.Duration) *PanelTheme {
	c := p.clone()
	c.PulseDuration = duration
	return c
}

func (p *PanelTheme) WithPadding(padding *widget.Insets) *PanelTheme {
	c := p.clone()
	c.Padding = padding
	return c
}

func (p *PanelTheme) WithSimplePadding(padding int) *PanelTheme {
	return p.WithPadding(widget.NewInsetsSimple(padding))
}

func (p *PanelTheme) WithSpacing(spacing int) *PanelTheme {
	c := p.clone()
	c.Spacing = spacing
	return c
}

func (p *PanelTheme) WithGradient(top, bottom Color) *PanelTheme {
	c := p.clone()
	c.GradientTop = top
	c.GradientBottom = bottom
	return c
}

func (p *PanelTheme) clone() *PanelTheme {
	return &PanelTheme{
		GradientTop:             p.GradientTop,
		GradientBottom:          p.GradientBottom,
		BackgroundImage:         p.BackgroundImage,
		BackgroundColorize:      p.BackgroundColorize,
		PulseBackgroundImage:    p.PulseBackgroundImage,
		PulseBackgroundColorize: p.PulseBackgroundColorize,
		PulseDuration:           p.PulseDuration,
		Padding:                 p.Padding,
		Spacing:                 p.Spacing,
	}
}
