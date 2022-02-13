package randx

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"log"
	"math"
	"math/rand"
	"reflect"
	"strings"
	"time"
)

type XNumber interface {
	RandomNumber(variable, condition string) (float64, error)
}

type defaultXNumber struct {
	min  float64
	max  float64
	step float64
}

func NewXNumber(min, max, step float64) XNumber {
	return &defaultXNumber{min, max, step}
}

func (gen *defaultXNumber) RandomNumber(variable, condition string) (float64, error) {
	// default
	if condition == "" {
		return randomNumber([][]float64{{gen.min, gen.max}}), nil
	}

	expression, err := govaluate.NewEvaluableExpression(condition)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	// constant expression
	if len(expression.Vars()) == 0 {
		result, err := expression.Evaluate(nil)
		if err != nil {
			return 0, err
		}
		return reflect.ValueOf(result).Float(), nil
	}

	// invalid find-x expressions
	if len(expression.Vars()) != 0 && len(variable) == 0 {
		return 0, fmt.Errorf("not indicated variable: %s", expression.String())
	}

	for _, v := range expression.Vars() {
		if v != variable {
			return 0, fmt.Errorf("variable '%s' not found: %s", variable, expression.String())
		}
	}

	// single find-x expression
	if isSingleFindXExpression(expression) {
		numberRange, err := gen.findExpressionRange(expression)
		if err != nil {
			return 0, err
		}
		return randomNumber([][]float64{numberRange}), nil
	}

	// multiple find-x expression
	expressions, err := splitToSingleExpressions(expression)
	if err != nil {
		return 0, err
	}

	var allRanges [][]float64
	for _, exp := range expressions {
		r, err := gen.findExpressionRange(exp)
		if err != nil {
			return 0, err
		}
		allRanges = append(allRanges, r)
	}
	finalRanges, err := evaluateRanges(allRanges, expression)
	if err != nil {
		return 0, err
	}
	return randomNumber(finalRanges), nil
}

/*
isSingleFindXExpression  checks if the expression has a variable and must be followed the rules.
-- e.g. TRUE --
x < 2
x >= 2 + 3 + 4

-- e.g. FALSE --
x > 0 && x < 10 // multiple
x - 2 < 5 // not allowed
3 = 9 - x
*/

func isSingleFindXExpression(exp *govaluate.EvaluableExpression) bool {
	if len(exp.Tokens()) < 3 {
		log.Printf("invalid expression: %s", exp.String())
		return false
	}

	if len(exp.Vars()) != 1 {
		return false
	}

	if exp.Tokens()[0].Kind != govaluate.VARIABLE || exp.Tokens()[1].Kind != govaluate.COMPARATOR {
		log.Printf("invalid expression: %s", exp.String())
		return false
	}

	return true
}

func (gen *defaultXNumber) findExpressionRange(exp *govaluate.EvaluableExpression) (expRange []float64, err error) {
	tokens := exp.Tokens()[2:] // ignore variable and first operator
	newExp, err := govaluate.NewEvaluableExpressionFromTokens(tokens)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	number, err := newExp.Evaluate(nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	operator := exp.Tokens()[1]
	switch fmt.Sprint(operator.Value) {
	case fmt.Sprint(govaluate.LT):
		expRange = []float64{gen.min, reflect.ValueOf(number).Float() - gen.step}
	case fmt.Sprint(govaluate.LTE):
		expRange = []float64{gen.min, reflect.ValueOf(number).Float()}
	case fmt.Sprint(govaluate.GT), fmt.Sprint(govaluate.NEQ):
		expRange = []float64{reflect.ValueOf(number).Float() + gen.step, gen.max}
	case fmt.Sprint(govaluate.GTE):
		expRange = []float64{reflect.ValueOf(number).Float(), gen.max}
	default:
		expRange = []float64{gen.min, gen.max}
	}
	return expRange, nil
}

func splitToSingleExpressions(exp *govaluate.EvaluableExpression) ([]*govaluate.EvaluableExpression, error) {
	processTokens := append([]govaluate.ExpressionToken{}, exp.Tokens()...)
	processTokens = append(processTokens, govaluate.ExpressionToken{Kind: govaluate.LOGICALOP})

	var expressions []*govaluate.EvaluableExpression
	var tokens []string
	for _, token := range processTokens {
		switch token.Kind {
		case govaluate.LOGICALOP:
			subExp, err := govaluate.NewEvaluableExpression(strings.Join(tokens, " "))
			if err != nil {
				log.Println(err)
				return nil, err
			}
			expressions = append(expressions, subExp)
			tokens = nil
		case govaluate.NUMERIC:
			tokens = append(tokens, fmt.Sprintf(" %f", token.Value))
		default:
			tokens = append(tokens, fmt.Sprint(token.Value))
		}
	}
	return expressions, nil
}

func randomNumber(numberRanges [][]float64) float64 {
	randNumbers := make([]float64, len(numberRanges))
	source := rand.NewSource(time.Now().UnixNano())
	for i, r := range numberRanges {
		randNumbers[i] = r[0] + rand.New(source).Float64()*(r[1]-r[0])
	}
	return randNumbers[rand.New(source).Intn(len(randNumbers))]
}

func evaluateRanges(ranges [][]float64, exp *govaluate.EvaluableExpression) ([][]float64, error) {
	if len(ranges) <= 1 {
		return ranges, nil
	}

	var logicalOperators []govaluate.ExpressionToken
	for _, token := range exp.Tokens() {
		if token.Kind == govaluate.LOGICALOP {
			logicalOperators = append(logicalOperators, token)
		}
	}
	if len(logicalOperators) != len(ranges)-1 {
		return nil, fmt.Errorf("invalid expression: %s", exp.String())
	}
	var finalRanges [][]float64
	for i := 0; i < len(ranges)-1; i++ {
		operator := reflect.ValueOf(logicalOperators[i].Value).String()

		for j := i + 1; j < len(ranges); j++ {
			minI := ranges[i][0]
			maxI := ranges[i][1]
			minJ := ranges[j][0]
			maxJ := ranges[j][1]

			switch operator {
			case fmt.Sprint(govaluate.AND):
				r := []float64{math.Max(minI, minJ), math.Min(maxI, maxJ)}
				if r[1]-r[0] >= 0 {
					finalRanges = [][]float64{r}
				}
			case fmt.Sprint(govaluate.OR):
				finalRanges = append(finalRanges, ranges[i], ranges[j])
			default:
				return nil, fmt.Errorf("unsupported operator: %s", operator)
			}
		}
	}
	if len(finalRanges) == 0 {
		return nil, fmt.Errorf("invalid expression: %s", exp.String())
	}
	return finalRanges, nil
}
