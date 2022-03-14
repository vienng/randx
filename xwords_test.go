package randx

import (
	"github.com/stretchr/testify/require"
	"log"
	"strings"
	"testing"
)

func TestXWords(t *testing.T) {
	generator := NewXWords("etc/vietnamese")
	generator.SetFallback([]string{"ERROR"})

	testCases := []struct {
		expression        string
		isValidExpression bool
	}{
		{"", true},
		{"length == 4", true},
		{"LENGTH == 4", true},
		{"begin == 'Harry'", true},
		{"begin with 'Harry'", false},
		{"begins == 'Harry'", true},
		{"begin == Harry", false},
		{"END == 'Harry'", true},
		{"end    == Harry", false},
		{"end=='Harry'", true},
		{"length 4", false},
		{"len == 4", false},
		{"length == 4 &&", false},
		{"4", false},
		{"begin Harry", false},
		{"begin > Harry", false},
	}
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.expression, func(t *testing.T) {
			words := generator.Random(tc.expression)
			if tc.isValidExpression {
				require.NotEqual(t, []string{"ERROR"}, words)
			} else {
				require.Equal(t, []string{"ERROR"}, words)
			}
			log.Println("[INFO]", tc.expression, "result ->", strings.Join(words.([]string), " "))
		})
	}
}
