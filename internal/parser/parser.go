package parser

import (
	"errors"
	"fmt"
	tz "vladcheck/opi/internal/tokenizer"
)

var (
	errNotANumber  = errors.New("not a number")
	errNotAString  = errors.New("not a string")
	errNotABool    = errors.New("not a bool")
	errNotALiteral = errors.New("not a literal")
)

// Parser consumes a slice of tokens and produces an AST.
type Parser struct {
	tokens []tz.Token
	// Tree is the AST produced by the most recent call to Parse.
	Tree    AST
	current int
}

// NewParser returns a Parser initialized to tokenize the given tokens.
func NewParser(tokens []tz.Token) Parser {
	return Parser{
		tokens: tokens,
		Tree: AST{
			Root: nil,
		},
		current: 0,
	}
}

// Parse consumes the parser's tokens and returns the resulting AST.
func (p *Parser) Parse() (*Node, error) {
	root := NewNode(NtProgram)

	node, err := p.parseExpression()
	if err != nil {
		return p.Tree.Root, err
	}

	root.Left = node
	p.Tree.Root = root

	return p.Tree.Root, nil
}

func (p *Parser) peek() tz.Token {
	return p.tokens[p.current]
}

func (p *Parser) peekType() tz.TokenType {
	return p.peek().TokenType
}

func (p *Parser) advance() {
	p.current++
}

func (p *Parser) announceError(err error) error {
	token := p.peek()

	return fmt.Errorf("parse error: %w, got %s at %s", err, token, token.Position())
}
