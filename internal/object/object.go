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
	ARRAY        = "ARRAY"
	HASH         = "HASH"
)

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type HashPair struct {
	Key   Object
	Value Object
}

type HashableProtocol interface {
	HashKey() HashKey
}

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

type Array struct {
	Elements []Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	} else {
		return HashKey{Type: b.Type(), Value: 0}
	}
}

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (n *Null) Type() ObjectType { return NULL }
func (n *Null) Inspect() string  { return "null" }

func (n *Void) Type() ObjectType { return VOID }
func (n *Void) Inspect() string  { return "void" }

func (n *ReturnValue) Type() ObjectType { return RETURN_VALUE }
func (n *ReturnValue) Inspect() string  { return n.Value.Inspect() }

func (n *Array) Type() ObjectType { return ARRAY }
func (n *Array) Inspect() string  { return "ARRAYYYYYY" }

func (n *Hash) Type() ObjectType { return HASH }
func (n *Hash) Inspect() string  { return "HASH" }
