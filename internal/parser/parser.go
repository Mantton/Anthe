package parser

import (
	"github.com/mantton/anthe/internal/ast"
	"github.com/mantton/anthe/internal/lexer"
	"github.com/mantton/anthe/internal/token"
)

type ExpPrecedence = byte

const (
	_ ExpPrecedence = iota
	LOWEST
	ASSIGN      // x = y
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         //+
	PRODUCT     //*
	PREFIX      //-Xor!X
	CALL        // myFunction(X)
	INDEX       // array[index]
)

type (
	prefixParseFn func() (ast.Expression, error)               // --5
	infixParseFn  func(ast.Expression) (ast.Expression, error) // 5 * 5
)

var precedences = map[token.TokenType]ExpPrecedence{
	token.ASSIGN:   ASSIGN,
	token.EQL:      EQUALS,
	token.NEQ:      EQUALS,
	token.LSS:      LESSGREATER,
	token.GTR:      LESSGREATER,
	token.GEQ:      LESSGREATER,
	token.LEQ:      LESSGREATER,
	token.ADD:      SUM,
	token.SUB:      SUM,
	token.QUO:      PRODUCT,
	token.MUL:      PRODUCT,
	token.LPAREN:   CALL,
	token.LBRACKET: INDEX,
}

type Parser struct {
	l *lexer.Lexer

	curToken  token.Token // the current token
	peekToken token.Token // the next token after the current token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// read two tokens, setting both current and peekToken
	p.next() // sets peek
	p.next() // sets current

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)

	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatingPointLiteral)
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)

	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)

	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.NULL, p.parseNullLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.ADD, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.QUO, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.EQL, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.LSS, p.parseInfixExpression)
	p.registerInfix(token.GTR, p.parseInfixExpression)
	p.registerInfix(token.GEQ, p.parseInfixExpression)
	p.registerInfix(token.LEQ, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignmentExpression)

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
			break
		}

		// if statement is valid, append

		program.Statements = append(program.Statements, stmt)

		p.next()
	}
	return program
}

func (p *Parser) registerPrefix(tok token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tok] = fn
}

func (p *Parser) registerInfix(tok token.TokenType, fn infixParseFn) {
	p.infixParseFns[tok] = fn
}
