// Package target realizes a closed native program using supplied grammar rules.
package target

import (
	"context"
	"slices"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Fragment[T any] struct {
	owner *Store[T]
	index uint32
}
type Child[T any] struct {
	Role     string
	Fragment Fragment[T]
}
type Symbol[T any] struct {
	owner       *Store[T]
	declaration uint32
	parameter   int
	name        string
}

func (s Symbol[T]) Name() string   { return s.name }
func (s Symbol[T]) IsFormal() bool { return s.parameter >= 0 }

type fragment[T any] struct {
	node     ast.Node[T]
	children []Child[T]
	symbol   Symbol[T]
	role     string
	external *g.ExternalUse
}
type Store[T any] struct {
	ctx     context.Context
	limits  g.Limits
	records []fragment[T]
	bytes   int
}

func New[T any](ctx context.Context, limits g.Limits) *Store[T] {
	return &Store[T]{ctx: ctx, limits: limits}
}
func (s *Store[T]) Node(n ast.Node[T], children ...Child[T]) (Fragment[T], error) {
	if err := s.ctx.Err(); err != nil {
		return Fragment[T]{}, err
	}
	if len(n.Children) != 0 || len(n.References) != 0 || n.Kind == "" {
		return Fragment[T]{}, g.ErrInvalid
	}
	if len(s.records) >= s.limits.OutputNodes {
		return Fragment[T]{}, g.ErrLimit
	}
	size := len(n.Kind)
	for _, a := range n.Attributes {
		size += len(a.Name) + len(a.Text)
	}
	for _, c := range children {
		if _, err := s.read(c.Fragment); err != nil {
			return Fragment[T]{}, err
		}
		size += len(c.Role)
	}
	if size > s.limits.OutputBytes-s.bytes {
		return Fragment[T]{}, g.ErrLimit
	}
	s.bytes += size
	n.Attributes = slices.Clone(n.Attributes)
	n.Origin.Nodes = slices.Clone(n.Origin.Nodes)
	n.Origin.Spans = slices.Clone(n.Origin.Spans)
	s.records = append(s.records, fragment[T]{node: n, children: slices.Clone(children)})
	return Fragment[T]{s, uint32(len(s.records))}, nil
}
func (s *Store[T]) read(f Fragment[T]) (fragment[T], error) {
	if f.owner != s || f.index == 0 || int(f.index) > len(s.records) {
		return fragment[T]{}, g.ErrScope
	}
	return s.records[f.index-1], nil
}
func (s *Store[T]) Inspect(f Fragment[T]) (ast.Node[T], []Child[T], error) {
	r, e := s.read(f)
	n := r.node
	n.Attributes = slices.Clone(n.Attributes)
	n.Origin.Nodes = slices.Clone(n.Origin.Nodes)
	n.Origin.Spans = slices.Clone(n.Origin.Spans)
	return n, slices.Clone(r.children), e
}
func (s *Store[T]) Link(f Fragment[T], symbol Symbol[T], role string) (Fragment[T], error) {
	r, e := s.read(f)
	if e != nil {
		return f, e
	}
	if symbol.owner != s || symbol.declaration == 0 || role == "" {
		return f, g.ErrScope
	}
	out, e := s.Node(r.node, r.children...)
	if e != nil {
		return out, e
	}
	s.records[out.index-1].symbol = symbol
	s.records[out.index-1].role = role
	return out, nil
}

func (s *Store[T]) External(f Fragment[T], use g.ExternalUse) (Fragment[T], error) {
	r, e := s.read(f)
	if e != nil {
		return f, e
	}
	if !use.At.Role.Valid() || use.At.Ordinal != 0 {
		return f, g.ErrInvalid
	}
	out, e := s.Node(r.node, r.children...)
	if e != nil {
		return out, e
	}
	s.records[out.index-1].external = &use
	return out, nil
}

func Attributes(pairs ...string) []ast.Attribute {
	var out []ast.Attribute
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, ast.Attribute{Name: pairs[i], Text: pairs[i+1]})
	}
	return out
}
func Node[T any](kind ast.Kind[T], attrs []ast.Attribute, origin ast.NodeRef) ast.Node[T] {
	n := ast.Node[T]{Kind: kind, Attributes: attrs}
	if origin.Tree != "" {
		n.Origin.Nodes = []ast.NodeRef{origin}
	}
	return n
}

// Lowering owns native form meaning and grammar only. The compiler below owns
// dependency scheduling, scope, relocation, limits and publication association.
type Lowering[T any] interface {
	Parameter(string, ast.NodeRef) (Fragment[T], error)
	Reference(Symbol[T], []Fragment[T], ast.NodeRef) (Fragment[T], error)
	Expression(g.Expression[T], []Fragment[T]) (Fragment[T], error)
	Declaration(g.NativeDeclaration[T], []Fragment[T], Fragment[T]) (Fragment[T], error)
	Finish([]Fragment[T]) (Fragment[T], error)
}
