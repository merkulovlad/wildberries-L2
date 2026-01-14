package env

import (
	"os"
)

type Expander struct {
}

// NewExpander creates a new environment variable expander
func NewExpander() *Expander {
	return &Expander{}
}

// Expand replaces environment variables in the input string
func (e *Expander) Expand(input string) string {
	return ""
}

// ExpandArgs expands environment variables in all arguments
func (e *Expander) ExpandArgs(args []string) []string {
	return nil
}

// Get retrieves an environment variable value
func (e *Expander) Get(key string) string {
	return os.Getenv(key)
}

// Set sets an environment variable
func (e *Expander) Set(key, value string) error {
	return nil
}
