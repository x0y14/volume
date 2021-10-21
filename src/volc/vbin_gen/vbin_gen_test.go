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

func TestNewVBinGen2(t *testing.T) {
	var tests = []struct {
		title  string
		path   string
		output string
	}{
		{
			"while loop",
			"../../../sample/vasm/while_print_expect.vol.s",
			"../../../sample/vbin/while_print_expect.vol.b",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			fmt.Printf("============== < %v > ==============\n", test.path)

			// read file
			dat, err := os.ReadFile(test.path)
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

			vbinCode := vbg.AsString()

			fmt.Printf("%v\n", vbinCode)
			fmt.Printf("\n\n")

			if err := vbg.Export(test.output); err != nil {
				t.Fatal(err)
			}

		})
	}
}
