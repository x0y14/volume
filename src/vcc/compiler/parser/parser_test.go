package parser

import (
	"github.com/x0y14/volume/src/vcc/compiler/misc"
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		title string
		path  string
	}{
		{
			"while_loop_print_range_variable",
			"../../../../sample/vcc/proj/while_loop_print_range_variable/script.vol",
		},
		{
			"define function",
			"./tests/define_func.vol",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			script := misc.Scan(test.path)
			tk := tokenizer.NewTokenizer(script)
			tokens, err := tk.Tokenize([]tokenizer.TokenType{tokenizer.WHITESPACE, tokenizer.COMMENT, tokenizer.NEWLINE})
			if err != nil {
				t.Fatal(err)
			}

			ps := NewParser(tokens)
			err = ps.Parse()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
