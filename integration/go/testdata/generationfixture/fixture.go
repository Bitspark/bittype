package generationfixture

import (
	"context"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Fixture struct {
	source       *Source
	native       *Native
	config       g.Config[SourceLanguage, s.Language]
	declarations map[string]g.Declaration[SourceLanguage]
	rootUses     map[string]g.Use[SourceLanguage]
	roots        []g.Declaration[SourceLanguage]
}

func (f *Fixture) Input() g.Input { return g.Input{Snapshot: f.source.snapshot} }
func (f *Fixture) Config() g.Config[SourceLanguage, s.Language] {
	c := f.config
	c.Plans = slices.Clone(c.Plans)
	return c
}
func (f *Fixture) Target() representation.Language                       { return f.native.language }
func (f *Fixture) Declaration(name string) g.Declaration[SourceLanguage] { return f.declarations[name] }
func (f *Fixture) Use(name string) g.Use[SourceLanguage]                 { return f.rootUses[name] }
func (f *Fixture) Roots() []g.Declaration[SourceLanguage]                { return slices.Clone(f.roots) }

// New composes independent source, target and mapping owners through public
// contracts. No engine registry or hardcoded product variant is involved.
func New(ctx context.Context) (*Fixture, error) {
	return newFixture(ctx, false)
}

// NewArgumentOnlySupport composes a deliberately incorrect Page candidate with
// the unchanged relation authority. Qualification requires its rejection.
func NewArgumentOnlySupport(ctx context.Context) (*Fixture, error) {
	return newFixture(ctx, true)
}

func newFixture(ctx context.Context, argumentOnlySupport bool) (*Fixture, error) {
	source, err := newSource(ctx)
	if err != nil {
		return nil, err
	}
	native, err := newNative()
	if err != nil {
		return nil, err
	}
	f := &Fixture{source: source, native: native, declarations: make(map[string]g.Declaration[SourceLanguage]), rootUses: make(map[string]g.Use[SourceLanguage])}
	names := []string{"Person", "Page", "Outer", "Nested", "People", "Inline", "External", "InlinePeople", "ExternalPeople"}
	for _, name := range names {
		generic := name == "Page" || name == "Outer" || name == "Nested" || name == "Inline" || name == "External"
		f.declarations[name] = source.declaration(generic)
	}
	person, page := f.declarations["Person"], f.declarations["Page"]
	personOpen, pageOpen := source.opens[person], source.opens[page]
	source.body(person, source.use(personOpen.Environment, "person-body", g.Declaration[SourceLanguage]{}, g.Binder[SourceLanguage]{}))
	source.body(page, source.use(pageOpen.Environment, "page-body", g.Declaration[SourceLanguage]{}, g.Binder[SourceLanguage]{}))
	outer := source.opens[f.declarations["Outer"]]
	source.body(outer.Declaration, source.call(outer.Environment, "Page", page, source.occurrence(outer.Formals[0].Subject)))
	nested := source.opens[f.declarations["Nested"]]
	inner := source.call(nested.Environment, "Page", page, source.occurrence(nested.Formals[0].Subject))
	source.body(nested.Declaration, source.call(nested.Environment, "Page", page, inner))
	for _, name := range []string{"Inline", "External"} {
		o := source.opens[f.declarations[name]]
		source.body(o.Declaration, o.Formals[0].Subject)
	}
	for _, pair := range [][2]string{{"People", "Page"}, {"InlinePeople", "Inline"}, {"ExternalPeople", "External"}} {
		o := source.opens[f.declarations[pair[0]]]
		personUse := source.call(o.Environment, "Person", person, g.Use[SourceLanguage]{})
		source.body(o.Declaration, source.call(o.Environment, pair[1], f.declarations[pair[1]], personUse))
	}
	rootPerson := source.call(source.root, "Person", person, g.Use[SourceLanguage]{})
	f.rootUses["Person"] = rootPerson
	for _, name := range []string{"Page", "Outer", "Nested", "Inline", "External"} {
		f.rootUses[name] = source.call(source.root, name, f.declarations[name], rootPerson)
	}
	innerRoot := source.call(source.root, "Page", page, rootPerson)
	f.rootUses["NestedPage"] = source.call(source.root, "Page", page, innerRoot)

	// Each entry has its own association owner. The inline library is an
	// ordinary additional configuration entry and contributes no core variant.
	rules := make(map[string]*Rules)
	var readers []*Rules
	for _, name := range append(slices.Clone(names), "person-body", "page-body") {
		rules[name] = newRules(native)
		readers = append(readers, rules[name])
	}
	read := supportReader{owners: readers, source: source}
	var plans []g.Plan[SourceLanguage, s.Language]
	for _, name := range names {
		decl := f.declarations[name]
		open := source.opens[decl]
		owned := name != "Inline" && name != "External"
		alias := name != "Person" && name != "Page"
		exportName, artifact, arity := name, "", 0
		if !owned {
			artifact, arity = "fixture-native-library-v1", 2
			if name == "External" {
				exportName = "ExternalBox"
			}
		}
		plan := g.Plan[SourceLanguage, s.Language]{Key: name, Declaration: decl, Associations: rules[name]}
		plan.Acquire = func(_ context.Context, site g.Site[SourceLanguage, s.Language]) (g.Destination[s.Language], bool, error) {
			return native.reserve(site, exportName, owned, alias, artifact, arity)
		}
		if len(open.Formals) != 0 {
			rules[name].applicationDeclaration = decl
			rules[name].applicationInline = name == "Inline"
			parameter := "T"
			if name == "Outer" {
				parameter = "U"
			}
			if name == "Nested" {
				parameter = "V"
			}
			mapping := &genericMapping{native: native, rules: rules[name], read: read, name: parameter, inline: name == "Inline", argumentOnlySupport: argumentOnlySupport && name == "Page"}
			plan.User, plan.Convention = mapping, mapping
		} else {
			plan.User = nonGeneric(native, rules[name])
			plan.Header = func(_ context.Context, site g.Site[SourceLanguage, s.Language], destination g.Destination[s.Language]) (g.Layout[s.Language], error) {
				return native.freeze(site, destination)
			}
		}
		if owned {
			plan.Publisher, plan.Definer = publisher{}, defineFunc(defineSource)
			f.roots = append(f.roots, decl)
		}
		plans = append(plans, plan)
	}
	plans = append(plans, g.Plan[SourceLanguage, s.Language]{Key: "bound", User: useFunc(func(_ context.Context, r g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
		return r.Bound()
	})})
	plans = append(plans, g.Plan[SourceLanguage, s.Language]{Key: "person-body", Associations: rules["person-body"], User: useFunc(func(ctx context.Context, r g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
		text, err := native.text()
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		body, err := native.record(nativeField{"Name", text})
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		token, err := rules["person-body"].issue(ctx, r.Want(), body, text)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		return r.Return(token)
	})})
	plans = append(plans, g.Plan[SourceLanguage, s.Language]{Key: "page-body", Associations: rules["page-body"], User: useFunc(func(ctx context.Context, r g.UseRequest[SourceLanguage, s.Language]) (g.Product[SourceLanguage, s.Language], error) {
		parameter, err := r.Child(pageOpen.Formals[0].Subject)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		support, err := read.readActual(ctx, r.Want(), pageOpen.Formals[0].Subject, parameter)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		items, err := native.slice(parameter.Payload())
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		helper, err := r.RequireHelper("cursor-helper")
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		cursor, err := native.reference(helper)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		body, err := native.record(nativeField{"Items", items}, nativeField{"Support", support}, nativeField{"Cursor", cursor})
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		token, err := rules["page-body"].issue(ctx, r.Want(), body, items)
		if err != nil {
			return g.Product[SourceLanguage, s.Language]{}, err
		}
		return r.Return(token)
	})})
	plans = append(plans, g.Plan[SourceLanguage, s.Language]{
		Key: "cursor-helper",
		Acquire: func(_ context.Context, site g.Site[SourceLanguage, s.Language]) (g.Destination[s.Language], bool, error) {
			return native.reserve(site, "Cursor", true, false, "", 0)
		},
		Header: func(_ context.Context, site g.Site[SourceLanguage, s.Language], destination g.Destination[s.Language]) (g.Layout[s.Language], error) {
			return native.freeze(site, destination)
		},
		Definer: defineFunc(func(_ context.Context, r g.DefineRequest[SourceLanguage, s.Language]) (g.Definition[SourceLanguage, s.Language], error) {
			text, err := native.text()
			if err != nil {
				return g.Definition[SourceLanguage, s.Language]{}, err
			}
			return r.CompleteBody(text)
		}),
	})
	f.config = g.Config[SourceLanguage, s.Language]{Source: source, Native: native, Target: native.language, Plans: plans, Select: source.selectPlan}
	return f, nil
}

// ExternalSource is the supplied library export qualified alongside the emitted
// source. The generator emits an unresolved native spelling for this export,
// never a made-up local NodeRef or a second local declaration.
const ExternalSource = "package fixture\ntype ExternalBox[T any, S any] struct { Value T; Support S }\n"
