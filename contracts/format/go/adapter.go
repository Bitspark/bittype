package format

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

var ErrInvalidFormatter = errors.New("invalid format provider")

type SyntaxFormatter struct {
	definition representation.Language
	run        func(context.Context, representation.AST) (representation.CST, error)
}

func NewSyntaxFormatter[L any](definition representation.Language, provider Provider[ast.View[L], cst.CST[L]]) (SyntaxFormatter, error) {
	if representation.IsNil(provider) {
		return SyntaxFormatter{}, ErrInvalidFormatter
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return SyntaxFormatter{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return SyntaxFormatter{}, ErrInvalidFormatter
	}
	result := SyntaxFormatter{definition: definition}
	result.run = func(ctx context.Context, input representation.AST) (representation.CST, error) {
		value, err := representation.Unpack[L](definition, input)
		if err != nil {
			return representation.CST{}, err
		}
		output, err := provider.Format(ctx, value)
		if err != nil {
			return representation.CST{}, err
		}
		if err := ctx.Err(); err != nil {
			return representation.CST{}, err
		}
		return representation.PackCST(definition, output)
	}
	return result, nil
}
func (p SyntaxFormatter) ID() language.ID                   { return p.definition.ID() }
func (p SyntaxFormatter) Name() string                      { return p.definition.Name() }
func (p SyntaxFormatter) Language() representation.Language { return p.definition }
func (p SyntaxFormatter) Validate() error {
	if p.run == nil {
		return ErrInvalidFormatter
	}
	return p.definition.Validate()
}
func (p SyntaxFormatter) Format(ctx context.Context, input representation.AST) (representation.CST, error) {
	if err := p.Validate(); err != nil {
		return representation.CST{}, err
	}
	if ctx == nil {
		return representation.CST{}, ErrInvalidFormatter
	}
	if err := ctx.Err(); err != nil {
		return representation.CST{}, err
	}
	return p.run(ctx, input)
}
