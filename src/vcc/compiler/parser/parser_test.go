package parser

import (
	"fmt"
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
			_, err = ps.Parse()
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
			"args: 1, ret: 0, body: 0",
			"func main(n int) {}",
		},
		{
			"args: 1, ret: 1, body: 0",
			"func main(n int) int {}",
		},
		{
			"args: 1, ret: (1), body: 0",
			"func main(n int) (int) {}",
		},
		{
			"args: 1, ret: 2, body: 0",
			"func main() (int, string) {}",
		},

		{
			"args: 0, ret: 0, body: 0",
			"func main() {}",
		},
		{
			"args: 0, ret: 1, body: 0",
			"func main() string {}",
		},
		{
			"args: 0, ret: (1), body: 0",
			"func main() (string) {}",
		},
		{
			"args: 0, ret: 2, body: 0",
			"func main() (string, int) {}",
		},
		{
			"import",
			"import \"stdio\"",
		},
		{
			"args: 0, ret: 0, body: 1",
			"func main() { var text = \"hello, world\" }",
		},
		{
			"args: 0, ret: 0, body: 2",
			"func main() { var text = \"hello, world1\"\nvar text2 = \"hello, world2\" }",
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
			nodes, err := ps.Parse()
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("\n[Nodes]\n")
			for _, node := range nodes {
				fmt.Printf("%v\n", node.String())
			}
			fmt.Println()
		})
	}
}

func TestParser_consumeRHS(t *testing.T) {
	var tests = []struct {
		title string
		in    string
		want  bool
	}{
		{
			"suc) 単項 int",
			"1",
			true,
		},
		{
			"suc) 2項の足し算",
			"1 + 2",
			true,
		},
		{
			"suc) 2項 int, ident",
			"3 + n",
			true,
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
			node, err := ps.consumeRHS()
			if err != nil {
				t.Fatal(err)
			}

			//fmt.Printf("%v\n", node.String())

			fmt.Printf("\n[Nodes]\n")
			for _, node := range node.childrenNode {
				fmt.Printf("%v\n", node.String())
			}
			fmt.Println()
		})
	}
}
