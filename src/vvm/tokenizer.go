package vvm

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

func NewTokenizerPath(path string) (*Tokenizer, error) {
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(dat)
	return &Tokenizer{
		raw:   text,
		runes: []rune(text),
		pos:   0,
	}, nil
}

func NewTokenizer(text string) *Tokenizer {
	return &Tokenizer{
		raw:   text,
		runes: []rune(text),
		pos:   0,
	}
}

type Tokenizer struct {
	raw   string
	runes []rune
	pos   int
}

// get rune
func (tk *Tokenizer) prev() rune {
	return tk.runes[tk.pos-1]
}

func (tk *Tokenizer) curt() rune {
	return tk.runes[tk.pos]
}

func (tk *Tokenizer) next() rune {
	return tk.runes[tk.pos+1]
}

// move on runes
func (tk *Tokenizer) goPrev() int {
	tk.pos--
	return tk.pos
}

func (tk *Tokenizer) goNext() int {
	tk.pos++
	return tk.pos
}

func (tk *Tokenizer) isEof() bool {
	return tk.pos >= len(tk.runes)
}

func (tk *Tokenizer) consumeKeyword() *Token {
	literal := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()

		if unicode.IsSpace(c) {
			break
		} else {
			literal += string(c)
		}
		tk.goNext()
	}
	ePos := tk.pos

	if isReservedKeyWord(literal) {
		switch CheckKeyWordType(literal) {
		case _POINTERKeyWord:
			return &Token{
				typ:  _POINTER,
				op:   _ILLEGALOpcode,
				lit:  literal,
				sPos: sPos,
				ePos: ePos,
			}
		case _REGISTERKeyWord:
			return &Token{
				typ:  _REGISTER,
				op:   _ILLEGALOpcode,
				lit:  literal,
				sPos: sPos,
				ePos: ePos,
			}
		default:
			return &Token{
				typ:  _ILLEGALToken,
				op:   _ILLEGALOpcode,
				lit:  literal,
				sPos: sPos,
				ePos: ePos,
			}
		}
	} else {
		return &Token{
			typ:  _OPERATION,
			op:   ConvertOpcode(literal),
			lit:  literal,
			sPos: sPos,
			ePos: ePos,
		}
	}
}

// ConsumeNumeric consume float or int
func (tk *Tokenizer) consumeNumeric() (*Token, error) {
	literal := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()

		if unicode.IsSpace(c) {
			break
		}

		if unicode.IsNumber(c) || c == '.' || c == '-' {
			literal += string(c)
		} else {
			return nil, fmt.Errorf("invalid data @ %v : %v", tk.pos, string(c))
		}
		tk.goNext()
	}
	ePos := tk.pos

	typ := _INT

	if strings.Contains(literal, ".") {
		typ = _FLOAT
		// [0] or [-1] == '.'
		if literal[0] == '.' || strings.Index(literal, ".") == len(literal)-1 {
			return nil, fmt.Errorf("invalid dot location")
		}

		if strings.Count(literal, ".") != 1 {
			return nil, fmt.Errorf("invalid number of dots")
		}
	}

	if strings.Contains(literal, "-") {
		if literal[0] != '-' {
			return nil, fmt.Errorf("invalid minus location")
		}
		if strings.Count(literal, "-") != 1 {
			return nil, fmt.Errorf("invalid number of minus")
		}
	}

	return &Token{
		typ:  typ,
		op:   _ILLEGALOpcode,
		lit:  literal,
		sPos: sPos,
		ePos: ePos,
	}, nil
}

// start with ' or "
func (tk *Tokenizer) consumeString() *Token {
	literal := ""
	sPos := tk.pos

	c := tk.curt()
	quoCount := 0
	var useDoubleQuo bool
	if c == '"' {
		useDoubleQuo = true
	} else {
		useDoubleQuo = false
	}

	for !tk.isEof() {
		c := tk.curt()

		// ダブルクオーテーション発見
		if c == '"' {
			// もしダブルクオーテーションをルートとして使用している
			if useDoubleQuo {
				if tk.prev() == '\\' {
					// エスケープされていたら、内容の一部なので追加
					literal += string(c)
				} else {
					// ルート、偶数だった場合、閉じカッコだと思うので、終了しちゃう。
					quoCount++
					if quoCount%2 == 0 {
						tk.goNext()
						break
					}
				}
			}
		} else if c == '\'' {
			// もしシングルクオーテーションをルートとして使用している
			if !useDoubleQuo {
				if tk.prev() == '\\' {
					// エスケープされていたら、内容の一部なので追加
					literal += string(c)
				} else {
					// ルート、偶数だった場合、閉じカッコだと思うので、終了しちゃう。
					quoCount++
					if quoCount%2 == 0 {
						tk.goNext()
						break
					}
				}
			}
		} else {
			literal += string(c)
		}
		tk.goNext()
	}
	ePos := tk.pos

	return &Token{
		typ:  _STRING,
		op:   _ILLEGALOpcode,
		lit:  literal,
		sPos: sPos,
		ePos: ePos,
	}
}

// start with #, while not newline
func (tk *Tokenizer) consumeComment() *Token {
	literal := ""
	sPos := tk.pos

	for !tk.isEof() {
		c := tk.curt()
		if c == '\n' {
			break
		} else {
			literal += string(c)
		}
		tk.goNext()
	}
	ePos := tk.pos

	return &Token{
		typ:  _COMMENT,
		op:   _ILLEGALOpcode,
		lit:  literal,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeAddr() *Token {
	literal := ""
	sPos := tk.pos
	for !tk.isEof() {
		c := tk.curt()
		if unicode.IsSpace(c) {
			break
		}

		if c == ']' {
			literal += string(c)
			tk.goNext()
			break
		} else {
			literal += string(c)
		}

		tk.goNext()
	}

	ePos := tk.pos

	return &Token{
		typ:  _ADDR,
		op:   _ILLEGALOpcode,
		lit:  literal,
		sPos: sPos,
		ePos: ePos,
	}
}

func (tk *Tokenizer) consumeWhitespace() {
	literal := ""
	for !tk.isEof() {
		c := tk.curt()
		if unicode.IsSpace(c) {
			literal += string(c)
		} else {
			break
		}
		tk.goNext()
	}
}

func (tk *Tokenizer) Tokenize() (*[]Token, error) {
	var tokens []Token

	for !tk.isEof() {
		c := tk.curt()
		if unicode.IsLetter(c) {
			tok := tk.consumeKeyword()
			tokens = append(tokens, *tok)
		} else if unicode.IsNumber(c) || c == '-' {
			tok, err := tk.consumeNumeric()
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, *tok)
		} else if c == '"' || c == '\'' {
			tok := tk.consumeString()
			tokens = append(tokens, *tok)
		} else if c == '#' {
			tok := tk.consumeComment()
			tokens = append(tokens, *tok)
		} else if c == '[' {
			tok := tk.consumeAddr()
			tokens = append(tokens, *tok)
		} else if unicode.IsSpace(c) {
			tk.consumeWhitespace()
		} else {
			return nil, fmt.Errorf("unexpected letter : %v", string(c))
		}
	}

	return &tokens, nil
}
