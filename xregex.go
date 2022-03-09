package randx

import regen "github.com/zach-klippenstein/goregen"

type XRegex interface {
	RandomString() string
}

type xRegex struct {
	regex string
	gen   regen.Generator
}

func NewXRegex(regex string) (XRegex, error) {
	generator, err := regen.NewGenerator(regex, nil)
	if err != nil {
		return nil, err
	}
	return &xRegex{regex: regex, gen: generator}, nil
}

func (g *xRegex) RandomString() string {
	return g.gen.Generate()
}
