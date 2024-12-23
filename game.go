package main

import (
	"image/color"
	"strconv"
	"time"

	lv "github.com/SolarCTP/path-on-paper/levels"
	eb "github.com/hajimehoshi/ebiten/v2"
	input "github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	PlayerRadius int = 1
)

type FPSCap struct {
	cap         int32
	maxInterval time.Duration
	lastDraw    time.Time
}

func NewFPSCap(cap int32) FPSCap {
	return FPSCap{
		cap:         cap,
		maxInterval: time.Second / time.Duration(cap),
		lastDraw:    time.Now(),
	}
}

type PlayState uint8

const (
	StateNotInLevel        PlayState = iota // there is no level loaded (e.g. while in menu)
	StateBeforeStart                        // before the cursor has been moved to the start pos
	StateStarted                            // after the cursor touched the start pos, game started
	StateGameOver                           // cursor has touched edges, game over
	StateTouchedFinishArea                  // the cursor has touched the finish area and won
)

// FrameTooEarly returns true if the time passed since the last
// draw is less than the interval specified by the frame cap
func (f *FPSCap) FrameTooEarly() bool {
	return time.Since(f.lastDraw) < f.maxInterval
}

type Game struct {
	lvl   *lv.LevelManager
	fps   FPSCap
	state PlayState
	font  *text.GoXFace
}

func (g *Game) Update() error {
	if input.IsKeyJustPressed(eb.KeyQ) {
		return eb.Termination
	}

	// DEBUG: reset game state
	if input.IsKeyJustPressed(eb.KeyR) {
		g.state = StateNotInLevel
		if g.lvl.LoadLevelByID(1) {
			g.state = StateBeforeStart
		}
	}

	cursorPos := lv.XYtoPoint(eb.CursorPosition())
	switch g.state {
	case StateNotInLevel:
	case StateBeforeStart:
		if g.lvl.ActiveLevel.TouchingStartPos(cursorPos) {
			g.state = StateStarted
		}
	case StateStarted:
		if g.lvl.ActiveLevel.TouchingEdge(PlayerRadius, cursorPos) {
			g.state = StateGameOver
		} else if g.lvl.ActiveLevel.TouchingFinishArea(cursorPos) {
			g.state = StateTouchedFinishArea
		}
	}

	return nil
}

func (g *Game) Draw(screen *eb.Image) {
	// Cap the FPS by skipping frames that want to be drawn too early
	if g.fps.FrameTooEarly() {
		return
	}

	cursorX, cursorY := eb.CursorPosition()
	vector.DrawFilledCircle(
		screen, float32(cursorX), float32(cursorY),
		float32(PlayerRadius), color.Black, false,
	)

	opts := &eb.DrawImageOptions{}
	screen.DrawImage(g.lvl.ActiveLevel.Img, opts)

	txtOpts := &text.DrawOptions{}
	txtOpts.GeoM.Translate(50, 50)
	txtOpts.GeoM.Scale(3, 3)
	text.Draw(screen, strconv.Itoa(int(g.state)), g.font, txtOpts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1920, 1080
}
