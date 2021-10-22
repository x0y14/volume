package vasm_gen

import (
	"fmt"
	"math/rand"
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
	config          []Node
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
	conf, nodes, err := parser.Parse()
	if err != nil {
		return err
	}
	vg.config = conf
	vg.nodes = nodes

	return nil
}

func (vg *VasmGen) LibNeedForBuild() []string {
	var lib []string
	for _, conf := range vg.config {
		if conf.typ == ImportLib {
			lib = append(lib, conf.tok.lit)
		}
	}
	return lib
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func (vg *VasmGen) GenerateCode() (map[string]string, error) {

	codes := map[string]string{}
	codes["root"] = "call main\nexit\n"

	for _, fRootNode := range vg.nodes {
		// set up prepare
		fName := fRootNode.tok.lit
		code := codes[fName]
		code += fmt.Sprintf("%v:\n", fName)
		code += "\t; == 戻り用 ==\n"
		code += "\tpush bp\n"
		code += "\tcp sp bp\n"
		code += "\t; ===========\n"

		// function arguments
		fArgs := fRootNode.children[0]
		if fArgs.children != nil {
			code += "\n\t; 引数をローカル変数として保存\n"
		}

		// todo : scope
		// 狭いscopeを使うのであれば、definedVariable: map[string][]stringにして、関数ごとに変数を定義させなければならない。
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
					return nil, NotYetImplementedErr("GenerateCode.VariableDefine.valDataNode", valDataNode.tokTyp.String())
				}

				code += "\t; ローカル変数に値を代入\n"
				code += fmt.Sprintf("\tcp %v [bp-%v]\n", valDataNode.tok.lit, valPos)

			case WhileLoop:
				id := randomString(10)
				thisLoopName := fmt.Sprintf("%v_while_loop_%v", fName, id)
				entryPointLabel := fmt.Sprintf("%v_entry", thisLoopName)
				conditionalExprLabel := fmt.Sprintf("%v_conditional_expr", thisLoopName)

				// loop entry_point
				codes[entryPointLabel] = ""
				codes[entryPointLabel] += "\t; == 戻り用 ==\n"
				codes[entryPointLabel] += "\tpush bp\n"
				codes[entryPointLabel] += "\tcp sp bp\n"
				codes[entryPointLabel] += "\t; ===========\n"
				codes[entryPointLabel] += "\t; 条件式へ飛ぶ。\n"
				codes[entryPointLabel] += fmt.Sprintf("\tjump %v\n", conditionalExprLabel)

				// loop expr
				codes[conditionalExprLabel] = ""
				codes[conditionalExprLabel] += "\t; == 関数本体 ==\n"
				opr := line.children[0]
				if !IsAllowedType([]TokenType{LT, GT, QUESTEq, EQUALEq}, opr.tokTyp) {
					return nil, NotYetImplementedErr("GenerateCode.whileLoop.Expr", opr.tokTyp.String())
				}
				//target := line.children[1]

			case Expr:

			}
		}

		// function return value
		//fRetVal := fRootNode.children[1]
		// todo : 戻り値

		code += "\t; == 呼び出し前の状態に復元 ==\n"
		code += "\tcp bp sp\n"
		code += "\tpop bp\n"
		code += "\tret\n"
		code += "\t; ========================\n"
	}

	return codes, nil
}
