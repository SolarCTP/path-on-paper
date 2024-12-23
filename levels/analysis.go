package levels

import (
	"image"
	"image/color"
)

// colors used in level analysis
var (
	COLOR_GREEN = color.RGBA{0, 255, 0, 255}
	COLOR_BLUE  = color.RGBA{0, 0, 255, 255}
	COLOR_BLACK = color.RGBA{0, 0, 0, 255}
)

func findPixelsByColors(img image.Image, colors ...color.RGBA) map[color.RGBA][]Point {
	var colorPixelMap = make(map[color.RGBA][]Point)
	bounds := img.Bounds()

	// Define the black color to compare against.

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for _, c := range colors {
				if img.At(x, y) == c {
					colorPixelMap[c] = append(colorPixelMap[c], image.Point{X: x, Y: y})
				}
			}
		}
	}
	return colorPixelMap
}
