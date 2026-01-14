package parser

type Parser struct {
	lexer   *Lexer
	current Token
}

// New creates a new parser for the given input string
func New(input string) *Parser {
	return &Parser{}
}

// Parse parses the input into an AST
func (p *Parser) Parse() (Node, error) {
	return nil, nil
}

// parseLogical parses logical operators (&& and ||)
func (p *Parser) parseLogical() (Node, error) {
	return nil, nil
}

// parsePipeline parses pipeline commands separated by |
func (p *Parser) parsePipeline() (Node, error) {
	return nil, nil
}

// parseCommand parses a single command with its arguments
func (p *Parser) parseCommand() (Node, error) {
	return nil, nil
}

// parseRedirect parses input/output redirection (< and >)
func (p *Parser) parseRedirect() (Node, error) {
	return nil, nil
}

// advance moves to the next token
func (p *Parser) advance() {
}

// expect checks if current token matches expected type and advances
func (p *Parser) expect(t TokenType) error {
	return nil
}
