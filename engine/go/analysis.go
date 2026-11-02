package engine

import (
	"context"
	"errors"
	"fmt"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

var ErrDuplicateAnalyzer = errors.New("semantic analyzer is already registered")

// RegisterAnalyzer explicitly attaches one analysis capability to an existing
// language entry. It neither replaces syntax/text capabilities nor finds a route
// through another language. Construction has already captured provider metadata.
func (e *Engine) RegisterAnalyzer(analyzer analysis.Analyzer) error {
	if err := analyzer.Validate(); err != nil {
		return err
	}
	definition := analyzer.Language()
	e.mu.Lock()
	defer e.mu.Unlock()
	r, ok := e.languages[definition.ID()]
	if !ok {
		return fmt.Errorf("%w: %s", ErrUnknownLanguage, definition.ID())
	}
	if !r.language.Compatible(definition) {
		return representation.ErrLanguageMismatch
	}
	if r.analyzer != nil {
		return fmt.Errorf("%w: %s", ErrDuplicateAnalyzer, definition.ID())
	}
	r.analyzer = &analyzer
	e.languages[definition.ID()] = r
	return nil
}

// Analyze selects an explicit analyzer by the input language. lookup releases
// the registry lock before admission checks or provider callbacks. Expected
// semantic rejection is a Result; non-nil errors always discard all output.
func (e *Engine) Analyze(ctx context.Context, input analysis.Input) (semantic.Result, error) {
	if ctx == nil || !input.Snapshot.Valid() {
		return semantic.Result{}, semantic.ErrInvalidRequest
	}
	r, err := e.lookup(input.Snapshot.LanguageID())
	if err != nil {
		return semantic.Result{}, err
	}
	if r.analyzer == nil {
		return semantic.Result{}, fmt.Errorf("%w: semantic analysis for %s", ErrMissingCapability, input.Snapshot.LanguageID())
	}
	result, err := r.analyzer.Analyze(ctx, input)
	if err != nil {
		return semantic.Result{}, err
	}
	return result, nil
}
