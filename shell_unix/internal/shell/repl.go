package shell

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type REPL struct {
	reader *bufio.Reader
}

// NewREPL creates a new REPL instance for the given shell
func NewREPL() *REPL {
	return &REPL{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadLine reads a single line of input from the user
func (r *REPL) ReadLine() (string, error) {
	if r.reader == nil {
		r.reader = bufio.NewReader(os.Stdin)
	}

	line, err := r.reader.ReadString('\n')
	line = strings.TrimRight(line, "\r\n")

	if err == nil {
		return line, nil
	}
	if err == io.EOF && line != "" {
		return line, nil
	}

	return "", err
}

// Prompt returns the shell prompt string to display
func (r *REPL) Prompt() string {
	return "$ "
}
