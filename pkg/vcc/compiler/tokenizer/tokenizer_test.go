package tokenizer

import (
	"fmt"
	"github.com/x0y14/volume/pkg/vcc/compiler/misc"
	"testing"
)

func TestTokenizer_TokenizeWithPath(t *testing.T) {
	var tests = []struct {
		title string
		path  string
	}{
		{
			"while_loop_print_range_variable",
			"../../../../sample/vcc/proj/while_loop_print_range_variable/script.vol",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			script := misc.Scan(test.path)
			tk := NewTokenizer(script)
			tokens, err := tk.Tokenize([]TokenType{WHITESPACE, COMMENT, NEWLINE})
			if err != nil {
				t.Fatal(err)
			}
			for i, tok := range tokens {
				fmt.Printf("%03d | %v\n", i, tok.String())
			}
		})
	}
}

func TestTokenizer_TokenizeWithText(t *testing.T) {
	var tests = []struct {
		title  string
		script string
	}{
		{
			"OR AND",
			"|| &&",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			tk := NewTokenizer(test.script)
			tokens, err := tk.Tokenize([]TokenType{WHITESPACE, COMMENT, NEWLINE})
			if err != nil {
				t.Fatal(err)
			}
			for i, tok := range tokens {
				fmt.Printf("%03d | %v\n", i, tok.String())
			}
		})
	}
}
