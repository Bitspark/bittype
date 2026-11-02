// Package provider supplies an explicitly registered, AST-only Horos provider.
// It has no text parser, formatter, printer or native-language transformer.
package provider

import (
	"context"

	analyze "github.com/Bitspark/bittype/analyze/horos/go"
	check "github.com/Bitspark/bittype/check/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Provider struct{ s.Language }

func (Provider) Check(view ast.View[s.Language]) error { return check.Check(view) }
func (Provider) Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	return analyze.Analyze(ctx, request)
}

var _ text.Checker[s.Language] = Provider{}
var _ analysis.Provider[s.Language] = Provider{}

func (Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return resolve.ResolveReferences(ctx, request)
}

var _ resolution.Provider[s.Language] = Provider{}
