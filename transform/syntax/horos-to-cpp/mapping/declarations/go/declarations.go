// Package declarations is the builtin horos-to-cpp mapping library for named
// horos declarations. It selects the C++ forms struct and alias for them.
package declarations

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	model "github.com/Bitspark/bittype/transform/syntax/cpp/model/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Library applies the horos declaration recipe with the forms struct and
// alias.
func Library() rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		return g.MappingDescription[h.Language, s.Language, source.Reader]{
			Profile:      rules.RepresentationProfile,
			Declarations: []g.DeclarationRecipe[h.Language, s.Language, source.Reader]{recipe.Recipe(model.Form(model.FormStruct), model.Form(model.FormAlias))},
		}, nil
	}
}
