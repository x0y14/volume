package vasm_gen

import "testing"

//func TestIsKeyword(t *testing.T) {
//	var test = []struct {
//		title string
//		value string
//		expect bool
//	}{
//		{
//			"for",
//			"for",
//			true,
//		},
//	}
//
//	//for _, tt := range test {
//	//	t.Run(tt.title, func(t *testing.T) {
//	//		if got := Ident2Keyword(tt.value); got != tt.expect {
//	//			t.Fatalf("IsKeyword(%q) = %t, expect %v", tt.value, got, tt.expect)
//	//		}
//	//	})
//	//}
//}

func TestOperator2Type(t *testing.T) {
	var tests = []struct {
		title  string
		value  string
		expect TokenType
	}{
		{
			"plus eq",
			"+=",
			PLUSEq,
		},
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			if got := Operator2Type(test.value); got != test.expect {
				t.Fatalf("Operator2Type(%q) = %v, expect %v", test.value, got, test.expect)
			}
		})
	}
}
