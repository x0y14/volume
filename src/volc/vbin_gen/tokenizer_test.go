package vbin_gen

import (
	"fmt"
	"testing"
)

func TestNewTokenizer(t *testing.T) {
	tokenizer := NewTokenizer("main:\n \tcp \"hello\", reg_a\n\tadd reg_b, [bp-3] ; src dst\n\tsub 1.3 [sp+1]\n\texit")
	tokens, err := tokenizer.Tokenize(nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, token := range *tokens {
		fmt.Printf("%v\n", token.String())
	}
}
