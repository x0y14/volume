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

func (ps *Parser) _consumeFuncArgs() (Node, error) {
	// "("
	if _lpArgs := ps.curt(); _lpArgs.Typ != tokenizer.LPAREN {
		return Node{}, SyntaxErr("consumeFunction", tokenizer.LPAREN.String(), _lpArgs.Typ.String())
	}
	ps.goNext()

	var arguments []Node

	// <args>
argsLoop:
	for !ps.isEof() {
		ident := ps.curt()
		if ident.Typ != tokenizer.IDENT {
			return Node{}, SyntaxErr("_consumeFuncArgs", tokenizer.IDENT.String(), ident.Typ.String())
		}
		ps.goNext()

		typ := ps.curt()
		if !tokenizer.IsMoldType(typ.Typ) {
			return Node{}, SyntaxErr("_consumeFuncArgs", "MOLD", typ.Typ.String())
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
	if _rpArgs := ps.curt(); _rpArgs.Typ != tokenizer.RPAREN {
		return Node{}, SyntaxErr("consumeFunction", tokenizer.RPAREN.String(), _rpArgs.Typ.String())
	}
	ps.goNext()

	argumentsNode, err := NewFuncFormalArgsNode(arguments)
	if err != nil {
		return Node{}, err
	}
	return argumentsNode, nil
}

func (ps *Parser) consumeFunction() (Node, error) {
	// func f( <args> ) ( <ret> ) { <body> }

	var fName tokenizer.Token
	var fArgs Node

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
	fArgs, err := ps._consumeFuncArgs()
	if err != nil {
		return Node{}, err
	}

	// todo : retType
	// todo : body

	functionDefine := NewFuncDefNode(fName, fArgs, Node{}, Node{})

	return functionDefine, nil
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
