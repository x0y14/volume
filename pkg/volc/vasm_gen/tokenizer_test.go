package vasm_gen

import (
	"fmt"
	"os"
	"testing"
)

func TestTokenizer_Tokenize(t *testing.T) {
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

			for _, tok := range tokens {
				fmt.Printf("%v\n", tok.String())
			}
		})
	}
}
