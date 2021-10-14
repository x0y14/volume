package vm

import (
	"strings"
)

type Opcode int

const (
	_ Opcode = iota
	_ILLEGALOpcode

	_SET // ?

	_ADD
	_SUB
	_CMP

	_JZ
	_JNZ

	_CALL
	_RET

	_CP

	_PUSH
	_POP

	_ADDsp
	_SUBsp

	_ECHO

	_EXIT
)

func (opc Opcode) String() string {
	switch opc {
	case _SET:
		return "set"
	case _ADD:
		return "add"
	case _SUB:
		return "sub"
	case _CMP:
		return "cmp"
	case _JZ:
		return "jz"
	case _JNZ:
		return "jnz"
	case _CALL:
		return "call"
	case _RET:
		return "ret"
	case _CP:
		return "cp"
	case _PUSH:
		return "push"
	case _POP:
		return "pop"
	case _ADDsp:
		return "add_sp"
	case _SUBsp:
		return "sub_sp"
	case _ECHO:
		return "echo"
	case _EXIT:
		return "exit"
	case _ILLEGALOpcode:
		return "illegal"
	default:
		return "illegal"
	}
}

func ConvertOpcode(code string) Opcode {
	switch strings.ToLower(code) {
	case "set":
		return _SET

	case "add":
		return _ADD
	case "sub":
		return _SUB
	case "cmp":
		return _CMP

	case "jz":
		return _JZ
	case "jnz":
		return _JNZ

	case "call":
		return _CALL
	case "ret":
		return _RET

	case "cp":
		return _CP

	case "push":
		return _PUSH
	case "pop":
		return _POP

	case "add_sp":
		return _ADDsp
	case "sub_sp":
		return _SUBsp

	case "echo":
		return _ECHO

	case "exit":
		return _EXIT
	default:
		return _ILLEGALOpcode
	}
}

func OperandHowManyHas(typ Opcode) int {
	//typ := ConvertOpcode(code)

	switch typ {
	case _SET:
		return 2

	case _ADD:
		return 2
	case _SUB:
		return 2
	case _CMP:
		return 2

	case _JZ:
		return 1
	case _JNZ:
		return 1

	case _CALL:
		return 1
	case _RET:
		return 0

	case _CP:
		return 2

	case _PUSH:
		return 1
	case _POP:
		return 1

	case _ADDsp:
		return 1
	case _SUBsp:
		return 1

	case _ECHO:
		return 1

	case _EXIT:
		return 0

	default:
		return 0
	}
}
