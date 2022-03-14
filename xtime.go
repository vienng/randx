package randx

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"reflect"
	"strings"
	"time"
)

// XTime implements interface X, XTime returns a random datetime satisfied inputted condition
type XTime struct {
	min      time.Time
	max      time.Time
	step     time.Duration
	fallback interface{}
}

// NewXTime makes a new instance for XTime
func NewXTime(min, max time.Time, step time.Duration) X {
	return &XTime{
		min:      min,
		max:      max,
		step:     step,
		fallback: time.Time{},
	}
}

func (xt *XTime) SetFallback(value interface{}) {
	xt.fallback = value
}

// BindOperator returns supported operator of XTime
func (xt XTime) BindOperator(expression string) XOP {
	if len(expression) == 0 {
		return Any
	}
	exp, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		log.Println(err)
		return Invalid
	}

	if len(exp.Vars()) == 0 {
		return Constant
	}
	if len(exp.Vars()) == 1 && len(exp.Tokens()) >= 3 {
		return FindX
	}
	if len(exp.Vars()) > 1 && len(uniqVars(exp.Vars())) == 1 {
		return FindXs
	}
	return Invalid
}

// Random returns a random datetime satisfied inputted condition
func (xt XTime) Random(expression string) interface{} {
	op := xt.BindOperator(expression)
	switch op {
	case Any:
		return xt.randomUnixTime("")
	case Constant, FindX, FindXs:
		newExp, err := xt.toTimestamp(expression)
		if err != nil {
			log.Println(err)
			return xt.fallback
		}
		return xt.randomUnixTime(newExp)
	default:
		return xt.fallback
	}
}

// randomUnixTime treats the timestamp as an int64 number. randomUnixTime uses XNumber to random a timestamp
func (xt XTime) randomUnixTime(exp string) interface{} {
	xNumber := NewXNumber(float64(xt.min.Unix()), float64(xt.max.Unix()), xt.step.Seconds())
	xNumber.SetFallback(-999)
	value := xNumber.Random(exp)
	if value == -999 {
		return xt.fallback
	}
	randomTimestamp := int64(reflect.ValueOf(value).Float())
	return time.Unix(randomTimestamp, 0)
}

// toTimestamp coverts the detected datetime (any format) into timestamp
func (xt XTime) toTimestamp(expStr string) (string, error) {
	exp, err := govaluate.NewEvaluableExpression(expStr)
	if err != nil {
		log.Printf("error create expression: %v", err)
		return "", err
	}

	elems := make([]string, len(exp.Tokens()))
	for i, token := range exp.Tokens() {
		if token.Kind == govaluate.TIME {
			constTimeExp, err := govaluate.NewEvaluableExpressionFromTokens([]govaluate.ExpressionToken{token})
			if err != nil {
				return "", fmt.Errorf("error create expression from: %v", token)
			}
			result, err := constTimeExp.Evaluate(nil)
			if err != nil {
				return "", fmt.Errorf("error evaluate expression: %v", constTimeExp)
			}
			switch result.(type) {
			case float64:
				elems[i] = fmt.Sprintf("%f", reflect.ValueOf(result).Float())
			default:
				return "", fmt.Errorf("unknown time format: %v", result)
			}
		} else {
			elems[i] = fmt.Sprint(token.Value)
		}
	}
	return strings.Join(elems, " "), nil
}
