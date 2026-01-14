package builtin

type CD struct {
}

// NewCD creates a new cd builtin command
func NewCD() *CD {
	return &CD{}
}

// Execute changes the current working directory
func (c *CD) Execute(args []string) error {
	return nil
}
