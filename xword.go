package randx

import (
	"encoding/csv"
	"github.com/Knetic/govaluate"
	"math/rand"
	"os"
	"strings"
	"time"
)

const (
	defaultWordSource = "/usr/share/dict/propernames" // Default proper names on a Unix machine
	length
	contain
	begin
)

//type XWord interface {
//	RandomWorlds(wordCount int) []string
//}

type XName struct {
	words []string
}

func NewXName(sourcePath string) (*XName, error) {
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
	return &XName{
		words: words,
	}, nil
}

func (xName XName) Random(x string, xExpression string) (interface{}, error) {
	switch x {
	case length:
		return xName.randomNWords("")
	}
	//if wordCount < 1 {
	//	return []string{}
	//}
	//
	//words := make([]string, wordCount)
	//for i := range words {
	//	idx := rand.New(source).Intn(len(gen.words))
	//	words[i] = gen.words[idx]
	//}
	//return words
}

func (xName XName) any(n int) []string {
	words := make([]string, n)
	source := rand.NewSource(time.Now().Unix())
	for i := range words {
		idx := rand.New(source).Intn(len(xName.words))
		words[i] = xName.words[idx]
	}
	return words
}

func (xName XName) randomNWords(expression string) ([]string, error) {
	var words []string
	if expression == "" {
		return xName.any(3), nil
	}
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}
	if len(exp.Tokens()) < 3 {
		return xName.any(3), nil
	}
	variable := exp.Tokens()[0]
	operator := exp.Tokens()[1]
	if variable.Value != length || operator.Kind != govaluate.COMPARATOR {
		return xName.any(3), nil
	}
	lengthConstantExp, err := govaluate.NewEvaluableExpressionFromTokens(exp.Tokens()[2:])
	if err != nil {
		return nil, err
	}
	expectedLength, err := lengthConstantExp.

	switch operator.Value {
	case govaluate.EQ:
		return
	}
	return words, nil
}

func (xName XName) randomNameContain() []string {

}

func (xName XName) randomNameBeginWith() []string {

}

//func (gen defaultXWord) RandomWorlds(wordCount int) []string {
//	source := rand.NewSource(time.Now().Unix())
//	if wordCount < 1 {
//		return []string{}
//	}
//
//	words := make([]string, wordCount)
//	for i := range words {
//		idx := rand.New(source).Intn(len(gen.words))
//		words[i] = gen.words[idx]
//	}
//	return words
//}
