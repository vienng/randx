package randx

/*
	Represents all valid operator of X
*/
const (
	Unknown XOP = iota
	Any
	Invalid

	Length
	Begin
	End

	Constant
	FindX
	FindXs

	Regex
)

// XOP defines operator of X
type XOP int

/*
	String method returns a string that describes the given XOP
	e.g., when passed the FindX XOP, this returns the string "FindX".
*/
func (op XOP) String() string {
	switch op {
	case Any:
		return "Any"
	case Invalid:
		return "Invalid"
	case Length:
		return "Length"
	case Begin:
		return "Begin"
	case End:
		return "End"
	case Constant:
		return "Constant"
	case FindX:
		return "FindX"
	case FindXs:
		return "FindXs"
	case Regex:
		return "Regex"
	default:
		return "Unknown"
	}
}
