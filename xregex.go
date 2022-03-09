package randx

import (
	regen "github.com/zach-klippenstein/goregen"
	"math/rand"
	"time"
)

// XRegex implements interface X, random a string with inputted regex
type XRegex struct{}

// NewXRegex makes a new instance for XRegex
func NewXRegex() X {
	return XRegex{}
}

// BindOperator returns Regex XOP
func (xr XRegex) BindOperator(regex string) XOP {
	return Regex
}

// Random returns a random string generated from inputted regex
func (xr XRegex) Random(regex string) (interface{}, error) {
	randAgr := &regen.GeneratorArgs{RngSource: rand.NewSource(time.Now().UnixNano())}
	generator, err := regen.NewGenerator(regex, randAgr)
	if err != nil {
		return nil, err
	}
	return generator.Generate(), nil
}
