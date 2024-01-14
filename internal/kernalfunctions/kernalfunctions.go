package kernalfunctions

import (
	"image"
	"image/color"
	"math"

	"github.com/jlowell000/utils"
	"github.com/shopspring/decimal"
	"jlowell000.github.io/init/internal/actionengine"
)

var (
	radius int64               = 5
	sigma  float64             = math.Max(float64(radius)/2.0, 1.0)
	kernal [][]decimal.Decimal = calculateGaussianKernal(radius, decimal.NewFromFloat(sigma))
)

/*GaussianHandle runs Gaussian Filter over image*/
func GaussianHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImagePixel(loadedImage, gaussianFilter)
}

func gaussianFilter(p image.Point, imageIn image.Image) color.Color {
	if pointKernalInBounds(p, imageIn, radius) {
		red, green, blue := decimal.Zero, decimal.Zero, decimal.Zero
		for i := -radius; i <= radius; i++ {
			for j := radius; j <= radius; j++ {
				x, y := p.X+int(i), p.Y+int(j)
				r, b, g, _ := imageIn.At(x, y).RGBA()
				k := kernal[i+radius][j+radius]
				red = red.Add(decimal.NewFromInt32(int32(r))).Mul(k)
				blue = blue.Add(decimal.NewFromInt32(int32(b))).Mul(k)
				green = green.Add(decimal.NewFromInt32(int32(g))).Mul(k)
				// alpha += float64(a) * k
			}
		}
		_, _, _, a := imageIn.At(p.X, p.Y).RGBA()
		return color.RGBA{
			R: uint8(red.IntPart()),
			B: uint8(blue.IntPart()),
			G: uint8(green.IntPart()),
			A: uint8(a),
		}

	}
	return GetPointColor(p, imageIn)
}

func calculateGaussianKernal(radius int64, sigma decimal.Decimal) [][]decimal.Decimal {
	kernalWidth := int64((2 * radius) + 1)
	kernal := make([][]decimal.Decimal, kernalWidth)
	for i := range kernal {
		kernal[i] = make([]decimal.Decimal, kernalWidth)
	}
	sum := decimal.Zero

	for i := int64(1); i <= kernalWidth; i++ {
		for j := int64(1); j <= kernalWidth; j++ {
			kernal[i-1][j-1] = gaussianKernalValue(
				decimal.NewFromInt(i),
				decimal.NewFromInt(j),
				decimal.NewFromInt(radius),
				sigma,
			)
			sum = sum.Add(kernal[i-1][j-1])
		}
	}
	for i, kCol := range kernal {
		for j, kValue := range kCol {
			kernal[i][j] = kValue.Div(sum)
		}
	}
	return kernal
}

/*
 * Assumes 1 <= (x,y) <= (2 * radius) + 1
 */
func gaussianKernalValue(x, y, radius, sigma decimal.Decimal) decimal.Decimal {
	one := decimal.NewFromInt(1)
	two := decimal.NewFromInt(2)
	minusOne := decimal.NewFromInt(-1)
	pi := decimal.NewFromFloat(math.Pi)

	exNum := x.Sub(radius.Add(one)).Pow(two).Add(y.Sub(radius.Add(one)).Pow(two)).Mul(minusOne)
	exDom := sigma.Pow(two).Mul(two)

	ex, _ := exNum.Div(exDom).ExpTaylor(3)
	dom := pi.Mul(exDom)
	return ex.Div(dom)
}

/*SobelHandle runs Sobel kernal over image*/
func SobelHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImageKernal(loadedImage, sobelFilter, 3)
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
		utils.ForSquare(kernalSize, func(i, j int) {
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

func pointKernalInBounds(p image.Point, imageIn image.Image, radius int64) bool {
	bounds := imageIn.Bounds()
	r := int(radius)
	return p.X-r < bounds.Min.X ||
		p.X+r < bounds.Max.X ||
		p.Y-r < bounds.Min.Y ||
		p.Y+r < bounds.Max.Y
}

func GetPointColor(p image.Point, imageIn image.Image) color.RGBA {
	red, green, blue, alpha := imageIn.At(p.X, p.Y).RGBA()
	return color.RGBA{
		R: uint8(red),
		B: uint8(green),
		G: uint8(blue),
		A: uint8(alpha),
	}
}
