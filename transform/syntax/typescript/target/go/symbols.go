package target

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Reserve fixes the native declaration form and a target-owned header builder.
// record means this language's interface/struct declaration, not a Horos kind.
func (b *Backend[F]) Reserve(site g.Site[F, s.Language], name string, record bool, origin ast.NodeRef) (g.Destination[s.Language], bool, error) {
	return b.Store.Reserve(site, name, origin, func(params []g.Fragment[s.Language], body g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
		children := edges(s.RoleTypeParameter, params)
		if record {
			n, fields, err := b.Inspect(body)
			if err != nil {
				return body, err
			}
			if n.Kind != s.KindObjectType {
				return body, g.ErrInvalid
			}
			children = append(children, fields...)
			return b.Node(m.Node(s.KindInterface, m.Attributes("name", name, "exported", "true"), origin), children...)
		}
		children = append(children, m.Child[s.Language]{Role: s.RoleType, Fragment: body})
		return b.Node(m.Node(s.KindTypeAlias, m.Attributes("name", name, "exported", "true"), origin), children...)
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
