package parser3

import (
	"github.com/x0y14/volume/pkg/vcc/compiler/tokenizer"
	"strconv"
)

// program    = stmt*
// stmt       = expr ";"
//            | "{" stmt* "}"
//            | "for" expr ";" expr ";" expr stmt
//            | "while" logically stmt
//            | "var" stmt
//            | "if" logically stmt ("else if" logically stmt)? ("else" stmt)?
//            | "return" stmt?
//            | "break" ";"
//            | "continue" ";"
//            | "func" ident "(" (ident kind)? | (ident kind ",")* ")" kind* stmt
//            | "import" string
// kind       = int | string | float | bool
// expr       = assign
// assign     = logically ("=" logically)?
// logically  = equality ("||" equality | "&&" equality)
// equality   = relational ("==" equality | "!=" equality)*
// relational = add (">" add | ">=" add | "<" add | "<=" add)*
// add        = mul ("+" mul | "-" mul)*
// mul        = unary ("*" unary | "/" unary)*
// unary      = ("+" | "-")? primary
// primary    = int | string | float | bool | ident | call | "(" expr ")"

func NewParser(tokens []tokenizer.Token) Parser {
	return Parser{
		tokens: tokens,
		pos:    0,
	}
}

type Parser struct {
	tokens []tokenizer.Token
	pos    int
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

func (ps *Parser) isEof() bool {
	return ps.curt().Typ == tokenizer.EOF
}

func (ps *Parser) consume(opr string) bool {
	if ps.next().Lit == opr {
		ps.goNext()
		return true
	}
	return false
}

func (ps *Parser) Program() Node {
	return Node{}
}

func (ps *Parser) Stmt() Node {
	return Node{}
}

func (ps *Parser) Expr() (Node, error) {
	return ps.Assign()
}

func (ps *Parser) Assign() (Node, error) {
	node, err := ps.Logically()
	if err != nil {
		return Node{}, err
	}

	if ps.consume("=") {
		d, err := ps.Logically()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdAssign, node, d), nil
	}

	return node, nil
}

func (ps *Parser) Logically() (Node, error) {
	node, err := ps.Equality()
	if err != nil {
		return Node{}, err
	}
	if ps.consume("&&") {
		d, err := ps.Equality()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdAnd, node, d), nil
	} else if ps.consume("||") {
		d, err := ps.Equality()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdOr, node, d), nil
	}
	return node, nil
}

func (ps *Parser) Equality() (Node, error) {
	node, err := ps.Relational()
	if err != nil {
		return Node{}, err
	}

	if ps.consume("==") {
		d, err := ps.Relational()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdEqual, node, d), nil
	} else if ps.consume("!=") {

	}

	return Node{}, nil
}

func (ps *Parser) Relational() (Node, error) {
	node, err := ps.Add()
	if err != nil {
		return Node{}, err
	}

	if ps.consume("<") {
		d, err := ps.Add()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdLt, node, d), nil
	} else if ps.consume("<=") {
		d, err := ps.Add()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdLte, node, d), nil
	} else if ps.consume(">") {
		d, err := ps.Add()
		if err != nil {
			return Node{}, err
		}

		return NewNode(NdGt, node, d), nil
	} else if ps.consume(">=") {
		d, err := ps.Add()
		if err != nil {
			return Node{}, err
		}

		return NewNode(NdGte, node, d), nil
	}

	return Node{}, nil
}

func (ps *Parser) Add() (Node, error) {
	node, err := ps.Mul()
	if err != nil {
		return Node{}, err
	}

	if ps.consume("+") {
		d, err := ps.Mul()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdAdd, node, d), nil
	} else if ps.consume("-") {
		d, err := ps.Mul()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdSub, node, d), nil
	}
	return node, nil
}

func (ps *Parser) Mul() (Node, error) {
	node, err := ps.Unary()
	if err != nil {
		return Node{}, err
	}
	if ps.consume("*") {
		d, err := ps.Unary()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdMul, node, d), nil
	} else if ps.consume("/") {
		d, err := ps.Unary()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdDiv, node, d), nil
	}
	return node, nil
}

func (ps *Parser) Unary() (Node, error) {
	if ps.consume("+") {
		return ps.Primary()
	} else {
		p, err := ps.Primary()
		if err != nil {
			return Node{}, err
		}
		return NewNode(NdSub, NewIntNode(0), p), nil
	}
}

func (ps *Parser) Primary() (Node, error) {
	switch ps.curt().Typ {
	case tokenizer.INT:
		// int
		i, err := strconv.Atoi(ps.curt().Lit)
		if err != nil {
			return Node{}, err
		}
		ps.goNext()
		return NewIntNode(i), nil
	case tokenizer.STRING:
		// string
		ps.goNext()
		return NewStringNode(ps.curt().Lit), nil
	case tokenizer.FLOAT:
		// float
		f, err := strconv.ParseFloat(ps.curt().Lit, 64)
		if err != nil {
			return Node{}, err
		}
		ps.goNext()
		return NewFloatNode(f), nil
	case tokenizer.BOOL:
		// bool
		if ps.curt().Lit == "true" {
			ps.goNext()
			return NewBoolNode(true), nil
		} else {
			ps.goNext()
			return NewBoolNode(false), nil
		}
	case tokenizer.IDENT:
		// call
		// ident
		if ps.next().Typ == tokenizer.LPAREN {
			c, err := ps.consumeCall()
			if err != nil {
				return Node{}, err
			}
			return c, nil
		} else {
			return NewIdentNode(ps.curt().Lit), nil
		}
	}
	return Node{}, nil
}

func (ps *Parser) consumeCall() (Node, error) {
	return Node{kind: NdCall}, nil
}
