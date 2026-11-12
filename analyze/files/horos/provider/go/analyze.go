// Package analyze consumes explicit Horos resolution evidence without resolving again.
package analyze

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
)

// bittype:scope docs/analysis.md#horos-profile
func Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if request.Contract().Profile != n.AnalysisProfile || len(request.Configuration()) != 0 {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	resolved, ok := request.Resolution()
	if !ok {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	publication := semantic.Publication{Diagnostics: resolved.Diagnostics()}
	known := make(map[semantics.Key]bool)
	for _, obligation := range request.Contract().Obligations {
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

// Provider supplies this role independently of other language capabilities.
type Provider struct{ s.Language }

func (Provider) Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	return Analyze(ctx, request)
}

var _ analysis.Provider[s.Language] = Provider{}
