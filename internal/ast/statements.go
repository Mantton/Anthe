package ast

import "github.com/mantton/anthe/internal/token"

type Statement interface {
	Node
	statementNode()
}

// LET
type LetStatement struct {
	Token token.Token // The token.LET token
	Name  *Identifier
	Value Expression
}

// RETURN
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

// conform
func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
