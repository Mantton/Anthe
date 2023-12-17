package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/mantton/anthe/internal/ast"
)

const PREFIX = "_an__"

// compile AST Statement
func (c *Compiler) compileStatement(node ast.Statement, block *ir.Block, table *SymbolTable) {
	fmt.Printf("\n%T", node)
	switch node := node.(type) {
	case *ast.NamedFunctionDeclaration:
		c.compileNamedFunctionDeclaration(node, block, table)
	case *ast.LetStatement:
		c.compileLetStatement(node, table)
	case *ast.ExpressionStatement:
		c.compileExpression(node.Expression, table)
	case *ast.BlockStatement:
		c.compileBlockStatement(node, block, table)
	case *ast.ReturnStatement:
		c.compileReturnStatement(node, table)
	default:
		panic("\nstatement compilation not implemented")
	}
}

func (c *Compiler) compileNamedFunctionDeclaration(node *ast.NamedFunctionDeclaration, block *ir.Block, table *SymbolTable) {
	isMain := node.Name == "main"
	// TODO: package check too.
	if isMain {
		fn := c.module.NewFunc("main", types.I64)
		entryBlock := fn.NewBlock("entry")
		c.currentBlock = entryBlock

		for _, s := range node.Fn.Body.Statements {
			c.compileStatement(s, entryBlock, table)
		}

		if c.currentBlock == entryBlock {
			entryBlock.NewRet(constant.NewInt(types.I64, 0))
		}

	} else {

		name := PREFIX + node.Name
		fnTable := NewSymbolTable(table)
		fnParams := []*ir.Param{}

		if node.Fn.Parameters != nil {
			for _, param := range node.Fn.Parameters {
				// allocate new param
				p := ir.NewParam(param.Value, types.I64)
				fnTable.Add(param.Value, SymbolInfo{Name: param.Value, Value: p, Type: p.Type(), IsParameter: true})
				fnParams = append(fnParams, p)
			}
		}

		fn := c.module.NewFunc(name, types.I64, fnParams...)
		table.Add(node.Name, SymbolInfo{Name: node.Name, Value: fn, Type: fn.Type()})

		entryBlock := fn.NewBlock("entry")
		c.currentBlock = entryBlock

		for _, s := range node.Fn.Body.Statements {
			c.compileStatement(s, entryBlock, fnTable)
		}

		if c.currentBlock == entryBlock {
			entryBlock.NewRet(constant.NewInt(types.I64, 0))
		}
	}

}

func (c *Compiler) compileLetStatement(node *ast.LetStatement, table *SymbolTable) {

	if c.currentBlock == nil {
		panic("nil block")
	}
	rhs := c.compileExpression(node.Value, table)

	// Allocate mem
	val := c.currentBlock.NewAlloca(rhs.Type())

	// Store
	c.currentBlock.NewStore(rhs, val)

	// add to symbol table
	table.Add(node.Name.Value, SymbolInfo{Name: node.Name.Value, Value: val, Type: rhs.Type()})
}

func (c *Compiler) compileBlockStatement(node *ast.BlockStatement, block *ir.Block, table *SymbolTable) {
	c.compileStatements(node.Statements, block, table)
}

func (c *Compiler) compileStatements(nodes []ast.Statement, block *ir.Block, table *SymbolTable) {

	for _, s := range nodes {
		c.compileStatement(s, block, table)
	}
}

func (c *Compiler) compileReturnStatement(node *ast.ReturnStatement, table *SymbolTable) {

	val := c.compileExpression(node.ReturnValue, table)

	c.currentBlock.NewRet(val)
}
