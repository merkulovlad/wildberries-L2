package executor

import (
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/builtin"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/process"
)

var ErrUnknownType = errors.New("unknown node type")

type Executor struct {
	builtins *builtin.Registry
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
	state    *foregroundState
}

type foregroundState struct {
	current *process.Group
}

// New creates a new command executor
func New() *Executor {
	return &Executor{
		builtins: builtin.NewRegistry(),
		stdin:    os.Stdin,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		state:    &foregroundState{},
	}
}

func (e *Executor) withIO(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Executor {
	if stdin == nil {
		stdin = e.stdin
	}
	if stdout == nil {
		stdout = e.stdout
	}
	if stderr == nil {
		stderr = e.stderr
	}

	return &Executor{
		builtins: e.builtins,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
		state:    e.state,
	}
}

func (e *Executor) setForeground(group *process.Group) {
	e.state.current = group
}

func (e *Executor) clearForeground() {
	e.state.current = nil
}

func (e *Executor) InterruptForeground() error {
	if e.state.current == nil {
		return nil
	}

	return e.state.current.KillAll()
}

// Execute executes the given AST node
func (e *Executor) Execute(node parser.Node) error {
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	case *parser.Command:
		return e.executeCommand(n)
	case *parser.Pipeline:
		return e.executePipeline(n)
	case *parser.Redirect:
		return e.executeRedirect(n)
	case *parser.LogicalOp:
		return e.executeLogical(n)
	default:
		return ErrUnknownType
	}
}

// executeCommand executes a single command (builtin or external)
func (e *Executor) executeCommand(cmd *parser.Command) error {
	if cmd == nil || cmd.Name == "" {
		return nil
	}

	if b, ok := e.builtins.Get(cmd.Name); ok {
		return b.Execute(commandArgs(cmd), e.stdin, e.stdout, e.stderr)
	}

	return e.executeExternal(cmd)
}

func commandArgs(cmd *parser.Command) []string {
	args := make([]string, 0, len(cmd.Args)+1)
	args = append(args, cmd.Name)
	args = append(args, cmd.Args...)
	return args
}

// executeLogical executes logical operations (&& and ||)
func (e *Executor) executeLogical(op *parser.LogicalOp) error {
	if op == nil {
		return nil
	}

	switch op.Operator {
	case parser.OpAnd:
		if err := e.Execute(op.Left); err != nil {
			return err
		}
		return e.Execute(op.Right)
	case parser.OpOr:
		if err := e.Execute(op.Left); err == nil {
			return nil
		}
		return e.Execute(op.Right)
	default:
		return ErrUnknownType
	}
}

// executeExternal executes an external command using fork/exec
func (e *Executor) executeExternal(cmd *parser.Command) error {
	c := exec.Command(cmd.Name, cmd.Args...)
	c.Stdin = e.stdin
	c.Stdout = e.stdout
	c.Stderr = e.stderr

	proc := process.New(c)
	group := e.state.current
	ownsForeground := false
	if group == nil {
		group = process.NewGroup()
		e.setForeground(group)
		ownsForeground = true
	}
	group.Add(proc)
	defer group.Remove(proc.PID())

	if ownsForeground {
		defer e.clearForeground()
	}

	if err := proc.Start(); err != nil {
		group.Remove(proc.PID())
		return err
	}

	return proc.Wait()
}
