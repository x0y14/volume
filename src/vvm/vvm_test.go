package vvm

import "testing"

func TestNewVM(t *testing.T) {
	vm := NewVM()
	vm.SetUp(250, "../../sample/vbin/add_test.vol.b")
	vm.Execute()
}

func TestVM_CP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_cp.vol.b")
	vm.Execute()
}

func TestVM_PUSH_POP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_push_pop.vol.b")
	vm.Execute()
}

func TestVM_SP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_sp.vol.b")
	vm.Execute()
}

func TestVM_ECHO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_echo.vol.b")
	vm.Execute()
}

func TestVM_CALL_RET(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_call_ret.vol.b")
	vm.Execute()
}

func TestVM_CMP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_cmp.vol.b")
	vm.Execute()
}

func TestVM_JUMP(t *testing.T) {
	vm := NewVM()
	vm.SetUp(10, "../../sample/vbin/test_jump_jz_jnz.vol.b")
	vm.Execute()
}

func TestVM_SAY_HELLO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/vbin/say_hello.vol.b.s")
	vm.Execute()
}

func TestVM_ADD(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/vbin/test_add.vol.b")
	vm.Execute()
}

func TestVM_SUB(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/vbin/test_sub.vol.b")
	vm.Execute()
}

func TestVM_COMMENT(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/vbin/test_comment.vol.b")
	vm.Execute()
}

func TestVM_FOT_LOOP_ECHO(t *testing.T) {
	vm := NewVM()
	vm.SetUp(20, "../../sample/vbin/for_loop_echo.vol.b")
	vm.Execute()
}

func TestVM_Script(t *testing.T) {
	var tests = []struct {
		title string
		path  string
	}{
		{
			"while print",
			"../../sample/vbin/while_print_expect.vol.b",
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			vm := NewVM()
			vm.SetUp(30, test.path)
			vm.Execute()
		})
	}
}
