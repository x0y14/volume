package tokenizer

import (
	"fmt"
	"strings"
	"unicode"
)

func NewTokenizer(script string) Tokenizer {
	return Tokenizer{
		pos:   0,
		runes: []rune(script),
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
		Lit:  lit,
		Typ:  NEWLINE,
		SPos: sPos,
		EPos: ePos,
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
		Lit:  lit,
		Typ:  WHITESPACE,
		SPos: sPos,
		EPos: ePos,
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
			//return nil, InvalidTokenErr("invalid dot location", SPos, EPos)
			return Token{"", ILLEGAL, -1, -1}, fmt.Errorf("invalid dot location")
		}
	} else {
		typ = INT
	}

	return Token{
		Lit:  lit,
		Typ:  typ,
		SPos: sPos,
		EPos: ePos,
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
			lit += string(c) // \"も追加してるよ。
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
		Lit:  lit,
		Typ:  STRING,
		SPos: sPos,
		EPos: ePos,
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
		Lit:  lit,
		Typ:  COMMENT,
		SPos: sPos,
		EPos: ePos,
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
		Lit:  lit,
		Typ:  COMMENT,
		SPos: sPos,
		EPos: ePos,
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
			Lit:  lit,
			Typ:  typ,
			SPos: sPos,
			EPos: ePos,
		}
	}

	return Token{
		Lit:  lit,
		Typ:  ILLEGAL,
		SPos: sPos,
		EPos: ePos,
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
			Lit:  lit,
			Typ:  typ,
			SPos: sPos,
			EPos: ePos,
		}
	}

	return Token{
		Lit:  lit,
		Typ:  ILLEGAL,
		SPos: sPos,
		EPos: ePos,
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
		Lit:  lit,
		Typ:  typ,
		SPos: sPos,
		EPos: ePos,
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
		Lit:  lit,
		Typ:  ILLEGAL,
		SPos: sPos,
		EPos: ePos,
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

		if !tk.isIgnoreTokenType(ignore, tok.Typ) {
			tokens = append(tokens, tok)
		}
	}

	tokens = append(tokens, Token{
		Lit:  "",
		Typ:  EOF,
		SPos: tk.pos,
		EPos: tk.pos + 1,
	})

	return tokens, nil
}
