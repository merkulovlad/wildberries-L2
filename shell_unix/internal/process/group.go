package process

type Group struct {
	processes []*Process
}

// NewGroup creates a new process group
func NewGroup() *Group {
	return &Group{}
}

// Add adds a process to the group
func (g *Group) Add(p *Process) {
}

// Remove removes a process from the group by PID
func (g *Group) Remove(pid int) {
}

// List returns all processes in the group
func (g *Group) List() []*Process {
	return nil
}

// KillAll terminates all processes in the group
func (g *Group) KillAll() error {
	return nil
}

// WaitAll waits for all processes in the group to complete
func (g *Group) WaitAll() error {
	return nil
}
