package ui

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var MainFont *text.GoXFace = nil

const mainFontPath string = "./ui/Roboto-Regular.ttf"

func init() {
	setupFont()
}

func setupFont() {
	fontData, err := os.ReadFile(mainFontPath)
	if err != nil {
		log.Print("Error while reading font file: ")
		log.Fatalln(err)
	}
	parsedFont, err2 := opentype.Parse(fontData)
	if err != nil {
		log.Print("Error while parsing font file: ")
		log.Fatalln(err2)
	}
	face, err3 := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Print("Error while creating font face: ")
		log.Fatalln(err3)
	}
	MainFont = text.NewGoXFace(face)
}
