package integration_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/boolean/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/declarations/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/list/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/optional/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/record/go"
	"github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/text/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// golangLibraries is the one place in this package that lists the whole builtin
// set of horos-to-golang mapping libraries. Tests that need the whole set call
// it, so a change to that set edits this file only. A test that composes a
// selection of its own names single libraries itself, as
// horos_generation_golang_optional_test.go does.
//
// The record library wraps an optional field in a helper and the optional
// library supplies that helper. Neither creates the helper key: the consumer
// creates it once and passes the same key to both.
func golangLibraries() []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	optionalKey := g.NewHelperKey("horos/optional")
	return []g.MappingLibrary[h.Language, s.Language, source.Reader]{
		text.Library(),
		boolean.Library(),
		list.Library(),
		record.Library(optionalKey),
		declarations.Library(),
		optional.Library(optionalKey),
	}
}

// golangConstructors lists the exported constructor of every builtin
// horos-to-golang mapping package, in the order of golangLibraries.
func golangConstructors() []any {
	return []any{text.Library, boolean.Library, list.Library, record.Library, declarations.Library, optional.Library}
}
