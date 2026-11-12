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
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library returns the declaration mapping: the horos declaration recipe with
// the two Go forms. It states the structural representation profile itself.
func Library() rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{
			Profile:      rules.RepresentationProfile,
			Declarations: []g.DeclarationRecipe[h.Language, s.Language, source.Reader]{recipe.Recipe(model.Form(model.FormDefined), model.Form(model.FormAlias))},
		}, nil
	}
}
