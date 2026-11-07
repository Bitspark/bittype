// Package target constructs cpp syntax. It contains no canonical type mappings.
package target

import (
	"context"
	check "github.com/Bitspark/bittype/check/syntax/cpp/provider/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	native "github.com/Bitspark/bittype/internal/cpp"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

type Config struct {
	TreeID    ast.TreeID
	Namespace string
}
type Backend[F any] struct {
	*m.Store[F, s.Language]
	config  Config
	headers []string
}

func New[F any](ctx context.Context, to representation.Language, config Config, limits g.Limits) (*Backend[F], error) {
	if to.ID() != (s.Language{}).ID() {
		return nil, g.ErrInvalid
	}
	if config.Namespace == "" || !native.Identifier(config.Namespace) || strings.HasPrefix(config.Namespace, "_") || config.Namespace == "std" {
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
	var segments []m.Child[s.Language]
	var attrs []ast.Attribute
	if !symbol.IsFormal() {
		ns, err := b.Node(m.Node(s.KindNameSegment, m.Attributes("name", b.config.Namespace), origin))
		if err != nil {
			return ns, err
		}
		segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: ns})
		attrs = m.Attributes("global", "true")
	}
	segment, err := b.Node(m.Node(s.KindNameSegment, m.Attributes("name", symbol.Name()), origin), edges(s.RoleArgument, args)...)
	if err != nil {
		return segment, err
	}
	segment, err = b.Link(segment, symbol, s.RoleTarget)
	if err != nil {
		return segment, err
	}
	segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: segment})
	return b.Node(m.Node(s.KindNamedType, attrs, origin), segments...)
}
func edges(role string, fragments []g.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}
func attribute(n s.Node, key string) string {
	for _, a := range n.Attributes {
		if a.Name == key {
			return a.Text
		}
	}
	return ""
}
func (b *Backend[F]) Include(header string) error {
	if header == "" || strings.ContainsAny(header, "\r\n<>") {
		return g.ErrInvalid
	}
	for _, old := range b.headers {
		if old == header {
			return nil
		}
	}
	b.headers = append(b.headers, header)
	return nil
}
