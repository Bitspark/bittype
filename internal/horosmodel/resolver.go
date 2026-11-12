// Package resolve adapts Horos declaration semantics to shared semantic facts.
// It resolves no names itself and never interprets descriptor constructors.
package horosmodel

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

func (f Frontend[L]) configuredLimits(settings []semantic.Setting) (decl.Limits, error) {
	limits := decl.DefaultLimits()
	if len(settings) == 0 {
		return limits, nil
	}
	if len(settings) != 1 || settings[0].Key != n.LimitsSetting() {
		return decl.Limits{}, fmt.Errorf("%w: unsupported Horos configuration", semantic.ErrInvalidRequest)
	}
	decoder := json.NewDecoder(bytes.NewReader(settings[0].Value))
	var values []int
	if err := decoder.Decode(&values); err != nil || len(values) != 4 {
		return decl.Limits{}, fmt.Errorf("%w: expected four Horos limits", semantic.ErrInvalidRequest)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return decl.Limits{}, fmt.Errorf("%w: trailing Horos configuration", semantic.ErrInvalidRequest)
	}
	// Configuration is a request precondition, even when source decoding will
	// also fail. These are the fixed resolution profile's resource-option bounds.
	if values[0] <= 0 || values[0] > 512 || values[1] <= 0 || values[2] <= 0 || values[3] <= 0 {
		return decl.Limits{}, fmt.Errorf("%w: invalid Horos resource limits", semantic.ErrInvalidRequest)
	}
	return decl.Limits{MaxDepth: values[0], MaxNodes: values[1], MaxBytes: values[2], MaxModules: values[3]}, nil
}

func (f Frontend[L]) reports(request semantic.Request[L], status semantic.Completion, verdict semantic.Verdict, complete bool) []semantic.Report {
	var result []semantic.Report
	known := make(map[semantics.Key]bool)
	for _, obligation := range n.ResolutionContract().Semantic.Obligations {
		r := semantic.Report{Obligation: obligation.Key, Status: semantic.Incomplete}
		if complete {
			r.Status = semantic.Complete
		}
		if obligation.Key == n.ReferenceResolution() && status == semantic.Complete && verdict == semantic.Holds {
			r.Status, r.Verdict = semantic.Complete, semantic.Holds
		}
		if obligation.Key == n.Validity() {
			r.Status, r.Verdict = status, verdict
		}
		result = append(result, r)
		known[obligation.Key] = true
	}
	for _, requested := range request.Requested() {
		if !known[requested] {
			result = append(result, semantic.Report{Obligation: requested, Status: semantic.Unsupported})
		}
	}
	return result
}

func (f Frontend[L]) failure(request semantic.Request[L], decoded *Decoded[L], err error) (semantic.Result, error) {
	status, verdict := semantic.Complete, semantic.Rejected
	diagnostic := semantic.Diagnostic{Code: n.Key("invalid-source"), Subject: n.Key("source"), Message: err.Error()}
	var detail *decl.Error
	var assertion *ValidationError
	var syntax *text.Diagnostic
	switch {
	case errors.As(err, &detail):
		if detail.Code == decl.InvalidOptions || detail.Subject == decl.SubjectConfiguration {
			return semantic.Result{}, fmt.Errorf("%w: %w", semantic.ErrInvalidRequest, err)
		}
		diagnostic.Code, diagnostic.Subject = n.Key(string(detail.Code)), n.Key(string(detail.Subject))
		if detail.Code == decl.MissingDependency || detail.Code == decl.Exhausted || detail.Code == decl.UnprovedObligation {
			status, verdict = semantic.Incomplete, semantic.Undetermined
		}
		if decoded != nil {
			diagnostic.At = decoded.Locate(detail)
		}
	case errors.As(err, &assertion):
		diagnostic.Code = n.Key("invalid-annotation")
		if assertion.Context {
			diagnostic.Subject = n.Key("context")
			if !assertion.Annotation {
				diagnostic.Code = n.Key("invalid-context")
			}
		}
		site := assertion.Site
		diagnostic.At = &site
	case errors.As(err, &syntax):
		diagnostic.At = syntax.Node
		if syntax.Code == text.LimitExceeded {
			status, verdict = semantic.Incomplete, semantic.Undetermined
			diagnostic.Code, diagnostic.Subject = n.Key("exhausted"), n.Key("request")
		}
	default:
		return semantic.Result{}, err
	}
	coverage := f.reports(request, status, verdict, false)
	if detail != nil && (detail.Code == decl.UnknownImport || detail.Code == decl.UnknownDeclaration) {
		for i := range coverage {
			if coverage[i].Obligation == n.ReferenceResolution() {
				coverage[i].Status, coverage[i].Verdict = semantic.Complete, semantic.Rejected
			}
		}
	}
	diagnostics := []semantic.Diagnostic{diagnostic}
	if decoded != nil && detail != nil {
		if site := decoded.CalleeDeclarationSite(detail); site != nil {
			subject := n.Key("source")
			target := detail.Target
			if detail.CalleeParameter.Declaration.Name != "" {
				target = detail.CalleeParameter.Declaration
			}
			root := decoded.tree.Nodes[decoded.tree.Root-1].References[0].Target.Node
			if (reader[L]{tree: decoded.tree}).attr(root, s.AttributePin) != target.Module.String() {
				subject = n.Key("context")
			}
			diagnostics = append(diagnostics, semantic.Diagnostic{Code: n.Key("callee-declaration"), Subject: subject, At: site, Message: "callee declaration or parameter contract"})
		}
	}
	return request.Publish(semantic.Publication{Reports: coverage, Diagnostics: diagnostics})
}

// ResolveReferences resolves the explicit entry Program against all supplied
// Program trees. Sites remain qualified by their original tree throughout.
func (f Frontend[L]) ResolveReferences(ctx context.Context, request semantic.Request[L]) (semantic.Result, error) {
	if ctx == nil {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if request.Contract().Profile != n.ResolutionProfile {
		return semantic.Result{}, fmt.Errorf("%w: unsupported Horos resolution profile", semantic.ErrInvalidRequest)
	}
	inputs := request.Inputs()
	entryFound := false
	for _, input := range inputs {
		if input.ID == request.Entry().Tree && input.Root == request.Entry().Node {
			entryFound = true
		}
	}
	if !entryFound {
		return semantic.Result{}, fmt.Errorf("%w: Horos resolution requires Programs and an explicit root entry", semantic.ErrInvalidRequest)
	}
	limits, err := f.configuredLimits(request.Configuration())
	if err != nil {
		return semantic.Result{}, err
	}
	decoded, err := f.DecodeCollection(inputs, request.Entry(), limits)
	if err != nil {
		return f.failure(request, nil, err)
	}
	program, err := decoded.Resolve(limits)
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if err != nil {
		return f.failure(request, decoded, err)
	}
	facts, err := f.project(ctx, decoded, program, request.Limits().Semantics)
	if err != nil {
		if !errors.Is(err, semantics.ErrLimit) {
			return semantic.Result{}, err
		}
		return request.Publish(semantic.Publication{Reports: f.reports(request, semantic.Complete, semantic.Holds, false), Diagnostics: []semantic.Diagnostic{{Code: n.Key("export-exhausted"), Subject: n.Key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	result, err := request.Publish(semantic.Publication{Facts: &facts, Reports: f.reports(request, semantic.Complete, semantic.Holds, true)})
	if errors.Is(err, semantics.ErrLimit) {
		return request.Publish(semantic.Publication{Reports: f.reports(request, semantic.Complete, semantic.Holds, false), Diagnostics: []semantic.Diagnostic{{Code: n.Key("export-exhausted"), Subject: n.Key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	return result, err
}
