package randx

import (
	"encoding/csv"
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"time"
)

const (
	defaultWordSource = "/usr/share/dict/propernames" // Default proper names on a Unix machine
)

// XWord implements interface X
type XWord struct {
	words []string
}

// NewXWord makes a new instance for XWord
func NewXWord(sourcePath string) X {
	if sourcePath == "" {
		sourcePath = defaultWordSource
	}
	f, err := os.Open(sourcePath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	words := make([]string, len(lines))
	for i, line := range lines {
		words[i] = strings.Join(line, " ")
	}
	return &XWord{
		words: words,
	}
}

/*
BindOperator returns the detected operator of XWord in the input condition:
"length (>=<) 4"
"begin = 'Marry'"
"end = 'Marry' "
*/
func (xw XWord) BindOperator(expression string) XOP {
	if len(expression) == 0 {
		return Any
	}
	lwExp := strings.ToLower(expression)
	if strings.Contains(lwExp, "length") {
		return Length
	}
	if strings.Contains(lwExp, "begin") {
		return Begin
	}
	if strings.Contains(lwExp, "end") {
		return End
	}
	return Invalid
}

// Random returns random words matching input condition
func (xw XWord) Random(xExpression string) (interface{}, error) {
	op := xw.BindOperator(xExpression)
	switch op {
	case Any:
		return xw.any(1), nil
	case Length:
		return xw.randomWordsWithLength(xExpression)
	case Begin:
		return xw.randomWordsBeginWith(xExpression)
	case End:
		return xw.randomWordsEndWith(xExpression)
	case Invalid:
		return xw.any(0), fmt.Errorf("invalid expression %s", xExpression)
	default:
		return xw.any(0), nil
	}
}

func (xw XWord) any(n int) []string {
	words := make([]string, n)
	source := rand.NewSource(time.Now().UnixNano())
	for i := range words {
		idx := rand.New(source).Intn(len(xw.words))
		words[i] = xw.words[idx]
	}
	return words
}

func (xw XWord) randomWordsWithLength(expression string) ([]string, error) {
	xn := NewXNumber(0, 1000, 1)
	length, err := xn.Random(expression)
	if err != nil {
		return nil, err
	}
	return xw.any(int(reflect.ValueOf(length).Float())), nil
}

func (xw XWord) randomWordsBeginWith(expression string) ([]string, error) {
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}
	if len(exp.Tokens()) != 3 {
		return nil, fmt.Errorf("invalid expression: %v", expression)
	}
	if exp.Tokens()[2].Kind != govaluate.STRING {
		return nil, fmt.Errorf("unknown begin word: %v", exp.Tokens()[2].Value)
	}
	return append([]string{reflect.ValueOf(exp.Tokens()[2].Value).String()}, xw.any(2)...), nil
}

func (xw XWord) randomWordsEndWith(expression string) ([]string, error) {
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return nil, err
	}
	if len(exp.Tokens()) != 3 {
		return nil, fmt.Errorf("invalid expression: %v", expression)
	}
	if exp.Tokens()[2].Kind != govaluate.STRING {
		return nil, fmt.Errorf("unknown begin word: %v", exp.Tokens()[2].Value)
	}
	return append(xw.any(2), reflect.ValueOf(exp.Tokens()[2].Value).String()), nil
}
