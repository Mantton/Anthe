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

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

type AssignmentExpression struct {
	Token  token.Token
	Target *IdentifierExpression
	Value  Expression
}

// conform
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return "Prefix " + pe.Token.Literal }

func (i *IdentifierExpression) expressionNode()      {}
func (i *IdentifierExpression) TokenLiteral() string { return "Ident " + i.Token.Literal }

func (i *InfixExpression) expressionNode()      {}
func (i *InfixExpression) TokenLiteral() string { return "Infix " + i.Token.Literal }

func (i *IfExpression) expressionNode()      {}
func (i *IfExpression) TokenLiteral() string { return "If " + i.Token.Literal }

func (i *CallExpression) expressionNode()      {}
func (i *CallExpression) TokenLiteral() string { return "Call " + i.Token.Literal }

func (i *IndexExpression) expressionNode()      {}
func (i *IndexExpression) TokenLiteral() string { return "Idx " + i.Token.Literal }

func (i *AssignmentExpression) expressionNode()      {}
func (i *AssignmentExpression) TokenLiteral() string { return "assign " + i.Token.Literal }
