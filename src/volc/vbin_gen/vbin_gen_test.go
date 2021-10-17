package vbin_gen

import (
	"fmt"
	"os"
	"testing"
)

func TestNewVBinGen(t *testing.T) {
	tests := []string{"../../../sample/vasm/for_loop_echo.vol.s", "../../../sample/vasm/say_hello.vol.s"}

	for _, tt := range tests {
		// title
		fmt.Printf("============== < %v > ==============\n", tt)

		// read file
		dat, err := os.ReadFile(tt)
		if err != nil {
			t.Fatal(err)
		}
		text := string(dat)

		// tokenizer
		tokenizer := NewTokenizer(text)
		tokens, err := tokenizer.Tokenize([]TokenType{COMMENT, NEWLINE, WHITESPACE})
		if err != nil {
			t.Fatal(err)
		}

		// parser
		parser := NewParser(*tokens)
		ops, err := parser.Parse()
		if err != nil {
			t.Fatal(err)
		}

		// show op
		for _, op := range *ops {
			if op.label != "" {
				fmt.Println()
			}
			fmt.Printf("%v\n", op.String())
		}

		// vbin-gen
		vbg := NewVBinGen(*ops)
		vbg.Scan()

		fmt.Printf("%v\n", vbg.AsString())

		fmt.Printf("\n\n")
	}
}
