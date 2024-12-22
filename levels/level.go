package levels

import (
	"image"
	"math"
)

type Point = image.Point

// Distance between points
func Dist(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(float64(p2.X)-float64(p1.X), 2) +
		math.Pow(float64(p2.Y)-float64(p1.Y), 2))
}

type LevelID uint64

type LevelData struct {
	edges []Point // the black pixels in the level's image
}

type Level struct {
	ID   LevelID
	data *LevelData
}

func (l *Level) load(data *LevelData) bool {
	panic("unimplemented")
}

func (l *Level) unload() {
	panic("unimplemented")
}

func (l *Level) TouchingEdge(playerRadius int, playerPos Point) bool {
	// octrees would be better, but I don't fucking know what they are
	for _, edge := range l.data.edges {
		if Dist(edge, playerPos) <= float64(playerRadius) {
			return true
		}
	}
	return false
}

func NewLevel(newId LevelID) *Level {
	return &Level{
		ID:   newId,
		data: nil,
	}
}
