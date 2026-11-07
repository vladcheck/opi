package internal

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
		line:    0,
	}
}

// Tokenize takes source code and turns it into a slice of tokens.
func (t *Tokenizer) Tokenize() ([]Token, error) {
	tokens := make([]Token, 0, 2) //nolint:mnd

	for !t.isEOF() {
		t.current++
	}

	tokens = append(tokens, newToken(ttEOF, t.line, "\\0"))

	return tokens, nil
}

func (t *Tokenizer) isEOF() bool {
	return t.current >= len(t.input)
}
