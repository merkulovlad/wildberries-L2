package main
import (
	"strings"
)

// processLine extracts selected fields from a line
// Returns the processed output and a boolean indicating if line should be skipped
func processLine(line string, fields []int, opts Options) (string, bool) {
	// TODO: Check if line contains delimiter
	if opts.Separated && !containsDelimiter(line, opts.Delimiter) {
		return "", true
	}
	// TODO: Split line by delimiter
	parts := splitLine(line, opts.Delimiter)
	// TODO: Extract requested fields
	selected := selectFields(parts, fields)
	output := joinFields(selected, opts.Delimiter)
	return output, false
}

// splitLine splits a line by delimiter
func splitLine(line string, delimiter string) []string {
	// TODO: Split line by delimiter
	return 	strings.Split(line, delimiter)
}

// selectFields extracts fields at specified indices from parts
func selectFields(parts []string, fields []int) []string {
	// TODO: Extract fields at requested indices
	var selected []string
	for _, idx := range fields {
		if idx-1 < len(parts) && idx-1 >= 0 {
			selected = append(selected, parts[idx-1])
		}
	}
	return selected
}

// joinFields joins selected fields with delimiter
func joinFields(fields []string, delimiter string) string {
	// TODO: Join fields with delimiter
	return strings.Join(fields, delimiter)
}

// containsDelimiter checks if line contains the delimiter
func containsDelimiter(line string, delimiter string) bool {
	// TODO: Check if delimiter is present in line
	if strings.Contains(line, delimiter) {
		return true
	}
	return false
}
