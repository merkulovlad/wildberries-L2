package builtin

import (
	"io"
	"os/exec"
)

type PS struct {
}

// NewPS creates a new ps builtin command
func NewPS() *PS {
	return &PS{}
}

// Execute lists all running processes
func (p *PS) Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := exec.Command("ps")
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}
