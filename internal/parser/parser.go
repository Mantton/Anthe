package parser

import (
	"errors"
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/lexer"
	"github.com/mantton/anthe/internal/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // the current token
	peekToken token.Token // the next token after the current token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// read two tokens, setting both current and peekToken
	p.next() // sets peek
	p.next() // sets current

	return p
}

func (p *Parser) next() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{} // init
	program.Statements = []ast.Statement{}
	program.Errors = []string{}

	for p.curToken.Type != token.EOF {
		// while not at eof token, parse statement
		stmt, err := p.parseStatement()

		if err != nil {
			program.Errors = append(program.Errors, err.Error())
		}

		// if statement is valid, append
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.next()
	}
	return program
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()

	default:
		return nil, errors.New("unknown statement declaration")
	}
}

/*
Parse let statements
called when the parser's current token is a token.LET
let statements take the following form

`identifier` `token.ASSIGN` `expression | literal | identifier`
*/
func (p *Parser) parseLetStatement() (*ast.LetStatement, error) {
	stmt := &ast.LetStatement{Token: p.curToken}
	// if the next token is not an identifier, it is not a valid statement
	if !p.consumeIfPeekMatches(token.IDENTIFIER) {
		return nil, fmt.Errorf("syntax error: expected an `identifier` got %s instead", p.peekToken.Literal)
	}

	// consumed, so now the current token matches the peek which we checked was an identifier
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// the next statement must be an assignment token to be a valid token, return nil if not
	if !p.consumeIfPeekMatches(token.ASSIGN) {
		return nil, errors.New("variables must be assigned immediately")
	}

	// TODO: Expressions till we match semi colon
	for !p.currentMatches(token.SEMICOLON) {
		p.next()
	}

	return stmt, nil

}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// Move to next token
	p.next()

	// TODO: expressions
	for !p.currentMatches(token.SEMICOLON) {
		p.next()
	}
	return stmt, nil
}
