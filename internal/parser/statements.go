package parser

import (
	"fmt"

	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/token"
)

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.curToken.Type {
	case token.FUNCTION:
		return p.parseFunctionDeclaration()
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

	if p.consumeIfPeekMatches(token.COLON) {
		// strong typing variable
		p.next() // move to type decl, parse type

		t, err := p.parseTypeDeclaration()

		if err != nil {
			return nil, err
		}
		stmt.Type = t
	}

	// the next statement must be an assignment token to be a valid token, return nil if not
	if !p.consumeIfPeekMatches(token.ASSIGN) {
		return nil, fmt.Errorf("expected variable assignment ('=') found %s instead", p.peekToken.Literal)
	}

	p.next()

	if p.curToken.Type == token.VOID {
		return nil, fmt.Errorf("variable cannot be assigned to void, use null instead") // TODO: move to type checker
	}

	v, err := p.parseExpression(LOWEST)

	if err != nil {
		return nil, err
	}

	stmt.Value = v

	for p.peekMatches(token.SEMICOLON) {
		p.next()
	}

	return stmt, nil

}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// Move to next token
	p.next()

	// TODO: return;
	// if p.peekMatches(token.SEMICOLON) {
	// 	// return;
	// 	p.next() // move to semi
	// 	p.next()
	// }

	v, err := p.parseExpression(LOWEST)

	if err != nil {
		return nil, err
	}

	stmt.ReturnValue = v
	for p.peekMatches(token.SEMICOLON) {
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
	// fmt.Printf("\n%T", stmt.Expression)

	if p.peekMatches(token.SEMICOLON) {
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

func (p *Parser) parseTypeDeclaration() (ast.TypeExpression, error) {

	if p.curToken.Type != token.IDENTIFIER {
		return nil, fmt.Errorf("unknown type identifier `%s`", p.curToken.Literal)
	}

	tok := token.LookUpBuiltInType(p.curToken.Literal)
	var gen []ast.TypeExpression
	var err error
	if tok == token.VOID || tok == token.NULL {
		return nil, fmt.Errorf("cannot declare variable as `%s`", p.curToken.Literal)
	}

	name := p.curToken.Literal
	if p.consumeIfPeekMatches(token.LSS) {

		gen, err = p.parseTypeGenericList(token.GTR, '>')
		if err != nil {
			return nil, err
		}
	}

	var t ast.TypeExpression

	// non generic types

	switch tok {
	// Literal
	case token.INT_T:
		t = &ast.LiteralIntegerType{}
	case token.STR_T:
		t = &ast.LiteralStringType{}
	case token.FLT_T:
		t = &ast.LiteralFloatType{}
	case token.BOOL_T:
		t = &ast.LiteralBooleanType{}

	// generics
	case token.OPTIONAL_T:
		if gen == nil || len(gen) != 1 {
			return nil, fmt.Errorf("generic type `%s` requires parameter definition: `%s`<T>", name, name)
		}
		t = &ast.OptionalType{Value: gen[0]}
	default:
		t = &ast.ScopeDefinedType{Values: gen, Name: name}
	}

	// is marked as optional, consume '?'
	if p.consumeIfPeekMatches(token.Q_MARK) {
		t = &ast.OptionalType{Value: t}
	}

	return t, nil
}

func (p *Parser) parseTypeGenericList(end token.TokenType, c rune) ([]ast.TypeExpression, error) {

	list := []ast.TypeExpression{}

	if p.peekMatches(end) {
		p.next()
		return list, nil
	}

	p.next()

	expr, err := p.parseTypeDeclaration()

	if err != nil {
		return nil, err
	}

	list = append(list, expr)

	for p.peekMatches(token.COMMA) {
		p.next()
		p.next()
		expr, err := p.parseTypeDeclaration()

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
