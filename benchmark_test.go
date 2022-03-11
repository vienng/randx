package randx

import (
	"testing"
	"time"
)

func BenchmarkNewXNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewXNumber(0, 1000, 1)
	}
}

func BenchmarkNewXWord(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewXWord("")
	}
}

func BenchmarkNewXTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewXTime(time.Now(), time.Now(), time.Second)
	}
}

func BenchmarkNewXRegex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewXRegex()
	}
}

func BenchmarkXNumber_BindOperator(b *testing.B) {
	xn := NewXNumber(0, 100, 1)
	for i := 0; i < b.N; i++ {
		xn.BindOperator("x > 0 && x < 100")
	}
}

func BenchmarkXWord_BindOperator(b *testing.B) {
	xw := NewXWord("")
	for i := 0; i < b.N; i++ {
		xw.BindOperator("begin == 'John'")
	}
}

func BenchmarkXTime_BindOperator(b *testing.B) {
	xt := NewXTime(time.Now().AddDate(-100, 0, 0), time.Now(), time.Second)
	for i := 0; i < b.N; i++ {
		xt.BindOperator("begin == 'John'")
	}
}

func BenchmarkXRegex_BindOperator(b *testing.B) {
	xr := NewXRegex()
	for i := 0; i < b.N; i++ {
		xr.BindOperator("xRegex doesn't need to bind operator")
	}
}

func BenchmarkXNumber_Random(b *testing.B) {
	xn := NewXNumber(0, 1000000, 1)
	for i := 0; i < b.N; i++ {
		_, err := xn.Random("x > 0 && x < 100")
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkXWord_Random(b *testing.B) {
	xw := NewXWord("")
	for i := 0; i < b.N; i++ {
		_, err := xw.Random("length == 3")
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkXTime_Random(b *testing.B) {
	xt := NewXTime(time.Now().AddDate(-100, 0, 0), time.Now(), 365*24*time.Hour)
	for i := 0; i < b.N; i++ {
		_, err := xt.Random("day > '1960-01-01' || day < '1990-01-01'")
		if err != nil {
			b.FailNow()
		}
	}
}

func BenchmarkXRegex_Random(b *testing.B) {
	xr := NewXRegex()
	for i := 0; i < b.N; i++ {
		_, err := xr.Random("^[a-z]{6}{[0-9]{3}$")
		if err != nil {
			b.FailNow()
		}
	}
}
