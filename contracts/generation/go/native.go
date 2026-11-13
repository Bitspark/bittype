package generation

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
)

// NativeDeclaration is target intent after source mapping and substitution.
// It has no source formal bindings, argument templates or mapping callbacks.
// Values are what the declaration's owner states about it, for example a
// naming intent or the access a helper requires. The kernel carries them to
// the declaration's inventory record and does not read them.
type NativeDeclaration[T any] struct {
	Name       string
	Form       Form[T]
	Origin     ast.NodeRef
	Parameters []NativeParameter
	Body       Expression[T]
	Values     []IssuedValue
}
type NativeParameter struct {
	Name   string
	Origin ast.NodeRef
}

func (d NativeDeclaration[T]) Clone() NativeDeclaration[T] {
	d.Parameters = slices.Clone(d.Parameters)
	d.Values = slices.Clone(d.Values)
	return d
}
