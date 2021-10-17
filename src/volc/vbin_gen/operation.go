package vbin_gen

import "fmt"
import "github.com/x0y14/volume/src/vvm"

type Operation struct {
	opcode   vvm.Opcode
	operands []Operand
	label    string
}

func (op *Operation) String() string {
	args := ""
	for i, arg := range op.operands {
		args += arg.String()
		if len(op.operands) != i {
			args += ", "
		}
	}

	lbl := ""
	if op.label != "" {
		lbl += " as " + op.label
	}

	return fmt.Sprintf("Operation { %10s, [ %20s ] } %v", op.opcode.String(), args, lbl)
}
