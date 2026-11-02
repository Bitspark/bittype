// Package analyze adapts Horos declaration semantics to shared analysis facts.
// It resolves no names itself and never interprets descriptor constructors.
package analyze

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	bridge "github.com/Bitspark/bittype/internal/horos"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

const Profile = "horos-declarations-1/analysis-1"

func key(name string) semantics.Key { return semantics.Key{Owner: "horos", Name: name} }
func Validity() semantics.Key { return key("program-validity") }
func Indexing() semantics.Key { return key("declaration-index") }
func Bindings() semantics.Key { return key("occurrence-bindings") }
func Applications() semantics.Key { return key("applications") }
func Dependencies() semantics.Key { return key("application-dependencies") }
func TargetRole() semantics.Key { return key("target") }
func ExpressionRole() semantics.Key { return key("expression") }
func ApplicationDependency() semantics.Key { return key("declaration-application") }

// LimitsSetting contains a JSON array [depth,nodes,bytes,modules]. All four
// positive limits are required. Its absence selects Horos's fixed defaults.
// This bounds authoritative resolution separately from shared graph publication.
func LimitsSetting() semantics.Key { return key("declaration-limits-json-1") }

func Contract() analysis.Contract {
	keys := []semantics.Key{Validity(), Indexing(), Bindings(), Applications(), Dependencies()}
	contract := analysis.Contract{Profile: Profile}
	for _, k := range keys {
		contract.Obligations = append(contract.Obligations, analysis.Obligation{Key: k, Domain: key("entry-transitive-module-closure")})
		verdict := analysis.Undetermined
		if k == Validity() { verdict = analysis.Holds }
		contract.ModelRequirements = append(contract.ModelRequirements, analysis.Requirement{Obligation: k, Verdict: verdict})
	}
	return contract
}

func configuredLimits(settings []analysis.Setting) (decl.Limits, error) {
	limits := decl.DefaultLimits()
	if len(settings) == 0 { return limits, nil }
	if len(settings) != 1 || settings[0].Key != LimitsSetting() { return decl.Limits{}, fmt.Errorf("%w: unsupported Horos configuration", analysis.ErrInvalidRequest) }
	decoder := json.NewDecoder(bytes.NewReader(settings[0].Value))
	var values []int
	if err := decoder.Decode(&values); err != nil || len(values) != 4 { return decl.Limits{}, fmt.Errorf("%w: expected four Horos limits", analysis.ErrInvalidRequest) }
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF { return decl.Limits{}, fmt.Errorf("%w: trailing Horos configuration", analysis.ErrInvalidRequest) }
	return decl.Limits{MaxDepth: values[0], MaxNodes: values[1], MaxBytes: values[2], MaxModules: values[3]}, nil
}

func reports(request analysis.Request[s.Language], status analysis.Completion, verdict analysis.Verdict, complete bool) []analysis.Report {
	var result []analysis.Report
	known := make(map[semantics.Key]bool)
	for _, obligation := range Contract().Obligations {
		r := analysis.Report{Obligation: obligation.Key, Status: analysis.Incomplete}
		if complete { r.Status = analysis.Complete }
		if obligation.Key == Validity() { r.Status, r.Verdict = status, verdict }
		result = append(result, r)
		known[obligation.Key] = true
	}
	for _, requested := range request.Requested() {
		if !known[requested] { result = append(result, analysis.Report{Obligation: requested, Status: analysis.Unsupported}) }
	}
	return result
}

func failure(request analysis.Request[s.Language], decoded *bridge.Decoded, err error) (analysis.Result, error) {
	status, verdict := analysis.Complete, analysis.Rejected
	diagnostic := analysis.Diagnostic{Code: key("invalid-source"), Subject: key("source"), Message: err.Error()}
	var detail *decl.Error
	var assertion *bridge.ValidationError
	var syntax *text.Diagnostic
	switch {
	case errors.As(err, &detail):
		if detail.Code == decl.InvalidOptions || detail.Subject == decl.SubjectConfiguration { return analysis.Result{}, fmt.Errorf("%w: %w", analysis.ErrInvalidRequest, err) }
		diagnostic.Code, diagnostic.Subject = key(string(detail.Code)), key(string(detail.Subject))
		if detail.Code == decl.MissingDependency || detail.Code == decl.Exhausted { status, verdict = analysis.Incomplete, analysis.Undetermined }
		if decoded != nil { diagnostic.At = decoded.Locate(detail) }
	case errors.As(err, &assertion):
		diagnostic.Code = key("invalid-annotation")
		if assertion.Context { diagnostic.Code, diagnostic.Subject = key("invalid-context"), key("context") }
		site := assertion.Site
		diagnostic.At = &site
	case errors.As(err, &syntax):
		diagnostic.At = syntax.Node
		if syntax.Code == text.LimitExceeded {
			status, verdict = analysis.Incomplete, analysis.Undetermined
			diagnostic.Code, diagnostic.Subject = key("exhausted"), key("request")
		}
	default:
		return analysis.Result{}, err
	}
	return request.Publish(analysis.Publication{Reports: reports(request, status, verdict, false), Diagnostics: []analysis.Diagnostic{diagnostic}})
}

// Analyze requires one Program tree and its root as the explicit entry. The
// structural root reference selects the module whose reachable closure is checked.
// Missing annotations are allowed; present ones are assertions, never repaired.
func Analyze(ctx context.Context, request analysis.Request[s.Language]) (analysis.Result, error) {
	if ctx == nil { return analysis.Result{}, analysis.ErrInvalidRequest }
	if err := ctx.Err(); err != nil { return analysis.Result{}, err }
	if request.Contract().Profile != Profile { return analysis.Result{}, fmt.Errorf("%w: unsupported Horos analysis profile", analysis.ErrInvalidRequest) }
	inputs := request.Inputs()
	if len(inputs) != 1 || request.Entry() != (ast.NodeRef{Tree: inputs[0].ID, Node: inputs[0].Root}) {
		return analysis.Result{}, fmt.Errorf("%w: Horos analysis requires one Program and its root entry", analysis.ErrInvalidRequest)
	}
	limits, err := configuredLimits(request.Configuration())
	if err != nil { return analysis.Result{}, err }
	decoded, err := bridge.Decode(inputs[0], limits)
	if err != nil { return failure(request, nil, err) }
	program, err := decoded.Resolve(limits)
	if err := ctx.Err(); err != nil { return analysis.Result{}, err }
	if err != nil { return failure(request, decoded, err) }
	facts, err := project(ctx, decoded, program, request.Limits().Model)
	if err != nil {
		if !errors.Is(err, semantics.ErrLimit) { return analysis.Result{}, err }
		return request.Publish(analysis.Publication{Reports: reports(request, analysis.Complete, analysis.Holds, false), Diagnostics: []analysis.Diagnostic{{Code: key("export-exhausted"), Subject: key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	result, err := request.Publish(analysis.Publication{Facts: &facts, Reports: reports(request, analysis.Complete, analysis.Holds, true)})
	if errors.Is(err, semantics.ErrLimit) {
		return request.Publish(analysis.Publication{Reports: reports(request, analysis.Complete, analysis.Holds, false), Diagnostics: []analysis.Diagnostic{{Code: key("export-exhausted"), Subject: key("request"), Message: "Horos semantic export exceeds its resource budget"}}})
	}
	return result, err
}
