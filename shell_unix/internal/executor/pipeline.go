package executor

import (
	"io"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

type PipelineExecutor struct {
	commands []*parser.Command
}

// NewPipelineExecutor creates a new pipeline executor for the given commands
func NewPipelineExecutor(commands []*parser.Command) *PipelineExecutor {
	return &PipelineExecutor{}
}

// Execute runs all commands in the pipeline with pipes connecting them
func (p *PipelineExecutor) Execute() error {
	return nil
}

// setupPipes creates and connects pipes between commands
func (p *PipelineExecutor) setupPipes() ([]io.ReadCloser, []io.WriteCloser, error) {
	return nil, nil, nil
}

// closePipes closes all pipe file descriptors
func (p *PipelineExecutor) closePipes(readers []io.ReadCloser, writers []io.WriteCloser) {
}
