package generationfixture

import (
	"context"
	"fmt"

	checkgo "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	target "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Fixture struct {
	source       *Source
	names        []string
	declarations map[string]g.Declaration[SourceLanguage]
	target       representation.Language
}

func New(ctx context.Context) (*Fixture, error) {
	source, err := newSource(ctx)
	if err != nil {
		return nil, err
	}
	native, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	if err != nil {
		return nil, err
	}
	f := &Fixture{source: source, names: []string{"Person", "Page", "Outer", "Nested", "People"}, declarations: map[string]g.Declaration[SourceLanguage]{}, target: native}
	for _, name := range f.names {
		f.declarations[name] = source.declaration(name == "Page" || name == "Outer" || name == "Nested")
	}
	person, page := f.declarations["Person"], f.declarations["Page"]
	source.body(person, source.use(source.opens[person].Environment, "person-body", g.Declaration[SourceLanguage]{}, g.Binder[SourceLanguage]{}))
	source.body(page, source.use(source.opens[page].Environment, "page-body", g.Declaration[SourceLanguage]{}, g.Binder[SourceLanguage]{}))
	outer, nested, people := source.opens[f.declarations["Outer"]], source.opens[f.declarations["Nested"]], source.opens[f.declarations["People"]]
	source.body(outer.Declaration, source.call(outer.Environment, "Page", page, source.occurrence(outer.Formals[0].Subject)))
	inner := source.call(nested.Environment, "Page", page, source.occurrence(nested.Formals[0].Subject))
	source.body(nested.Declaration, source.call(nested.Environment, "Page", page, inner))
	personUse := source.call(people.Environment, "Person", person, g.Use[SourceLanguage]{})
	source.body(people.Declaration, source.call(people.Environment, "Page", page, personUse))
	return f, nil
}
func (f *Fixture) Target() representation.Language { return f.target }

type reader struct {
	*Source
	f *Fixture
}

func (r *reader) Declarations() []g.Declaration[SourceLanguage] {
	var out []g.Declaration[SourceLanguage]
	for _, name := range r.f.names {
		out = append(out, r.f.declarations[name])
	}
	return out
}
func (r *reader) DeclarationSelector(d g.Declaration[SourceLanguage]) (g.Selector, bool, error) {
	for name, value := range r.f.declarations {
		if value == d {
			key, err := g.NewSelector(r.snapshot.Language(), name)
			return key, true, err
		}
	}
	return g.Selector{}, false, g.ErrScope
}
func (r *reader) UseSelector(u g.UseData[SourceLanguage]) (g.Selector, error) {
	key, err := r.selectPlan(u)
	if err != nil {
		return g.Selector{}, err
	}
	return g.NewSelector(r.snapshot.Language(), key)
}
func (r *reader) ResolveRoots(ctx context.Context, sites []ast.NodeRef) ([]g.Declaration[SourceLanguage], error) {
	var out []g.Declaration[SourceLanguage]
	for _, site := range sites {
		if site.Tree != r.snapshot.Entry().Tree || site.Node < 2 || int(site.Node) > len(r.f.names)+1 {
			return nil, g.ErrScope
		}
		out = append(out, r.f.declarations[r.f.names[site.Node-2]])
	}
	return out, nil
}
func (r *Source) CheckChild(ctx context.Context, env g.Environment[SourceLanguage], parent, child g.Use[SourceLanguage]) error {
	p, err := r.Admit(ctx, env, parent)
	if err != nil {
		return err
	}
	if _, err = r.Admit(ctx, env, child); err != nil {
		return err
	}
	for _, a := range p.Actuals {
		if a.Use == child {
			return nil
		}
	}
	for _, open := range r.opens {
		if open.Environment == env && open.Body == parent {
			for _, formal := range open.Formals {
				if formal.Subject == child {
					return nil
				}
			}
		}
	}
	return g.ErrScope
}

type sourceProvider struct{ f *Fixture }

func (sourceProvider) Contract() transform.Contract { return transform.Contract{} }
func (p sourceProvider) Open(ctx context.Context, q transform.Request[SourceLanguage]) (*reader, error) {
	if !q.Snapshot().Same(p.f.source.snapshot) {
		return nil, g.ErrScope
	}
	return &reader{p.f.source, p.f}, nil
}

type library struct{}

func (library) Describe(from, to representation.Language) (g.MappingDescription[SourceLanguage, s.Language, *reader], error) {
	cursor := g.NewHelperKey("independent/cursor")
	person, _ := g.NewSelector(from, "person-body")
	page, _ := g.NewSelector(from, "page-body")
	return g.MappingDescription[SourceLanguage, s.Language, *reader]{
		Profile: "independent/1",
		Bindings: []g.BindingDescription[SourceLanguage, s.Language, *reader]{
			{ID: "person", Input: g.Pattern[SourceLanguage, *reader]{Selector: person}, Output: model.Record(model.Field("Name", false, model.Primitive("string")))},
			{ID: "page", Input: g.Pattern[SourceLanguage, *reader]{Selector: page}, Expand: func(ctx context.Context, q g.Expansion[SourceLanguage, s.Language, *reader]) (g.Expression[s.Language], error) {
				r, err := q.Source()
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				use, err := q.Use()
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				open := r.opens[r.f.declarations["Page"]]
				if use.Environment != open.Environment {
					return g.Expression[s.Language]{}, g.ErrScope
				}
				value, err := q.Mapped(open.Formals[0].Subject)
				if err != nil {
					return g.Expression[s.Language]{}, err
				}
				return model.Record(model.Field("Items", false, model.Collection(value)), model.Field("Storage", false, g.Parameter[s.Language](1)), model.Field("Cursor", false, g.Helper[s.Language](cursor))), nil
			}},
		},
		Declarations: []g.DeclarationRecipe[SourceLanguage, s.Language, *reader]{{ID: "named", Expand: func(ctx context.Context, q g.DeclarationRequest[SourceLanguage, s.Language, *reader]) (g.DeclarationDescription[s.Language], error) {
			r, err := q.Source()
			if err != nil {
				return g.DeclarationDescription[s.Language]{}, err
			}
			open, err := q.Open()
			if err != nil {
				return g.DeclarationDescription[s.Language]{}, err
			}
			key, _, err := r.DeclarationSelector(open.Declaration)
			if err != nil {
				return g.DeclarationDescription[s.Language]{}, err
			}
			name := key.Key()
			form := model.Form("alias")
			if name == "Person" || name == "Page" {
				form = model.Form("defined")
			}
			d := g.DeclarationDescription[s.Language]{Name: name, Form: form}
			for i := range open.Formals {
				d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: fmt.Sprintf("T%d", i)})
				d.Generic.Bindings = append(d.Generic.Bindings, g.Parameter[s.Language](i))
				d.Generic.Arguments = append(d.Generic.Arguments, g.Argument[s.Language](i))
			}
			if name == "Page" {
				d.Generic.Parameters = append(d.Generic.Parameters, g.ParameterDescription{Name: "Storage"})
				d.Generic.Arguments = append(d.Generic.Arguments, model.Collection(g.Argument[s.Language](0)))
			}
			d.Body, err = q.Mapped(open.Body)
			return d, err
		}}},
		Helpers: []g.HelperRecipe[s.Language]{{Key: cursor, Declaration: g.NativeDeclaration[s.Language]{Name: "Cursor", Form: model.Form("alias"), Body: model.Primitive("string")}}},
	}, nil
}
func (f *Fixture) Plan(ctx context.Context) (g.GenerationPlan, error) {
	var registry c.Registry
	from := f.source.snapshot.Language()
	if err := c.RegisterSource[SourceLanguage, *reader](&registry, from, sourceProvider{f}); err != nil {
		return g.GenerationPlan{}, err
	}
	if err := c.RegisterTarget(&registry, f.target, target.Provider{Package: "generated"}); err != nil {
		return g.GenerationPlan{}, err
	}
	if _, err := c.RegisterMapping(&registry, from, f.target, library{}); err != nil {
		return g.GenerationPlan{}, err
	}
	var roots []ast.NodeRef
	for i := range f.names {
		roots = append(roots, ast.NodeRef{Tree: f.source.snapshot.Entry().Tree, Node: ast.NodeID(i + 2)})
	}
	checker, err := check.NewChecker[s.Language](f.target, checkgo.Provider{})
	if err != nil {
		return g.GenerationPlan{}, err
	}
	return registry.PlanGeneration(ctx, c.Request{Input: transform.Input{Source: f.source.snapshot}, To: f.target, Roots: roots, Output: "independent", Profile: "independent/1", Limits: g.DefaultLimits()}, checker)
}
