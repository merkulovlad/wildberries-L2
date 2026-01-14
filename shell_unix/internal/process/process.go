package process

import (
	"os/exec"
)

type Process struct {
	cmd    *exec.Cmd
	pid    int
	status ProcessStatus
}

type ProcessStatus int

const (
	StatusRunning ProcessStatus = iota
	StatusStopped
	StatusDone
)

// New creates a new process wrapper for the given command
func New(cmd *exec.Cmd) *Process {
	return &Process{}
}

// Start begins execution of the process
func (p *Process) Start() error {
	return nil
}

// Wait blocks until the process completes
func (p *Process) Wait() error {
	return nil
}

// Kill terminates the process
func (p *Process) Kill() error {
	return nil
}

// PID returns the process ID
func (p *Process) PID() int {
	return 0
}

// Status returns the current status of the process
func (p *Process) Status() ProcessStatus {
	return StatusRunning
}
