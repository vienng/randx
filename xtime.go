package randx

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"reflect"
	"strings"
	"time"
)

type XTime struct {
	min  time.Time
	max  time.Time
	step time.Duration
}

func NewXTime(min, max time.Time, step time.Duration) (X, error) {
	return &XTime{
		min: min,
		max: max,
		//layout: layout,
		step: step,
	}, nil
}

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

func (xt XTime) Random(expression string) (interface{}, error) {
	op := xt.BindOperator(expression)
	switch op {
	case Any:
		return xt.randomTime("")
	case Constant, FindX, FindXs:
		newExp, err := xt.toTimestamp(expression)
		if err != nil {
			return nil, err
		}
		return xt.randomTime(newExp)
	default:
		return time.Time{}, fmt.Errorf("invalid expression %s", expression)
	}
}

func (xt XTime) randomTime(exp string) (interface{}, error) {
	xNumber := NewXNumber(float64(xt.min.Unix()), float64(xt.max.Unix()), xt.step.Seconds())
	value, err := xNumber.Random(exp)
	if err != nil {
		return time.Time{}, err
	}
	randomTimestamp := int64(reflect.ValueOf(value).Float())
	return time.Unix(randomTimestamp, 0), nil
}

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
			//case string:
			//	tm, err := time.Parse(xt.layout, reflect.ValueOf(result).String())
			//	if err != nil {
			//		return "", fmt.Errorf("error parse time: %v", result)
			//	}
			//	elems[i] = fmt.Sprintf("%d", tm.Unix())
			default:
				return "", fmt.Errorf("unknown time format: %v", result)
			}
		} else {
			elems[i] = fmt.Sprint(token.Value)
		}
	}
	log.Println("DEVVVV", strings.Join(elems, " "))
	return strings.Join(elems, " "), nil
}
