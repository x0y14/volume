package vm

type TokenType int

const (
	_ TokenType = iota
	_ILLEGALToken

	_OPERATION // call, push, pop, ...
	_STRING    // \"aaa\", ....
	_INT       // 123, ...
	_FLOAT
	_ADDR    // [bp+1], [sp+1], [sp-1], ...
	_KEYWORD // bp, sp, reg_a, ...
	_COMMENT
)

func (tkTyp TokenType) String() string {
	switch tkTyp {
	case _OPERATION:
		return "OPERATION"
	case _STRING:
		return "STRING"
	case _INT:
		return "INT"
	case _FLOAT:
		return "FLOAT"
	case _ADDR:
		return "ADDR"
	case _KEYWORD:
		return "KEYWORD"
	case _COMMENT:
		return "COMMENT"
	default:
		return "ILLEGALToken"
	}
}
