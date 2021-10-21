package vvm

import (
	"fmt"
	"regexp"
	"strconv"
)

func NewVM() *VVM {
	vm := &VVM{
		pc:         0,
		bp:         0,
		sp:         0,
		zf:         0,
		regA:       nil,
		regB:       nil,
		regC:       nil,
		mainMemory: nil,
		stack:      nil,
	}
	return vm
}

type VVM struct {
	pc int
	bp int
	sp int
	zf int

	regA *Token
	regB *Token
	regC *Token

	mainMemory []Token
	stack      []*Token
	stream     []string
}

func (vvm *VVM) loadAsm(path string) {
	tokenizer, err := NewTokenizerPath(path)
	if err != nil {
		panic(err)
	}
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		panic(err)
	}
	vvm.mainMemory = *tokens
}

func (vvm *VVM) reserveStack(n int) {
	var stack []*Token

	for i := 0; i < n; i++ {
		stack = append(stack, nil)
	}
	vvm.stack = stack
}

func (vvm *VVM) setPointer() {
	vvm.sp = len(vvm.stack) - 1
	vvm.bp = len(vvm.stack) - 1
}

func (vvm *VVM) SetUp(stackSize int, programPath string) {
	vvm.reserveStack(stackSize)
	vvm.loadAsm(programPath)
	vvm.setPointer()
}

func (vvm *VVM) writeStream(text string) {
	vvm.stream = append(vvm.stream, text)
}

func (vvm *VVM) isProgEof() bool {
	return vvm.pc >= len(vvm.mainMemory)
}

func (vvm *VVM) movePc(n int) {
	vvm.pc += n
}

func (vvm *VVM) curtProg() Token {
	return vvm.mainMemory[vvm.pc]
}

func (vvm *VVM) operands(n int) []Token {
	var tokens []Token
	for i := 0; i < n; i++ {
		tokens = append(tokens, vvm.mainMemory[vvm.pc+i+1])
	}
	return tokens
}

func (vvm *VVM) addSp(diff int) error {
	if (vvm.sp+diff) < 0 || len(vvm.stack)-1 < (vvm.sp+diff) {
		return StackAccessErr("addSp", len(vvm.stack)-1, vvm.sp+diff)
	}
	vvm.sp += diff
	return nil
}
func (vvm *VVM) subSp(diff int) error {
	// [0, 1, 2, 3] : len() => 4
	if (vvm.sp-diff) < 0 || len(vvm.stack)-1 < (vvm.sp-diff) {
		return StackAccessErr("subSp", len(vvm.stack)-1, vvm.sp-diff)
	}
	vvm.sp -= diff
	return nil
}

func (vvm *VVM) IsValidPointerLocation(pointerType PointerType, pointer int) bool {
	switch pointerType {
	case _BasePointer, _StackPointer:
		return 0 <= pointer && pointer <= len(vvm.stack)-1
	default:
		return false
	}
}

//func (vvm *VVM) pickTokenUsingPointer(pointerType PointerType, pointer int) (Token, error) {
//	switch pointerType {
//	case _BasePointer:
//
//	}
//}

func (vvm *VVM) Execute() []string {
	lineNo := 1

	for !vvm.isProgEof() {
		program := vvm.curtProg()
		if program.op != ILLEGALOpcode {
			howManyOperands := OperandHowManyHas(program.op)
			operands := vvm.operands(howManyOperands)

			// print information
			fmt.Printf("%v) %v ", lineNo, program.lit)
			for _, operand := range operands {
				fmt.Printf("%v ", operand.lit)
			}
			fmt.Println()
			lineNo++

			// exec
			exit, err := vvm.executeOpcode(program.op, operands)
			if err != nil {
				panic(err)
			}

			if exit {
				break
			}
		} else {
			// read comment?
			vvm.movePc(1)
		}
	}

	fmt.Printf("stream:\n")
	fmt.Printf("[ ")
	for i, text := range vvm.stream {
		fmt.Printf("%v", text)
		if len(vvm.stream)-1 != i {
			fmt.Printf(", ")
		}
	}
	fmt.Printf(" ]\n")
	return vvm.stream
}

func (vvm *VVM) executeOpcode(opcode Opcode, args []Token) (exit bool, err error) {
	exit = false

	switch opcode {
	case NOP:
		vvm._nop()
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case SET:
		//vvm._set(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case ADD:
		err = vvm._add(&args[0], &args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case SUB:
		err = vvm._sub(&args[0], &args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case CMP:
		err = vvm._cmp(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case LT:
		err = vvm._lt(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case GT:
		err = vvm._gt(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case LTE:
		err = vvm._lte(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case GTE:
		err = vvm._gte(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case JMP:
		err = vvm._jump(args[0])
	case JZ:
		err = vvm._jz(args[0])
	case JNZ:
		err = vvm._jnz(args[0])

	case CALL:
		err = vvm._call(args[0])
	case RET:
		err = vvm._ret()

	case CP:
		err = vvm._cp(args[0], args[1])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case PUSH:
		err = vvm._push(args[0])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case POP:
		err = vvm._pop(args[0])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case ADDsp:
		err = vvm._addSp(args[0])
		vvm.movePc(1 + OperandHowManyHas(opcode))
	case SUBsp:
		err = vvm._subSp(args[0])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case ECHO:
		err = vvm._echo(args[0])
		vvm.movePc(1 + OperandHowManyHas(opcode))

	case EXIT:
		vvm._exit()
		exit = true

	default:
		err = fmt.Errorf("unknown opcode : %v", opcode.String())
	}

	return exit, err
}

func (vvm *VVM) pickUpRegister(lit string) (*Token, error) {
	switch lit {
	case "reg_a":
		return vvm.regA, nil
	case "reg_b":
		return vvm.regB, nil
	case "reg_c":
		return vvm.regC, nil
	}
	return nil, UnDefinedRegisterErr(lit)
}

func (vvm *VVM) pickUpPointer(lit string) (PointerType, int, error) {
	switch lit {
	case "bp":
		return _BasePointer, vvm.bp, nil
	case "sp":
		return _StackPointer, vvm.sp, nil
	}
	return _IllegalPointer, 0, UndefinedPointerErr(lit)
}

func (vvm *VVM) isSameTokenType(t1 TokenType, t2 TokenType) bool {
	return t1 == t2
}

func (vvm *VVM) addrToPointer(addrLiteral string) (PointerType, int, error) {
	reg := regexp.MustCompile(`\[(bp|sp)([+\-])([0-9]+)]`)
	matches := reg.FindStringSubmatch(addrLiteral)

	//fmt.Printf("`%v`\n", matches)

	diffStr := matches[3]
	diff, _ := strconv.Atoi(diffStr)

	if matches[2] == "-" {
		diff = diff * -1
	}

	switch matches[1] {
	case "bp":
		p := vvm.bp + diff
		if vvm.IsValidPointerLocation(_BasePointer, p) {
			return _BasePointer, p, nil
		} else {
			return _BasePointer, -1, StackAccessErr("bp", len(vvm.stack), p)
		}
	case "sp":
		p := vvm.sp + diff
		if vvm.IsValidPointerLocation(_StackPointer, p) {
			return _StackPointer, vvm.sp + diff, nil
		} else {
			return _StackPointer, -1, StackAccessErr("sp", len(vvm.stack), p)
		}
	}
	return _IllegalPointer, -1, UnexpectedKPointerTypeErr("addrToPointer", []PointerType{_BasePointer, _StackPointer}, _IllegalPointer)
}

func (vvm *VVM) addrToToken(addrLiteral string) (*Token, error) {
	pointerType, pos, err := vvm.addrToPointer(addrLiteral)
	if err != nil {
		return nil, err
	}
	if pointerType == _IllegalPointer {
		return nil, UnexpectedKPointerTypeErr("addrToToken", []PointerType{_BasePointer, _StackPointer}, pointerType)
	}

	return vvm.stack[pos], nil
}

func (vvm *VVM) registerToToken(lit string) (*Token, error) {
	switch lit {
	case "reg_a":
		return vvm.regA, nil
	case "reg_b":
		return vvm.regB, nil
	case "reg_c":
		return vvm.regC, nil
	}
	return nil, UnDefinedRegisterErr(lit)
}

func (vvm *VVM) _nop() {
	// nop :-)
}

//func (vvm *VVM) _set(src Token, dst Token) {
//	// cpでよくね?
//}

func (vvm *VVM) _add(src *Token, dst *Token) error {
	// src: [registers, addr, int, float]
	// dst: [registers, addr] as (int or float)
	var srcTok *Token
	switch src.typ {
	case _REGISTER:
		tok, err := vvm.registerToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = tok
	case _ADDR:
		tok, err := vvm.addrToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = tok
	case _INT, _FLOAT:
		srcTok = src
	}

	switch dst.typ {
	case _REGISTER:
		switch dst.lit {
		case "reg_a":
			if !vvm.isSameTokenType(srcTok.typ, vvm.regA.typ) {
				return DoseNotMatchTokenTypeErr(srcTok.typ, vvm.regA.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vvm.regA.LoadAsInt()
				if err != nil {
					return err
				}
				vvm.regA = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target+diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vvm.regA.LoadAsFloat()
				if err != nil {
					return err
				}
				vvm.regA = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
			}
		case "reg_b":
			if !vvm.isSameTokenType(srcTok.typ, vvm.regB.typ) {
				return DoseNotMatchTokenTypeErr(srcTok.typ, vvm.regB.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vvm.regB.LoadAsInt()
				if err != nil {
					return err
				}
				vvm.regB = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target+diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vvm.regB.LoadAsFloat()
				if err != nil {
					return err
				}
				vvm.regB = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
			}
		case "reg_c":
			if !vvm.isSameTokenType(srcTok.typ, vvm.regC.typ) {
				return DoseNotMatchTokenTypeErr(srcTok.typ, vvm.regC.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vvm.regC.LoadAsInt()
				if err != nil {
					return err
				}
				vvm.regC = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target+diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vvm.regC.LoadAsFloat()
				if err != nil {
					return err
				}
				vvm.regC = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
			}
		}
	case _ADDR:
		_, pointer, err := vvm.addrToPointer(dst.lit)
		if err != nil {
			return err
		}
		if !vvm.isSameTokenType(srcTok.typ, vvm.stack[pointer].typ) {
			return DoseNotMatchTokenTypeErr(srcTok.typ, vvm.stack[pointer].typ)
		}
		switch vvm.stack[pointer].typ {
		case _INT:
			diff, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			target, err := vvm.stack[pointer].LoadAsInt()
			if err != nil {
				return err
			}
			vvm.stack[pointer] = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target+diff))
		case _FLOAT:
			diff, err := srcTok.LoadAsFloat()
			if err != nil {
				return err
			}
			target, err := vvm.stack[pointer].LoadAsFloat()
			if err != nil {
				return err
			}
			vvm.stack[pointer] = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
		}
	}

	return nil
}
func (vvm *VVM) _sub(src *Token, dst *Token) error {
	// src: [registers, pointers] as (int or float)
	// dst: [int, float]
	var srcTok *Token
	switch src.typ {
	case _REGISTER:
		tok, err := vvm.registerToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = tok
	case _ADDR:
		tok, err := vvm.addrToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = tok
	case _INT, _FLOAT:
		srcTok = src
	}

	switch dst.typ {
	case _REGISTER:
		switch dst.lit {
		case "reg_a":
			if !vvm.isSameTokenType(src.typ, vvm.regA.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vvm.regA.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vvm.regA.LoadAsInt()
				if err != nil {
					return err
				}
				vvm.regA = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target-diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vvm.regA.LoadAsFloat()
				if err != nil {
					return err
				}
				vvm.regA = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
			}
		case "reg_b":
			if !vvm.isSameTokenType(src.typ, vvm.regB.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vvm.regB.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vvm.regB.LoadAsInt()
				if err != nil {
					return err
				}
				vvm.regB = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target-diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vvm.regB.LoadAsFloat()
				if err != nil {
					return err
				}
				vvm.regB = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
			}
		case "reg_c":
			if !vvm.isSameTokenType(src.typ, vvm.regC.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vvm.regC.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vvm.regC.LoadAsInt()
				if err != nil {
					return err
				}
				vvm.regC = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target-diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vvm.regC.LoadAsFloat()
				if err != nil {
					return err
				}
				vvm.regC = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
			}
		}
	case _ADDR:
		_, pointer, err := vvm.addrToPointer(dst.lit)
		if err != nil {
			return err
		}
		if !vvm.isSameTokenType(srcTok.typ, vvm.stack[pointer].typ) {
			return DoseNotMatchTokenTypeErr(srcTok.typ, vvm.stack[pointer].typ)
		}
		switch vvm.stack[pointer].typ {
		case _INT:
			diff, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			target, err := vvm.stack[pointer].LoadAsInt()
			if err != nil {
				return err
			}
			vvm.stack[pointer] = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(target-diff))
		case _FLOAT:
			diff, err := srcTok.LoadAsFloat()
			if err != nil {
				return err
			}
			target, err := vvm.stack[pointer].LoadAsFloat()
			if err != nil {
				return err
			}
			vvm.stack[pointer] = NewToken(_FLOAT, ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
		}
	}

	return nil
}
func (vvm *VVM) _cmp(data1 Token, data2 Token) error {
	// data1: [registers, addr, string, int, float]
	// data2: [registers, addr, string, int, float]

	// 型と、litしか見ない。

	var tok1 Token
	var tok2 Token

	switch data1.typ {
	case _REGISTER:
		switch data1.lit {
		case "reg_a":
			tok1 = *vvm.regA
		case "reg_b":
			tok1 = *vvm.regB
		case "reg_c":
			tok1 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data1.lit)
		if err != nil {
			return err
		}
		tok1 = *t
	case _STRING, _INT, _FLOAT:
		tok1 = data1
	}

	switch data2.typ {
	case _REGISTER:
		switch data2.lit {
		case "reg_a":
			tok2 = *vvm.regA
		case "reg_b":
			tok2 = *vvm.regB
		case "reg_c":
			tok2 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data2.lit)
		if err != nil {
			return err
		}
		tok2 = *t
	case _STRING, _INT, _FLOAT:
		tok2 = data2
	}

	// 型おなじ?
	if !vvm.isSameTokenType(tok1.typ, tok2.typ) {
		vvm.zf = 0
	}

	if tok1.lit == tok2.lit {
		vvm.zf = 1
	} else {
		vvm.zf = 0
	}

	return nil
}

func (vvm *VVM) _lt(data1 Token, data2 Token) error {
	// data1: [registers, addr, int, float]
	// data2: [registers, addr, int, float]

	// 型と、litしか見ない。

	var tok1 Token
	var tok2 Token

	switch data1.typ {
	case _REGISTER:
		switch data1.lit {
		case "reg_a":
			tok1 = *vvm.regA
		case "reg_b":
			tok1 = *vvm.regB
		case "reg_c":
			tok1 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data1.lit)
		if err != nil {
			return err
		}
		tok1 = *t
	case _INT, _FLOAT:
		tok1 = data1
	}

	switch data2.typ {
	case _REGISTER:
		switch data2.lit {
		case "reg_a":
			tok2 = *vvm.regA
		case "reg_b":
			tok2 = *vvm.regB
		case "reg_c":
			tok2 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data2.lit)
		if err != nil {
			return err
		}
		tok2 = *t
	case _INT, _FLOAT:
		tok2 = data2
	}

	// 型おなじ?
	if !vvm.isSameTokenType(tok1.typ, tok2.typ) {
		vvm.zf = 0
	}

	if tok1.typ != _INT && tok1.typ != _FLOAT {
		return UnexpectedTokenTypeErr("lt", []TokenType{_INT, _FLOAT}, tok1.typ)
	}

	var result bool

	if tok1.typ == _FLOAT {
		a1, err := tok1.LoadAsFloat()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsFloat()
		if err != nil {
			return err
		}
		result = a1 < a2
	} else {
		a1, err := tok1.LoadAsInt()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsInt()
		if err != nil {
			return err
		}
		result = a1 < a2
	}

	if result {
		vvm.zf = 1
	} else {
		vvm.zf = 0
	}

	return nil
}

func (vvm *VVM) _lte(data1 Token, data2 Token) error {
	// data1: [registers, addr, int, float]
	// data2: [registers, addr, int, float]

	// 型と、litしか見ない。

	var tok1 Token
	var tok2 Token

	switch data1.typ {
	case _REGISTER:
		switch data1.lit {
		case "reg_a":
			tok1 = *vvm.regA
		case "reg_b":
			tok1 = *vvm.regB
		case "reg_c":
			tok1 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data1.lit)
		if err != nil {
			return err
		}
		tok1 = *t
	case _INT, _FLOAT:
		tok1 = data1
	}

	switch data2.typ {
	case _REGISTER:
		switch data2.lit {
		case "reg_a":
			tok2 = *vvm.regA
		case "reg_b":
			tok2 = *vvm.regB
		case "reg_c":
			tok2 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data2.lit)
		if err != nil {
			return err
		}
		tok2 = *t
	case _INT, _FLOAT:
		tok2 = data2
	}

	// 型おなじ?
	if !vvm.isSameTokenType(tok1.typ, tok2.typ) {
		vvm.zf = 0
	}

	if tok1.typ != _INT && tok1.typ != _FLOAT {
		return UnexpectedTokenTypeErr("lt", []TokenType{_INT, _FLOAT}, tok1.typ)
	}

	var result bool

	if tok1.typ == _FLOAT {
		a1, err := tok1.LoadAsFloat()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsFloat()
		if err != nil {
			return err
		}
		result = a1 <= a2
	} else {
		a1, err := tok1.LoadAsInt()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsInt()
		if err != nil {
			return err
		}
		result = a1 <= a2
	}

	if result {
		vvm.zf = 1
	} else {
		vvm.zf = 0
	}

	return nil
}

func (vvm *VVM) _gt(data1 Token, data2 Token) error {
	// data1: [registers, addr, int, float]
	// data2: [registers, addr, int, float]

	// 型と、litしか見ない。

	var tok1 Token
	var tok2 Token

	switch data1.typ {
	case _REGISTER:
		switch data1.lit {
		case "reg_a":
			tok1 = *vvm.regA
		case "reg_b":
			tok1 = *vvm.regB
		case "reg_c":
			tok1 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data1.lit)
		if err != nil {
			return err
		}
		tok1 = *t
	case _INT, _FLOAT:
		tok1 = data1
	}

	switch data2.typ {
	case _REGISTER:
		switch data2.lit {
		case "reg_a":
			tok2 = *vvm.regA
		case "reg_b":
			tok2 = *vvm.regB
		case "reg_c":
			tok2 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data2.lit)
		if err != nil {
			return err
		}
		tok2 = *t
	case _INT, _FLOAT:
		tok2 = data2
	}

	// 型おなじ?
	if !vvm.isSameTokenType(tok1.typ, tok2.typ) {
		vvm.zf = 0
	}

	if tok1.typ != _INT && tok1.typ != _FLOAT {
		return UnexpectedTokenTypeErr("lt", []TokenType{_INT, _FLOAT}, tok1.typ)
	}

	var result bool

	if tok1.typ == _FLOAT {
		a1, err := tok1.LoadAsFloat()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsFloat()
		if err != nil {
			return err
		}
		result = a1 > a2
	} else {
		a1, err := tok1.LoadAsInt()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsInt()
		if err != nil {
			return err
		}
		result = a1 > a2
	}

	if result {
		vvm.zf = 1
	} else {
		vvm.zf = 0
	}

	return nil
}

func (vvm *VVM) _gte(data1 Token, data2 Token) error {
	// data1: [registers, addr, int, float]
	// data2: [registers, addr, int, float]

	// 型と、litしか見ない。

	var tok1 Token
	var tok2 Token

	switch data1.typ {
	case _REGISTER:
		switch data1.lit {
		case "reg_a":
			tok1 = *vvm.regA
		case "reg_b":
			tok1 = *vvm.regB
		case "reg_c":
			tok1 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data1.lit)
		if err != nil {
			return err
		}
		tok1 = *t
	case _INT, _FLOAT:
		tok1 = data1
	}

	switch data2.typ {
	case _REGISTER:
		switch data2.lit {
		case "reg_a":
			tok2 = *vvm.regA
		case "reg_b":
			tok2 = *vvm.regB
		case "reg_c":
			tok2 = *vvm.regC
		}
	case _ADDR:
		t, err := vvm.addrToToken(data2.lit)
		if err != nil {
			return err
		}
		tok2 = *t
	case _INT, _FLOAT:
		tok2 = data2
	}

	// 型おなじ?
	if !vvm.isSameTokenType(tok1.typ, tok2.typ) {
		vvm.zf = 0
	}

	if tok1.typ != _INT && tok1.typ != _FLOAT {
		return UnexpectedTokenTypeErr("lt", []TokenType{_INT, _FLOAT}, tok1.typ)
	}

	var result bool

	if tok1.typ == _FLOAT {
		a1, err := tok1.LoadAsFloat()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsFloat()
		if err != nil {
			return err
		}
		result = a1 >= a2
	} else {
		a1, err := tok1.LoadAsInt()
		if err != nil {
			return err
		}
		a2, err := tok2.LoadAsInt()
		if err != nil {
			return err
		}
		result = a1 >= a2
	}

	if result {
		vvm.zf = 1
	} else {
		vvm.zf = 0
	}

	return nil
}

func (vvm *VVM) _jump(to Token) error {
	newPc, err := to.LoadAsInt()
	if err != nil {
		return err
	}
	vvm.pc = newPc
	return nil
}
func (vvm *VVM) _jz(to Token) error {
	// to: [int] as pc
	if to.typ != _INT {
		return UnexpectedTokenTypeErr("jz", []TokenType{_INT}, to.typ)
	}

	if vvm.zf == 0 {
		newPc, err := to.LoadAsInt()
		if err != nil {
			return err
		}
		vvm.pc = newPc
	} else {
		vvm.movePc(1 + OperandHowManyHas(JZ))
	}

	return nil
}
func (vvm *VVM) _jnz(to Token) error {
	// to: [int] as pc
	if to.typ != _INT {
		return UnexpectedTokenTypeErr("jnz", []TokenType{_INT}, to.typ)
	}

	if vvm.zf != 0 {
		newPc, err := to.LoadAsInt()
		if err != nil {
			return err
		}
		vvm.pc = newPc
	} else {
		vvm.movePc(1 + OperandHowManyHas(JNZ))
	}

	return nil
}

//func (vvm *VVM) _je(a1 Token, a2 Token) error {
//	return nil
//}

func (vvm *VVM) _call(op Token) error {
	// op: [int] as pc
	if err := vvm.subSp(1); err != nil {
		return err
	}
	vvm.stack[vvm.sp] = NewToken(_RTNAddr, ILLEGALOpcode, strconv.Itoa(vvm.pc+2))
	addrWeAreGoing := op
	newPc, err := addrWeAreGoing.LoadAsInt()
	if err != nil {
		return err
	}
	vvm.pc = newPc

	return nil
}
func (vvm *VVM) _ret() error {
	rtnAddr := vvm.stack[vvm.sp]
	if err := vvm.addSp(1); err != nil {
		return err
	}
	if rtnAddr.typ != _RTNAddr {
		return UnexpectedTokenTypeErr("ret", []TokenType{_RTNAddr}, rtnAddr.typ)
	}
	newPc, err := rtnAddr.LoadAsInt()
	if err != nil {
		return err
	}
	vvm.pc = newPc

	return nil
}

func (vvm *VVM) _cp(src Token, dst Token) error {
	// type check
	srcAllowed := []TokenType{_REGISTER, _ADDR, _POINTER, _STRING, _FLOAT, _INT}
	dstAllowed := []TokenType{_REGISTER, _ADDR, _POINTER}
	if !IsAllowedTokenType(srcAllowed, src.typ) {
		return UnexpectedTokenTypeErr("cp.src", srcAllowed, src.typ)
	}
	if !IsAllowedTokenType(dstAllowed, dst.typ) {
		return UnexpectedTokenTypeErr("cp.dst", dstAllowed, dst.typ)
	}

	// コピー元
	var srcTok *Token
	switch src.typ {
	case _REGISTER:
		t, err := vvm.pickUpRegister(src.lit)
		if err != nil {
			return err
		}
		srcTok = t
	case _ADDR:
		t, err := vvm.addrToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = t
	case _STRING, _FLOAT, _INT:
		srcTok = &src
	case _POINTER:
		switch src.lit {
		case "bp":
			srcTok = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(vvm.bp))
		case "sp":
			srcTok = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(vvm.sp))
		}
	}

	// コピー先
	switch dst.typ {
	case _REGISTER:
		switch dst.lit {
		case "reg_a":
			vvm.regA = srcTok
		case "reg_b":
			vvm.regB = srcTok
		case "reg_c":
			vvm.regC = srcTok
		}
	case _ADDR:
		pointerType, pointer, err := vvm.addrToPointer(dst.lit)
		if err != nil {
			return err
		}
		switch pointerType {
		case _BasePointer, _StackPointer:
			vvm.stack[pointer] = srcTok
		}
	case _POINTER:
		if srcTok.typ != _INT {
			return UnexpectedTokenTypeErr("cp[to pointer]", []TokenType{_INT}, srcTok.typ)
		}

		switch dst.lit {
		case "bp":
			newBp, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			vvm.bp = newBp
		case "sp":
			newSp, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			vvm.sp = newSp
		}
	}

	return nil
}

func (vvm *VVM) _push(data Token) error {
	// data: [registers, pointer, addr, string, int, float]
	// addr to raw data

	var dataTok *Token
	switch data.typ {
	case _REGISTER:
		switch data.lit {
		case "reg_a":
			dataTok = vvm.regA
		case "reg_b":
			dataTok = vvm.regB
		case "reg_c":
			dataTok = vvm.regC
		}
	case _ADDR:
		tok, err := vvm.addrToToken(data.lit)
		if err != nil {
			return err
		}
		dataTok = tok
	case _STRING, _INT, _FLOAT:
		dataTok = &data
	case _POINTER:
		switch data.lit {
		case "bp":
			dataTok = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(vvm.bp))
		case "sp":
			dataTok = NewToken(_INT, ILLEGALOpcode, strconv.Itoa(vvm.sp))
		}
	default:
		return UnexpectedTokenTypeErr("push", []TokenType{_REGISTER, _ADDR, _STRING, _INT, _FLOAT}, data.typ)
	}

	if err := vvm.subSp(1); err != nil {
		return err
	}

	vvm.stack[vvm.sp] = dataTok

	return nil
}
func (vvm *VVM) _pop(popTo Token) error {
	// popTo: [registers, addr]
	// addr to pos
	data := vvm.stack[vvm.sp]

	switch popTo.typ {
	case _REGISTER:
		switch popTo.lit {
		case "reg_a":
			vvm.regA = data
		case "reg_b":
			vvm.regB = data
		case "reg_c":
			vvm.regC = data
		}
	case _ADDR:
		pointerType, pointer, err := vvm.addrToPointer(popTo.lit)
		if err != nil {
			return err
		}

		switch pointerType {
		case _BasePointer, _StackPointer:
			vvm.stack[pointer] = data
		default:
			return UnexpectedKPointerTypeErr("pop", []PointerType{_BasePointer, _StackPointer}, pointerType)
		}
	case _POINTER:
		switch popTo.lit {
		case "bp":
			newBp, err := data.LoadAsInt()
			if err != nil {
				return err
			}
			vvm.bp = newBp
		case "sp":
			newSp, err := data.LoadAsInt()
			if err != nil {
				return err
			}
			vvm.sp = newSp
		}
	}

	if err := vvm.addSp(1); err != nil {
		return err
	}

	return nil
}

func (vvm *VVM) _addSp(n Token) error {
	// n: [int]
	if n.typ != _INT {
		return UnexpectedTokenTypeErr("_addSp", []TokenType{_INT}, n.typ)
	}
	diff, err := n.LoadAsInt()
	if err != nil {
		return err
	}
	err = vvm.addSp(diff)
	if err != nil {
		return err
	}

	return nil
}

func (vvm *VVM) _subSp(n Token) error {
	// n: [int]
	if n.typ != _INT {
		return UnexpectedTokenTypeErr("n.typ", []TokenType{_INT}, n.typ)
	}
	diff, err := n.LoadAsInt()
	if err != nil {
		return err
	}
	err = vvm.subSp(diff)
	if err != nil {
		return err
	}

	return nil
}

func (vvm *VVM) _echo(data Token) error {
	// data: [registers, pointer, addr, string, int, float]
	var dataTok Token
	switch data.typ {
	case _REGISTER:
		switch data.lit {
		case "reg_a":
			dataTok = *vvm.regA
		case "reg_b":
			dataTok = *vvm.regB
		case "reg_c":
			dataTok = *vvm.regC
		}
	case _POINTER, _STRING, _INT, _FLOAT:
		dataTok = data
	case _ADDR:
		tok, err := vvm.addrToToken(data.lit)
		if err != nil {
			return err
		}
		dataTok = *tok
	default:
		if data.lit == "zf" {
			dataTok = *NewToken(_ILLEGALToken, ILLEGALOpcode, strconv.Itoa(vvm.zf))
		}
	}

	vvm.writeStream(dataTok.lit)
	return nil
}

// いる？これ？
// もしかしたら、終了前の最終サービスみたいな感じで使えるかもしれないけど
// そんな大掛かりな機能を追加する予定はない
func (vvm *VVM) _exit() {}
