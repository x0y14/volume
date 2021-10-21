package vasm_gen

import "fmt"

type NodeType int

const (
	_ NodeType = iota

	ImportLib

	FuncDefine
	FuncArgs
	FuncReturnValueCase
	FuncReturnValue
	FuncBody
	FuncReturn

	CallFunc

	Expr
	ExprVal
	Operator

	ForLoop
	ForExpr
	WhileLoop
	WhileExpr

	LoopBody

	VariableDefine
	VariableName
	VariableData
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
	case FuncDefine:
		// name - ident
		fmt.Printf("func %q", nod.tok.lit)

		// [args?, ret-val?, body]

		// 3(args, ret-val, body)
		// 2(args, body) or 2(ret-val, body)
		// 1(body)

		switch len(nod.children) {
		case 3:
			//fmt.Printf("")

			args := nod.children[0]
			retVal := nod.children[1]
			body := nod.children[2]

			// args
			fmt.Printf(" ( ")
			for i, arg := range args.children {
				fmt.Printf("%q %v", arg.tok.lit, arg.tokTyp.String())
				if len(args.children)-1 != i {
					fmt.Printf(", ")
				}
			}
			fmt.Printf(" )")

			// return values
			fmt.Printf(" ( ")
			for i, ret := range retVal.children {
				fmt.Printf("%v", ret.tokTyp.String())
				if len(retVal.children)-1 != i {
					fmt.Printf(", ")
				}
			}
			fmt.Printf(" )")

			// body
			fmt.Printf(" {\n")
			for line_i, line := range body.children {
				// line
				//for tok_i, tok := range line.children {
				//	fmt.Printf("%v", tok)
				//	if len(line.children)-1 != tok_i {
				//		fmt.Printf(", ")
				//	}
				//}
				//fmt.Printf("%v", tok)
				line.Status()
				if len(body.children)-1 != line_i {
					fmt.Printf("\n")
				}
			}
			fmt.Printf("}")

		case 2:
			fmt.Printf("2 type")

		case 1:
			//fmt.Printf("body only")
			fmt.Printf(" ()")
			body := nod.children[0]
			fmt.Printf(" {\n")
			for line_i, line := range body.children {
				line.Status()
				if len(body.children)-1 != line_i {
					fmt.Printf("\n")
				}
			}
			fmt.Printf("}")

		default:
			fmt.Printf("function syntax error")
		}

	case VariableDefine:
		ident := nod.tok
		data := nod.children[0]

		fmt.Printf("\tvar %q = %q(%v)", ident.lit, data.tok.lit, data.tokTyp.String())

	case CallFunc:
		ident := nod.tok
		args := nod.children

		fmt.Printf("\t%q( ", ident.lit)

		for i, arg := range args {
			fmt.Printf("%q", arg.tok.lit)
			if len(args)-1 != i {
				fmt.Printf(", ")
			}
		}

		fmt.Printf(" )")

	case WhileLoop:
		expr := nod.children[0]
		fmt.Printf("\t while ")
		expr.Status()

	case Expr:
		fmt.Printf("%q", nod.tok.lit)

		for _, child := range nod.children {
			//fmt.Printf("%v")
			if child.typ == Operator {
				fmt.Printf(" %v ", child.tok.lit)
			} else {
				fmt.Printf("%q", child.tok.lit)
			}
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
}
