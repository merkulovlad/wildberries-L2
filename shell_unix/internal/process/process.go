package process

import (
	"fmt"
	"os/exec"
	"syscall"
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
	return &Process{
		cmd: cmd,
	}
}

// Start begins execution of the process
func (p *Process) Start() error {
	err := p.cmd.Start()
	if err != nil {
		return err
	}
	p.pid = p.cmd.Process.Pid
	p.status = StatusRunning
	return nil
}

// Wait blocks until the process completes
func (p *Process) Wait() error {
	err := p.cmd.Wait()
	p.status = StatusDone
	return err
}

// Kill terminates the process
func (p *Process) Kill() error {
	err := p.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("Error killing process: %w", err)
	}
	p.status = StatusDone
	return nil
}

// PID returns the process ID
func (p *Process) PID() int {
	return p.pid
}

// Status returns the current status of the process
func (p *Process) Status() ProcessStatus {
	return p.status
}
