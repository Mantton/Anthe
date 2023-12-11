package ast

import "github.com/mantton/anthe/internal/token"

type Expression interface {
	Node
	expressionNode()
}

type IdentifierExpression struct {
	Token token.Token
	Value string
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Action      *BlockStatement
	Alternative *BlockStatement
}

// conform
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (i *IdentifierExpression) expressionNode()      {}
func (i *IdentifierExpression) TokenLiteral() string { return i.Token.Literal }

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return i.Token.Literal }

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
