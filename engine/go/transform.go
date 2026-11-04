package engine

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/semantic/go"
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

func (e *Engine) TransformTo(ctx context.Context, input transform.Input, targetID language.ID) (representation.ASTSet, error) {
	plan, err := e.Plan(input.Source.LanguageID(), targetID)
	if err != nil {
		return representation.ASTSet{}, err
	}
	return e.Execute(ctx, input, plan)
}

// Execute uses one route and never retries another mapping after an error.
// Providers receive owned snapshots. Failure returns a zero envelope, even if a
// plugin returned a partial result alongside its error.
func (e *Engine) Execute(ctx context.Context, input transform.Input, plan Plan) (representation.ASTSet, error) {
	if ctx == nil {
		return representation.ASTSet{}, semantic.ErrInvalidRequest
	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	e.mu.RLock()
	registrations, steps, err := e.resolve(plan.route)
	e.mu.RUnlock()
	if err != nil {
		return representation.ASTSet{}, err
	}
	// Require all admission capabilities before invoking any provider or edge.
	for _, registered := range registrations {
		if registered.checker == nil {
			return representation.ASTSet{}, fmt.Errorf("%w: checking for %s", ErrMissingCapability, registered.language.ID())
		}
	}
	current, err := (representation.ASTSet{Entry: input.Source.Entry(), Trees: input.Source.Inputs()}).Snapshot()
	if err != nil {
		return representation.ASTSet{}, err
	}
	for _, tree := range current.Trees {
		if err := registrations[0].checker.Check(tree); err != nil {
			return representation.ASTSet{}, err
		}
	}
	for i, step := range steps {
		fail := func(err error) (representation.ASTSet, error) {
			return representation.ASTSet{}, &StepError{Index: i, From: step.from, To: step.to, Err: err}
		}
		if step.transformer.FromID() != step.from || step.transformer.ToID() != step.to {
			return fail(fmt.Errorf("%w: endpoint metadata changed after registration", transform.ErrInvalidTransformer))
		}
		if err := ctx.Err(); err != nil {
			return fail(err)
		}
		result, err := step.transformer.Transform(ctx, input)
		if err != nil {
			return fail(err)
		}
		current, err = result.Snapshot()
		if err != nil {
			return fail(err)
		}
		captured, err := semantic.Capture(ctx, *registrations[i+1].checker, semantic.Input{Entry: current.Entry, Trees: current.Trees})
		if err != nil {
			return fail(err)
		}
		// Receipts never survive a route leg. The next transformer must receive its
		// own required evidence explicitly; the engine does not discover analyzers.
		input = transform.Input{Source: captured}

	}
	if err := ctx.Err(); err != nil {
		return representation.ASTSet{}, err
	}
	return current, nil
}
