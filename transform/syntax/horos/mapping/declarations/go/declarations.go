// Package declarations holds the recipe by which a named horos declaration
// becomes a native declaration: its name, the module prefix of a declaration
// outside the root module, the choice between the record form and the alias
// form, and its generic parameters. The horos side of a mapping owns this
// recipe, and it is the same for every target. A target's declarations mapping
// owns only the two native forms it passes in.
//
// The recipe states no profile. A mapping library puts it into its own
// description, under the profile that the library states. The recipe gives a
// native declaration bare generic parameters: it carries no kind constraint
// into the native type.
package declarations

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	naming "github.com/Bitspark/bittype/transform/syntax/horos/mapping/naming/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
)

// Recipe returns the declaration recipe. record is the native form of a
// declaration whose body is a horos record; alias is the form of every other
// declaration.
func Recipe[T any](record, alias g.Form[T]) g.DeclarationRecipe[h.Language, T, source.Reader] {
	return g.DeclarationRecipe[h.Language, T, source.Reader]{ID: "declarations", Expand: func(ctx context.Context, r g.DeclarationRequest[h.Language, T, source.Reader]) (g.DeclarationDescription[T], error) {
		zero := g.DeclarationDescription[T]{}
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
		name, err := naming.Exported(decl.Name)
		if err != nil {
			return zero, src.Unsupported(open.Body, err.Error())
		}
		if !decl.Root {
			name = "M" + decl.ModulePin[:16] + "_" + name
		}
		body, descriptor, err := src.Descriptor(open.Body)
		if err != nil {
			return zero, err
		}
		form := alias
		if descriptor && body.Constructor == "record" {
			form = record
		}
		desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: decl.Origin}
		for i, f := range open.Formals {
			parameter, err := src.Parameter(f)
			if err != nil {
				return zero, err
			}
			pname, err := naming.Exported(parameter.Name)
			if err != nil {
				return zero, err
			}
			desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: pname, Origin: parameter.Origin, RequiredUse: true})
			desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
			desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
		}
		desc.Body, err = r.Mapped(open.Body)
		return desc, err
	}}
}
