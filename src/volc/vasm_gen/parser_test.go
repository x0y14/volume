package vasm_gen

import (
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
			if err := parser.Parse(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
