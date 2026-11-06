// Package target owns shared native-template storage, not a native grammar or
// a canonical representation catalogue. Language adapters supply syntax checks,
// declaration headers and output-unit construction.
package target

import (
	"context"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"slices"
)

type Child[T any] struct {
	Role     string
	Fragment g.Fragment[T]
}
type Symbol[T any] struct {
	destination g.Destination[T]
	layout      g.Layout[T]
	index       int
	name        string
}

func (s Symbol[T]) Name() string   { return s.name }
func (s Symbol[T]) IsFormal() bool { return s.layout.Valid() }

type fragment[T any] struct {
	node          ast.Node[T]
	children      []Child[T]
	symbol        Symbol[T]
	referenceRole string
}
type identity struct {
	run       g.RunID
	config    g.ConfigID
	frame     g.FrameID
	signature g.SignatureID
}

func at[F, T any](site g.Site[F, T]) identity {
	return identity{site.Run(), site.Config(), site.Frame(), site.Signature()}
}

type destination[T any] struct {
	identity identity
	name     string
	origin   ast.NodeRef
	build    func([]g.Fragment[T], g.Fragment[T]) (g.Fragment[T], error)
}
type layout[T any] struct {
	identity    identity
	destination g.Destination[T]
	names       []string
	parameters  []g.Fragment[T]
}

type Store[F, T any] struct {
	ctx                context.Context
	language           representation.Language
	id                 ast.TreeID
	limits             g.Limits
	destinations       g.Space[g.DestinationRole[T]]
	layouts            g.Space[g.LayoutRole[T]]
	fragments          g.Space[g.FragmentRole[T]]
	destinationRecords map[g.Destination[T]]destination[T]
	layoutRecords      map[g.Layout[T]]layout[T]
	fragmentRecords    map[g.Fragment[T]]fragment[T]
	headers            map[g.Destination[T]]g.Layout[T]
	occupied           map[string]bool
	formalNames        map[string]bool
	bytes              int
	Check              func(ast.View[T]) error
	Finish             func([]g.Fragment[T]) (g.Fragment[T], error)
}

func New[F, T any](ctx context.Context, language representation.Language, id ast.TreeID, limits g.Limits) (*Store[F, T], error) {
	if ctx == nil || id == "" || !limits.Valid() {
		return nil, g.ErrInvalid
	}
	if err := representation.CheckMarker[T](language); err != nil {
		return nil, err
	}
	return &Store[F, T]{ctx: ctx, language: language, id: id, limits: limits, destinations: g.NewSpace[g.DestinationRole[T]](), layouts: g.NewSpace[g.LayoutRole[T]](), fragments: g.NewSpace[g.FragmentRole[T]](), destinationRecords: map[g.Destination[T]]destination[T]{}, layoutRecords: map[g.Layout[T]]layout[T]{}, fragmentRecords: map[g.Fragment[T]]fragment[T]{}, headers: map[g.Destination[T]]g.Layout[T]{}, occupied: map[string]bool{}, formalNames: map[string]bool{}}, nil
}

// Node captures native syntax atoms and symbolic child relocations. Callers
// cannot supply live arena children/references or mutate retained storage.
func (s *Store[F, T]) Node(node ast.Node[T], children ...Child[T]) (g.Fragment[T], error) {
	if err := s.ctx.Err(); err != nil {
		return g.Fragment[T]{}, err
	}
	if len(node.Children) != 0 || len(node.References) != 0 || node.Kind == "" {
		return g.Fragment[T]{}, g.ErrInvalid
	}
	if len(s.fragmentRecords) >= s.limits.OutputNodes {
		return g.Fragment[T]{}, g.ErrLimit
	}
	bytes := len(node.Kind)
	for _, a := range node.Attributes {
		bytes += len(a.Name) + len(a.Text)
	}
	for _, c := range children {
		bytes += len(c.Role)
		if _, ok := s.fragmentRecords[c.Fragment]; !ok {
			return g.Fragment[T]{}, g.ErrScope
		}
	}
	if bytes > s.limits.OutputBytes-s.bytes {
		return g.Fragment[T]{}, g.ErrLimit
	}
	s.bytes += bytes
	node.Attributes = slices.Clone(node.Attributes)
	node.Origin.Nodes = slices.Clone(node.Origin.Nodes)
	node.Origin.Spans = slices.Clone(node.Origin.Spans)
	f, _ := s.fragments.Ref(uint64(len(s.fragmentRecords) + 1))
	s.fragmentRecords[f] = fragment[T]{node: node, children: slices.Clone(children)}
	return f, nil
}

func (s *Store[F, T]) Link(f g.Fragment[T], symbol Symbol[T], role string) (g.Fragment[T], error) {
	r, ok := s.fragmentRecords[f]
	if !ok || role == "" || r.referenceRole != "" {
		return g.Fragment[T]{}, g.ErrScope
	}
	if _, err := s.symbol(symbol); err != nil {
		return g.Fragment[T]{}, err
	}
	copy, err := s.Node(r.node, r.children...)
	if err != nil {
		return copy, err
	}
	r = s.fragmentRecords[copy]
	r.symbol = symbol
	r.referenceRole = role
	s.fragmentRecords[copy] = r
	return copy, nil
}
func (s *Store[F, T]) Inspect(f g.Fragment[T]) (ast.Node[T], []Child[T], error) {
	r, ok := s.fragmentRecords[f]
	if !ok {
		return ast.Node[T]{}, nil, g.ErrScope
	}
	n := r.node
	n.Attributes = slices.Clone(n.Attributes)
	n.Origin.Nodes = slices.Clone(n.Origin.Nodes)
	n.Origin.Spans = slices.Clone(n.Origin.Spans)
	return n, slices.Clone(r.children), nil
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
