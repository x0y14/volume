package vasm_gen

import (
	"fmt"
	"os"
)

func NewVasmGenWithPath(path string) VasmGen {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	text := string(data)
	return VasmGen{text: text}
}

type VasmGen struct {
	text            string
	nodes           []Node
	definedVariable []string
}

func (vg *VasmGen) Prepare() error {
	tokenizer := NewTokenizer(vg.text)
	tokens, err := tokenizer.Tokenize([]TokenType{WHITESPACE, NEWLINE})
	if err != nil {
		return err
	}

	parser := NewParser(tokens)
	nodes, err := parser.Parse()
	if err != nil {
		return err
	}
	vg.nodes = nodes

	return nil
}

func (vg *VasmGen) GenerateCode() (string, error) {
	code := "call main\nexit\n"

	for _, fRootNode := range vg.nodes {

		code += "\n"

		// set up prepare
		fName := fRootNode.tok.lit
		code += fmt.Sprintf("%v:\n", fName)
		code += "\t; 呼び出し元に戻れるように現状保存\n"
		code += "\tpush bp\n"
		code += "\tcp sp bp\n"

		// function arguments
		fArgs := fRootNode.children[0]
		if fArgs.children != nil {
			code += "\n\t; 引数をローカル変数として保存\n"
		}

		// todo : scope
		// 狭いscopeを使うのであれば、map[string][]stringにして、関数ごとに変数を定義させなければならない。
		// んで、狭いスコープを実装する予定だが、とりあえず、全てグローバルスコープな実装にする。

		argNo := 2
		for _, argument := range fArgs.children {
			vg.definedVariable = append(vg.definedVariable, argument.tok.lit)
			variablePos := len(vg.definedVariable)
			code += fmt.Sprintf("\t; arg: %v\n", argument.tok.lit)
			code += "\tsub 1 sp\n"
			// 引数が「+」に、ローカル変数が「-」にあった気がする。
			code += fmt.Sprintf("\tcp [bp+%v] [bp-%v]\n", argNo, variablePos)
			argNo++
		}

		// function body
		fBody := fRootNode.children[2]
		for _, line := range fBody.children {
			switch line.typ {
			case CallFunc:
				// todo : impl

			case VariableDefine:
				vg.definedVariable = append(vg.definedVariable, line.tok.lit)
				code += fmt.Sprintf("\n\t; 引数準備: %v\n", line.tok.lit)
				code += "\tsub 1 sp\n"

				valPos := len(vg.definedVariable)
				valDataNode := line.children[0]

				if valDataNode.tokTyp != INT && valDataNode.tokTyp != STRING && valDataNode.tokTyp != FLOAT {
					return "", NotYetImplementedErr("GenerateCode", valDataNode.tokTyp.String())
				}

				code += "\t; ローカル変数に値を代入\n"
				code += fmt.Sprintf("\tcp %v [bp-%v]\n", valDataNode.tok.lit, valPos)

			case WhileLoop:
			case Expr:

			}
		}

		// function return value
		//fRetVal := fRootNode.children[1]
		// todo : 戻り値
	}

	return code, nil
}
