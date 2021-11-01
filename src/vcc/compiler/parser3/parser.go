package parser3

import (
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
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

func Program() Node {
	return Node{}
}

func Stmt() Node {
	return Node{}
}

func Expr() Node {
	return Node{}
}

func Assign() Node {
	return Node{}
}

func Logically() Node {
	return Node{}
}

func Equality() Node {
	return Node{}
}

func Relational() Node {
	return Node{}
}

func Add() Node {
	return Node{}
}

func Mul() Node {
	return Node{}
}

func Unary() Node {
	return Node{}
}

func Primary() Node {
	return Node{}
}
