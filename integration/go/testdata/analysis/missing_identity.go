package missingidentity

import (
	"context"
	"github.com/Bitspark/bittype/contracts/analysis/go"
)

type Language struct{}
type Provider struct{}

func (Provider) Analyze(context.Context, analysis.Request[Language]) (analysis.Result, error) {
	return analysis.Result{}, nil
}

var _ analysis.Provider[Language] = Provider{}
