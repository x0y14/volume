package vbin_gen

import (
	"fmt"
	"testing"
)

func TestNewTokenizer(t *testing.T) {
	tokenizer := NewTokenizer("main:\n \tcp \"hello\", reg_a\n\tadd 1, reg_b ; src dst\n\tsub 1.3 reg_c\n\texit")
	tokens, err := tokenizer.Tokenize()
	if err != nil {
		t.Fatal(err)
	}

	for _, token := range *tokens {
		fmt.Printf("%v\n", token.String())
	}
}
