package internal

const (
	ttEOF = iota
)

// Token represents a smallest unit of opi language.
type Token struct {
	TokenType, Line int
	Value           any
}

func newToken(tokenType, line int, value any) Token {
	return Token{
		TokenType: tokenType,
		Line:      line,
		Value:     value,
	}
}
