package builtin

import (
	"fmt"
	"io"
	"strings"
)

type Echo struct {
}

// NewEcho creates a new echo builtin command
func NewEcho() *Echo {
	return &Echo{}
}

// Execute prints the arguments to stdout
func (e *Echo) Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	values := []string{}
	if len(args) > 1 {
		values = args[1:]
	}

	_, err := fmt.Fprintln(stdout, strings.Join(values, " "))
	return err
}
