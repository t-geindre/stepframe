package theme

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Theme struct {
	*widget.Theme
	PanelTheme    *PanelTheme
	IconsTheme    *IconsTheme
	MainMenuTheme *MainMenuTheme
}

type MainMenuTheme struct {
	ButtonImage   *widget.ButtonImage
	ButtonPadding *widget.Insets
	IconSpacing   int
	Font          *text.Face
	TextColor     color.Color
}

type IconsTheme struct {
	Add         *widget.GraphicImage
	Camera      *widget.GraphicImage
	CameraSmall *widget.GraphicImage
	Center      *widget.GraphicImage
	ChunkSmall  *widget.GraphicImage
	Delete      *widget.GraphicImage
	File        *widget.GraphicImage
	Fullscreen  *widget.GraphicImage
	Noise       *widget.GraphicImage
	NoiseSmall  *widget.GraphicImage
	Open        *widget.GraphicImage
	Save        *widget.GraphicImage
	StatsSmall  *widget.GraphicImage
	Zoom        *widget.GraphicImage
	ZoomSmall   *widget.GraphicImage

	NormalAccent Icons
	Normal       Icons
	Small        Icons
}

type PanelTheme struct {
	BackgroundImage *image.NineSlice
	ForegroundImage *image.NineSlice
	Padding         *widget.Insets
	Spacing         int
}
