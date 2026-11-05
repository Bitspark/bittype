// Package analyze consumes explicit Horos resolution evidence without resolving again.
package analyze

import (
	"context"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
)

const Profile = "horos-declarations-4/analysis-1"

func Contract() analysis.Contract { return ContractWithResolutionConfiguration(nil) }

func ContractWithResolutionConfiguration(settings []semantic.Setting) analysis.Contract {
	resolved := resolve.Contract()
	requirement := resolution.Requirement{Profile: resolved.Semantic.Profile, Configuration: settings, Reports: resolved.Semantic.SemanticsRequirements}.Clone()
	contract := semantic.Contract{Profile: Profile}
	for _, obligation := range resolved.Semantic.Obligations {
		if obligation.Key != resolved.Bindings {
			contract.Obligations = append(contract.Obligations, obligation)
		}
	}
	for _, required := range resolved.Semantic.SemanticsRequirements {
		if required.Obligation != resolved.Bindings {
			contract.SemanticsRequirements = append(contract.SemanticsRequirements, required)
		}
	}
	return analysis.Contract{Semantic: contract, Resolution: &requirement}
}

func Analyze(ctx context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if request.Contract().Profile != Profile || len(request.Configuration()) != 0 {
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
