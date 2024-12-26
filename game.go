package main

import (
	"image/color"
	"log"
	"strconv"

	lv "github.com/SolarCTP/path-on-paper/levels"
	"github.com/SolarCTP/path-on-paper/ui"
	eb "github.com/hajimehoshi/ebiten/v2"
	input "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	PlayerRadius     float32 = 3  // used in edge collision calculation
	PlayerAuraRadius float32 = 10 // used in edge collision calculation

	AntiCheatMaxMouseDelta float32 = 200

	// Window resolution used in Game.Layout()
	LogicalWinResX int = 1920
	LogicalWinResY int = 1080
)

var (
	PlayerColor     color.RGBA = color.RGBA{255, 255, 255, 255}
	PlayerAuraColor color.RGBA = color.RGBA{196, 66, 49, 10}
)

type PlayState uint8

const (
	StateNotInLevel        PlayState = iota // there is no level loaded (e.g. while in menu)
	StateBeforeStart                        // before the cursor has been moved to the start pos
	StateStarted                            // after the cursor touched the start pos, game started
	StateGameOver                           // cursor has touched edges, game over
	StateTouchedFinishArea                  // the cursor has touched the finish area and won
)

type Game struct {
	lvl           *lv.LevelManager
	state         PlayState
	settings      Settings
	lastCursorPos lv.Point
}

func (g *Game) Update() error {
	if input.IsKeyJustPressed(eb.KeyQ) {
		return eb.Termination
	}

	ui.Manager.EUI.Update()

	// DEBUG: reset game state
	if input.IsKeyJustPressed(eb.KeyR) {
		g.state = StateNotInLevel
		g.lvl.LoadLevelByID(g.lvl.ActiveLevel.ID, func() {
			g.state = StateBeforeStart
		}, func() {})
	}

	// toggle fullscreen
	if input.IsKeyJustPressed(eb.KeyF11) {
		g.settings.Fullscreen = !g.settings.Fullscreen
		eb.SetFullscreen(g.settings.Fullscreen)
	}

	// end here if there is no active level (otherwise nil dereference)
	if g.lvl.ActiveLevel == nil {
		return nil
	}

	// DEBUG: switch between levels with arrow keys
	checkLevelSwitchKeys := func() {
		if input.IsKeyJustPressed(eb.KeyArrowLeft) {
			g.lvl.LoadPrevAvailableLevel()
			g.state = StateBeforeStart
		} else if input.IsKeyJustReleased(eb.KeyArrowRight) {
			g.lvl.LoadNextAvailableLevel()
			g.state = StateBeforeStart
		}
	}

	cursorPos := lv.XYtoPoint(eb.CursorPosition())
	switch g.state {
	case StateNotInLevel:
	case StateBeforeStart:
		checkLevelSwitchKeys()

		// start the game once the player moves the cursor to the start area
		if g.lvl.ActiveLevel.TouchingStartPos(cursorPos) {
			g.state = StateStarted
		}
	case StateStarted:
		// if g.lvl.ActiveLevel.TouchingEdge(int(PlayerRadius), cursorPos) {
		// 	g.state = StateGameOver
		antiCheatLargeDistance := lv.Dist(cursorPos, g.lastCursorPos) > float64(AntiCheatMaxMouseDelta)
		if g.lvl.ActiveLevel.TouchingEdgeV2(int(PlayerRadius), cursorPos, g.lastCursorPos) ||
			antiCheatLargeDistance {
			if antiCheatLargeDistance {
				log.Println("Cheating detected. You might be using a touchscreen or drawing tablet.")
			}
			g.state = StateGameOver
		} else if g.lvl.ActiveLevel.TouchingFinishArea(cursorPos) {
			g.state = StateTouchedFinishArea
		}
	case StateTouchedFinishArea:
		checkLevelSwitchKeys()
	}

	g.lastCursorPos.X, g.lastCursorPos.Y = eb.CursorPosition()
	return nil
}

func (g *Game) Draw(screen *eb.Image) {
	// if g.lvl.ActiveLevel == nil {
	// 	screen.Fill(color.Black)
	// 	text.Draw(screen, "Loading...", ui.MainFontWithSize(100),
	// 		ui.DefaultTxtOptsAt(200, 200))
	// 	return
	// }

	opts := &eb.DrawImageOptions{}
	screen.DrawImage(g.lvl.ActiveLevel.Img, opts)

	// Draw the cursor (a pixel surrounded by a circle)
	cursorX, cursorY := eb.CursorPosition()
	// vector.DrawFilledCircle(
	// 	screen, float32(cursorX), float32(cursorY),
	// 	PlayerAuraRadius, PlayerAuraColor, true,
	// )
	vector.DrawFilledCircle(
		screen, float32(cursorX), float32(cursorY),
		PlayerRadius, PlayerColor, true,
	)

	// DEBUG text: Game state, current level ID, FPS
	text.Draw(screen, g.gameStateText(), ui.MainFontWithSize(40),
		ui.DefaultTxtOptsAt(20, float64(LogicalWinResY)-50))
	text.Draw(screen, strconv.Itoa(int(g.lvl.ActiveLevel.ID)), ui.MainFontWithSize(40),
		ui.DefaultTxtOptsAt(float64(LogicalWinResX)-30-text.Advance(
			strconv.Itoa(LogicalWinResX), ui.MainFontWithSize(40),
		), float64(LogicalWinResY)-50))
	text.Draw(screen, strconv.Itoa(int(eb.ActualFPS())), ui.MainFontWithSize(40),
		ui.DefaultTxtOptsAt(float64(LogicalWinResX)-100, float64(LogicalWinResY)-100))

	// Draw the UI on top of everything else
	ui.Manager.EUI.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return LogicalWinResX, LogicalWinResY
}

func (g *Game) gameStateText() string {
	switch g.state {
	case StateBeforeStart:
		return "Game not started"
	case StateStarted:
		return "Playing"
	case StateGameOver:
		return "Game over"
	case StateTouchedFinishArea:
		return "You win!"
	case StateNotInLevel:
		return "N/A"
	}
	return "<UNKNOWN>"
}
