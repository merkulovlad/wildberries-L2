package executor

import (
	"io"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

type RedirectExecutor struct {
	redirect *parser.Redirect
}

// NewRedirectExecutor creates a new redirect executor
func NewRedirectExecutor(redirect *parser.Redirect) *RedirectExecutor {
	return &RedirectExecutor{}
}

// Execute performs the redirection and executes the command
func (r *RedirectExecutor) Execute() error {
	return nil
}

// setupInput redirects stdin from the specified file
func (r *RedirectExecutor) setupInput(target string) (io.ReadCloser, error) {
	return nil, nil
}

// setupOutput redirects stdout to the specified file
func (r *RedirectExecutor) setupOutput(target string) (io.WriteCloser, error) {
	return nil, nil
}

// restore restores original stdin/stdout after redirection
func (r *RedirectExecutor) restore() {
}
