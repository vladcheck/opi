package tokenizer

import "fmt"

// TokenType represents a type of a token, say a number, an identifier, or something different.
type TokenType int

const (
	ttEOF           TokenType = iota // \0
	ttNumber                         // literal number
	ttString                         // literal string
	ttNull                           // null
	ttPlus                           // +
	ttPlusEquals                     // +=
	ttPlusPlus                       // ++
	ttMinus                          // -
	ttMinusEquals                    // -=
	ttMinusMinus                     // --
	ttStar                           // *
	ttStarEquals                     // *=
	ttSlash                          // /
	ttSlashEquals                    // /=
	ttPercent                        // %
	ttPercentEquals                  // %=
	ttEquals                         // =
	ttEqualsEquals                   // ==
	ttLeftParen                      // (
	ttRightParen                     // )
	ttLeftBracket                    // [
	ttRightBracket                   // ]
	ttLeftBrace                      // {
	ttRightBrace                     // }
)

// Token represents a smallest unit of opi language.
type Token struct {
	TokenType TokenType
	Line      int
	Value     string
}

func newToken(tokenType TokenType, line int, value string) Token {
	return Token{
		TokenType: tokenType,
		Line:      line,
		Value:     value,
	}
}

// String represents a Token in plain string format.
func (t Token) String() string {
	return fmt.Sprintf("[line %d]: %s (%q)", t.Line, tokenTypeToString(t.TokenType), t.Value)
}

func tokenTypeToString(tokenType TokenType) string {
	ttStringMap := map[TokenType]string{
		ttEOF:           "EOF",
		ttNumber:        "Number",
		ttString:        "String",
		ttNull:          "Null",
		ttPlus:          "Plus",
		ttPlusEquals:    "PlusEquals",
		ttPlusPlus:      "PlusPlus",
		ttMinus:         "Minus",
		ttMinusEquals:   "MinusEquals",
		ttMinusMinus:    "MinusMinus",
		ttStar:          "Star",
		ttStarEquals:    "StarEquals",
		ttPercent:       "Percent",
		ttPercentEquals: "PercentEquals",
		ttSlash:         "Slash",
		ttSlashEquals:   "SlashEquals",
		ttEquals:        "Equals",
		ttEqualsEquals:  "EqualsEquals",
		ttLeftParen:     "LeftParen",
		ttRightParen:    "RightParen",
		ttLeftBracket:   "LeftBracket",
		ttRightBracket:  "RightBracket",
		ttLeftBrace:     "LeftBrace",
		ttRightBrace:    "RightBrace",
	}

	value, ok := ttStringMap[tokenType]
	if !ok {
		return "UNMATCHED"
	}

	return value
}
