package missingidentity

import (
	"context"
	"github.com/Bitspark/bittype/contracts/analysis/go"
 "github.com/Bitspark/bittype/contracts/semantic/go"
)

type Language struct{}
type Provider struct{}

func (Provider) Analyze(context.Context, analysis.Request[Language]) (semantic.Result, error) {
	return semantic.Result{}, nil
}

var _ analysis.Provider[Language] = Provider{}
