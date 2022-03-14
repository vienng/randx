package randx

import (
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
			x := generator.Random(tc.regex)
			log.Println("[INFO]", tc.name, "result ->", x)
		})
	}
}
