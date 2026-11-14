package main

import (
	"fmt"
	"vladcheck/opi/internal/parser"
	"vladcheck/opi/internal/tokenizer"
)

func executionLoop(content string, debug bool) error {
	t := tokenizer.NewTokenizer(content)

	tokens, err := t.Tokenize()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if debug {
		tokenizer.DebugPrintAllTokens(tokens)
	}

	p := parser.NewParser(tokens)

	root, err := p.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if debug {
		parser.DebugPrintAST(root)
	}

	return nil
}
