package parser3

import (
	"fmt"
	"github.com/x0y14/volume/src/vcc/compiler/tokenizer"
)

func UnexpectedTokenError(expect []tokenizer.TokenType, actual tokenizer.TokenType) error {
	expectStr := ""
	for i, exp := range expect {
		expectStr += exp.String()
		if len(expect)-1 != i {
			expectStr += ", "
		}
	}
	return fmt.Errorf("unexpected token-type: expected[ %v ], actual[ %v ]", expectStr, actual.String())
}
