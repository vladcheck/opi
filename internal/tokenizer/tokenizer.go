package tokenizer

import (
	"errors"
	"fmt"
	"strings"
)

var errUnknownCharacter = errors.New("syntax error: unknown character")

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

//nolint:cyclop,funlen // This function must be a giant switch statement
func (t *Tokenizer) match() (Token, error) {
	switch {
	case t.isEOL():
		t.advanceToNextLine()

		return newSentinelToken(), nil
	case IsHorizontalWhitespace(t.peek()):
		t.skipWhitespace()

		return newSentinelToken(), nil
	case IsNumeric(t.peek()):
		t.advance()

		for IsNumeric(t.peek()) {
			t.advance()
		}

		return t.emitWithValue(TtNumber), nil
	case t.isPeekMatches("\""):
		return t.matchString(), nil
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
	case IsAlpha(t.peek()):
		return t.matchSymbolOrKeyword(), nil
	default:
		err := fmt.Errorf("%w %q at %s", errUnknownCharacter, t.peek(), t.position())

		return Token{}, err
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

	return t.emitWithValue(TtString)
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
		if IsHorizontalWhitespace(t.peek()) {
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

	t.advanceToNextLine()
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

	for IsAlphanumeric(t.peek()) {
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
	for IsHorizontalWhitespace(t.peek()) {
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

func (t *Tokenizer) position() string {
	return fmt.Sprintf("line %d, column %d", t.line, t.col)
}
