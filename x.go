package randx

// X is an interface providing the random method with user condition
type X interface {
	Random(expression string) (interface{}, error)
	BindOperator(expression string) XOP
}
