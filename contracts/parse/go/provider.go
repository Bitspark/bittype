// Package parse defines source inputs and text-to-CST providers.
package parse

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Provider parses source text into a concrete tree with its stored AST view.
type Provider[L any] interface {
	language.Provider
	Parse(Input) (cst.CST[L], error)
}
