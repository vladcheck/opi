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
			Root: NewNode(NtProgram),
		},
		current: 0,
	}
}

// Parse consumes the parser's tokens and returns the resulting AST.
func (p *Parser) Parse() (*Node, error) {
	node, err := p.parseExpression()
	if err != nil {
		return p.Tree.Root, err
	}

	p.Tree.Root.Left = node

	return p.Tree.Root, nil
}

func (p *Parser) parseExpression() (*Node, error) {
	node := NewNode(NtExpression)

	current := node

	for p.peekType() != tz.TtNewLine {
		newNode, err := p.parseOr()
		if err != nil {
			return nil, err
		}

		current.Left = newNode
		current = current.Left
	}

	return node, nil
}

// level 8.
func (p *Parser) parseOr() (*Node, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.peekType() == tz.TtOr {
		p.advance() // consume TtOr

		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}

		n := NewNode(NtOr)
		n.Left = left
		n.Right = right
		left = n
	}

	return left, nil
}

// level 7.
func (p *Parser) parseAnd() (*Node, error) {
	left, err := p.parseEquality()
	if err != nil {
		return nil, err
	}

	for p.peekType() == tz.TtAnd {
		p.advance()

		right, err := p.parseEquality()
		if err != nil {
			return nil, err
		}

		n := NewNode(NtAnd)
		n.Left = left
		n.Right = right
		left = n
	}

	return left, nil
}

// level 6.
func (p *Parser) parseEquality() (*Node, error) {
	left, err := p.parseLiteral()
	if err != nil {
		return nil, err
	}

	for p.peekType() == tz.TtEqualsEquals || p.peekType() == tz.TtNotEquals {
		tt := p.peekType()
		p.advance()

		right, err := p.parseLiteral()
		if err != nil {
			return nil, err
		}

		var node *Node

		switch tt {
		case tz.TtEqualsEquals:
			node = NewNode(NtEquals)
		case tz.TtNotEquals:
			node = NewNode(NtNotEquals)
		}

		node.Left = left
		node.Right = right
		left = node
	}

	return left, nil
}

func (p *Parser) parseLiteral() (*Node, error) {
	switch p.peekType() {
	case tz.TtNumber:
		return p.parseNumber()
	case tz.TtString:
		return p.parseString()
	case tz.TtFalse, tz.TtTrue:
		return p.parseBool()
	default:
		return nil, p.announceError(errNotALiteral)
	}
}

func (p *Parser) parseNumber() (*Node, error) {
	tt := p.peekType()

	return p.parseTerminalWithValue(tt == tz.TtNumber, NtNumber, p.announceError(errNotANumber))
}

func (p *Parser) parseString() (*Node, error) {
	tt := p.peekType()

	return p.parseTerminalWithValue(tt == tz.TtString, NtString, p.announceError(errNotAString))
}

func (p *Parser) parseBool() (*Node, error) {
	tt := p.peekType()

	return p.parseTerminalWithValue(tt == tz.TtTrue || tt == tz.TtFalse, NtBool, p.announceError(errNotABool))
}

func (p *Parser) parseTerminalWithValue(condition bool, nt NodeType, err error) (*Node, error) {
	if condition {
		v := p.peek().Value
		p.advance()

		return NewNodeWithValue(nt, v), nil
	}

	return nil, err
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
