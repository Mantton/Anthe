package ast

import "github.com/mantton/anthe/internal/token"

type Statement interface {
	Node
	statementNode()
}

// LET
type LetStatement struct {
	Token token.Token // The token.LET token
	Name  *IdentifierExpression
	Value Expression
}

// RETURN
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// EXPRESSION
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

// conform
func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return "Let " + s.Token.Literal }

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return "Return  " + s.Token.Literal }

func (s *ExpressionStatement) statementNode() {}
func (s *ExpressionStatement) TokenLiteral() string {
	return "Expression " + s.Token.Literal + s.Expression.TokenLiteral()
}

func (s *BlockStatement) statementNode()       {}
func (s *BlockStatement) TokenLiteral() string { return s.Token.Literal }
