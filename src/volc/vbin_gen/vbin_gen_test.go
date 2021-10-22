package vbin_gen

import (
	"fmt"
	"testing"
)

func TestNewVBinGen(t *testing.T) {
	var tests = []struct {
		title  string
		path   string
		output string
	}{
		{
			"while loop ",
			"../../../sample/vasm/while_print_expect.vol.s",
			"../../../sample/vbin/while_print_expect.vol.b",
		},
		{
			"while loop n as local",
			"../../../sample/vasm/while_print_local_var.vol.s",
			"../../../sample/vbin/while_print_local_var.vol.b",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			fmt.Printf("============== < %v > ==============\n", test.path)

			// vbin-gen
			vbg := NewVBinGen()
			err := vbg.Prepare(test.path)
			if err != nil {
				t.Fatal(err)
			}
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
