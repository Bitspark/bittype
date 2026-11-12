// Package declarations is the horos-to-typescript mapping library for named
// horos declarations. A consumer registers it on its own; it imports no other
// mapping library.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
)

// Library applies the horos declaration recipe with the TypeScript forms
// "interface" for records and "alias" for every other declaration.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return recipe.Library(model.Form(model.FormInterface), model.Form(model.FormAlias))
}
