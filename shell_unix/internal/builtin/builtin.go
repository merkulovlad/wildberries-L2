package builtin

import "io"

type Command interface {
	Execute(args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error
}

type Registry struct {
	commands map[string]Command
}

// NewRegistry creates a new builtin command registry
func NewRegistry() *Registry {
	registry := &Registry{
		commands: make(map[string]Command),
	}

	registry.Register("cd", NewCD())
	registry.Register("echo", NewEcho())
	registry.Register("kill", NewKill())
	registry.Register("ps", NewPS())
	registry.Register("pwd", NewPWD())

	return registry
}

// Register adds a builtin command to the registry
func (r *Registry) Register(name string, cmd Command) {
	r.commands[name] = cmd
}

// Get retrieves a builtin command by name
func (r *Registry) Get(name string) (Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

// IsBuiltin checks if a command name is a builtin
func (r *Registry) IsBuiltin(name string) bool {
	_, ok := r.Get(name)
	return ok
}
