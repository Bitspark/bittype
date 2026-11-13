// Package resolve supplies the captured Rust compiler profile.
package resolve

import (
	"context"
	"github.com/Bitspark/bittype/trees/representation/go"

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
	HelperPath  string
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

// bittype:scope docs/native-semantics.md#first-rust-profile
func (p Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil || request.Contract().Profile != n.ResolutionProfile {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	// The request states its environment. Settings that no supported
	// environment gives are refused here, before the compiler runs.
	environment, err := n.EnvironmentOf(request.Configuration())
	if err != nil {
		return semantic.Result{}, err
	}
	if err := p.Compiler.Validate(); err != nil {
		return semantic.Result{}, err
	}
	out, err := rustsem.Check(ctx, p.Compiler.HelperPath, p.Compiler.SysrootPath, rustsem.Environment{Edition: environment.Edition, Target: environment.Target, CrateName: environment.Crate.Name, CrateKind: environment.Crate.Kind}, request.Inputs(), request.Files())
	if err != nil {
		return semantic.Result{}, err
	}
	publication := semantic.Publication{Facts: out.Facts, Diagnostics: out.Diagnostics}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Incomplete})
	}
	switch out.Status {
	case rustsem.StatusOK:
		for i := range publication.Reports {
			publication.Reports[i].Status = semantic.Complete
			publication.Reports[i].Verdict = semantic.Holds
		}
	case rustsem.StatusInvalid:
		publication.Reports[0] = semantic.Report{Obligation: n.Coverage(), Status: semantic.Complete, Verdict: semantic.Holds}
		publication.Reports[2] = semantic.Report{Obligation: n.Validity(), Status: semantic.Complete, Verdict: semantic.Rejected}
	case rustsem.StatusUnsupported:
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
func New(definition representation.Language, compiler Compiler) (resolution.Resolver, error) {
	if err := compiler.Validate(); err != nil {
		return resolution.Resolver{}, err
	}
	return resolution.NewResolver[s.Language](definition, n.ResolutionContract(), Provider{Compiler: compiler})
}
