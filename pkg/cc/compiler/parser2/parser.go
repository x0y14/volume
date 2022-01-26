package parser2

import (
	"fmt"
	"github.com/x0y14/volume/pkg/cc/compiler/tokenizer"
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

func (ps *Parser) consumeImport() (Node, error) {
	// import "lib"
	var libName tokenizer.Token

	// "import"
	ps.goNext()

	if _lib := ps.curt(); _lib.Typ != tokenizer.STRING {
		return Node{}, SyntaxErr("consumeImport", tokenizer.STRING.String(), _lib.Typ.String())
	} else {
		libName = _lib
		// "lib"
		ps.goNext()
	}

	imp := NewImportNode(&libName)

	return imp, nil
}

func (ps *Parser) consumeFuncDef() (Node, error) {
	// func main( params ) (returnValues) { scripts }

	var name tokenizer.Token
	var params Node
	var returnValues Node
	var scripts Node

	// "func"
	if _func := ps.curt(); _func.Typ != tokenizer.FUNC {
		return Node{}, SyntaxErr("consumeFuncDef", tokenizer.FUNC.String(), _func.Typ.String())
	} else {
		ps.goNext()
	}

	// "main"
	if ident := ps.curt(); ident.Typ != tokenizer.IDENT {
		return Node{}, SyntaxErr("consumeFuncDef", tokenizer.IDENT.String(), ident.Typ.String())
	} else {
		name = ident
		ps.goNext()
	}

	// ( params )
	p, err := ps.consumeFuncParam()
	if err != nil {
		return Node{}, err
	}
	params = p

	// (returnValues)

	r, err := ps.consumeFuncReturnValues()
	if err != nil {
		return Node{}, err
	}
	returnValues = r

	// { scripts }
	s, err := ps.consumeScripts()
	if err != nil {
		return Node{}, err
	}
	scripts = s

	node := NewFuncDefNode(&name, params, returnValues, scripts)

	return node, nil
}
func (ps *Parser) consumeFuncParam() (Node, error) {
	// ( args )

	// "("
	if lp := ps.curt(); lp.Typ != tokenizer.LPAREN {
		return Node{}, SyntaxErr("consumeFuncParam", tokenizer.LPAREN.String(), lp.Typ.String())
	}
	ps.goNext()

	var params []Node

	for ps.curt().Typ != tokenizer.RPAREN {
		name := ps.curt()
		if name.Typ != tokenizer.IDENT {
			return Node{}, SyntaxErr("consumeFuncParam", tokenizer.IDENT.String(), name.Typ.String())
		}
		ps.goNext()

		dataTok := ps.curt()
		data, err := NewDataTypeNode(&dataTok)
		if err != nil {
			return Node{}, err
		}
		ps.goNext()

		param := NewFuncParamsItemNode(&name, data)
		params = append(params, param)

		if comma := ps.curt(); comma.Typ == tokenizer.COMMA {
			ps.curt()
			continue
		} else if rp := ps.curt(); rp.Typ == tokenizer.RPAREN {
			break
		} else {
			return Node{}, SyntaxErr("consumeFuncParam", "comma or )", ps.curt().Typ.String())
		}
	}

	// ")"
	if rp := ps.curt(); rp.Typ != tokenizer.RPAREN {
		return Node{}, SyntaxErr("consumeFuncParam", tokenizer.RPAREN.String(), rp.Typ.String())
	}
	ps.goNext()

	node := NewFuncParamsNode(params)

	return node, nil
}
func (ps *Parser) consumeFuncReturnValues() (Node, error) {
	// 単体の場合は、NewDataTypeNodeで成功する。
	// 単体の場合、次のtokenは、"{"である
	// "("で始まっていれば、複数の値。
	// "("で始まっていないにも関わらず、NewNodeでエラーが出るのであれば、不正な構文

	var values []Node

	if ps.curt().Typ == tokenizer.LBRACE {
		return NewFuncReturnValuesNode(values), nil
	}

	// multi
	if lp := ps.curt(); lp.Typ == tokenizer.LPAREN {
		// "("
		ps.goNext()
	retsLoop:
		for ps.curt().Typ != tokenizer.RPAREN {
			ret := ps.curt()
			nod, err := NewDataTypeNode(&ret)
			if err != nil {
				return Node{}, err
			}
			values = append(values, nod)
			ps.goNext()

			switch ps.curt().Typ {
			case tokenizer.COMMA:
				ps.goNext()
				continue
			case tokenizer.RPAREN:
				break retsLoop
			default:
				return Node{}, SyntaxErr("consumeFuncReturnValues", "comma or }", ps.curt().Typ.String())
			}
		}
		if rp := ps.curt(); rp.Typ != tokenizer.RPAREN {
			return Node{}, SyntaxErr("consumeFuncReturnValues", tokenizer.RPAREN.String(), rp.Typ.String())
		}
		// ")"
		ps.goNext()
	} else {
		// single
		if typ := ps.curt(); !tokenizer.IsMoldType(typ.Typ) {
			// 戻り値が記述されていないパターン
			if typ.Typ != tokenizer.LBRACE {
				return Node{}, SyntaxErr("_consumeFuncRetTypes", "MOLD", typ.Typ.String())
			}
		} else {
			nod, err := NewDataTypeNode(&typ)
			if err != nil {
				return Node{}, err
			}
			values = append(values, nod)
			ps.goNext()
		}
	}

	node := NewFuncReturnValuesNode(values)

	return node, nil
}

func (ps *Parser) consumeScripts() (Node, error) {
	if lb := ps.curt(); lb.Typ != tokenizer.LBRACE {
		return Node{}, SyntaxErr("consumeScripts", tokenizer.LBRACE.String(), lb.Typ.String())
	}
	// "{"
	ps.goNext()

	var lines []Node

	for ps.curt().Typ != tokenizer.RBRACE {
		lineTokens, err := ps._consumeTokenAsLine()
		if err != nil {
			return Node{}, err
		}

		line, err := ps._analyzeLine(lineTokens)
		if err != nil {
			return Node{}, err
		}

		lines = append(lines, line)
	}

	if rb := ps.curt(); rb.Typ != tokenizer.RBRACE {
		return Node{}, SyntaxErr("consumeScripts", tokenizer.RBRACE.String(), rb.Typ.String())
	}
	// "}"
	ps.goNext()

	return NewScriptsNode(lines), nil
}

func (ps *Parser) _consumeTokenAsLine() ([]tokenizer.Token, error) {
	var tokens []tokenizer.Token

	for ps.curt().Typ != tokenizer.SEMI && ps.curt().Typ != tokenizer.RBRACE {
		tok := ps.curt()
		tokens = append(tokens, tok)
		ps.goNext()
	}

	// ";"
	if semi := ps.curt(); semi.Typ != tokenizer.SEMI {
		return nil, SyntaxErr("_consumeLine", tokenizer.SEMI.String(), semi.Typ.String())
	}
	ps.goNext()

	return tokens, nil
}

func (ps *Parser) _analyzeLine(tokens []tokenizer.Token) (Node, error) {
	fmt.Println("Line<<<")
	for _, tok := range tokens {
		fmt.Printf("  `%v` [ %v ]\n", tok.Lit, tok.Typ.String())
	}
	fmt.Printf(">>>\n")
	return Node{}, nil
}

func (ps *Parser) Parse() ([]Node, error) {

	var nodes []Node

	for !ps.isEof() {
		switch ps.curt().Typ {
		case tokenizer.IMPORT:
			nod, err := ps.consumeImport()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, nod)
		case tokenizer.FUNC:
			nod, err := ps.consumeFuncDef()
			if err != nil {
				return nodes, err
			}
			nodes = append(nodes, nod)
		default:
			return nodes, NotYetImplErr("Parse", ps.curt().Typ.String())
		}
	}

	return nodes, nil
}
