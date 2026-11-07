package generation

import (
	"context"
	"errors"
	"fmt"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type composition[F, T any] struct {
	Config   g.Config[F, T]
	Catalog  g.Catalog[F]
	Bindings bindingTable[F, T]
	Roots    []g.Declaration[F]
	Helpers  []g.Plan[F, T]
	Limits   g.Limits
}

// Run selects exact mappings once, publishes roots, drains acquired definitions
// to a fixed point, and commits the complete output. It never interprets syntax,
// invokes analysis or retries a failed mapping.
func run[F, T any](ctx context.Context, input g.Input, c composition[F, T]) (representation.ASTSet, error) {
	fail := func(err error) (representation.ASTSet, error) { return representation.ASTSet{}, err }
	if ctx == nil || ctx.Err() != nil {
		if ctx != nil {
			return fail(ctx.Err())
		}
		return fail(g.ErrInvalid)
	}
	if !input.Snapshot.Language().Compatible(c.Bindings.language) || representation.IsNil(c.Catalog) || !c.Limits.Valid() || c.Config.Select != nil || len(c.Config.Plans) != 0 {
		return fail(g.ErrInvalid)
	}
	plans := make([]g.Plan[F, T], 0)
	declarations := map[g.Declaration[F]]string{}
	direct := map[string]string{}
	appendPlan := func(m g.Mapping[F, T], d g.Declaration[F]) (string, error) {
		if len(plans) >= c.Limits.Plans {
			return "", g.ErrLimit
		}
		p, err := m.Plan(d)
		if err != nil {
			return "", err
		}
		if p.Declaration != d {
			return "", g.ErrScope
		}
		p.Key = fmt.Sprintf("mapping/%d", len(plans))
		plans = append(plans, p)
		return p.Key, nil
	}
	for _, binding := range c.Bindings.entries {
		// Exact declaration keys are instantiated below, not as fake direct plans.
		direct[binding.Selector.Key()] = ""
	}
	seen := map[g.Declaration[F]]bool{}
	declKeys := map[string]bool{}
	inventory := c.Catalog.Declarations()
	if len(inventory) > c.Limits.References {
		return fail(g.ErrLimit)
	}
	for _, d := range inventory {
		if err := ctx.Err(); err != nil {
			return fail(err)
		}
		if !d.Valid() || seen[d] {
			return fail(g.ErrInvalid)
		}
		seen[d] = true
		key, eligible, err := c.Catalog.DeclarationSelector(d)
		if err != nil {
			return fail(err)
		}
		if !key.Valid() || !input.Snapshot.Language().Compatible(key.Language()) || declKeys[key.Key()] {
			return fail(g.ErrInvalid)
		}
		declKeys[key.Key()] = true
		m, err := c.Bindings.selectMapping(key, eligible)
		if errors.Is(err, g.ErrUnsupported) {
			continue // Missing mappings fail if this declaration is demanded.
		}
		if err != nil {
			return fail(err)
		}
		plan, err := appendPlan(m, d)
		if err != nil {
			return fail(err)
		}
		declarations[d] = plan
	}
	for _, b := range c.Bindings.entries {
		if declKeys[b.Selector.Key()] {
			continue
		}
		plan, err := appendPlan(b.Mapping, g.Declaration[F]{})
		if err != nil {
			return fail(err)
		}
		direct[b.Selector.Key()] = plan
	}
	for _, helper := range c.Helpers {
		if len(plans) >= c.Limits.Plans {
			return fail(g.ErrLimit)
		}
		if helper.Declaration.Valid() || helper.Key == "" {
			return fail(g.ErrInvalid)
		}
		plans = append(plans, helper)
	}
	c.Config.Plans = plans
	c.Config.Select = func(use g.UseData[F]) (string, error) {
		if use.Callee.Valid() {
			if key := declarations[use.Callee]; key != "" {
				return key, nil
			}
			return "", g.ErrUnsupported
		}
		key, err := c.Catalog.UseSelector(use)
		if err != nil {
			return "", err
		}
		if !key.Valid() || !input.Snapshot.Language().Compatible(key.Language()) {
			return "", g.ErrInvalid
		}
		if value := direct[key.Key()]; value != "" {
			return value, nil
		}
		return "", fmt.Errorf("%w: no representation for %s", g.ErrUnsupported, key.Key())
	}
	run, err := g.Begin(ctx, input, c.Config, c.Limits)
	if err != nil {
		return fail(err)
	}
	for _, root := range c.Roots {
		if _, err = run.Publish(root); err != nil {
			return fail(err)
		}
	}
	for {
		pending, err := run.PendingDefinitions()
		if err != nil {
			return fail(err)
		}
		if len(pending) == 0 {
			break
		}
		for _, callee := range pending {
			if _, err = run.Define(callee); err != nil {
				return fail(err)
			}
		}
	}
	return run.Finish()
}
