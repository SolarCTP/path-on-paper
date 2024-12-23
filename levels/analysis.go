package levels

import (
	"image"
	"image/color"
	"log"
)

// If you add more types, make sure to update the slice
// in findPixelsByColors()
const (
	PIXELS_STARTPOS = iota
	PIXELS_EDGES
	PIXELS_FINISHAREA
)

// ptToColor is the map from pixel type to pixel color
func ptToColor(pt uint) color.RGBA {
	switch pt {
	case PIXELS_STARTPOS:
		return color.RGBA{0, 0, 255, 255}
	case PIXELS_EDGES:
		return color.RGBA{0, 0, 0, 255}
	case PIXELS_FINISHAREA:
		return color.RGBA{0, 255, 0, 255}
	default:
		log.Fatal("Invalid pixel type", pt)
	}
	return color.RGBA{0, 0, 0, 0}
}

// func colorPixelType(c color.RGBA) uint {
// 	switch c {
// 	case color.RGBA{0, 0, 255, 255}:
// 		return PIXELS_STARTPOS
// 	case color.RGBA{0, 0, 0, 255}:
// 		return PIXELS_EDGES
// 	case color.RGBA{0, 255, 0, 255}:
// 		return PIXELS_FINISHAREA
// 	default:
// 		log.Fatal("Invalid color", c)
// 	}
// 	return 999
// }

// getPixelMap analyzes a level image pixel by pixel and returns a
// map of (pixel type, pixel slice) pairs. Each pixel type has its
// own color (defined in pixelTypeColor()) and purpose.
func getPixelMap(img image.Image) map[uint][]Point {
	var pixelMap = make(map[uint][]Point)
	pixelTypes := []uint{
		PIXELS_STARTPOS,
		PIXELS_EDGES,
		PIXELS_FINISHAREA,
	}

	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for _, pt := range pixelTypes {
				if img.At(x, y) == ptToColor(pt) {
					pixelMap[pt] = append(pixelMap[pt], image.Point{X: x, Y: y})
				}
			}
		}
	}
	return pixelMap
}

// verifyPixelMap checks that every component is present in the pixel map,
// and returns true if all of them are present, or false otherwise (and logs
// what is missing to the console)
func verifyPixelMap(pMap map[uint][]Point) bool {
	_, okEdges := pMap[PIXELS_EDGES]
	_, okStartPos := pMap[PIXELS_STARTPOS]
	_, okFinishArea := pMap[PIXELS_FINISHAREA]

	if !(okEdges && okStartPos && okFinishArea) {
		log.Println("The level is missing the following parts (or analysis failed for):")
		if !okEdges {
			log.Println("- Edges")
		}
		if !okEdges {
			log.Println("- Starting position (blue)")
		}
		if !okEdges {
			log.Println("- Finish area (green)")
		}
		return false
	}
	return true
}
