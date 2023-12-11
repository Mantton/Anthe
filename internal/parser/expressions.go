package parser

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/token"
)

func (p *Parser) parseExpression(prec ExpPrecedence) (ast.Expression, error) {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		return nil, fmt.Errorf("no prefix parse method for %s found", p.curToken.Literal)
	}

	lhs, err := prefix()

	if err != nil {
		return nil, err
	}

	for !p.peekMatches(token.SEMICOLON) && prec < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return lhs, nil
		}

		p.next()

		lhs, err = infix(lhs)

		if err != nil {
			return nil, err
		}
	}
	return lhs, nil
}

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.IdentifierExpression{Token: p.curToken, Value: p.curToken.Literal}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	expr := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.next()

	rhs, err := p.parseExpression(PREFIX)

	if err != nil {
		return nil, err
	}

	expr.Right = rhs

	return expr, nil
}

func (p *Parser) parseInfixExpression(lhs ast.Expression) (ast.Expression, error) {
	expr := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     lhs,
	}

	presc := p.currentPrecedence()
	p.next()

	rhs, err := p.parseExpression(presc)

	if err != nil {
		return nil, err
	}
	expr.Right = rhs
	return expr, nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.next()

	exp, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	if !p.consumeIfPeekMatches(token.RPAREN) {
		return nil, fmt.Errorf("expected ')' got %s instead", p.peekToken.Literal)
	}

	return exp, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	expr := &ast.IfExpression{Token: p.curToken}

	hasLParen := p.consumeIfPeekMatches(token.LPAREN)
	// on either if token or bracket token, move once more to expression
	p.next()

	condition, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}

	expr.Condition = condition

	if hasLParen {
		if !p.consumeIfPeekMatches(token.RPAREN) {
			return nil, fmt.Errorf("expected ')', got %s", p.peekToken.Literal)
		}
	}

	// now either at '{' or at last expression
	if !p.consumeIfPeekMatches(token.LBRACE) {
		return nil, fmt.Errorf("expected '{', got %s", p.peekToken.Literal)
	}

	action, err := p.parseBlockStatement()

	if err != nil {
		return nil, err
	}

	expr.Action = action

	// Next token is else statement
	if p.peekMatches(token.ELSE) {
		p.next()

		if !p.consumeIfPeekMatches(token.LBRACE) {
			return nil, fmt.Errorf("expected '{' got %s instead", p.peekToken.Literal)
		}

		alt, err := p.parseBlockStatement()

		if err != nil {
			return nil, err
		}

		expr.Alternative = alt
	}

	return expr, nil

}

func (p *Parser) parseCallExpression(function ast.Expression) (ast.Expression, error) {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}

	args, err := p.parseExpressionList(token.RPAREN, ')')
	if err != nil {
		return nil, err
	}

	exp.Arguments = args

	return exp, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.next() // move away from '['
	idx, err := p.parseExpression(LOWEST)

	if err != nil {
		return nil, err
	}

	exp.Index = idx

	if !p.consumeIfPeekMatches(token.RBRACKET) {
		return nil, fmt.Errorf("expected ']' after index got %s", p.peekToken.Literal)
	}

	return exp, nil
}
