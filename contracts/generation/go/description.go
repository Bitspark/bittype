package generation

import (
	"fmt"
	"slices"

	"github.com/Bitspark/bittype/trees/ast/go"
)

type ParameterDescription struct {
	Name        string
	Origin      ast.NodeRef
	RequiredUse bool
}

// GenericConvention describes both sides of a generic, including reordered,
// omitted and additional target parameters. Bindings is indexed by source
// parameter; Arguments is indexed by target parameter. The contribution that
// supplies the declaration states it (row 2).
type GenericConvention[T any] struct {
	Parameters []ParameterDescription
	Bindings   []Expression[T]
	Arguments  []Expression[T]
}

// DeclarationDescription is what the projector of a declaration contribution
// returns: its description of one target declaration. Values are what the
// contribution states about it; see TargetDeclaration.
type DeclarationDescription[T any] struct {
	Name    string
	Form    Form[T]
	Origin  ast.NodeRef
	Generic GenericConvention[T]
	Body    Expression[T]
	Values  []IssuedValue
}

func (d DeclarationDescription[T]) Clone() DeclarationDescription[T] {
	d.Generic.Parameters = slices.Clone(d.Generic.Parameters)
	d.Generic.Bindings = slices.Clone(d.Generic.Bindings)
	d.Generic.Arguments = slices.Clone(d.Generic.Arguments)
	d.Values = slices.Clone(d.Values)
	return d
}

// The two refusals of a request that names another use than the kernel
// minted for the subject, distinguished for diagnosis.
var (
	errNotChild = fmt.Errorf("%w: a request must name a child the kernel minted for the subject", ErrScope)
	errNotBody  = fmt.Errorf("%w: a declaration's rule may request only its body", ErrScope)
)

// childOf reports whether child is a structural child the kernel minted for parent.
func childOf[F any](child, parent Use[F]) bool {
	return child.o != nil && parent.o != nil && child.o.parent == parent.o
}
