// Package readonlyprofile is a test contribution: one horos-to-typescript
// mapping library under a generation profile of its own. It is written the way
// an independent author would write one. It implements the generation contract
// itself, states its own profile identity, and uses the horos source reader,
// the common horos patterns and naming, and the TypeScript model.
//
// The family differs from the builtin structural family, so that a request
// shows which of the two it received:
//
//   - a horos list becomes ReadonlyArray<E> and not E[];
//   - every declaration becomes an alias, and a generic declaration is refused;
//   - records are not mapped at all.
//
// The horos mapping package and the generation registry know nothing of this
// package or of its profile.
package readonlyprofile

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	patterns "github.com/Bitspark/bittype/transform/syntax/horos/mapping/pattern/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// Profile is this contribution's own promise: lists are read-only views. Its
// author states it and answers for it.
const Profile = "integration-fixture/readonly-collections/1"

// The binding identity of the list mapping, for a test that reads the plan.
const ListBinding = "integration-fixture/readonly-list/typescript/1"

type library struct{}

// Library returns the contribution. It takes no profile: the profile is part
// of what the library describes.
func Library() g.MappingLibrary[h.Language, s.Language, source.Reader] { return library{} }

func (library) Describe(from, to representation.Language) (g.MappingDescription[h.Language, s.Language, source.Reader], error) {
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
	return g.MappingDescription[h.Language, s.Language, source.Reader]{
		Profile: Profile,
		Bindings: []g.BindingDescription[h.Language, s.Language, source.Reader]{
			{ID: "integration-fixture/readonly-text/typescript/1", Input: text, Output: model.Primitive(model.PrimitiveString)},
			{ID: "integration-fixture/readonly-bool/typescript/1", Input: boolean, Output: model.Primitive(model.PrimitiveBool)},
			{ID: ListBinding, Input: list, Output: model.External("", "ReadonlyArray", g.Mapped[s.Language](element))},
		},
		Declarations: []g.DeclarationRecipe[h.Language, s.Language, source.Reader]{{ID: "integration-fixture/readonly-aliases/typescript/1", Expand: alias}},
	}, nil
}

// alias makes a named, non-generic declaration of the entry module an alias of
// its mapped body.
func alias(ctx context.Context, r g.DeclarationRequest[h.Language, s.Language, source.Reader]) (g.DeclarationDescription[s.Language], error) {
	zero := g.DeclarationDescription[s.Language]{}
	src, err := r.Reader()
	if err != nil {
		return zero, err
	}
	open, err := r.Declaration()
	if err != nil {
		return zero, err
	}
	decl, err := src.Declaration(open.Declaration)
	if err != nil {
		return zero, err
	}
	if len(open.Formals) != 0 {
		return zero, src.Unsupported(open.Body, "the read-only fixture profile maps no generic declaration")
	}
	if !decl.Root {
		return zero, src.Unsupported(open.Body, "the read-only fixture profile maps no imported declaration")
	}
	name, err := naming.Exported(decl.Name)
	if err != nil {
		return zero, src.Unsupported(open.Body, err.Error())
	}
	desc := g.DeclarationDescription[s.Language]{Name: name, Form: model.Form(model.FormAlias), Origin: decl.Origin}
	desc.Body, err = r.Mapped(open.Body)
	return desc, err
}
