package tokenizer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	alphaRE          = regexp.MustCompile(`[_a-zA-Z]+`)
	isAlphanumericRE = regexp.MustCompile(`[_0-9a-zA-Z]+`)
	isNumericRE      = regexp.MustCompile("[0-9]+")
)

var (
	errUnknownCharacter   = errors.New("syntax error: unknown character")
	errUnterminatedString = errors.New("syntax error: unterminated string")
)

//nolint:cyclop,funlen // This function must be a giant switch statement
func (t *Tokenizer) match() (Token, error) {
	switch {
	case t.isEOL():
		t.advanceToNextLine()

		return newSentinelToken(), nil
	case isHorizontalWhitespace(t.peek()):
		t.skipWhitespace()

		return newSentinelToken(), nil
	case isNumeric(t.peek()):
		t.advance()

		for isNumeric(t.peek()) {
			t.advance()
		}

		return t.emitWithValue(TtNumber), nil
	case t.isPeekMatches("\""):
		return t.matchString()
	case t.isPeekMatches("."):
		return t.emitSingle(TtDot), nil
	case t.isPeekMatches(","):
		return t.emitSingle(TtComma), nil
	case t.isPeekMatches(":"):
		return t.emitSingle(TtColon), nil
	case t.isPeekMatches("!"):
		return t.matchExclamation(), nil
	case t.isPeekMatches("+"):
		return t.matchPlus(), nil
	case t.isPeekMatches("-"):
		return t.matchMinus(), nil
	case t.isPeekMatches("*"):
		return t.matchStar(), nil
	case t.isPeekMatches("/"):
		return t.matchSlash(), nil
	case t.isPeekMatches("%"):
		return t.matchPercent(), nil
	case t.isPeekMatches("="):
		return t.matchEquals(), nil
	case t.isPeekMatches(">"):
		return t.matchRArrow(), nil
	case t.isPeekMatches("<"):
		return t.matchLArrow(), nil
	case t.isPeekMatches("("):
		return t.emitSingle(TtLParen), nil
	case t.isPeekMatches(")"):
		return t.emitSingle(TtRParen), nil
	case t.isPeekMatches("["):
		return t.emitSingle(TtLBracket), nil
	case t.isPeekMatches("]"):
		return t.emitSingle(TtRBracket), nil
	case t.isPeekMatches("{"):
		return t.emitSingle(TtLBrace), nil
	case t.isPeekMatches("}"):
		return t.emitSingle(TtRBrace), nil
	case isAlpha(t.peek()):
		return t.matchSymbolOrKeyword(), nil
	default:
		err := fmt.Errorf("%w %q at %s", errUnknownCharacter, t.peek(), t.position())

		return Token{}, err
	}
}

func (t *Tokenizer) matchString() (Token, error) {
	t.advance()

	possibleEscape := false

	for !t.isEOL() && !t.isEOF() {
		if t.isPeekMatches("\\") {
			possibleEscape = true

			t.advance()

			continue
		}

		if t.isPeekMatches("\"") {
			if !possibleEscape {
				break // string has ended
			}
		}

		t.advance()

		if possibleEscape {
			possibleEscape = false
		}
	}

	if t.isEOF() || t.isEOL() {
		return newSentinelToken(), fmt.Errorf("%w at %s", errUnterminatedString, t.position())
	}

	t.advance()

	return t.emitWithValue(TtString), nil
}

func (t *Tokenizer) matchExclamation() Token {
	return t.matchWithEquals(TtExclamation, TtNotEquals)
}

func (t *Tokenizer) matchPlus() Token {
	next, ok := t.peekNext()
	if !ok {
		return t.emitSingle(TtPlus)
	}

	switch {
	case strings.EqualFold(next, "+"):
		t.advanceTwice()

		return t.emit(TtPlusPlus)
	case strings.EqualFold(next, "="):
		t.advanceTwice()

		return t.emit(TtPlusEquals)
	default:
		return t.emitSingle(TtPlus)
	}
}

func (t *Tokenizer) matchMinus() Token {
	next, ok := t.peekNext()
	if !ok {
		return t.emitSingle(TtMinus)
	}

	switch {
	case strings.EqualFold(next, "-"):
		t.advanceTwice()

		return t.emit(TtMinusMinus)
	case strings.EqualFold(next, "="):
		t.advanceTwice()

		return t.emit(TtMinusEquals)
	default:
		return t.emitSingle(TtMinus)
	}
}

func (t *Tokenizer) matchEquals() Token {
	return t.matchWithEquals(TtEquals, TtEqualsEquals)
}

func (t *Tokenizer) matchStar() Token {
	return t.matchWithEquals(TtStar, TtStarEquals)
}

func (t *Tokenizer) matchSlash() Token {
	next, ok := t.peekNext()
	if !ok {
		return t.emitSingle(TtSlash)
	}

	if strings.EqualFold(next, "/") {
		t.skipInlineComment()

		return newSentinelToken()
	} else if strings.EqualFold(next, "*") {
		t.skipMultilineComments()

		return newSentinelToken()
	}

	return t.matchWithEquals(TtSlash, TtSlashEquals)
}

func (t *Tokenizer) skipMultilineComments() {
	t.inComment = true
	t.advanceTwice()
	t.commentDepth = 1

	for t.commentDepth > 0 {
		if isHorizontalWhitespace(t.peek()) {
			t.skipWhitespace()

			continue
		}

		if t.isPeekMatches("\n") {
			t.advanceToNextLine()

			continue
		}

		matches, ok := t.isTwoCharsMatch("/", "*")
		if !ok {
			return
		}

		if matches {
			t.commentDepth++
		}

		matches, ok = t.isTwoCharsMatch("*", "/")
		if !ok {
			return
		}

		if matches {
			t.commentDepth--
		}

		t.advance()
	}

	t.advanceTwice()
	t.inComment = false
}

func (t *Tokenizer) skipInlineComment() {
	t.inComment = true
	for !t.isEOL() {
		t.advance()
	}

	t.inComment = false
}

func (t *Tokenizer) matchPercent() Token {
	return t.matchWithEquals(TtPercent, TtPercentEquals)
}

func (t *Tokenizer) matchRArrow() Token {
	return t.matchWithEquals(TtRArrow, TtMoreThanOrEqual)
}

func (t *Tokenizer) matchLArrow() Token {
	return t.matchWithEquals(TtLArrow, TtLessThanOrEqual)
}

func (t *Tokenizer) matchSymbolOrKeyword() Token {
	t.advance()

	for isAlphanumeric(t.peek()) {
		t.advance()
	}

	Tt, match := tokenTypeBySubstring(t.see())
	if !match {
		// It is a symbol
		return t.emitWithValue(TtSymbol)
	}

	// It is a keyword
	return t.emit(Tt)
}

func (t *Tokenizer) matchWithEquals(char TokenType, charEquals TokenType) Token {
	next, ok := t.peekNext()

	if !ok {
		return t.emitSingle(char)
	}

	if strings.EqualFold(next, "=") {
		t.advanceTwice()

		return t.emit(charEquals)
	}

	return t.emitSingle(char)
}

func (t *Tokenizer) isTwoCharsMatch(ch1, ch2 string) (bool, bool) {
	if t.isPeekMatches(ch1) {
		next, ok := t.peekNext()
		if !ok {
			return false, false
		}

		return strings.EqualFold(next, ch2), true
	}

	return false, true
}

func (t *Tokenizer) skipWhitespace() {
	for isHorizontalWhitespace(t.peek()) {
		t.advance()
	}

	t.consume()
}

func (t *Tokenizer) see() string {
	return t.input[t.start:t.current]
}

func (t *Tokenizer) isEOL() bool {
	return !t.isEOF() && t.isPeekMatches("\n")
}

func (t *Tokenizer) isEOF() bool {
	return t.current >= len(t.input)
}

func (t *Tokenizer) isPeekMatches(char string) bool {
	return strings.EqualFold(t.peek(), char)
}

// isHorizontalWhitespace checks if substring contains only whitespace characters except '\n'.
func isHorizontalWhitespace(s string) bool {
	if s == "" {
		return false
	}

	for _, r := range s {
		if r != ' ' && r != '\t' {
			return false
		}
	}

	return true
}

// isAlpha checks if substring is alphabetic.
func isAlpha(substring string) bool {
	return alphaRE.MatchString(substring)
}

// isAlphanumeric checks if substring is alphanumeric.
func isAlphanumeric(substring string) bool {
	return isAlphanumericRE.MatchString(substring)
}

// isNumeric checks if substring contains only digits.
func isNumeric(substring string) bool {
	return isNumericRE.MatchString(substring)
}
