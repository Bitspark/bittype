package integration_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	boolean "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/boolean/go"
	declarations "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/declarations/go"
	list "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/list/go"
	record "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/record/go"
	text "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/text/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// cppLibraries is the one place in this package that lists the builtin
// horos-to-cpp mapping libraries. Tests that need the whole set call it, so a
// change to how that pair is composed edits this file only.
func cppLibraries() []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{
		text.Library(),
		boolean.Library(),
		list.Library(),
		record.Library(),
		declarations.Library(),
	}
}

// cppConstructors lists the exported constructor of every builtin
// horos-to-cpp mapping package, in the order of cppLibraries.
func cppConstructors() []any {
	return []any{text.Library, boolean.Library, list.Library, record.Library, declarations.Library}
}
