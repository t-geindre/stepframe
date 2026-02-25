package theme

import (
	"bytes"
	_ "embed"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// SRC: https://fonts.google.com/specimen/Roboto
//
//go:embed font.ttf
var fontSrc []byte

//go:embed font-bold.ttf
var fontBoldSrc []byte

var fontCache = make(map[string]*text.Face)

func getFontFace(size float64, bold bool) *text.Face {
	key := strconv.Itoa(int(size)) + "_" + strconv.FormatBool(bold)
	if face, ok := fontCache[key]; ok {
		return face
	}

	ftSrc := fontSrc
	if bold {
		ftSrc = fontBoldSrc
	}

	src, err := text.NewGoTextFaceSource(bytes.NewReader(ftSrc))
	if err != nil {
		panic(err)
	}

	var face text.Face
	face = &text.GoTextFace{
		Source: src,
		Size:   size,
	}

	fontCache[key] = &face

	return &face
}
