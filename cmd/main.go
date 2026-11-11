//nolint:gochecknoglobals
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"vladcheck/opi/internal/parser"
	"vladcheck/opi/internal/tokenizer"
)

const extension = ".opi"

var (
	errLogger = log.New(os.Stderr, "", 0)

	ErrPathRequired   = errors.New("filepath is not specified")
	ErrDoesNotExist   = errors.New("path does not exist")
	ErrNotADirectory  = errors.New("path must not be a directory")
	ErrWrongExtension = errors.New("wrong file extension")
)

func main() {
	err := run()
	if err != nil {
		errLogger.Fatal(err)
	}
}

func run() error {
	flag.Usage = func() {
		usage()
		flag.PrintDefaults()
	}

	var execPtr string
	flag.StringVar(&execPtr, "exec", "", "--exec <code>")
	flag.StringVar(&execPtr, "e", "", "-e <code>")

	flag.Parse()

	if len(execPtr) > 0 {
		err := executionLoop(execPtr)
		if err != nil {
			return err
		}

		return nil
	}

	if len(flag.Args()) == 0 {
		return fmt.Errorf("%w", ErrPathRequired)
	}

	path := filepath.Clean(flag.Arg(0))

	err := validatePath(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	err = executionLoop(string(content))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

func executionLoop(content string) error {
	t := tokenizer.NewTokenizer(content)

	tokens, err := t.Tokenize()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tokenizer.DebugPrintAllTokens(tokens)

	p := parser.NewParser(tokens)

	ast, err := p.Parse()
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	parser.DebugPrintAST(ast)

	return nil
}

func usage() {
	_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s path/to/a/file\n", filepath.Base(os.Args[0]))
	if err != nil {
		errLogger.Fatal(err)
	}
}

func validatePath(path string) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	defer func() { _ = file.Close() }()

	info, err := file.Stat()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("%w: %s", ErrDoesNotExist, path)
		}

		return fmt.Errorf("%w", err)
	} else if info.IsDir() {
		return fmt.Errorf("%w: %s", ErrNotADirectory, path)
	} else if got := filepath.Ext(filepath.Clean(path)); !strings.EqualFold(got, extension) {
		return fmt.Errorf("%w: want %q, got %q", ErrWrongExtension, extension, got)
	}

	return nil
}
