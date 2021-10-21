package vbin_gen

import (
	"os"
	"strconv"
)

func NewVBinGen(ops []Operation) *VBinGen {
	return &VBinGen{
		operations: ops,
		labelTable: map[string]int{},
	}
}

type VBinGen struct {
	operations []Operation
	labelTable map[string]int
}

func (vbg *VBinGen) Scan() {
	pc := 0
	for _, operation := range vbg.operations {
		if operation.label != "" {
			vbg.labelTable[operation.label] = pc
		}
		pc += 1 + len(operation.operands)
	}
}

func (vbg *VBinGen) replace() []Operation {
	var ops []Operation

	for _, originalOp := range vbg.operations {

		// 新規作成
		newOp := Operation{}
		// 命令隊は変わらないので、複製
		newOp.opcode = originalOp.opcode

		var newOperands []Operand

		// 引数のなかに、ラベル定義されたものがないかチェック
		for _, originalOperand := range originalOp.operands {
			if pc, ok := vbg.labelTable[originalOperand.lit]; ok {
				// 差し替えたものを挿入
				newOperands = append(newOperands, Operand{lit: strconv.Itoa(pc)})
			} else {
				// そのままを挿入
				newOperands = append(newOperands, originalOperand)
			}
		}
		newOp.operands = newOperands

		ops = append(ops, newOp)
	}

	return ops
}

func (vbg *VBinGen) AsString() string {
	str := ""
	ops := vbg.replace()
	for i, op := range ops {
		str += op.Line()
		if i != len(ops)-1 {
			str += "\n"
		}
	}

	return str
}

func (vbg *VBinGen) AsLine() []string {
	var line []string
	ops := vbg.replace()
	for _, op := range ops {
		line = append(line, op.Line())
	}

	return line
}

func (vbg *VBinGen) Export(path string) error {
	var file *os.File

	if _, err := os.Stat(path); err == nil {
		if err := os.Remove(path); err != nil {
			return err
		}
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		file = f
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		f, err := os.Create(path)
		if err != nil {
			return err
		}
		file = f
	} else {
		return err
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if _, err := file.WriteString(vbg.AsString()); err != nil {
		return err
	}

	return nil
}
