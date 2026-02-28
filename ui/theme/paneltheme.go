package theme

import (
	"time"

	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

type PanelTheme struct {
	GradientTop               Color
	GradientBottom            Color
	IdleBackgroundImage       Image
	IdleBackgroundColorize    Color
	HoverBackgroundImage      Image
	HoverBackgroundColorize   Color
	PressedBackgroundImage    Image
	PressedBackgroundColorize Color
	PulseBackgroundImage      Image
	PulseBackgroundColorize   Color
	PulseDuration             time.Duration
	Padding                   *widget.Insets
	Spacing                   int
	Cursor                    string
}

func NewPanelTheme() *PanelTheme {
	return &PanelTheme{
		IdleBackgroundColorize:    ColorNone,
		HoverBackgroundColorize:   ColorNone,
		PressedBackgroundColorize: ColorNone,
		IdleBackgroundImage:       ImageNone,
		HoverBackgroundImage:      ImageNone,
		PressedBackgroundImage:    ImageNone,
		PulseBackgroundColorize:   ColorNone,
		Cursor:                    input.CURSOR_NONE,
	}
}

func (p *PanelTheme) WithIdleBackgroundImage(img Image) *PanelTheme {
	c := p.clone()
	c.IdleBackgroundImage = img
	return c
}

func (p *PanelTheme) WithIdleBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.IdleBackgroundColorize = color
	return c
}

func (p *PanelTheme) WithHoverBackgroundImage(img Image) *PanelTheme {
	c := p.clone()
	c.HoverBackgroundImage = img
	return c
}

func (p *PanelTheme) WithHoverBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.HoverBackgroundColorize = color
	return c
}

func (p *PanelTheme) WithPressedBackgroundImage(img Image) *PanelTheme {
	c := p.clone()
	c.PressedBackgroundImage = img
	return c
}

func (p *PanelTheme) WithPressedBackgroundColorize(color Color) *PanelTheme {
	c := p.clone()
	c.PressedBackgroundColorize = color
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

func (p *PanelTheme) WithCursor(cursor string) *PanelTheme {
	c := p.clone()
	c.Cursor = cursor
	return c
}

func (p *PanelTheme) clone() *PanelTheme {
	return &PanelTheme{
		GradientTop:               p.GradientTop,
		GradientBottom:            p.GradientBottom,
		IdleBackgroundImage:       p.IdleBackgroundImage,
		IdleBackgroundColorize:    p.IdleBackgroundColorize,
		HoverBackgroundImage:      p.HoverBackgroundImage,
		HoverBackgroundColorize:   p.HoverBackgroundColorize,
		PressedBackgroundImage:    p.PressedBackgroundImage,
		PressedBackgroundColorize: p.PressedBackgroundColorize,
		PulseBackgroundImage:      p.PulseBackgroundImage,
		PulseBackgroundColorize:   p.PulseBackgroundColorize,
		PulseDuration:             p.PulseDuration,
		Padding:                   p.Padding,
		Spacing:                   p.Spacing,
		Cursor:                    p.Cursor,
	}
}
