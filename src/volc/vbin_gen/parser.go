package vbin_gen

import (
	"fmt"
	"github.com/x0y14/volume/src/vvm"
)

func NewParser(tokens []Token) *Parser {
	return &Parser{
		pos:    0,
		tokens: tokens,
	}
}

type Parser struct {
	pos    int
	tokens []Token
}

func (ps *Parser) curt() *Token {
	return &ps.tokens[ps.pos]
}
func (ps *Parser) next() (*Token, error) {
	if ps.pos+1 >= len(ps.tokens) {
		return nil, InvalidAccess(ps.pos + 1)
	}
	return &ps.tokens[ps.pos+1], nil
}

func (ps *Parser) goNext() {
	ps.pos++
}

func (ps *Parser) isEof() bool {
	return ps.pos >= len(ps.tokens)
}

func (ps *Parser) consumeAddr() *Operand {
	// "[" + "bp" + "+" + "1" + "]"
	lit := ""
	for !ps.isEof() {
		opr := ps.curt()
		lit += opr.lit
		if opr.typ == SYMBOL && opr.lit == "]" {
			break
		}
		ps.goNext()
	}

	return &Operand{lit: lit}
}

func (ps *Parser) consumeLabeledOperation() (*Operation, error) {
	label := ps.curt().lit
	ps.goNext()

	colon := ps.curt()
	if !(colon.typ == SYMBOL && colon.lit == ":") {
		return nil, SyntaxErr(":", SYMBOL, colon.lit, colon.typ)
	}
	ps.goNext()

	op, err := ps.consumeOperation()
	if err != nil {
		return nil, err
	}
	op.label = label

	return op, nil
}

func (ps *Parser) consumeOperation() (*Operation, error) {
	opcodeTok := ps.curt()
	opcode := vvm.ConvertOpcode(opcodeTok.lit)
	if opcode == vvm.ILLEGALOpcode {
		return nil, SyntaxIllegalOpcodeErr(*opcodeTok)
	}
	ps.goNext()

	operandCount := vvm.OperandHowManyHas(opcode)
	var operands []Operand
	i := 0
	for i < operandCount {
		oprTok := ps.curt()
		if oprTok.typ == SYMBOL && oprTok.lit == "," {
			ps.goNext()
			continue
		} else if oprTok.typ == SYMBOL && oprTok.lit == "[" {
			opr := ps.consumeAddr()
			operands = append(operands, *opr)
			ps.goNext()
			i++
			continue
		}
		opr := Operand{lit: oprTok.lit}
		operands = append(operands, opr)
		ps.goNext()
		i++
	}

	//fmt.Printf(">>> %v\n", ps.curt().String())
	return &Operation{
		opcode:   opcode,
		operands: operands,
		label:    "",
	}, nil
}

func (ps *Parser) Parse() (*[]Operation, error) {
	var ops []Operation

	for !ps.isEof() {
		t := ps.curt()
		if t.typ == IDENT {
			nx, err := ps.next()
			if err == nil && nx.typ == SYMBOL && nx.lit == ":" {
				o, err := ps.consumeLabeledOperation()
				if err != nil {
					return nil, err
				}
				ops = append(ops, *o)
			} else {
				o, err := ps.consumeOperation()
				if err != nil {
					return nil, err
				}
				ops = append(ops, *o)
			}
		} else {
			// ?
			fmt.Printf("skipped : %v\n", t.String())
			ps.goNext()
		}
	}

	return &ops, nil
}
