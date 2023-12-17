package compiler

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type SymbolInfo struct {
	Name  string
	Value value.Value
	Type  types.Type
}

type SymbolTable struct {
	symbols map[string]SymbolInfo
	parent  *SymbolTable
}

func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]SymbolInfo),
		parent:  parent,
	}
}

func (s *SymbolTable) Add(name string, info SymbolInfo) {
	s.symbols[name] = info
}

func (s *SymbolTable) Lookup(name string) (SymbolInfo, bool) {
	info, exists := s.symbols[name]
	if !exists && s.parent != nil {
		return s.parent.Lookup(name)
	}
	return info, exists
}
