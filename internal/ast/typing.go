package ast

import "fmt"

type TypeExpression interface {
	typeNode()
	Type() string
}

// Primitives
type LiteralIntegerType struct{}
type LiteralStringType struct{}
type LiteralFloatType struct{}
type LiteralBooleanType struct{}

func (t *LiteralIntegerType) typeNode() {}
func (t *LiteralIntegerType) Type() string {
	return "int"
}

func (t *LiteralStringType) typeNode() {}
func (t *LiteralStringType) Type() string {
	return "string"
}

func (t *LiteralFloatType) typeNode() {}
func (t *LiteralFloatType) Type() string {
	return "float"
}

func (t *LiteralBooleanType) typeNode() {}
func (t *LiteralBooleanType) Type() string {
	return "boolean"
}

// More Advanced
type OptionalType struct {
	Value TypeExpression
}

type ScopeDefinedType struct {
	Name   string
	Values []TypeExpression
}

func (t *OptionalType) typeNode() {}
func (t *OptionalType) Type() string {
	return fmt.Sprintf("optional<%s>", t.Value.Type())
}

func (t *ScopeDefinedType) typeNode() {}
func (t *ScopeDefinedType) Type() string {
	if t.Values == nil {
		return t.Name
	} else {
		val := t.Name + "<"

		for i, exp := range t.Values {
			val += exp.Type()

			if i != len(t.Values)-1 {
				val += ", "
			}
		}

		val += ">"
		return val
	}
}
