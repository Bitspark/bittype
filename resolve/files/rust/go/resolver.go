// Package resolve supplies the captured Rust compiler profile.
package resolve

import (
	"context"

	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/rustsem"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

// Compiler names explicitly provisioned executables. Constructors do no discovery
// or installation; SysrootPath selects the matching pinned compiler sysroot.
type Compiler struct {
	HelperPath       string
	SysrootPath string
}

func (c Compiler) Validate() error {
	if !rustsem.ValidCompiler(c.HelperPath, c.SysrootPath) {
		return semantic.ErrInvalidOperation
	}
	return nil
}

type Provider struct {
	s.Language
	Compiler Compiler
}

func Contract() resolution.Contract {
	c := semantic.Contract{Profile: n.ResolutionProfile, SemanticsRequirements: n.Requirements()}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		c.Obligations = append(c.Obligations, semantic.Obligation{Key: key, Domain: n.Key("captured-program")})
	}
	return resolution.Contract{Semantic: c, Bindings: n.Bindings()}
}
func (p Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil || request.Contract().Profile != n.ResolutionProfile || !n.MatchesConfiguration(request.Configuration()) {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := p.Compiler.Validate(); err != nil {
		return semantic.Result{}, err
	}
	out, err := rustsem.Check(ctx, p.Compiler.HelperPath, p.Compiler.SysrootPath, request.Inputs(), request.Files())
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
