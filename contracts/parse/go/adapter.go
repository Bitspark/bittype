package parse

import (
	"context"
	"errors"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

var ErrInvalidParser = errors.New("invalid parse provider")

type SyntaxParser struct {
	definition representation.Language
	run        func(context.Context, Input) (representation.CST, error)
}

func NewSyntaxParser[L any](definition representation.Language, provider Provider[Input, cst.CST[L]]) (SyntaxParser, error) {
	if representation.IsNil(provider) {
		return SyntaxParser{}, ErrInvalidParser
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return SyntaxParser{}, err
	}
	if provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return SyntaxParser{}, ErrInvalidParser
	}
	result := SyntaxParser{definition: definition}
	result.run = func(ctx context.Context, input Input) (representation.CST, error) {

		output, err := provider.Parse(ctx, input)
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
func (p SyntaxParser) ID() language.ID                   { return p.definition.ID() }
func (p SyntaxParser) Name() string                      { return p.definition.Name() }
func (p SyntaxParser) Language() representation.Language { return p.definition }
func (p SyntaxParser) Validate() error {
	if p.run == nil {
		return ErrInvalidParser
	}
	return p.definition.Validate()
}
func (p SyntaxParser) Parse(ctx context.Context, input Input) (representation.CST, error) {
	if err := p.Validate(); err != nil {
		return representation.CST{}, err
	}
	if ctx == nil {
		return representation.CST{}, ErrInvalidParser
	}
	if err := ctx.Err(); err != nil {
		return representation.CST{}, err
	}
	return p.run(ctx, input)
}
