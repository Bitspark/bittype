package generation

// PROTOTYPE (bittype#196): registration and selection under the contribution
// contract. A consumer's composition registers each contribution, each payload
// codec, each facet and each product it selects. Nothing is registered apart
// from a consumer's composition (O11).

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contributionEntry struct {
	from, to representation.Language
	api      reflect.Type
	profile  string
	identity g.ContributionIdentity
	domain   string
	value    any
	plan     func(context.Context, ContributionRequest, sourceEntry, targetEntry, []contributionEntry, []g.PayloadCodec, []g.Facet, []g.Product, check.Checker) (g.ContributionPlan, error)
}

type contributions struct {
	entries      []contributionEntry
	codecs       []g.PayloadCodec
	facets       []g.Facet
	products     []g.Product
	replacements [][2]g.ContributionIdentity
}

// ContributionRequest is a generation request under the contribution
// contract. Products names the products asked for besides native types (O5).
// Configuration is the target configuration the plan records (row 35).
type ContributionRequest struct {
	Request
	Products      []string
	Configuration []g.Setting
}

// RegisterContribution registers one contribution. The local key is unique
// for its publisher across targets (O4): one publisher with one key twice
// fails, and the same key under two publishers registers.
func RegisterContribution[F, T, R any](r *Registry, from, to representation.Language, c g.Contribution[F, T, R]) error {
	if r == nil || representation.CheckMarker[F](from) != nil || representation.CheckMarker[T](to) != nil || !c.Reference.Valid() || c.Profile == "" {
		return g.ErrInvalid
	}
	if c.Target != string(to.ID()) {
		return fmt.Errorf("%w: contribution %s states target %q and is registered for %q", g.ErrScope, c.Reference, c.Target, to.ID())
	}
	c.Facts, c.Owned = slices.Clone(c.Facts), slices.Clone(c.Owned)
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, e := range r.contributions.entries {
		if e.identity == c.Reference.Identity {
			return fmt.Errorf("%w: publisher %q registers the local key %q twice", g.ErrInvalid, e.identity.Publisher, e.identity.Key)
		}
	}
	r.contributions.entries = append(r.contributions.entries, contributionEntry{from, to, reflect.TypeFor[R](), c.Profile, c.Reference.Identity, c.Domain.Key(), c, planContributions[F, T, R]})
	return nil
}

// RegisterCodec registers the codec of one payload schema.
func (r *Registry) RegisterCodec(c g.PayloadCodec) error {
	if r == nil || representation.IsNil(c) || !c.Schema().Valid() {
		return g.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, old := range r.contributions.codecs {
		if old.Schema() == c.Schema() {
			return g.ErrAmbiguous
		}
	}
	r.contributions.codecs = append(r.contributions.codecs, c)
	return nil
}

// RegisterFacet registers one facet. The kernel admits a facet that names a
// registered contribution of the same publisher (O8).
func (r *Registry) RegisterFacet(f g.Facet) error {
	if r == nil || !f.Valid() {
		return g.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	known := slices.ContainsFunc(r.contributions.entries, func(e contributionEntry) bool { return e.identity == f.Contribution() })
	if !known || f.Publisher() != f.Contribution().Publisher {
		return fmt.Errorf("%w: a facet names no registered contribution of its publisher", g.ErrScope)
	}
	for _, old := range r.contributions.facets {
		if old.Contribution() == f.Contribution() && old.Kind().Identity == f.Kind().Identity && old.Target() == f.Target() {
			return g.ErrAmbiguous
		}
	}
	r.contributions.facets = append(r.contributions.facets, f)
	return nil
}

// RegisterProduct registers a product that a request may name.
func (r *Registry) RegisterProduct(p g.Product) error {
	if r == nil || p.Identity == "" || p.Identity == g.TypesProduct {
		return g.ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if slices.ContainsFunc(r.contributions.products, func(old g.Product) bool { return old.Identity == p.Identity }) {
		return g.ErrAmbiguous
	}
	r.contributions.products = append(r.contributions.products, p)
	return nil
}

// ReplaceContribution is the consumer's explicit composition of two
// contributions that claim one domain under one profile (row 15): expected is
// not selected, and with answers the domain.
func (r *Registry) ReplaceContribution(expected, with g.ContributionIdentity) error {
	if r == nil || expected == with {
		return g.ErrScope
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	var old, replacement *contributionEntry
	for i, e := range r.contributions.entries {
		if e.identity == expected {
			old = &r.contributions.entries[i]
		}
		if e.identity == with {
			replacement = &r.contributions.entries[i]
		}
	}
	if old == nil || replacement == nil || old.domain != replacement.domain || old.profile != replacement.profile || !old.to.Compatible(replacement.to) {
		return fmt.Errorf("%w: a replacement names two registered contributions of one domain, profile and target", g.ErrScope)
	}
	for _, pair := range r.contributions.replacements {
		if pair[0] == expected || pair[0] == with || pair[1] == expected {
			return g.ErrAmbiguous
		}
	}
	r.contributions.replacements = append(r.contributions.replacements, [2]g.ContributionIdentity{expected, with})
	return nil
}

// PlanContributions plans a request with every contribution that states the
// request's profile and target.
func (r *Registry) PlanContributions(ctx context.Context, q ContributionRequest, checker check.Checker) (g.ContributionPlan, error) {
	fail := func(err error) (g.ContributionPlan, error) { return g.ContributionPlan{}, err }
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
		if root.Node == 0 || int(root.Node) > trees[root.Tree] {
			return fail(g.ErrScope)
		}
	}
	r.mu.RLock()
	src, sok := r.sources[string(q.Input.Source.LanguageID())]
	target, tok := r.targets[string(q.To.ID())]
	replaced := map[g.ContributionIdentity]bool{}
	for _, pair := range r.contributions.replacements {
		replaced[pair[0]] = true
	}
	var entries []contributionEntry
	selected := map[g.ContributionIdentity]bool{}
	for _, e := range r.contributions.entries {
		if e.profile == q.Profile && e.from.Compatible(q.Input.Source.Language()) && e.to.Compatible(q.To) && !replaced[e.identity] {
			entries = append(entries, e)
			selected[e.identity] = true
		}
	}
	codecs := slices.Clone(r.contributions.codecs)
	var facets []g.Facet
	for _, f := range r.contributions.facets {
		if selected[f.Contribution()] {
			facets = append(facets, f)
		}
	}
	var products []g.Product
	missing := ""
	for _, name := range q.Products {
		i := slices.IndexFunc(r.contributions.products, func(p g.Product) bool { return p.Identity == name })
		if i < 0 {
			missing = name
			break
		}
		products = append(products, r.contributions.products[i])
	}
	r.mu.RUnlock()
	if missing != "" {
		return fail(fmt.Errorf("%w: no product %q is registered", g.ErrUnsupported, missing))
	}
	if !sok || !tok || !src.language.Compatible(q.Input.Source.Language()) || !target.language.Compatible(q.To) {
		return fail(g.ErrScope)
	}
	if len(entries) == 0 {
		return fail(&g.Refusal{Class: g.RefusalMissingContribution, Owner: g.KernelOwner, Reason: "no contribution states the profile " + q.Profile, Profile: q.Profile, Target: string(q.To.ID())})
	}
	for _, e := range entries {
		if e.api != src.api {
			return fail(g.ErrScope)
		}
	}
	return entries[0].plan(ctx, q, src, target, entries, codecs, facets, products, checker)
}

func planContributions[F, T, R any](ctx context.Context, q ContributionRequest, src sourceEntry, target targetEntry, entries []contributionEntry, codecs []g.PayloadCodec, facets []g.Facet, products []g.Product, checker check.Checker) (g.ContributionPlan, error) {
	provider, ok := src.provider.(g.SourceProvider[F, R])
	native, ok2 := target.provider.(g.Target[T])
	if !ok || !ok2 {
		return g.ContributionPlan{}, g.ErrScope
	}
	var selected []g.Contribution[F, T, R]
	needed := map[g.SchemaReference]bool{}
	for _, e := range entries {
		c, ok := e.value.(g.Contribution[F, T, R])
		if !ok {
			return g.ContributionPlan{}, g.ErrScope
		}
		selected = append(selected, c)
		needed[c.Schema] = true
	}
	// A plan is given the codecs of the schemas that its contributions state.
	codecs = slices.DeleteFunc(codecs, func(c g.PayloadCodec) bool { return !needed[c.Schema()] })
	return g.PlanContributions(ctx, g.ContributionInput[F, T, R]{From: src.language, Input: q.Input, Source: provider, Roots: slices.Clone(q.Roots), Profile: q.Profile, Contributions: selected, Codecs: codecs, Facets: facets, Products: products, Configuration: slices.Clone(q.Configuration), Execution: g.Execution[T]{Target: native, Checker: checker, Language: q.To, Output: q.Output, Limits: q.Limits}})
}
