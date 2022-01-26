package volc

import (
	"github.com/stretchr/testify/assert"
	"github.com/x0y14/volume/pkg/volc/vbin_gen"
	"github.com/x0y14/volume/pkg/vvm"
	"testing"
)

func TestNewVolc(t *testing.T) {
	volc := NewVolc()
	volc.Build("../../sample/volume/simple_while_loop.vol")
}

func TestRunVasm(t *testing.T) {
	var tests = []struct {
		title   string
		asmPath string
		binPath string
		stream  []string
	}{
		{
			"while print local var",
			"../../sample/vasm/while_print_local_var.vol.s",
			"../../sample/vbin/while_print_local_var.vol.b",
			[]string{"5", "4", "3", "2", "1"},
		},
		{
			"simple while loop",
			"../../sample/vasm/simple_while_loop_exp.vol.s",
			"../../sample/vbin/simple_while_loop_exp.vol.b",
			[]string{"0", "1", "2"},
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			vbg := vbin_gen.NewVBinGen()
			err := vbg.Prepare(test.asmPath)
			if err != nil {
				t.Fatal(err)
			}
			vbg.Scan()
			err = vbg.Export(test.binPath)
			if err != nil {
				t.Fatal(err)
			}
			vm := vvm.NewVM()
			vm.SetUp(40, test.binPath)
			stream := vm.Execute()
			assert.Equal(t, test.stream, stream)
		})
	}
}
