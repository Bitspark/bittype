// Package declarations is the horos-to-golang mapping library for named
// declarations: a horos record declaration becomes a Go defined type and any
// other declaration a Go alias.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Library returns the declaration mapping. The horos declaration recipe
// states the structural representation profile.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return recipe.Library(model.Form(model.FormDefined), model.Form(model.FormAlias))
}
