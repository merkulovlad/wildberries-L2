package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	opts := parseFlags()

	if err := cut(os.Stdin, os.Stdout, opts); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// parseFlags parses command-line arguments and returns Options
func parseFlags() Options {
	var opts Options

	// TODO: Parse -f flag for field specification
	flag.StringVar(&opts.Fields, "f", "", "")

	// TODO: Parse -d flag for delimiter
	flag.StringVar(&opts.Delimiter, "d", "\t", "")

	// TODO: Parse -s flag for separated mode
	flag.BoolVar(&opts.Separated, "s", false, "")

	flag.Parse()

	return opts
}
