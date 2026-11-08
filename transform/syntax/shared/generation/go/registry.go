// Package generation composes registered mapping descriptions into semantic plans.
package generation

import (
	"context"
	"fmt"
	"reflect"
	"slices"
	"sync"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type MappingRef struct {
	owner *Registry
	id    uint64
}
type BindingRef struct {
	Provider MappingRef
	ID       string
}
type Replacement struct {
	Selector       g.Selector
	Expected, With BindingRef
}
type Request struct {
	Input   transform.Input
	To      representation.Language
	Roots   []ast.NodeRef
	Output  ast.TreeID
	Profile string
	Limits  g.Limits
}
type sourceEntry struct {
	language representation.Language
	api      reflect.Type
	contract transform.Contract
	provider any
}
type targetEntry struct {
	language representation.Language
	provider any
}
type mappingEntry struct {
	from, to representation.Language
	api      reflect.Type
	profile  string
	value    any
	plan     func(context.Context, Request, sourceEntry, targetEntry, []mappingEntry, []Replacement, check.Checker) (g.GenerationPlan, error)
	ref      MappingRef
}
type Registry struct {
	mu           sync.RWMutex
	next         uint64
	sources      map[string]sourceEntry
	targets      map[string]targetEntry
	mappings     []mappingEntry
	replacements []Replacement
}

func RegisterSource[F, R any](r *Registry, lang representation.Language, p g.SourceProvider[F, R]) error {
	if r == nil || representation.IsNil(p) || representation.CheckMarker[F](lang) != nil {
		return g.ErrInvalid
	}
	c := p.Contract()
	if err := semantic.ValidateEvidenceRequirements(c.Evidence); err != nil {
		return err
	}
	graphs := 0
	for _, requirement := range c.Evidence {
		if requirement.Semantics {
			graphs++
		}
	}
	if graphs != 1 {
		// The kernel walks exactly the graph one admitted result publishes.
		return g.ErrInvalid
	}
	c.Evidence = slices.Clone(c.Evidence)
	for i := range c.Evidence {
		c.Evidence[i] = c.Evidence[i].Clone()
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.sources == nil {
		r.sources = map[string]sourceEntry{}
	}
	if _, ok := r.sources[string(lang.ID())]; ok {
		return g.ErrAmbiguous
	}
	r.sources[string(lang.ID())] = sourceEntry{lang, reflect.TypeFor[R](), c, p}
	return nil
}
func RegisterTarget[T any](r *Registry, lang representation.Language, p g.Target[T]) error {
	if r == nil || representation.IsNil(p) || representation.CheckMarker[T](lang) != nil {
		return g.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.targets == nil {
		r.targets = map[string]targetEntry{}
	}
	if _, ok := r.targets[string(lang.ID())]; ok {
		return g.ErrAmbiguous
	}
	r.targets[string(lang.ID())] = targetEntry{lang, p}
	return nil
}
func RegisterMapping[F, T, R any](r *Registry, from, to representation.Language, p g.MappingLibrary[F, T, R]) (MappingRef, error) {
	if r == nil || representation.IsNil(p) || representation.CheckMarker[F](from) != nil || representation.CheckMarker[T](to) != nil {
		return MappingRef{}, g.ErrInvalid
	}
	// Provider code runs outside the registry lock and only at registration.
	desc, err := p.Describe(from, to)
	if err != nil {
		return MappingRef{}, err
	}
	desc, err = g.FreezeDescription(from, desc, g.DefaultLimits())
	if err != nil {
		return MappingRef{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.next++
	ref := MappingRef{r, r.next}
	r.mappings = append(r.mappings, mappingEntry{from, to, reflect.TypeFor[R](), desc.Profile, desc, plan[F, T, R], ref})
	return ref, nil
}
func (r *Registry) Replace(v Replacement) error {
	if r == nil || !v.Selector.Valid() || v.Expected.Provider.owner != r || v.With.Provider.owner != r || v.Expected.ID == "" || v.With.ID == "" || v.Expected == v.With {
		return g.ErrScope
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	// Chained/overlapping policy is ambiguous rather than registration-order based.
	for _, old := range r.replacements {
		if old.Selector.Key() == v.Selector.Key() && old.Selector.Language().Compatible(v.Selector.Language()) {
			return g.ErrAmbiguous
		}
	}
	r.replacements = append(r.replacements, v)
	return nil
}
func (r *Registry) PlanGeneration(ctx context.Context, q Request, checker check.Checker) (g.GenerationPlan, error) {
	fail := func(err error) (g.GenerationPlan, error) { return g.GenerationPlan{}, err }
	if r == nil || ctx == nil || !q.Input.Source.Valid() || q.Profile == "" || q.Output == "" || len(q.Roots) == 0 || !q.Limits.Valid() {
		return fail(g.ErrInvalid)
	}
	if err := ctx.Err(); err != nil {
		return fail(err)
	}
	if checker.Validate() != nil || !checker.Language().Compatible(q.To) {
		return fail(g.ErrScope)
	}
	if len(q.Roots) > q.Limits.References || len(q.Input.Evidence) > semantic.DefaultLimits().MaxInputs {
		return fail(g.ErrLimit)
	}
	q.Input.Evidence = slices.Clone(q.Input.Evidence)
	q.Roots = slices.Clone(q.Roots)
	trees := map[ast.TreeID]int{}
	seen := map[ast.NodeRef]bool{}
	for _, t := range q.Input.Source.Inputs() {
		info, err := t.Inspect()
		if err != nil {
			return fail(err)
		}
		if info.ID == q.Output {
			return fail(g.ErrScope)
		}
		trees[info.ID] = info.Nodes
	}
	for _, root := range q.Roots {
		if root.Node == 0 || int(root.Node) > trees[root.Tree] || seen[root] {
			return fail(g.ErrScope)
		}
		seen[root] = true
	}
	r.mu.RLock()
	src, sok := r.sources[string(q.Input.Source.LanguageID())]
	target, tok := r.targets[string(q.To.ID())]
	mappings := []mappingEntry{}
	for _, m := range r.mappings {
		if m.profile == q.Profile && m.from.Compatible(q.Input.Source.Language()) && m.to.Compatible(q.To) {
			mappings = append(mappings, m)
		}
	}
	replacements := slices.Clone(r.replacements)
	r.mu.RUnlock()
	if !sok || !tok || !src.language.Compatible(q.Input.Source.Language()) || !target.language.Compatible(q.To) {
		return fail(g.ErrScope)
	}
	if len(mappings) == 0 {
		return fail(g.ErrUnsupported)
	}
	if len(mappings) > q.Limits.Plans {
		return fail(g.ErrLimit)
	}
	for _, m := range mappings {
		if m.api != src.api {
			return fail(g.ErrScope)
		}
	}
	p, err := mappings[0].plan(ctx, q, src, target, mappings, replacements, checker)
	if err != nil {
		return fail(diagnostic(q.Input, err))
	}
	return p, nil
}
func (r *Registry) Generate(ctx context.Context, q Request, checker check.Checker) (g.Output, error) {
	p, err := r.PlanGeneration(ctx, q, checker)
	if err != nil {
		return g.Output{}, err
	}
	out, err := p.Run(ctx)
	if err != nil {
		return g.Output{}, diagnostic(q.Input, err)
	}
	return out, nil
}
func plan[F, T, R any](ctx context.Context, q Request, src sourceEntry, target targetEntry, entries []mappingEntry, replacements []Replacement, checker check.Checker) (g.GenerationPlan, error) {
	fail := func(err error) (g.GenerationPlan, error) { return g.GenerationPlan{}, err }
	provider, ok := src.provider.(g.SourceProvider[F, R])
	if !ok {
		return fail(g.ErrScope)
	}
	native, ok := target.provider.(g.Target[T])
	if !ok {
		return fail(g.ErrScope)
	}
	captured, err := g.AdmitInput[F](src.language, src.contract, q.Input)
	if err != nil {
		return fail(err)
	}
	var parts []g.SelectedDescription[F, T, R]
	for _, e := range entries {
		d, ok := e.value.(g.MappingDescription[F, T, R])
		if !ok {
			return fail(g.ErrScope)
		}
		d, err = g.FreezeDescription(src.language, d, q.Limits)
		if err != nil {
			return fail(err)
		}
		parts = append(parts, g.SelectedDescription[F, T, R]{ID: fmt.Sprintf("library-%d", e.ref.id), Description: d})
	}
	// Explicit replacements refer to both original and replacement offers.
	for _, v := range replacements {
		if !v.Selector.Language().Compatible(src.language) {
			continue
		}
		expected, with := -1, -1
		oldIndex, newIndex := -1, -1
		for i, e := range entries {
			for j, b := range parts[i].Description.Bindings {
				if b.Input.Selector.Key() != v.Selector.Key() {
					continue
				}
				if e.ref == v.Expected.Provider && b.ID == v.Expected.ID {
					expected, oldIndex = i, j
				}
				if e.ref == v.With.Provider && b.ID == v.With.ID {
					with, newIndex = i, j
				}
			}
		}
		if expected < 0 && with < 0 {
			continue
		}
		if expected < 0 || with < 0 || oldIndex < 0 || newIndex < 0 {
			return fail(g.ErrScope)
		}
		parts[expected].Description.Bindings = append(parts[expected].Description.Bindings[:oldIndex:oldIndex], parts[expected].Description.Bindings[oldIndex+1:]...)
	}
	source, err := provider.Open(ctx, captured)
	if err != nil {
		return fail(err)
	}
	if representation.IsNil(source.Interpretation) {
		return fail(g.ErrInvalid)
	}
	return g.Plan(ctx, g.PlanningInput[F, T, R]{Captured: captured, Source: source, Roots: slices.Clone(q.Roots), Mappings: parts, Execution: g.Execution[T]{Target: native, Checker: checker, Language: q.To, Output: q.Output, Limits: q.Limits}})
}
