package internal

import (
	"sort"
)

func CheckIsSorted(records []*Record, opts Options) bool {
	for i := 1; i < len(records); i++ {
		if less(records[i], records[i-1], opts) {
			return false
		}
	}

	return true
}

func SortRecords(records []*Record, opts Options) {
	sort.SliceStable(records, func(i, j int) bool {
		return less(records[i], records[j], opts)
	})
}

func less(a, b *Record, opts Options) bool {
	reverse := opts.Reverse
	if reverse {
		return !less(a, b, Options{
			Column:         opts.Column,
			Numeric:        opts.Numeric,
			Month:          opts.Month,
			HumanNumeric:   opts.HumanNumeric,
			IgnoreTrailing: opts.IgnoreTrailing,
		})
	}

	if opts.Numeric {
		if a.NumValid && b.NumValid {
			return a.NumericValue < b.NumericValue
		} else if a.NumValid {
			return true
		} else if b.NumValid {
			return false
		}
	}

	if opts.Month {
		if a.MonthValid && b.MonthValid {
			return a.MonthValue < b.MonthValue
		} else if a.MonthValid {
			return true
		} else if b.MonthValid {
			return false
		}
	}

	if opts.HumanNumeric {
		if a.HumanNumericValid && b.HumanNumericValid {
			return a.HumanNumericValue < b.HumanNumericValue
		} else if a.HumanNumericValid {
			return true
		} else if b.HumanNumericValid {
			return false
		}
	}

	return a.TrimmedKey < b.TrimmedKey
}
