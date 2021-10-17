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

	SJOIN

	JUMP
	JZ
	JNZ

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
	case SJOIN:
		return "sjoin"
	case JUMP:
		return "jump"
	case JZ:
		return "jz"
	case JNZ:
		return "jnz"
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

	case "sjoin":
		return SJOIN

	case "jump":
		return JUMP
	case "jz":
		return JZ
	case "jnz":
		return JNZ

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

	case SJOIN:
		return 2

	case JUMP:
		return 1
	case JZ:
		return 1
	case JNZ:
		return 1

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
