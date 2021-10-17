package vbin_gen

import "fmt"

func InvalidTokenErr(msg string, sPos int, ePos int) error {
	return fmt.Errorf("(@%03d-%03d) %v", sPos, ePos, msg)
}

func SyntaxErr(expectText string, expectType TokenType, actualText string, actualType TokenType) error {
	return fmt.Errorf("SyntaxError: expect %10q(%v), but found %10q(%v)", expectText, expectType.String(), actualText, actualType.String())
}
