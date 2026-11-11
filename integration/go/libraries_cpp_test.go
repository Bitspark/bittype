package integration_test

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	mapping "github.com/Bitspark/bittype/transform/syntax/horos-to-cpp/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// cppLibraries is the one place in this package that lists the builtin
// horos-to-cpp mapping libraries. Tests that need the whole set call it, so a
// change to how that pair is composed edits this file only.
func cppLibraries(profile string) []g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return mapping.Libraries(profile)
}
