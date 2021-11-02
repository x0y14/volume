package parser3

type NodeKind int

// program    = stmt*
// stmt       = expr ";"
//            | "{" stmt* "}"
//            | "for" expr ";" expr ";" expr stmt
//            | "while" logically stmt
//            | "var" stmt
//            | "if" logically stmt ("else if" logically stmt)? ("else" stmt)?
//            | "return" stmt?
//            | "break" ";"
//            | "continue" ";"
//            | "func" ident "(" (ident kind)? | (ident kind ",")* ")" kind* stmt
//            | "import" string
// kind       = int | string | float | bool
// expr       = assign
// assign     = logically ("=" logically)?
// logically  = equality ("||" equality | "&&" equality)
// equality   = relational ("==" equality | "!=" equality)*
// relational = add (">" add | ">=" add | "<" add | "<=" add)*
// add        = mul ("+" mul | "-" mul)*
// mul        = unary ("*" unary | "/" unary)*
// unary      = ("+" | "-")? primary
// primary    = int | string | float | bool | ident | call | "(" expr ")"

const (
	_ NodeKind = iota
	NdEof

	NdProgram

	NdStmtExpr

	NdBlock
	NdFor
	NdWhile
	NdVariable
	NdIf
	NdReturn
	NdBreak
	NdContinue
	NdFunction
	NdImport

	NdExpr
	NdAssign
	NdLogically

	NdEquality
	NdEqual
	NdNotEqual

	NdRelational
	NdLt
	NdLte
	NdGt
	NdGte

	NdAdd
	NdSub
	NdMul
	NdDiv
	NdUnary

	NdPrimary
	NdInt
	NdString
	NdFloat
	NdBool
	NdIdent
	NdCall
)

type Node struct {
	kind NodeKind

	lhs      *Node
	rhs      *Node
	children []Node

	varInt    int
	varFloat  float64
	varString string
	varBool   bool
}

func NewProgramNode(children []Node) Node {
	return Node{
		kind:      NdProgram,
		lhs:       nil,
		rhs:       nil,
		children:  children,
		varInt:    0,
		varFloat:  0,
		varString: "",
		varBool:   false,
	}
}

func NewStmtExpr() Node {
	return Node{
		kind: NdStmtExpr,
	}
}

func NewStmtBlock(stmt []Node) Node {
	return Node{
		kind:     0,
		lhs:      nil,
		rhs:      nil,
		children: stmt,
	}
}

func NewImportNode(lib string) Node {
	return Node{
		kind:      NdImport,
		varString: lib,
	}
}

func NewNode(kind NodeKind, lhs Node, rhs Node) Node {
	return Node{
		kind: kind,
		lhs:  &lhs,
		rhs:  &rhs,
	}
}

// primary

func NewIntNode(val int) Node {
	return Node{
		kind:   NdInt,
		varInt: val,
	}
}

func NewFloatNode(val float64) Node {
	return Node{
		kind:     NdFloat,
		varFloat: val,
	}
}

func NewStringNode(val string) Node {
	return Node{
		kind:      NdString,
		varString: val,
	}
}

func NewBoolNode(val bool) Node {
	return Node{
		kind:    NdBool,
		varBool: val,
	}
}

func NewIdentNode(val string) Node {
	return Node{
		kind:      NdIdent,
		varString: val,
	}
}
