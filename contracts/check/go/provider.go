// Package check defines syntax admission over AST views.
package check

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Provider checks AST structure and supported language syntax through a borrowed
// read-only view. Successful checking does not establish reference resolution
// or semantic validity. No text parsing, formatting or printing is required.
type Provider[L any] interface {
	language.Provider
	Check(ast.View[L]) error
}
