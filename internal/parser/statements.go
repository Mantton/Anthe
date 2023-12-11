package parser

import (
	"errors"
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/token"
)

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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
	stmt.Name = &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}

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

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	expr, err := p.parseExpression(LOWEST)

	if err != nil {
		return nil, err
	}

	stmt.Expression = expr
	if p.consumeIfPeekMatches(token.SEMICOLON) {
		p.next()
	}
	return stmt, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{Token: p.curToken}

	block.Statements = []ast.Statement{}

	p.next()

	for !p.currentMatches(token.RBRACE) && !p.currentMatches(token.EOF) {
		stmt, err := p.parseStatement()

		if err != nil {
			return nil, err
		}

		block.Statements = append(block.Statements, stmt)
		p.next()
	}

	return block, nil
}
