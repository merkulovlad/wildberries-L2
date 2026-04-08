package parser

type TokenType int

const (
	TokenWord TokenType = iota
	TokenPipe
	TokenRedirectIn
	TokenRedirectOut
	TokenAnd
	TokenOr
	TokenAmpersand
	TokenSemicolon
	TokenEOF
)

type QuoteMode string

const (
	NonQuotes    QuoteMode = ""
	SingleQuotes QuoteMode = "single"
	DoubleQuotes QuoteMode = "double"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input     string
	pos       int
	quoteMode QuoteMode
}

// NewLexer creates a new lexer for tokenizing input strings
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:     input,
		pos:       0,
		quoteMode: NonQuotes,
	}
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	symbol := l.peek()
	n := 0

	defer func() {
		for i := 0; i < n; i++ {
			l.advance()
		}
	}()

	switch symbol {
	case 0:
		return Token{
			Type:  TokenEOF,
			Value: "",
		}

	case '|':
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '|' {
			n = 2
			return Token{
				Type:  TokenOr,
				Value: "||",
			}
		}
		n = 1
		return Token{
			Type:  TokenPipe,
			Value: "|",
		}

	case '&':
		if l.pos+1 < len(l.input) && l.input[l.pos+1] == '&' {
			n = 2
			return Token{
				Type:  TokenAnd,
				Value: "&&",
			}
		}
		n = 1
		return Token{
			Type:  TokenAmpersand,
			Value: "&",
		}

	case ';':
		n = 1
		return Token{
			Type:  TokenSemicolon,
			Value: ";",
		}

	case '<':
		n = 1
		return Token{
			Type:  TokenRedirectIn,
			Value: "<",
		}

	case '>':
		n = 1
		return Token{
			Type:  TokenRedirectOut,
			Value: ">",
		}

	default:
		start := l.pos
		for {
			ch := l.peek()
			if ch == 0 ||
				ch == ' ' ||
				ch == '\t' ||
				ch == '\n' ||
				ch == '\r' ||
				ch == '|' ||
				ch == '&' ||
				ch == ';' ||
				ch == '<' ||
				ch == '>' {
				break
			}
			l.advance()
		}
		return Token{
			Type:  TokenWord,
			Value: l.input[start:l.pos],
		}
	}
}

// peek returns the current character without advancing
func (l *Lexer) peek() byte {
	if l.pos < len(l.input) {
		return l.input[l.pos]
	}
	return 0
}

// advance moves to the next character in the input
func (l *Lexer) advance() {
	l.pos++
}

// skipWhitespace advances past all whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.quoteMode == NonQuotes {
		ch := l.peek()
		if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
			break
		}
		l.advance()
	}
}