package levels

import (
	"image"
	"image/color"
)

func FindBlackPixels(img image.Image) []image.Point {
	var blackPixels []image.Point
	bounds := img.Bounds()

	// Define the black color to compare against.
	black := color.RGBA{0, 0, 0, 255}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if img.At(x, y) == black {
				blackPixels = append(blackPixels, image.Point{X: x, Y: y})
			}
		}
	}
	return blackPixels
}
