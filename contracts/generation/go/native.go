package generation

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	"slices"
)

// NativeDeclaration is target intent after source mapping and substitution.
// It has no source formal bindings, argument templates or mapping callbacks.
type NativeDeclaration[T any] struct {
	Name       string
	Form       Form[T]
	Origin     ast.NodeRef
	Parameters []NativeParameter
	Body       Expression[T]
}
type NativeParameter struct {
	Name   string
	Origin ast.NodeRef
}

func (d NativeDeclaration[T]) Clone() NativeDeclaration[T] {
	d.Parameters = slices.Clone(d.Parameters)
	return d
}
