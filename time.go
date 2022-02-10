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
		min:    fromDate,
		max:    toDate,
		layout: layout,
		step:   step,
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

	newExp, err := gen.toTimestamp(exp)
	if err != nil {
		return time.Time{}, err
	}
	expressionStr, err := exportExpressionStringFromTokens(newExp.Tokens())
	if err != nil {
		return time.Time{}, err
	}

	randTimestamp, err := xNumber.RandomNumber(variable, expressionStr)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(int64(randTimestamp), 0), nil
}

func (gen *defaultXTime) toTimestamp(exp *govaluate.EvaluableExpression) (*govaluate.EvaluableExpression, error) {
	tokens := make([]govaluate.ExpressionToken, len(exp.Tokens()))
	for i, token := range exp.Tokens() {
		if token.Kind == govaluate.TIME {
			constTimeExp, err := govaluate.NewEvaluableExpressionFromTokens([]govaluate.ExpressionToken{token})
			if err != nil {
				return nil, fmt.Errorf("error create expression from: %v", token)
			}
			result, err := constTimeExp.Evaluate(nil)
			if err != nil {
				return nil, fmt.Errorf("error evaluate expression: %v", constTimeExp)
			}
			switch result.(type) {
			case float64:
				tokens[i] = govaluate.ExpressionToken{
					Kind:  govaluate.NUMERIC,
					Value: result,
				}
			case string:
				tm, err := time.Parse(gen.layout, reflect.ValueOf(result).String())
				if err != nil {
					return nil, fmt.Errorf("error parse time: %v", result)
				}
				tokens[i] = govaluate.ExpressionToken{
					Kind:  govaluate.NUMERIC,
					Value: tm.Unix(),
				}
			default:
				return nil, fmt.Errorf("unknown format: %v", result)
			}
		} else {
			tokens[i] = token
		}
	}
	return govaluate.NewEvaluableExpressionFromTokens(tokens)
}
