// Package declarations is the horos-to-typescript mapping library for named
// horos declarations. A consumer registers it on its own; it imports no other
// mapping library.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library applies the horos declaration recipe with the TypeScript forms
// "interface" for records and "alias" for every other declaration.
func Library() rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{
			Profile:      rules.RepresentationProfile,
			Declarations: []g.DeclarationRecipe[h.Language, s.Language, source.Reader]{recipe.Recipe(model.Form(model.FormInterface), model.Form(model.FormAlias))},
		}, nil
	}
}
