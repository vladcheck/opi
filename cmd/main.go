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
)

const extension = ".opi"

var (
	errLogger = log.New(os.Stderr, "ERROR: ", 0)

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

	flag.Parse()

	if len(flag.Args()) < 1 {
		return fmt.Errorf("%w", ErrPathRequired)
	}

	path := filepath.Clean(flag.Arg(0))
	if err := validatePath(path); err != nil {
		return err
	}

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
