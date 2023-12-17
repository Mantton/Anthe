package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/mantton/anthe/internal/compiler"
	"github.com/mantton/anthe/internal/evaluator"
	"github.com/mantton/anthe/internal/lexer"
	"github.com/mantton/anthe/internal/parser"
)

const PROMPT = ">> "

func main() {

	allArgs := os.Args[1:]
	argCount := len(allArgs)

	e := evaluator.New()

	if argCount > 1 { // not enough args provided
		fmt.Println("Usage: anthe [script]")
		os.Exit(64)
	} else if argCount == 1 {

		path := allArgs[0]

		data, err := os.ReadFile(path)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		l := lexer.New(string(data), "repl.an")
		p := parser.New(l)

		prog := p.ParseProgram()

		c := compiler.New()
		result, err := c.Compile(prog)
		if err != nil {
			fmt.Println(err.Error())
			return

		}

		fmt.Println("\n\n\n")
		fmt.Println(result)
		fmt.Println("done")

	} else {
		fmt.Println("Anthe REPL")

		for {
			fmt.Print("\n>>> ")

			reader := bufio.NewReader(os.Stdin)

			line, err := reader.ReadString('\n')

			if err != nil {
				if err == io.EOF {
					return
				}
				panic(err)
			}

			if len(line) <= 1 {
				fmt.Println("-")
				continue
			}

			if line == "exit()\n" {
				fmt.Println("Bye!")
				return
			}

			l := lexer.New(line, "repl.an")
			p := parser.New(l)

			prog := p.ParseProgram()

			if prog == nil {
				fmt.Println("program failed to parse")
				continue
			}

			if len(prog.Errors) > 0 {
				for _, err := range prog.Errors {
					fmt.Println(err)
				}
				continue
			}

			// checker := typing.New(prog.Statements)

			// ok, t_err := checker.CheckAll()

			// if !ok {
			// 	fmt.Println("Type Checker : Errors")
			// 	for _, msg := range t_err {
			// 		fmt.Println(msg)
			// 	}
			// }

			result, err := e.RunProgram(prog)

			if err != nil {
				fmt.Println(err.Error())

			}

			if result != nil && result.Type() != "void" {
				fmt.Println("\nOUTPUT: " + result.Inspect())
			}
		}
	}

}
