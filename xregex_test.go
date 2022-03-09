package randx

import (
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestXRegex(t *testing.T) {
	testCases := []struct {
		name    string
		regex   string
		isError bool
	}{
		{"error", "[}", true},
		{"empty", "", false},
		{"any", "^.{3}$", false},
		{"valid", "^[0-9]{3}[a-z]{3}$", false},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			generator := NewXRegex()
			x, err := generator.Random(tc.regex)
			if tc.isError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Regexp(t, tc.regex, x)
			}
			log.Println("[INFO]", tc.name, "result ->", x)
		})
	}
}
