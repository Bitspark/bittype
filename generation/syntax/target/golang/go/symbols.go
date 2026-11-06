package target

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	m "github.com/Bitspark/bittype/generation/syntax/target/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Reserve fixes the native declaration form and a target-owned header builder.
// record means this language's interface/struct declaration, not a Horos kind.
func (b *Backend[F]) Reserve(site g.Site[F, s.Language], name string, record bool, origin ast.NodeRef) (g.Destination[s.Language], bool, error) {
	return b.Store.Reserve(site, name, origin, func(params []g.Fragment[s.Language], body g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
		kind := s.KindTypeAlias
		if record {
			kind = s.KindTypeDeclaration
		}
		children := append(edges(s.RoleTypeParameter, params), m.Child[s.Language]{Role: s.RoleType, Fragment: body})
		return b.Node(m.Node(kind, m.Attributes("name", name), origin), children...)
	})
}
func (b *Backend[F]) Header(site g.Site[F, s.Language], d g.Destination[s.Language], names []string, origins []ast.NodeRef) (g.Layout[s.Language], error) {
	if len(names) != len(origins) {
		return g.Layout[s.Language]{}, g.ErrInvalid
	}
	var params []g.Fragment[s.Language]
	for i, name := range names {
		p, err := b.Parameter(name, origins[i])
		if err != nil {
			return g.Layout[s.Language]{}, err
		}
		params = append(params, p)
	}
	return b.Freeze(site, d, names, params)
}
