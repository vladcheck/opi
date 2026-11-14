package parser

import (
	tz "vladcheck/opi/internal/tokenizer"
)

// parseExpressionAsGroup is a decorator over parseExpression
// that replaces NodeType NtExpr with NtGroup before returning.
//
// Helpful for debugging.
func (p *Parser) parseExpressionAsGroup() (*Node, error) {
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	// Signifying that this expression is a group
	node.NodeType = NtGroup

	return node, nil
}

func (p *Parser) parseExpression() (*Node, error) {
	node := NewNode(NtExpression)

	newNode, err := p.parseOr()
	if err != nil {
		return nil, err
	}

	node.Left = newNode
	node = node.Left

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
	left, err := p.parseRelation()
	if err != nil {
		return nil, err
	}

	for p.peekType() == tz.TtEqualsEquals || p.peekType() == tz.TtNotEquals {
		tt := p.peekType()
		p.advance()

		right, err := p.parseRelation()
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

// level 5.
func (p *Parser) parseRelation() (*Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}

	tokenType := p.peekType()
	if tz.IsRelationTokenType(tokenType) {
		p.advance()

		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}

		var node *Node

		switch tokenType {
		case tz.TtLArrow:
			node = NewNode(NtLessThan)
		case tz.TtNotEquals:
			node = NewNode(NtMoreThan)
		case tz.TtLessThanOrEqual:
			node = NewNode(NtLessThanOrEqual)
		case tz.TtMoreThanOrEqual:
			node = NewNode(NtMoreThanOrEqual)
		}

		node.Left = left
		node.Right = right
		left = node
	}

	return left, nil
}

// level 4.
func (p *Parser) parseTerm() (*Node, error) {
	left, err := p.parseFactor()
	if err != nil {
		return nil, err
	}

	for p.peekType() == tz.TtPlus || p.peekType() == tz.TtMinus {
		tt := p.peekType()
		p.advance()

		right, err := p.parseFactor()
		if err != nil {
			return nil, err
		}

		var node *Node

		switch tt {
		case tz.TtPlus:
			node = NewNode(NtAdd)
		case tz.TtMinus:
			node = NewNode(NtSub)
		}

		node.Left = left
		node.Right = right
		left = node
	}

	return left, nil
}

// level 3.
func (p *Parser) parseFactor() (*Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}

	for p.peekType() == tz.TtStar || p.peekType() == tz.TtSlash || p.peekType() == tz.TtPercent {
		tt := p.peekType()
		p.advance()

		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}

		var node *Node

		switch tt {
		case tz.TtStar:
			node = NewNode(NtMul)
		case tz.TtSlash:
			node = NewNode(NtDiv)
		case tz.TtPercent:
			node = NewNode(NtMod)
		}

		node.Left = left
		node.Right = right
		left = node
	}

	return left, nil
}

// level 2.
func (p *Parser) parseUnary() (*Node, error) {
	left, err := p.parsePostfix()
	if err != nil {
		return nil, err
	}

	return left, nil
}

// level 1.
func (p *Parser) parsePostfix() (*Node, error) {
	left, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}

	return left, nil
}

// level 0.
func (p *Parser) parsePrimary() (*Node, error) {
	var (
		left *Node
		err  error
	)

	if left, err = p.parseLiteral(); err == nil {
		return left, nil
	}

	if left, err = p.parseSymbol(); err == nil {
		return left, nil
	}

	if left, err = p.parseGroup(); err == nil {
		return left, nil
	}

	return nil, p.announceError(errExpectedSymbolLiteralOrExpr)
}

func (p *Parser) parseGroup() (*Node, error) {
	if p.peekType() != tz.TtLParen {
		return nil, p.announceError(errExpectedLParen)
	}

	p.advance()

	left, err := p.parseExpressionAsGroup()
	if err != nil {
		return nil, err
	}

	if p.peekType() != tz.TtRParen {
		return nil, p.announceError(errExpectedRParen)
	}

	p.advance()

	return left, nil
}

func (p *Parser) parseSymbol() (*Node, error) {
	if p.peekType() != tz.TtSymbol {
		return nil, p.announceError(errNotASymbol)
	}

	n := NewNodeWithValue(NodeType(p.peekType()), p.peek().Value)

	return n, nil
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
