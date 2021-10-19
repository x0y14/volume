package vasm_gen

import "fmt"

type NodeType int

const (
	_ NodeType = iota

	FUNCTION
	FuncArgs
	FuncReturnValueCase
	FuncReturnValue
	FuncBody

	VARIABLE
	ValName
	ValData
)

type Node struct {
	typ      NodeType
	children []Node

	tokTyp TokenType
	tok    Token
}

// func args
// Node{  }

func (nod Node) Status() {
	switch nod.typ {
	case FUNCTION:
		// name - ident
		fmt.Printf("func %q", nod.tok.lit)

		// [args?, ret-val?, body]
		switch len(nod.children) {
		case 3:
			//fmt.Printf("")

			args := nod.children[0]
			//retVal := nod.children[1]
			//body := nod.children[2]

			fmt.Printf(" ( ")
			for i, arg := range args.children {
				fmt.Printf("%q %v", arg.tok.lit, arg.tokTyp.String())
				if len(args.children)-1 != i {
					fmt.Printf(", ")
				}
			}
			fmt.Printf(" )")

		case 2:
			fmt.Printf("2 type")

		case 1:
			fmt.Printf("body only")
		default:
			fmt.Printf("function syntax error")
		}

		// 3(args, ret-val, body)
		// 2(args, body) or 2(ret-val, body)
		// 1(body)

	}
}
