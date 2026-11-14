package main

import (
	"fmt"
	"vladcheck/opi/internal/parser"
	"vladcheck/opi/internal/tokenizer"
)

func executionLoop(content string) error {
	t := tokenizer.NewTokenizer(content)

	tokens, err := t.Tokenize()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tokenizer.DebugPrintAllTokens(tokens)

	p := parser.NewParser(tokens)

	root, err := p.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	parser.DebugPrintAST(root)

	return nil
}
