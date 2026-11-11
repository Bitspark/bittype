package integration_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	mapping "github.com/Bitspark/bittype/transform/syntax/horos-to-golang/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// golangLibraries is the one place in this package that lists the builtin
// horos-to-golang mapping libraries. Tests that need the whole set call it, so a
// change to how that pair is composed edits this file only.
func golangLibraries(profile string) []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return mapping.Libraries(profile)
}
