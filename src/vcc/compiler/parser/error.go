package parser

import "fmt"

func SyntaxErr(section string, expected string, actual string) error {
	return fmt.Errorf("syntax error(%v): expected %v, but got %v", section, expected, actual)
}

func NotYetImplErr(section string, text string) error {
	return fmt.Errorf("not yet implemented(%v): %v", section, text)
}
