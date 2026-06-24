package env

import (
	"os"

	"github.com/merkulovlad/wildberries-L2/shell_unix/internal/parser"
)

type Expander struct {
}

// NewExpander creates a new environment variable expander
func NewExpander() *Expander {
	return &Expander{}
}

// Expand replaces environment variables in the input string
func (e *Expander) Expand(input string) string {
	return os.ExpandEnv(input)
}

// ExpandArgs expands environment variables in all arguments
func (e *Expander) ExpandArgs(args []string) []string {
	if args == nil {
		return nil
	}

	ans := make([]string, 0, len(args))
	for _, arg := range args {
		ans = append(ans, e.Expand(arg))
	}

	return ans
}

// ExpandNode expands environment variables in string fields across the AST.
func (e *Expander) ExpandNode(node parser.Node) parser.Node {
	switch n := node.(type) {
	case nil:
		return nil
	case *parser.Command:
		return &parser.Command{
			Name: e.Expand(n.Name),
			Args: e.ExpandArgs(n.Args),
		}
	case *parser.Pipeline:
		nodes := make([]parser.Node, 0, len(n.Commands))
		for _, command := range n.Commands {
			nodes = append(nodes, e.ExpandNode(command))
		}
		return &parser.Pipeline{Commands: nodes}
	case *parser.Redirect:
		return &parser.Redirect{
			Type:   n.Type,
			Target: e.Expand(n.Target),
			Node:   e.ExpandNode(n.Node),
		}
	case *parser.LogicalOp:
		return &parser.LogicalOp{
			Operator: n.Operator,
			Left:     e.ExpandNode(n.Left),
			Right:    e.ExpandNode(n.Right),
		}
	default:
		return node
	}
}
