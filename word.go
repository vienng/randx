package randx

import (
	"encoding/csv"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	defaultWordSource = "/usr/share/dict/propernames" // Default proper names on a Unix machine
)

type XWord interface {
	RandomWorlds(wordCount int) []string
}

type defaultXWord struct {
	words []string
}

func NewXWord(sourcePath string) (XWord, error) {
	if sourcePath == "" {
		sourcePath = defaultWordSource
	}
	f, err := os.Open(sourcePath)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	words := make([]string, len(lines))
	for i, line := range lines {
		words[i] = strings.Join(line, " ")
	}
	return &defaultXWord{
		words: words,
	}, nil
}

func (gen defaultXWord) RandomWorlds(wordCount int) []string {
	source := rand.NewSource(time.Now().Unix())
	if wordCount < 1 {
		return []string{}
	}

	words := make([]string, wordCount)
	for i := range words {
		idx := rand.New(source).Intn(len(gen.words))
		words[i] = gen.words[idx]
	}
	return words
}
