package theme

import (
	"bytes"
	_ "embed"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// SRC: https://fonts.google.com/specimen/Roboto
//
//go:embed font.ttf
var fontSrc []byte

func getFontFace(size float64) *text.Face {
	src, err := text.NewGoTextFaceSource(bytes.NewReader(fontSrc))
	if err != nil {
		panic(err)
	}

	var face text.Face
	face = &text.GoTextFace{
		Source: src,
		Size:   size,
	}

	return &face
}
