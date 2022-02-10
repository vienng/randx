package randx

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestRandomTimeEmptyInputs(t *testing.T) {
	layout := "2006-01-02"
	generator, err := NewXTime(
		time.Now().AddDate(-100, 0, 0),
		time.Now(),
		layout,
		24*time.Hour,
	)
	require.NoError(t, err)
	require.NotNil(t, generator)

	x, err := generator.RandomTime("", "")
	require.NoError(t, err)
	require.NotZero(t, x)

	y, err := generator.RandomTime("ssss", "")
	require.NoError(t, err)
	require.NotZero(t, y)
	z, err := generator.RandomTime("", "z > '2021-21-12'")
	require.Error(t, err)
	require.Zero(t, z)

	log.Println("x", x.Format(layout))
	log.Println("y", y.Format(layout))
	log.Println("z", z.Format(layout))
}

func TestRandomTimeNoVariable(t *testing.T) {
	layout := "2006-01-02"
	generator, err := NewXTime(
		time.Now().AddDate(-100, 0, 0),
		time.Now(),
		layout,
		24*time.Hour,
	)
	require.NoError(t, err)
	require.NotNil(t, generator)

	x, err := generator.RandomTime("", "'2013-01-02'")
	require.NoError(t, err)
	require.Equal(t, x.Format(layout), "2013-01-02")

	y, err := generator.RandomTime("", "'2014-01-01 23:59:59'")
	require.NoError(t, err)
	require.Equal(t, y.Format(layout), "2014-01-01")

	z, err := generator.RandomTime("", "'2014-01-01T15:04:05Z07:00'")
	require.Error(t, err)
	require.Zero(t, z)

	log.Println("x", x.Format(layout))
	log.Println("y", y.Format(layout))
	log.Println("z", z.Format(layout))
}

func TestRandomTimeSingleVariableExpression(t *testing.T) {
	layout := "2006-01-02"
	generator, err := NewXTime(
		time.Now().AddDate(-100, 0, 0),
		time.Now(),
		layout,
		24*time.Hour,
	)

	x, err := generator.RandomTime("x", "x < '2018-01-02'")
	require.NoError(t, err)

	y, err := generator.RandomTime("y", "y >= '2018-01-02'")
	require.NoError(t, err)

	z, err := generator.RandomTime("z", "z != '2021-01-02'")
	require.NoError(t, err)

	log.Println("x", x.Format(layout))
	log.Println("y", y.Format(layout))
	log.Println("z", z.Format(layout))
}

func TestRandomTimeMultipleVariableExpression(t *testing.T) {
	layout := "2006-01-02"
	generator, err := NewXTime(
		time.Now().AddDate(-100, 0, 0),
		time.Now(),
		layout,
		24*time.Hour,
	)

	x, err := generator.RandomTime("x", "x > '2000-01-01' && x < '2010-01-01'")
	require.NoError(t, err)

	y, err := generator.RandomTime("y", "y < '2000-01-01' || y > '2010-01-01'")
	require.NoError(t, err)

	z, err := generator.RandomTime("z", "z > '2000-01-01' || z > '2010-01-01' && z < '2020-01-01'")
	require.NoError(t, err)

	h, err := generator.RandomTime("h", "h < '2000-01-01' && h > '2010-01-01'")
	require.Error(t, err)

	log.Println("x", x.Format(layout))
	log.Println("y", y.Format(layout))
	log.Println("z", z.Format(layout))
	log.Println("h", h.Format(layout))
}
