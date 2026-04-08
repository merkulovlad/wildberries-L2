package builtin

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"syscall"
)

var ErrMissingPID = errors.New("kill: missing pid")

type Kill struct {
}

// NewKill creates a new kill builtin command
func NewKill() *Kill {
	return &Kill{}
}

// Execute sends a signal to the specified process
func (k *Kill) Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	if len(args) > 2 {
		return fmt.Errorf("kill: %w", ErrTooManyValues)
	}

	if len(args) < 2 {
		return ErrMissingPID
	}

	intNumProc, err := strconv.Atoi(args[1])

	if err != nil {
		return err
	}

	proc, err := os.FindProcess(intNumProc)

	if err != nil {
		return err
	}

	err = proc.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}
