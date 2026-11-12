package integration_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	boolean "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/boolean/go"
	declarations "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/declarations/go"
	list "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/list/go"
	record "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/record/go"
	text "github.com/Bitspark/bittype/transform/syntax/horos-to-rust/mapping/text/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// rustLibraries is the one place in this package that lists the builtin
// horos-to-rust mapping libraries. Tests that need the whole set call it, so a
// change to how that pair is composed edits this file only.
func rustLibraries() []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{
		text.Library(),
		boolean.Library(),
		list.Library(),
		record.Library(),
		declarations.Library(),
	}
}

// rustConstructors lists the exported constructor of every builtin
// horos-to-rust mapping package, in the order of rustLibraries.
func rustConstructors() []any {
	return []any{text.Library, boolean.Library, list.Library, record.Library, declarations.Library}
}
