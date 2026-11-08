package tokenizer

import (
	"log"
	"strings"
)

// Tokenizer is responsible for taking source code and turning it into a slice of tokens.
type Tokenizer struct {
	input                              string
	start, current, line, commentDepth int
	inString                           bool
}

// NewTokenizer creates a new Tokenizer.
func NewTokenizer(input string) Tokenizer {
	return Tokenizer{
		input:        input,
		start:        0,
		current:      0,
		line:         1,
		commentDepth: 0,
		inString:     false,
	}
}

// Tokenize takes source code and turns it into a slice of tokens.
func (t *Tokenizer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0, 2)

	for !t.isEOF() {
		token, ok := t.match()

		if ok {
			tokens = append(tokens, token)
		}
	}

	tokens = append(tokens, newToken(ttEOF, t.line, "\x00"))

	return tokens, nil
}

//nolint:cyclop,funlen // This function must be a giant switch statement
func (t *Tokenizer) match() (Token, bool) {
	switch {
	case t.isEOL():
		t.advanceToNextLine()

		return Token{}, false
	case IsWhitespace(t.peek()):
		t.skipWhitespace()

		return Token{}, false
	case IsNumeric(t.peek()):
		t.advance()

		for IsNumeric(t.peek()) {
			t.advance()
		}

		return t.emit(ttNumber), true
	case t.isPeekMatches("\""):
		return t.matchString(), true
	case t.isPeekMatches("."):
		return t.emitSingle(ttDot), true
	case t.isPeekMatches(","):
		return t.emitSingle(ttComma), true
	case t.isPeekMatches(":"):
		return t.emitSingle(ttColon), true
	case t.isPeekMatches("!"):
		return t.matchExclamation(), true
	case t.isPeekMatches("+"):
		return t.matchPlus(), true
	case t.isPeekMatches("-"):
		return t.matchMinus(), true
	case t.isPeekMatches("*"):
		return t.matchStar(), true
	case t.isPeekMatches("/"):
		return t.matchSlash()
	case t.isPeekMatches("%"):
		return t.matchPercent(), true
	case t.isPeekMatches("="):
		return t.matchEquals(), true
	case t.isPeekMatches(">"):
		return t.matchRArrow(), true
	case t.isPeekMatches("<"):
		return t.matchLArrow(), true
	case t.isPeekMatches("("):
		return t.emitSingle(ttLParen), true
	case t.isPeekMatches(")"):
		return t.emitSingle(ttRParen), true
	case t.isPeekMatches("["):
		return t.emitSingle(ttLBracket), true
	case t.isPeekMatches("]"):
		return t.emitSingle(ttRBracket), true
	case t.isPeekMatches("{"):
		return t.emitSingle(ttLBrace), true
	case t.isPeekMatches("}"):
		return t.emitSingle(ttRBrace), true
	case IsAlpha(t.peek()):
		return t.matchSymbolOrKeyword()
	default:
		log.Fatalf("Can't lex: %q", t.peek())

		return Token{}, false
	}
}

func (t *Tokenizer) matchString() Token {
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

	t.advance()

	return t.emit(ttString)
}

func (t *Tokenizer) matchExclamation() Token {
	return t.matchWithEquals(ttExclamation, ttNotEquals)
}

func (t *Tokenizer) matchPlus() Token {
	next, ok := t.peekNext()
	if !ok {
		return t.emitSingle(ttPlus)
	}

	switch {
	case strings.EqualFold(next, "+"):
		t.advanceTwice()

		return t.emit(ttPlusPlus)
	case strings.EqualFold(next, "="):
		t.advanceTwice()

		return t.emit(ttPlusEquals)
	default:
		return t.emitSingle(ttPlus)
	}
}

func (t *Tokenizer) matchMinus() Token {
	next, ok := t.peekNext()
	if !ok {
		return t.emitSingle(ttMinus)
	}

	switch {
	case strings.EqualFold(next, "-"):
		t.advanceTwice()

		return t.emit(ttMinusMinus)
	case strings.EqualFold(next, "="):
		t.advanceTwice()

		return t.emit(ttMinusEquals)
	default:
		return t.emitSingle(ttMinus)
	}
}

func (t *Tokenizer) matchEquals() Token {
	return t.matchWithEquals(ttEquals, ttEqualsEquals)
}

func (t *Tokenizer) matchStar() Token {
	return t.matchWithEquals(ttStar, ttStarEquals)
}

func (t *Tokenizer) matchSlash() (Token, bool) {
	next, ok := t.peekNext()
	if !ok {
		return t.emitSingle(ttSlash), true
	}

	if strings.EqualFold(next, "/") {
		t.skipInlineComment()

		return Token{}, false
	} else if strings.EqualFold(next, "*") {
		t.advanceTwice()
		t.commentDepth = 1
		t.skipMultilineComments()

		return Token{}, false
	}

	return t.matchWithEquals(ttSlash, ttSlashEquals), true
}

func (t *Tokenizer) skipMultilineComments() {
	for t.commentDepth > 0 {
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
}

func (t *Tokenizer) skipInlineComment() {
	for !t.isEOL() {
		t.advance()
	}

	t.advanceToNextLine()
}

func (t *Tokenizer) matchPercent() Token {
	return t.matchWithEquals(ttPercent, ttPercentEquals)
}

func (t *Tokenizer) matchRArrow() Token {
	return t.matchWithEquals(ttRArrow, ttMoreThanOrEqual)
}

func (t *Tokenizer) matchLArrow() Token {
	return t.matchWithEquals(ttLArrow, ttLessThanOrEqual)
}

func (t *Tokenizer) matchSymbolOrKeyword() (Token, bool) {
	t.advance()

	for IsAlphanumeric(t.peek()) {
		t.advance()
	}

	tt, match := tokenTypeBySubstring(t.see())
	if !match {
		return t.emit(ttSymbol), true
	}

	return t.emit(tt), true
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

// emit creates a token in range `[t.start:t.current]`, consumes and returns it.
func (t *Tokenizer) emit(tt TokenType) Token {
	token := newToken(tt, t.line, t.input[t.start:t.current])
	t.consume()

	return token
}

func (t *Tokenizer) emitSingle(tt TokenType) Token {
	t.advance()

	return t.emit(tt)
}

func (t *Tokenizer) peek() string {
	if t.isEOF() {
		return ""
	}

	return string(t.input[t.current])
}

func (t *Tokenizer) peekNext() (string, bool) {
	if t.current+1 >= len(t.input) {
		return "", false
	}

	return string(t.input[t.current+1]), true
}

func (t *Tokenizer) see() string {
	return t.input[t.start:t.current]
}

func (t *Tokenizer) skipWhitespace() {
	for IsWhitespace(t.peek()) {
		t.advance()
	}

	t.consume()
}

func (t *Tokenizer) isEOF() bool {
	return t.current >= len(t.input)
}

func (t *Tokenizer) isEOL() bool {
	return !t.isEOF() && t.isPeekMatches("\n")
}

func (t *Tokenizer) isPeekMatches(char string) bool {
	return strings.EqualFold(t.peek(), char)
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

func (t *Tokenizer) advanceToNextLine() {
	t.line++
	t.advance()
	t.consume()
}

func (t *Tokenizer) advanceTwice() {
	t.advance()
	t.advance()
}

func (t *Tokenizer) advance() {
	t.current++
}

func (t *Tokenizer) consume() {
	t.start = t.current
}
