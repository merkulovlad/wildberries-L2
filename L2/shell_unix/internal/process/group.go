package process

import "sync"

type Group struct {
	processes []*Process
	mu        sync.Mutex
}

// NewGroup creates a new process group
func NewGroup() *Group {
	return &Group{
		processes: make([]*Process, 0, 4),
	}
}

// Add adds a process to the group
func (g *Group) Add(p *Process) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.processes = append(g.processes, p)
}

// Remove removes a process from the group by PID
func (g *Group) Remove(pid int) {
	g.mu.Lock()
	defer g.mu.Unlock()

	for i, proc := range g.processes {
		if proc.pid == pid {
			last := len(g.processes) - 1
			g.processes[i] = g.processes[last]
			g.processes[last] = nil
			g.processes = g.processes[:last]
			return
		}
	}
}

// List returns all processes in the group
func (g *Group) List() []*Process {
	g.mu.Lock()
	defer g.mu.Unlock()

	res := make([]*Process, len(g.processes))
	copy(res, g.processes)
	return res
}

// KillAll terminates all processes in the group
func (g *Group) KillAll() error {
	for _, proc := range g.List() {
		err := proc.Kill()
		if err != nil {
			return err
		}
	}

	return nil
}

// WaitAll waits for all processes in the group to complete
func (g *Group) WaitAll() error {
	for _, proc := range g.List() {
		if err := proc.Wait(); err != nil {
			return err
		}
	}
	return nil
}
