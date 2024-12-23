package levels

import (
	"image"
	"math"
)

// a pair of X, Y int coordinates
type Point = image.Point

func XYtoPoint(x, y int) Point {
	return image.Point{X: x, Y: y}
}

// Distance between points
func Dist(p1, p2 Point) float64 {
	return math.Sqrt(math.Pow(float64(p2.X)-float64(p1.X), 2) +
		math.Pow(float64(p2.Y)-float64(p1.Y), 2))
}
