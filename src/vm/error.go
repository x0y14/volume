package vm

import "fmt"

func LoadInvalidTypeErr(actualType TokenType, orderedType TokenType) error {
	return fmt.Errorf("this token type is %v, but i ordered to load as %v", actualType.String(), orderedType.String())
}

func UnexpectedTokenTypeErr(expected []TokenType, actual TokenType) error {
	expectedStr := ""
	for _, ex := range expected {
		expectedStr += ex.String() + ", "
	}
	return fmt.Errorf("expcted token type: %vbut actual token type: %v", expectedStr, actual.String())
}

func UnexpectedKeyWordErr(expected []KeyWordType, actual KeyWordType) error {
	expectedStr := ""
	for _, ex := range expected {
		expectedStr += ex.String() + ", "
	}
	return fmt.Errorf("expcted keyword type: %vbut actual keyword type: %v", expectedStr, actual.String())
}

func UnexpectedKPointerTypeErr(expected []PointerType, actual PointerType) error {
	expectedStr := ""
	for _, ex := range expected {
		expectedStr += ex.String() + ", "
	}
	return fmt.Errorf("expcted pointer type: %vbut actual pointer type: %v", expectedStr, actual.String())
}

func StackAccessErr(max int, pointer int) error {
	return fmt.Errorf("stack invalid access: %v, you can access: %v-%v", pointer, 0, max)
}

func UndefinedKeyWordErr(text string) error {
	return fmt.Errorf("undefined keyword: %v", text)
}

func UnDefinedRegisterErr(text string) error {
	return fmt.Errorf("undefined register: %v", text)
}

func UndefinedPointerErr(text string) error {
	return fmt.Errorf("undefined pointer: %v", text)
}