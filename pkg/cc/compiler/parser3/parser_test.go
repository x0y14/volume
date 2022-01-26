package parser3

import (
	"fmt"
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
