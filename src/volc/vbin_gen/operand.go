package vbin_gen

type Operand struct {
	lit string
}

func (opr *Operand) String() string {
	return opr.lit
}
