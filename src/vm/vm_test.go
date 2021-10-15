package vm

import "testing"

func TestNewVM(t *testing.T) {
	vm := NewVM()
	vm.SetUp(250, "../../sample/machine_lang/add_test.vol.s.b")
	vm.Execute()
}

func TestVM_CP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_cp.vol.s.b")
	vm.Execute()
}

func TestVM_PUSH_POP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_push_pop.vol.s.b")
	vm.Execute()
}

func TestVM_SP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_sp.vol.s.b")
	vm.Execute()
}

func TestVM_ECHO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_echo.vol.s.b")
	vm.Execute()
}

func TestVM_CALL_RET(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_call_ret.vol.s.b")
	vm.Execute()
}

func TestVM_CMP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_cmp.vol.s.b")
	vm.Execute()
}

func TestVM_JUMP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/machine_lang/test_jump_jz_jnz.vol.s.b")
	vm.Execute()
}

func TestVM_SAY_HELLO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/machine_lang/say_hello.vol.s.b.s")
	vm.Execute()
}

func TestVM_ADD(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/machine_lang/test_add.vol.s.b")
	vm.Execute()
}

func TestVM_SUB(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/machine_lang/test_sub.vol.s.b")
	vm.Execute()
}

func TestVM_COMMENT(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/machine_lang/test_comment.vol.s.b")
	vm.Execute()
}

func TestVM_FOT_LOOP_ECHO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/machine_lang/for_loop_echo.vol.s.b")
	vm.Execute()
}
