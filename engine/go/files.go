package engine

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/format/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func (e *Engine) RegisterFilesParser(p parse.FilesRegistration) error {
	if representation.IsNil(p) {
		return parse.ErrInvalidParser
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return e.attach(p.Language(), func(r *registration) error {
		if r.filesParser != nil {
			return ErrDuplicateFilesParser
		}
		r.filesParser = p
		return nil
	})
}
func (e *Engine) RegisterFilesFormatter(p format.FilesRegistration) error {
	if representation.IsNil(p) {
		return format.ErrInvalidFormatter
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return e.attach(p.Language(), func(r *registration) error {
		if r.filesFormatter != nil {
			return ErrDuplicateFilesFormatter
		}
		r.filesFormatter = p
		return nil
	})
}
func (e *Engine) RegisterFilesPrinter(p print.FilesRegistration) error {
	if representation.IsNil(p) {
		return print.ErrInvalidPrinter
	}
	if err := p.Validate(); err != nil {
		return err
	}
	return e.attach(p.Language(), func(r *registration) error {
		if r.filesPrinter != nil {
			return ErrDuplicateFilesPrinter
		}
		r.filesPrinter = p
		return nil
	})
}
func (e *Engine) ParseFiles(ctx context.Context, id language.ID, input parse.FilesInput) (source.Envelope, error) {
	r, err := e.lookup(id)
	if err != nil {
		return source.Envelope{}, err
	}
	if r.filesParser == nil {
		return source.Envelope{}, fmt.Errorf("%w: Files parsing for %s", ErrMissingCapability, id)
	}
	output, err := r.filesParser.ParseFiles(ctx, input)
	if err != nil {
		return source.Envelope{}, err
	}
	if !output.Valid() || !r.language.Compatible(output.Language()) {
		return source.Envelope{}, representation.ErrLanguageMismatch
	}
	return output, nil
}
func (e *Engine) FormatFiles(ctx context.Context, input format.FilesValue) (source.Envelope, error) {
	r, err := e.lookup(input.Abstract.LanguageID())
	if err != nil {
		return source.Envelope{}, err
	}
	if r.filesFormatter == nil {
		return source.Envelope{}, fmt.Errorf("%w: Files formatting for %s", ErrMissingCapability, input.Abstract.LanguageID())
	}
	output, err := r.filesFormatter.FormatFiles(ctx, input)
	if err != nil {
		return source.Envelope{}, err
	}
	if !output.Valid() || !r.language.Compatible(output.Language()) {
		return source.Envelope{}, representation.ErrLanguageMismatch
	}
	return output, nil
}
func (e *Engine) PrintFiles(ctx context.Context, input source.Envelope) (files.Files, error) {
	r, err := e.lookup(input.Language().ID())
	if err != nil {
		return files.Files{}, err
	}
	if r.filesPrinter == nil {
		return files.Files{}, fmt.Errorf("%w: Files printing for %s", ErrMissingCapability, input.Language().ID())
	}
	output, err := r.filesPrinter.PrintFiles(ctx, input)
	if err != nil {
		return files.Files{}, err
	}
	return output.Clone(), nil
}
