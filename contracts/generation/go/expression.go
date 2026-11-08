package generation

import (
	"slices"

	"github.com/Bitspark/bittype/trees/ast/go"
)

// Form is a target-owned native semantic constructor. Only its issuer can
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

// Slot identity is independent of its display name. One description declares
// its slots, and captures supply the exact source occurrences for them.
type Slot struct {
	owner *identity
	name  string
}

func TypeSlot(name string) Slot { return Slot{fresh(), name} }
func (s Slot) Name() string     { return s.name }
func (s Slot) Valid() bool      { return s.owner != nil && s.name != "" }

type HelperKey struct {
	owner *identity
	name  string
}

func NewHelperKey(name string) HelperKey { return HelperKey{fresh(), name} }
func (k HelperKey) Name() string         { return k.name }
func (k HelperKey) Valid() bool          { return k.owner != nil && k.name != "" }

type ExpressionKind uint8

const (
	NativeExpression ExpressionKind = iota + 1
	CapturedExpression
	ParameterExpression
	ArgumentExpression
	HelperExpression
	DeclarationExpression
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
	slot        Slot
	index       int
	helper      HelperKey
	declaration uint32
}

func Native[T any](form Form[T], data []string, args ...Expression[T]) Expression[T] {
	return Expression[T]{kind: NativeExpression, form: form, data: slices.Clone(data), args: slices.Clone(args)}
}
func Mapped[T any](slot Slot) Expression[T] {
	return Expression[T]{kind: CapturedExpression, slot: slot}
}
func Parameter[T any](index int) Expression[T] {
	return Expression[T]{kind: ParameterExpression, index: index}
}
func Argument[T any](index int) Expression[T] {
	return Expression[T]{kind: ArgumentExpression, index: index}
}
func Helper[T any](key HelperKey, args ...Expression[T]) Expression[T] {
	return Expression[T]{kind: HelperExpression, helper: key, args: slices.Clone(args)}
}
func (e Expression[T]) Kind() ExpressionKind       { return e.kind }
func (e Expression[T]) Form() Form[T]              { return e.form }
func (e Expression[T]) Data() []string             { return slices.Clone(e.data) }
func (e Expression[T]) Arguments() []Expression[T] { return slices.Clone(e.args) }
func (e Expression[T]) Index() int                 { return e.index }
func (e Expression[T]) Declaration() uint32        { return e.declaration }
func (e Expression[T]) Slot() Slot                 { return e.slot }
func (e Expression[T]) Helper() HelperKey          { return e.helper }

func At[T any](e Expression[T], origin ast.NodeRef) Expression[T] { e.origin = origin; return e }
func (e Expression[T]) Origin() ast.NodeRef                       { return e.origin }
