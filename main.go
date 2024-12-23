package main

import (
	"log"

	lv "github.com/SolarCTP/path-on-paper/levels"
	"github.com/SolarCTP/path-on-paper/ui"
	eb "github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// game object creation
	game := &Game{
		fps:   NewFPSCap(240),
		lvl:   lv.NewLevelManager(),
		state: StateNotInLevel,
		font:  ui.MainFont,
	}
	// initial game settings
	eb.SetWindowSize(1280, 720)
	eb.SetWindowTitle("Path on Paper")
	eb.SetTPS(120)
	eb.SetVsyncEnabled(false)

	if game.lvl.LoadLevelByID(1) {
		game.state = StateBeforeStart
	} else {
		log.Fatal("Level load failed")
	}

	if err := eb.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
