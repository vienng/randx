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
		NewXWords("etc/vietnamese")
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
	xw := NewXWords("etc/vietnamese")
	for i := 0; i < b.N; i++ {
		xw.BindOperator("begin == 'John'")
	}
}

func BenchmarkXTime_BindOperator(b *testing.B) {
	xt := NewXTime(time.Now().AddDate(-100, 0, 0), time.Now(), time.Second)
	for i := 0; i < b.N; i++ {
		xt.BindOperator("created_at > '01-01-2021'")
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
		xn.Random("x > 0 && x < 100")
	}
}

func BenchmarkXWord_Random(b *testing.B) {
	xw := NewXWords("etc/vietnamese")
	for i := 0; i < b.N; i++ {
		xw.Random("length == 3")
	}
}

func BenchmarkXTime_Random(b *testing.B) {
	xt := NewXTime(time.Now().AddDate(-100, 0, 0), time.Now(), 365*24*time.Hour)
	for i := 0; i < b.N; i++ {
		xt.Random("day > '1960-01-01' || day < '1990-01-01'")
	}
}

func BenchmarkXRegex_Random(b *testing.B) {
	xr := NewXRegex()
	for i := 0; i < b.N; i++ {
		xr.Random("^[a-z]{6}{[0-9]{3}$")
	}
}
