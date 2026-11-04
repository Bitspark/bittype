// Package format defines the format role at the Syntax and Files levels.
package format

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
)

type Provider[I, O any] interface {
	language.Provider
	Format(context.Context, I) (O, error)
}
