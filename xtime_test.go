package randx

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestXTime(t *testing.T) {
	startTime := time.Now().AddDate(-100, 0, 0)
	endTime := time.Now()
	step := 24 * time.Hour
	generator := NewXTime(startTime, endTime, step)

	testCases := []struct {
		name              string
		condition         string
		isValidExpression bool
	}{
		{"any", "", true},
		{"constant1", "'2018-01-01'", true},
		//{"constant2", "2018-01-01", false}, TODO
		{"constant3", "'2018-01-01 19:36:43'", true},
		{"constant4", "1646644467", false},
		{"constant5", "'1646644467'", true},
		{"expression1", "x == '2018-01-01 19:36:43'", true},
		{"expression2", "x > '2014-01-09'", true},
		{"expression3", " x <    '2014-01-09'", true},
		{"expression4", "x >= 1646644467", false},
		{"expression5", "x >= '1646644467'", true},
		{"expression6", "x != '2014-01-09'", true},
		{"expression7", "x > '2014-01-09' && x < '2017-01-09'", true},
		{"expression8", "x < '2014-01-09' || x > '2017-01-09'", true},
		{"expression9", "x > '2014-01-09' || x > '2016-01-09' && x < '2018-01-09'", true},
		{"invalid", "x < '2014-01-09' && x > '2011-01-09'", true},
		{"invalid1", "x < '2014-01-09' &&", false},
		{"invalid2", "x <", false},
		{"invalid3", "< 100", false},
		{"invalid4", "x", false},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			x := generator.Random(tc.condition)
			if tc.isValidExpression {
				require.NotZero(t, x)
			} else {
				require.Zero(t, x)
			}
			log.Println("[INFO]", tc.name, tc.condition, "result ->", x)
		})
	}
}
