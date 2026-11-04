package print

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

var ErrInvalidPrinter = errors.New("invalid print provider")

type SyntaxPrinter struct {
	definition representation.Language
	run        func(context.Context, representation.CST) (string, error)
}

func NewSyntaxPrinter[L any](definition representation.Language, provider Provider[cst.CST[L], string]) (SyntaxPrinter, error) {
	if representation.IsNil(provider) {
		return SyntaxPrinter{}, ErrInvalidPrinter
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return SyntaxPrinter{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return SyntaxPrinter{}, ErrInvalidPrinter
	}
	result := SyntaxPrinter{definition: definition}
	result.run = func(ctx context.Context, input representation.CST) (string, error) {
		value, err := representation.UnpackCST[L](definition, input)
		if err != nil {
			return "", err
		}
		output, err := provider.Print(ctx, value)
		if err != nil {
			return "", err
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
		return output, nil
	}
	return result, nil
}
func (p SyntaxPrinter) ID() language.ID                   { return p.definition.ID() }
func (p SyntaxPrinter) Name() string                      { return p.definition.Name() }
func (p SyntaxPrinter) Language() representation.Language { return p.definition }
func (p SyntaxPrinter) Validate() error {
	if p.run == nil {
		return ErrInvalidPrinter
	}
	return p.definition.Validate()
}
func (p SyntaxPrinter) Print(ctx context.Context, input representation.CST) (string, error) {
	if err := p.Validate(); err != nil {
		return "", err
	}
	if ctx == nil {
		return "", ErrInvalidPrinter
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return p.run(ctx, input)
}
