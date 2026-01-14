package parser

type TokenType int

const (
	TokenWord TokenType = iota
	TokenPipe
	TokenRedirectIn
	TokenRedirectOut
	TokenAnd
	TokenOr
	TokenEOF
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input string
	pos   int
}

// NewLexer creates a new lexer for tokenizing input strings
func NewLexer(input string) *Lexer {
	return &Lexer{}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	return Token{}
}

// peek returns the current character without advancing
func (l *Lexer) peek() byte {
	return 0
}

// advance moves to the next character in the input
func (l *Lexer) advance() {
}

// skipWhitespace advances past all whitespace characters
func (l *Lexer) skipWhitespace() {
}
