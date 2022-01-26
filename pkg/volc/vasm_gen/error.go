package vasm_gen

import "fmt"

func SyntaxErr(title string, expect string, actual string) error {
	return fmt.Errorf("syntax error(%v): expected %v, but found %v", title, expect, actual)
}

func NotYetImplementedErr(title string, about string) error {
	return fmt.Errorf("not yet implemented(%v): %v", title, about)
}
