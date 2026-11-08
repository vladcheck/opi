package internal

import "fmt"

// TokenType represents a type of a token, say a number, an identifier, or something different.
type TokenType int

const (
	ttEOF        TokenType = iota // \0
	ttNumber                      // literal number
	ttString                      // literal string
	ttNull                        // null
	ttPlus                        // +
	ttMinus                       // -
	ttStar                        // *
	ttPercent                     // %
	ttLeftParen                   // (
	ttRightParen                  // )
)

// Token represents a smallest unit of opi language.
type Token struct {
	TokenType TokenType
	Line      int
	Value     any
}

func newToken(tokenType TokenType, line int, value any) Token {
	return Token{
		TokenType: tokenType,
		Line:      line,
		Value:     value,
	}
}

// String represents a Token in plain string format.
func (t Token) String() string {
	return fmt.Sprintf("[line %d]: %d (%q)", t.Line, int(t.TokenType), t.Value)
}
