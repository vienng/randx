package randx

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"reflect"
	"time"
)

type XTime interface {
	RandomTime(variable string, condition string) (time.Time, error)
}

type defaultXTime struct {
	min    time.Time
	max    time.Time
	layout string
	step   time.Duration
}

func NewXTime(min, max time.Time, layout string, step time.Duration) (XTime, error) {
	fromDate, err := time.Parse(layout, min.Format(layout))
	if err != nil {
		return nil, err
	}
	toDate, err := time.Parse(layout, max.Format(layout))
	if err != nil {
		return nil, err
	}
	return &defaultXTime{
		min:  fromDate,
		max:  toDate,
		step: step,
	}, nil
}

func (gen *defaultXTime) RandomTime(variable string, condition string) (time.Time, error) {
	xNumber := NewXNumber(float64(gen.min.Unix()), float64(gen.max.Unix()), gen.step.Seconds())

	if condition == "" {
		randTimestamp, err := xNumber.RandomNumber("", "")
		if err != nil {
			return time.Time{}, err
		}
		return time.Unix(int64(randTimestamp), 0), nil
	}

	exp, err := govaluate.NewEvaluableExpression(condition)
	if err != nil {
		log.Printf("error create expression: %v", err)
		return time.Time{}, err
	}

	if len(exp.Vars()) == 0 {
		result, err := exp.Evaluate(nil)
		if err != nil {
			log.Printf("error evaluate expression: %v", err)
			return time.Time{}, err
		}
		switch result.(type) {
		case float64:
			return time.Unix(int64(reflect.ValueOf(result).Float()), 0), nil
		case string:
			return time.Parse(gen.layout, reflect.ValueOf(result).String())
		default:
			return time.Time{}, fmt.Errorf("unknown format: %v", result)
		}
	}

	expressionStr, err := toTimestamp(exp)
	if err != nil {
		return time.Time{}, err
	}

	randTimestamp, err := xNumber.RandomNumber(variable, expressionStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(randTimestamp), 0), nil
}

func toTimestamp(exp *govaluate.EvaluableExpression) (string, error) {
	expression := ""
	for _, token := range exp.Tokens() {
		switch token.Kind {
		case govaluate.NUMERIC:
			expression = expression + fmt.Sprintf(" %f", token.Value)
		case govaluate.TIME:
			constTimeExp, err := govaluate.NewEvaluableExpressionFromTokens([]govaluate.ExpressionToken{token})
			if err != nil {
				return "", fmt.Errorf("error create expression from: %v", token)
			}
			result, err := constTimeExp.Evaluate(nil)
			expression = expression + fmt.Sprintf(" %f", reflect.ValueOf(result).Float())
		default:
			expression = expression + fmt.Sprintf(" %s", token.Value)
		}
	}
	return expression, nil
}
