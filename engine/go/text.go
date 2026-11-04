package engine

import (
	"context"
	"fmt"

	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func (e *Engine) ParseSyntax(ctx context.Context, id language.ID, input parse.Input) (representation.CST, error) {
	r, err := e.lookup(id)
	if err != nil {
		return representation.CST{}, err
	}
	if r.parser == nil {
		return representation.CST{}, fmt.Errorf("%w: parsing for %s", ErrMissingCapability, id)
	}
	return r.parser.Parse(ctx, input)
}
func (e *Engine) FormatSyntax(ctx context.Context, input representation.AST) (representation.CST, error) {
	r, err := e.lookup(input.LanguageID())
	if err != nil {
		return representation.CST{}, err
	}
	if r.formatter == nil {
		return representation.CST{}, fmt.Errorf("%w: formatting for %s", ErrMissingCapability, input.LanguageID())
	}
	return r.formatter.Format(ctx, input)
}
func (e *Engine) PrintSyntax(ctx context.Context, input representation.CST) (string, error) {
	r, err := e.lookup(input.LanguageID())
	if err != nil {
		return "", err
	}
	if r.printer == nil {
		return "", fmt.Errorf("%w: printing for %s", ErrMissingCapability, input.LanguageID())
	}
	return r.printer.Print(ctx, input)
}

// Render captures both required capabilities before any callback. It does not
// parse or invoke a registered checker, and missing printing fails before format.
func (e *Engine) RenderSyntax(ctx context.Context, input representation.AST) (string, error) {
	r, err := e.lookup(input.LanguageID())
	if err != nil {
		return "", err
	}
	if r.formatter == nil || r.printer == nil {
		return "", fmt.Errorf("%w: rendering for %s", ErrMissingCapability, input.LanguageID())
	}
	formatted, err := r.formatter.Format(ctx, input)
	if err != nil {
		return "", err
	}
	return r.printer.Print(ctx, formatted)
}
