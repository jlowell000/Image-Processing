package pixalfunctions

import (
	"image"
	"image/color"
	"runtime"

	"jlowell000.github.io/init/internal/actionengine"
)

const (
	highThreshold, lowThreshold float32 = 0.90, 0.10
)

/*GreyscaleHandle converts loadedImage image.Image to greyscale version*/
func GreyscaleHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImagePixel(loadedImage, greyscale, runtime.GOMAXPROCS(0))
}

func greyscale(p image.Point, imageOld image.Image) color.Color {
	return color.Gray16Model.Convert(imageOld.At(p.X, p.Y))
}

/*InvertColor inverts image colors*/
func InvertColor(loadedImage image.Image) image.Image {
	return actionengine.ActOnImagePixel(loadedImage, invertColor, runtime.GOMAXPROCS(0))
}
func invertColor(p image.Point, imageOld image.Image) color.Color {
	red, green, blue, alpha := imageOld.At(p.X, p.Y).RGBA()
	return color.RGBA64{uint16(alpha - red), uint16(alpha - green), uint16(alpha - blue), uint16(alpha)}
}

/*DoubleThresholdHandle removes low intensity pixals*/
func DoubleThresholdHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImagePixel(loadedImage, doubleThreshold, runtime.GOMAXPROCS(0))
}
func doubleThreshold(p image.Point, imageOld image.Image) color.Color {
	pixal := imageOld.At(p.X, p.Y)
	red, green, blue, alpha := pixal.RGBA()
	iToA := float32(red+green+blue) / float32(3*alpha)
	if iToA > highThreshold || iToA < lowThreshold {
		return color.RGBA64{0, 0, 0, uint16(alpha)}
	} else {
		return pixal
	}
}

/*FillInGapsHandle connects some lines to make more solid edges*/
func FillInGapsHandle(loadedImage image.Image) image.Image {
	return actionengine.ActOnImagePixel(loadedImage, fillInGaps, runtime.GOMAXPROCS(0))
}

func fillInGaps(p image.Point, imageOld image.Image) color.Color {

	pr, _, _, _ := imageOld.At(p.X, p.Y).RGBA()
	if pr > 0 {
		return imageOld.At(p.X, p.Y)
	}

	up := p.X-1 >= imageOld.Bounds().Min.X
	left := p.Y-1 >= imageOld.Bounds().Min.Y
	down := p.X+1 < imageOld.Bounds().Max.X
	right := p.Y+1 < imageOld.Bounds().Max.Y

	var sumR, countR uint32

	cacl := func(p image.Point, imageIn image.Image, sumRIn, countRIn uint32) (sumR, countR uint32) {
		r, _, _, _ := imageIn.At(p.X, p.Y-1).RGBA()
		sumR += r
		if r > 0 {
			countR++
		}
		return sumRIn, countRIn
	}

	if up {
		sumR, countR = cacl(image.Point{p.X - 1, p.Y}, imageOld, sumR, countR)
	}
	if down {
		sumR, countR = cacl(image.Point{p.X + 1, p.Y}, imageOld, sumR, countR)
	}
	if left {
		sumR, countR = cacl(image.Point{p.X, p.Y - 1}, imageOld, sumR, countR)
	}
	if right {
		sumR, countR = cacl(image.Point{p.X, p.Y + 1}, imageOld, sumR, countR)
	}
	if up && left {
		sumR, countR = cacl(image.Point{p.X - 1, p.Y - 1}, imageOld, sumR, countR)
	}
	if up && right {
		sumR, countR = cacl(image.Point{p.X - 1, p.Y + 1}, imageOld, sumR, countR)
	}
	if down && left {
		sumR, countR = cacl(image.Point{p.X + 1, p.Y - 1}, imageOld, sumR, countR)
	}
	if down && right {
		sumR, countR = cacl(image.Point{p.X + 1, p.Y + 1}, imageOld, sumR, countR)
	}
	if countR < 3 {
		return color.Gray16{0}
	}
	return color.Gray16{uint16(sumR / countR)}
}
