package internal

// Options contains configuration for the cut utility
type Options struct {
	Fields    string // Field specification (e.g., "1,3-5")
	Delimiter string // Field delimiter (default: tab)
	Separated bool   // Output only lines containing delimiter
}

func NewOptions(fields, delimiter string, separated bool) Options {
	return Options{
		Fields:    fields,
		Delimiter: delimiter,
		Separated: separated,
	}
}
