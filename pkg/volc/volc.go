package volc

import (
	"fmt"
	"github.com/x0y14/volume/pkg/volc/vasm_gen"
	_ "github.com/x0y14/volume/pkg/volc/vbin_gen"
)

func NewVolc() Volc {
	return Volc{}
}

type Volc struct {
}

func (v *Volc) Build(volPath string) {
	vasmGen := vasm_gen.NewVasmGenWithPath(volPath)
	err := vasmGen.Prepare()
	if err != nil {
		panic(err)
	}
	libs := vasmGen.LibNeedForBuild()
	asm, err := vasmGen.GenerateCode()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n\n", libs)
	fmt.Printf("%v\n", asm)
}
