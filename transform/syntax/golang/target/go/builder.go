// Package target constructs golang syntax. It contains no canonical type mappings.
package target

import (
	"context"
	check "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"go/token"
)

type Config struct {
	TreeID  ast.TreeID
	Package string
}
type Backend[F any] struct {
	*m.Store[F, s.Language]
	config Config
}

func New[F any](ctx context.Context, to representation.Language, config Config, limits g.Limits) (*Backend[F], error) {
	if to.ID() != (s.Language{}).ID() {
		return nil, g.ErrInvalid
	}
	if config.Package == "" || token.Lookup(config.Package).IsKeyword() {
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
	constraint, err := b.Node(m.Node(s.KindInterfaceType, nil, origin))
	if err != nil {
		return constraint, err
	}
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin), m.Child[s.Language]{Role: s.RoleConstraint, Fragment: constraint})
}
func (b *Backend[F]) Reference(symbol m.Symbol[s.Language], args []g.Fragment[s.Language], origin ast.NodeRef) (g.Fragment[s.Language], error) {
	node, err := b.Node(m.Node(s.KindNamedType, m.Attributes("name", symbol.Name()), origin))
	if err != nil {
		return node, err
	}
	node, err = b.Link(node, symbol, s.RoleTarget)
	if err != nil {
		return node, err
	}
	if len(args) == 0 {
		return node, nil
	}
	children := []m.Child[s.Language]{{Role: s.RoleConstructor, Fragment: node}}
	for _, a := range args {
		children = append(children, m.Child[s.Language]{Role: s.RoleArgument, Fragment: a})
	}
	return b.Node(m.Node(s.KindTypeApplication, nil, origin), children...)
}
func edges(role string, fragments []g.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}
