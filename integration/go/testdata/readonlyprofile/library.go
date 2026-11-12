// Package readonlyprofile is a test contribution: one horos-to-typescript
// mapping library under a generation profile of its own. It is written the way
// an author of a second profile would write one. It states its own profile
// identity and builds its description from the parts that the horos mapping
// side offers to every profile: the function type mapping.Library, the common
// patterns, the common declaration recipe, the horos source reader and the
// TypeScript model.
//
// The family differs from the builtin structural family, so that a request
// shows which of the two it received:
//
//   - a horos list becomes ReadonlyArray<E> and not E[];
//   - every declaration becomes an alias, also one whose body is a record;
//   - records are not mapped at all.
//
// The horos mapping package and the generation registry know nothing of this
// package or of its profile.
package readonlyprofile

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	recipe "github.com/Bitspark/bittype/transform/syntax/horos/mapping/declarations/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Profile is this contribution's own promise: lists are read-only views. Its
// author states it and answers for it.
const Profile = "integration-fixture/readonly-collections/1"

// ListBinding is the identity of the list binding. A test reads it in the
// inspectable plan.
const ListBinding = "integration-fixture/readonly-list/typescript/1"

// Library returns the contribution as a mapping.Library function value, the
// same form that the builtin mappings have. It takes no profile: the profile
// is part of what the library describes.
func Library() rules.Library[s.Language] {
	return func(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
		zero := g.MappingDescription[h.Language, s.Language, source.Reader]{}
		text, err := patterns.Head(from, "scalar", "utf8-text")
		if err != nil {
			return zero, err
		}
		boolean, err := patterns.Head(from, "scalar", "bool")
		if err != nil {
			return zero, err
		}
		element := g.TypeSlot("E")
		list, err := patterns.Element(from, element)
		if err != nil {
			return zero, err
		}
		alias := model.Form(model.FormAlias)
		return g.MappingDescription[h.Language, s.Language, source.Reader]{
			Profile: Profile,
			Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{
				{ID: "integration-fixture/readonly-text/typescript/1", Input: text, Output: model.Primitive(model.PrimitiveString)},
				{ID: "integration-fixture/readonly-bool/typescript/1", Input: boolean, Output: model.Primitive(model.PrimitiveBool)},
				{ID: ListBinding, Input: list, Output: model.External("", "ReadonlyArray", g.Mapped[s.Language](element))},
			},
			// The common recipe states no profile. It is used here with the
			// alias form for every declaration.
			Declarations: []g.DeclarationRecipe[h.Language, s.Language, source.Reader]{recipe.Recipe(alias, alias)},
		}, nil
	}
}
