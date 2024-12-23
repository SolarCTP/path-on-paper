package levels

import (
	"image/png"
	"log"
	"os"
	"strconv"
)

const (
	IMAGE_FORMAT string = "png"
	LEVELS_PATH  string = "./level_files"
)

type LevelID uint64

type LevelData struct {
	edges      []Point // the black pixels in the level's image
	startPos   Point
	finishArea []Point
}

type Level struct {
	ID   LevelID
	data *LevelData
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

// setStartPos sets the startPos as the average of the coords
// of blue pixels
func (l *Level) setStartPos(startPosArea []Point) {
	for _, px := range startPosArea {
		l.data.startPos = l.data.startPos.Add(px)
	}
	l.data.startPos = l.data.startPos.Div(len(startPosArea))
}

// Load loads the level and returns whether it succeded.
// Optionally, you can use it concurrently by passing a
// non-nil bool channel, on which success or failure
// will be sent
func (l *Level) Load(loadProgress chan bool) bool {
	// quick function to report success or failure to
	// the loadProgress channel if it is valid
	chanReport := func(s bool) {
		if loadProgress != nil {
			loadProgress <- s
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
		strconv.FormatUint(uint64(l.ID), 10) + "." + IMAGE_FORMAT, // e.g. 1.png
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

	colorPxCoordsMap := findPixelsByColors(levelImage,
		COLOR_BLACK, COLOR_GREEN, COLOR_BLUE)

	// get edges (black pixels)
	l.data.edges = colorPxCoordsMap[COLOR_BLACK]

	// set start pos (calculated in place)
	l.setStartPos(colorPxCoordsMap[COLOR_BLUE])

	// get finish area (green pixels)
	l.data.finishArea = colorPxCoordsMap[COLOR_GREEN]

	chanReport(true)
	return true
}

// Unload clears the level data
func (l *Level) Unload() {
	clear(l.data.edges)
	clear(l.data.finishArea)
	l.data = nil // shoutout to the GC
}

// NewLevel returns the level with the specified ID, unloaded
func NewLevel(newId LevelID) *Level {
	newLvl := &Level{
		ID:   newId,
		data: nil,
	}
	return newLvl
}
