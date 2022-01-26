package vvm

import (
	"fmt"
	"strconv"
)

func IsAllowedTokenType(types []TokenType, typ TokenType) bool {
	for _, tokType := range types {
		if tokType == typ {
			return true
		}
	}
	return false
}

func NewToken(typ TokenType, op Opcode, lit string) *Token {
	return &Token{
		typ:  typ,
		op:   op,
		lit:  lit,
		sPos: -1,
		ePos: -1,
	}
}

type Token struct {
	typ TokenType
	op  Opcode

	lit  string
	sPos int
	ePos int
}

func (tk *Token) String() string {
	return fmt.Sprintf("Token{ %03d~%03d | %20s | %20s | `%20s`}", tk.sPos, tk.ePos, tk.typ.String(), tk.op.String(), tk.lit)
}

func (tk *Token) LoadAsString() (string, error) {
	if tk.typ != _STRING {
		return "", LoadInvalidTypeErr(tk.typ, _STRING)
	}
	return tk.lit, nil
}

func (tk *Token) LoadAsFloat() (float64, error) {
	if tk.typ != _FLOAT {
		return 0, LoadInvalidTypeErr(tk.typ, _FLOAT)
	}

	fl, err := strconv.ParseFloat(tk.lit, 64)
	if err != nil {
		return 0, err
	}
	return fl, nil
}

func (tk *Token) LoadAsInt() (int, error) {
	if tk.typ != _INT && tk.typ != _RTNAddr {
		return 0, LoadInvalidTypeErr(tk.typ, _INT)
	}
	in, err := strconv.Atoi(tk.lit)
	if err != nil {
		return 0, err
	}
	return in, nil
}
