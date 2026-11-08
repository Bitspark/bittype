// Package generationfixture is an independent, deliberately bounded source for
// the public generation SPI. It publishes an ordinary semantic graph through its
// own analysis and supplies only its own interpretation, reader and mapping
// library: the kernel walks its graph exactly as it walks horos.
package generationfixture

import (
	"context"
	"errors"
	"fmt"

	checkgo "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/semantics/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	target "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// interpretation classifies fixture declarations, names fixture expressions
// and judges occurrences from the fixture analysis facts alone.
type interpretation struct {
	captured g.Captured[SourceLanguage]
	graph    semantics.Semantics[SourceLanguage]
	language representation.Language
	// foreign qualifies the variant's Foreign selector in another profile.
	foreign representation.Language
	variant SourceVariant
	// rejected is the admission report's verdict, when the variant reads one.
	rejected bool
}

func (i interpretation) selector(key string) (g.Selector, error) {
	if key == i.variant.Foreign {
		return g.NewSelector(i.foreign, key)
	}
	return g.NewSelector(i.language, key)
}

func (i interpretation) Declaration(ctx context.Context, d g.Declaration[SourceLanguage]) (g.Selector, bool, error) {
	id, ok := d.In(i.captured)
	if !ok {
		return g.Selector{}, false, g.ErrScope
	}
	e, _ := i.graph.Entity(id)
	if e.Kind != declarationKind {
		return g.Selector{}, false, nil
	}
	key, err := i.selector("declaration/" + e.Name)
	return key, true, err
}
func (i interpretation) Expression(ctx context.Context, u g.Use[SourceLanguage]) (g.Selector, error) {
	id, ok := u.In(i.captured)
	if !ok {
		return g.Selector{}, g.ErrScope
	}
	t, _ := i.graph.Term(id)
	if t.Operator.Owner != declarationKind.Owner {
		return g.Selector{}, g.ErrInvalid
	}
	return i.selector("expression/" + t.Operator.Name)
}
func (i interpretation) JudgeDeclaration(ctx context.Context, d g.Declaration[SourceLanguage]) (g.Judgment, error) {
	id, ok := d.In(i.captured)
	if !ok {
		return g.Judgment{}, g.ErrScope
	}
	if i.rejected {
		return g.Judgment{Verdict: g.Rejected, Diagnostic: errors.New("the admission report rejects the program")}, nil
	}
	if len(i.graph.Facts(semantics.Subject{Entity: id}, refuted)) > 0 {
		return g.Judgment{Verdict: g.Rejected, Diagnostic: fmt.Errorf("fixture declaration %d is refuted", id)}, nil
	}
	return g.Judgment{Verdict: g.Admitted}, ctx.Err()
}

// JudgeUse: refuted is rejected; judged is admitted; neither is outside the
// coverage of the fixture analysis.
func (i interpretation) JudgeUse(ctx context.Context, u g.Use[SourceLanguage]) (g.Judgment, error) {
	id, ok := u.In(i.captured)
	if !ok {
		return g.Judgment{}, g.ErrScope
	}
	subject := semantics.Subject{Term: id}
	switch {
	case len(i.graph.Facts(subject, refuted)) > 0:
		return g.Judgment{Verdict: g.Rejected, Diagnostic: fmt.Errorf("fixture term %d is refuted", id)}, nil
	case len(i.graph.Facts(subject, zero)) > 0:
		return g.Judgment{}, nil
	case len(i.graph.Facts(subject, judged)) > 0:
		return g.Judgment{Verdict: g.Admitted}, nil
	default:
		return g.Judgment{Verdict: g.Undetermined}, nil
	}
}

// Reader is the fixture reading API for its mapping library.
type Reader interface {
	Name(g.Declaration[SourceLanguage]) (string, error)
	ParameterName(g.Binder[SourceLanguage]) (string, error)
}
type reader struct {
	captured g.Captured[SourceLanguage]
	graph    semantics.Semantics[SourceLanguage]
}

func (r reader) Name(d g.Declaration[SourceLanguage]) (string, error) {
	id, ok := d.In(r.captured)
	if !ok {
		return "", g.ErrScope
	}
	e, _ := r.graph.Entity(id)
	return e.Name, nil
}
func (r reader) ParameterName(b g.Binder[SourceLanguage]) (string, error) {
	id, ok := b.In(r.captured)
	if !ok {
		return "", g.ErrScope
	}
	e, _ := r.graph.Entity(id)
	return e.Name, nil
}

// SourceVariant varies the fixture's source integration for negative
// qualification. The zero value is the standard source.
type SourceVariant struct {
	// Foreign is a selector key the source qualifies in ForeignIn, which
	// defaults to another profile of the source language.
	Foreign   string
	ForeignIn representation.Language
	// Admission makes the contract require a second analysis report; a
	// Rejected report rejects every declaration.
	Admission bool
	// reuse, when set, keeps the first opened source and serves it for every
	// later capture, as a caching integration would.
	reuse *g.Source[SourceLanguage, Reader]
}

// Reusing is the variant that serves its first opened source for every later
// capture. That interpretation can only refuse the later runs' handles.
func (v SourceVariant) Reusing() SourceVariant {
	v.reuse = &g.Source[SourceLanguage, Reader]{}
	return v
}

const admissionProfile = "fixture-admission"

var admissible = key("admissible")

type sourceProvider struct{ variant SourceVariant }

func (p sourceProvider) Contract() transform.Contract {
	evidence := []semantic.EvidenceRequirement{Evidence()}
	if p.variant.Admission {
		evidence = append(evidence, semantic.EvidenceRequirement{Operation: analysis.Kind(), Profile: admissionProfile})
	}
	return transform.Contract{Evidence: evidence}
}
func (p sourceProvider) Open(ctx context.Context, captured g.Captured[SourceLanguage]) (g.Source[SourceLanguage, Reader], error) {
	if err := ctx.Err(); err != nil {
		return g.Source[SourceLanguage, Reader]{}, err
	}
	if p.variant.reuse != nil && p.variant.reuse.Interpretation != nil {
		return *p.variant.reuse, nil
	}
	foreign := p.variant.ForeignIn
	if foreign.Validate() != nil {
		var err error
		if foreign, err = representation.NewLanguage[SourceLanguage](SourceLanguage{}, "fixture-foreign"); err != nil {
			return g.Source[SourceLanguage, Reader]{}, err
		}
	}
	graph := captured.Graph()
	i := interpretation{captured: captured, graph: graph, language: captured.Request().Snapshot().Language(), foreign: foreign, variant: p.variant}
	if p.variant.Admission {
		receipt, ok := captured.Request().Evidence(1)
		if !ok {
			return g.Source[SourceLanguage, Reader]{}, semantic.ErrEvidenceUnavailable
		}
		report, ok := receipt.Report(admissible)
		if !ok || report.Status != semantic.Complete {
			return g.Source[SourceLanguage, Reader]{}, semantic.ErrEvidenceUnavailable
		}
		i.rejected = report.Verdict == semantic.Rejected
	}
	source := g.Source[SourceLanguage, Reader]{Interpretation: i, Reader: reader{captured, graph}}
	if p.variant.reuse != nil {
		*p.variant.reuse = source
	}
	return source, nil
}

// otherSource is a different source language identity.
type otherSource struct{}

func (otherSource) ID() language.ID { return "generation-fixture-other" }
func (otherSource) Name() string    { return "Other generation fixture" }

// otherMarker is a different Go marker for the fixture's language identity.
type otherMarker struct{}

// ForeignLanguages are languages that differ from the fixture's source language
// in exactly one part of its identity: "profile", "id" or "marker".
func (f *Fixture) ForeignLanguages() (map[string]representation.Language, error) {
	profile, err := representation.NewLanguage[SourceLanguage](SourceLanguage{}, "fixture-foreign")
	if err != nil {
		return nil, err
	}
	id, err := representation.NewLanguage[SourceLanguage](otherSource{}, f.captured.language.Profile())
	if err != nil {
		return nil, err
	}
	marker, err := representation.NewLanguage[otherMarker](SourceLanguage{}, f.captured.language.Profile())
	if err != nil {
		return nil, err
	}
	return map[string]representation.Language{"profile": profile, "id": id, "marker": marker}, nil
}

// admission publishes the second report a SourceVariant with Admission reads.
type admission struct {
	SourceLanguage
	verdict semantic.Verdict
}

func (a admission) Analyze(_ context.Context, req analysis.Request[SourceLanguage]) (semantic.Result, error) {
	return req.Publish(semantic.Publication{Reports: []semantic.Report{{Obligation: admissible, Status: semantic.Complete, Verdict: a.verdict}}})
}

// Mode selects a library variant for negative qualification.
type Mode int

const (
	Standard Mode = iota
	// StrayCapture maps the occurrence itself instead of one of its children.
	StrayCapture
	// ForeignHandle maps a handle kept from an earlier plan.
	ForeignHandle
	// PatternStray captures the Cell occurrence itself instead of its actual.
	PatternStray
	// RecipeStray maps a kept non-body use from Outer's declaration recipe.
	RecipeStray
)

type library struct {
	mode Mode
	kept *g.Use[SourceLanguage]
}

func (l library) Describe(from, to representation.Language) (g.MappingDescription[SourceLanguage, s.Language, Reader], error) {
	cursor := g.NewHelperKey("independent/cursor")
	person, _ := g.NewSelector(from, "expression/person-body")
	page, _ := g.NewSelector(from, "expression/page-body")
	cell, _ := g.NewSelector(from, "declaration/Cell")
	ignore, _ := g.NewSelector(from, "declaration/Ignore")
	shadow, _ := g.NewSelector(from, "declaration/Shadow")
	v := g.TypeSlot("V")
	return g.MappingDescription[SourceLanguage, s.Language, Reader]{
		Profile: "independent/1",
		Bindings: []g.BindingDescription[SourceLanguage, s.Language, Reader]{
			{ID: "person", Input: g.Pattern[SourceLanguage, Reader]{Selector: person}, Output: model.Record(model.Field("Name", false, model.Primitive("string")))},
			{ID: "page", Input: g.Pattern[SourceLanguage, Reader]{Selector: page}, Expand: func(ctx context.Context, q g.Expansion[SourceLanguage, s.Language, Reader]) (g.Expression[s.Language], error) {
				view, err := q.Use()
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				element, ok := view.Operand(0)
				if !ok {
					return g.Expression[s.Language]{}, g.ErrScope
				}
				switch l.mode {
				case StrayCapture:
					element = view.Use
				case ForeignHandle:
					if l.kept.Valid() {
						element = *l.kept
					} else {
						*l.kept = element
					}
				case RecipeStray:
					*l.kept = element
				}
				value, err := q.Mapped(element)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				return model.Record(model.Field("Items", false, model.Collection(value)), model.Field("Storage", false, g.Parameter[s.Language](1)), model.Field("Cursor", false, g.Helper[s.Language](cursor))), nil
			}},
			// Cell's library represents every Cell use natively; Cell's body is never opened.
			{ID: "cell", Input: g.Pattern[SourceLanguage, Reader]{Selector: cell, Slots: []g.Slot{v}, Capture: func(ctx context.Context, r Reader, view g.UseView[SourceLanguage]) ([]g.Capture[SourceLanguage], error) {
				if l.mode == PatternStray {
					return []g.Capture[SourceLanguage]{{Slot: v, Use: view.Use}}, nil
				}
				for _, a := range view.Actuals {
					if name, err := r.ParameterName(a.Binder); err != nil {
						return nil, err
					} else if name == "V" {
						return []g.Capture[SourceLanguage]{{Slot: v, Use: a.Use}}, nil
					}
				}
				return nil, g.ErrScope
			}}, Output: model.Collection(g.Mapped[s.Language](v))},
		},
		Declarations: []g.DeclarationRecipe[SourceLanguage, s.Language, Reader]{
			{ID: "named", Expand: func(ctx context.Context, q g.DeclarationRequest[SourceLanguage, s.Language, Reader]) (g.DeclarationDescription[s.Language], error) {
				r, err := q.Reader()
				if err != nil {
					return g.DeclarationDescription[s.Language]{}, err
				}
				view, err := q.Declaration()
				if err != nil {
					return g.DeclarationDescription[s.Language]{}, err
				}
				name, err := r.Name(view.Declaration)
				if err != nil {
					return g.DeclarationDescription[s.Language]{}, err
				}
				form := model.Form("alias")
				if name == "Person" || name == "Page" {
					form = model.Form("defined")
				}
				d := g.DeclarationDescription[s.Language]{Name: name, Form: form}
				for i := range view.Formals {
					d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: fmt.Sprintf("T%d", i)})
					d.Generic.Bindings = append(d.Generic.Bindings, g.Parameter[s.Language](i))
					d.Generic.Arguments = append(d.Generic.Arguments, g.Argument[s.Language](i))
				}
				if name == "Page" {
					d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: "Storage"})
					d.Generic.Arguments = append(d.Generic.Arguments, model.Collection(g.Argument[s.Language](0)))
				}
				body := view.Body
				if l.mode == RecipeStray && name == "Outer" && l.kept.Valid() {
					body = *l.kept
				}
				d.Body, err = q.Mapped(body)
				return d, err
			}},
			// Ignore omits its source parameter natively; uses must still admit it.
			{ID: "erase", Selector: ignore, Expand: erase("Ignore")},
			// Shadow's recipe never represents its body; the kernel still admits it.
			{ID: "erase-shadow", Selector: shadow, Expand: erase("Shadow")},
		},
		Helpers: []g.HelperRecipe[s.Language]{{Key: cursor, Declaration: g.NativeDeclaration[s.Language]{Name: "Cursor", Form: model.Form("alias"), Body: model.Primitive("string")}}},
	}, nil
}

// erase describes a declaration as string, omitting its parameter and body.
func erase(name string) func(context.Context, g.DeclarationRequest[SourceLanguage, s.Language, Reader]) (g.DeclarationDescription[s.Language], error) {
	return func(ctx context.Context, q g.DeclarationRequest[SourceLanguage, s.Language, Reader]) (g.DeclarationDescription[s.Language], error) {
		return g.DeclarationDescription[s.Language]{Name: name, Form: model.Form("alias"), Body: model.Primitive("string"), Generic: g.GenericConvention[s.Language]{Bindings: []g.Expression[s.Language]{model.Primitive("string")}}}, nil
	}
}

// Fixture is one captured fixture program and its Go target.
type Fixture struct {
	captured captured
	target   representation.Language
	kept     g.Use[SourceLanguage]
}

// Roots are the declarations the standard plan generates.
var Roots = []string{"Person", "Page", "Outer", "Nested", "People", "Cells"}

func New(ctx context.Context) (*Fixture, error) { return newFixture(ctx, program()) }

// NewSharedSites is the fixture whose analysis records the binder reference
// that Lend and Borrow share at the syntax of both declarations, so the graph
// does not say which site an occurrence of that term is.
func NewSharedSites(ctx context.Context) (*Fixture, error) {
	p := program()
	lend := p.data.Entities[p.names["Lend"]-1]
	borrow := p.data.Entities[p.names["Borrow"]-1]
	if lend.Body != borrow.Body {
		return nil, errors.New("Lend and Borrow no longer share a body")
	}
	shared := p.data.Terms[lend.Body-1].Arguments[0]
	for _, e := range []semantics.Entity{lend, borrow} {
		p.data.Expressions = append(p.data.Expressions, semantics.Expression{At: semantics.Occurrence{Site: e.Declarations[0], Role: key("body-argument")}, Term: shared})
	}
	return newFixture(ctx, p)
}

func newFixture(ctx context.Context, p *graph) (*Fixture, error) {
	captured, err := capture(ctx, p)
	if err != nil {
		return nil, err
	}
	native, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	if err != nil {
		return nil, err
	}
	return &Fixture{captured: captured, target: native}, nil
}
func (f *Fixture) Target() representation.Language { return f.target }

// Plan generates the standard roots with the standard library.
func (f *Fixture) Plan(ctx context.Context) (g.GenerationPlan, error) {
	return f.PlanRoots(ctx, Standard, Roots...)
}

// PlanRoots generates the named declarations through the public registry.
func (f *Fixture) PlanRoots(ctx context.Context, mode Mode, names ...string) (g.GenerationPlan, error) {
	return f.plan(ctx, mode, SourceVariant{}, f.Input(), names...)
}

// PlanSites generates the declarations at the given syntax sites.
func (f *Fixture) PlanSites(ctx context.Context, sites ...ast.NodeRef) (g.GenerationPlan, error) {
	return f.planAt(ctx, Standard, SourceVariant{}, f.Input(), sites)
}

// PlanSource generates the named declarations from input through the public
// registry, with a source variant.
func (f *Fixture) PlanSource(ctx context.Context, variant SourceVariant, input transform.Input, names ...string) (g.GenerationPlan, error) {
	return f.plan(ctx, Standard, variant, input, names...)
}

// PlanDirect generates the named declarations from input through the exported
// planner, without the registry, so a test reaches the public planning
// boundary itself.
func (f *Fixture) PlanDirect(ctx context.Context, variant SourceVariant, input transform.Input, names ...string) (g.GenerationPlan, error) {
	roots, err := f.roots(names)
	if err != nil {
		return g.GenerationPlan{}, err
	}
	d, err := library{kept: &f.kept}.Describe(f.captured.language, f.target)
	if err != nil {
		return g.GenerationPlan{}, err
	}
	checker, err := check.NewChecker[s.Language](f.target, checkgo.Provider{})
	if err != nil {
		return g.GenerationPlan{}, err
	}
	return g.Plan(ctx, g.PlanningInput[SourceLanguage, s.Language, Reader]{
		From: f.captured.language, Input: input, Source: sourceProvider{variant}, Roots: roots,
		Mappings:  []g.SelectedDescription[SourceLanguage, s.Language, Reader]{{ID: "fixture", Description: d}},
		Execution: g.Execution[s.Language]{Target: target.Provider{Package: "generated"}, Checker: checker, Language: f.target, Output: "independent", Limits: g.DefaultLimits()},
	})
}

// AdmissionInput is the fixture input with a second report over the same
// snapshot and graph, for a SourceVariant with Admission: Holds admits the
// program and Rejected refutes it.
func (f *Fixture) AdmissionInput(ctx context.Context, verdict semantic.Verdict) (transform.Input, error) {
	contract := analysis.Contract{Semantic: semantic.Contract{Profile: admissionProfile, Obligations: []semantic.Obligation{{Key: admissible, Domain: key("program")}}}}
	a, err := analysis.NewAnalyzer[SourceLanguage](f.captured.language, contract, admission{verdict: verdict})
	if err != nil {
		return transform.Input{}, err
	}
	report, err := a.Analyze(ctx, analysis.Input{Snapshot: f.captured.snapshot})
	if err != nil {
		return transform.Input{}, err
	}
	in := f.Input()
	in.Evidence = append(in.Evidence, report)
	return in, nil
}

// Site is the syntax site of a fixture entity, or the program node for "".
func (f *Fixture) Site(name string) ast.NodeRef {
	if name == "" {
		return ast.NodeRef{Tree: programTree, Node: 1}
	}
	return f.captured.sites[name]
}

// PlanInput generates the named declarations from caller-supplied input, for
// evidence and snapshot negative cases.
func (f *Fixture) PlanInput(ctx context.Context, input transform.Input, names ...string) (g.GenerationPlan, error) {
	return f.plan(ctx, Standard, SourceVariant{}, input, names...)
}

// Input is the fixture's admitted input.
func (f *Fixture) Input() transform.Input {
	return transform.Input{Source: f.captured.snapshot, Evidence: []semantic.Result{f.captured.evidence}}
}

func (f *Fixture) roots(names []string) ([]ast.NodeRef, error) {
	var roots []ast.NodeRef
	for _, name := range names {
		site, ok := f.captured.sites[name]
		if !ok {
			return nil, fmt.Errorf("unknown fixture declaration %q", name)
		}
		roots = append(roots, site)
	}
	return roots, nil
}

func (f *Fixture) plan(ctx context.Context, mode Mode, variant SourceVariant, input transform.Input, names ...string) (g.GenerationPlan, error) {
	roots, err := f.roots(names)
	if err != nil {
		return g.GenerationPlan{}, err
	}
	return f.planAt(ctx, mode, variant, input, roots)
}

func (f *Fixture) planAt(ctx context.Context, mode Mode, variant SourceVariant, input transform.Input, roots []ast.NodeRef) (g.GenerationPlan, error) {
	var registry c.Registry
	from := f.captured.language
	if err := c.RegisterSource[SourceLanguage, Reader](&registry, from, sourceProvider{variant}); err != nil {
		return g.GenerationPlan{}, err
	}
	if err := c.RegisterTarget(&registry, f.target, target.Provider{Package: "generated"}); err != nil {
		return g.GenerationPlan{}, err
	}
	if _, err := c.RegisterMapping(&registry, from, f.target, library{mode: mode, kept: &f.kept}); err != nil {
		return g.GenerationPlan{}, err
	}
	checker, err := check.NewChecker[s.Language](f.target, checkgo.Provider{})
	if err != nil {
		return g.GenerationPlan{}, err
	}
	return registry.PlanGeneration(ctx, c.Request{Input: input, To: f.target, Roots: roots, Output: "independent", Profile: "independent/1", Limits: g.DefaultLimits()}, checker)
}
