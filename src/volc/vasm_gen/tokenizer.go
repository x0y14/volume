package vasm_gen

import (
	"fmt"
	"strings"
	"unicode"
)

func NewTokenizer(text string) Tokenizer {
	return Tokenizer{
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

func (tk *Tokenizer) consumeNewline() Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if IsNewline(c) {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}
	ePos := tk.pos

	return Token{
		lit:  lit,
		typ:  NEWLINE,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeWhitespace() Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if IsWhitespace(c) {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}
	ePos := tk.pos

	return Token{
		lit:  lit,
		typ:  WHITESPACE,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeNumeric() (Token, error) {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if unicode.IsDigit(c) || c == '.' {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}
	ePos := tk.pos

	var typ TokenType

	if strings.Contains(lit, ".") {
		typ = FLOAT
		// head or tail
		if lit[0] == '.' || strings.Index(lit, ".") == len(lit)-1 {
			//return nil, InvalidTokenErr("invalid dot location", sPos, ePos)
			return Token{"", ILLEGAL, -1, -1}, fmt.Errorf("invalid dot location")
		}
	} else {
		typ = INT
	}

	return Token{
		lit:  lit,
		typ:  typ,
		sPos: sPos,
		ePos: ePos,
	}, nil
}

func (tk *Tokenizer) consumeString() Token {
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

	return Token{
		lit:  lit,
		typ:  STRING,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeOneLineComment() Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if IsNewline(c) {
			break
		} else {
			lit += string(c)
		}
		tk.goNext()
	}
	ePos := tk.pos

	return Token{
		lit:  lit,
		typ:  COMMENT,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeMultiLineComment() Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		n := tk.next()
		if c == '*' && n == '/' {
			// comment end
			lit += string(c)
			tk.goNext()
			lit += string(tk.curt())
			tk.goNext()
			break
		} else {
			lit += string(c)
		}
		tk.goNext()
	}
	ePos := tk.pos

	return Token{
		lit:  lit,
		typ:  COMMENT,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeOperator() Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if IsOperator(c) {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}
	ePos := tk.pos

	if typ := Operator2Type(lit); typ != ILLEGAL {
		return Token{
			lit:  lit,
			typ:  typ,
			sPos: sPos,
			ePos: ePos,
		}
	}

	return Token{
		lit:  lit,
		typ:  ILLEGAL,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeSymbol() Token {
	lit := ""
	sPos := tk.pos

	c := tk.curt()
	lit += string(c)
	tk.goNext()

	ePos := tk.pos

	if typ := Symbol2Type(lit); typ != ILLEGAL {
		return Token{
			lit:  lit,
			typ:  typ,
			sPos: sPos,
			ePos: ePos,
		}
	}

	return Token{
		lit:  lit,
		typ:  ILLEGAL,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeIdent() Token {
	lit := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if unicode.IsLetter(c) || unicode.IsDigit(c) || c == '_' {
			lit += string(c)
		} else {
			break
		}
		tk.goNext()
	}
	ePos := tk.pos

	typ := Ident2Keyword(lit)

	return Token{
		lit:  lit,
		typ:  typ,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeIllegal() Token {
	lit := ""
	sPos := tk.pos

	c := tk.curt()
	lit += string(c)
	tk.goNext()

	ePos := tk.pos

	return Token{
		lit:  lit,
		typ:  ILLEGAL,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) isIgnoreTokenType(ignore []TokenType, actual TokenType) bool {
	for _, ig := range ignore {
		if ig == actual {
			return true
		}
	}
	return false
}

func (tk *Tokenizer) Tokenize(ignore []TokenType) ([]Token, error) {
	var tokens []Token

	for !tk.isEof() {
		var tok Token
		c := tk.curt()

		if IsNewline(c) {
			tok = tk.consumeNewline()
		} else if IsWhitespace(c) {
			tok = tk.consumeWhitespace()
		} else if unicode.IsDigit(c) {
			t, err := tk.consumeNumeric()
			if err != nil {
				return nil, err
			}
			tok = t
		} else if c == '"' {
			tok = tk.consumeString()
		} else if c == '/' && tk.next() == '/' {
			tok = tk.consumeOneLineComment()
		} else if c == '/' && tk.next() == '*' {
			tok = tk.consumeMultiLineComment()
		} else if unicode.IsLetter(c) || c == '_' {
			tok = tk.consumeIdent()
		} else if IsOperator(c) {
			tok = tk.consumeOperator()
		} else if IsSymbol(c) {
			tok = tk.consumeSymbol()
		} else {
			tok = tk.consumeIllegal()
		}

		if !tk.isIgnoreTokenType(ignore, tok.typ) {
			tokens = append(tokens, tok)
		}
	}

	tokens = append(tokens, Token{
		lit:  "",
		typ:  EOF,
		sPos: tk.pos,
		ePos: tk.pos + 1,
	})

	return tokens, nil
}
