package parser

type Node interface {
	node()
}

type Command struct {
	Name string
	Args []string
}

func (c *Command) node() {}

type Pipeline struct {
	Commands []*Command
}

func (p *Pipeline) node() {}

type Redirect struct {
	Type   RedirectType
	Target string
	Node   Node
}

func (r *Redirect) node() {}

type RedirectType int

const (
	RedirectInput RedirectType = iota
	RedirectOutput
)

type LogicalOp struct {
	Operator OpType
	Left     Node
	Right    Node
}

func (l *LogicalOp) node() {}

type OpType int

const (
	OpAnd OpType = iota
	OpOr
)
