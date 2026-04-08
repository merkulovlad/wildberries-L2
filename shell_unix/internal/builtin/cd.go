package builtin

import (
	"fmt"
	"io"
	"os"
)

type CD struct {
}

// NewCD creates a new cd builtin command
func NewCD() *CD {
	return &CD{}
}

// Execute changes the current working directory
func (c *CD) Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if len(args) < 2 {
		return fmt.Errorf("cd: %w", ErrMissingPath)
	}
	if len(args) > 2 {
		return fmt.Errorf("cd: %w", ErrTooManyValues)
	}

	if err := os.Chdir(args[1]); err != nil {
		return err
	}

	return nil
}
