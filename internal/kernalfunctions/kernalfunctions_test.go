package kernalfunctions

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_calculateGaussianKernal(t *testing.T) {
	k := int64(2)
	size := (2 * k) + 1
	expected := make([][]decimal.Decimal, size)
	tolarance := decimal.NewFromFloat(.012)
	mod := decimal.NewFromInt(273)
	for i, col := range [][]float64{
		{1, 4, 7, 4, 1},
		{4, 16, 26, 16, 4},
		{7, 26, 41, 26, 7},
		{4, 16, 26, 16, 4},
		{1, 4, 7, 4, 1},
	} {
		expected[i] = make([]decimal.Decimal, size)
		for j, v := range col {
			expected[i][j] = decimal.NewFromFloat(v).Div(mod)
		}
	}

	actual := calculateGaussianKernal(k, decimal.NewFromInt(1))
	matrixRadiallySymetric := func(arr [][]decimal.Decimal) bool {
		return assert.Equal(t, arr[0][0], arr[0][4], "{0,0},{0,4}") &&
			assert.Equal(t, arr[0][0], arr[4][0], "{0,0},{4,0}") &&
			assert.Equal(t, arr[0][0], arr[4][4], "{0,0},{4,4}") &&

			assert.Equal(t, arr[0][1], arr[0][3], "{0,1},{0,3}") &&
			assert.Equal(t, arr[0][1], arr[1][0], "{0,1},{1,0}") &&
			assert.Equal(t, arr[0][1], arr[1][4], "{0,1},{1,4}") &&
			assert.Equal(t, arr[0][1], arr[3][0], "{0,1},{3,0}") &&
			assert.Equal(t, arr[0][1], arr[3][4], "{0,1},{3,4}") &&
			assert.Equal(t, arr[0][1], arr[4][1], "{0,1},{4,1}") &&
			assert.Equal(t, arr[0][1], arr[4][3], "{0,1},{4,3}") &&

			assert.Equal(t, arr[0][2], arr[2][0], "{0,2},{2,0}") &&
			assert.Equal(t, arr[0][2], arr[2][4], "{0,2},{2,4}") &&
			assert.Equal(t, arr[0][2], arr[4][2], "{0,2},{4,2}") &&

			assert.Equal(t, arr[1][1], arr[1][3], "{1,1},{1,3}") &&
			assert.Equal(t, arr[1][1], arr[3][1], "{1,1},{3,1}") &&
			assert.Equal(t, arr[1][1], arr[3][3], "{1,1},{3,3}") &&

			assert.Equal(t, arr[1][2], arr[2][1], "{1,2},{2,1}") &&
			assert.Equal(t, arr[1][2], arr[2][3], "{1,2},{2,3}") &&
			assert.Equal(t, arr[1][2], arr[3][2], "{1,2},{3,2}")
	}
	assert.True(t, matrixRadiallySymetric(expected), "expected not radially symetric")
	assert.True(t, matrixRadiallySymetric(actual), "actual not radially symetric")

	assert.Equal(t, len(expected), len(actual), "array i dim do not match")
	for i, expectedCol := range expected {
		assert.Equal(t, len(expectedCol), len(actual[i]), "array j dim do not match")

		for j, expectedV := range expectedCol {
			assert.True(
				t,
				expectedV.Sub(tolarance).LessThan(actual[i][j]) && expectedV.Add(tolarance).GreaterThan(actual[i][j]),
				fmt.Sprintf("value [%d,%d] outside of tolarance ex: %s, ac: %s", i, j, expectedV.String(), actual[i][j].String()),
			)
		}
	}
}
