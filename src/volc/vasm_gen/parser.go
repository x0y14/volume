package vasm_gen

import (
	"fmt"
)

func NewParser(tokens []Token) Parser {
	return Parser{
		pos:    0,
		tokens: tokens,
	}
}

type Parser struct {
	pos    int
	tokens []Token
}

func (ps *Parser) prev() Token {
	return ps.tokens[ps.pos-1]
}
func (ps *Parser) curt() Token {
	return ps.tokens[ps.pos]
}
func (ps *Parser) next() Token {
	return ps.tokens[ps.pos+1]
}

func (ps *Parser) goNext() {
	ps.pos++
}

func (ps *Parser) isEof() bool {
	return ps.curt().typ == EOF || ps.pos >= len(ps.tokens)
}

func (ps *Parser) consumeFuncArgs() (Node, error) {
	// consume "("
	ps.goNext()

	var args []Node

	for ps.curt().typ != RPAREN {

		nod := Node{
			typ:      FuncArgs,
			children: nil,
			tokTyp:   0,
			tok:      Token{},
		}

		if ident := ps.curt(); ident.typ == IDENT {
			nod.tok = ident
			ps.goNext()
		} else {
			return Node{}, fmt.Errorf("syntax error: expect ident, but found %v", ident.String())
		}

		if typ := ps.curt(); IsMoldType(typ.typ) {
			nod.tokTyp = typ.typ
			ps.goNext()
		} else {
			return Node{}, fmt.Errorf("syntax error: expect mold, but found %v", typ.String())
		}

		args = append(args, nod)

		if comma := ps.curt(); comma.typ == COMMA {
			ps.goNext()
		} else {
			break
		}
	}

	// consume ")"
	ps.goNext()

	//fmt.Println("[func args]:")
	//for _, arg := range args {
	//	fmt.Printf("\t(%v) %v\n", arg.tokTyp.String(), arg.tok.String())
	//}

	funcArgsNode := Node{
		typ:      FuncArgs,
		children: args,
		tokTyp:   0,
		tok:      Token{},
	}

	return funcArgsNode, nil
}
func (ps *Parser) consumeFuncReturnValue() Node {

	var rets []Node

	for ps.curt().typ != LBRACE {
		tok := ps.curt()
		rets = append(rets, Node{
			typ:      FuncReturnValue,
			children: nil,
			tok:      Token{},
			tokTyp:   tok.typ,
		})
		ps.goNext()
	}

	//tok := ps.curt()

	//ps.goNext()
	//
	//retNode := Node{
	//	typ:      FuncReturnValue,
	//	children: nil,
	//	tok:      Token{},
	//	tokTyp:   tok.typ,
	//}
	//
	nod := Node{
		typ:      FuncReturnValueCase,
		children: rets,
		tokTyp:   0,
		tok:      Token{},
	}

	return nod
}

func (ps *Parser) consumeFuncMultiReturnValue() (Node, error) {
	var rets []Node
	// } が来るまで
	// consume "("
	ps.goNext()
argLoop:
	for ps.curt().typ != RBRACE {
		tok := ps.curt()
		// 型だったら
		if IsMoldType(tok.typ) {
			retNode := Node{
				typ:      FuncReturnValue,
				children: nil,
				tok:      Token{},
				tokTyp:   tok.typ,
			}
			rets = append(rets, retNode)
			// consume mold
			ps.goNext()
		}

		switch ps.curt().typ {
		case RPAREN:
			// consume ")"
			ps.goNext()
			break argLoop
		case COMMA:
			// consume ","
			ps.goNext()
			continue
		default:
			curt := ps.curt()
			return Node{}, fmt.Errorf("syntax error: %v", curt.String())
		}
	}

	//fmt.Println("[func return value]:")
	//for _, r := range rets {
	//	fmt.Printf("\t%v\n", r.String())
	//}

	nod := Node{
		typ:      FuncReturnValueCase,
		children: rets,
		tokTyp:   0,
		tok:      Token{},
	}

	return nod, nil
}

func (ps *Parser) consumeDefineVariable() (Node, error) {
	// consume "var"
	if defVar := ps.curt(); defVar.typ != VAR {
		return Node{}, SyntaxErr("consumeDefineVariable", VAR.String(), defVar.typ.String())
	}
	ps.goNext()

	// VariableDefine:
	// - tokType: IDENT
	// - tok: ident
	// - children[0]: Node{typ: VariableData, tokType: any, tok: any}
	variable := Node{
		typ: VariableDefine,
	}

	// consume ident
	if ident := ps.curt(); ident.typ != IDENT {
		return Node{}, SyntaxErr("consumeDefineVariable", IDENT.String(), ident.String())
	} else {
		variable.tok = ident
		variable.tokTyp = IDENT
		ps.goNext()
	}

	// consume "="
	if eq := ps.curt(); eq.typ != EQUAL {
		return Node{}, SyntaxErr("consumeDefineVariable", EQUAL.String(), eq.String())
	}
	ps.goNext()

	// consume data
	if data := ps.curt(); IsLiteralType(data.typ) == false {
		return Node{}, SyntaxErr("consumeDefineVariable", "LITERAL", data.typ.String())
	} else {
		dataNode := Node{
			typ:    VariableData,
			tokTyp: data.typ,
			tok:    data,
		}
		variable.children = append(variable.children, dataNode)
		ps.goNext()
	}

	return variable, nil
}

func (ps *Parser) consumeCallFunc() (Node, error) {

	var callFunc Node
	var arguments []Node

	// f(a1, a2)

	// consume f
	if ident := ps.curt(); ident.typ != IDENT {
		return Node{}, SyntaxErr("consumeCallFunc", IDENT.String(), ident.typ.String())
	} else {
		callFunc = Node{
			typ:    CallFunc,
			tokTyp: IDENT,
			tok:    ident,
		}
		ps.goNext()
	}

	// consume "("
	if lp := ps.curt(); lp.typ != LPAREN {
		return Node{}, SyntaxErr("consumeCallFunc", LPAREN.String(), lp.typ.String())
	}
	ps.goNext()

	// consume arguments
	for ps.curt().typ != RPAREN {
		tok := ps.curt()
		if IsLiteralType(tok.typ) || tok.typ == IDENT {
			arguments = append(arguments, Node{
				typ:    FuncArgs,
				tokTyp: tok.typ,
				tok:    tok,
			})
			ps.goNext()
		} else {
			return Node{}, SyntaxErr("consumeCallFunc", "literal or ident", tok.typ.String())
		}

		if comma := ps.curt(); comma.typ != COMMA {
			break
		} else {
			ps.goNext()
		}

	}

	// consume ")"
	ps.goNext()

	callFunc.children = arguments
	return callFunc, nil
}

func (ps *Parser) consumeExpr() (Node, error) {
	allowed := []TokenType{INCREMENT, DECREMENT, LT}

	var expr Node

	// a++

	// a
	ident := ps.curt()
	ps.goNext()

	expr = Node{
		typ:    Expr,
		tok:    ident,
		tokTyp: IDENT,
	}

	if opr := ps.curt(); !IsAllowedType(allowed, opr.typ) {
		return Node{}, NotYetImplementedErr("consumeExpr", opr.lit)
	} else {
		expr.children = append(expr.children, Node{
			tok:    opr,
			tokTyp: opr.typ,
			typ:    Operator,
		})
		ps.goNext()
		if opr.typ == INCREMENT || opr.typ == DECREMENT {
			return expr, nil
		}
	}

	if val := ps.curt(); !IsAllowedType([]TokenType{INT}, val.typ) {
		return Node{}, NotYetImplementedErr("consumeExpr", val.lit)
	} else {
		expr.children = append(expr.children, Node{
			tok:    val,
			tokTyp: val.typ,
			typ:    ExprVal,
		})
		ps.goNext()
	}

	return expr, nil
}

func (ps *Parser) consumeSubstitution() (Node, error) {
	return Node{}, nil
}

func (ps *Parser) consumeWhile() (Node, error) {
	// consume "while"
	ps.goNext()

	whileNod := Node{
		typ: WhileLoop,
	}

	expr, err := ps.consumeExpr()
	if err != nil {
		return Node{}, err
	}

	// get expr
	whileNod.children = append(whileNod.children, expr)

	// body

	//// consume "{"
	//ps.goNext()

	body, err := ps.consumeFuncBody()
	if err != nil {
		return Node{}, err
	}

	body.typ = LoopBody
	whileNod.children = append(whileNod.children, body)

	// consume "}"
	ps.goNext()

	return whileNod, nil
}

func (ps *Parser) consumeFuncBody() (Node, error) {
	braceCount := 0

	var nodes []Node

	for !ps.isEof() {
		c := ps.curt()
		n := ps.next()

		switch c.typ {
		case VAR:
			nod, err := ps.consumeDefineVariable()
			if err != nil {
				return Node{}, err
			}
			nodes = append(nodes, nod)
		case IDENT:
			if IsOperatorType(n.typ) && n.typ == EQUAL {
				// 代入
				nod, err := ps.consumeSubstitution()
				if err != nil {
					return Node{}, err
				}
				nodes = append(nodes, nod)
			}
			if IsOperatorType(n.typ) && n.typ != EQUAL {
				// +=, ++, --, ...?
				nod, err := ps.consumeExpr()
				if err != nil {
					return Node{}, err
				}
				nodes = append(nodes, nod)
			}
			if n.typ == LPAREN {
				nod, err := ps.consumeCallFunc()
				if err != nil {
					return Node{}, err
				}
				nodes = append(nodes, nod)
			}
		case WHILE:
			nod, err := ps.consumeWhile()
			if err != nil {
				return Node{}, err
			}
			nodes = append(nodes, nod)
		default:
			ps.goNext()
		}

		if c.typ == LBRACE {
			braceCount++
		} else if c.typ == RBRACE {
			braceCount--
			if braceCount == 0 {
				break
			}
		}
	}

	return Node{
		typ:      FuncBody,
		children: nodes,
		tokTyp:   0,
		tok:      Token{},
	}, nil
}
func (ps *Parser) consumeFunc() (Node, error) {
	// "func" [ident] "(" [args]? ")" [return value]? "{" [body]? "}"

	funcNode := Node{
		typ:      FuncDefine,
		children: nil,
		tokTyp:   0,
		tok:      Token{},
	}

	// children
	// [ Args?, ReturnValue?, Body? ]

	if funcDefineTok := ps.curt(); funcDefineTok.typ != FUNC {
		return Node{}, fmt.Errorf("syntax error")
	}
	// consume funcDefineTok
	ps.goNext()

	//returnValue := Node{typ: FuncDefine, children: nil}

	funcName := ps.curt()
	//fmt.Println("[func name]:")
	//fmt.Printf("\t%v\n", funcName.String())
	// consume funcName
	ps.goNext()

	// set func-name
	funcNode.tokTyp = IDENT
	funcNode.tok = funcName

	args, err := ps.consumeFuncArgs()
	if err != nil {
		return Node{}, err
	}

	funcNode.children = append(funcNode.children, args)

	switch ps.curt().typ {
	//case LBRACE:
	//	// no return-value
	//	funcBody, err := ps.consumeFuncBody()
	//	if err != nil {
	//		return Node{}, err
	//	}
	//	funcNode.children = append(funcNode.children, funcBody)
	case LBRACE, INT, FLOAT, STRING, BOOL, MAP, LIST:
		returnValue := ps.consumeFuncReturnValue()
		funcNode.children = append(funcNode.children, returnValue)
		funcBody, err := ps.consumeFuncBody()
		if err != nil {
			return Node{}, err
		}
		funcNode.children = append(funcNode.children, funcBody)
	case LPAREN:
		returnValue, err := ps.consumeFuncMultiReturnValue()
		if err != nil {
			return Node{}, err
		}
		funcNode.children = append(funcNode.children, returnValue)
		funcBody, err := ps.consumeFuncBody()
		if err != nil {
			return Node{}, err
		}
		funcNode.children = append(funcNode.children, funcBody)
	default:
		return Node{}, fmt.Errorf("syntax error")
	}

	return funcNode, nil
}

func (ps *Parser) Parse() (nodes []Node, err error) {
tokenLoop:
	for !ps.isEof() {
		tok := ps.curt()
		switch tok.typ {
		case FUNC:
			var nod Node
			nod, err = ps.consumeFunc()
			//fmt.Printf("%v\n", nod)
			nod.Status()
			nodes = append(nodes, nod)
		//case VAR:
		//	var nod Node
		//	nod, err = ps.consumeDefineVariable()
		//	nod.Status()
		default:
			fmt.Printf("%v\n", tok.String())
			break tokenLoop
		}
	}

	return nodes, err
}
