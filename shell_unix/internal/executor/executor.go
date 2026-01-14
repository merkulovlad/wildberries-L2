package executor

import (
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/builtin"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

type Executor struct {
	builtins *builtin.Registry
}

// New creates a new command executor
func New() *Executor {
	return &Executor{}
}

// Execute executes the given AST node
func (e *Executor) Execute(node parser.Node) error {
	return nil
}

// executeCommand executes a single command (builtin or external)
func (e *Executor) executeCommand(cmd *parser.Command) error {
	return nil
}

// executePipeline executes a pipeline of commands
func (e *Executor) executePipeline(pipe *parser.Pipeline) error {
	return nil
}

// executeRedirect executes a command with input/output redirection
func (e *Executor) executeRedirect(redir *parser.Redirect) error {
	return nil
}

// executeLogical executes logical operations (&& and ||)
func (e *Executor) executeLogical(op *parser.LogicalOp) error {
	return nil
}

// executeExternal executes an external command using fork/exec
func (e *Executor) executeExternal(cmd *parser.Command) error {
	return nil
}
