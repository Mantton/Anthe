package parser

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/token"
)

/*
KEYWORD | IDENTIFIER | EXPRESSION LIST | BOCK_STATEMENT
func myFunc() {}
*/
func (p *Parser) parseFunctionDeclaration() (*ast.NamedFunctionDeclaration, error) {
	expr := &ast.NamedFunctionDeclaration{}
	// on the func keyword

	if !p.consumeIfPeekMatches(token.IDENTIFIER) {
		return nil, fmt.Errorf("expected function name got %s instead", p.peekToken.Literal)
	}

	// currently on identifier, parse
	ident, err := p.parseIdentifier()
	expr.Name = ident.(*ast.IdentifierExpression).Value

	if err != nil {
		return nil, err
	}

	// at this point, it becomes a regular FunctionLiteral, parse that

	fn, err := p.parseFunctionLiteral()

	if err != nil {
		return nil, err
	}

	switch fn := fn.(type) {
	case *ast.FunctionLiteral:
		expr.Fn = fn
	default:
		return nil, fmt.Errorf("expected function literal got %s instead", p.curToken.Literal)
	}

	if p.peekMatches(token.SEMICOLON) {
		p.next()
	}

	return expr, nil

}
