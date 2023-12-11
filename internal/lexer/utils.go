package lexer

import "fmt"

// checks if the given character is a letter
func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// checks if the given character is a digit
func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

// moves cursor to last letter character, returns the resulting substring
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.next()
	}
	return string(l.input[position:l.position])
}

// moves pointer till the point where the current character is not a new line
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.next()
	}
}

// moves cursor to last digit character, returns the resulting substring
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.next()
	}
	return string(l.input[position:l.position])
}

func (l *Lexer) debug() {
	fmt.Printf("\nCurrently at : %s", string(l.ch))
}
