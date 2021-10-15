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
		vm._sub(args[0], args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))
	case _CMP:
		vm._cmp(args[0], args[1])
		vm.movePc(1 + OperandHowManyHas(opcode))

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
		vm._echo(args[0])
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

//func (vm *VM) keywordToToken(keyword string) (*Token, error) {
//	switch CheckKeyWordType(keyword) {
//	case _REGISTER:
//		tok, err := vm.pickUpRegister(keyword)
//		if err != nil {
//			return nil, err
//		}
//		return tok, nil
//	case _POINTER:
//		return nil, UnexpectedKeyWordErr([]KeyWordType{_REGISTER}, _POINTER)
//	default:
//		return nil, UndefinedKeyWordErr(keyword)
//	}
//}

func (vm *VM) _nop() {
	// nop :-)
}

func (vm *VM) _set(src Token, dst Token) {
	// cpでよくね?
}

func (vm *VM) _add(src *Token, dst *Token) error {
	// src: [registers, addr, int, float]
	// dst: [registers, addr] as (int or float)

	//// is keyword or addr?
	//if dst.typ != _KEYWORD && dst.typ != _ADDR {
	//	return UnexpectedTokenTypeErr("dst.typ", []TokenType{_KEYWORD, _ADDR}, dst.typ)
	//}
	//
	//// src
	//var srcToken *Token
	//switch src.typ {
	//case _ADDR:
	//	tok, err := vm.addrToToken(src.lit)
	//	if err != nil {
	//		return err
	//	}
	//	srcToken = tok
	//case _KEYWORD:
	//	tok, err := vm.keywordToToken(src.lit)
	//	if err != nil {
	//		return err
	//	}
	//	srcToken = tok
	//case _INT, _FLOAT:
	//	srcToken = src
	//default:
	//	return UnexpectedTokenTypeErr("src.typ", []TokenType{_ADDR, _KEYWORD, _INT, _FLOAT}, srcToken.typ)
	//}
	//
	//// dst
	//var dstToken *Token
	//if dst.typ == _ADDR {
	//	tok, err := vm.addrToToken(dst.lit)
	//	if err != nil {
	//		return err
	//	}
	//	dstToken = tok
	//} else {
	//	tok, err := vm.keywordToToken(dst.lit)
	//	if err != nil {
	//		return err
	//	}
	//	dstToken = tok
	//}
	//
	//// is float or int?
	//if dstToken.typ != _INT && dstToken.typ != _FLOAT {
	//	return UnexpectedTokenTypeErr("dstToken.typ", []TokenType{_INT, _FLOAT}, dstToken.typ)
	//}
	//
	//// is same type?
	//if !vm.isSameTokenType(srcToken.typ, dstToken.typ) {
	//	return DoseNotMatchTokenTypeErr(srcToken.typ, dstToken.typ)
	//}
	//
	//switch dstToken.typ {
	//case _INT:
	//	target, err := strconv.Atoi(dstToken.lit)
	//	if err != nil {
	//		return err
	//	}
	//
	//	diff, err := strconv.Atoi(srcToken.lit)
	//	if err != nil {
	//		return err
	//	}
	//
	//	dstToken.lit = strconv.Itoa(target + diff)
	//case _FLOAT:
	//	target, err := strconv.ParseFloat(dstToken.lit, 64)
	//	if err != nil {
	//		return err
	//	}
	//
	//	diff, err := strconv.ParseFloat(srcToken.lit, 64)
	//	if err != nil {
	//		return err
	//	}
	//
	//	dstToken.lit = strconv.FormatFloat(target+diff, 'f', -1, 64)
	//}

	return nil
}
func (vm *VM) _sub(src Token, dst Token) {
	// src: [registers, pointers] as (int or float)
	// dst: [int, float]
}
func (vm *VM) _cmp(data1 Token, data2 Token) {
	// data1: [any]
	// data2: [any]
}

func (vm *VM) _jz(to Token) error {
	// to: [int] as pc
	//if to.typ != _INT {
	//	return UnexpectedTokenTypeErr("to.typ", []TokenType{_INT}, to.typ)
	//}
	//
	//if vm.zf == 0 {
	//	newPc, err := to.LoadAsInt()
	//	if err != nil {
	//		return err
	//	}
	//	vm.pc = newPc
	//}
	return nil
}
func (vm *VM) _jnz(to Token) error {
	// to: [int] as pc
	//if to.typ != _INT {
	//	return UnexpectedTokenTypeErr("to.typ", []TokenType{_INT}, to.typ)
	//}
	//
	//if vm.zf != 0 {
	//	newPc, err := to.LoadAsInt()
	//	if err != nil {
	//		return err
	//	}
	//	vm.pc = newPc
	//}

	return nil
}

func (vm *VM) _call(op Token) error {
	// op: [int] as pc
	//vm.sp--
	//vm.stack[vm.sp] = NewToken(_INT, _ILLEGALOpcode, strconv.Itoa(vm.pc+2))
	//addrWeAreGoing := op
	//newPc, err := addrWeAreGoing.LoadAsInt()
	//if err != nil {
	//	return err
	//}
	//vm.pc = newPc

	return nil
}
func (vm *VM) _ret() error {
	//returnAddr := vm.stack[vm.pc+1]
	//vm.sp++
	//newPc, err := returnAddr.LoadAsInt()
	//if err != nil {
	//	return err
	//}
	//vm.pc = newPc

	return nil
}

func (vm *VM) _cp(src Token, dst Token) error {
	// type check
	srcAllowed := []TokenType{_REGISTER, _ADDR, _STRING, _FLOAT, _INT}
	dstAllowed := []TokenType{_REGISTER, _ADDR}
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
	}

	return nil
}

func (vm *VM) _push(data Token) error {
	// data: [registers, addr, string, int, float]
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
