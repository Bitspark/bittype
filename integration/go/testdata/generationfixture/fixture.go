// Package generationfixture is an independent, deliberately bounded source for
// the public generation SPI. It publishes an ordinary semantic graph through its
// own analysis and supplies only its own interpretation, reader and mapping
// library: the kernel walks its graph exactly as it walks horos.
package generationfixture

import (
	"context"
	"fmt"

	checkgo "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
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
	graph    semantics.Semantics[SourceLanguage]
	language representation.Language
}

func (i interpretation) Declaration(ctx context.Context, d g.Declaration[SourceLanguage]) (g.Selector, bool, error) {
	id, ok := d.In(i.graph)
	if !ok {
		return g.Selector{}, false, g.ErrScope
	}
	e, _ := i.graph.Entity(id)
	if e.Kind != declarationKind {
		return g.Selector{}, false, nil
	}
	key, err := g.NewSelector(i.language, "declaration/"+e.Name)
	return key, true, err
}
func (i interpretation) Expression(ctx context.Context, u g.Use[SourceLanguage]) (g.Selector, error) {
	id, ok := u.In(i.graph)
	if !ok {
		return g.Selector{}, g.ErrScope
	}
	t, _ := i.graph.Term(id)
	if t.Operator.Owner != declarationKind.Owner {
		return g.Selector{}, g.ErrInvalid
	}
	return g.NewSelector(i.language, "expression/"+t.Operator.Name)
}
func (i interpretation) JudgeDeclaration(ctx context.Context, d g.Declaration[SourceLanguage]) (g.Judgment, error) {
	return g.Judgment{Verdict: g.Admitted}, ctx.Err()
}

// JudgeUse: refuted is rejected; judged is admitted; neither is outside the
// coverage of the fixture analysis.
func (i interpretation) JudgeUse(ctx context.Context, u g.Use[SourceLanguage]) (g.Judgment, error) {
	id, ok := u.In(i.graph)
	if !ok {
		return g.Judgment{}, g.ErrScope
	}
	subject := semantics.Subject{Term: id}
	switch {
	case len(i.graph.Facts(subject, refuted)) > 0:
		return g.Judgment{Verdict: g.Rejected, Diagnostic: fmt.Errorf("fixture term %d is refuted", id)}, nil
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
	graph semantics.Semantics[SourceLanguage]
}

func (r reader) Name(d g.Declaration[SourceLanguage]) (string, error) {
	id, ok := d.In(r.graph)
	if !ok {
		return "", g.ErrScope
	}
	e, _ := r.graph.Entity(id)
	return e.Name, nil
}
func (r reader) ParameterName(b g.Binder[SourceLanguage]) (string, error) {
	id, ok := b.In(r.graph)
	if !ok {
		return "", g.ErrScope
	}
	e, _ := r.graph.Entity(id)
	return e.Name, nil
}

type sourceProvider struct{}

func (sourceProvider) Contract() transform.Contract {
	return transform.Contract{Evidence: []semantic.EvidenceRequirement{Evidence()}}
}
func (sourceProvider) Open(ctx context.Context, captured g.Captured[SourceLanguage]) (g.Source[SourceLanguage, Reader], error) {
	if err := ctx.Err(); err != nil {
		return g.Source[SourceLanguage, Reader]{}, err
	}
	graph := captured.Graph()
	return g.Source[SourceLanguage, Reader]{Interpretation: interpretation{graph, captured.Request().Snapshot().Language()}, Reader: reader{graph}}, nil
}

// Mode selects a library variant for negative qualification.
type Mode int

const (
	Standard Mode = iota
	// StrayCapture maps the occurrence itself instead of one of its children.
	StrayCapture
	// ForeignHandle maps a handle kept from an earlier plan.
	ForeignHandle
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
				}
				value, err := q.Mapped(element)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				return model.Record(model.Field("Items", false, model.Collection(value)), model.Field("Storage", false, g.Parameter[s.Language](1)), model.Field("Cursor", false, g.Helper[s.Language](cursor))), nil
			}},
			// Cell's library represents every Cell use natively; Cell's body is never opened.
			{ID: "cell", Input: g.Pattern[SourceLanguage, Reader]{Selector: cell, Slots: []g.Slot{v}, Capture: func(ctx context.Context, r Reader, view g.UseView[SourceLanguage]) ([]g.Capture[SourceLanguage], error) {
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
				d.Body, err = q.Mapped(view.Body)
				return d, err
			}},
			// Ignore omits its source parameter natively; uses must still admit it.
			{ID: "erase", Selector: ignore, Expand: func(ctx context.Context, q g.DeclarationRequest[SourceLanguage, s.Language, Reader]) (g.DeclarationDescription[s.Language], error) {
				return g.DeclarationDescription[s.Language]{Name: "Ignore", Form: model.Form("alias"), Body: model.Primitive("string"), Generic: g.GenericConvention[s.Language]{Bindings: []g.Expression[s.Language]{model.Primitive("string")}}}, nil
			}},
		},
		Helpers: []g.HelperRecipe[s.Language]{{Key: cursor, Declaration: g.NativeDeclaration[s.Language]{Name: "Cursor", Form: model.Form("alias"), Body: model.Primitive("string")}}},
	}, nil
}

// Fixture is one captured fixture program and its Go target.
type Fixture struct {
	captured captured
	target   representation.Language
	kept     g.Use[SourceLanguage]
}

// Roots are the declarations the standard plan generates.
var Roots = []string{"Person", "Page", "Outer", "Nested", "People", "Cells"}

func New(ctx context.Context) (*Fixture, error) {
	captured, err := capture(ctx, program())
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
	return f.plan(ctx, mode, transform.Input{Source: f.captured.snapshot, Evidence: []semantic.Result{f.captured.evidence}}, names...)
}

// PlanInput generates the named declarations from caller-supplied input, for
// evidence and snapshot negative cases.
func (f *Fixture) PlanInput(ctx context.Context, input transform.Input, names ...string) (g.GenerationPlan, error) {
	return f.plan(ctx, Standard, input, names...)
}

// Input is the fixture's admitted input.
func (f *Fixture) Input() transform.Input {
	return transform.Input{Source: f.captured.snapshot, Evidence: []semantic.Result{f.captured.evidence}}
}

func (f *Fixture) plan(ctx context.Context, mode Mode, input transform.Input, names ...string) (g.GenerationPlan, error) {
	var registry c.Registry
	from := f.captured.language
	if err := c.RegisterSource[SourceLanguage, Reader](&registry, from, sourceProvider{}); err != nil {
		return g.GenerationPlan{}, err
	}
	if err := c.RegisterTarget(&registry, f.target, target.Provider{Package: "generated"}); err != nil {
		return g.GenerationPlan{}, err
	}
	if _, err := c.RegisterMapping(&registry, from, f.target, library{mode: mode, kept: &f.kept}); err != nil {
		return g.GenerationPlan{}, err
	}
	var roots []ast.NodeRef
	for _, name := range names {
		site, ok := f.captured.sites[name]
		if !ok {
			return g.GenerationPlan{}, fmt.Errorf("unknown fixture declaration %q", name)
		}
		roots = append(roots, site)
	}
	checker, err := check.NewChecker[s.Language](f.target, checkgo.Provider{})
	if err != nil {
		return g.GenerationPlan{}, err
	}
	return registry.PlanGeneration(ctx, c.Request{Input: input, To: f.target, Roots: roots, Output: "independent", Profile: "independent/1", Limits: g.DefaultLimits()}, checker)
}
