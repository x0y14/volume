package parser3

import (
	"fmt"
	"github.com/x0y14/volume/pkg/cc/compiler/misc"
	"github.com/x0y14/volume/pkg/cc/compiler/tokenizer"
	"testing"
)

func TestParser_Expr(t *testing.T) {
	var tests = []struct {
		title string
		in    string
	}{
		{
			"a",
			"ident = 6",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			tk := tokenizer.NewTokenizer(test.in)
			tokens, err := tk.Tokenize([]tokenizer.TokenType{tokenizer.WHITESPACE, tokenizer.COMMENT, tokenizer.NEWLINE})
			if err != nil {
				t.Fatal(err)
			}
			ps := NewParser(tokens)
			node, err := ps.Expr()
			fmt.Printf("%v\n", node)
		})
	}
}

func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		title    string
		filepath string
	}{
		{
			"a",
			"../../../../sample/cc/compiler/println",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			tk := tokenizer.NewTokenizer(misc.Scan(test.filepath))
			tokens, err := tk.Tokenize([]tokenizer.TokenType{tokenizer.WHITESPACE, tokenizer.COMMENT, tokenizer.NEWLINE})
			if err != nil {
				t.Fatal(err)
			}
			ps := NewParser(tokens)
			node, err := ps.Expr()
			fmt.Printf("%v\n", node)
		})
	}
}
