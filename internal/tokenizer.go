package internal

import (
	"strings"
)

// Tokenizer is responsible for taking source code and turning it into a slice of tokens.
type Tokenizer struct {
	input                string
	start, current, line int
}

// NewTokenizer creates a new Tokenizer.
func NewTokenizer(input string) Tokenizer {
	return Tokenizer{
		input:   input,
		start:   0,
		current: 0,
		line:    1,
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

//nolint:cyclop // This function must be a giant switch statement
func (t *Tokenizer) match() (Token, bool) {
	switch {
	case t.isEOL():
		t.line++
		t.advance()
		t.consume()

		return Token{}, false
	case isWhitespace(t.peek()):
		t.skipWhitespace()

		return Token{}, false
	case isDigit(t.peek()):
		t.advance()

		for isDigit(t.peek()) {
			t.advance()
		}

		return t.emit(ttNumber), true
	case strings.EqualFold(t.peek(), "+"):
		return t.matchPlus(), true
	case strings.EqualFold(t.peek(), "-"):
		return t.matchMinus(), true
	case strings.EqualFold(t.peek(), "*"):
		return t.matchStar(), true
	case strings.EqualFold(t.peek(), "/"):
		return t.matchSlash(), true
	case strings.EqualFold(t.peek(), "%"):
		return t.matchPercent(), true
	case strings.EqualFold(t.peek(), "="):
		return t.matchEquals(), true
	case strings.EqualFold(t.peek(), "("):
		return t.emitSingle(ttLeftParen), true
	case strings.EqualFold(t.peek(), ")"):
		return t.emitSingle(ttRightParen), true
	case strings.EqualFold(t.peek(), "["):
		return t.emitSingle(ttLeftBracket), true
	case strings.EqualFold(t.peek(), "]"):
		return t.emitSingle(ttRightBracket), true
	case strings.EqualFold(t.peek(), "{"):
		return t.emitSingle(ttLeftBrace), true
	case strings.EqualFold(t.peek(), "}"):
		return t.emitSingle(ttRightBrace), true
	default:
		return Token{}, false
	}
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

func (t *Tokenizer) matchSlash() Token {
	return t.matchWithEquals(ttSlash, ttSlashEquals)
}

func (t *Tokenizer) matchPercent() Token {
	return t.matchWithEquals(ttPercent, ttPercentEquals)
}

func (t *Tokenizer) matchWithEquals(
	char TokenType,
	charEquals TokenType,
) Token {
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

func (t *Tokenizer) skipWhitespace() {
	for isWhitespace(t.peek()) {
		t.advance()
	}

	t.consume()
}

func (t *Tokenizer) isEOF() bool {
	return t.current >= len(t.input)
}

func (t *Tokenizer) isEOL() bool {
	return !t.isEOF() && strings.EqualFold(t.peek(), "\n")
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
