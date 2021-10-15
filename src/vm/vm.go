package vm

import (
	"fmt"
	"regexp"
	"strconv"
)

func NewVM() *VM {
	vm := &VM{
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

type VM struct {
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

func (vm *VM) loadAsm(path string) {
	tokenizer, err := NewTokenizerPath(path)
	if err != nil {
		panic(err)
	}
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		panic(err)
	}
	vm.mainMemory = *tokens
}

func (vm *VM) reserveStack(n int) {
	var stack []*Token

	for i := 0; i < n; i++ {
		stack = append(stack, nil)
	}
	vm.stack = stack
}

func (vm *VM) setPointer() {
	vm.sp = len(vm.stack) - 1
	vm.bp = len(vm.stack) - 1
}

func (vm *VM) SetUp(stackSize int, programPath string) {
	vm.reserveStack(stackSize)
	vm.loadAsm(programPath)
	vm.setPointer()
}

func (vm *VM) writeStream(text string) {
	vm.stream = append(vm.stream, text)
}

func (vm *VM) isProgEof() bool {
	return vm.pc >= len(vm.mainMemory)
}

func (vm *VM) movePc(n int) {
	vm.pc += n
}

func (vm *VM) curtProg() Token {
	return vm.mainMemory[vm.pc]
}

func (vm *VM) operands(n int) []Token {
	var tokens []Token
	for i := 0; i < n; i++ {
		tokens = append(tokens, vm.mainMemory[vm.pc+i+1])
	}
	return tokens
}

func (vm *VM) addSp(diff int) error {
	if (vm.sp+diff) < 0 || len(vm.stack)-1 < (vm.sp+diff) {
		return StackAccessErr("addSp", len(vm.stack)-1, vm.sp+diff)
	}
	vm.sp += diff
	return nil
}
func (vm *VM) subSp(diff int) error {
	// [0, 1, 2, 3] : len() => 4
	if (vm.sp-diff) < 0 || len(vm.stack)-1 < (vm.sp-diff) {
		return StackAccessErr("subSp", len(vm.stack)-1, vm.sp-diff)
	}
	vm.sp -= diff
	return nil
}

func (vm *VM) IsValidPointerLocation(pointerType PointerType, pointer int) bool {
	switch pointerType {
	case _BasePointer, _StackPointer:
		return 0 <= pointer && pointer <= len(vm.stack)-1
	default:
		return false
	}
}

//func (vm *VM) pickTokenUsingPointer(pointerType PointerType, pointer int) (Token, error) {
//	switch pointerType {
//	case _BasePointer:
//
//	}
//}

func (vm *VM) Execute() {
	lineNo := 1

	for !vm.isProgEof() {
		program := vm.curtProg()
		if program.op != _ILLEGALOpcode {
			howManyOperands := OperandHowManyHas(program.op)
			operands := vm.operands(howManyOperands)

			// print information
			fmt.Printf("%v) %v ", lineNo, program.lit)
			for _, operand := range operands {
				fmt.Printf("%v ", operand.lit)
			}
			fmt.Println()
			lineNo++

			// exec
			exit, err := vm.executeOpcode(program.op, operands)
			if err != nil {
				panic(err)
			}

			if exit {
				break
			}
		} else {
			// read comment?
			vm.movePc(1)
		}
	}
}

func (vm *VM) executeOpcode(opcode Opcode, args []Token) (exit bool, err error) {
	exit = false

	switch opcode {
	case _NOP:
		vm._nop()
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _SET:
		vm._set(args[0], args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _ADD:
		err = vm._add(&args[0], &args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))
	case _SUB:
		err = vm._sub(&args[0], &args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))
	case _CMP:
		err = vm._cmp(args[0], args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _JUMP:
		err = vm._jump(args[0])
	case _JZ:
		err = vm._jz(args[0])
	case _JNZ:
		err = vm._jnz(args[0])

	case _CALL:
		err = vm._call(args[0])
	case _RET:
		err = vm._ret()

	case _CP:
		err = vm._cp(args[0], args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _PUSH:
		err = vm._push(args[0])
		vm.movePc(1 + OperandHowManyHas(opcode))
	case _POP:
		err = vm._pop(args[0])
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _ADDsp:
		err = vm._addSp(args[0])
		vm.movePc(1 + OperandHowManyHas(opcode))
	case _SUBsp:
		err = vm._subSp(args[0])
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _ECHO:
		err = vm._echo(args[0])
		vm.movePc(1 + OperandHowManyHas(opcode))

	case _EXIT:
		vm._exit()
		exit = true

	default:
		err = fmt.Errorf("unknown opcode : %v", opcode.String())
	}

	return exit, err
}

func (vm *VM) pickUpRegister(lit string) (*Token, error) {
	switch lit {
	case "reg_a":
		return vm.regA, nil
	case "reg_b":
		return vm.regB, nil
	case "reg_c":
		return vm.regC, nil
	}
	return nil, UnDefinedRegisterErr(lit)
}

func (vm *VM) pickUpPointer(lit string) (PointerType, int, error) {
	switch lit {
	case "bp":
		return _BasePointer, vm.bp, nil
	case "sp":
		return _StackPointer, vm.sp, nil
	}
	return _IllegalPointer, 0, UndefinedPointerErr(lit)
}

func (vm *VM) isSameTokenType(t1 TokenType, t2 TokenType) bool {
	return t1 == t2
}

func (vm *VM) addrToPointer(addrLiteral string) (PointerType, int, error) {
	reg := regexp.MustCompile(`\[(bp|sp)([+\-])([0-9]+)\]`)
	matches := reg.FindStringSubmatch(addrLiteral)

	//fmt.Printf("`%v`\n", matches)

	diffStr := matches[3]
	diff, _ := strconv.Atoi(diffStr)

	if matches[2] == "-" {
		diff = diff * -1
	}

	switch matches[1] {
	case "bp":
		p := vm.bp + diff
		if vm.IsValidPointerLocation(_BasePointer, p) {
			return _BasePointer, p, nil
		} else {
			return _BasePointer, -1, StackAccessErr("bp", len(vm.stack), p)
		}
	case "sp":
		p := vm.sp + diff
		if vm.IsValidPointerLocation(_StackPointer, p) {
			return _StackPointer, vm.sp + diff, nil
		} else {
			return _StackPointer, -1, StackAccessErr("sp", len(vm.stack), p)
		}
	}
	return _IllegalPointer, -1, UnexpectedKPointerTypeErr("addrToPointer", []PointerType{_BasePointer, _StackPointer}, _IllegalPointer)
}

func (vm *VM) addrToToken(addrLiteral string) (*Token, error) {
	pointerType, pos, err := vm.addrToPointer(addrLiteral)
	if err != nil {
		return nil, err
	}
	if pointerType == _IllegalPointer {
		return nil, UnexpectedKPointerTypeErr("addrToToken", []PointerType{_BasePointer, _StackPointer}, pointerType)
	}

	return vm.stack[pos], nil
}

func (vm *VM) registerToToken(lit string) (*Token, error) {
	switch lit {
	case "reg_a":
		return vm.regA, nil
	case "reg_b":
		return vm.regB, nil
	case "reg_c":
		return vm.regC, nil
	}
	return nil, UnDefinedRegisterErr(lit)
}

func (vm *VM) _nop() {
	// nop :-)
}

func (vm *VM) _set(src Token, dst Token) {
	// cpでよくね?
}

func (vm *VM) _add(src *Token, dst *Token) error {
	// src: [registers, addr, int, float]
	// dst: [registers, addr] as (int or float)
	var srcTok *Token
	switch src.typ {
	case _REGISTER:
		tok, err := vm.registerToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = tok
	case _ADDR:
		tok, err := vm.addrToToken(src.lit)
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
			if !vm.isSameTokenType(src.typ, vm.regA.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vm.regA.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vm.regA.LoadAsInt()
				if err != nil {
					return err
				}
				vm.regA = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target+diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vm.regA.LoadAsFloat()
				if err != nil {
					return err
				}
				vm.regA = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
			}
		case "reg_b":
			if !vm.isSameTokenType(src.typ, vm.regB.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vm.regB.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vm.regB.LoadAsInt()
				if err != nil {
					return err
				}
				vm.regB = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target+diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vm.regB.LoadAsFloat()
				if err != nil {
					return err
				}
				vm.regB = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
			}
		case "reg_c":
			if !vm.isSameTokenType(src.typ, vm.regC.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vm.regC.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vm.regC.LoadAsInt()
				if err != nil {
					return err
				}
				vm.regC = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target+diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vm.regC.LoadAsFloat()
				if err != nil {
					return err
				}
				vm.regC = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
			}
		}
	case _ADDR:
		_, pointer, err := vm.addrToPointer(dst.lit)
		if err != nil {
			return err
		}
		if !vm.isSameTokenType(srcTok.typ, vm.stack[pointer].typ) {
			return DoseNotMatchTokenTypeErr(srcTok.typ, vm.stack[pointer].typ)
		}
		switch vm.stack[pointer].typ {
		case _INT:
			diff, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			target, err := vm.stack[pointer].LoadAsInt()
			if err != nil {
				return err
			}
			vm.stack[pointer] = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target+diff))
		case _FLOAT:
			diff, err := srcTok.LoadAsFloat()
			if err != nil {
				return err
			}
			target, err := vm.stack[pointer].LoadAsFloat()
			if err != nil {
				return err
			}
			vm.stack[pointer] = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target+diff, 'f', -1, 64))
		}
	}

	return nil
}
func (vm *VM) _sub(src *Token, dst *Token) error {
	// src: [registers, pointers] as (int or float)
	// dst: [int, float]
	var srcTok *Token
	switch src.typ {
	case _REGISTER:
		tok, err := vm.registerToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = tok
	case _ADDR:
		tok, err := vm.addrToToken(src.lit)
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
			if !vm.isSameTokenType(src.typ, vm.regA.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vm.regA.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vm.regA.LoadAsInt()
				if err != nil {
					return err
				}
				vm.regA = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target-diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vm.regA.LoadAsFloat()
				if err != nil {
					return err
				}
				vm.regA = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
			}
		case "reg_b":
			if !vm.isSameTokenType(src.typ, vm.regB.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vm.regB.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vm.regB.LoadAsInt()
				if err != nil {
					return err
				}
				vm.regB = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target-diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vm.regB.LoadAsFloat()
				if err != nil {
					return err
				}
				vm.regB = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
			}
		case "reg_c":
			if !vm.isSameTokenType(src.typ, vm.regC.typ) {
				return DoseNotMatchTokenTypeErr(src.typ, vm.regC.typ)
			}
			switch srcTok.typ {
			case _INT:
				diff, err := srcTok.LoadAsInt()
				if err != nil {
					return err
				}
				target, err := vm.regC.LoadAsInt()
				if err != nil {
					return err
				}
				vm.regC = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target-diff))
			case _FLOAT:
				diff, err := srcTok.LoadAsFloat()
				if err != nil {
					return err
				}
				target, err := vm.regC.LoadAsFloat()
				if err != nil {
					return err
				}
				vm.regC = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
			}
		}
	case _ADDR:
		_, pointer, err := vm.addrToPointer(dst.lit)
		if err != nil {
			return err
		}
		if !vm.isSameTokenType(srcTok.typ, vm.stack[pointer].typ) {
			return DoseNotMatchTokenTypeErr(srcTok.typ, vm.stack[pointer].typ)
		}
		switch vm.stack[pointer].typ {
		case _INT:
			diff, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			target, err := vm.stack[pointer].LoadAsInt()
			if err != nil {
				return err
			}
			vm.stack[pointer] = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(target-diff))
		case _FLOAT:
			diff, err := srcTok.LoadAsFloat()
			if err != nil {
				return err
			}
			target, err := vm.stack[pointer].LoadAsFloat()
			if err != nil {
				return err
			}
			vm.stack[pointer] = NewToken(_FLOAT, _ILLEGALOpcode, strconv.FormatFloat(target-diff, 'f', -1, 64))
		}
	}

	return nil
}
func (vm *VM) _cmp(data1 Token, data2 Token) error {
	// data1: [registers, addr, string, int, float]
	// data2: [registers, addr, string, int, float]

	// 型と、litしか見ない。

	var tok1 Token
	var tok2 Token

	switch data1.typ {
	case _REGISTER:
		switch data1.lit {
		case "reg_a":
			tok1 = *vm.regA
		case "reg_b":
			tok1 = *vm.regB
		case "reg_c":
			tok1 = *vm.regC
		}
	case _ADDR:
		t, err := vm.addrToToken(data1.lit)
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
			tok2 = *vm.regA
		case "reg_b":
			tok2 = *vm.regB
		case "reg_c":
			tok2 = *vm.regC
		}
	case _ADDR:
		t, err := vm.addrToToken(data2.lit)
		if err != nil {
			return err
		}
		tok2 = *t
	case _STRING, _INT, _FLOAT:
		tok2 = data2
	}

	// 型おなじ?
	if !vm.isSameTokenType(tok1.typ, tok2.typ) {
		vm.zf = 0
	}

	if tok1.lit == tok2.lit {
		vm.zf = 1
	} else {
		vm.zf = 0
	}

	return nil
}

func (vm *VM) _jump(to Token) error {
	newPc, err := to.LoadAsInt()
	if err != nil {
		return err
	}
	vm.pc = newPc
	return nil
}
func (vm *VM) _jz(to Token) error {
	// to: [int] as pc
	if to.typ != _INT {
		return UnexpectedTokenTypeErr("jz", []TokenType{_INT}, to.typ)
	}

	if vm.zf == 0 {
		newPc, err := to.LoadAsInt()
		if err != nil {
			return err
		}
		vm.pc = newPc
	} else {
		vm.movePc(1 + OperandHowManyHas(_JZ))
	}

	return nil
}
func (vm *VM) _jnz(to Token) error {
	// to: [int] as pc
	if to.typ != _INT {
		return UnexpectedTokenTypeErr("jnz", []TokenType{_INT}, to.typ)
	}

	if vm.zf != 0 {
		newPc, err := to.LoadAsInt()
		if err != nil {
			return err
		}
		vm.pc = newPc
	} else {
		vm.movePc(1 + OperandHowManyHas(_JNZ))
	}

	return nil
}

func (vm *VM) _call(op Token) error {
	// op: [int] as pc
	if err := vm.subSp(1); err != nil {
		return err
	}
	vm.stack[vm.sp] = NewToken(_RTNAddr, _ILLEGALOpcode, strconv.Itoa(vm.pc+2))
	addrWeAreGoing := op
	newPc, err := addrWeAreGoing.LoadAsInt()
	if err != nil {
		return err
	}
	vm.pc = newPc

	return nil
}
func (vm *VM) _ret() error {
	rtnAddr := vm.stack[vm.sp]
	if err := vm.addSp(1); err != nil {
		return err
	}
	if rtnAddr.typ != _RTNAddr {
		return UnexpectedTokenTypeErr("ret", []TokenType{_RTNAddr}, rtnAddr.typ)
	}
	newPc, err := rtnAddr.LoadAsInt()
	if err != nil {
		return err
	}
	vm.pc = newPc

	return nil
}

func (vm *VM) _cp(src Token, dst Token) error {
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
		t, err := vm.pickUpRegister(src.lit)
		if err != nil {
			return err
		}
		srcTok = t
	case _ADDR:
		t, err := vm.addrToToken(src.lit)
		if err != nil {
			return err
		}
		srcTok = t
	case _STRING, _FLOAT, _INT:
		srcTok = &src
	case _POINTER:
		switch src.lit {
		case "bp":
			srcTok = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(vm.bp))
		case "sp":
			srcTok = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(vm.sp))
		}
	}

	// コピー先
	switch dst.typ {
	case _REGISTER:
		switch dst.lit {
		case "reg_a":
			vm.regA = srcTok
		case "reg_b":
			vm.regB = srcTok
		case "reg_c":
			vm.regC = srcTok
		}
	case _ADDR:
		pointerType, pointer, err := vm.addrToPointer(dst.lit)
		if err != nil {
			return err
		}
		switch pointerType {
		case _BasePointer, _StackPointer:
			vm.stack[pointer] = srcTok
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
			vm.bp = newBp
		case "sp":
			newSp, err := srcTok.LoadAsInt()
			if err != nil {
				return err
			}
			vm.sp = newSp
		}
	}

	return nil
}

func (vm *VM) _push(data Token) error {
	// data: [registers, pointer, addr, string, int, float]
	// addr to raw data

	var dataTok *Token
	switch data.typ {
	case _REGISTER:
		switch data.lit {
		case "reg_a":
			dataTok = vm.regA
		case "reg_b":
			dataTok = vm.regB
		case "reg_c":
			dataTok = vm.regC
		}
	case _ADDR:
		tok, err := vm.addrToToken(data.lit)
		if err != nil {
			return err
		}
		dataTok = tok
	case _STRING, _INT, _FLOAT:
		dataTok = &data
	case _POINTER:
		switch data.lit {
		case "bp":
			dataTok = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(vm.bp))
		case "sp":
			dataTok = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(vm.sp))
		}
	default:
		return UnexpectedTokenTypeErr("push", []TokenType{_REGISTER, _ADDR, _STRING, _INT, _FLOAT}, data.typ)
	}

	if err := vm.subSp(1); err != nil {
		return err
	}

	vm.stack[vm.sp] = dataTok

	return nil
}
func (vm *VM) _pop(popTo Token) error {
	// popTo: [registers, addr]
	// addr to pos
	data := vm.stack[vm.sp]

	switch popTo.typ {
	case _REGISTER:
		switch popTo.lit {
		case "reg_a":
			vm.regA = data
		case "reg_b":
			vm.regB = data
		case "reg_c":
			vm.regC = data
		}
	case _ADDR:
		pointerType, pointer, err := vm.addrToPointer(popTo.lit)
		if err != nil {
			return err
		}

		switch pointerType {
		case _BasePointer, _StackPointer:
			vm.stack[pointer] = data
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
			vm.bp = newBp
		case "sp":
			newSp, err := data.LoadAsInt()
			if err != nil {
				return err
			}
			vm.sp = newSp
		}
	}

	if err := vm.addSp(1); err != nil {
		return err
	}

	return nil
}

func (vm *VM) _addSp(n Token) error {
	// n: [int]
	if n.typ != _INT {
		return UnexpectedTokenTypeErr("_addSp", []TokenType{_INT}, n.typ)
	}
	diff, err := n.LoadAsInt()
	if err != nil {
		return err
	}
	err = vm.addSp(diff)
	if err != nil {
		return err
	}

	return nil
}

func (vm *VM) _subSp(n Token) error {
	// n: [int]
	if n.typ != _INT {
		return UnexpectedTokenTypeErr("n.typ", []TokenType{_INT}, n.typ)
	}
	diff, err := n.LoadAsInt()
	if err != nil {
		return err
	}
	err = vm.subSp(diff)
	if err != nil {
		return err
	}

	return nil
}

func (vm *VM) _echo(data Token) error {
	// data: [registers, pointer, addr, string, int, float]
	var dataTok Token
	switch data.typ {
	case _REGISTER:
		switch data.lit {
		case "reg_a":
			dataTok = *vm.regA
		case "reg_b":
			dataTok = *vm.regB
		case "reg_c":
			dataTok = *vm.regC
		}
	case _POINTER, _STRING, _INT, _FLOAT:
		dataTok = data
	case _ADDR:
		tok, err := vm.addrToToken(data.lit)
		if err != nil {
			return err
		}
		dataTok = *tok
	}

	vm.writeStream(dataTok.lit)
	return nil
}

// いる？これ？
// もしかしたら、終了前の最終サービスみたいな感じで使えるかもしれないけど
// そんな大掛かりな機能を追加する予定はない
func (vm *VM) _exit() {}
