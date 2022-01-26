package parser

import (
	"fmt"
	"github.com/x0y14/volume/pkg/vcc/compiler/tokenizer"
)

type NodeType int

const (
	ILLEGAL NodeType = iota

	FuncDef        // 関数定義
	FuncFormalArgs // 仮引数
	FuncArg
	FuncRetTypes // 戻り値型定義
	FuncBody     // 関数本体

	Contents // Bodyに格納される。スクリプト本体。
	Line

	CallFunc     // 関数呼び出し
	CallFuncArgs // 関数呼び出し引数

	VarDef   // 変数定義
	VarSubst // 代入
	VarData  // 変数の中身
	RHS      // 右辺

	ExprOperator
	CalcExpr      // 式, +, -, -, *, とか
	CondExpr      //条件分岐用の式, Boolを求める。
	AssignExpr    // += -=
	ControlExpr   // ++, --,
	CondExprGroup // CondExprのグループ。 CondExpr || CondExprとか、CondExpr && CondExprとか。
	// experimental

	Cond // 条件分岐
	If
	Elif
	Else
	IfBody

	Return // 関数戻り値

	WhileLoop // while
	ForLoop   // for
	LoopBody
	Break

	Import
)

var nodes = [...]string{
	ILLEGAL: "ILLEGAL",

	FuncDef:        "FuncDef",
	FuncFormalArgs: "FuncFormalArgs",
	FuncArg:        "FuncArg",
	FuncRetTypes:   "FuncRetTypes",
	FuncBody:       "FuncBody",

	Contents: "Contents",
	Line:     "Line",

	CallFunc:     "CallFunc",
	CallFuncArgs: "CallFuncArgs",

	VarDef:   "VarDef",
	VarSubst: "VarSubst",
	VarData:  "VarData",

	RHS: "RHS",

	ExprOperator:  "ExprOperator",
	AssignExpr:    "AssignExpr",
	ControlExpr:   "ControlExpr",
	CalcExpr:      "CalcExpr",
	CondExpr:      "CondExpr",
	CondExprGroup: "CondExprGroup",

	Cond:   "Cond",
	If:     "If",
	Elif:   "Elif",
	Else:   "Else",
	IfBody: "IfBody",

	Return: "Return",

	WhileLoop: "WhileLoop",
	ForLoop:   "ForLoop",
	LoopBody:  "LoopBody",
	Break:     "Break",

	Import: "Import",
}

//func NewNode(typ NodeType, cNode []Node, cToken []tokenizer.Token) Node {
//	return Node{
//		typ:           typ,
//		childrenNode:  cNode,
//		childrenToken: cToken,
//	}
//}

type Node struct {
	typ           NodeType
	childrenNode  []Node
	childrenToken []tokenizer.Token
}

func _whitespace(n int) string {
	var space string
	for i := 0; i < n; i++ {
		space += " "
	}
	return space
}

func (nod *Node) String() string {
	var str string

	switch nod.typ {
	case Import:
		str = fmt.Sprintf("Import : %v ", nod.childrenToken[0].Lit)
	case FuncDef:
		name := nod.childrenToken[0].Lit
		formalArgs := nod.childrenNode[0]
		retTypes := nod.childrenNode[1]
		body := nod.childrenNode[2]

		str = fmt.Sprintf("FuncDef : %v\n", name)
		str += _whitespace(3) + "FuncFormalArgs :\n"
		for i, arg := range formalArgs.childrenNode {
			str += _whitespace(5) + fmt.Sprintf("(%2d) %v (%v)\n", i, arg.childrenToken[0].Lit, arg.childrenToken[1].Typ.String())
		}
		str += _whitespace(3) + "FuncRetTypes :\n"
		for i, ret := range retTypes.childrenToken {
			str += _whitespace(5) + fmt.Sprintf("(%2d) %v\n", i, ret.Typ.String())
		}
		str += _whitespace(3) + "FuncBody :\n"
		for _, content := range body.childrenNode {
			str += _whitespace(6) + content.String()
		}
	case VarDef:
		name := nod.childrenToken[0]
		data := nod.childrenNode[0].childrenToken[0]
		str = fmt.Sprintf("Variable : %v = %v (%v)\n", name.Lit, data.Lit, data.Typ.String())
	default:
		str = fmt.Sprintf("Node { %v }", nodes[nod.typ])
	}

	return str
}

/* 関数 */

func NewFuncDefNode(ident tokenizer.Token, formalArgs Node, retType Node, body Node) Node {

	return Node{
		typ:           FuncDef,
		childrenToken: []tokenizer.Token{ident},
		childrenNode:  []Node{formalArgs, retType, body},
	}
}

func NewFuncFormalArgsNode(args []Node) (Node, error) {
	// check arg type
	for _, arg := range args {
		if arg.typ != FuncArg {
			return Node{}, SyntaxErr("NewFuncFormalArgsNode> arg.typ", "FuncArg", nodes[arg.typ])
		}
	}

	return Node{
		typ:           FuncFormalArgs,
		childrenToken: nil,
		childrenNode:  args,
	}, nil
}

func NewFuncArgNode(ident tokenizer.Token, typ tokenizer.Token) (Node, error) {
	if ident.Typ != tokenizer.IDENT {
		return Node{}, SyntaxErr("NewFuncArgNode> FuncArg.ident", tokenizer.IDENT.String(), ident.Typ.String())
	}

	if !tokenizer.IsMoldType(typ.Typ) {
		return Node{}, SyntaxErr("NewFuncArgNode> FuncArg.typ", "mold", typ.Typ.String())
	}

	return Node{
		typ:           FuncArg,
		childrenToken: []tokenizer.Token{ident, typ},
		childrenNode:  nil,
	}, nil
}

func NewFuncRetTypeNode(types []tokenizer.Token) Node {
	return Node{
		typ:           FuncRetTypes,
		childrenToken: types,
		childrenNode:  nil,
	}
}

func NewFuncBodyNode(lines []Node) Node {
	return Node{
		typ:           FuncBody,
		childrenToken: nil,
		childrenNode:  lines,
	}
}

func NewImportNode(lib tokenizer.Token) Node {
	return Node{
		typ:           Import,
		childrenToken: []tokenizer.Token{lib},
		childrenNode:  nil,
	}
}

func NewVarDefNode(ident tokenizer.Token, data Node) Node {
	return Node{
		typ:           VarDef,
		childrenToken: []tokenizer.Token{ident},
		childrenNode:  []Node{data},
	}
}

func NewVarDataNode(data tokenizer.Token) Node {
	return Node{
		typ:           VarData,
		childrenToken: []tokenizer.Token{data},
		childrenNode:  nil,
	}
}

func NewContentsNode(contents []Node) Node {
	return Node{
		typ:           Contents,
		childrenToken: nil,
		childrenNode:  contents,
	}
}

func NewControlExprNode(controlOp tokenizer.Token, target tokenizer.Token, diff tokenizer.Token) Node {
	return Node{
		typ:           ControlExpr,
		childrenToken: []tokenizer.Token{controlOp, target, diff},
		childrenNode:  nil,
	}
}

func NewCalcExprNode(calcOp tokenizer.Token, args []tokenizer.Token) Node {
	var tokens []tokenizer.Token
	tokens = append(tokens, calcOp)
	tokens = append(tokens, args...)
	return Node{
		typ:           CalcExpr,
		childrenToken: tokens,
		childrenNode:  nil,
	}
}

func NewCondExprNode(condOp tokenizer.Token, a1 tokenizer.Token, a2 tokenizer.Token) Node {
	return Node{
		typ:           CondExpr,
		childrenToken: []tokenizer.Token{condOp, a1, a2},
		childrenNode:  nil,
	}
}

func NewCondExprGroupNode(logicalOps []tokenizer.Token, experiments []Node) Node {
	return Node{
		typ:           CondExprGroup,
		childrenToken: logicalOps,
		childrenNode:  experiments,
	}
}

func NewRHSNode(tokens []tokenizer.Token) (Node, error) {
	for _, tok := range tokens {
		fmt.Printf(">> %v\n", tok.String())
	}
	return Node{}, nil
}
