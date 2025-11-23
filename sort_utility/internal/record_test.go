package internal

import (
	"math"
	"testing"
)

func TestNewRecordColumnSelectionAndTrim(t *testing.T) {
	opts := Options{
		Column:         3,
		IgnoreTrailing: true,
	}

	rec := NewRecord("field1\tfield2\tvalue  \textra", opts)

	if rec.TrimmedKey != "value" {
		t.Fatalf("expected third column without trailing spaces, got %q", rec.TrimmedKey)
	}

	opts.Column = 5 // beyond available fields; should fallback to whole line

	rec = NewRecord("a\tb\tc", opts)
	if rec.TrimmedKey != "a\tb\tc" {
		t.Fatalf("expected entire line to be used as key when column missing, got %q", rec.TrimmedKey)
	}
}

func TestNewRecordNumericAndMonthValidation(t *testing.T) {
	numOpts := Options{Numeric: true}

	numRecord := NewRecord("10", numOpts)
	if !numRecord.NumValid || numRecord.NumericValue != 10 {
		t.Fatalf("expected numeric value 10, got %.1f (valid=%v)", numRecord.NumericValue, numRecord.NumValid)
	}

	numInvalid := NewRecord("not-num", numOpts)
	if numInvalid.NumValid {
		t.Fatalf("expected invalid numeric to be flagged")
	}

	monthOpts := Options{Month: true}

	monthRecord := NewRecord("February", monthOpts)
	if !monthRecord.MonthValid || monthRecord.MonthValue != 2 {
		t.Fatalf("expected month February to parse to 2, got (%v, %d)", monthRecord.MonthValid, monthRecord.MonthValue)
	}

	monthInvalid := NewRecord("F", monthOpts)
	if monthInvalid.MonthValid {
		t.Fatalf("expected short month string to be invalid")
	}
}

func TestNumericNan(t *testing.T) {
	opts := Options{Numeric: true}
	lines := []string{"10", "NaN", "001"}
	records := make([]*Record, 0, len(lines))

	for _, line := range lines {
		records = append(records, NewRecord(line, opts))
	}

	SortRecords(records, opts)
	got := []string{records[0].OriginalLine, records[1].OriginalLine, records[2].OriginalLine}
	want := []string{"001", "10", "NaN"}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("numeric sort order mismatch: got %v want %v", got, want)
		}
	}
}

func TestNewRecordHumanNumeric(t *testing.T) {
	opts := Options{HumanNumeric: true}

	valid := NewRecord("1.25G", opts)
	want := 1.25 * 1024 * 1024 * 1024

	if !valid.HumanNumericValid {
		t.Fatalf("expected human numeric value to be valid")
	}

	if math.Abs(valid.HumanNumericValue-want) > 1e-6 {
		t.Fatalf("expected %f, got %f", want, valid.HumanNumericValue)
	}

	invalid := NewRecord("abcZ", opts)
	if invalid.HumanNumericValid {
		t.Fatalf("expected invalid human numeric string to be rejected")
	}
}

func TestParseHumanNumeric(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"10", 10},
		{"2K", 2 * 1024},
		{"3m", 3 * 1024 * 1024},
		{"0.5T", 0.5 * 1024 * 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		got, err := parseHumanNumeric(tt.input)
		if err != nil {
			t.Fatalf("parseHumanNumeric(%q) unexpected error: %v", tt.input, err)
		}

		if math.Abs(got-tt.want) > 1e-9 {
			t.Fatalf("parseHumanNumeric(%q) = %f, want %f", tt.input, got, tt.want)
		}
	}

	invalid := []string{"", "abc", "12Q"}
	for _, input := range invalid {
		if _, err := parseHumanNumeric(input); err == nil {
			t.Fatalf("parseHumanNumeric(%q) expected error", input)
		}
	}
}
