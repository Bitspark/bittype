package integration_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	boolean "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/boolean/go"
	declarations "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/declarations/go"
	list "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/list/go"
	record "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/record/go"
	text "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/text/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// typescriptLibraries is the one place in this package that lists the whole
// builtin set of horos-to-typescript mapping libraries. Tests that need the
// whole set call it, so a change to that set edits this file only. A test that
// composes a selection of its own names single libraries itself, as
// horos_generation_composition_test.go does.
func typescriptLibraries(profile string) []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{
		text.Library(profile),
		boolean.Library(profile),
		list.Library(profile),
		record.Library(profile),
		declarations.Library(profile),
	}
}
