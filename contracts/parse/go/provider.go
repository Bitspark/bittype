// Package parse defines the parse role at the Syntax and Files levels.
package parse

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
)

type Provider[I, O any] interface {
	language.Provider
	Parse(context.Context, I) (O, error)
}
