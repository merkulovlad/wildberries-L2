package shell

import (
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/executor"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

type Shell struct {
	parser   *parser.Parser
	executor *executor.Executor
}

// New creates and initializes a new Shell instance
func New() *Shell {
	return &Shell{}
}

// Run starts the main shell loop, handling user input and command execution
func (s *Shell) Run() error {
	return nil
}
