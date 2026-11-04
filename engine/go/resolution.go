package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrDuplicateResolver = errors.New("resolver is already registered")

func (e *Engine) Capture(ctx context.Context, input semantic.Input) (semantic.Snapshot, error) {
	if ctx == nil || len(input.Trees) == 0 {
		return semantic.Snapshot{}, semantic.ErrInvalidRequest
	}
	r, err := e.lookup(input.LanguageID())
	if err != nil {
		return semantic.Snapshot{}, err
	}
	if r.checker == nil {
		return semantic.Snapshot{}, fmt.Errorf("%w: checking for %s", ErrMissingCapability, input.LanguageID())
	}
	return semantic.Capture(ctx, *r.checker, input)
}

func (e *Engine) RegisterResolver(resolver resolution.Resolver) error {
	if err := resolver.Validate(); err != nil {
		return err
	}
	definition := resolver.Language()
	e.mu.Lock()
	defer e.mu.Unlock()
	r, ok := e.languages[definition.ID()]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownLanguage, definition.ID())
	}
	if !r.language.Compatible(definition) {
		return representation.ErrLanguageMismatch
	}
	if r.resolver != nil {
		return ErrDuplicateResolver
	}
	r.resolver = &resolver
	e.languages[definition.ID()] = r
	return nil
}

// Resolver exposes captured capability metadata without installing a default.
func (e *Engine) Resolver(id language.ID) (resolution.Resolver, error) {
	r, err := e.lookup(id)
	if err != nil {
		return resolution.Resolver{}, err
	}
	if r.resolver == nil {
		return resolution.Resolver{}, fmt.Errorf("%w: resolution for %s", ErrMissingCapability, id)
	}
	return *r.resolver, nil
}

func (e *Engine) Resolve(ctx context.Context, input semantic.Snapshot, options semantic.Options) (semantic.Result, error) {
	if ctx == nil || !input.Valid() {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	r, err := e.Resolver(input.LanguageID())
	if err != nil {
		return semantic.Result{}, err
	}
	return r.Resolve(ctx, input, options)
}
