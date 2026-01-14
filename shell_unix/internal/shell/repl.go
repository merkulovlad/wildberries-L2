package shell

type REPL struct {
	shell *Shell
}

// NewREPL creates a new REPL instance for the given shell
func NewREPL(shell *Shell) *REPL {
	return &REPL{}
}

// Start begins the read-eval-print loop
func (r *REPL) Start() error {
	return nil
}

// ReadLine reads a single line of input from the user
func (r *REPL) ReadLine() (string, error) {
	return "", nil
}

// Prompt returns the shell prompt string to display
func (r *REPL) Prompt() string {
	return ""
}
