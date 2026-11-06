package print

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	printcontract "github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Provider[L any] struct{ identity language.Provider }

func New[L any](identity language.Provider) (Provider[L], error) {
	if representation.IsNil(identity) {
		return Provider[L]{}, printcontract.ErrInvalidPrinter
	}
	return Provider[L]{identity: identity}, nil
}
// bittype:scope docs/syntax.md#identity-and-registration
func (p Provider[L]) ID() language.ID { return p.identity.ID() }
// bittype:scope docs/syntax.md#identity-and-registration
func (p Provider[L]) Name() string    { return p.identity.Name() }
func (p Provider[L]) Print(ctx context.Context, input cst.CST[L]) (string, error) {
	return Print(ctx, input)
}
