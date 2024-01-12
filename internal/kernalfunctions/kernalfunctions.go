package kernalfunctions

import (
	"image"
	"image/color"
	"math"
	"runtime"

	"jlowell000.github.io/init/internal/actionengine"
)

const (
	sigma = 1.6
)

/*GaussianHandle runs Gaussian Filter over image*/
func GaussianHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImageKernal(loadedImage, gaussianFilter, 7, runtime.GOMAXPROCS(0)/5)
}

func gaussianFilter(p image.Point, imageOld image.Image, kernalSize int) [][]color.Color {

	k := (kernalSize - 1) / 2
	p = boundsKernal(p, imageOld, kernalSize)

	value := float64(0)
	forKernal(kernalSize, func(i, j int) {
		aa := math.Pow(float64(i-(k+1)), 2)
		bb := math.Pow(float64(j-(k+1)), 2)
		cc := 2 * math.Pow(sigma, 2)
		dd := 2 * math.Pi * math.Pow(sigma, 2)

		h := (1 / dd) * math.Exp(-((aa + bb) / cc))

		r, _, _, _ := imageOld.At(p.X+i, p.Y+j).RGBA()
		value += float64(r) * h
	})

	output := make([][]color.Color, kernalSize)
	for i := range output {
		output[i] = make([]color.Color, kernalSize)
		for j := range output[i] {
			output[i][j] = color.Gray16{uint16(value)}
		}
	}
	return output
}

/*SobelHandle runs Sobel kernal over image*/
func SobelHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImageKernal(loadedImage, sobelFilter, 3, runtime.GOMAXPROCS(0)/5)
}
func sobelFilter(p image.Point, imageOld image.Image, kernalSize int) [][]color.Color {

	gradientGValue := func(p image.Point, imageOld image.Image, kernalSize int) (gX, gY float64) {

		kernalX := [][]float64{
			{1, 2, 1},
			{0, 0, 0},
			{-1, -2, -1}}
		kernalY := [][]float64{
			{1, 0, -1},
			{2, 0, -2},
			{1, 0, -1}}
		// gX, gY := float64(0), float64(0)
		forKernal(kernalSize, func(i, j int) {
			r, _, _, _ := imageOld.At(p.X+i, p.Y+j).RGBA()

			gX += kernalX[i][j] * float64(r)
			gY += kernalY[i][j] * float64(r)
		})
		return gX, gY
	}

	pointsToCheck := func(theta float64, p image.Point) (posPoint, negPoint image.Point) {
		step := float64(45 / 2)
		if step <= theta && theta < 45.0+step {
			return image.Point{p.X + 1, p.Y + 1}, image.Point{p.X - 1, p.Y - 1}
		} else if 45.0+step <= theta && theta < 90.0+step {
			return image.Point{p.X + 1, p.Y}, image.Point{p.X + 1, p.Y}
		} else if 90.0+step <= theta && theta < 135.0+step {
			return image.Point{p.X + 1, p.Y - 1}, image.Point{p.X - 1, p.Y + 1}
		}
		return image.Point{p.X, p.Y + 1}, image.Point{p.X, p.Y - 1}
	}

	p = boundsKernal(p, imageOld, kernalSize)

	output := make([][]color.Color, kernalSize)
	gX, gY := gradientGValue(p, imageOld, kernalSize)

	for i := range output {
		output[i] = make([]color.Color, kernalSize)
		for j := range output[i] {

			g := math.Sqrt(math.Pow(gX, 2) + math.Pow(gY, 2))
			posP, negP := pointsToCheck(math.Atan2(gY, gX), p)
			testGX1, testGY1 := gradientGValue(posP, imageOld, kernalSize)
			testGX2, testGY2 := gradientGValue(negP, imageOld, kernalSize)
			if math.Sqrt(math.Pow(testGX1, 2)+math.Pow(testGY1, 2)) < g && math.Sqrt(math.Pow(testGX2, 2)+math.Pow(testGY2, 2)) < g {
				output[i][j] = color.Gray16{uint16(g)}
			} else {
				output[i][j] = color.Gray16{uint16(0)}
			}
		}
	}
	return output
}

func boundsKernal(p image.Point, imageIn image.Image, kernalSize int) image.Point {

	if p.X < imageIn.Bounds().Min.X {
		p.X = imageIn.Bounds().Min.X
	} else if p.X+kernalSize >= imageIn.Bounds().Max.X {
		p.X = imageIn.Bounds().Max.X - kernalSize
	}
	if p.Y < imageIn.Bounds().Min.Y {
		p.Y = imageIn.Bounds().Min.Y
	} else if p.Y+kernalSize >= imageIn.Bounds().Max.Y {
		p.Y = imageIn.Bounds().Max.Y - kernalSize
	}
	return p
}

func forKernal(kernalSize int, action func(int, int)) {
	for i := 0; i < kernalSize; i++ {
		for j := 0; j < kernalSize; j++ {
			action(i, j)
		}
	}
}
