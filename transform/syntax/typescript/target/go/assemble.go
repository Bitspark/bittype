package target

import (
	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func (b *Backend[F]) assemble(declarations []g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindSourceFile, nil, ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
}
