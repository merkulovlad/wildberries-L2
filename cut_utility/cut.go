package main

import (
	"bufio"
	"fmt"
	"io"
)

// cut reads from r, processes lines according to opts, and writes to w
func cut(r io.Reader, w io.Writer, opts Options) error {
	// TODO: Validate options

	// TODO: Parse field specification into usable format
	fields, err := parseFields(opts.Fields)
	if err != nil {
		return err
	}

	// TODO: Process input line by line
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		// TODO: Process line and extract selected fields
		output, skip := processLine(line, fields, opts)
		if skip {
			continue
		}

		// TODO: Write output
		if _, err := fmt.Fprintln(w, output); err != nil {
			return err
		}
	}

	// TODO: Handle scanner errors
	return scanner.Err()
}
