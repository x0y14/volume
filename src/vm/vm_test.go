package vm

import "testing"

func TestNewVM(t *testing.T) {
	vm := NewVM()
	vm.SetUp(250, "../../sample/asm/add_test.vol.s")
	vm.Execute()
}

func TestVM_CP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_cp.vol.s")
	vm.Execute()
}

func TestVM_PUSH_POP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_push_pop.vol.s")
	vm.Execute()
}
