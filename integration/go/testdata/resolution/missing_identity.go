package missingidentity

import (
	"context"
	"github.com/Bitspark/bittype/contracts/resolution/go"
 "github.com/Bitspark/bittype/contracts/semantic/go"
)

type Language struct{}
type Provider struct{}

func (Provider) Resolve(context.Context, semantic.Request[Language]) (semantic.Result, error) {
	return semantic.Result{}, nil
}

var _ resolution.ResolverProvider[Language] = Provider{}
