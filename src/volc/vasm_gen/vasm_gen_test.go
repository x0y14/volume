package vasm_gen

import (
	"fmt"
	"testing"
)

func TestNewVasmGenWithPath(t *testing.T) {
	//vg := NewVasmGenWithPath("../../../sample/vasm/for_loop_echo.vol.s")
	var tests = []struct {
		title string
		path  string
	}{
		{
			"while print",
			"../../../sample/volume/while_print.vol",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			vg := NewVasmGenWithPath(test.path)
			err := vg.Prepare()
			if err != nil {
				t.Fatal(err)
			}
			code, err := vg.GenerateCode()
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("< %v\n", test.path)
			fmt.Printf("%v", code)
			fmt.Printf(">\n\n")
		})
	}
}
