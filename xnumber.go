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

// XNumber implements interface X, XNumber returns a random number satisfied inputted condition
type XNumber struct {
	min    float64
	max    float64
	step   float64
	errorN float64
}

// NewXNumber creates a new instance for XNumber
func NewXNumber(min, max, step float64) X {
	return &XNumber{min, max, step, min - step}
}

// BindOperator returns supported operator of XNumber
func (xn XNumber) BindOperator(expression string) XOP {
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
	return Unknown
}

func uniqVars(slice []string) []string {
	var uniqSlice []string
	for _, s := range slice {
		if !strings.Contains(strings.Join(uniqSlice, " "), s) {
			uniqSlice = append(uniqSlice, s)
		}
	}
	return uniqSlice
}

// Random generates a random number x that satisfies given conditions
func (xn XNumber) Random(condition string) (interface{}, error) {
	kind := xn.BindOperator(condition)
	switch kind {
	case Any:
		return randomNumber([][]float64{{xn.min, xn.max}}), nil
	case Constant:
		return xn.constantNumber(condition)
	case FindX:
		return xn.findSingleX(condition)
	case FindXs:
		return xn.findMultipleX(condition)
	default:
		return xn.errorN, fmt.Errorf("invalid or not supported expression")
	}
}

func (xn XNumber) constantNumber(expStr string) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(expStr)
	if err != nil {
		return xn.errorN, err
	}
	result, err := expression.Evaluate(nil)
	if err != nil {
		return xn.errorN, err
	}
	return reflect.ValueOf(result).Float(), nil
}

func (xn XNumber) findSingleX(expStr string) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(expStr)
	if err != nil {
		return xn.errorN, err
	}
	numberRange, err := xn.findExpressionRange(expression)
	if err != nil {
		return xn.errorN, err
	}
	return randomNumber([][]float64{numberRange}), nil
}

func (xn XNumber) findMultipleX(expStr string) (interface{}, error) {
	expression, err := govaluate.NewEvaluableExpression(expStr)
	if err != nil {
		return xn.errorN, err
	}
	expressions, err := splitToSingleExpressions(expression)
	if err != nil {
		return xn.errorN, err
	}

	var allRanges [][]float64
	for _, exp := range expressions {
		r, err := xn.findExpressionRange(exp)
		if err != nil {
			return xn.errorN, err
		}
		allRanges = append(allRanges, r)
	}
	finalRanges, err := evaluateRanges(allRanges, expression)
	if err != nil {
		return xn.errorN, err
	}
	return randomNumber(finalRanges), nil
}

/*
findExpressionRange returns the possible range of X.
Input of this function must be a SINGLE X expression. In case of multiple X, use splitToSingleExpressions first.

E.g. if the XNumber is initialized with min=0,max=1000,step=1
"x > 0" --> possible range is [1,1000]
"x < 100" --> possible range is [0,99]
*/
func (xn XNumber) findExpressionRange(exp *govaluate.EvaluableExpression) (expRange []float64, err error) {
	tokens := exp.Tokens()[2:] // ignore variable and first operator
	newExp, err := govaluate.NewEvaluableExpressionFromTokens(tokens)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	val, err := newExp.Evaluate(nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	operator := exp.Tokens()[1]
	number := reflect.ValueOf(val).Float()
	switch operator.Value {
	case "==":
		expRange = []float64{number, number}
	case "<":
		expRange = []float64{xn.min, number - xn.step}
	case "<=":
		expRange = []float64{xn.min, number}
	case ">", "!=":
		expRange = []float64{number + xn.step, xn.max}
	case ">=":
		expRange = []float64{number, xn.max}
	default:
		expRange = []float64{xn.min, xn.max}
	}
	return expRange, nil
}

/*
splitToSingleExpression splits an expression with multiple x into the smaller ones.
the separator when splitting is the govaluate.LOGICALOP
*/
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
			number := reflect.ValueOf(token.Value).Float()
			tokens = append(tokens, fmt.Sprintf("%f", number))
		default:
			tokens = append(tokens, fmt.Sprint(token.Value))
		}
	}
	return expressions, nil
}

/*
randomNumber picks a random number in every given ranges [min1, max1] [min2, max2]
then pick a random element from picked slices.
*/
func randomNumber(numberRanges [][]float64) float64 {
	randNumbers := make([]float64, len(numberRanges))
	source := rand.NewSource(time.Now().UnixNano())
	for i, r := range numberRanges {
		randNumbers[i] = r[0] + rand.New(source).Float64()*(r[1]-r[0])
	}
	return randNumbers[rand.New(source).Intn(len(randNumbers))]
}

/*
evaluateRanges returns the final ranges of X.
"x > 0 && x < 100" --> possible ranges [1,1000] && [0,99] --> final range is [1,99]
"x < 100 || x > 200" --> possible ranges [1,99] || [201,1000] --> final ranges are [1,99] [201,1000]
*/
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
