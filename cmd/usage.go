package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func usage() {
	_, err := fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s path/to/a/file\n", filepath.Base(os.Args[0]))
	if err != nil {
		errLogger.Fatal(err)
	}
}
