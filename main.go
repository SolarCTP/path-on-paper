package main

import (
	"log"

	lv "github.com/SolarCTP/path-on-paper/levels"
	eb "github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// game object creation
	game := &Game{
		lvl:   lv.NewLevelManager(),
		state: StateNotInLevel,
	}
	// initial game settings
	eb.SetWindowSize(1280, 720)
	eb.SetWindowTitle("Path on Paper")
	eb.SetTPS(120)
	eb.SetVsyncEnabled(true)
	eb.SetCursorMode(eb.CursorModeHidden)

	game.lvl.LoadLevelByID(1, func() {
		game.state = StateBeforeStart
	}, nil)

	if err := eb.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
