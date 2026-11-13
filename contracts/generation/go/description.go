package generation

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Capture[F any] struct {
	Slot Slot
	Use  Use[F]
}

// Pattern fixes selection and declares its captures. Capture reads the source
// meaning of the kernel's view through the source reader R and names children
// of that view; the kernel checks that each captured use is one of them.
type Pattern[F, R any] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, R, UseView[F]) ([]Capture[F], error)
}
type BindingDescription[F, T, R any] struct {
	ID     string
	Input  Pattern[F, R]
	Output Expression[T]
	Expand func(context.Context, Expansion[F, T, R]) (Expression[T], error)
}
type ParameterDescription struct {
	Name        string
	Origin      ast.NodeRef
	RequiredUse bool
}

// Convention describes both sides of a generic, including reordered, omitted
// and additional native parameters. Bindings is indexed by source parameter;
// Arguments is indexed by target parameter.
type GenericConvention[T any] struct {
	Parameters []ParameterDescription
	Bindings   []Expression[T]
	Arguments  []Expression[T]
}

// DeclarationDescription is a recipe's description of one native declaration.
// Values are what the recipe's owner states about it; see NativeDeclaration.
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

type DeclarationRecipe[F, T, R any] struct {
	ID string
	// A zero selector explicitly supplies the structural declaration policy.
	Selector Selector
	Expand   func(context.Context, DeclarationRequest[F, T, R]) (DeclarationDescription[T], error)
}

// HelperRecipe describes a declaration that a mapping demands by key. Its
// declaration states what it requires of the Files stage, for example its
// access, as issued values.
type HelperRecipe[T any] struct {
	Key         HelperKey
	Declaration NativeDeclaration[T]
}

type MappingDescription[F, T, R any] struct {
	Profile      string
	Bindings     []BindingDescription[F, T, R]
	Declarations []DeclarationRecipe[F, T, R]
	Helpers      []HelperRecipe[T]
}

// Describe is called at registration. The registry freezes all returned data.
// No provider receives a target AST builder or is called during Run.
// bittype:dimension Syntax
type MappingLibrary[F, T, R any] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, R], error)
}

// expansionState is one callback's capability. It is active only while the
// callback runs; afterwards every accessor and capture fails.
type expansionState[F, T, R any] struct {
	mu          sync.Mutex
	active      bool
	ctx         context.Context
	reader      R
	use         *UseView[F]
	declaration *DeclarationView[F]
	captures    map[Slot]Use[F]
	limits      Limits
	failure     error
}
type Expansion[F, T, R any] struct{ state *expansionState[F, T, R] }

func (r Expansion[F, T, R]) live() error {
	if r.state == nil {
		return ErrScope
	}
	r.state.mu.Lock()
	defer r.state.mu.Unlock()
	if !r.state.active {
		return ErrScope
	}
	return nil
}

// Reader is the source's interpretation reader, available during the callback.
func (r Expansion[F, T, R]) Reader() (R, error) {
	if err := r.live(); err != nil {
		var z R
		return z, err
	}
	return r.state.reader, nil
}

// Use is the kernel's view of the occurrence a binding rule maps.
func (r Expansion[F, T, R]) Use() (UseView[F], error) {
	if err := r.live(); err != nil {
		return UseView[F]{}, err
	}
	if r.state.use == nil {
		return UseView[F]{}, ErrScope
	}
	return r.state.use.Clone(), nil
}

// Mapped captures a use for the planner to represent after the callback. A
// binding rule may capture only children of its occurrence; a declaration
// recipe only its body.
func (r Expansion[F, T, R]) Mapped(use Use[F]) (Expression[T], error) {
	s := r.state
	if s == nil {
		return Expression[T]{}, ErrScope
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.active {
		return Expression[T]{}, ErrScope
	}
	if err := s.ctx.Err(); err != nil {
		s.failure = err
		return Expression[T]{}, err
	}
	if len(s.captures) >= s.limits.References {
		s.failure = ErrLimit
		return Expression[T]{}, ErrLimit
	}
	allowed := false
	if s.declaration != nil {
		allowed = use == s.declaration.Body
	} else if s.use != nil {
		allowed = childOf(use, s.use.Use)
	}
	if !allowed {
		err := errNotChild
		if s.declaration != nil {
			err = errNotBody
		}
		s.failure = err
		return Expression[T]{}, err
	}
	slot := TypeSlot("capture")
	s.captures[slot] = use
	return Mapped[T](slot), nil
}

type DeclarationRequest[F, T, R any] struct{ Expansion[F, T, R] }

// Declaration is the kernel's view of the declaration being planned.
func (r DeclarationRequest[F, T, R]) Declaration() (DeclarationView[F], error) {
	if err := r.live(); err != nil {
		return DeclarationView[F]{}, err
	}
	if r.state.declaration == nil {
		return DeclarationView[F]{}, ErrScope
	}
	return r.state.declaration.Clone(), nil
}

// The two capture refusals, distinguished for diagnosis.
var (
	errNotChild = fmt.Errorf("%w: a capture must be a child the kernel minted for the mapped occurrence", ErrScope)
	errNotBody  = fmt.Errorf("%w: a declaration recipe may map only its body", ErrScope)
)

// childOf reports whether child is a structural child the kernel minted for parent.
func childOf[F any](child, parent Use[F]) bool {
	return child.o != nil && parent.o != nil && child.o.parent == parent.o
}
