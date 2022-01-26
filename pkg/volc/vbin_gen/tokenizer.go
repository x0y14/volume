package vbin_gen

import (
	"strings"
	"unicode"
)

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		pos:   0,
		runes: []rune(text),
	}
}

type Tokenizer struct {
	pos   int
	runes []rune
}

func (tk *Tokenizer) isEof() bool {
	return tk.pos >= len(tk.runes)
}

// read
func (tk *Tokenizer) prev() rune {
	return tk.runes[tk.pos-1]
}

func (tk *Tokenizer) curt() rune {
	return tk.runes[tk.pos]
}

func (tk *Tokenizer) next() rune {
	return tk.runes[tk.pos+1]
}

// move
func (tk *Tokenizer) goPrev() {
	tk.pos--
}

func (tk *Tokenizer) goNext() {
	tk.pos++
}

// is type of rune
func (tk *Tokenizer) isNewLine(c rune) bool {
	return c == '\n' || c == '\r'
}

func (tk *Tokenizer) isWhitespace(c rune) bool {
	return unicode.IsSpace(c) && c != '\n' && c != '\r'
}

func (tk *Tokenizer) isSymbol(c rune) bool {
	// exclude double quotation("), underscore(_) and semicolon(;).
	return strings.Contains(".!#$%&'()*+,-/:<=>?@[\\]^`{|}~", string(c))
}

func (tk *Tokenizer) isNumber(c rune) bool {
	return unicode.IsNumber(c)
}

// consume
func (tk *Tokenizer) consumeNewLine() *Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if tk.isNewLine(c) {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}

	ePos := tk.pos

	return NewToken(lit, NEWLINE, sPos, ePos)
}

func (tk *Tokenizer) consumeWhitespace() *Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if tk.isWhitespace(c) {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}

	ePos := tk.pos

	return NewToken(lit, WHITESPACE, sPos, ePos)
}

func (tk *Tokenizer) consumeNumeric() (*Token, error) {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if tk.isNumber(c) || c == '.' {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}

	ePos := tk.pos

	var typ TokenType

	// check dot's location
	if strings.Contains(lit, ".") {
		typ = FLOAT
		// head or tail
		if lit[0] == '.' || strings.Index(lit, ".") == len(lit)-1 {
			return nil, InvalidTokenErr("invalid dot location", sPos, ePos)
		}
	} else {
		typ = INT
	}

	return NewToken(lit, typ, sPos, ePos), nil
}

func (tk *Tokenizer) consumeString() *Token {
	lit := ""
	sPos := tk.pos
	quoCount := 0

	for !tk.isEof() {
		c := tk.curt()
		if c == '"' {
			quoCount++
			lit += string(c)
			if quoCount%2 == 0 {
				tk.goNext()
				break
			}
		} else {
			lit += string(c)
		}
		tk.goNext()
	}

	ePos := tk.pos

	return NewToken(lit, STRING, sPos, ePos)
}

func (tk *Tokenizer) consumeComment() *Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if tk.isNewLine(c) {
			break
		} else {
			lit += string(c)
		}
		tk.goNext()
	}

	ePos := tk.pos

	return NewToken(lit, COMMENT, sPos, ePos)
}

func (tk *Tokenizer) consumeSymbol() *Token {
	lit := ""
	sPos := tk.pos

	c := tk.curt()
	lit += string(c)
	tk.goNext()

	ePos := tk.pos

	return NewToken(lit, SYMBOL, sPos, ePos)
}

func (tk *Tokenizer) consumeIdent() *Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if unicode.IsLetter(c) || unicode.IsNumber(c) || c == '_' {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}

	ePos := tk.pos

	return NewToken(lit, IDENT, sPos, ePos)
}

func (tk *Tokenizer) consumeIllegal() *Token {
	lit := ""
	sPos := tk.pos

	c := tk.curt()
	lit += string(c)
	tk.goNext()

	ePos := tk.pos

	return NewToken(lit, ILLEGAL, sPos, ePos)
}

func (tk *Tokenizer) isIgnoreTokenType(ignore []TokenType, actual TokenType) bool {
	for _, ig := range ignore {
		if ig == actual {
			return true
		}
	}
	return false
}

func (tk *Tokenizer) Tokenize(ignore []TokenType) (*[]Token, error) {
	var tokens []Token

	for !tk.isEof() {
		var tok *Token
		c := tk.curt()

		if tk.isNewLine(c) {
			tok = tk.consumeNewLine()
		} else if tk.isWhitespace(c) {
			tok = tk.consumeWhitespace()
		} else if tk.isNumber(c) {
			t, err := tk.consumeNumeric()
			if err != nil {
				return nil, err
			}
			tok = t
		} else if c == '"' {
			tok = tk.consumeString()
		} else if tk.isSymbol(c) {
			tok = tk.consumeSymbol()
		} else if c == ';' {
			tok = tk.consumeComment()
		} else if unicode.IsLetter(c) || c == '_' {
			tok = tk.consumeIdent()
		} else {
			tok = tk.consumeIllegal()
		}

		if !tk.isIgnoreTokenType(ignore, tok.typ) {
			tokens = append(tokens, *tok)
		}
	}

	return &tokens, nil
}
