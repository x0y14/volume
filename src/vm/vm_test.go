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

func TestVM_SP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_sp.vol.s")
	vm.Execute()
}

func TestVM_ECHO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_echo.vol.s")
	vm.Execute()
}

func TestVM_CALL_RET(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_call_ret.vol.s")
	vm.Execute()
}

func TestVM_CMP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_cmp.vol.s")
	vm.Execute()
}

func TestVM_JUMP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/asm/test_jump_jz_jnz.vol.s")
	vm.Execute()
}

func TestVM_SAY_HELLO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/asm/say_hello.vol.s")
	vm.Execute()
}
