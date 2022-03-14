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

// XWords implements interface X
type XWords struct {
	words    []string
	fallback interface{}
}

// NewXWords makes a new instance for XWords
func NewXWords(sourcePath string) X {
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
	return &XWords{
		words:    words,
		fallback: []string{},
	}
}

/*
SetFallback sets the value to be returned when the random function troubled.
*/
func (xw *XWords) SetFallback(fallback interface{}) {
	xw.fallback = fallback
}

/*
BindOperator returns the detected operator of XWords in the input condition:
"length (>=<) 4"
"begin = 'Marry'"
"end = 'Marry' "
*/
func (xw XWords) BindOperator(expression string) XOP {
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
func (xw XWords) Random(xExpression string) interface{} {
	op := xw.BindOperator(xExpression)
	switch op {
	case Any:
		return xw.any(1)
	case Length:
		return xw.randomWordsWithLength(xExpression)
	case Begin:
		return xw.randomWordsBeginWith(xExpression)
	case End:
		return xw.randomWordsEndWith(xExpression)
	case Invalid:
		return xw.fallback
	default:
		return xw.any(0)
	}
}

func (xw XWords) any(n int) []string {
	words := make([]string, n)
	source := rand.NewSource(time.Now().UnixNano())
	for i := range words {
		idx := rand.New(source).Intn(len(xw.words))
		words[i] = xw.words[idx]
	}
	return words
}

func (xw XWords) randomWordsWithLength(expression string) interface{} {
	xn := NewXNumber(0, 1000, 1)
	xn.SetFallback(-999)
	length := xn.Random(expression)
	if length == -999 {
		return xw.fallback
	}
	return xw.any(int(reflect.ValueOf(length).Float()))
}

func (xw XWords) randomWordsBeginWith(expression string) interface{} {
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		log.Println(err)
		return xw.fallback
	}
	if len(exp.Tokens()) != 3 {
		log.Println(fmt.Errorf("invalid expression: %v", expression))
		return xw.fallback
	}
	if exp.Tokens()[2].Kind != govaluate.STRING {
		log.Println(fmt.Errorf("unknown begin word: %v", exp.Tokens()[2].Value))
		return xw.fallback
	}
	return append([]string{reflect.ValueOf(exp.Tokens()[2].Value).String()}, xw.any(2)...)
}

func (xw XWords) randomWordsEndWith(expression string) interface{} {
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		log.Println(exp)
		return xw.fallback
	}
	if len(exp.Tokens()) != 3 {
		log.Println(fmt.Errorf("invalid expression: %v", expression))
		return xw.fallback
	}
	if exp.Tokens()[2].Kind != govaluate.STRING {
		log.Println(fmt.Errorf("unknown begin word: %v", exp.Tokens()[2].Value))
		return xw.fallback
	}
	return append(xw.any(2), reflect.ValueOf(exp.Tokens()[2].Value).String())
}
