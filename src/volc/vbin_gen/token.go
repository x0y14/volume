package vbin_gen

import "fmt"

type TokenType int

const (
	_ TokenType = iota
	ILLEGAL

	STRING // "abc"
	INT    // 123
	FLOAT  // 123.4
	NULL

	IDENT // var, main, ...

	NEWLINE
	WHITESPACE

	SYMBOL // +-., ...

	COMMENT // start with semicolon
)

func (typ TokenType) String() string {
	switch typ {
	case STRING:
		return "STRING"
	case INT:
		return "INT"
	case FLOAT:
		return "FLOAT"
	case NULL:
		return "NULL"
	case IDENT:
		return "IDENT"
	case NEWLINE:
		return "NEWLINE"
	case WHITESPACE:
		return "WHITESPACE"
	case SYMBOL:
		return "SYMBOL"
	case COMMENT:
		return "COMMENT"
	case ILLEGAL:
		return "ILLEGAL"
	default:
		return "ILLEGAL"
	}
}

type Token struct {
	lit  string
	typ  TokenType
	sPos int
	ePos int
}

func (tok Token) String() string {
	return fmt.Sprintf("Token( pos: %03d <= ... < %03d ) { typ: %20s, lit: %20q }", tok.sPos, tok.ePos, tok.typ.String(), tok.lit)
}

func NewToken(lit string, typ TokenType, sPos int, ePos int) *Token {
	return &Token{
		lit:  lit,
		typ:  typ,
		sPos: sPos,
		ePos: ePos,
	}
}
