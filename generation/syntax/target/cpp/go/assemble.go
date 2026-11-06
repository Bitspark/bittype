package target

import (
 "slices"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	m "github.com/Bitspark/bittype/generation/syntax/target/go"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func (b *Backend[F]) assemble(declarations []g.Fragment[s.Language]) (g.Fragment[s.Language], error) {
	var children []m.Child[s.Language]
	headers := slices.Clone(b.headers)
 slices.Sort(headers)
 for _, header := range headers {
		n, err := b.Node(m.Node(s.KindInclude, m.Attributes("path", header, "system", "true"), ast.NodeRef{}))
		if err != nil {
			return n, err
		}
		children = append(children, m.Child[s.Language]{Role: s.RoleInclude, Fragment: n})
	}
	ns, err := b.Node(m.Node(s.KindNamespace, m.Attributes("name", b.config.Namespace), ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
	if err != nil {
		return ns, err
	}
	children = append(children, m.Child[s.Language]{Role: s.RoleDeclaration, Fragment: ns})
	return b.Node(m.Node(s.KindTranslationUnit, nil, ast.NodeRef{}), children...)
}
