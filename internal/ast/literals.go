package ast

import (
	"github.com/mantton/anthe/internal/token"
)

type LiteralExpression interface {
	Expression
	literalNode()
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

type FloatLiteral struct {
	Token token.Token
	Value float64
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

type NullLiteral struct {
	Token token.Token
}

// conform
func (il *IntegerLiteral) expressionNode()      {}
func (n *IntegerLiteral) literalNode()          {}
func (il *IntegerLiteral) TokenLiteral() string { return "IntLit " + il.Token.Literal }

func (b *BooleanLiteral) expressionNode()      {}
func (n *BooleanLiteral) literalNode()         {}
func (b *BooleanLiteral) TokenLiteral() string { return "BoolLit " + b.Token.Literal }

func (b *FunctionLiteral) expressionNode()      {}
func (n *FunctionLiteral) literalNode()         {}
func (b *FunctionLiteral) TokenLiteral() string { return "FuncLit " + b.Token.Literal }

func (b *ArrayLiteral) expressionNode()      {}
func (n *ArrayLiteral) literalNode()         {}
func (b *ArrayLiteral) TokenLiteral() string { return "ArrLit " + b.Token.Literal }

func (b *HashLiteral) expressionNode()      {}
func (n *HashLiteral) literalNode()         {}
func (b *HashLiteral) TokenLiteral() string { return "HashLit " + b.Token.Literal }

func (b *StringLiteral) expressionNode()      {}
func (n *StringLiteral) literalNode()         {}
func (b *StringLiteral) TokenLiteral() string { return "StringLit " + b.Token.Literal }

func (b *FloatLiteral) expressionNode()      {}
func (n *FloatLiteral) literalNode()         {}
func (b *FloatLiteral) TokenLiteral() string { return "Float Lit " + b.Token.Literal }

func (b *NullLiteral) expressionNode()      {}
func (n *NullLiteral) literalNode()         {}
func (b *NullLiteral) TokenLiteral() string { return "null Lit " + b.Token.Literal }
