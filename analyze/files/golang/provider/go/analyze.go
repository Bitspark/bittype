// Package analyze consumes Go resolution evidence.
package analyze

import (
	"context"
	"reflect"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

// bittype:scope docs/native-semantics.md#first-go-profile
func Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil || request.Contract().Profile != n.AnalysisProfile || !n.MatchesConfiguration(request.Configuration()) {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	resolved, ok := request.Resolution()
	// Analysis reuses the resolution's graph, so it must claim exactly the
	// environment that graph was computed under.
	if !ok || !reflect.DeepEqual(request.Configuration(), resolved.Configuration()) {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	publication := semantic.Publication{Diagnostics: resolved.Diagnostics()}
	known := map[semantics.Key]bool{}
	for _, obligation := range n.AnalysisContract().Semantic.Obligations {
		report, found := resolved.Report(obligation.Key)
		if !found {
			report = semantic.Report{Obligation: obligation.Key, Status: semantic.Incomplete}
		}
		publication.Reports = append(publication.Reports, report)
		known[obligation.Key] = true
	}
	for _, key := range request.Requested() {
		if !known[key] {
			publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Unsupported})
		}
	}
	if resolved.HasSemantics() {
		publication.ReuseSemanticsFrom = &resolved
	}
	return request.Publish(publication)
}

type Provider struct{ s.Language }

func (Provider) Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	return Analyze(ctx, request)
}

var _ analysis.Provider[s.Language] = Provider{}
