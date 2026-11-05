// Package print prints admitted units through an explicitly supplied Syntax
// printer. It never selects paths, reflows syntax or accesses a filesystem.
package print

import (
	"context"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

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
