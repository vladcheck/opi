package tokenizer

import "fmt"

// TokenType represents a type of a token, say a number, an identifier, or something different.
type TokenType int

const (
	ttEOF           TokenType = iota // \0
	ttSymbol                         //
	ttNumber                         // literal number
	ttString                         // literal string
	ttNull                           // null
	ttExclamation                    // !
	ttNotEquals                      // !=
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
	ttLParen                         // (
	ttRParen                         // )
	ttLBracket                       // [
	ttRBracket                       // ]
	ttLBrace                         // {
	ttRBrace                         // }
	ttTrue                           // true
	ttFalse                          // false
	ttAnd                            // and
	ttOr                             // or
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

func tokenTypeBySubstring(substring string) (TokenType, bool) {
	ttKeywordMap := map[string]TokenType{
		"true":  ttTrue,
		"false": ttFalse,
		"and":   ttAnd,
		"or":    ttOr,
	}
	tt, ok := ttKeywordMap[substring]

	return tt, ok
}

func tokenTypeToString(tokenType TokenType) string {
	ttStringMap := map[TokenType]string{
		ttEOF:           "EOF",
		ttSymbol:        "Symbol",
		ttNumber:        "Number",
		ttString:        "String",
		ttNull:          "Null",
		ttExclamation:   "Exclamation",
		ttNotEquals:     "NotEquals",
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
		ttLParen:        "LParen",
		ttRParen:        "RParen",
		ttLBracket:      "LBracket",
		ttRBracket:      "RBracket",
		ttLBrace:        "LBrace",
		ttRBrace:        "RBrace",
		ttTrue:          "True",
		ttFalse:         "False",
		ttAnd:           "And",
		ttOr:            "Or",
	}

	value, ok := ttStringMap[tokenType]
	if !ok {
		return "UNMATCHED (edit tokenTypeToString())"
	}

	return value
}
