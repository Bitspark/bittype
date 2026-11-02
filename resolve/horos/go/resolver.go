// Package resolve adapts Horos declaration semantics to shared semantic facts.
// It resolves no names itself and never interprets descriptor constructors.
package resolve

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
	bridge "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

const Profile = "horos-declarations-1/resolution-1"

func ReferenceResolution() semantics.Key { return key("reference-resolution") }

func key(name string) semantics.Key        { return semantics.Key{Owner: "horos", Name: name} }
func Validity() semantics.Key              { return key("program-validity") }
func Indexing() semantics.Key              { return key("declaration-index") }
func Bindings() semantics.Key              { return key("occurrence-bindings") }
func Applications() semantics.Key          { return key("applications") }
func Dependencies() semantics.Key          { return key("application-dependencies") }
func TargetRole() semantics.Key            { return key("target") }
func ExpressionRole() semantics.Key        { return key("expression") }
func ApplicationDependency() semantics.Key { return key("declaration-application") }

// LimitsSetting contains a JSON array [depth,nodes,bytes,modules]. All four
// positive limits are required; depth must not exceed 512. Its absence selects
// Horos's fixed defaults.
// This bounds authoritative resolution separately from shared graph publication.
func LimitsSetting() semantics.Key { return key("declaration-limits-json-1") }

func Contract() resolution.Contract {
	keys := []semantics.Key{ReferenceResolution(), Validity(), Indexing(), Bindings(), Applications(), Dependencies()}
	contract := semantic.Contract{Profile: Profile}
	for _, k := range keys {
		contract.Obligations = append(contract.Obligations, semantic.Obligation{Key: k, Domain: key("entry-transitive-module-closure")})
		verdict := semantic.Undetermined
		if k == Validity() || k == ReferenceResolution() {
			verdict = semantic.Holds
		}
		contract.ModelRequirements = append(contract.ModelRequirements, semantic.Requirement{Obligation: k, Verdict: verdict})
	}
	return resolution.Contract{Semantic: contract, Bindings: ReferenceResolution()}
}

func configuredLimits(settings []semantic.Setting) (decl.Limits, error) {
	limits := decl.DefaultLimits()
	if len(settings) == 0 {
		return limits, nil
	}
	if len(settings) != 1 || settings[0].Key != LimitsSetting() {
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

func reports(request semantic.Request[s.Language], status semantic.Completion, verdict semantic.Verdict, complete bool) []semantic.Report {
	var result []semantic.Report
	known := make(map[semantics.Key]bool)
	for _, obligation := range Contract().Semantic.Obligations {
		r := semantic.Report{Obligation: obligation.Key, Status: semantic.Incomplete}
		if complete {
			r.Status = semantic.Complete
		}
		if obligation.Key == ReferenceResolution() && status == semantic.Complete && verdict == semantic.Holds {
			r.Status, r.Verdict = semantic.Complete, semantic.Holds
		}
		if obligation.Key == Validity() {
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

func failure(request semantic.Request[s.Language], decoded *bridge.Decoded, err error) (semantic.Result, error) {
	status, verdict := semantic.Complete, semantic.Rejected
	diagnostic := semantic.Diagnostic{Code: key("invalid-source"), Subject: key("source"), Message: err.Error()}
	var detail *decl.Error
	var assertion *bridge.ValidationError
	var syntax *text.Diagnostic
	switch {
	case errors.As(err, &detail):
		if detail.Code == decl.InvalidOptions || detail.Subject == decl.SubjectConfiguration {
			return semantic.Result{}, fmt.Errorf("%w: %w", semantic.ErrInvalidRequest, err)
		}
		diagnostic.Code, diagnostic.Subject = key(string(detail.Code)), key(string(detail.Subject))
		if detail.Code == decl.MissingDependency || detail.Code == decl.Exhausted {
			status, verdict = semantic.Incomplete, semantic.Undetermined
		}
		if decoded != nil {
			diagnostic.At = decoded.Locate(detail)
		}
	case errors.As(err, &assertion):
		diagnostic.Code = key("invalid-annotation")
		if assertion.Context {
			diagnostic.Subject = key("context")
			if !assertion.Annotation {
				diagnostic.Code = key("invalid-context")
			}
		}
		site := assertion.Site
		diagnostic.At = &site
	case errors.As(err, &syntax):
		diagnostic.At = syntax.Node
		if syntax.Code == text.LimitExceeded {
			status, verdict = semantic.Incomplete, semantic.Undetermined
			diagnostic.Code, diagnostic.Subject = key("exhausted"), key("request")
		}
	default:
		return semantic.Result{}, err
	}
	coverage := reports(request, status, verdict, false)
	if detail != nil && (detail.Code == decl.UnknownImport || detail.Code == decl.UnknownDeclaration) {
		for i := range coverage {
			if coverage[i].Obligation == ReferenceResolution() {
				coverage[i].Status, coverage[i].Verdict = semantic.Complete, semantic.Rejected
			}
		}
	}
	return request.Publish(semantic.Publication{Reports: coverage, Diagnostics: []semantic.Diagnostic{diagnostic}})
}

// ResolveReferences requires one Program tree and its root as the explicit entry. The
// structural root reference selects the module whose reachable closure is checked.
// Missing annotations are allowed; present ones are assertions, never repaired.
func ResolveReferences(ctx context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	if ctx == nil {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if request.Contract().Profile != Profile {
		return semantic.Result{}, fmt.Errorf("%w: unsupported Horos resolution profile", semantic.ErrInvalidRequest)
	}
	inputs := request.Inputs()
	if len(inputs) != 1 || request.Entry() != (ast.NodeRef{Tree: inputs[0].ID, Node: inputs[0].Root}) {
		return semantic.Result{}, fmt.Errorf("%w: Horos resolution requires one Program and its root entry", semantic.ErrInvalidRequest)
	}
	limits, err := configuredLimits(request.Configuration())
	if err != nil {
		return semantic.Result{}, err
	}
	decoded, err := bridge.Decode(inputs[0], limits)
	if err != nil {
		return failure(request, nil, err)
	}
	program, err := decoded.Resolve(limits)
	if err := ctx.Err(); err != nil {
		return semantic.Result{}, err
	}
	if err != nil {
		return failure(request, decoded, err)
	}
	facts, err := project(ctx, decoded, program, request.Limits().Model)
	if err != nil {
		if !errors.Is(err, semantics.ErrLimit) {
			return semantic.Result{}, err
		}
		return request.Publish(semantic.Publication{Reports: reports(request, semantic.Complete, semantic.Holds, false), Diagnostics: []semantic.Diagnostic{{Code: key("export-exhausted"), Subject: key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	result, err := request.Publish(semantic.Publication{Facts: &facts, Reports: reports(request, semantic.Complete, semantic.Holds, true)})
	if errors.Is(err, semantics.ErrLimit) {
		return request.Publish(semantic.Publication{Reports: reports(request, semantic.Complete, semantic.Holds, false), Diagnostics: []semantic.Diagnostic{{Code: key("export-exhausted"), Subject: key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	return result, err
}
