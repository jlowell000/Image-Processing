package pixalfunctions

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/jlowell000/utils"
	"github.com/stretchr/testify/assert"
)

var (
	testImage     image.Image     = &TestImage{}
	testMax       int             = 9
	testMinPoint  image.Point     = image.Point{X: 0, Y: 0}
	testMaxPoint  image.Point     = image.Point{X: testMax, Y: testMax}
	testRectangle image.Rectangle = image.Rectangle{Min: testMinPoint, Max: testMaxPoint}
)

func Test_greyscale(t *testing.T) {
	utils.ForSquare(testMax+1, func(i, j int) {
		expected := color.Gray16Model.Convert(testImage.At(i, j))
		actual := greyscale(image.Point{X: i, Y: j}, testImage)

		assert.Equal(t, expected, actual, fmt.Sprintf("test intensity %d, %d expected: %d; actual: %d", i, j, expected, actual))

		red, green, blue, _ := expected.RGBA()
		intensity := red + green + blue
		assert.Equal(t, red*3, intensity, fmt.Sprintf("test intensity %d, %d expected: %d; actual: %d; color: %d", i, j, red, intensity, expected))
		assert.Equal(t, blue*3, intensity, fmt.Sprintf("test intensity %d, %d expected: %d; actual: %d; color: %d", i, j, blue, intensity, expected))
		assert.Equal(t, green*3, intensity, fmt.Sprintf("test intensity %d, %d expected: %d; actual: %d; color: %d", i, j, green, intensity, expected))
	})
}

func Test_invertColor(t *testing.T) {
	utils.ForSquare(testMax+1, func(i, j int) {

		red, green, blue, alpha := testImage.At(i, j).RGBA()
		expected := color.RGBA64{uint16(alpha - red), uint16(alpha - green), uint16(alpha - blue), uint16(alpha)}
		actual := invertColor(image.Point{X: i, Y: j}, testImage)

		assert.Equal(t, expected, actual, fmt.Sprintf("test intensity %d, %d expected: %d; actual: %d", i, j, expected, actual))
	})
}

func Test_doubleThreshold(t *testing.T) {
	utils.ForSquare(testMax+1, func(i, j int) {
		pixal := testImage.At(i, j)
		red, green, blue, alpha := pixal.RGBA()
		iToA := float32(red+green+blue) / float32(3*alpha)
		var expected color.Color
		if iToA > highThreshold || iToA < lowThreshold {
			expected = color.RGBA64{0, 0, 0, uint16(alpha)}
		} else {
			expected = pixal
		}

		actual := doubleThreshold(image.Point{X: i, Y: j}, testImage)
		assert.Equal(t, expected, actual, fmt.Sprintf("test intensity %d, %d expected: %d; actual: %d", i, j, expected, actual))
	})
}

type TestImage struct{}

func (i *TestImage) ColorModel() color.Model {
	return color.RGBA64Model
}

func (i *TestImage) Bounds() image.Rectangle {
	return testRectangle
}

func (i *TestImage) At(x, y int) color.Color {
	return color.RGBA{
		R: scale(x), G: scale(y), B: scale((x * y) / 2), A: 255,
	}
}

func scale(i int) uint8 {
	return uint8((255 * i) / testMax)
}
