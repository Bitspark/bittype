// Package print defines the print role at the Syntax and Files levels.
package print

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
)

type Provider[I, O any] interface {
	language.Provider
	Print(context.Context, I) (O, error)
}
