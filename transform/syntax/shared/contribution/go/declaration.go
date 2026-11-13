// Package contribution holds what the contributions of any pair are built
// from: mechanisms over the generation contract that read no source and know
// no language. A contribution of a pair instantiates one with its own types,
// identities and target forms. The policy of a family, which reads its source
// and decides, stays with the source side of its pair.
package contribution

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
)

// Declaration is the description of a target declaration whose target
// parameters are its source binders, one for one and in their order: the
// projector's part of a declaration contribution that keeps the source's
// generic convention. name is the declaration's name in the target, parameters
// are the names of its target parameters, form is the target's declaration
// form, and the body is the target type of the one dependency, the
// declaration's body. Each parameter must be used. It reads the projection
// and nothing of a source.
func Declaration[T any](p g.Projection[T], name string, parameters []string, form g.Form[T]) (g.DeclarationDescription[T], error) {
	if len(parameters) != p.Binders() || p.Dependencies() != 1 {
		return g.DeclarationDescription[T]{}, g.ErrScope
	}
	desc := g.DeclarationDescription[T]{Name: name, Form: form, Origin: p.Origin(), Body: p.Dependency(0).Type}
	for i, name := range parameters {
		desc.Generic.Parameters = append(desc.Generic.Parameters, g.ParameterDescription{Name: name, Origin: p.BinderOrigin(i), RequiredUse: true})
		desc.Generic.Bindings = append(desc.Generic.Bindings, g.Parameter[T](i))
		desc.Generic.Arguments = append(desc.Generic.Arguments, g.Argument[T](i))
	}
	return desc, nil
}
