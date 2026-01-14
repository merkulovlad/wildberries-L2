package builtin

type Kill struct {
}

// NewKill creates a new kill builtin command
func NewKill() *Kill {
	return &Kill{}
}

// Execute sends a signal to the specified process
func (k *Kill) Execute(args []string) error {
	return nil
}
