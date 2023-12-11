package parser

import "github.com/mantton/anthe/internal/token"

// bool indicating the current token is of the specified type
func (p *Parser) currentMatches(t token.TokenType) bool {
	return p.curToken.Type == t
}

// bool indicating the peek/next token is of the specified type
func (p *Parser) peekMatches(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// consumes a token if the peek matches a specified token, returns true if the peek matches
func (p *Parser) consumeIfPeekMatches(t token.TokenType) bool {
	if p.peekMatches(t) {
		p.next()
		return true
	} else {
		return false
	}
}
