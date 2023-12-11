package parser

import (
	"fmt"
	"strconv"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/token"
)

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 10, 64)

	if err != nil {
		return nil, err
	}

	lit.Value = value

	return lit, nil

}

func (p *Parser) parseBooleanLiteral() (ast.Expression, error) {
	return &ast.BooleanLiteral{Token: p.curToken, Value: p.currentMatches(token.TRUE)}, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.peekMatches(token.LPAREN) {
		return nil, fmt.Errorf("expected '(' found %s instead", p.peekToken.Literal)
	}

	params, err := p.parseFunctionParameters()

	if err != nil {
		return nil, err
	}

	lit.Parameters = params

	body, err := p.parseBlockStatement()

	if err != nil {
		return nil, err
	}

	lit.Body = body
	return lit, nil
}

func (p *Parser) parseFunctionParameters() ([]*ast.IdentifierExpression, error) {
	identifiers := []*ast.IdentifierExpression{}

	if p.peekMatches(token.RPAREN) {
		p.next()
		return identifiers, nil
	}

	p.next()

	ident := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekMatches(token.COMMA) {
		p.next() // move to comma
		p.next() // move to token after comma

		ident := &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.consumeIfPeekMatches(token.RPAREN) {
		return nil, fmt.Errorf("expected ')' after parameter list found %s", p.peekToken.Literal)
	}

	return identifiers, nil
}

func (p *Parser) parseArrayLiteral() (ast.Expression, error) {
	array := &ast.ArrayLiteral{Token: p.curToken}

	elems, err := p.parseExpressionList(token.RBRACKET)

	if err != nil {
		return nil, err
	}

	array.Elements = elems
	return array, nil
}

func (p *Parser) parseExpressionList(end token.TokenType) ([]ast.Expression, error) {

	list := []ast.Expression{}

	if p.peekMatches(end) {
		p.next()
		return list, nil
	}

	p.next()

	expr, err := p.parseExpression(LOWEST)

	if err != nil {
		return nil, err
	}

	list = append(list, expr)

	for p.peekMatches(token.COMMA) {
		p.next()
		p.next()
		expr, err := p.parseExpression(LOWEST)

		if err != nil {
			return nil, err
		}
		list = append(list, expr)

	}

	if !p.consumeIfPeekMatches(end) {
		return nil, fmt.Errorf("expected '%d' at end of expression list", end) // TODO: lookup token
	}

	return list, nil
}
