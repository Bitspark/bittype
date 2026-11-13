package generation

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
)

// TargetDeclaration is target intent after source mapping and substitution.
// It has no source formal bindings, argument templates or mapping callbacks.
// Values are what the declaration's owner states about it, for example a
// naming intent or the access a helper requires. The kernel carries them to
// the declaration's inventory record and does not read them.
type TargetDeclaration[T any] struct {
	Name       string
	Form       Form[T]
	Origin     ast.NodeRef
	Parameters []TargetParameter
	Body       Expression[T]
	Values     []IssuedValue
}
type TargetParameter struct {
	Name   string
	Origin ast.NodeRef
}

func (d TargetDeclaration[T]) Clone() TargetDeclaration[T] {
	d.Parameters = slices.Clone(d.Parameters)
	d.Values = slices.Clone(d.Values)
	return d
}
