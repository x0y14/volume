package vasm_gen

import "fmt"

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
	return ps.pos >= len(ps.tokens)
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

		if typ := ps.curt(); IsMold(typ.typ) {
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
	tok := ps.curt()
	//fmt.Println("[func return value]:")
	//fmt.Printf("\t%v\n", tok.String())
	ps.goNext()

	retNode := Node{
		typ:      FuncReturnValue,
		children: nil,
		tok:      Token{},
		tokTyp:   tok.typ,
	}

	nod := Node{
		typ:      FuncReturnValueCase,
		children: []Node{retNode},
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
		if IsMold(tok.typ) {
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

func (ps *Parser) consumeFuncBody() Node {
	return Node{
		typ:      FuncBody,
		children: nil,
		tokTyp:   0,
		tok:      Token{},
	}
}
func (ps *Parser) consumeFunc() (Node, error) {
	// "func" [ident] "(" [args]? ")" [return value]? "{" [body]? "}"

	funcNode := Node{
		typ:      FUNCTION,
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

	//returnValue := Node{typ: FUNCTION, children: nil}

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
	case LBRACE:
		// no return-value
		funcBody := ps.consumeFuncBody()
		funcNode.children = append(funcNode.children, funcBody)
	case INT, FLOAT, STRING, BOOL, MAP, LIST:
		returnValue := ps.consumeFuncReturnValue()
		funcNode.children = append(funcNode.children, returnValue)
		funcBody := ps.consumeFuncBody()
		funcNode.children = append(funcNode.children, funcBody)
	case LPAREN:
		returnValue, err := ps.consumeFuncMultiReturnValue()
		if err != nil {
			return Node{}, err
		}
		funcNode.children = append(funcNode.children, returnValue)
		funcBody := ps.consumeFuncBody()
		funcNode.children = append(funcNode.children, funcBody)
	default:
		return Node{}, fmt.Errorf("syntax error")
	}

	return funcNode, nil
}

func (ps *Parser) Parse() (err error) {
tokenLoop:
	for !ps.isEof() {
		tok := ps.curt()
		switch tok.typ {
		case FUNC:
			var nod Node
			nod, err = ps.consumeFunc()
			//fmt.Printf("%v\n", nod)
			nod.Status()
		default:
			break tokenLoop
		}
	}

	return err
}
