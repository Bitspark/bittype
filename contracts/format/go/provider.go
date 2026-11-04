// Package format defines AST-to-CST providers.
package format

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Provider constructs a concrete tree from a borrowed, read-only AST view.
type Provider[L any] interface {
	language.Provider
	Format(ast.View[L]) (cst.CST[L], error)
}
