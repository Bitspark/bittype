package generation

import (
	"context"
	"github.com/Bitspark/bittype/contracts/check/go"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func adapter[F, T any, S g.SourceSession[F], N g.Native[F, T]](request Request, src sourceEntry, target targetEntry, mappings []mappingEntry, checker check.Checker) (runner, error) {
	sp, ok := src.provider.(g.SourceProvider[F, S])
	if !ok {
		return nil, g.ErrScope
	}
	tp, ok := target.provider.(g.TargetProvider[F, T, N])
	if !ok {
		return nil, g.ErrScope
	}
	providers := make([]g.MappingProvider[F, T, S, N], len(mappings))
	for i, m := range mappings {
		providers[i], ok = m.provider.(g.MappingProvider[F, T, S, N])
		if !ok {
			return nil, g.ErrScope
		}
	}
	return func(ctx context.Context, input transform.Input) (g.Output, error) {
		zero := g.Output{}
		r, err := transform.AdmitRequest[F](src.language, src.contract, input)
		if err != nil {
			return zero, err
		}
		if !r.Snapshot().Same(request.Input.Source) {
			return zero, g.ErrScope
		}
		for i, expected := range request.Input.Evidence {
			actual, exists := r.Evidence(i)
			if !exists || actual != expected {
				return zero, g.ErrScope
			}
		}
		if _, extra := r.Evidence(len(request.Input.Evidence)); extra {
			return zero, g.ErrScope
		}
		construction, err := g.NewConstruction(ctx, g.TargetInfo{From: src.language, To: request.To, Output: request.Output, Limits: request.Limits})
		if err != nil {
			return zero, err
		}
		defer construction.Close()
		sh, err := g.OpenSource(construction, sp, r)
		if err != nil {
			return zero, err
		}
		source, err := sh.Value(construction)
		if err != nil {
			return zero, err
		}
		roots, err := source.ResolveRoots(ctx, slices.Clone(request.Roots))
		if err != nil {
			return zero, err
		}
		if len(roots) != len(request.Roots) {
			return zero, g.ErrScope
		}
		roots = slices.Clone(roots)
		seen := map[g.Declaration[F]]bool{}
		for _, root := range roots {
			if !root.Valid() || seen[root] {
				return zero, g.ErrScope
			}
			seen[root] = true
		}
		if err = source.Authenticate(ctx, g.References[F]{Declarations: roots}); err != nil {
			return zero, err
		}
		nh, err := g.OpenTarget(construction, tp)
		if err != nil {
			return zero, err
		}
		native, err := nh.Value(construction)
		if err != nil {
			return zero, err
		}
		contributions := make([]g.Contribution[F, T], len(providers))
		total := 0
		for i, p := range providers {
			part, err := g.OpenMapping(construction, p, sh, nh)
			if err != nil {
				return zero, err
			}
			total += len(part.Bindings) + len(part.Helpers)
			if part.Structural != nil {
				total++
			}
			if total > request.Limits.Plans {
				return zero, g.ErrLimit
			}
			part.Bindings = slices.Clone(part.Bindings)
			part.Helpers = slices.Clone(part.Helpers)
			if part.Structural != nil {
				copy := *part.Structural
				part.Structural = &copy
			}
			contributions[i] = part
		}
		bindings, helpers, err := collect[F, T](src.language, request.Selection, contributions)
		if err != nil {
			return zero, err
		}
		var evidence []semantic.Result
		for i := 0; ; i++ {
			e, exists := r.Evidence(i)
			if !exists {
				break
			}
			evidence = append(evidence, e)
		}
		output, err := run(ctx, g.Input{Snapshot: r.Snapshot(), Evidence: evidence}, composition[F, T]{Config: g.Config[F, T]{Source: source, Native: native, Checker: checker, Target: request.To, Evidence: src.contract.Evidence}, Catalog: source, Bindings: bindings, Roots: roots, Helpers: helpers, Limits: request.Limits})
		if err != nil {
			return zero, err
		}
		if output.Snapshot().Entry().Tree != request.Output {
			return zero, g.ErrScope
		}
		inputTrees := map[ast.TreeID]bool{}
		for _, tree := range r.Inputs() {
			info, err := tree.Inspect()
			if err != nil {
				return zero, err
			}
			inputTrees[info.ID] = true
		}
		for _, tree := range output.Snapshot().Inputs() {
			info, err := tree.Inspect()
			if err != nil {
				return zero, err
			}
			if inputTrees[info.ID] {
				return zero, g.ErrScope
			}
		}
		return output, nil
	}, nil
}

func collect[F, T any](from representation.Language, selection Selection, parts []g.Contribution[F, T]) (bindingTable[F, T], []g.Plan[F, T], error) {
	type candidate struct {
		ref     BindingRef
		binding g.RepresentationBinding[F, T]
	}
	fail := func(err error) (bindingTable[F, T], []g.Plan[F, T], error) { return bindingTable[F, T]{}, nil, err }
	groups := map[string][]candidate{}
	var order []string
	var structural *g.Mapping[F, T]
	var helpers []g.Plan[F, T]
	helperKeys := map[string]bool{}
	for i, part := range parts {
		ref := selection.Mappings[i]
		ids := map[string]bool{}
		qualify := func(m g.Mapping[F, T]) (g.Mapping[F, T], error) {
			if m.ID == "" || m.PluginID == "" || m.Plan == nil || ids[m.ID] {
				return m, g.ErrInvalid
			}
			ids[m.ID] = true

			return m, nil
		}
		for _, b := range part.Bindings {
			if !b.Selector.Valid() || !from.Compatible(b.Selector.Language()) {
				return fail(g.ErrScope)
			}
			id := b.Mapping.ID
			m, err := qualify(b.Mapping)
			if err != nil {
				return fail(err)
			}
			b.Mapping = m
			if len(groups[b.Selector.Key()]) == 0 {
				order = append(order, b.Selector.Key())
			}
			groups[b.Selector.Key()] = append(groups[b.Selector.Key()], candidate{BindingRef{ref, id}, b})
		}
		if part.Structural != nil {
			if structural != nil {
				return fail(g.ErrInvalid)
			}
			m, err := qualify(*part.Structural)
			if err != nil {
				return fail(err)
			}
			structural = &m
		}
		for _, p := range part.Helpers {
			if p.Key == "" || helperKeys[p.Key] || p.Declaration.Valid() {
				return fail(g.ErrInvalid)
			}
			helperKeys[p.Key] = true
			helpers = append(helpers, p)
		}
	}
	replacements := map[string]Replacement{}
	for _, r := range selection.Replace {
		if _, exists := replacements[r.Selector.Key()]; exists {
			return fail(g.ErrInvalid)
		}
		replacements[r.Selector.Key()] = r
	}
	var entries []g.RepresentationBinding[F, T]
	for _, key := range order {
		group := groups[key]
		replacement, replace := replacements[key]
		if !replace {
			if len(group) != 1 {
				return fail(g.ErrInvalid)
			}
			entries = append(entries, group[0].binding)
			continue
		}
		if len(group) != 2 {
			return fail(g.ErrInvalid)
		}
		old, new := -1, -1
		for i, c := range group {
			if c.ref == replacement.Expected {
				old = i
			}
			if c.ref == replacement.With {
				new = i
			}
		}
		if old < 0 || new < 0 || old == new {
			return fail(g.ErrInvalid)
		}
		entries = append(entries, group[new].binding)
		delete(replacements, key)
	}
	if len(replacements) != 0 {
		return fail(g.ErrInvalid)
	}
	return bindingTable[F, T]{from, entries, structural}, helpers, nil
}
