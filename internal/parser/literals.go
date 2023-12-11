package parser

import (
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
