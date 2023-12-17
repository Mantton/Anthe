package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/mantton/anthe/internal/ast"
)

// compile AST Statement
func (c *Compiler) compileStatement(node ast.Statement, block *ir.Block, table *SymbolTable) {
	fmt.Printf("\n%T", node)
	switch node := node.(type) {
	case *ast.NamedFunctionDeclaration:
		c.compileNamedFunctionDeclaration(node, block, table)
	case *ast.LetStatement:
		c.compileLetStatement(node, block, table)
	}
}

func (c *Compiler) compileNamedFunctionDeclaration(node *ast.NamedFunctionDeclaration, block *ir.Block, table *SymbolTable) {
	isMain := node.Name == "main"
	// TODO: package check too.
	if isMain {
		fn := c.module.NewFunc(node.Name, types.I32)
		entryBlock := fn.NewBlock("entry")

		for _, s := range node.Fn.Body.Statements {
			c.compileStatement(s, entryBlock, table)
		}

		entryBlock.NewRet(constant.NewInt(types.I32, 0))

	} else {
		fn := c.module.NewFunc(node.Name, types.I32)
		entryBlock := fn.NewBlock("entry")

		for _, s := range node.Fn.Body.Statements {
			c.compileStatement(s, entryBlock, table)
		}

		entryBlock.NewRet(constant.NewInt(types.I32, 0))

	}

}

func (c *Compiler) compileLetStatement(node *ast.LetStatement, block *ir.Block, table *SymbolTable) {

	if block == nil {
		panic("nil block")
	}
	rhs := c.compileExpression(node.Value, table, block)

	// Allocate mem
	val := block.NewAlloca(rhs.Type())

	// Store
	block.NewStore(rhs, val)

	// add to symbol table
	table.Add(node.Name.Value, SymbolInfo{Name: node.Name.Value, Value: val, Type: rhs.Type()})
}
