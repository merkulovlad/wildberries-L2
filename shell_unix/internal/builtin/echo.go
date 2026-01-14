package builtin

type Echo struct {
}

// NewEcho creates a new echo builtin command
func NewEcho() *Echo {
	return &Echo{}
}

// Execute prints the arguments to stdout
func (e *Echo) Execute(args []string) error {
	return nil
}
