package ui

import (
	"bytes"
	_ "embed"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

var mainFontSource *text.GoTextFaceSource = nil

//go:embed Roboto-Regular.ttf
var robotoRegular []byte

var DefaultTextOpts text.DrawOptions = text.DrawOptions{}

func MainFontWithSize(size float64) *text.GoTextFace {
	return &text.GoTextFace{
		Source:   mainFontSource,
		Size:     size,
		Language: language.Italian,
	}
}

// DefaultTxtOptsAt returns the default text options with the specified translation applied.
//
// If you need other controls, copy DefaultTextOpts (don't reference it directly)
// and apply your changes
func DefaultTxtOptsAt(x, y float64) *text.DrawOptions {
	DefaultTextOpts.GeoM.Reset()
	DefaultTextOpts.GeoM.Translate(x, y)
	return &DefaultTextOpts
}

func init() {
	setupFont()
	setupDefaultTextOpts()
}

func setupFont() {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(robotoRegular))
	if err != nil {
		log.Fatalln("Could not load font:", err)
	}
	mainFontSource = fontSource
}

func setupDefaultTextOpts() {
	DefaultTextOpts.Filter = 1 // Linear filtering
	DefaultTextOpts.ColorScale.ScaleWithColor(color.Black)
}
