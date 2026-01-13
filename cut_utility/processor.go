package main

// processLine extracts selected fields from a line
// Returns the processed output and a boolean indicating if line should be skipped
func processLine(line string, fields []int, opts Options) (string, bool) {
	// TODO: Check if line contains delimiter
	// TODO: If -s flag is set and no delimiter found, skip line

	// TODO: Split line by delimiter
	// TODO: Extract requested fields
	// TODO: Handle field indices that exceed available columns
	// TODO: Join selected fields with delimiter

	return "", false
}

// splitLine splits a line by delimiter
func splitLine(line string, delimiter string) []string {
	// TODO: Split line by delimiter
	// TODO: Handle special case when delimiter is tab
	return nil
}

// selectFields extracts fields at specified indices from parts
func selectFields(parts []string, fields []int) []string {
	// TODO: Extract fields at requested indices
	// TODO: Ignore indices that exceed available parts
	// TODO: Maintain order of requested fields
	return nil
}

// joinFields joins selected fields with delimiter
func joinFields(fields []string, delimiter string) string {
	// TODO: Join fields with delimiter
	return ""
}

// containsDelimiter checks if line contains the delimiter
func containsDelimiter(line string, delimiter string) bool {
	// TODO: Check if delimiter is present in line
	return false
}
