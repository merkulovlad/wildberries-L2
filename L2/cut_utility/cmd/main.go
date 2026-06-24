package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/merkulovlad/wildberries-L2/cut_utility/internal"
)

func main() {
	opts := parseFlags()

	if err := internal.Cut(os.Stdin, os.Stdout, opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags parses command-line arguments and returns Options
func parseFlags() internal.Options {
	var Fields string

	var Delimiter string

	var Separated bool

	// TODO: Parse -f flag for field specification
	flag.StringVar(&Fields, "f", "", "Comma-separated list of fields or ranges (e.g. 1,3-5)")

	// TODO: Parse -d flag for delimiter
	flag.StringVar(&Delimiter, "d", "\t", "")

	// TODO: Parse -s flag for separated mode
	flag.BoolVar(&Separated, "s", false, "Suppress lines without delimiter")

	flag.Parse()

	return internal.NewOptions(Fields, Delimiter, Separated)
}
