package main

import (
	"bytes"
	"strings"
	"testing"
)

// TestCutSingleField tests extracting a single field
func TestCutSingleField(t *testing.T) {
	input := "field1\tfield2\tfield3\n"
	expected := "field2\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "2",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutMultipleFields tests extracting multiple non-consecutive fields
func TestCutMultipleFields(t *testing.T) {
	input := "a\tb\tc\td\te\n"
	expected := "a\tc\te\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "1,3,5",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutFieldRanges tests extracting field ranges
func TestCutFieldRanges(t *testing.T) {
	input := "one\ttwo\tthree\tfour\tfive\n"
	expected := "two\tthree\tfour\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "2-4",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutCombinedFieldSpec tests combination of fields and ranges
func TestCutCombinedFieldSpec(t *testing.T) {
	input := "1\t2\t3\t4\t5\t6\t7\t8\n"
	expected := "1\t3\t4\t5\t7\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "1,3-5,7",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutCustomDelimiter tests using a custom delimiter
func TestCutCustomDelimiter(t *testing.T) {
	input := "apple,banana,cherry,date\n"
	expected := "banana,date\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "2,4",
		Delimiter: ",",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutSeparatedFlag tests the -s flag
func TestCutSeparatedFlag(t *testing.T) {
	input := "a\tb\tc\nno-delimiter-here\nx\ty\tz\n"
	expected := "a\nx\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "1",
		Delimiter: "\t",
		Separated: true,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutFieldsExceedColumns tests when requested fields exceed available columns
func TestCutFieldsExceedColumns(t *testing.T) {
	input := "a\tb\tc\n"
	expected := "b\tc\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "2,3,4,5,10",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutEmptyInput tests handling of empty input
func TestCutEmptyInput(t *testing.T) {
	input := ""
	expected := ""

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "1",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestCutDefaultDelimiter tests default tab delimiter
func TestCutDefaultDelimiter(t *testing.T) {
	input := "first\tsecond\tthird\n"
	expected := "first\tthird\n"

	reader := strings.NewReader(input)
	writer := &bytes.Buffer{}

	opts := Options{
		Fields:    "1,3",
		Delimiter: "\t",
		Separated: false,
	}

	err := cut(reader, writer, opts)
	if err != nil {
		t.Fatalf("cut() error = %v", err)
	}

	if writer.String() != expected {
		t.Errorf("cut() = %q, want %q", writer.String(), expected)
	}
}

// TestParseFields tests field specification parsing
func TestParseFields(t *testing.T) {
	tests := []struct {
		name     string
		spec     string
		expected []int
		wantErr  bool
	}{
		{"single field", "1", []int{1}, false},
		{"multiple fields", "1,3,5", []int{1, 3, 5}, false},
		{"range", "2-4", []int{2, 3, 4}, false},
		{"combination", "1,3-5,7", []int{1, 3, 4, 5, 7}, false},
		{"invalid spec", "abc", nil, true},       // Returns error
		{"negative field", "-1", nil, true},      // Returns error
		{"invalid range", "5-2", []int{}, false}, // Empty result for invalid range
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFields(tt.spec)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !equalSlices(got, tt.expected) {
				t.Errorf("parseFields() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestSplitLine tests line splitting logic
func TestSplitLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		delimiter string
		expected  []string
	}{
		{"tab delimiter", "a\tb\tc", "\t", []string{"a", "b", "c"}},
		{"comma delimiter", "x,y,z", ",", []string{"x", "y", "z"}},
		{"multi-char delimiter", "foo::bar::baz", "::", []string{"foo", "bar", "baz"}},
		{"no delimiter", "single", "\t", []string{"single"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitLine(tt.line, tt.delimiter)
			if !equalSlices(got, tt.expected) {
				t.Errorf("splitLine() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestSelectFields tests field selection from parts
func TestSelectFields(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		fields   []int
		expected []string
	}{
		{"valid indices", []string{"a", "b", "c", "d"}, []int{1, 3}, []string{"a", "c"}},
		{"exceed parts", []string{"a", "b"}, []int{1, 2, 3, 4}, []string{"a", "b"}},
		{"ordered fields", []string{"x", "y", "z"}, []int{3, 1, 2}, []string{"z", "x", "y"}},
		{"empty fields", []string{"a", "b", "c"}, []int{}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := selectFields(tt.parts, tt.fields)
			if !equalSlices(got, tt.expected) {
				t.Errorf("selectFields() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestContainsDelimiter tests delimiter detection
func TestContainsDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		delimiter string
		expected  bool
	}{
		{"has delimiter", "a\tb\tc", "\t", true},
		{"no delimiter", "abc", "\t", false},
		{"empty line", "", "\t", false},
		{"multi-char delimiter present", "foo::bar", "::", true},
		{"multi-char delimiter absent", "foobar", "::", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := containsDelimiter(tt.line, tt.delimiter)
			if got != tt.expected {
				t.Errorf("containsDelimiter() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Helper function to compare slices
func equalSlices[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
