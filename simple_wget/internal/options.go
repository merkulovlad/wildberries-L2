package internal

import (
	"time"

	"github.com/spf13/pflag"
)

type Options struct {
	Depth       int
	OutputDir   string
	Concurrency int
	Timeout     time.Duration
	UserAgent   string
	SameHost    bool
}

func ParseFlags() Options {
	var opts Options

	pflag.IntVarP(&opts.Depth, "level", "l", 1, "Maximum recursion depth (0 means only the given page).")
	pflag.StringVarP(&opts.OutputDir, "directory-prefix", "P", "site", "Directory where downloaded files will be saved.")
	pflag.IntVar(&opts.Concurrency, "concurrency", 4, "Number of concurrent downloads.")
	pflag.DurationVar(&opts.Timeout, "timeout", 30*time.Second, "Per-request timeout.")
	pflag.StringVarP(&opts.UserAgent, "user-agent", "U", "simple_wget/1.0", "User-Agent header value.")
	pflag.BoolVar(&opts.SameHost, "same-host", true, "Restrict crawling to the initial host.")
	pflag.Parse()

	return opts
}
