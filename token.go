package ptime

type TokenType int

const (
	End TokenType = iota
	Number
	Text
	Indicator
)

func (t TokenType) String() string {
	switch t {
	case Number:
		return "number"
	case Text:
		return "text"
	case Indicator:
		return "indicator"
	}
	return "invalid"
}

type Token struct {
	Type TokenType
	Val  string
	Pos  int
}
