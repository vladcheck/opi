package tokenizer

import (
	"fmt"
	"strings"
)

// TokenType represents a type of a token, say a number, an identifier, or something different.
type TokenType int

const (
	// TtEOF \0.
	TtEOF TokenType = iota
	// TtNewLine is equal to \n.
	TtNewLine
	// TtSymbol represents an identifier.
	TtSymbol
	// TtNumber is a literal number.
	TtNumber
	// TtString is a literal string.
	TtString
	// TtNull represents null.
	TtNull
	// TtDot is the . symbol.
	TtDot
	// TtComma is the , symbol.
	TtComma
	// TtColon is the : symbol.
	TtColon
	// TtLArrow is the < symbol.
	TtLArrow
	// TtLessThanOrEqual is the <= operator.
	TtLessThanOrEqual
	// TtRArrow is the > symbol.
	TtRArrow
	// TtMoreThanOrEqual is the >= operator.
	TtMoreThanOrEqual
	// TtDoubleQuote is the " symbol.
	TtDoubleQuote
	// TtExclamation is the ! symbol.
	TtExclamation
	// TtNotEquals is the != operator.
	TtNotEquals
	// TtPlus is the + operator.
	TtPlus
	// TtPlusEquals is the += operator.
	TtPlusEquals
	// TtPlusPlus is the ++ operator.
	TtPlusPlus
	// TtMinus is the - operator.
	TtMinus
	// TtMinusEquals is the -= operator.
	TtMinusEquals
	// TtMinusMinus is the -- operator.
	TtMinusMinus
	// TtStar is the * operator.
	TtStar
	// TtStarEquals is the *= operator.
	TtStarEquals
	// TtSlash is the / operator.
	TtSlash
	// TtSlashEquals is the /= operator.
	TtSlashEquals
	// TtPercent is the % operator.
	TtPercent
	// TtPercentEquals is the %= operator.
	TtPercentEquals
	// TtEquals is the = operator.
	TtEquals
	// TtEqualsEquals is the == operator.
	TtEqualsEquals
	// TtLParen is the ( symbol.
	TtLParen
	// TtRParen is the ) symbol.
	TtRParen
	// TtLBracket is the [ symbol.
	TtLBracket
	// TtRBracket is the ] symbol.
	TtRBracket
	// TtLBrace is the { symbol.
	TtLBrace
	// TtRBrace is the } symbol.
	TtRBrace
	// TtTrue is the true keyword.
	TtTrue
	// TtFalse is the false keyword.
	TtFalse
	// TtAnd is the and keyword.
	TtAnd
	// TtOr is the or keyword.
	TtOr
	// TtIf is the if keyword.
	TtIf
	// TtElse is the else keyword.
	TtElse
	// TtSwitch is the switch keyword.
	TtSwitch
	// TtFall works like fallthrough in Go.
	TtFall
	// TtBreak is the break keyword.
	TtBreak
	// TtContinue is the continue keyword.
	TtContinue
	// TtFunc is the func keyword.
	TtFunc
	// TtRet is the ret keyword.
	TtRet
	// TtDefer is the defer keyword.
	TtDefer
	// TtMut is the mut keyword, inspired by Rust.
	TtMut
	// TtFor is the for keyword.
	TtFor
	// TtWhile is the while keyword.
	TtWhile
	// TtDo is the do keyword.
	TtDo
)

// Token represents a smallest unit of opi language.
type Token struct {
	TokenType TokenType
	Line, Col int
	Value     *string
}

func newToken(tokenType TokenType, line int, col int) Token {
	return Token{
		TokenType: tokenType,
		Line:      line,
		Col:       col,
		Value:     nil,
	}
}

func newTokenWithValue(tokenType TokenType, line int, col int, value string) Token {
	return Token{
		TokenType: tokenType,
		Line:      line,
		Col:       col,
		Value:     &value,
	}
}

// String represents a Token in plain string format.
func (t Token) String() string {
	if t.Value == nil {
		return fmt.Sprintf("%s: %s", t.position(), tokenTypeToString(t.TokenType))
	}

	if strings.EqualFold(string((*t.Value)[0]), "\"") {
		return fmt.Sprintf("%s: %s (%s)", t.position(), tokenTypeToString(t.TokenType), *t.Value)
	}

	return fmt.Sprintf("%s: %s (%q)", t.position(), tokenTypeToString(t.TokenType), *t.Value)
}

func (t Token) position() string {
	return fmt.Sprintf("[%d:%d]", t.Line, t.Col)
}

func tokenTypeBySubstring(substring string) (TokenType, bool) {
	TtKeywordMap := map[string]TokenType{
		"true":     TtTrue,
		"false":    TtFalse,
		"and":      TtAnd,
		"or":       TtOr,
		"if":       TtIf,
		"else":     TtElse,
		"switch":   TtSwitch,
		"fall":     TtFall,
		"break":    TtBreak,
		"continue": TtContinue,
		"func":     TtFunc,
		"ret":      TtRet,
		"defer":    TtDefer,
		"mut":      TtMut,
		"for":      TtFor,
		"while":    TtWhile,
		"do":       TtDo,
	}
	Tt, ok := TtKeywordMap[substring]

	return Tt, ok
}

//nolint:funlen // May be long
func tokenTypeToString(tokenType TokenType) string {
	TtStringMap := map[TokenType]string{
		TtEOF:             "EOF",
		TtNewLine:         "NewLine",
		TtSymbol:          "Symbol",
		TtNumber:          "Number",
		TtString:          "String",
		TtNull:            "Null",
		TtDot:             "Dot",
		TtComma:           "Comma",
		TtColon:           "Colon",
		TtLArrow:          "LArrow",
		TtLessThanOrEqual: "LessThanOrEqual",
		TtRArrow:          "RArrow",
		TtMoreThanOrEqual: "MoreThanOrEqual",
		TtDoubleQuote:     "DoubleQuote",
		TtExclamation:     "Exclamation",
		TtNotEquals:       "NotEquals",
		TtPlus:            "Plus",
		TtPlusEquals:      "PlusEquals",
		TtPlusPlus:        "PlusPlus",
		TtMinus:           "Minus",
		TtMinusEquals:     "MinusEquals",
		TtMinusMinus:      "MinusMinus",
		TtStar:            "Star",
		TtStarEquals:      "StarEquals",
		TtPercent:         "Percent",
		TtPercentEquals:   "PercentEquals",
		TtSlash:           "Slash",
		TtSlashEquals:     "SlashEquals",
		TtEquals:          "Equals",
		TtEqualsEquals:    "EqualsEquals",
		TtLParen:          "LParen",
		TtRParen:          "RParen",
		TtLBracket:        "LBracket",
		TtRBracket:        "RBracket",
		TtLBrace:          "LBrace",
		TtRBrace:          "RBrace",
		TtTrue:            "True",
		TtFalse:           "False",
		TtAnd:             "And",
		TtOr:              "Or",
		TtIf:              "If",
		TtElse:            "Else",
		TtSwitch:          "Switch",
		TtFall:            "Fall",
		TtBreak:           "Break",
		TtContinue:        "Continue",
		TtFunc:            "Func",
		TtRet:             "Ret",
		TtDefer:           "Defer",
		TtMut:             "Mut",
		TtFor:             "For",
		TtWhile:           "While",
		TtDo:              "Do",
	}

	value, ok := TtStringMap[tokenType]
	if !ok {
		return "UNMATCHED (edit tokenTypeToString())"
	}

	return value
}
