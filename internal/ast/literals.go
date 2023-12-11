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

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

type StringLiteral struct {
	Token token.Token
	Value string
}

// conform
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return "IntLit " + il.Token.Literal }

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return "BoolLit " + b.Token.Literal }

func (b *FunctionLiteral) expressionNode()      {}
func (b *FunctionLiteral) TokenLiteral() string { return "FuncLit " + b.Token.Literal }

func (b *ArrayLiteral) expressionNode()      {}
func (b *ArrayLiteral) TokenLiteral() string { return "ArrLit " + b.Token.Literal }

func (b *HashLiteral) expressionNode()      {}
func (b *HashLiteral) TokenLiteral() string { return "HashLit " + b.Token.Literal }

func (b *StringLiteral) expressionNode()      {}
func (b *StringLiteral) TokenLiteral() string { return "StringLit " + b.Token.Literal }
