package generation

// Registration and selection under the contribution contract (bittype#196). A consumer's composition registers each contribution, each payload
// codec, each facet and each product it selects. Nothing is registered apart
// from a consumer's composition (O11).

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type contributionEntry struct {
	from, to representation.Language
	api      reflect.Type
	profile  string
	identity g.ContributionIdentity
	domain   string
	value    any
	plan     func(context.Context, Request, sourceEntry, targetEntry, []contributionEntry, []g.PayloadCodec, []g.Facet, []g.Product, check.Checker) (g.GenerationPlan, error)
}

type contributions struct {
	entries      []contributionEntry
	codecs       []g.PayloadCodec
	facets       []g.Facet
	products     []g.Product
	replacements [][2]g.ContributionIdentity
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
	if err := c.Validate(); err != nil {
		return err
	}
	// The registry keeps its own copy: nothing the caller holds reaches a plan
	// after registration.
	c.Facts, c.Owned = slices.Clone(c.Facts), slices.Clone(c.Owned)
	for i := range c.Owned {
		c.Owned[i].Declaration = c.Owned[i].Declaration.Clone()
	}
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

// selection is what the registry selected for one request: every contribution
// that states the request's profile and target and that no replacement set
// aside, the registered codecs, the facets registered for those
// contributions and the products the request names.
type selection struct {
	entries  []contributionEntry
	codecs   []g.PayloadCodec
	facets   []g.Facet
	products []g.Product
}

// selectFor selects for a request. The caller holds the registry's lock.
func (c *contributions) selectFor(q Request) (selection, error) {
	replaced := map[g.ContributionIdentity]bool{}
	for _, pair := range c.replacements {
		replaced[pair[0]] = true
	}
	var out selection
	selected := map[g.ContributionIdentity]bool{}
	for _, e := range c.entries {
		if e.profile == q.Profile && e.from.Compatible(q.Input.Source.Language()) && e.to.Compatible(q.To) && !replaced[e.identity] {
			out.entries = append(out.entries, e)
			selected[e.identity] = true
		}
	}
	out.codecs = slices.Clone(c.codecs)
	for _, f := range c.facets {
		if selected[f.Contribution()] {
			out.facets = append(out.facets, f)
		}
	}
	for _, name := range q.Products {
		i := slices.IndexFunc(c.products, func(p g.Product) bool { return p.Identity == name })
		if i < 0 {
			return selection{}, fmt.Errorf("%w: no product %q is registered", g.ErrUnsupported, name)
		}
		out.products = append(out.products, c.products[i])
	}
	return out, nil
}

// plan plans a request with what was selected for it.
func (s selection) plan(ctx context.Context, q Request, src sourceEntry, target targetEntry, checker check.Checker) (g.GenerationPlan, error) {
	if len(s.entries) > q.Limits.Plans {
		return g.GenerationPlan{}, g.ErrLimit
	}
	for _, e := range s.entries {
		if e.api != src.api {
			return g.GenerationPlan{}, g.ErrScope
		}
	}
	return s.entries[0].plan(ctx, q, src, target, s.entries, s.codecs, s.facets, s.products, checker)
}

func planContributions[F, T, R any](ctx context.Context, q Request, src sourceEntry, target targetEntry, entries []contributionEntry, codecs []g.PayloadCodec, facets []g.Facet, products []g.Product, checker check.Checker) (g.GenerationPlan, error) {
	provider, ok := src.provider.(g.SourceProvider[F, R])
	built, ok2 := target.provider.(g.Target[T])
	if !ok || !ok2 {
		return g.GenerationPlan{}, g.ErrScope
	}
	var selected []g.Contribution[F, T, R]
	needed := map[g.SchemaReference]bool{}
	for _, e := range entries {
		c, ok := e.value.(g.Contribution[F, T, R])
		if !ok {
			return g.GenerationPlan{}, g.ErrScope
		}
		selected = append(selected, c)
		needed[c.Schema] = true
	}
	// A plan is given the codecs of the schemas that its contributions state.
	codecs = slices.DeleteFunc(codecs, func(c g.PayloadCodec) bool { return !needed[c.Schema()] })
	return g.PlanContributions(ctx, g.ContributionInput[F, T, R]{From: src.language, Input: q.Input, Source: provider, Roots: slices.Clone(q.Roots), Profile: q.Profile, Contributions: selected, Codecs: codecs, Facets: facets, Products: products, Configuration: slices.Clone(q.Configuration), Execution: g.Execution[T]{Target: built, Checker: checker, Language: q.To, Output: q.Output, Limits: q.Limits}})
}
