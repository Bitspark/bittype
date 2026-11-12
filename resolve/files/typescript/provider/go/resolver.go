// Package resolve supplies the captured TypeScript compiler profile.
package resolve

import (
	"context"
	"github.com/Bitspark/bittype/trees/representation/go"

	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/typescriptsem"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

type Provider struct {
	s.Language
	Compiler n.Compiler
}

// bittype:scope docs/native-semantics.md#first-typescript-profile
func (p Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil || request.Contract().Profile != n.ResolutionProfile || !n.MatchesConfiguration(request.Configuration()) {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := p.Compiler.Validate(); err != nil {
		return semantic.Result{}, err
	}
	out, err := typescriptsem.Check(ctx, p.Compiler.NodePath, p.Compiler.TypeScriptPath, request.Inputs(), request.Files(), nil, request.Configuration())
	if err != nil {
		return semantic.Result{}, err
	}
	publication := semantic.Publication{Facts: out.Facts, Diagnostics: out.Diagnostics}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Incomplete})
	}
	switch out.Status {
	case "ok":
		for i := range publication.Reports {
			publication.Reports[i].Status = semantic.Complete
			publication.Reports[i].Verdict = semantic.Holds
		}
	case "invalid":
		publication.Reports[0] = semantic.Report{Obligation: n.Coverage(), Status: semantic.Complete, Verdict: semantic.Holds}
		publication.Reports[2] = semantic.Report{Obligation: n.Validity(), Status: semantic.Complete, Verdict: semantic.Rejected}
	case "unsupported":
		publication.Reports[0].Status = semantic.Unsupported
	}
	known := map[semantics.Key]bool{n.Coverage(): true, n.Bindings(): true, n.Validity(): true}
	for _, key := range request.Requested() {
		if !known[key] {
			publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Unsupported})
		}
	}
	return request.Publish(publication)
}

var _ resolution.Provider[s.Language] = Provider{}

// New captures an explicitly provisioned compiler and constructs only resolution.
func New(definition representation.Language, compiler n.Compiler) (resolution.Resolver, error) {
	if err := compiler.Validate(); err != nil {
		return resolution.Resolver{}, err
	}
	return resolution.NewResolver[s.Language](definition, n.ResolutionContract(), Provider{Compiler: compiler})
}
