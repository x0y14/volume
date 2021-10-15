package vm

import "fmt"

func LoadInvalidTypeErr(actualType TokenType, orderedType TokenType) error {
	return fmt.Errorf("this token type is %v, but i ordered to load as %v", actualType.String(), orderedType.String())
}

func UnexpectedTokenTypeErr(sectionName string, expected []TokenType, actual TokenType) error {
	expectedStr := ""
	for _, ex := range expected {
		expectedStr += ex.String() + ", "
	}
	return fmt.Errorf("(%v) expcted token type: %vbut actual token type: %v", sectionName, expectedStr, actual.String())
}

func UnexpectedKeyWordErr(expected []KeyWordType, actual KeyWordType) error {
	expectedStr := ""
	for _, ex := range expected {
		expectedStr += ex.String() + ", "
	}
	return fmt.Errorf("expcted keyword type: %vbut actual keyword type: %v", expectedStr, actual.String())
}

func UnexpectedKPointerTypeErr(sectionName string, expected []PointerType, actual PointerType) error {
	expectedStr := ""
	for _, ex := range expected {
		expectedStr += ex.String() + ", "
	}
	return fmt.Errorf("(%v) expcted pointer type: %vbut actual pointer type: %v", sectionName, expectedStr, actual.String())
}

func DoseNotMatchTokenTypeErr(tok1 TokenType, tok2 TokenType) error {
	return fmt.Errorf("two token's type dosen't match: %v, %v", tok1.String(), tok2.String())
}

func StackAccessErr(sectionName string, max int, pointer int) error {
	return fmt.Errorf("(%v) stack invalid access: %v, you can access: %v <= X <= %v", sectionName, pointer, 0, max)
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
