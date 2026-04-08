package builtin

import (
	"fmt"
	"io"
	"os"
)

type PWD struct {
}

// NewPWD creates a new pwd builtin command
func NewPWD() *PWD {
	return &PWD{}
}

// Execute prints the current working directory
func (p *PWD) Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	_, err = fmt.Fprintln(stdout, dir)
	return err
}
