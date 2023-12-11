package evaluator

import (
	"fmt"
	"testing"

	"github.com/mantton/anthe/internal/lexer"
	"github.com/mantton/anthe/internal/object"
	"github.com/mantton/anthe/internal/parser"
)

func TestEvaluator(t *testing.T) {
	line := "let x = 10; let y = 20; return x + y;"
	l := lexer.New(line, "test.an")
	p := parser.New(l)

	prog := p.ParseProgram()

	if prog == nil {
		fmt.Println("program failed to parse")
	}

	if len(prog.Errors) > 0 {
		for _, err := range prog.Errors {
			fmt.Println(err)
		}
	}

	_, err := Eval(prog, object.New(nil))

	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(evaluator.Inspect())
}
