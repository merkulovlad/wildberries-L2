package internal

import "github.com/spf13/pflag"

type Options struct {
	LinesAfterFound  int
	LinesBeforeFound int
	CountOnly        bool
	IgnoreCase       bool
	InvertMatch      bool
	ExactMatch       bool
	WriteLineNumbers bool
}

func ParseFlags() Options {
	var opts Options

	pflag.IntVarP(&opts.LinesAfterFound, "after-context", "A", 0, "Print NUM lines of trailing context after matching lines.")
	pflag.IntVarP(&opts.LinesBeforeFound, "before-context", "B", 0, "Print NUM lines of leading context before matching lines.")
	pflag.BoolVarP(&opts.CountOnly, "count", "c", false, "Only print a count of matching lines per FILE.")
	pflag.BoolVarP(&opts.IgnoreCase, "ignore-case", "i", false, "Ignore case distinctions in patterns and data.")
	pflag.BoolVarP(&opts.InvertMatch, "invert-match", "v", false, "Select non-matching lines.")
	pflag.BoolVarP(&opts.ExactMatch, "fixed-strings", "F", false, "Interpret PATTERN as a list of fixed strings, separated by newlines, any of which is to be matched.")
	pflag.BoolVarP(&opts.WriteLineNumbers, "line-number", "n", false, "Prefix each line of output with the 1-based line number within its input file.")
	pflag.Parse()

	return opts
}
