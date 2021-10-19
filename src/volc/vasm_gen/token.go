package vasm_gen

import (
	"fmt"
	"unicode"
)

// symbols : ! " # $ % & ' ( ) * + , - . / : ; < = > ? @ [ \ ] ^ _ ` { | } ~
// letters : [a-zA-Z]
// digit   : [0-9]

// STRING  : letters + digit + _
// [note]  start with letters

// INT     : digit

// FLOAT   : digit + .

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	COMMENT
	NEWLINE
	WHITESPACE

	literal_beg
	IDENT  //  main
	STRING // "hello"
	INT    // 123
	FLOAT  // 123.4
	BOOL
	MAP
	LIST
	TRUE
	FALSE
	NULL
	literal_end

	operator_beg

	PERCENT // %
	AST     // *
	PLUS    // +
	MINUS   // -
	SLASH   // /
	EQUAL   // =
	VERBAR  // |
	EXCL    // !
	QUEST   // ?
	COLON   // :
	AMP     // &

	LT // <
	GT // >

	EQUALEq // ==
	QUESTEq // !=
	PLUSEq  // +=
	MINUSEq // -=
	COLONEq // :=
	LTEq    // <=
	GTEq    // >=
	operator_end

	symbol_beg
	DOLLAR     // $
	NUM        // #
	COMMA      // ,
	PERIOD     // .
	SEMI       // ;
	AT         // @
	BSLASH     // \
	CIRC       // ^
	UNDERSCORE // _
	GRAVE      // `
	TILDE      // ~

	LPAREN // (
	LBRACK // [
	LBRACE // {

	RPAREN // )
	RBRACK // ]
	RBRACE // }

	SQUO // '
	DQUO // "
	symbol_end

	keyword_beg
	VAR

	FUNC
	RETURN

	FOR
	WHILE
	BREAK

	IF
	ELSE
	CASE
	keyword_end
)

var tokenTypes = [...]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	WHITESPACE: "WHITESPACE",
	NEWLINE:    "NEWLINE",

	IDENT:  "IDENT",
	STRING: "string",
	INT:    "int",
	FLOAT:  "float",
	BOOL:   "bool",
	MAP:    "map",
	LIST:   "list",
	TRUE:   "true",
	FALSE:  "false",
	NULL:   "null",

	EXCL:       "!",
	NUM:        "#",
	DOLLAR:     "$",
	PERCENT:    "%",
	AMP:        "&",
	AST:        "*",
	PLUS:       "+",
	COMMA:      ",",
	MINUS:      "-",
	PERIOD:     ".",
	SLASH:      "/",
	COLON:      ":",
	SEMI:       ";",
	EQUAL:      "=",
	QUEST:      "?",
	AT:         "@",
	BSLASH:     "\\",
	CIRC:       "^",
	UNDERSCORE: "_",
	GRAVE:      "`",
	VERBAR:     "|",
	TILDE:      "~",

	LT: "<",
	GT: ">",

	EQUALEq: "==",
	QUESTEq: "!=",
	PLUSEq:  "+=",
	MINUSEq: "-=",
	COLONEq: ":=",
	LTEq:    "<=",
	GTEq:    ">=",

	LPAREN: "(",
	LBRACK: "[",
	LBRACE: "{",

	RPAREN: ")",
	RBRACK: "]",
	RBRACE: "}",

	SQUO: "'",
	DQUO: "\"",

	VAR:    "var",
	FUNC:   "func",
	RETURN: "return",
	FOR:    "for",
	WHILE:  "while",
	BREAK:  "break",
	IF:     "if",
	ELSE:   "else",
	CASE:   "case",
}

func (tokenType TokenType) String() string {
	return tokenTypes[tokenType]
}

var keywords map[string]TokenType

func init() {
	keywords = map[string]TokenType{}
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokenTypes[i]] = i
	}
}

func NewToken(typ TokenType, lit string, s int, e int) *Token {
	return &Token{
		lit:  lit,
		typ:  typ,
		sPos: s,
		ePos: e,
	}
}

type Token struct {
	lit  string
	typ  TokenType
	sPos int
	ePos int
}

func (tok *Token) String() string {
	return fmt.Sprintf("Token @ (%03d <= ... < %03d) {  %15q | %15q  }", tok.sPos, tok.ePos, tokenTypes[tok.typ], tok.lit)
}

func IsWhitespace(c rune) bool {
	return unicode.IsSpace(c) && c != '\n'
}

func IsNewline(c rune) bool {
	return c == '\n'
}

func Ident2Keyword(str string) TokenType {
	if typ, ok := keywords[str]; ok {
		return typ
	}

	for i := literal_beg + 1; i < literal_end; i++ {
		if str == tokenTypes[i] {
			return i
		}
	}

	return IDENT

	//switch str {
	//case "true":
	//	return TRUE
	//case "false":
	//	return FALSE
	//case "null":
	//	return NULL
	//default:
	//	return IDENT
	//}
}

func IsOperator(c rune) bool {
	// exclude '/' for comment
	// exclude '_' for ident
	// original = !"#$%&'()*+,.-/:;<=>?@[\]^_`{|}~

	for i := operator_beg + 1; i < operator_end; i++ {
		if string(c) == tokenTypes[i] {
			return true
		}
	}

	//return strings.Contains("!\"#$%&'()*+,.-:;<=>?@[\\]^`{|}~", string(c))
	return false
}

func Operator2Type(str string) TokenType {
	for _, c := range []rune(str) {
		if !IsOperator(c) {
			return ILLEGAL
		}
	}

	for i := operator_beg + 1; i < operator_end; i++ {
		if tokenTypes[i] == str {
			return i
		}
	}

	return ILLEGAL
}

func IsSymbol(c rune) bool {
	for i := symbol_beg + 1; i < symbol_end; i++ {
		if string(c) == tokenTypes[i] {
			return true
		}
	}
	return false
}

func Symbol2Type(str string) TokenType {
	for _, c := range []rune(str) {
		if !IsSymbol(c) {
			return ILLEGAL
		}
	}

	for i := symbol_beg + 1; i < symbol_end; i++ {
		if tokenTypes[i] == str {
			return i
		}
	}

	return ILLEGAL
}

func IsMold(typ TokenType) bool {
	//STRING // "hello"
	//INT    // 123
	//FLOAT  // 123.4
	//BOOL
	//MAP
	//LIST
	for _, mold := range []TokenType{STRING, INT, FLOAT, BOOL, MAP, LIST} {
		if mold == typ {
			return true
		}
	}
	return false
}
