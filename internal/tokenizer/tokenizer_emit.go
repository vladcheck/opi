package tokenizer

// emit creates a token in range `[t.start:t.current]`, consumes and returns it.
func (t *Tokenizer) emit(tt TokenType) Token {
	token := newToken(tt, t.line, t.col)
	t.consume()

	return token
}

func (t *Tokenizer) emitWithValue(tt TokenType) Token {
	token := newTokenWithValue(tt, t.line, t.col, t.input[t.start:t.current])
	t.consume()

	return token
}

func (t *Tokenizer) emitSingle(tt TokenType) Token {
	t.advance()

	return t.emit(tt)
}
