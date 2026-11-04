// Package print defines CST-to-text providers.
package print

import (
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/cst/go"
)

// Provider emits source text from a concrete tree.
type Provider[L any] interface {
	language.Provider
	Print(cst.CST[L]) (string, error)
}
