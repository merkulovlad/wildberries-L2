package parser

import (
	"errors"
	"testing"
)

func TestParseRejectsInvalidInput(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		input string
	}{
		{
			name:  "trailing and",
			input: "echo hi &&",
		},
		{
			name:  "trailing or",
			input: "echo hi ||",
		},
		{
			name:  "trailing pipe",
			input: "echo hi |",
		},
		{
			name:  "output redirect without target",
			input: "echo >",
		},
		{
			name:  "input redirect without target",
			input: "cat <",
		},
		{
			name:  "leading pipe",
			input: "| grep hi",
		},
		{
			name:  "leading logical operator",
			input: "&& echo hi",
		},
		{
			name:  "unsupported semicolon operator",
			input: "echo hi ; pwd",
		},
		{
			name:  "unsupported single ampersand operator",
			input: "echo hi & pwd",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p := New(tc.input)
			node, err := p.Parse()
			if !errors.Is(err, ErrUnexpectedToken) {
				t.Fatalf("Parse(%q) error = %v, want %v", tc.input, err, ErrUnexpectedToken)
			}
			if node != nil {
				t.Fatalf("Parse(%q) returned non-nil node on invalid input: %#v", tc.input, node)
			}
		})
	}
}
