package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/merkulovlad/wildberries-L2/grep_utility/internal"
	"github.com/spf13/pflag"
)

func main() {
	opts := internal.ParseFlags()
	args := pflag.Args()
	pattern := args[0]
	files := args[1:]

	if pattern == "" {
		fmt.Fprintln(os.Stderr, "grep: empty pattern")
		os.Exit(2)
	}

	lines := make([]string, 0, 500)

	if len(files) == 0 {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "grep: error reading stdin:", err)
			os.Exit(2)
		}
	} else {
		for _, file := range files {
			// #nosecG304 - file path comes from user input
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %v\n", file, err)
				os.Exit(2)
			}

			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "grep: error reading %s: %v\n", file, err)
				os.Exit(2)
			}

			_ = f.Close()
		}
	}

	answer, err := internal.GrepLines(lines, pattern, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "grep: error during matching:", err)
		os.Exit(2)
	}

	printGrepResults(answer, opts)
}

func printGrepResults(answer []*internal.Line, opts internal.Options) {
	if opts.CountOnly {
		fmt.Println(len(answer))
	} else {
		for _, line := range answer {
			{
				if opts.WriteLineNumbers {
					fmt.Printf("%d:%s\n", line.Number, line.Text)
				} else {
					fmt.Println(line.Text)
				}
			}
		}
	}
}
