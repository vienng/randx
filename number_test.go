package randx

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestRandomNumberEmptyInputs(t *testing.T) {
	generator := NewXNumber(0, 1, 0.01)

	x, err := generator.RandomNumber("", "")
	require.NoError(t, err)
	require.NotZero(t, x)

	y, err := generator.RandomNumber("ssss", "")
	require.NoError(t, err)
	require.NotZero(t, y)
	require.NotEqual(t, x, y)

	z, err := generator.RandomNumber("", "z > 0")
	require.Error(t, err)
	require.Zero(t, z)

	log.Println("x", x)
	log.Println("y", y)
	log.Println("z", z)
}

func TestRandomNumberNoVariable(t *testing.T) {
	generator := NewXNumber(0, 100000, 5)

	x, err := generator.RandomNumber("", "2 + 3")
	require.NoError(t, err)
	require.EqualValues(t, x, 5)

	y, err := generator.RandomNumber("", "10")
	require.NoError(t, err)
	require.EqualValues(t, y, 10)

	z, err := generator.RandomNumber("xxx", "4*(2+3)")
	require.NoError(t, err)
	require.EqualValues(t, z, 20)

	log.Println("x", x)
	log.Println("y", y)
	log.Println("z", z)
}

func TestVariableDetection(t *testing.T) {
	generator := NewXNumber(0, 100000, 1)

	x, err := generator.RandomNumber("x", "X > 0")
	require.Error(t, err)
	require.Zero(t, x, "must be equal")

	y, err := generator.RandomNumber("y", "yy > 0")
	require.Error(t, err)
	require.Zero(t, y, "must be equal")

	z, err := generator.RandomNumber("z", " z>0")
	require.NoError(t, err)
	require.True(t, z > 0)

	log.Println("x", x)
	log.Println("y", y)
	log.Println("z", z)
}

func TestRandomNumberSingleVariableExpression(t *testing.T) {
	generator := NewXNumber(0, 100000, 1)

	x, err := generator.RandomNumber("x", "x < 2")
	require.NoError(t, err)
	require.True(t, x < 2)

	y, err := generator.RandomNumber("y", "y >= 2*3")
	require.NoError(t, err)
	require.True(t, y >= 6)

	z, err := generator.RandomNumber("z", "z != 100")
	require.NoError(t, err)
	require.True(t, z != 100)

	log.Println("x", x)
	log.Println("y", y)
	log.Println("z", z)
}

func TestRandomNumberMultipleVariableExpression(t *testing.T) {
	generator := NewXNumber(0, 100000, 1)

	x, err := generator.RandomNumber("x", "x > 0 && x < 100")
	require.NoError(t, err)

	y, err := generator.RandomNumber("y", "y < 100 || y > 200")
	require.NoError(t, err)

	z, err := generator.RandomNumber("z", "z > 5 || z > 15 && z < 100")
	require.NoError(t, err)

	h, err := generator.RandomNumber("h", "h < 100 && h > 200")
	require.Error(t, err)

	log.Println("x", x)
	log.Println("y", y)
	log.Println("z", z)
	log.Println("h", h)
}
