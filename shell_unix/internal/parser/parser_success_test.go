package parser

import "testing"

func assertConsumedAllInput(t *testing.T, p *Parser, input string) {
	t.Helper()

	if p.current.Type != TokenEOF {
		t.Fatalf("Parse(%q) stopped before EOF: current=%v value=%q", input, p.current.Type, p.current.Value)
	}

	if p.lexer.pos != len(p.lexer.input) {
		t.Fatalf("Parse(%q) did not consume the full input: pos=%d len=%d", input, p.lexer.pos, len(p.lexer.input))
	}
}

func TestParseConsumesAllInputForValidCommands(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"echo hi",
		"echo hi | grep hi",
		"cat < input.txt",
		"echo hi > out.txt",
		"echo hi && pwd || whoami",
	}

	for _, input := range inputs {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			p := New(input)
			node, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse(%q) returned unexpected error: %v", input, err)
			}
			if node == nil {
				t.Fatalf("Parse(%q) returned nil node", input)
			}

			assertConsumedAllInput(t, p, input)
		})
	}
}

func TestParseEmptyInput(t *testing.T) {
	t.Parallel()

	inputs := []string{
		"",
		" \t\r\n ",
	}

	for _, input := range inputs {
		input := input
		t.Run(input, func(t *testing.T) {
			t.Parallel()

			p := New(input)
			node, err := p.Parse()
			if err != nil {
				t.Fatalf("Parse(%q) returned unexpected error: %v", input, err)
			}
			if node != nil {
				t.Fatalf("Parse(%q) returned non-nil node for empty input: %#v", input, node)
			}

			assertConsumedAllInput(t, p, input)
		})
	}
}
