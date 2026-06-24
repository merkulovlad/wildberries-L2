package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/merkulovlad/wildberries-L2/simple_wget/internal"
	"github.com/spf13/pflag"
)

func main() {
	opts := internal.ParseFlags()
	args := pflag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "usage: simple_wget [flags] URL")
		os.Exit(2)
	}

	start, err := url.Parse(args[0])
	if err != nil || (start.Scheme != "http" && start.Scheme != "https") {
		fmt.Fprintf(os.Stderr, "simple_wget: invalid URL %q\n", args[0])
		os.Exit(2)
	}

	crawler := internal.NewCrawler(opts, start)
	crawler.Run(start)
	fmt.Println(crawler.Summary())
}
