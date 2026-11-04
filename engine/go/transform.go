package engine

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
)

func (e *Engine) Check(input representation.AST) error {
	registration, err := e.lookup(input.LanguageID())
	if err != nil {
		return err
	}
	if registration.checker == nil {
		return fmt.Errorf("%w: checking for %s", ErrMissingCapability, input.LanguageID())
	}
	return registration.checker.Check(input)
}

func (e *Engine) TransformTo(ctx context.Context, input transform.Input, targetID language.ID) (representation.AST, error) {
	plan, err := e.Plan(input.Source.LanguageID(), targetID)
	if err != nil {
		return representation.AST{}, err
	}
	return e.Execute(ctx, input, plan)
}

// Execute uses one route and never retries another mapping after an error.
// Providers receive owned snapshots. Failure returns a zero envelope, even if a
// plugin returned a partial result alongside its error.
func (e *Engine) Execute(ctx context.Context, input transform.Input, plan Plan) (representation.AST, error) {
	if ctx == nil { return representation.AST{}, semantic.ErrInvalidRequest }
	if err := ctx.Err(); err != nil { return representation.AST{}, err }
	e.mu.RLock()
	registrations, steps, err := e.resolve(plan.route)
	e.mu.RUnlock()
	if err != nil {
		return representation.AST{}, err
	}
	// Require all admission capabilities before invoking any provider or edge.
	for _, registered := range registrations {
		if registered.checker == nil {
			return representation.AST{}, fmt.Errorf("%w: checking for %s", ErrMissingCapability, registered.language.ID())
		}
	}
	source, err := input.Source.EntryTree()
	if err != nil { return representation.AST{}, err }
	if err := registrations[0].checker.Check(source); err != nil {
		return representation.AST{}, err
	}
	current, err := source.Snapshot()
	if err != nil {
		return representation.AST{}, err
	}
	for i, step := range steps {
		fail := func(err error) (representation.AST, error) {
			return representation.AST{}, &StepError{Index: i, From: step.from, To: step.to, Err: err}
		}
		if step.transformer.FromID() != step.from || step.transformer.ToID() != step.to {
			return fail(fmt.Errorf("%w: endpoint metadata changed after registration", transform.ErrInvalidTransformer))
		}
		if err := ctx.Err(); err != nil { return fail(err) }
		result, err := step.transformer.Transform(ctx, input)
		if err != nil {
			return fail(err)
		}
		if i + 1 == len(steps) {
			if err := registrations[i+1].checker.Check(result); err != nil { return fail(err) }
		}
		current, err = result.Snapshot()
		if err != nil {
			return fail(err)
		}
		if i + 1 < len(steps) {
			info, err := current.Inspect()
			if err != nil { return fail(err) }
			captured, err := semantic.Capture(ctx, *registrations[i+1].checker, semantic.Input{Entry: ast.NodeRef{Tree: info.ID, Node: info.Root}, Trees: []representation.AST{current}})
			if err != nil { return fail(err) }
			input = transform.Input{Source: captured}
		}
	}
	if err := ctx.Err(); err != nil { return representation.AST{}, err }
	return current, nil
}
