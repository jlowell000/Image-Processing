package actionengine

import (
	"image"
	"image/color"

	"github.com/jlowell000/utils"
)

/*
 * ActOnImagePixel runs provided function on image. Acts with a worker pool
 */
func ActOnImagePixel(imageOld image.Image, pixelAction func(image.Point, image.Image) color.Color) *image.RGBA {
	imageNew := MakeRGBASpace(imageOld)
	utils.ForEachWG(
		ImageToArray(imageOld),
		func(p image.Point) {
			imageNew.Set(p.X, p.Y, pixelAction(p, imageOld))
		},
	)
	return imageNew
}

/*
 * ActOnImageKernal runs provided function on image. Acts a worker pool
 */
func ActOnImageKernal(imageOld image.Image, kernalAction func(p image.Point, imageOld image.Image, kernalSize int) [][]color.Color, kernalSize int) *image.RGBA {
	imageNew := MakeRGBASpace(imageOld)
	utils.ForEachWG(
		ImageToArray(imageOld),
		func(p image.Point) {
			kernalResult := kernalAction(p, imageOld, kernalSize)
			utils.ForSquare(kernalSize, func(x, y int) { imageNew.Set(x+p.X, y+p.Y, kernalResult[x][y]) })
		},
	)
	return imageNew
}

/*
 * Convert input image, a 2D array of points into a 1D array of points.
 */
func ImageToArray(input image.Image) (result []image.Point) {
	for y := input.Bounds().Min.Y; y < input.Bounds().Max.Y; y++ {
		for x := input.Bounds().Min.X; x < input.Bounds().Max.X; x++ {
			result = append(result, image.Point{X: x, Y: y})
		}
	}
	return
}

/*
 * Makes an RGBA space the same dimensions as the input
 */
func MakeRGBASpace(input image.Image) *image.RGBA {
	return image.NewRGBA(
		image.Rect(
			input.Bounds().Min.X,
			input.Bounds().Min.Y,
			input.Bounds().Max.X,
			input.Bounds().Max.Y,
		),
	)
}
