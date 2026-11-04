package semantic

import (
	"errors"
	"fmt"
	"slices"
	"unicode/utf8"

	"github.com/Bitspark/bittype/internal/tree"
	"github.com/Bitspark/bittype/semantics/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// capture is also the invocation identity. It is private, nonempty storage;
// equal entry addresses in different captures do not establish compatibility.
type capture struct {
	source        *snapshotState
	operation     semantics.Key
	dependencies  []Dependency
	contract      Contract
	requested     []semantics.Key
	configuration []Setting
	limits        Limits
}

func (c *capture) contains(site ast.NodeRef) bool {
	return site.Tree != "" && site.Node > 0 && site.Node <= c.source.catalogue[site.Tree]
}

type requestState[L any] struct {
	capture *capture
	inputs  []ast.AST[L]
}

// Request owns a captured syntax/context snapshot. Its zero value is invalid.
// Accessors return owned copies. No accessor fetches dependencies or resolves names.
type Request[L any] struct{ state *requestState[L] }

func (r Request[L]) Entry() ast.NodeRef {
	if r.state == nil {
		return ast.NodeRef{}
	}
	return r.state.capture.source.entry
}

func (r Request[L]) Inputs() []ast.AST[L] {
	if r.state == nil {
		return nil
	}
	result := make([]ast.AST[L], len(r.state.inputs))
	for i, input := range r.state.inputs {
		result[i] = tree.Clone(input)
	}
	return result
}

func (r Request[L]) Requested() []semantics.Key {
	if r.state == nil {
		return nil
	}
	return slices.Clone(r.state.capture.requested)
}

func (r Request[L]) Configuration() []Setting {
	if r.state == nil {
		return nil
	}
	return cloneSettings(r.state.capture.configuration)
}

func (r Request[L]) Contract() Contract {
	if r.state == nil {
		return Contract{}
	}
	return r.state.capture.contract.Clone()
}

func (r Request[L]) Limits() Limits {
	if r.state == nil {
		return Limits{}
	}
	return r.state.capture.limits
}

// Publish seals output while the provider owns its mutable construction storage.
// It does not attach this request's identity to a previously published result.
// Re-publication is allowed, but returning an earlier invocation's Result is not.
func (r Request[L]) Publish(publication Publication) (Result, error) {
	if r.state == nil {
		return Result{}, ErrInvalidRequest
	}
	return publish[L](r.state.capture, publication)
}

type resultState struct {
	owner        *capture
	graph        representation.Semantics
	hasSemantics bool
	reports      []Report
	diagnostics  []Diagnostic
}

// Result is an immutable request-bound publication. A valid result can contain
// diagnostics and reports without a model. Valid does not mean semantic validity.
type Result struct{ state *resultState }

func (r Result) Valid() bool        { return r.state != nil }
func (r Result) HasSemantics() bool { return r.state != nil && r.state.hasSemantics }
func (r Result) Semantics() (representation.Semantics, bool) {
	if !r.HasSemantics() {
		return representation.Semantics{}, false
	}
	return r.state.graph, true
}
func (r Result) Entry() ast.NodeRef {
	if r.state == nil {
		return ast.NodeRef{}
	}
	return r.state.owner.source.entry
}
func (r Result) Contract() Contract {
	if r.state == nil {
		return Contract{}
	}
	return r.state.owner.contract.Clone()
}
func (r Result) Requested() []semantics.Key {
	if r.state == nil {
		return nil
	}
	return slices.Clone(r.state.owner.requested)
}
func (r Result) Configuration() []Setting {
	if r.state == nil {
		return nil
	}
	return cloneSettings(r.state.owner.configuration)
}
func (r Result) Limits() Limits {
	if r.state == nil {
		return Limits{}
	}
	return r.state.owner.limits
}
func (r Result) Reports() []Report {
	if r.state == nil {
		return nil
	}
	return slices.Clone(r.state.reports)
}
func (r Result) Report(key semantics.Key) (Report, bool) {
	if r.state != nil {
		for _, report := range r.state.reports {
			if report.Obligation == key {
				return report, true
			}
		}
	}
	return Report{}, false
}
func (r Result) Diagnostics() []Diagnostic {
	if r.state == nil {
		return nil
	}
	return cloneDiagnostics(r.state.diagnostics)
}

func cloneSettings(input []Setting) []Setting {
	result := slices.Clone(input)
	for i := range result {
		result[i].Value = slices.Clone(result[i].Value)
	}
	return result
}

func cloneDiagnostics(input []Diagnostic) []Diagnostic {
	result := slices.Clone(input)
	for i := range result {
		if result[i].At != nil {
			site := *result[i].At
			result[i].At = &site
		}
	}
	return result
}

func publish[L any](c *capture, p Publication) (Result, error) {
	bad := func(reason string) (Result, error) {
		return Result{}, fmt.Errorf("%w: %s", ErrInvalidResult, reason)
	}
	if p.Facts != nil && p.ReuseSemanticsFrom != nil {
		return bad("new facts and reused model are mutually exclusive")
	}
	known := make(map[semantics.Key]bool, len(c.contract.Obligations))
	allowed := make(map[semantics.Key]bool)
	for _, obligation := range c.contract.Obligations {
		known[obligation.Key] = true
		allowed[obligation.Key] = true
	}
	for _, key := range c.requested {
		allowed[key] = true
	}
	if len(p.Reports) > len(allowed) {
		return bad("too many reports")
	}
	reports := make(map[semantics.Key]Report, len(p.Reports))
	for _, report := range p.Reports {
		if _, exists := reports[report.Obligation]; exists {
			return bad("duplicate obligation report")
		}
		if !allowed[report.Obligation] || report.Status < Complete || report.Status > Unsupported || report.Verdict > Rejected {
			return bad("invalid obligation report")
		}
		if report.Status != Complete && report.Verdict != Undetermined {
			return bad("incomplete or unsupported report has a verdict")
		}
		if !known[report.Obligation] && report.Status != Unsupported {
			return bad("unadvertised obligation must be unsupported")
		}
		reports[report.Obligation] = report
	}
	for _, requested := range c.requested {
		if _, exists := reports[requested]; !exists {
			return bad("requested obligation was not reported")
		}
	}
	if len(p.Diagnostics) > c.limits.MaxDiagnostics {
		return Result{}, ErrLimit
	}
	bytes := c.limits.MaxDiagnosticBytes
	for _, diagnostic := range p.Diagnostics {
		if !diagnostic.Code.Valid() || !diagnostic.Subject.Valid() || !utf8.ValidString(diagnostic.Message) {
			return bad("invalid diagnostic")
		}
		if !keyBytes(&bytes, diagnostic.Code) || !keyBytes(&bytes, diagnostic.Subject) || !take(&bytes, len(diagnostic.Message)) {
			return Result{}, ErrLimit
		}
		if diagnostic.At != nil {
			if !c.contains(*diagnostic.At) {
				return bad("diagnostic site is outside captured inputs")
			}
			if !take(&bytes, len(diagnostic.At.Tree)) {
				return Result{}, ErrLimit
			}
		}
	}
	result := &resultState{owner: c}
	if p.Facts != nil || p.ReuseSemanticsFrom != nil {
		for _, dependency := range c.dependencies {
			if !dependency.reportsEstablished() {
				return bad("dependency requirements are not established")
			}
		}
		for _, requirement := range c.contract.SemanticsRequirements {
			report, found := reports[requirement.Obligation]
			if !found || report.Status != Complete || report.Verdict != requirement.Verdict {
				return bad("model requirements are not established")
			}
		}
		var model representation.Semantics
		var err error
		if p.ReuseSemanticsFrom != nil {
			admitted := false
			for _, dependency := range c.dependencies {
				admitted = admitted || dependency.Result.state == p.ReuseSemanticsFrom.state
			}
			if !admitted || !p.ReuseSemanticsFrom.HasSemantics() {
				return bad("model is not an admitted dependency")
			}
			model = p.ReuseSemanticsFrom.state.graph
			err = model.CheckLimits(c.limits.Semantics)
		} else {
			var graph semantics.Semantics[L]
			graph, err = semantics.New[L](*p.Facts, c.limits.Semantics, c.contains)
			if err == nil {
				model, err = representation.PackSemantics(c.source.syntax, graph)
			}
		}
		if err != nil {
			if errors.Is(err, semantics.ErrLimit) {
				return Result{}, fmt.Errorf("%w: %w", ErrLimit, err)
			}
			return Result{}, fmt.Errorf("%w: %w", ErrInvalidResult, err)
		}
		result.graph, result.hasSemantics = model, true
	}
	result.reports = slices.Clone(p.Reports)
	result.diagnostics = cloneDiagnostics(p.Diagnostics)
	return Result{state: result}, nil
}
