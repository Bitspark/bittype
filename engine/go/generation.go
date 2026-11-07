package engine

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	generation "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Generation returns this engine's additive construction registry. Concrete
// providers remain explicitly supplied through its typed registration functions.
func (e *Engine) Generation() *generation.Registry { return &e.generation }

// Generate checks both endpoints, composes the selected providers once, then
// admits the complete output. It performs no resolution, analysis or routing.
func (e *Engine) Generate(ctx context.Context, r generation.Request) (representation.ASTSet, error) {
	zero := representation.ASTSet{}
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
	out, err := e.generation.Generate(ctx, r)
	if err != nil {
		return zero, err
	}
	captured, err := semantic.Capture(ctx, *to.checker, semantic.Input{Entry: out.Entry, Trees: out.Trees})
	if err != nil {
		return zero, err
	}
	if err := ctx.Err(); err != nil {
		return zero, err
	}
	return representation.ASTSet{Entry: captured.Entry(), Trees: captured.Inputs()}, nil
}
