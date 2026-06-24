package env

import (
	"reflect"
	"testing"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

func TestExpandArgsPreservesNil(t *testing.T) {
	t.Parallel()

	expander := NewExpander()

	if got := expander.ExpandArgs(nil); got != nil {
		t.Fatalf("ExpandArgs(nil) = %#v, want nil", got)
	}
}

func TestExpandNodeCommand(t *testing.T) {
	t.Setenv("EXPAND_NODE_CMD", "echo")
	t.Setenv("EXPAND_NODE_ARG", "hello")

	expander := NewExpander()
	node := &parser.Command{
		Name: "$EXPAND_NODE_CMD",
		Args: []string{"$EXPAND_NODE_ARG", "world"},
	}

	got := expander.ExpandNode(node)
	want := &parser.Command{
		Name: "echo",
		Args: []string{"hello", "world"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpandNode(command) = %#v, want %#v", got, want)
	}

	if !reflect.DeepEqual(node, &parser.Command{
		Name: "$EXPAND_NODE_CMD",
		Args: []string{"$EXPAND_NODE_ARG", "world"},
	}) {
		t.Fatalf("ExpandNode mutated the original command: %#v", node)
	}
}

func TestExpandNodeRedirect(t *testing.T) {
	t.Setenv("EXPAND_NODE_FILE", "out.txt")
	t.Setenv("EXPAND_NODE_TEXT", "hi")

	expander := NewExpander()
	node := &parser.Redirect{
		Type:   parser.RedirectOutput,
		Target: "$EXPAND_NODE_FILE",
		Node: &parser.Command{
			Name: "echo",
			Args: []string{"$EXPAND_NODE_TEXT"},
		},
	}

	got := expander.ExpandNode(node)
	want := &parser.Redirect{
		Type:   parser.RedirectOutput,
		Target: "out.txt",
		Node: &parser.Command{
			Name: "echo",
			Args: []string{"hi"},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpandNode(redirect) = %#v, want %#v", got, want)
	}
}

func TestExpandNodePipelineAndLogical(t *testing.T) {
	t.Setenv("EXPAND_LEFT", "hello")
	t.Setenv("EXPAND_RIGHT", "world")
	t.Setenv("EXPAND_CMD", "cat")
	t.Setenv("EXPAND_INPUT", "input.txt")

	expander := NewExpander()
	node := &parser.LogicalOp{
		Operator: parser.OpAnd,
		Left: &parser.Pipeline{
			Commands: []parser.Node{
				&parser.Command{
					Name: "echo",
					Args: []string{"$EXPAND_LEFT"},
				},
				&parser.Command{
					Name: "echo",
					Args: []string{"$EXPAND_RIGHT"},
				},
			},
		},
		Right: &parser.Redirect{
			Type:   parser.RedirectInput,
			Target: "$EXPAND_INPUT",
			Node: &parser.Command{
				Name: "$EXPAND_CMD",
			},
		},
	}

	got := expander.ExpandNode(node)
	want := &parser.LogicalOp{
		Operator: parser.OpAnd,
		Left: &parser.Pipeline{
			Commands: []parser.Node{
				&parser.Command{
					Name: "echo",
					Args: []string{"hello"},
				},
				&parser.Command{
					Name: "echo",
					Args: []string{"world"},
				},
			},
		},
		Right: &parser.Redirect{
			Type:   parser.RedirectInput,
			Target: "input.txt",
			Node: &parser.Command{
				Name: "cat",
				Args: nil,
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpandNode(logical) = %#v, want %#v", got, want)
	}
}
