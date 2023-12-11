package builtins

import (
	"fmt"

	"github.com/mantton/anthe/internal/object"
)

var (
	NULL  = &object.Null{}
	VOID  = &object.Void{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

var BuiltInFunctions = map[string]*object.Builtin{
	"print": {
		Name: "print",
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println("\n" + arg.Inspect())
			}

			return VOID
		},
	},

	"type": {
		Name: "type",
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return VOID
			}

			fmt.Println(args[0].Inspect(), args[0].Type())
			return VOID
		},
	},
}
