// Package declarations is the builtin horos-to-cpp mapping library for named
// horos declarations. It selects the C++ forms struct and alias for them.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Library applies the horos declaration recipe with the forms struct and
// alias.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return recipe.Library(model.Form(model.FormStruct), model.Form(model.FormAlias))
}
