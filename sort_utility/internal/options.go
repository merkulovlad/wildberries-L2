package internal

import (
	"github.com/spf13/pflag"
)

type Options struct {
	Column         int  // -k N
	Numeric        bool // -n
	Reverse        bool // -r
	Unique         bool // -u
	Month          bool // -M
	IgnoreTrailing bool // -b
	CheckOnly      bool // -c
	HumanNumeric   bool // -hu
}

func ParseFlags() Options {
	opts := &Options{}

	pflag.IntVarP(&opts.Column, "key", "k", 1, "specify column for sorting (1-based index)")
	pflag.BoolVarP(&opts.Numeric, "numeric", "n", false, "compare according to string numerical value")
	pflag.BoolVarP(&opts.Reverse, "reverse", "r", false, "sort in reverse order")
	pflag.BoolVarP(&opts.Unique, "unique", "u", false, "output only the first of an equal run")
	pflag.BoolVarP(&opts.Month, "month", "M", false, "compare according to month name")
	pflag.BoolVarP(&opts.IgnoreTrailing, "ignore-trailing-blanks", "b", false, "ignore trailing blanks in sort keys")
	pflag.BoolVarP(&opts.CheckOnly, "check", "c", false, "check if the file is sorted")
	pflag.BoolVarP(&opts.HumanNumeric, "human-numeric", "hu", false, "compare according to human readable sizes (e.g., 2K, 1G)")

	pflag.Parse()

	return *opts
}
