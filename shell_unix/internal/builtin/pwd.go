package builtin

type PWD struct {
}

// NewPWD creates a new pwd builtin command
func NewPWD() *PWD {
	return &PWD{}
}

// Execute prints the current working directory
func (p *PWD) Execute(args []string) error {
	return nil
}
