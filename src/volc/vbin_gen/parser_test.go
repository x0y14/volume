package vbin_gen

import (
	"fmt"
	"os"
	"testing"
)

func TestNewParser(t *testing.T) {
	tokenizer := NewTokenizer("main:\n \tcp \"hello\", reg_a\n\tadd reg_b, [bp-3] ; src dst\n\tsub 1.3 [sp+1]\n\texit")
	tokens, err := tokenizer.Tokenize([]TokenType{COMMENT, NEWLINE, WHITESPACE})
	if err != nil {
		t.Fatal(err)
	}

	for _, token := range *tokens {
		fmt.Printf("%v\n", token.String())
	}
	parser := NewParser(*tokens)
	ops, err := parser.Parse()
	if err != nil {
		t.Fatal(err)
	}

	for _, op := range *ops {
		fmt.Printf("%v\n", op.String())
	}
}

func TestParser_Parse(t *testing.T) {
	tests := []string{"../../../sample/vasm/for_loop_echo.vol.s", "../../../sample/vasm/say_hello.vol.s"}

	for _, tt := range tests {
		fmt.Printf("< %40s >\n", tt)

		dat, err := os.ReadFile(tt)
		if err != nil {
			t.Fatal(err)
		}
		text := string(dat)

		tokenizer := NewTokenizer(text)
		tokens, err := tokenizer.Tokenize([]TokenType{COMMENT, NEWLINE, WHITESPACE})
		if err != nil {
			t.Fatal(err)
		}
		//for _, token := range *tokens {
		//	fmt.Printf("%v\n", token.String())
		//}
		parser := NewParser(*tokens)
		ops, err := parser.Parse()
		if err != nil {
			t.Fatal(err)
		}

		for _, op := range *ops {
			if op.label != "" {
				fmt.Println()
			}
			fmt.Printf("%v\n", op.String())
		}
		fmt.Printf("\n\n")
	}
}
