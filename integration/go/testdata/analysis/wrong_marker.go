package wrongmarker

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/language/go"
)

type A struct{}
type B struct{}
type Provider struct{}
func (Provider) ID() language.ID { return "same-name" }
func (Provider) Name() string { return "Same name" }
func (Provider) Analyze(context.Context, analysis.Request[A]) (analysis.Result, error) { return analysis.Result{}, nil }
var _ analysis.Provider[B] = Provider{}
