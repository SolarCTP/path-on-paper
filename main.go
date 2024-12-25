package main

import (
	"bytes"
	"image"
	"image/png"
	"log"

	_ "embed"

	lv "github.com/SolarCTP/path-on-paper/levels"
	eb "github.com/hajimehoshi/ebiten/v2"
)

//go:embed logo.png
var logoRaw []byte

func main() {
	// game object creation
	game := &Game{
		lvl:   lv.NewLevelManager(),
		state: StateNotInLevel,
	}
	// load game icon
	logo, err := png.Decode(bytes.NewReader(logoRaw))
	if err != nil {
		log.Fatalln("Error loading icon:", err)
	}
	// initial game settings
	eb.SetWindowSize(1280, 720)
	eb.SetWindowTitle("Path on Paper")
	eb.SetTPS(120)
	eb.SetVsyncEnabled(true)
	eb.SetCursorMode(eb.CursorModeHidden)
	eb.SetWindowIcon([]image.Image{logo})

	game.lvl.LoadLevelByID(1, func() {
		game.state = StateBeforeStart
	}, nil)

	if err := eb.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
