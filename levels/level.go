package levels

import (
	"image/png"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	IMAGE_FORMAT string = "png"
	LEVELS_PATH  string = "./levels/level_files/"
)

type LevelID uint64

type LevelData struct {
	edges        []Point
	startPosArea []Point
	finishArea   []Point
}

type Level struct {
	ID   LevelID
	data *LevelData
	Img  *ebiten.Image
}

// TouchingEdge reports whether the player is currently less than playerRadius away
// from any point in the edges region
func (l *Level) TouchingEdge(playerRadius int, playerPos Point) bool {
	// octrees would be better, but I don't fucking know what they are
	for _, edge := range l.data.edges {
		if Dist(edge, playerPos) <= float64(playerRadius) {
			return true
		}
	}
	return false
}

// TouchingEdgeV2 reports whether the player has crossed a point between
// lastPos and currentPos such that its distance to any edge point is
// smaller than playerRadius (basically the other one but less prone to cheating).
//
// This is more expensive than V1. TODO divide edge pixels into chunks to
// avoid having to check every single one every frame
func (l *Level) TouchingEdgeV2(playerRadius int, currentPos, lastPos Point) bool {
	if lastPos.Eq(Point{X: 0, Y: 0}) {
		return false
	}

	m := float64(currentPos.Y-lastPos.Y) / float64(currentPos.X-lastPos.X)
	y := float64(lastPos.Y)
	for x := math.Min(
		float64(lastPos.X), float64(currentPos.X),
	); x <= math.Max(
		float64(lastPos.X), float64(currentPos.X),
	); x++ {
		y += m
		roundedY := int(math.Round(y))
		for _, p := range l.data.edges {
			if Dist(p, XYtoPoint(int(x), roundedY)) <= float64(playerRadius) {
				return true
			}
		}
	}

	return false
}

// TouchingFinishArea reports whether the player is touching any point
// inside the finish area
func (l *Level) TouchingFinishArea(playerPos Point) bool {
	for _, finishPoint := range l.data.finishArea {
		if playerPos == finishPoint {
			return true
		}
	}
	return false
}

func (l *Level) TouchingStartPos(playerPos Point) bool {
	for _, startPoint := range l.data.startPosArea {
		if playerPos == startPoint {
			return true
		}
	}
	return false
}

// Load loads the level and returns whether it succeded.
// Optionally, you can use it concurrently by passing a
// non-nil bool channel, on which success or failure
// will be sent
func (l *Level) Load(loadProgressOrNil chan bool) bool {
	// quick function to report success or failure to
	// the loadProgress channel if it is valid
	chanReport := func(s bool) {
		if loadProgressOrNil != nil {
			loadProgressOrNil <- s
			close(loadProgressOrNil)
		}
	}

	if l.ID == 0 {
		log.Println("Level was not properly initialized, since its ID is 0. " +
			"Could not load level")
		chanReport(false)
		return false
	}

	// TODO: cache result and checksum of image inside a binary file, to avoid
	// repeating the image analysis at each level load

	// load the image file
	levelImageFile, fileErr := os.Open(
		LEVELS_PATH + strconv.FormatUint(uint64(l.ID), 10) + "." + IMAGE_FORMAT, // e.g. .../1.png
	)
	if fileErr != nil {
		log.Println(fileErr)
		chanReport(false)
		return false
	}

	// read the image file inside an image.Image
	levelImage, decodeErr := png.Decode(levelImageFile)
	if decodeErr != nil {
		log.Println(fileErr)
		chanReport(false)
		return false
	}
	levelImageFile.Close() // not needed anymore

	pixelMap := getPixelMap(levelImage)

	// check that all areas actually exist
	if !verifyPixelMap(pixelMap) {
		chanReport(false)
		return false
	}

	// get edges (black pixels)
	l.data.edges = append(l.data.edges, pixelMap[PIXELS_EDGES]...)

	// set start pos (calculated in place)
	l.data.startPosArea = append(l.data.startPosArea, pixelMap[PIXELS_STARTPOS]...)

	// get finish area (green pixels)
	l.data.finishArea = append(l.data.finishArea, pixelMap[PIXELS_FINISHAREA]...)

	clear(pixelMap)

	// save the image to draw it during game execution
	l.Img = ebiten.NewImageFromImage(levelImage)

	// report back to caller, and channel if called as thread
	chanReport(true)
	return true
}

// Unload clears the level data
func (l *Level) Unload() {
	clear(l.data.edges)
	clear(l.data.finishArea)
	l.Img = nil
	l.data = nil // shoutout to the GC
}

// NewLevel returns the level with the specified ID, unloaded
func NewLevel(newId LevelID) *Level {
	newLvl := &Level{
		ID: newId,
		data: &LevelData{
			edges:      nil,
			finishArea: nil,
		},
		Img: nil,
	}
	return newLvl
}
