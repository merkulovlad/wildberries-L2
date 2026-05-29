package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/merkulovlad/wildberries-L2/telnet_clinet/internal"
)

func main() {
	timeout := flag.Duration("timeout", 10*time.Second, "connection timeout")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: go-telnet [--timeout=DURATION] host port\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		os.Exit(2)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := internal.Run(ctx, args[0], args[1], *timeout); err != nil {
		fmt.Fprintf(os.Stderr, "go-telnet: %v\n", err)
		os.Exit(1)
	}
}
