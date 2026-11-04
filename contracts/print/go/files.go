package print

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

type FilesPrinter[L any] struct {
	definition representation.Language
	provider   Provider[source.Source[L], files.Files]
}
type FilesRegistration interface {
	language.Provider
	Language() representation.Language
	Validate() error
	PrintFiles(context.Context, source.Envelope) (files.Files, error)
}

func NewFilesPrinter[L any](definition representation.Language, provider Provider[source.Source[L], files.Files]) (FilesPrinter[L], error) {
	if representation.IsNil(provider) || provider.ID() != definition.ID() || strings.TrimSpace(provider.Name()) == "" {
		return FilesPrinter[L]{}, ErrInvalidPrinter
	}
	if err := representation.CheckMarker[L](definition); err != nil {
		return FilesPrinter[L]{}, err
	}
	return FilesPrinter[L]{definition, provider}, nil
}
func (p FilesPrinter[L]) ID() language.ID                   { return p.definition.ID() }
func (p FilesPrinter[L]) Name() string                      { return p.definition.Name() }
func (p FilesPrinter[L]) Language() representation.Language { return p.definition }
func (p FilesPrinter[L]) Validate() error {
	if representation.IsNil(p.provider) {
		return ErrInvalidPrinter
	}
	return p.definition.Validate()
}
func (p FilesPrinter[L]) Print(ctx context.Context, input source.Source[L]) (files.Files, error) {
	if err := p.Validate(); err != nil {
		return files.Files{}, err
	}
	if ctx == nil {
		return files.Files{}, source.ErrInvalid
	}
	if err := ctx.Err(); err != nil {
		return files.Files{}, err
	}
	if !input.Valid() || !p.definition.Compatible(input.Snapshot().Language()) {
		return files.Files{}, source.ErrInvalid
	}
	output, err := p.provider.Print(ctx, input)
	if err != nil {
		return files.Files{}, err
	}
	if err := ctx.Err(); err != nil {
		return files.Files{}, err
	}
	captured := input.Snapshot().Files()
	if captured == nil || !source.EqualFiles(captured.Files, output) {
		return files.Files{}, source.ErrInvalid
	}
	return output.Clone(), nil
}
func (p FilesPrinter[L]) PrintFiles(ctx context.Context, input source.Envelope) (files.Files, error) {
	typed, err := source.Unpack[L](p.definition, input)
	if err != nil {
		return files.Files{}, err
	}
	return p.Print(ctx, typed)
}
