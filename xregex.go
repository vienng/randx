package randx

import (
	regen "github.com/zach-klippenstein/goregen"
	"log"
	"math/rand"
	"time"
)

// XRegex implements interface X, random a string with inputted regex
type XRegex struct {
	defaultRegex interface{}
}

// NewXRegex makes a new instance for XRegex
func NewXRegex() X {
	return &XRegex{defaultRegex: "^[9]{3}$"}
}

func (xr *XRegex) SetFallback(defaultRegex interface{}) {
	xr.defaultRegex = defaultRegex
}

// BindOperator returns Regex XOP
func (xr XRegex) BindOperator(regex string) XOP {
	return Regex
}

// Random returns a random string generated from inputted regex
func (xr XRegex) Random(regex string) interface{} {
	randAgr := &regen.GeneratorArgs{RngSource: rand.NewSource(time.Now().UnixNano())}
	if len(regex) == 0 {
		regex = xr.defaultRegex.(string)
	}
	generator, err := regen.NewGenerator(regex, randAgr)
	if err != nil {
		log.Println(err)
		return ""
	}
	return generator.Generate()
}
