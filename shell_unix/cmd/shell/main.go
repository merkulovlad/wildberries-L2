package main

import (
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/shell"
)

// main initializes and starts the shell REPL
func main() {
	sh := shell.New()
	sh.Run()
}
