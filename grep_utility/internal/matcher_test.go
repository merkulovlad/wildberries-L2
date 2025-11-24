package internal

import "testing"

type test struct {
	name     string
	pattern  string
	opts     Options
	lines    []string
	expected []int // expected line numbers that match
}

var tests = []test{
	{
		name:    "Simple match",
		pattern: "hello",
		opts:    Options{},
		lines: []string{
			"hello world",
			"goodbye world",
			"hello again",
		},
		expected: []int{1, 3},
	},
	{
		name:    "Ignore case",
		pattern: "HELLO",
		opts:    Options{IgnoreCase: true},
		lines: []string{
			"hello world",
			"goodbye world",
			"HELLO again",
		},
		expected: []int{1, 3},
	},
	{
		name:    "Invert match",
		pattern: "hello",
		opts:    Options{InvertMatch: true},
		lines: []string{
			"hello world",
			"goodbye world",
			"hello again",
		},
		expected: []int{2},
	},
	{
		name:    "Exact match",
		pattern: "hello",
		opts:    Options{ExactMatch: true},
		lines: []string{
			"hello world",
			"say hello",
			"helloagain",
		},
		expected: []int{1, 2, 3},
	},
	{
		name:    "No matches",
		pattern: "xyz",
		opts: Options{
			WriteLineNumbers: true,
		},
		lines: []string{
			"hello world",
			"goodbye world",
			"hello again",
		},
		expected: []int{},
	},
}

func TestGrepLines(t *testing.T) {
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GrepLines(tc.lines, tc.pattern, tc.opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != len(tc.expected) {
				t.Fatalf("expected %d matches, got %d", len(tc.expected), len(result))
			}

			for i, line := range result {
				if line.Number != tc.expected[i] {
					t.Errorf("expected line number %d, got %d", tc.expected[i], line.Number)
				}
			}
		})
	}
}
func TestGrepLines_RegexMatch(t *testing.T) {
	opts := Options{
		// ExactMatch = false → используем regexp
		IgnoreCase: true,
	}

	lines := []string{
		"hello world",
		"hxllo world",
		"heLLo",
		"world",
	}

	// шаблон: h.anythingllo (регулярка)
	pattern := `h.*llo`

	result, err := GrepLines(lines, pattern, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ожидаем, что совпали строки 1, 2, 3 (по твоей нумерации 1-based)
	wantNums := []int{1, 2, 3}
	gotNums := make([]int, 0, len(result))

	for _, line := range result {
		if line.Match {
			gotNums = append(gotNums, line.Number)
		}
	}

	if len(gotNums) != len(wantNums) {
		t.Fatalf("expected %d matches, got %d", len(wantNums), len(gotNums))
	}

	for i, n := range wantNums {
		if gotNums[i] != n {
			t.Errorf("expected matched line %d, got %d", n, gotNums[i])
		}
	}
}

func TestGrepLines_ExactMatchIgnoreCase(t *testing.T) {
	opts := Options{
		ExactMatch: true,
		IgnoreCase: true,
	}

	lines := []string{
		"HELLO world",
		"say hello",
		"heLLoKitty",
		"no match here",
	}

	pattern := "hello"

	result, err := GrepLines(lines, pattern, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// здесь твой ExactMatch = fixed substring, поэтому все три строки с hello должны матчиться
	wantNums := []int{1, 2, 3}
	gotNums := make([]int, 0, len(result))

	for _, line := range result {
		if line.Match {
			gotNums = append(gotNums, line.Number)
		}
	}

	if len(gotNums) != len(wantNums) {
		t.Fatalf("expected %d matches, got %d", len(wantNums), len(gotNums))
	}

	for i, n := range wantNums {
		if gotNums[i] != n {
			t.Errorf("expected matched line %d, got %d", n, gotNums[i])
		}
	}
}

func TestGrepLines_ContextAfter(t *testing.T) {
	opts := Options{
		LinesAfterFound: 2, // -A 2
	}

	lines := []string{
		"00",
		"11 MATCH",
		"22",
		"33",
		"44",
	}

	pattern := "MATCH"

	result, err := GrepLines(lines, pattern, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// ожидаем, что:
	// Match на строке 2
	// контекст -A2 → строки 3 и 4
	// Printed должны быть 2,3,4
	wantPrinted := []int{2, 3, 4}
	gotPrinted := make([]int, 0, len(result))

	for _, line := range result {
		if line.Printed {
			gotPrinted = append(gotPrinted, line.Number)
		}
	}

	if len(gotPrinted) != len(wantPrinted) {
		t.Fatalf("expected %d printed lines, got %d", len(wantPrinted), len(gotPrinted))
	}

	for i, n := range wantPrinted {
		if gotPrinted[i] != n {
			t.Errorf("expected printed line %d, got %d", n, gotPrinted[i])
		}
	}
}

func TestGrepLines_ContextOverlappingAfter(t *testing.T) {
	opts := Options{
		LinesAfterFound: 2, // -A 2
	}

	lines := []string{
		"00",
		"11 MATCH",
		"22 MATCH",
		"33",
		"44",
		"55",
	}

	pattern := "MATCH"

	result, err := GrepLines(lines, pattern, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Должно получиться:
	// MATCH на 2 → контекст 3,4
	// MATCH на 3 → контекст 4,5
	// итого Printed: 2,3,4,5
	wantPrinted := []int{2, 3, 4, 5}
	gotPrinted := make([]int, 0, len(result))

	for _, line := range result {
		if line.Printed {
			gotPrinted = append(gotPrinted, line.Number)
		}
	}

	if len(gotPrinted) != len(wantPrinted) {
		t.Fatalf("expected %d printed lines, got %d", len(wantPrinted), len(gotPrinted))
	}

	for i, n := range wantPrinted {
		if gotPrinted[i] != n {
			t.Errorf("expected printed line %d, got %d", n, gotPrinted[i])
		}
	}
}

func TestGrepLines_ContextBefore(t *testing.T) {
	opts := Options{
		LinesBeforeFound: 2, // -B 2
	}

	lines := []string{
		"00",
		"11",
		"22 MATCH",
		"33",
	}

	pattern := "MATCH"

	result, err := GrepLines(lines, pattern, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// MATCH на 3 → контекст -B2 → 1,2
	// Printed: 1,2,3
	wantPrinted := []int{1, 2, 3}
	gotPrinted := make([]int, 0, len(result))

	for _, line := range result {
		if line.Printed {
			gotPrinted = append(gotPrinted, line.Number)
		}
	}

	if len(gotPrinted) != len(wantPrinted) {
		t.Fatalf("expected %d printed lines, got %d", len(wantPrinted), len(gotPrinted))
	}

	for i, n := range wantPrinted {
		if gotPrinted[i] != n {
			t.Errorf("expected printed line %d, got %d", n, gotPrinted[i])
		}
	}
}
