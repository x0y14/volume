package parser

import (
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
)

func NewParser(tokens []tokenizer.Token) Parser {
	return Parser{
		pos:    0,
		tokens: tokens,
	}
}

type Parser struct {
	pos    int
	tokens []tokenizer.Token
}

func (ps *Parser) isEof() bool {
	return ps.curt().Typ == tokenizer.EOF
}
func (ps *Parser) prev() tokenizer.Token {
	return ps.tokens[ps.pos-1]
}
func (ps *Parser) curt() tokenizer.Token {
	return ps.tokens[ps.pos]
}
func (ps *Parser) next() tokenizer.Token {
	return ps.tokens[ps.pos+1]
}
func (ps *Parser) goNext() {
	ps.pos++
}

/* 関数
func f(x int) (int) {}
*/
func (ps *Parser) _consumeFuncFormalArgs() (Node, error) {
	// "("
	if _lp := ps.curt(); _lp.Typ != tokenizer.LPAREN {
		return Node{}, SyntaxErr("consumeFunction", tokenizer.LPAREN.String(), _lp.Typ.String())
	}
	ps.goNext()

	var arguments []Node

	// FuncArg
argsLoop:
	for ps.curt().Typ != tokenizer.RPAREN {
		ident := ps.curt()
		if ident.Typ != tokenizer.IDENT {
			return Node{}, SyntaxErr("_consumeFuncFormalArgs", tokenizer.IDENT.String(), ident.Typ.String())
		}
		ps.goNext()

		typ := ps.curt()
		if !tokenizer.IsMoldType(typ.Typ) {
			return Node{}, SyntaxErr("_consumeFuncFormalArgs", "MOLD", typ.Typ.String())
		}
		ps.goNext()

		arg, err := NewFuncArgNode(ident, typ)
		if err != nil {
			return Node{}, err
		}
		arguments = append(arguments, arg)

		switch ps.curt().Typ {
		case tokenizer.COMMA:
			ps.goNext()
			continue
		case tokenizer.RPAREN:
			break argsLoop
		default:
			return Node{}, SyntaxErr("consumeFunction", "COMMA, RPAREN", ps.curt().Typ.String())
		}
	}

	// ")"
	if _rp := ps.curt(); _rp.Typ != tokenizer.RPAREN {
		return Node{}, SyntaxErr("consumeFunction", tokenizer.RPAREN.String(), _rp.Typ.String())
	}
	ps.goNext()

	argumentsNode, err := NewFuncFormalArgsNode(arguments)
	if err != nil {
		return Node{}, err
	}
	return argumentsNode, nil
}

func (ps *Parser) _consumeFuncRetTypes() (Node, error) {
	// 複数の戻り値をサポートする予定はないが、構文的にはサポートしておく。
	// もし、"("があったら、複数の値を持っている可能性がある。
	// また、基本は型が一つ配置されている。
	// single : int
	// multi  : (int, int, string)
	var fRetTypes []tokenizer.Token

	if _lp := ps.curt(); _lp.Typ == tokenizer.LPAREN {
		// multi
		// "("
		ps.goNext()

		// retTypes
	retTypesLoop:
		for ps.curt().Typ != tokenizer.RPAREN {
			ret := ps.curt()
			if !tokenizer.IsMoldType(ret.Typ) {
				return Node{}, SyntaxErr("_consumeFuncRetTypes", "MOLD", ret.Typ.String())
			}
			ps.goNext()

			fRetTypes = append(fRetTypes, ret)

			switch ps.curt().Typ {
			case tokenizer.COMMA:
				ps.goNext()
				continue
			case tokenizer.RPAREN:
				break retTypesLoop
			default:
				return Node{}, SyntaxErr("_consumeFuncRetTypes", "COMMA, RPAREN", ps.curt().Typ.String())
			}

		}

		// ")"
		if _rp := ps.curt(); _rp.Typ != tokenizer.RPAREN {
			return Node{}, SyntaxErr("_consumeFuncRetTypes", tokenizer.RPAREN.String(), _rp.Typ.String())
		}
		ps.goNext()

	} else {
		// single
		if typ := ps.curt(); !tokenizer.IsMoldType(typ.Typ) {
			// 戻り値が記述されていないパターン
			if typ.Typ != tokenizer.LBRACE {
				return Node{}, SyntaxErr("_consumeFuncRetTypes", "MOLD", typ.Typ.String())
			}
		} else {
			fRetTypes = append(fRetTypes, typ)
			ps.goNext()
		}
	}

	fRetTypesNode := NewFuncRetTypeNode(fRetTypes)

	return fRetTypesNode, nil
}

func (ps *Parser) _consumeFuncBody() (Node, error) {
	// todo : body
	if _lb := ps.curt(); _lb.Typ != tokenizer.LBRACE {
		return Node{}, SyntaxErr("_consumeFuncBody", tokenizer.LBRACE.String(), _lb.Typ.String())
	}
	ps.goNext()

	if _rb := ps.curt(); _rb.Typ != tokenizer.RBRACE {
		return Node{}, SyntaxErr("_consumeFuncBody", tokenizer.RBRACE.String(), _rb.Typ.String())
	}
	ps.goNext()

	fBody := NewFuncBodyNode(nil)

	return fBody, nil
}

func (ps *Parser) consumeVarDef2() (Node, error) {
	// "var"
	if _var := ps.curt(); _var.Typ != tokenizer.VAR {
		return Node{}, SyntaxErr("consumeVarDef2", tokenizer.VAR.String(), _var.Typ.String())
	}
	ps.goNext()

	var identTok tokenizer.Token

	// ident
	if _ident := ps.curt(); _ident.Typ != tokenizer.IDENT {
		return Node{}, SyntaxErr("consumeVarDef2", tokenizer.IDENT.String(), _ident.Typ.String())
	} else {
		identTok = _ident
	}

	var rhs []tokenizer.Token

	// RHS 右辺
	for ps.curt().Typ != tokenizer.SEMI {
		rhs = append(rhs, ps.curt())
		ps.goNext()
	}
	// ";"
	ps.goNext()

	varDataNode, err := NewRHSNode(rhs)
	if err != nil {
		return Node{}, err
	}

	return Node{
		typ:           VarDef,
		childrenToken: []tokenizer.Token{identTok},
		childrenNode:  []Node{varDataNode},
	}, nil
}

func (ps *Parser) consumeStartWithIdent() (Node, error) {
	ident := ps.curt()
	ps.goNext()

	// ident;

	// ++, --
	if _opr := ps.curt(); tokenizer.IsAllowedType([]tokenizer.TokenType{tokenizer.INCREMENT, tokenizer.DECREMENT}, _opr.Typ) {
		// ident++なのに;がない
		if ps.next().Typ != tokenizer.SEMI {
			return Node{}, SyntaxErr("", "", "")
		}
		// opr
		ps.goNext()
		// semi
		ps.goNext()
		return Node{
			typ:           ControlExpr,
			childrenToken: []tokenizer.Token{_opr, ident},
			childrenNode:  nil,
		}, nil
	}

	// += -=
	if _opr := ps.curt(); tokenizer.IsAllowedType([]tokenizer.TokenType{tokenizer.PLUSEq, tokenizer.MINUSEq}, _opr.Typ) {
		// += -=
		ps.goNext()
		var rhs []tokenizer.Token

		for ps.curt().Typ != tokenizer.SEMI {
			rhs = append(rhs, ps.curt())
			ps.goNext()
		}
		// ";"
		ps.goNext()

		rhsNode, err := NewRHSNode(rhs)
		if err != nil {
			return Node{}, err
		}

		return Node{
			typ:           AssignExpr,
			childrenToken: []tokenizer.Token{_opr, ident},
			childrenNode:  []Node{rhsNode},
		}, nil
	}

	// :=
	if _opr := ps.curt(); tokenizer.IsAllowedType([]tokenizer.TokenType{tokenizer.COLONEq}, _opr.Typ) {
		// :=
		ps.goNext()

		var rhs []tokenizer.Token
		for ps.curt().Typ != tokenizer.SEMI {
			rhs = append(rhs, ps.curt())
			ps.goNext()
		}
		// ;
		ps.goNext()

	}

	// ident()
	if _opr := ps.curt(); _opr.Typ == tokenizer.LPAREN {
		// (
		ps.goNext()
	}

	// else
	// ident + ...

	return Node{}, nil
}

func (ps *Parser) consumeLine() (Node, error) {
	//for ps.curt().Typ != tokenizer.SEMI {
	//	tok := ps.curt()
	//	//fmt.Printf("%v\n", tok.String())
	//	//ps.goNext()
	//
	//}
	//fmt.Printf("End of line\n\n")
	//// ;
	//ps.goNext()

	switch ps.curt().Typ {
	case tokenizer.VAR:
		nod, err := ps.consumeVarDef2()
		if err != nil {
			return Node{}, err
		}
		return nod, nil
	case tokenizer.IDENT:
		nod, err := ps.consumeStartWithIdent()
		if err != nil {
			return Node{}, err
		}
		return nod, nil
	default:
		return Node{}, NotYetImplErr("consumeLine", ps.curt().Typ.String())
	}
}

func (ps *Parser) consumeContents() (Node, error) {
	var contents []Node

	if _lb := ps.curt(); _lb.Typ != tokenizer.LBRACE {
		return Node{}, SyntaxErr("consumeLines", tokenizer.LBRACE.String(), _lb.Typ.String())
	}
	ps.goNext()

	for ps.curt().Typ != tokenizer.RBRACE {
		_, err := ps.consumeLine()
		if err != nil {
			return Node{}, err
		}
	}

	if _rb := ps.curt(); _rb.Typ != tokenizer.RBRACE {
		return Node{}, SyntaxErr("consumeLines", tokenizer.RBRACE.String(), _rb.Typ.String())
	}
	ps.goNext()

	return NewContentsNode(contents), nil
}

func (ps *Parser) consumeFunction() (funcNode Node, err error) {
	// func f( <args> ) ( <ret> ) { <body> }

	var fName tokenizer.Token
	var fArgs Node
	var fRetTypes Node
	var fBody Node

	// "func"
	if _func := ps.curt(); _func.Typ != tokenizer.FUNC {
		return Node{}, SyntaxErr("consumeFunction", tokenizer.FUNC.String(), _func.Typ.String())
	}
	ps.goNext()

	// "f"
	if _ident := ps.curt(); _ident.Typ != tokenizer.IDENT {
		return Node{}, SyntaxErr("consumeFunction", tokenizer.IDENT.String(), _ident.Typ.String())
	} else {
		fName = _ident
	}
	ps.goNext()

	// "(" <args> ")"
	fArgs, err = ps._consumeFuncFormalArgs()
	if err != nil {
		return Node{}, err
	}
	// 関数内で、")"を読み飛ばした後なので、goNextはいらない。

	// "("? <ret> ")"?
	fRetTypes, err = ps._consumeFuncRetTypes()
	if err != nil {
		return Node{}, err
	}
	// 関数内で、戻り値を読み飛ばしているので、goNextはいらない。

	// "{" <body> "}"
	// contentsに差し替えてみる.
	//fBody, err = ps._consumeFuncBody()
	fBody, err = ps.consumeContents()
	if err != nil {
		return Node{}, err
	}
	// 関数内で、"}"を読み飛ばした後なので、goNextはいらない。

	// function
	funcNode = NewFuncDefNode(fName, fArgs, fRetTypes, fBody)

	return funcNode, err
}

/* Import
import "lib"
*/
func (ps *Parser) consumeImport() (Node, error) {
	// import "lib"
	var library tokenizer.Token

	if _imp := ps.curt(); _imp.Typ != tokenizer.IMPORT {
		return Node{}, SyntaxErr("consumeImport", tokenizer.IMPORT.String(), _imp.Typ.String())
	}
	ps.goNext()

	if _lib := ps.curt(); _lib.Typ != tokenizer.STRING {
		return Node{}, SyntaxErr("consumeImport", tokenizer.STRING.String(), _lib.Typ.String())
	} else {
		library = _lib
	}
	ps.goNext()

	return NewImportNode(library), nil
}

/* 変数定義
var ident = any
*/
func (ps *Parser) consumeVarDef() (Node, error) {
	var varName tokenizer.Token
	// IDENT
	var varData Node
	varDataTypeExpected := []tokenizer.TokenType{
		tokenizer.IDENT, tokenizer.STRING, tokenizer.INT, tokenizer.FLOAT,
		tokenizer.TRUE, tokenizer.FALSE, tokenizer.MAP, tokenizer.LIST, tokenizer.NULL,
	}

	if _var := ps.curt(); _var.Typ != tokenizer.VAR {
		return Node{}, SyntaxErr("consumeVarDef", tokenizer.VAR.String(), _var.Typ.String())
	}
	ps.goNext()

	if _ident := ps.curt(); _ident.Typ != tokenizer.IDENT {
		return Node{}, SyntaxErr("consumeVarDef", tokenizer.IDENT.String(), _ident.Typ.String())
	} else {
		varName = _ident
	}
	ps.goNext()

	if _eq := ps.curt(); _eq.Typ != tokenizer.EQUAL {
		return Node{}, SyntaxErr("consumeVarDef", tokenizer.EQUAL.String(), _eq.Typ.String())
	}
	ps.goNext()

	// todo : 式が入る可能性があるので、個別の解析が必要
	// todo : consumeRHS (RHS = 右辺)
	if _data := ps.curt(); !tokenizer.IsAllowedType(varDataTypeExpected, _data.Typ) {
		return Node{}, SyntaxErr(
			"consumeVarDef",
			"[INDENT, STRING, INT, FLOAT, TRUE, FALSE, MAP, LIST, NULL]",
			_data.Typ.String())
	} else {
		// callFuncかどうか
		if ps.next().Typ == tokenizer.LPAREN {
			nod, err := ps.consumeCallFunc()
			if err != nil {
				return Node{}, err
			}
			varData = nod
		} else {
			varData = NewVarDataNode(_data)
		}
	}
	//varData, err := ps.consumeRHS()
	//if err != nil {
	//	return Node{}, err
	//}
	ps.goNext()

	return NewVarDefNode(varName, varData), nil
}

func (ps *Parser) consumeCalcExpr() {
}

func (ps *Parser) consumeFormula() {}

func (ps *Parser) consumeCallFunc() (Node, error) {
	return Node{}, nil
}

func (ps *Parser) consumeRHS() (Node, error) {
	// 予想される右辺について。
	// G1 = [INDENT, STRING, INT, FLOAT, TRUE, FALSE, MAP, LIST, NULL]
	// G2 = [CallFunc]
	// G3 = [INDENT, STRING, INT, FLOAT, CallFunc]

	// 1. G1, G2,の単体
	// 2. G3を使用した式 (全ての型が一致している必要がある。)
	// 再起的な解析が必要。

	var variables []Node

	valExpectTypes := []tokenizer.TokenType{
		tokenizer.IDENT, tokenizer.STRING, tokenizer.INT, tokenizer.FLOAT,
	}
	calcExpectTypes := []tokenizer.TokenType{
		tokenizer.PERCENT, tokenizer.AST, tokenizer.PLUS, tokenizer.MINUS, tokenizer.SLASH,
	}

rhsLoop:
	for !ps.isEof() {
		tok := ps.curt()
		switch tok.Typ {
		case tokenizer.IDENT, tokenizer.STRING, tokenizer.INT, tokenizer.FLOAT:
			// 項
			if tok.Typ == tokenizer.IDENT && ps.next().Typ == tokenizer.LPAREN {
				// 関数呼び出し
				nod, err := ps.consumeCallFunc()
				if err != nil {
					return Node{}, err
				}
				variables = append(variables, nod)
			} else {
				nod := Node{
					typ:           VarData,
					childrenToken: []tokenizer.Token{tok},
					childrenNode:  nil,
				}
				variables = append(variables, nod)
				ps.goNext()
			}

			// 演算子
			//if _opr := ps.curt(); _opr.Typ != _opr {
			//}

			_tok2 := ps.curt()
			if !tokenizer.IsAllowedType(calcExpectTypes, _tok2.Typ) {
				break rhsLoop
			}
			exprOpr := Node{
				typ:           ExprOperator,
				childrenToken: []tokenizer.Token{_tok2},
				childrenNode:  nil,
			}
			variables = append(variables, exprOpr)
			ps.goNext()

			if !tokenizer.IsAllowedType(valExpectTypes, ps.curt().Typ) {
				return Node{}, SyntaxErr("consumeRHS", "[INDENT, STRING, INT, FLOAT]", ps.curt().Typ.String())
			}

		case tokenizer.TRUE, tokenizer.FALSE, tokenizer.MAP, tokenizer.LIST, tokenizer.NULL:
			ps.goNext()
			return Node{
				typ:           VarData,
				childrenToken: []tokenizer.Token{tok},
				childrenNode:  nil,
			}, nil
		default:
			return Node{}, NotYetImplErr("consumeRHS", tok.Typ.String())
		}
	}
	// 式解析

	return Node{
		typ:           RHS,
		childrenToken: nil,
		childrenNode:  variables,
	}, nil
}

func (ps *Parser) Parse() ([]Node, error) {
	var nodes []Node
	for !ps.isEof() {
		tok := ps.curt()
		switch tok.Typ {
		case tokenizer.FUNC:
			funcNode, err := ps.consumeFunction()
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, funcNode)
		case tokenizer.IMPORT:
			importNode, err := ps.consumeImport()
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, importNode)
		}
	}
	return nodes, nil
}
