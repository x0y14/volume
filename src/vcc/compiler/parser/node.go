package parser

import (
	"fmt"
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
)

type NodeType int

const (
	ILLEGAL NodeType = iota

	FuncDef        // 関数定義
	FuncFormalArgs // 仮引数
	FuncArg
	FuncRetTypes // 戻り値型定義
	FuncBody     // 関数本体

	Return // 関数戻り値

	CallFunc     // 関数呼び出し
	CallFuncArgs // 関数呼び出し引数

	VarDef // 変数定義
	Subst  // 代入

	Expr // 式
	CondExpr

	IfBody
	If
	Elif
	Else

	WhileLoop // while
	ForLoop   // for
	LoopBody
	Break
)

var nodes = [...]string{
	ILLEGAL: "ILLEGAL",

	FuncDef:        "FuncDef",
	FuncFormalArgs: "FuncFormalArgs",
	FuncArg:        "FuncArg",
	FuncRetTypes:   "FuncRetTypes",
	FuncBody:       "FuncBody",

	Return: "Return",

	CallFunc:     "CallFunc",
	CallFuncArgs: "CallFuncArgs",

	VarDef: "VarDef",
	Subst:  "Subst",

	Expr:     "Expr",
	CondExpr: "CondExpr",

	IfBody: "IfBody",
	If:     "If",
	Elif:   "Elif",
	Else:   "Else",

	WhileLoop: "WhileLoop",
	ForLoop:   "ForLoop",
	LoopBody:  "LoopBody",
	Break:     "Break",
}

func NewNode(typ NodeType, cNode []Node, cToken []tokenizer.Token) Node {
	return Node{
		typ:           typ,
		childrenNode:  cNode,
		childrenToken: cToken,
	}
}

type Node struct {
	typ           NodeType
	childrenNode  []Node
	childrenToken []tokenizer.Token
}

func (nod *Node) String() string {
	str := ""
	str += fmt.Sprintf("Node{ typ: %v, ", nodes[nod.typ])

	if nod.childrenNode != nil {
		str += "cNode[ "
		for _, n := range nod.childrenNode {
			str += fmt.Sprintf("%v, ", nodes[n.typ])
		}
		str += "], "
	}

	if nod.childrenToken != nil {
		str += "cToken[ "
		for _, tok := range nod.childrenToken {
			str += fmt.Sprintf("%v, ", tok.String())
		}
		str += "], "
	}

	str += "}"
	return str
}

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
