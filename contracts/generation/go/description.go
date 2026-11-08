package generation

import (
	"context"
	"slices"
	"sync"

	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Capture[F any] struct {
	Slot Slot
	Use  Use[F]
}

// Pattern fixes selection and declares its captures. Capture interprets the
// selected source meaning only; the kernel checks each child relation.
type Pattern[F any, S SourceSession[F]] struct {
	Selector Selector
	Slots    []Slot
	Capture  func(context.Context, S, UseData[F]) ([]Capture[F], error)
}
type BindingDescription[F, T any, S SourceSession[F]] struct {
	ID     string
	Input  Pattern[F, S]
	Output Expression[T]
	Expand func(context.Context, Expansion[F, T, S]) (Expression[T], error)
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
type DeclarationDescription[T any] struct {
	Name    string
	Form    Form[T]
	Origin  ast.NodeRef
	Generic GenericConvention[T]
	Body    Expression[T]
}

func (d DeclarationDescription[T]) Clone() DeclarationDescription[T] {
	d.Generic.Parameters = slices.Clone(d.Generic.Parameters)
	d.Generic.Bindings = slices.Clone(d.Generic.Bindings)
	d.Generic.Arguments = slices.Clone(d.Generic.Arguments)
	return d
}

type DeclarationRecipe[F, T any, S SourceSession[F]] struct {
	ID string
	// A zero selector explicitly supplies the structural declaration policy.
	Selector Selector
	Expand   func(context.Context, DeclarationRequest[F, T, S]) (DeclarationDescription[T], error)
}
type HelperRecipe[T any] struct {
	Key         HelperKey
	Declaration NativeDeclaration[T]
}

type MappingDescription[F, T any, S SourceSession[F]] struct {
	Profile      string
	Bindings     []BindingDescription[F, T, S]
	Declarations []DeclarationRecipe[F, T, S]
	Helpers      []HelperRecipe[T]
}

// Describe is called at registration. The registry freezes all returned data.
// No provider receives a target AST builder or is called during Run.
// bittype:dimension Syntax
type MappingLibrary[F, T any, S SourceSession[F]] interface {
	Describe(from, to representation.Language) (MappingDescription[F, T, S], error)
}

type expansionState[F, T any, S SourceSession[F]] struct {
	mu       sync.Mutex
	active   bool
	ctx      context.Context
	source   S
	use      UseData[F]
	open     OpenData[F]
	captures map[Slot]Use[F]
	limits   Limits
	failure  error
}
type Expansion[F, T any, S SourceSession[F]] struct{ state *expansionState[F, T, S] }

func (r Expansion[F, T, S]) Source() (S, error) {
	if r.state == nil {
		var z S
		return z, ErrScope
	}
	r.state.mu.Lock()
	defer r.state.mu.Unlock()
	if !r.state.active {
		var z S
		return z, ErrScope
	}
	return r.state.source, nil
}
func (r Expansion[F, T, S]) Use() (UseData[F], error) {
	if _, err := r.Source(); err != nil {
		return UseData[F]{}, err
	}
	return r.state.use.Clone(), nil
}
func (r Expansion[F, T, S]) Mapped(use Use[F]) (Expression[T], error) {
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
	var err error
	if s.open.Declaration.Valid() {
		if use != s.open.Body {
			err = ErrScope
		} // declaration recipes can map their body only
	} else {
		err = s.source.CheckChild(s.ctx, s.use.Environment, s.use.Use, use)
	}
	if err != nil {
		s.failure = err
		return Expression[T]{}, err
	}
	slot := TypeSlot("capture")
	s.captures[slot] = use
	return Mapped[T](slot), nil
}

type DeclarationRequest[F, T any, S SourceSession[F]] struct{ Expansion[F, T, S] }

func (r DeclarationRequest[F, T, S]) Open() (OpenData[F], error) {
	if _, err := r.Source(); err != nil {
		return OpenData[F]{}, err
	}
	return r.state.open.Clone(), nil
}
