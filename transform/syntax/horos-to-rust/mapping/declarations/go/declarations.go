// Package declarations is the builtin horos-to-rust mapping library for named
// horos declarations. A consumer registers it on its own.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library applies the horos declaration recipe with the Rust model forms
// struct, for a record body, and alias, for every other body.
func Library() rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{
			Profile:      rules.RepresentationProfile,
			Declarations: []g.DeclarationRecipe[h.Language, s.Language, source.Reader]{recipe.Recipe(model.Form(model.FormStruct), model.Form(model.FormAlias))},
		}, nil
	}
}
