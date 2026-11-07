// Package target constructs rust syntax. It contains no canonical type mappings.
package target

import (
	"context"

	check "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Config struct{ TreeID ast.TreeID }
type Backend[F any] struct {
	*m.Store[F, s.Language]
	config Config
}

func New[F any](ctx context.Context, to representation.Language, config Config, limits g.Limits) (*Backend[F], error) {
	if to.ID() != (s.Language{}).ID() {
		return nil, g.ErrInvalid
	}
	store, err := m.New[F, s.Language](ctx, to, config.TreeID, limits)
	if err != nil {
		return nil, err
	}
	b := &Backend[F]{Store: store, config: config}
	store.Check = check.Check
	store.Finish = b.assemble
	return b, nil
}
func (b *Backend[F]) Parameter(name string, origin ast.NodeRef) (g.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin))
}
func (b *Backend[F]) Reference(symbol m.Symbol[s.Language], args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	segment, err := b.Node(m.Node(s.KindPathSegment, m.Attributes("name", symbol.Name()), origin), edges(s.RoleArgument, args)...)
	if err != nil {
		return segment, err
	}
	segment, err = b.Link(segment, symbol, s.RoleTarget)
	if err != nil {
		return segment, err
	}
	return b.Node(m.Node(s.KindPathType, nil, origin), m.Child[s.Language]{Role: s.RoleSegment, Fragment: segment})
}
func edges(role string, fragments []g.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}
