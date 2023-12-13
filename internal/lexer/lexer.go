package lexer

import (
	"github.com/mantton/anthe/internal/token"
)

type Lexer struct {
	filename     string // name of the file being read
	input        []rune // an array of each rune in the file
	position     int    // current position in input ( points to current character )
	readPosition int    // current reading position in input (after current character)
	ch           rune   // current character under examination, `position` points to this char in the input i.e input[position] = ch

	line int // current line
	col  int // current column
}

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1
)

func New(input string, filename string) *Lexer {
	l := &Lexer{input: []rune(input), filename: filename}
	l.ch = ' '
	l.position = 0
	l.readPosition = 0
	l.line = 1
	l.col = 1

	l.next()
	if l.ch == bom {
		l.next() //ignore BOM at file beginning
	}

	return l
}

// indicates the lexer is at end of the file
func (l *Lexer) isAtEnd() bool {
	return l.readPosition >= len(l.input)
}

// indicates that the next character matches the specified char, consume
func (l *Lexer) matchAndConsume(c rune) bool {
	if l.isAtEnd() {
		return false
	}
	if c != l.peek() {
		return false
	}

	l.next()
	return true
}

// move pointers a single character
func (l *Lexer) next() {
	if l.readPosition >= len(l.input) {
		// reached EOF
		l.ch = eof
	} else {
		// Set the ch to the next position
		l.ch = l.input[l.readPosition]

		if l.ch == '\n' {
			// moved to new line, reset col to 1 and increment line
			l.col = 1
			l.line++
		} else {
			// moved col by 1 position
			l.col += 1
		}
	}

	// update the current position to the next position
	l.position = l.readPosition

	// increment the next pointer
	l.readPosition += 1

}

// looks ahead at the next char
func (l *Lexer) peek() rune {
	if l.readPosition >= len(l.input) {
		return eof
	} else {
		return l.input[l.readPosition]
	}
}

// reads the current token from the current character and moves cursor to the next character after it
func (l *Lexer) NextToken() token.Token {
	// move to next non whitespace character
	l.skipWhitespace()

	var tok token.Token
	// is EOF character
	if l.ch == eof {
		tok = token.Token{Literal: "EOF", Type: token.EOF}
	} else if token.IsSymbol(l.ch) {
		tok = l.nextSymbolToken()

	} else {
		tok = l.nextNonSymbolToken()

		if tok.Type != token.ILLEGAL && tok.Type != token.STRING {
			// already shifted cursor dont shift again
			return tok
		}

	}

	l.next()
	return tok
}

// returns a new token
func newRuneToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// returns a new token
func newStringToken(tokenType token.TokenType, str string) token.Token {
	return token.Token{Type: tokenType, Literal: str}
}

// consumes a single symbol token
func (l *Lexer) nextSymbolToken() token.Token {
	var tok token.Token
	switch l.ch {

	case '=':
		if l.matchAndConsume('=') {
			// '=='
			tok = newStringToken(token.EQL, "==")
		} else {
			tok = newRuneToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.matchAndConsume('=') {
			tok = newStringToken(token.NEQ, "!=")
		} else {
			tok = newRuneToken(token.NOT, l.ch)
		}

	// arithmetic
	case '+':
		tok = newRuneToken(token.ADD, l.ch)
	case '-':
		tok = newRuneToken(token.SUB, l.ch)
	case '*':
		tok = newRuneToken(token.MUL, l.ch)
	case '/':
		tok = newRuneToken(token.QUO, l.ch)

	// boolean
	case '>':
		tok = newRuneToken(token.GTR, l.ch)
	case '<':
		tok = newRuneToken(token.LSS, l.ch)

	// delim
	case ';':
		tok = newRuneToken(token.SEMICOLON, l.ch)
	case ',':
		tok = newRuneToken(token.COMMA, l.ch)
	case '(':
		tok = newRuneToken(token.LPAREN, l.ch)
	case ')':
		tok = newRuneToken(token.RPAREN, l.ch)
	case '{':
		tok = newRuneToken(token.LBRACE, l.ch)
	case '}':
		tok = newRuneToken(token.RBRACE, l.ch)
	case '[':
		tok = newRuneToken(token.LBRACKET, l.ch)
	case ']':
		tok = newRuneToken(token.RBRACKET, l.ch)

	case ':':
		tok = newRuneToken(token.COLON, l.ch)
	case '?':
		tok = newRuneToken(token.Q_MARK, l.ch)
	default:
		tok = newRuneToken(token.ILLEGAL, l.ch)
	}

	return tok
}
func (l *Lexer) nextNonSymbolToken() token.Token {
	switch {
	case l.ch == '"' || l.ch == '\'':
		tok := token.Token{Literal: l.readString(), Type: token.STRING}

		if l.ch == eof {
			panic("TODO: invalid string declaration")
		}
		return tok
	case isLetter(l.ch):
		ident := l.readIdentifier()
		return token.Token{Literal: ident, Type: token.LookupIdent(ident)}
	case isDigit(l.ch):

		val, isFloat := l.readNumber()

		if isFloat {
			return token.Token{Literal: val, Type: token.FLOAT}
		} else {
			return token.Token{Literal: val, Type: token.INTEGER}
		}

	}
	return token.Token{Literal: string(l.ch), Type: token.ILLEGAL}
}
