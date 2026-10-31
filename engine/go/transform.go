package engine

import (
	"fmt"

	"github.com/Bitspark/bittype/language/go"
	"github.com/Bitspark/bittype/transform/go"
)

func (e *Engine) Check(input transform.AST) error {
	registration, err := e.lookup(input.LanguageID())
	if err != nil {
		return err
	}
	return registration.language.Check(input)
}

func (e *Engine) TransformTo(input transform.AST, targetID language.ID) (transform.AST, error) {
	plan, err := e.Plan(input.LanguageID(), targetID)
	if err != nil {
		return transform.AST{}, err
	}
	return e.Execute(input, plan)
}

// Execute uses one route and never retries another mapping after an error.
// Providers receive owned snapshots. Failure returns a zero envelope, even if a
// plugin returned a partial result alongside its error.
func (e *Engine) Execute(input transform.AST, plan Plan) (transform.AST, error) {
	e.mu.RLock()
	registrations, steps, err := e.resolve(plan.route)
	e.mu.RUnlock()
	if err != nil {
		return transform.AST{}, err
	}
	if err := registrations[0].language.Check(input); err != nil {
		return transform.AST{}, err
	}
	current, err := input.Snapshot()
	if err != nil {
		return transform.AST{}, err
	}
	for i, step := range steps {
		fail := func(err error) (transform.AST, error) {
			return transform.AST{}, &StepError{Index: i, From: step.from, To: step.to, Err: err}
		}
		if step.transformer.FromID() != step.from || step.transformer.ToID() != step.to {
			return fail(fmt.Errorf("%w: endpoint metadata changed after registration", transform.ErrInvalidTransformer))
		}
		result, err := step.transformer.Transform(current)
		if err != nil {
			return fail(err)
		}
		if err := registrations[i+1].language.Check(result); err != nil {
			return fail(err)
		}
		current, err = result.Snapshot()
		if err != nil {
			return fail(err)
		}
	}
	return current, nil
}
