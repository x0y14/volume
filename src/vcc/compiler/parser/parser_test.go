package parser

import (
	"github.com/x0y14/volume/src/vcc/compiler/misc"
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
	"testing"
)

func TestParser_ParseWithPath(t *testing.T) {
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

func TestParser_ParseWithText(t *testing.T) {
	var tests = []struct {
		title  string
		script string
	}{
		{
			"args: 1, ret: 0, body: 1",
			"func main(n int) {}",
		},
		{
			"args: 1, ret: 1, body: 1",
			"func main(n int) int {}",
		},
		{
			"args: 1, ret: (1), body: 1",
			"func main(n int) (int) {}",
		},
		{
			"args: 1, ret: 2, body: 1",
			"func main() (int, string) {}",
		},

		{
			"args: 0, ret: 0, body: 1",
			"func main() {}",
		},
		{
			"args: 0, ret: 1, body: 1",
			"func main() string {}",
		},
		{
			"args: 0, ret: (1), body: 1",
			"func main() (string) {}",
		},
		{
			"args: 0, ret: 2, body: 1",
			"func main() (string, int) {}",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			tk := tokenizer.NewTokenizer(test.script)
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
