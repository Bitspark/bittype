package generation

import (
	"slices"

	"github.com/Bitspark/bittype/trees/ast/go"
)

// Form is a target-owned semantic constructor. Only its issuer can
// interpret its key. This is not an AST kind or a globally registered type kind.
type Form[T any] struct {
	owner *identity
	name  string
}
type FormIssuer[T any] struct{ owner *identity }

func NewFormIssuer[T any]() FormIssuer[T] { return FormIssuer[T]{fresh()} }
func (i FormIssuer[T]) Form(name string) Form[T] {
	if i.owner == nil || name == "" {
		return Form[T]{}
	}
	return Form[T]{i.owner, name}
}
func (i FormIssuer[T]) Read(f Form[T]) (string, bool) {
	return f.name, i.owner != nil && i.owner == f.owner && f.name != ""
}
func (f Form[T]) Name() string { return f.name }
func (f Form[T]) Valid() bool  { return f.owner != nil && f.name != "" }

type ExpressionKind uint8

const (
	FormExpression ExpressionKind = iota + 1
	ParameterExpression
	ArgumentExpression
	DeclarationExpression
	// BinderExpression is the planner's placeholder for a binder of the
	// declaration being planned (bittype#196, O28). The planner replaces
	// it by the target meaning that the declaration's convention gives the
	// binder before a target sees the declaration.
	BinderExpression
)

// Expression is an immutable template. Parameter and Argument indices only have
// meaning in their declaration's single generic convention. DeclarationExpression
// is issued by the planner after it authenticates the source callee.
type Expression[T any] struct {
	origin      ast.NodeRef
	kind        ExpressionKind
	form        Form[T]
	data        []string
	args        []Expression[T]
	index       int
	declaration uint32
}

func Apply[T any](form Form[T], data []string, args ...Expression[T]) Expression[T] {
	return Expression[T]{kind: FormExpression, form: form, data: slices.Clone(data), args: slices.Clone(args)}
}
func Parameter[T any](index int) Expression[T] {
	return Expression[T]{kind: ParameterExpression, index: index}
}
func Argument[T any](index int) Expression[T] {
	return Expression[T]{kind: ArgumentExpression, index: index}
}
func (e Expression[T]) Kind() ExpressionKind       { return e.kind }
func (e Expression[T]) Form() Form[T]              { return e.form }
func (e Expression[T]) Data() []string             { return slices.Clone(e.data) }
func (e Expression[T]) Arguments() []Expression[T] { return slices.Clone(e.args) }
func (e Expression[T]) Index() int                 { return e.index }
func (e Expression[T]) Declaration() uint32        { return e.declaration }

func At[T any](e Expression[T], origin ast.NodeRef) Expression[T] { e.origin = origin; return e }
func (e Expression[T]) Origin() ast.NodeRef                       { return e.origin }

// identity is an allocation-unique owner for issuer-scoped values.
type identity struct{ value byte }

func fresh() *identity { return &identity{value: 1} }
