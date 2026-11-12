// Package declarations is the builtin horos-to-rust mapping library for named
// horos declarations. A consumer registers it on its own.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
)

// Library applies the horos declaration recipe with the Rust model forms
// struct, for a record body, and alias, for every other body.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] {
	return recipe.Library(model.Form(model.FormStruct), model.Form(model.FormAlias))
}
