package generationfixture

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type catalog struct{ f *Fixture }

func (catalog catalog) Declarations() []g.Declaration[SourceLanguage] {
	var out []g.Declaration[SourceLanguage]
	for _, p := range catalog.f.config.Plans {
		if p.Declaration.Valid() {
			out = append(out, p.Declaration)
		}
	}
	return out
}
func (v catalog) DeclarationSelector(d g.Declaration[SourceLanguage]) (g.Selector, bool, error) {
	for _, p := range v.f.config.Plans {
		if p.Declaration == d {
			key, err := g.NewSelector(v.f.source.snapshot.Language(), p.Key)
			return key, true, err
		}
	}
	return g.Selector{}, false, g.ErrScope
}
func (v catalog) UseSelector(use g.UseData[SourceLanguage]) (g.Selector, error) {
	name, err := v.f.source.selectPlan(use)
	if err != nil {
		return g.Selector{}, err
	}
	return g.NewSelector(v.f.source.snapshot.Language(), name)
}

type automaticSource struct {
	*Source
	catalog
}

func (s *automaticSource) ResolveRoots(ctx context.Context, sites []ast.NodeRef) ([]g.Declaration[SourceLanguage], error) {
	out := make([]g.Declaration[SourceLanguage], 0, len(sites))
	for _, site := range sites {
		if site.Tree != s.snapshot.Entry().Tree || site.Node < 2 {
			return nil, g.ErrScope
		}
		d, ok := s.declarations.Ref(uint64(site.Node - 1))
		if !ok {
			return nil, g.ErrScope
		}
		if _, exists := s.opens[d]; !exists {
			return nil, g.ErrScope
		}
		out = append(out, d)
	}
	return out, nil
}

type sourceProvider struct{ f *Fixture }

func (sourceProvider) Contract() transform.Contract { return transform.Contract{} }
func (p sourceProvider) Open(_ context.Context, r transform.Request[SourceLanguage]) (*automaticSource, error) {
	if !r.Snapshot().Same(p.f.source.snapshot) {
		return nil, g.ErrScope
	}
	return &automaticSource{p.f.source, catalog{p.f}}, nil
}

type targetProvider struct{ f *Fixture }

func (p targetProvider) Open(_ context.Context, r g.TargetRequest) (*Native, error) {
	info, err := r.Info()
	if err != nil {
		return nil, err
	}
	if info.Output != "generated-fixture" {
		return nil, g.ErrScope
	}
	return p.f.native, nil
}

type mappingProvider struct{ f *Fixture }

func (p mappingProvider) Open(_ context.Context, r g.MappingRequest[*automaticSource, *Native]) (g.Contribution[SourceLanguage, s.Language], error) {
	source, target, err := r.Opened()
	if err != nil {
		return g.Contribution[SourceLanguage, s.Language]{}, err
	}
	f := p.f
	if source.Source != f.source || target != f.native {
		return g.Contribution[SourceLanguage, s.Language]{}, g.ErrScope
	}
	var bindings []g.RepresentationBinding[SourceLanguage, s.Language]
	var helpers []g.Plan[SourceLanguage, s.Language]
	for _, plan := range f.config.Plans {
		if plan.Key == "cursor-helper" {
			helpers = append(helpers, plan)
			continue
		}
		key, err := g.NewSelector(f.source.snapshot.Language(), plan.Key)
		if err != nil {
			return g.Contribution[SourceLanguage, s.Language]{}, err
		}
		bindings = append(bindings, g.RepresentationBinding[SourceLanguage, s.Language]{Selector: key, Mapping: g.Mapping[SourceLanguage, s.Language]{ID: plan.Key, PluginID: "independent-fixture", Plan: func(d g.Declaration[SourceLanguage]) (g.Plan[SourceLanguage, s.Language], error) {
			if d != plan.Declaration {
				return g.Plan[SourceLanguage, s.Language]{}, g.ErrScope
			}
			return plan, nil
		}}})
	}
	return g.Contribution[SourceLanguage, s.Language]{Bindings: bindings, Helpers: helpers}, nil
}

// Automatic uses the same registry for an unrelated source language, independent
// association owners, an external generic and a native-only helper.
func (f *Fixture) Automatic(ctx context.Context) (g.Output, error) {
	var registry c.Registry
	from, to := f.source.snapshot.Language(), f.Target()
	source, err := c.RegisterSource[SourceLanguage, *automaticSource](&registry, from, sourceProvider{f})
	if err != nil {
		return g.Output{}, err
	}
	target, err := c.RegisterTarget[SourceLanguage, s.Language, *Native](&registry, from, to, targetProvider{f})
	if err != nil {
		return g.Output{}, err
	}
	mapping, err := c.RegisterMapping[SourceLanguage, s.Language, *automaticSource, *Native](&registry, from, to, mappingProvider{f})
	if err != nil {
		return g.Output{}, err
	}
	var roots []ast.NodeRef
	for _, d := range f.Roots() {
		index, _ := f.source.declarations.Index(d)
		roots = append(roots, ast.NodeRef{Tree: f.source.snapshot.Entry().Tree, Node: ast.NodeID(index + 1)})
	}
	return registry.Generate(ctx, c.Request{Input: transform.Input{Source: f.source.snapshot}, To: to, Output: "generated-fixture", Roots: roots, Selection: c.Selection{Source: source, Target: target, Mappings: []c.MappingRef{mapping}}, Limits: g.DefaultLimits()}, f.config.Checker)
}
