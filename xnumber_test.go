package randx

import (
	"github.com/Knetic/govaluate"
	"github.com/stretchr/testify/require"
	"log"
	"reflect"
	"testing"
)

func TestXNumber(t *testing.T) {
	generator := NewXNumber(0, 1000, 1)

	testCases := []struct {
		name               string
		condition          string
		isValidExpression  bool
		expectedExpression string
	}{
		{"any", "", true, "x >= 0 && x <= 1000"},
		{"constant1", "10", true, "x == 10"},
		{"constant2", "2+3", true, "x == 5"},
		{"constant3", "4*(2+3)", true, "x == 20"},
		{"expression1", "x == 100", true, "x == 100"},
		{"expression2", "x > 100", true, "x > 100 && x < 1000"},
		{"expression3", " x <    100", true, "x >= 0 && x < 100"},
		{"expression4", "x >= 2*3", true, "x >= 6 && x < 1000"},
		{"expression5", "x <= 2*3+5", true, "x >= 0 && x <= 11"},
		{"expression6", "x != 100", true, "x > 100 && x < 1000"},
		{"expression7", "x > 0 && x < 100", true, "x > 0 && x < 100"},
		{"expression8", "x < 100 || x > 200", true, "x < 100 || x > 200"},
		{"expression9", "x > 5 || x > 15 && x < 100", true, "x > 5 || x > 15 && x < 100"},
		{"expression10", "x < 1389200400 || x > 1483894800", true, "x < 1389200400 || x > 1483894800"},
		{"invalid1", "x < 100 &&", false, "x == -1"},
		{"invalid2", "x <", false, "x == -1"},
		{"invalid3", "< 100", false, "x == -1"},
		{"invalid4", "x", false, "x == -1"},
		{"invalid5", "x < 100 && x > 200", false, "x == -1"},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			x := generator.Random(tc.condition)
			exp, err := govaluate.NewEvaluableExpression(tc.expectedExpression)
			require.NoError(t, err)

			params := make(map[string]interface{}, 1)
			params["x"] = x
			result, err := exp.Evaluate(params)
			require.NoError(t, err)
			require.True(t, reflect.ValueOf(result).Bool(), "unsatisfied x", exp.String(), x)
			log.Println("[INFO]", tc.name, tc.condition, "result ->", x)
		})
	}
}
