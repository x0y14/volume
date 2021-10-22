package parser

import "fmt"

func SyntaxErr(section string, expected string, actual string) error {
	return fmt.Errorf("syntax error(%v): expected %v, but got %v", section, expected, actual)
}
