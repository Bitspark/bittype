package engine

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	generation "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
)

// Generation returns this engine's additive construction registry. Concrete
// providers remain explicitly supplied through its typed registration functions.
func (e *Engine) Generation() *generation.Registry { return &e.generation }

// PlanGeneration checks both endpoints and resolves the selected descriptions
// into native intent. It performs no AST construction or implicit analysis.
func (e *Engine) PlanGeneration(ctx context.Context, r generation.Request) (g.GenerationPlan, error) {
	zero := g.GenerationPlan{}
	if ctx == nil {
		return zero, g.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return zero, err
	}
	e.mu.RLock()
	from, fok := e.languages[r.Input.Source.LanguageID()]
	to, tok := e.languages[r.To.ID()]
	e.mu.RUnlock()
	if !fok || !tok {
		return zero, ErrUnknownLanguage
	}
	if !from.language.Compatible(r.Input.Source.Language()) || !to.language.Compatible(r.To) {
		return zero, g.ErrScope
	}
	if from.checker == nil || to.checker == nil {
		return zero, ErrMissingCapability
	}
	for _, tree := range r.Input.Source.Inputs() {
		if err := from.checker.Check(tree); err != nil {
			return zero, err
		}
	}
	return e.generation.PlanGeneration(ctx, r, *to.checker)
}

// Generate is the one-call composition of semantic planning and target construction.
func (e *Engine) Generate(ctx context.Context, r generation.Request) (g.Output, error) {
	plan, err := e.PlanGeneration(ctx, r)
	if err != nil {
		return g.Output{}, err
	}
	return plan.Run(ctx)
}
