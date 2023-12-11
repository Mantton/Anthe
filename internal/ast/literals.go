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

// conform
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }

func (b *BooleanLiteral) expressionNode()      {}
func (b *BooleanLiteral) TokenLiteral() string { return b.Token.Literal }
