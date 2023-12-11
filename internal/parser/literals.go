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

func (p *Parser) parseFloatingPointLiteral() (ast.Expression, error) {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)

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

	if !p.consumeIfPeekMatches(token.LPAREN) {
		return nil, fmt.Errorf("expected '(' found %s instead", p.peekToken.Literal)
	}

	params, err := p.parseFunctionParameters()

	if err != nil {
		return nil, err
	}

	lit.Parameters = params

	if !p.consumeIfPeekMatches(token.LBRACE) {
		return nil, fmt.Errorf("expected function body found %s instead", p.peekToken.Literal)
	}

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

	elems, err := p.parseExpressionList(token.RBRACKET, ']')

	if err != nil {
		return nil, err
	}

	array.Elements = elems
	return array, nil
}

func (p *Parser) parseExpressionList(end token.TokenType, c rune) ([]ast.Expression, error) {

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
		return nil, fmt.Errorf("expected '%c' at end of expression list got %s", c, p.peekToken.Literal) // TODO: lookup token
	}

	return list, nil
}

func (p *Parser) parseHashLiteral() (ast.Expression, error) {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)

	for !p.peekMatches(token.RBRACE) {
		p.next()
		key, err := p.parseExpression(LOWEST)

		if err != nil {
			return nil, err
		}

		if !p.consumeIfPeekMatches(token.COLON) {
			return nil, fmt.Errorf("expected ':' after key found %s", p.peekToken.Literal)
		}

		p.next()

		value, err := p.parseExpression(LOWEST)

		if err != nil {
			return nil, err
		}

		hash.Pairs[key] = value

		if !p.peekMatches(token.RBRACE) && !p.consumeIfPeekMatches(token.COMMA) {
			return nil, fmt.Errorf("invalid object expression")
		}
	}

	if !p.consumeIfPeekMatches(token.RBRACE) {
		return nil, fmt.Errorf("expected '}' found %s", p.peekToken.Literal)
	}

	return hash, nil
}

func (p *Parser) parseStringLiteral() (ast.Expression, error) {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}, nil
}
