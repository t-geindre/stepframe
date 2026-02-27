package theme

import "github.com/ebitenui/ebitenui/widget"

type PanelTheme struct {
	BackgroundImage *OffsetImage
	Padding         *widget.Insets
	Spacing         int
}

func NewPanelTheme() *PanelTheme {
	return &PanelTheme{}
}

func (p *PanelTheme) WithBackgroundImage(img *OffsetImage) *PanelTheme {
	return &PanelTheme{
		BackgroundImage: img,
		Padding:         p.Padding,
		Spacing:         p.Spacing,
	}
}

func (p *PanelTheme) WithPadding(padding *widget.Insets) *PanelTheme {
	return &PanelTheme{
		BackgroundImage: p.BackgroundImage,
		Padding:         padding,
		Spacing:         p.Spacing,
	}
}

func (p *PanelTheme) WithSpacing(spacing int) *PanelTheme {
	return &PanelTheme{
		BackgroundImage: p.BackgroundImage,
		Padding:         p.Padding,
		Spacing:         spacing,
	}
}
