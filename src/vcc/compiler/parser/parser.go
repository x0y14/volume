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
	fBody, err = ps._consumeFuncBody()
	if err != nil {
		return Node{}, err
	}
	// 関数内で、"}"を読み飛ばした後なので、goNextはいらない。

	// function
	funcNode = NewFuncDefNode(fName, fArgs, fRetTypes, fBody)

	return funcNode, err
}

func (ps *Parser) Parse() error {
	var nodes []Node
	for !ps.isEof() {
		tok := ps.curt()
		switch tok.Typ {
		case tokenizer.FUNC:
			funcNode, err := ps.consumeFunction()
			if err != nil {
				return err
			}
			nodes = append(nodes, funcNode)
		}
	}
	return nil
}
