package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/mantton/anthe/internal/evaluator"
	"github.com/mantton/anthe/internal/lexer"
	"github.com/mantton/anthe/internal/object"
	"github.com/mantton/anthe/internal/parser"
)

const PROMPT = ">> "

func main() {

	allArgs := os.Args[1:]
	argCount := len(allArgs)

	if argCount > 1 { // not enough args provided
		fmt.Println("Usage: anthe [script]")
		os.Exit(64)
	} else if argCount == 1 {

		// filePath := allArgs[0]

		// if err != nil {
		// 	fmt.Printf("\nErrors:\n%s", err.Error())
		// }

	} else {
		fmt.Println("Anthe REPL")
		test := object.New(nil)

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

			evaluator, err := evaluator.Eval(prog, test)

			if err != nil {
				fmt.Println(err.Error())

			}

			if evaluator != nil {
				fmt.Println("\nOUTPUT: " + evaluator.Inspect())
			}
		}
	}

}
