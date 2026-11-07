package target

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func (b *Backend[F]) assemble(declarations []g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindFile, nil, ast.NodeRef{}), edges(s.RoleItem, declarations)...)
}
