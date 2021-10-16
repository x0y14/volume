package vbin_gen

import "fmt"

func InvalidTokenErr(msg string, sPos int, ePos int) error {
	return fmt.Errorf("(@%03d-%03d) %v", sPos, ePos, msg)
}
