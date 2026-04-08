package parser

import "errors"

var ErrUnexpectedToken error = errors.New("Unexpected token")


type Parser struct {
	lexer   *Lexer
	current Token
}

// New creates a new parser for the given input string
func New(input string) *Parser {
	l := NewLexer(input)
	return &Parser{
		lexer: l,
		current: l.NextToken(),
	}
}

// Parse parses the input into an AST
func (p *Parser) Parse() (Node, error) {
	if p.current.Type == TokenEOF {
		return nil, nil
	}

	node, err := p.parseLogical()
	if err != nil {
		return nil, err
	}

	if p.current.Type != TokenEOF {
		return nil, ErrUnexpectedToken
	}

	return node, nil
}

// parseLogical parses logical operators (&& and ||)
func (p *Parser) parseLogical() (Node, error) {
	node, err := p.parsePipeline()
	if err != nil {
		return nil, err
	}
	
	for p.current.Type == TokenOr || p.current.Type == TokenAnd {
		typeLogic := p.current.Type == TokenOr // true if Or, false if and
		
		p.advance()
		rightNode, err := p.parsePipeline()
		
		if err != nil {
			return nil, err
		}
		
		if typeLogic {
			node = &LogicalOp{
				Operator: OpOr,
				Left: node,
				Right: rightNode,
			}
		} else {
			node = &LogicalOp{
				Operator: OpAnd,
				Left: node,
				Right: rightNode,
			}
		}
	}
	return node, nil
}

// parsePipeline parses pipeline commands separated by |
func (p *Parser) parsePipeline() (Node, error) {
	nodes := make([]Node, 0, 2)
	node, err := p.parseRedirect()
	if err != nil {
		return nil, err
	}
	nodes = append(nodes, node)
	for p.current.Type == TokenPipe {
		p.advance()
		node, err = p.parseRedirect()
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	
	if len(nodes) == 1 {
		return nodes[0], nil
	}
	
	return &Pipeline{Commands: nodes}, nil
}

// parseRedirect parses input/output redirection (< and >)
func (p *Parser) parseRedirect() (Node, error) {
	node, err := p.parseCommand()
	if err != nil {
		return nil, err
	}
	for p.current.Type == TokenRedirectIn ||  p.current.Type == TokenRedirectOut {
		current := p.current.Type == TokenRedirectIn // 1 if Redirect in, 0 RedirectOut
		p.advance()
		value := p.current.Value
		if err = p.expect(TokenWord); err != nil {
			return nil, err
		}
		if current {
			node = &Redirect{
				Type: RedirectInput,
				Target: value,
				Node: node,
			}
		} else {
			node = &Redirect{
				Type: RedirectOutput,
				Target: value,
				Node: node,
			}
		}
	}
	
	return node, nil
}

// parseCommand parses a single command with its arguments
func (p *Parser) parseCommand() (Node, error) {
	var words []string
	for p.current.Type == TokenWord {
		words = append(words, p.current.Value)
		p.advance()
	}
	if len(words) == 0 {
		return nil, ErrUnexpectedToken
	}
	return &Command{
		Name: words[0],
		Args: words[1:],
	}, nil
}

// expect checks if current token matches expected type and advances
func (p *Parser) expect(t TokenType) error {
	if t != p.current.Type {
		return ErrUnexpectedToken
	}
	p.advance()
	return nil
}

// advance moves to the next token
func (p *Parser) advance() {
	p.current = p.lexer.NextToken()
}

