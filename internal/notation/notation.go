package notation

type NotationStyle int

const (
	DotNotation = iota
	BracketsNotation
)

func GetNotationStyle(useBrackets bool) NotationStyle {
	if useBrackets {
		return BracketsNotation
	}

	return DotNotation
}
