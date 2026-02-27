package theme

import (
	"time"

	"github.com/ebitenui/ebitenui/widget"
)

type PanelTheme struct {
	BackgroundImage         *OffsetImage
	BackgroundColorize      Color
	PulseBackgroundImage    *OffsetImage
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

func (p *PanelTheme) WithBackgroundImage(img *OffsetImage) *PanelTheme {
	c := p.clone()
	c.BackgroundImage = img
	return c
}

func (p *PanelTheme) WithBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.BackgroundColorize = color
	return c
}

func (p *PanelTheme) WithPulseBackgroundImage(img *OffsetImage) *PanelTheme {
	c := p.clone()
	c.PulseBackgroundImage = img
	return c
}

func (p *PanelTheme) WithPulseBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.PulseBackgroundColorize = color
	return c
}

func (p *PanelTheme) WithPadding(padding *widget.Insets) *PanelTheme {
	c := p.clone()
	c.Padding = padding
	return c
}

func (p *PanelTheme) WithSpacing(spacing int) *PanelTheme {
	c := p.clone()
	c.Spacing = spacing
	return c
}

func (p *PanelTheme) clone() *PanelTheme {
	return &PanelTheme{
		BackgroundImage:         p.BackgroundImage,
		BackgroundColorize:      p.BackgroundColorize,
		PulseBackgroundImage:    p.PulseBackgroundImage,
		PulseBackgroundColorize: p.PulseBackgroundColorize,
		PulseDuration:           p.PulseDuration,
		Padding:                 p.Padding,
		Spacing:                 p.Spacing,
	}
}
