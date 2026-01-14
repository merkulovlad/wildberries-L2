package builtin

type PS struct {
}

// NewPS creates a new ps builtin command
func NewPS() *PS {
	return &PS{}
}

// Execute lists all running processes
func (p *PS) Execute(args []string) error {
	return nil
}
