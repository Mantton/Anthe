package ast

import "github.com/mantton/anthe/internal/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*IdentifierExpression
	Body       *BlockStatement
}

// conform
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return "IntLit " + il.Token.Literal }

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return "BoolLit " + b.Token.Literal }

func (b *FunctionLiteral) expressionNode()      {}
func (b *FunctionLiteral) TokenLiteral() string { return "FuncLit " + b.Token.Literal }
