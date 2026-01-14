package builtin

type Command interface {
	Execute(args []string) error
}

type Registry struct {
	commands map[string]Command
}

// NewRegistry creates a new builtin command registry
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds a builtin command to the registry
func (r *Registry) Register(name string, cmd Command) {
}

// Get retrieves a builtin command by name
func (r *Registry) Get(name string) (Command, bool) {
	return nil, false
}

// IsBuiltin checks if a command name is a builtin
func (r *Registry) IsBuiltin(name string) bool {
	return false
}
