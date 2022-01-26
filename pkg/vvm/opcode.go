package vvm

import (
	"strings"
)

type Opcode int

const (
	_ Opcode = iota
	ILLEGALOpcode

	NOP

	SET // ?

	ADD
	SUB

	CMP
	LT  // <
	GT  // >
	LTE // <=
	GTE // >=

	SJOIN

	JMP
	JZ  // jump zf == 0
	JNZ // jump zf != 0
	JE  // jump a1 == a2
	JNE // jump a1 != a2
	JL  // jump a1 < a2
	JLE // jump a1 <= a2
	JG  // jump a1 > a2
	JGE // jump a1 >= a2

	CALL
	RET

	CP

	PUSH
	POP

	ADDsp
	SUBsp

	ECHO

	EXIT
)

func (opc Opcode) String() string {
	switch opc {
	case NOP:
		return "nop"
	case SET:
		return "set"

	case ADD:
		return "add"
	case SUB:
		return "sub"

	case CMP:
		return "cmp"
	case LT:
		return "lt"
	case GT:
		return "gt"
	case LTE:
		return "lte"
	case GTE:
		return "gte"

	case SJOIN:
		return "sjoin"

	case JMP:
		return "jmp"
	case JZ:
		return "jz"
	case JNZ:
		return "jnz"
	case JE:
		return "je"
	case JNE:
		return "jne"
	case JL:
		return "jl"
	case JLE:
		return "jle"
	case JG:
		return "jg"
	case JGE:
		return "jge"
	case CALL:
		return "call"
	case RET:
		return "ret"
	case CP:
		return "cp"
	case PUSH:
		return "push"
	case POP:
		return "pop"
	case ADDsp:
		return "add_sp"
	case SUBsp:
		return "sub_sp"
	case ECHO:
		return "echo"
	case EXIT:
		return "exit"
	case ILLEGALOpcode:
		return "illegal"
	default:
		return "illegal"
	}
}

func ConvertOpcode(code string) Opcode {
	switch strings.ToLower(code) {
	case "nop":
		return NOP

	case "set":
		return SET

	case "add":
		return ADD
	case "sub":
		return SUB
	case "cmp":
		return CMP

	case "lt":
		return LT
	case "gt":
		return GT
	case "lte":
		return LTE
	case "gte":
		return GTE

	case "sjoin":
		return SJOIN

	case "jmp":
		return JMP
	case "jz":
		return JZ
	case "jnz":
		return JNZ
	case "je":
		return JE
	case "jne":
		return JNE
	case "jl":
		return JL
	case "jle":
		return JLE
	case "jg":
		return JG
	case "jge":
		return JGE
	case "call":
		return CALL
	case "ret":
		return RET

	case "cp":
		return CP

	case "push":
		return PUSH
	case "pop":
		return POP

	case "add_sp":
		return ADDsp
	case "sub_sp":
		return SUBsp

	case "echo":
		return ECHO

	case "exit":
		return EXIT
	default:
		return ILLEGALOpcode
	}
}

func OperandHowManyHas(typ Opcode) int {
	//typ := ConvertOpcode(code)

	switch typ {
	case NOP:
		return 0

	case SET:
		return 2

	case ADD:
		return 2
	case SUB:
		return 2

	case CMP:
		return 2
	case LT:
		return 2
	case GT:
		return 2
	case LTE:
		return 2
	case GTE:
		return 2

	case SJOIN:
		return 2

	case JMP:
		return 1
	case JZ:
		return 1
	case JNZ:
		return 1
	case JE:
		return 2
	case JNE:
		return 2
	case JL:
		return 2
	case JLE:
		return 2
	case JG:
		return 2
	case JGE:
		return 2

	case CALL:
		return 1
	case RET:
		return 0

	case CP:
		return 2

	case PUSH:
		return 1
	case POP:
		return 1

	case ADDsp:
		return 1
	case SUBsp:
		return 1

	case ECHO:
		return 1

	case EXIT:
		return 0

	default:
		return 0
	}
}
