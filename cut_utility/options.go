package main

// Options contains configuration for the cut utility
type Options struct {
	Fields    string // Field specification (e.g., "1,3-5")
	Delimiter string // Field delimiter (default: tab)
	Separated bool   // Output only lines containing delimiter
}
