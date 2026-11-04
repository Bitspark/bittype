// Package filesprint prints admitted units through an explicitly supplied Syntax
// printer. It never selects paths, reflows syntax or accesses a filesystem.
package filesprint

import (
	"context"
	"github.com/Bitspark/bittype/contracts/language/go"
	printcontract "github.com/Bitspark/bittype/contracts/print/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Provider[L any] struct {
	identity language.Provider
	syntax   printcontract.Provider[cst.CST[L], string]
}

func New[L any](identity language.Provider, syntax printcontract.Provider[cst.CST[L], string]) (Provider[L], error) {
	if representation.IsNil(identity) || representation.IsNil(syntax) || identity.ID() != syntax.ID() {
		return Provider[L]{}, printcontract.ErrInvalidPrinter
	}
	return Provider[L]{identity, syntax}, nil
}
func (p Provider[L]) ID() language.ID { return p.identity.ID() }
func (p Provider[L]) Name() string    { return p.identity.Name() }
func (p Provider[L]) Print(ctx context.Context, input source.Source[L]) (files.Files, error) {
	if ctx == nil || !input.Valid() || representation.IsNil(p.syntax) {
		return files.Files{}, source.ErrInvalid
	}
	output := files.Files{}
	for _, unit := range input.Units() {
		if err := ctx.Err(); err != nil {
			return files.Files{}, err
		}
		text, err := p.syntax.Print(ctx, unit.Syntax)
		if err != nil {
			return files.Files{}, err
		}
		output.Entries = append(output.Entries, files.File{ID: unit.File, Path: unit.Path, Text: text})
	}
	return output, nil
}
