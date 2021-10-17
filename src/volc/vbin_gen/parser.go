package vbin_gen

import "github.com/x0y14/volume/src/vvm"

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
	return ps.pos >= len(ps.tokens)-1
}

func (ps *Parser) consumeLabeledOperation() (*Operation, error) {
	label := ps.curt().lit
	ps.goNext()

	colon := ps.curt()
	if !(colon.typ == SYMBOL && colon.lit == ":") {
		return nil, SyntaxErr(":", SYMBOL, colon.lit, colon.typ)
	}

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
	}

	return nil, nil
}
