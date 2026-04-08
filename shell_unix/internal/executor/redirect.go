package executor

import (
	"os"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

// executeRedirect executes a command with input/output redirection
func (e *Executor) executeRedirect(redir *parser.Redirect) error {
	if redir == nil || redir.Node == nil {
		return nil
	}

	switch redir.Type {
	case parser.RedirectInput:
		file, err := os.Open(redir.Target)
		if err != nil {
			return err
		}
		defer file.Close()

		child := e.withIO(file, e.stdout, e.stderr)
		return child.Execute(redir.Node)

	case parser.RedirectOutput:
		file, err := os.Create(redir.Target)
		if err != nil {
			return err
		}
		defer file.Close()

		child := e.withIO(e.stdin, file, e.stderr)
		return child.Execute(redir.Node)

	default:
		return ErrUnknownType
	}
}
