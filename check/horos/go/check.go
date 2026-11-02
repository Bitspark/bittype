// Package check validates the Horos AST profile through Horos's declaration
// model. It checks syntax, pins and binding edges, not application values.
package check

import (
	bridge "github.com/Bitspark/bittype/internal/horos"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func Check(view ast.View[s.Language]) error {
	_, err := bridge.Read(view)
	return err
}
