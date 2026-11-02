package wrongmarker

import (
	"context"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
)

type A struct{}
type B struct{}
type Provider struct{}

func (Provider) ID() language.ID { return "same-name" }
func (Provider) Name() string    { return "Same name" }
func (Provider) Resolve(context.Context, semantic.Request[A]) (semantic.Result, error) {
	return semantic.Result{}, nil
}

var _ resolution.Provider[B] = Provider{}
