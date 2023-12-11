package ast

type Node interface {
	TokenLiteral() string
}

type Program struct {
	Statements []Statement
	Errors     []string
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
