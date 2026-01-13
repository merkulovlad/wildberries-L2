package internal

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// parseFields converts field specification string into a list of field indices
// Example: "1,3-5" -> []int{1, 3, 4, 5}
var ErrInvalidFieldSpec = errors.New("invalid field specification")

func parseFields(spec string) ([]int, error) {
	// TODO: Parse field specification
	tokens := strings.Split(spec, ",")
	set := make(map[int]struct{})

	for _, token := range tokens {
		if strings.Contains(token, "-") {
			// Handle range
			parts := strings.SplitN(token, "-", 2)

			start, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, ErrInvalidFieldSpec
			}

			end, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, ErrInvalidFieldSpec
			}

			for i := start; i <= end; i++ {
				// Append each field in range
				set[i] = struct{}{}
			}
		} else {
			// Handle single field
			field, err := strconv.Atoi(token)
			if err != nil {
				return nil, ErrInvalidFieldSpec
			}

			set[field] = struct{}{}
		}
	}
	// Convert set to sorted slice
	var fields []int

	for field := range set {
		if field <= 0 {
			continue // Ignore non-positive field numbers
		}

		fields = append(fields, field)
	}

	sort.Ints(fields)

	return fields, nil
}
