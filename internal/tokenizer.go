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
	tokens := make([]Token, 0, 2) //nolint:mnd

	for !t.isEOF() {
		token, ok := t.match()

		if ok {
			tokens = append(tokens, token)
		}
	}

	tokens = append(tokens, newToken(ttEOF, t.line, "\x00"))

	return tokens, nil
}

func (t *Tokenizer) match() (Token, bool) {
	// fmt.Printf("start: '%v' current: '%v' peek: '%s'\n", t.start, t.current, t.peek())
	switch {
	case strings.EqualFold(t.peek(), "\n"):
		t.line++
		t.advance()
		t.consume()

		return Token{}, false
	case isWhitespace(t.peek()):
		t.skipWhitespace()

		return Token{}, false
	case isDigit(t.peek()):
		for isDigit(t.peek()) {
			t.advance()
		}

		token := newToken(ttNumber, t.line, t.input[t.start:t.current])
		t.consume()

		return token, true
	case strings.EqualFold(t.peek(), "+"):
		token := newToken(ttPlus, t.line, "+")
		t.advance()
		t.consume()

		return token, true
	default:
		return Token{}, false
	}
}

func (t *Tokenizer) peek() string {
	if t.isEOF() {
		return ""
	}

	return string(t.input[t.current])
}

// func (t *Tokenizer) peekNext() string {
// 	if t.isEOF() || t.current+1 >= len(t.input) {
// 		return ""
// 	}

// 	return string(t.input[t.current+1])
// }

func (t *Tokenizer) skipWhitespace() {
	for isWhitespace(t.peek()) {
		t.advance()
	}

	t.consume()
}

func (t *Tokenizer) isEOF() bool {
	return t.current >= len(t.input)
}

func (t *Tokenizer) advance() {
	t.current++
}

func (t *Tokenizer) consume() {
	t.start = t.current
}
