package vasm_gen

import (
	"fmt"
	"os"
	"testing"
)

func TestParser_Parse(t *testing.T) {
	var tests = []struct {
		title string
		path  string
	}{
		{
			"comment",
			"../../../sample/volume/comment.vol",
		},
		{
			"main only",
			"../../../sample/volume/main_println.vol",
		},
		{
			"main arg",
			"../../../sample/volume/main_arg.vol",
		},
		{
			"main args",
			"../../../sample/volume/main_multi_args.vol",
		},
		{
			"two func",
			"../../../sample/volume/two_func.vol",
		},
		{
			"for_print",
			"../../../sample/volume/for_print.vol",
		},
		{
			"variable",
			"../../../sample/volume/variable.vol",
		},
		{
			"simple main",
			"../../../sample/volume/simple_main.vol",
		},
		{
			"variable call func",
			"../../../sample/volume/variable_call_func.vol",
		},
		{
			"while print",
			"../../../sample/volume/while_print.vol",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			data, err := os.ReadFile(test.path)
			if err != nil {
				t.Fatal(err)
			}
			text := string(data)

			tokenizer := NewTokenizer(text)
			tokens, err := tokenizer.Tokenize([]TokenType{WHITESPACE, NEWLINE})
			if err != nil {
				t.Fatal(err)
			}

			parser := NewParser(tokens)
			if nodes, err := parser.Parse(); err != nil {
				t.Fatal(err)
			} else {
				fmt.Printf("%v\n", nodes)
			}
		})
	}
}
