// Package resolve supplies bounded Go resolution and its coupled native checking.
package resolve

import (
	"context"
	"errors"

	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/internal/golangsem"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

func Contract() resolution.Contract {
	c := semantic.Contract{Profile: n.ResolutionProfile, SemanticsRequirements: n.Requirements()}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		c.Obligations = append(c.Obligations, semantic.Obligation{Key: key, Domain: n.Key("captured-package")})
	}
	return resolution.Contract{Semantic: c, Bindings: n.Bindings()}
}
// bittype:scope docs/native-semantics.md#first-go-profile
func Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil || request.Contract().Profile != n.ResolutionProfile || !n.MatchesConfiguration(request.Configuration()) {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	checked, err := golangsem.Check(ctx, request.Inputs())
	publication := semantic.Publication{}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Complete, Verdict: semantic.Holds})
	}
	if err == nil {
		publication.Facts = &checked.Facts
	} else {
		if errors.Is(err, semantic.ErrInvalidRequest) || ctx.Err() != nil {
			return semantic.Result{}, err
		}
		for i := range publication.Reports {
			publication.Reports[i].Status = semantic.Incomplete
			publication.Reports[i].Verdict = semantic.Undetermined
		}
		diagnostic := semantic.Diagnostic{Code: n.Key("incomplete"), Subject: n.Validity(), Message: err.Error()}
		var failure *golangsem.Failure
		if errors.As(err, &failure) {
			if failure.At.Tree != "" {
				at := failure.At
				diagnostic.At = &at
			}
			if failure.Invalid {
				publication.Reports[0] = semantic.Report{Obligation: n.Coverage(), Status: semantic.Complete, Verdict: semantic.Holds}
				publication.Reports[2] = semantic.Report{Obligation: n.Validity(), Status: semantic.Complete, Verdict: semantic.Rejected}
				diagnostic.Code = n.Key("type-error")
			}
		}
		if errors.Is(err, n.ErrUnsupported) {
			publication.Reports[0].Status = semantic.Unsupported
			diagnostic.Code = n.Key("unsupported")
			diagnostic.Subject = n.Coverage()
		} else if errors.Is(err, semantic.ErrLimit) {
			diagnostic.Subject = n.Coverage()
		} else if failure == nil {
			return semantic.Result{}, err
		}
		publication.Diagnostics = []semantic.Diagnostic{diagnostic}
	}
	known := map[semantics.Key]bool{n.Coverage(): true, n.Bindings(): true, n.Validity(): true}
	for _, key := range request.Requested() {
		if !known[key] {
			publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Unsupported})
		}
	}
	return request.Publish(publication)
}

type Provider struct{ s.Language }

func (Provider) Resolve(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	return Resolve(ctx, request)
}

var _ resolution.Provider[s.Language] = Provider{}
