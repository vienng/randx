package randx

// X is an interface providing the random method with user condition
type X interface {
	Random(expression string) interface{}
	BindOperator(expression string) XOP
	SetFallback(value interface{})
}
