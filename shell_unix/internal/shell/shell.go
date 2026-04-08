package shell

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/env"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/executor"
	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

type Shell struct {
	repl     *REPL
	executor *executor.Executor
	expander *env.Expander
}

// New creates and initializes a new Shell instance
func New() *Shell {
	return &Shell{
		repl:     NewREPL(),
		executor: executor.New(),
		expander: env.NewExpander(),
	}
}

// Run starts the main shell loop, handling user input and command execution
func (s *Shell) Run() error {
	sigHandler := NewSignalHandler()
	sigHandler.SetupWithInterrupt()
	defer sigHandler.Stop()

	go func() {
		for range sigHandler.C() {
			_ = s.executor.InterruptForeground()
		}
	}()

	for {
		fmt.Print(s.repl.Prompt())

		line, err := s.repl.ReadLine()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		node, err := parser.New(line).Parse()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		node = s.expander.ExpandNode(node)

		if err := s.executor.Execute(node); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
