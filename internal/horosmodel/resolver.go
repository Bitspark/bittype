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

	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

func (f Frontend[L]) ReferenceResolution() semantics.Key { return f.key("reference-resolution") }

func (f Frontend[L]) key(name string) semantics.Key        { return semantics.Key{Owner: f.Owner, Name: name} }
func (f Frontend[L]) Validity() semantics.Key              { return f.key("program-validity") }
func (f Frontend[L]) Indexing() semantics.Key              { return f.key("declaration-index") }
func (f Frontend[L]) Bindings() semantics.Key              { return f.key("occurrence-bindings") }
func (f Frontend[L]) Applications() semantics.Key          { return f.key("applications") }
func (f Frontend[L]) Dependencies() semantics.Key          { return f.key("reference-dependencies") }
func (f Frontend[L]) TypeFacts() semantics.Key             { return f.key("type-facts") }
func (f Frontend[L]) TargetRole() semantics.Key            { return f.key("target") }
func (f Frontend[L]) ExpressionRole() semantics.Key        { return f.key("expression") }
func (f Frontend[L]) ReferenceDependency() semantics.Key   { return f.key("declaration-reference") }
func (f Frontend[L]) KindRelation() semantics.Key          { return f.key("kind") }
func (f Frontend[L]) InterfaceRelation() semantics.Key     { return f.key("interface") }
func (f Frontend[L]) DataStatusRelation() semantics.Key    { return f.key("data-status") }
func (f Frontend[L]) DataSummaryRelation() semantics.Key   { return f.key("data-summary") }
func (f Frontend[L]) IdentityRelation() semantics.Key      { return f.key("canonical-identity") }

// LimitsSetting contains a JSON array [depth,nodes,bytes,modules]. All four
// positive limits are required; depth must not exceed 512. Its absence selects
// Horos's fixed defaults.
// This bounds authoritative resolution separately from shared graph publication.
func (f Frontend[L]) LimitsSetting() semantics.Key { return f.key("declaration-limits-json-1") }

func (f Frontend[L]) Contract() resolution.Contract {
	keys := []semantics.Key{f.ReferenceResolution(), f.Validity(), f.Indexing(), f.Bindings(), f.Applications(), f.Dependencies(), f.TypeFacts()}
	contract := semantic.Contract{Profile: f.ResolutionProfile}
	for _, k := range keys {
		contract.Obligations = append(contract.Obligations, semantic.Obligation{Key: k, Domain: f.key("entry-transitive-module-closure")})
		verdict := semantic.Undetermined
		if k == f.Validity() || k == f.ReferenceResolution() {
			verdict = semantic.Holds
		}
		contract.ModelRequirements = append(contract.ModelRequirements, semantic.Requirement{Obligation: k, Verdict: verdict})
	}
	return resolution.Contract{Semantic: contract, Bindings: f.ReferenceResolution()}
}

func (f Frontend[L]) configuredLimits(settings []semantic.Setting) (decl.Limits, error) {
	limits := decl.DefaultLimits()
	if len(settings) == 0 {
		return limits, nil
	}
	if len(settings) != 1 || settings[0].Key != f.LimitsSetting() {
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
	for _, obligation := range f.Contract().Semantic.Obligations {
		r := semantic.Report{Obligation: obligation.Key, Status: semantic.Incomplete}
		if complete {
			r.Status = semantic.Complete
		}
		if obligation.Key == f.ReferenceResolution() && status == semantic.Complete && verdict == semantic.Holds {
			r.Status, r.Verdict = semantic.Complete, semantic.Holds
		}
		if obligation.Key == f.Validity() {
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
	diagnostic := semantic.Diagnostic{Code: f.key("invalid-source"), Subject: f.key("source"), Message: err.Error()}
	var detail *decl.Error
	var assertion *ValidationError
	var syntax *text.Diagnostic
	switch {
	case errors.As(err, &detail):
		if detail.Code == decl.InvalidOptions || detail.Subject == decl.SubjectConfiguration {
			return semantic.Result{}, fmt.Errorf("%w: %w", semantic.ErrInvalidRequest, err)
		}
		diagnostic.Code, diagnostic.Subject = f.key(string(detail.Code)), f.key(string(detail.Subject))
		if detail.Code == decl.MissingDependency || detail.Code == decl.Exhausted {
			status, verdict = semantic.Incomplete, semantic.Undetermined
		}
		if decoded != nil {
			diagnostic.At = decoded.Locate(detail)
		}
	case errors.As(err, &assertion):
		diagnostic.Code = f.key("invalid-annotation")
		if assertion.Context {
			diagnostic.Subject = f.key("context")
			if !assertion.Annotation {
				diagnostic.Code = f.key("invalid-context")
			}
		}
		site := assertion.Site
		diagnostic.At = &site
	case errors.As(err, &syntax):
		diagnostic.At = syntax.Node
		if syntax.Code == text.LimitExceeded {
			status, verdict = semantic.Incomplete, semantic.Undetermined
			diagnostic.Code, diagnostic.Subject = f.key("exhausted"), f.key("request")
		}
	default:
		return semantic.Result{}, err
	}
	coverage := f.reports(request, status, verdict, false)
	if detail != nil && (detail.Code == decl.UnknownImport || detail.Code == decl.UnknownDeclaration) {
		for i := range coverage {
			if coverage[i].Obligation == f.ReferenceResolution() {
				coverage[i].Status, coverage[i].Verdict = semantic.Complete, semantic.Rejected
			}
		}
	}
	diagnostics := []semantic.Diagnostic{diagnostic}
	if decoded != nil && detail != nil {
		if site, context := decoded.CalleeDeclarationSite(detail); site != nil {
			subject := f.key("source")
			if context {
				subject = f.key("context")
			}
			diagnostics = append(diagnostics, semantic.Diagnostic{Code: f.key("callee-declaration"), Subject: subject, At: site, Message: "The callee's declared interface is established here."})
		}
	}
	return request.Publish(semantic.Publication{Reports: coverage, Diagnostics: diagnostics})
}

// ResolveReferences requires one Program tree and its root as the explicit entry. The
// structural root reference selects the module whose reachable closure is checked.
// Missing annotations are allowed; present ones are assertions, never repaired.
func (f Frontend[L]) ResolveReferences(ctx context.Context, request semantic.Request[L]) (semantic.Result, error) {
	if ctx == nil {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if request.Contract().Profile != f.ResolutionProfile {
		return semantic.Result{}, fmt.Errorf("%w: unsupported Horos resolution profile", semantic.ErrInvalidRequest)
	}
	inputs := request.Inputs()
	if len(inputs) != 1 || request.Entry() != (ast.NodeRef{Tree: inputs[0].ID, Node: inputs[0].Root}) {
		return semantic.Result{}, fmt.Errorf("%w: Horos resolution requires one Program and its root entry", semantic.ErrInvalidRequest)
	}
	limits, err := f.configuredLimits(request.Configuration())
	if err != nil {
		return semantic.Result{}, err
	}
	decoded, err := f.Decode(inputs[0], limits)
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
	facts, err := f.project(ctx, decoded, program, request.Limits().Model)
	if err != nil {
		if !errors.Is(err, semantics.ErrLimit) {
			return semantic.Result{}, err
		}
		return request.Publish(semantic.Publication{Reports: f.reports(request, semantic.Complete, semantic.Holds, false), Diagnostics: []semantic.Diagnostic{{Code: f.key("export-exhausted"), Subject: f.key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	result, err := request.Publish(semantic.Publication{Facts: &facts, Reports: f.reports(request, semantic.Complete, semantic.Holds, true)})
	if errors.Is(err, semantics.ErrLimit) {
		return request.Publish(semantic.Publication{Reports: f.reports(request, semantic.Complete, semantic.Holds, false), Diagnostics: []semantic.Diagnostic{{Code: f.key("export-exhausted"), Subject: f.key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	return result, err
}
