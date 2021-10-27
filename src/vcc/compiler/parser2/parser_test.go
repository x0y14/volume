package parser2

import (
	"fmt"
	"github.com/x0y14/volume/src/vcc/compiler/misc"
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
	"testing"
)

func TestParser_ParseWithText(t *testing.T) {
	var tests = []struct {
		title string
		in    string
	}{
		{
			"func",
			"func main(a int) int {var text = \"hello\";}",
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
			node, err := ps.Parse()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("%v\n", node)
		})
	}
}

func TestParser_ParseWithPath(t *testing.T) {
	var tests = []struct {
		title string
		in    string
	}{
		{
			"all",
			"./test/syntax_all.vol",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			script := misc.Scan(test.in)
			tk := tokenizer.NewTokenizer(script)
			tokens, err := tk.Tokenize([]tokenizer.TokenType{tokenizer.WHITESPACE, tokenizer.COMMENT, tokenizer.NEWLINE})
			if err != nil {
				t.Fatal(err)
			}

			ps := NewParser(tokens)
			_, err = ps.Parse()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
