// Package generation is the registry of the generation kernel: a consumer's
// composition registers a source, a target and each contribution, payload
// codec, facet and product it selects, and the registry plans a request with
// what states the request's profile.
package generation

import (
	"context"
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

// Request is a generation request: the roots, the target, the profile and the
// products. Products names the products asked for besides target types (O5):
// with none the request asks for types alone. Configuration is the target
// configuration that the plan records (row 35).
type Request struct {
	Input         transform.Input
	To            representation.Language
	Roots         []ast.NodeRef
	Output        ast.TreeID
	Profile       string
	Limits        g.Limits
	Products      []string
	Configuration []g.Setting
}
type sourceEntry struct {
	language representation.Language
	api      reflect.Type
	provider any
}
type targetEntry struct {
	language representation.Language
	provider any
}
type Registry struct {
	mu      sync.RWMutex
	sources map[string]sourceEntry
	targets map[string]targetEntry
	// contributions is the registration state of the contribution contract
	// (bittype#196).
	contributions contributions
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
	// This rejects an unusable contract early. Each Plan admits the evidence
	// against the provider's Contract at that time and checks it again.
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.sources == nil {
		r.sources = map[string]sourceEntry{}
	}
	if _, ok := r.sources[string(lang.ID())]; ok {
		return g.ErrAmbiguous
	}
	r.sources[string(lang.ID())] = sourceEntry{lang, reflect.TypeFor[R](), p}
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
	selected, selectionErr := r.contributions.selectFor(q)
	r.mu.RUnlock()
	if !sok || !tok || !src.language.Compatible(q.Input.Source.Language()) || !target.language.Compatible(q.To) {
		return fail(g.ErrScope)
	}
	if selectionErr != nil {
		return fail(diagnostic(q.Input, selectionErr))
	}
	if len(selected.entries) == 0 {
		// No contribution states the request's profile for this source and
		// target: there is nothing to select (O17).
		return fail(&g.Refusal{Class: g.RefusalMissingContribution, Owner: g.KernelOwner, Reason: g.ReasonNoProfile + " " + q.Profile, Profile: q.Profile, Target: string(q.To.ID())})
	}
	p, err := selected.plan(ctx, q, src, target, checker)
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
