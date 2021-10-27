package parser2

import (
	"fmt"
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
)

type NodeKind int

const (
	Illegal NodeKind = iota

	DString
	DInt
	DFloat
	DBool
	DNull

	Import // root

	FuncDef // root
	FuncParams
	FuncParamsItem
	FuncReturnValues

	Scripts

	ExprItem
	ConditionalExpr

	If
	Elseif
	Else

	For
	While

	CallFunc
	CallFuncArgs
	CallFuncArgsItem

	VarDef
	VarRHS
	VarRHSItem
)

type Node struct {
	kind NodeKind
	tok  *tokenizer.Token
	nods []Node
}

// DataTypes - 型
// identは、仮引数とか名前をつけるときに。
// nodsは、データを実際に入れるときに使う？

func NewDStringNode(name *tokenizer.Token) Node {
	return Node{
		kind: DString,
		tok:  name,
		nods: nil,
	}
}

func NewDIntNode(name *tokenizer.Token) Node {
	return Node{
		kind: DInt,
		tok:  name,
		nods: nil,
	}
}

func NewDFloatNode(name *tokenizer.Token) Node {
	return Node{
		kind: DFloat,
		tok:  name,
		nods: nil,
	}
}

func NewDBoolNode(name *tokenizer.Token) Node {
	return Node{
		kind: DBool,
		tok:  name,
		nods: nil,
	}
}

func NewDNullNode(name *tokenizer.Token) Node {
	return Node{
		kind: DNull,
		tok:  name,
		nods: nil,
	}
}

// Import - ライブラリインポート系

func NewImportNode(lib *tokenizer.Token) Node {
	return Node{
		kind: Import,
		tok:  lib,
		nods: nil,
	}
}

// Function - 関数関連

func NewFuncDefNode(name *tokenizer.Token, params Node, returnValues Node, scripts Node) Node {
	// todo
	return Node{
		kind: FuncDef,
		tok:  name,
		nods: []Node{params, returnValues, scripts},
	}
}

func NewFuncParamsNode(params []Node) Node {
	return Node{
		kind: FuncParams,
		tok:  nil,
		nods: params,
	}
}

func NewFuncParamsItemNode(name *tokenizer.Token, data Node) Node {
	return Node{
		kind: FuncParamsItem,
		tok:  name,
		nods: []Node{data},
	}
}

func NewFuncReturnValuesNode(kinds []Node) Node {
	// D_STRING | D_INT | D_FLOAT | D_BOOL | D_NULL
	return Node{
		kind: FuncReturnValues,
		tok:  nil,
		nods: kinds,
	}
}

// Scripts - １行１行のプログラムの集合体

func NewScriptsNode(scripts []Node) Node {
	return Node{
		kind: Scripts,
		tok:  nil,
		nods: scripts,
	}
}

// Expression - 式関連

func NewExprItemNode(item *tokenizer.Token) Node {
	return Node{
		kind: ExprItem,
		tok:  item,
		nods: nil,
	}
}

func NewConditionalExprNode(items []Node) Node {
	return Node{
		kind: ConditionalExpr,
		tok:  nil,
		nods: items,
	}
}

// Conditional Branch - 条件分岐

func NewIfNode(expr Node, scripts Node) Node {
	return Node{
		kind: If,
		tok:  nil,
		nods: []Node{expr, scripts},
	}
}

func NewElseifBNode(expr Node, scripts Node) Node {
	return Node{
		kind: Elseif,
		tok:  nil,
		nods: []Node{expr, scripts},
	}
}

func NewElseNode(scripts Node) Node {
	return Node{
		kind: Else,
		tok:  nil,
		nods: []Node{scripts},
	}
}

// Loop - ループ

func NewForNode(expr Node, scripts Node) Node {
	return Node{
		kind: For,
		tok:  nil,
		nods: []Node{expr, scripts},
	}
}

func NewWhileNode(expr Node, scripts Node) Node {
	return Node{
		kind: While,
		tok:  nil,
		nods: []Node{expr, scripts},
	}
}

func NewCallFuncNode(funcName *tokenizer.Token, args []Node) Node {
	return Node{
		kind: CallFunc,
		tok:  funcName,
		nods: args,
	}
}

func NewCallFuncArgsNode(items []Node) Node {
	return Node{
		kind: CallFuncArgs,
		tok:  nil,
		nods: items,
	}
}

func NewCallFuncItemNode(item *tokenizer.Token) Node {
	return Node{
		kind: CallFuncArgsItem,
		tok:  item,
		nods: nil,
	}
}

// 変数

func NewVarDefNode(name *tokenizer.Token, rhs Node) Node {
	return Node{
		kind: VarDef,
		tok:  nil,
		nods: []Node{rhs},
	}
}
func NewVarRHSNode(items []Node) Node {
	return Node{
		kind: VarRHS,
		tok:  nil,
		nods: items,
	}
}
func NewVarRHSItemNode(item *tokenizer.Token) Node {
	return Node{
		kind: VarRHSItem,
		tok:  item,
		nods: nil,
	}
}

func NewDataTypeNode(data *tokenizer.Token) (Node, error) {
	var nKind NodeKind
	switch data.Typ {
	case tokenizer.STRING:
		nKind = DString
	case tokenizer.INT:
		nKind = DInt
	case tokenizer.FLOAT:
		nKind = DFloat
	case tokenizer.TRUE, tokenizer.FALSE:
		nKind = DBool
	case tokenizer.NULL:
		nKind = DNull
	default:
		return Node{}, SyntaxErr("consumeFuncParam", "DataTypes", data.Typ.String())
	}
	return Node{
		kind: nKind,
		tok:  data,
		nods: nil,
	}, nil
}

func (nod Node) String() string {
	switch nod.kind {
	case DString, DInt, DFloat, DBool, DNull:
		return fmt.Sprintf("")
	}
	return ""
}
