package tokenizer

import (
	"fmt"
)

// Tokenizer is responsible for taking source code and turning it into a slice of tokens.
type Tokenizer struct {
	input          string
	start, current int
	line, col      int
	commentDepth   int
	inComment      bool
	tokens         []Token
}

// NewTokenizer creates a new Tokenizer.
func NewTokenizer(input string) *Tokenizer {
	return new(Tokenizer{
		input:        input,
		start:        0,
		current:      0,
		line:         1,
		col:          1,
		commentDepth: 0,
		inComment:    false,
		tokens:       make([]Token, 0, 2),
	})
}

// Tokenize takes source code and turns it into a slice of tokens.
func (t *Tokenizer) Tokenize() ([]Token, error) {
	for !t.isEOF() {
		token, err := t.match()
		if err != nil {
			return t.tokens, err
		}

		if token.TokenType != TtSentinel {
			t.tokens = append(t.tokens, token)
		}
	}

	t.tokens = append(t.tokens, newToken(TtEOF, t.line, t.col))

	return t.tokens, nil
}

func (t *Tokenizer) skipWhitespace() {
	for isHorizontalWhitespace(t.peek()) {
		t.advance()
	}

	t.consume()
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

func (t *Tokenizer) advanceToNextLine() {
	if !t.inComment {
		t.tokens = append(t.tokens, newToken(TtNewLine, t.line, t.col))
	}

	t.resetColumn()
	t.line++
	t.consume()
	t.advance()
}

func (t *Tokenizer) advanceTwice() {
	t.advance()
	t.advance()
}

func (t *Tokenizer) advance() {
	t.col++
	t.current++
}

func (t *Tokenizer) consume() {
	t.start = t.current
}

func (t *Tokenizer) resetColumn() {
	t.col = 0
}

func (t *Tokenizer) see() string {
	return t.input[t.start:t.current]
}

func (t *Tokenizer) position() string {
	return fmt.Sprintf("line %d, column %d", t.line, t.col)
}
