package vvm

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTokenizer_ExitOnly(t *testing.T) {
	tk, err := NewTokenizerPath("../../sample/vbin/exit_only.vol.s.b")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "exit", tk.raw)
}

func TestNewTokenizer_EchoExit(t *testing.T) {
	tk, err := NewTokenizerPath("../../sample/vbin/echo_exit.vol.s.b")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "echo \"hello\"\nexit", tk.raw)
}

func TestTokenizer_Tokenize_Single(t *testing.T) {
	tkr, err := NewTokenizerPath("../../sample/vbin/say_hello.vol.s.b.s")
	if err != nil {
		t.Fatal(err)
	}

	tokens, err := tkr.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	for _, t := range *tokens {
		fmt.Printf("%.10s | %.3d-%.3d | %v\n", t.typ.String(), t.sPos, t.ePos, t.lit)
	}
}
