package parser

import "vladcheck/opi/internal/tokenizer"

// Parser consumes a slice of tokens and produces an AST.
type Parser struct {
	tokens []tokenizer.Token
	// Tree is the AST produced by the most recent call to Parse.
	Tree AST
}

// NewParser returns a Parser initialized to tokenize the given tokens.
func NewParser(tokens []tokenizer.Token) Parser {
	return Parser{
		tokens: tokens,
		Tree: AST{
			Root: NewNode(NTProgram),
		},
	}
}

// Parse consumes the parser's tokens and returns the resulting AST.
func (p *Parser) Parse() (AST, error) {
	return p.Tree, nil
}
