package randx

import (
	"github.com/stretchr/testify/require"
	"log"
	"strings"
	"testing"
)

func TestDefaultWordGenerator(t *testing.T) {
	wordGenerator, err := NewXWord("")
	require.NoError(t, err)

	words0 := wordGenerator.RandomWorlds(0)
	require.Len(t, words0, 0)

	words1 := wordGenerator.RandomWorlds(1)
	require.Len(t, words1, 1)
	log.Println(strings.Join(words1, " "))

	words4 := wordGenerator.RandomWorlds(4)
	require.Len(t, words4, 4)
	log.Println(strings.Join(words4, " "))
}

func TestLocalizedWordGenerator(t *testing.T) {
	wordGenerator, err := NewXWord("etc/vietnamese")
	require.NoError(t, err)

	words0 := wordGenerator.RandomWorlds(0)
	require.Len(t, words0, 0)

	words1 := wordGenerator.RandomWorlds(1)
	require.Len(t, words1, 1)
	log.Println(strings.Join(words1, " "))

	words4 := wordGenerator.RandomWorlds(4)
	require.Len(t, words4, 4)
	log.Println(strings.Join(words4, " "))
}
