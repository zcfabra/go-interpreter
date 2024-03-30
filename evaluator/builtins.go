package evaluator

import (
	"fmt"
	"lang/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					return NULL
				}
				return arg.Elements[0]
			case *object.String:
				if len(arg.Value) > 0 {
					return NULL
				}
				return &object.String{Value: string(arg.Value[0])}
			default:
				return newError("argument to `first` not supported, got %s", args[0].Type())
			}
		},
	},
	"rest": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				if len(arg.Elements) == 1 {
					return &object.Array{Elements: []object.Object{}}
				}

				return &object.Array{Elements: arg.Elements[1:]}
			case *object.String:
				if len(arg.Value) == 0 {
					return NULL
				}
				if len(arg.Value) == 1 {
					return &object.String{Value: ""}
				}
				return &object.String{Value: arg.Value[1:]}
			default:
				return newError("argument to `rest` not supported, got %s", args[0].Type())
			}

		},
	},
	"push": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			arr, ok := args[0].(*object.Array)
			if !ok {
				return newError("argument to `push` not supported, got %s", args[0].Type())

			}
			return &object.Array{Elements: append(arr.Elements, args[1])}

		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
