package object

import "fmt"

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER      = "INTEGER"
	BOOLEAN      = "BOOLEAN"
	NULL         = "NULL"
	VOID         = "VOID"
	RETURN_VALUE = "RETURN_VALUE"
)

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type Null struct{}
type Void struct{}

type ReturnValue struct {
	Value Object
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

func (n *Void) Type() ObjectType { return VOID }
func (n *Void) Inspect() string  { return "void" }

func (n *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (n *ReturnValue) Inspect() string  { return n.Value.Inspect() }
