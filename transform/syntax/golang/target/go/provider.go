// Package target constructs native ASTs from a closed native semantic plan.
package target

import (
	"context"
	"go/token"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	model "github.com/Bitspark/bittype/transform/syntax/golang/model/go"
	"github.com/Bitspark/bittype/transform/syntax/golang/symbol/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Provider struct{ Package string }
type Backend struct {
	*m.Store[s.Language]
	config   Provider
	language representation.Language
}

func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	return model.CheckExpression(ctx, e)
}
func (p Provider) CheckDeclaration(ctx context.Context, d g.NativeDeclaration[s.Language]) error {
	return model.CheckDeclaration(ctx, d)
}
func (p Provider) Build(ctx context.Context, in g.BuildInput[s.Language]) (g.Assembly, error) {
	if p.Package == "" || token.Lookup(p.Package).IsKeyword() {
		return g.Assembly{}, g.ErrInvalid
	}
	b := p.Backend(ctx, in)
	return m.Compile(ctx, in, b.Store, b)
}

// Backend returns the lowering of one build: the forms of the Go model over a
// fresh store, with the configuration of p and the language of the build. A
// use of a predeclared type records a symbol of that language, so a lowering
// that is constructed in another way cannot write one.
func (p Provider) Backend(ctx context.Context, in g.BuildInput[s.Language]) *Backend {
	return &Backend{Store: m.New[s.Language](ctx, in.Limits), config: p, language: in.Language}
}
func (b *Backend) Parameter(name string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	constraint, err := b.Node(m.Node(s.KindInterfaceType, nil, origin))
	if err != nil {
		return constraint, err
	}
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin), m.Child[s.Language]{Role: s.RoleConstraint, Fragment: constraint})
}
func (b *Backend) Reference(symbol m.Symbol[s.Language], args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
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
func edges(role string, fragments []m.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}

// predeclared writes a use of a type of the universe scope. The identifier is
// a name that a declaration or a type parameter can capture, so the use
// records its destination: the predeclared symbol.
func (b *Backend) predeclared(name string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	destination, err := symbol.Predeclared(name)
	if err != nil {
		return m.Fragment[s.Language]{}, err
	}
	node, err := b.Node(m.Node(s.KindNamedType, m.Attributes("name", name), origin))
	if err != nil {
		return node, err
	}
	use, err := symbol.ExternalUse(b.language, ast.NodeRef{}, destination)
	if err != nil {
		return node, err
	}
	return b.External(node, use)
}
func (b *Backend) Declaration(d g.NativeDeclaration[s.Language], params []m.Fragment[s.Language], body m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	name, origin := d.Name, d.Origin
	form, _ := model.Name(d.Form)
	record := form != model.FormAlias

	kind := s.KindTypeAlias
	if record {
		kind = s.KindTypeDeclaration
	}
	children := append(edges(s.RoleTypeParameter, params), m.Child[s.Language]{Role: s.RoleType, Fragment: body})
	return b.Node(m.Node(kind, m.Attributes("name", name), origin), children...)
}
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindFile, m.Attributes("package", b.config.Package), ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
}
func (b *Backend) Expression(e g.Expression[s.Language], args []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if err := model.CheckExpression(context.Background(), e); err != nil {
		return m.Fragment[s.Language]{}, err
	}
	name, _ := model.Name(e.Form())
	data, origin := e.Data(), e.Origin()
	switch name {
	case model.FormPrimitive:
		return b.predeclared(data[0], origin)
	case model.FormSlice:
		return b.Node(m.Node(s.KindSliceType, nil, origin), m.Child[s.Language]{Role: s.RoleElement, Fragment: args[0]})
	case model.FormField:
		return b.Node(m.Node(s.KindField, m.Attributes("name", data[0]), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case model.FormRecord:
		return b.Node(m.Node(s.KindStructType, nil, origin), edges(s.RoleField, args)...)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

func (p Provider) CheckProgram(ctx context.Context, program g.Program[s.Language]) error {
	return m.CheckSingleScope(ctx, program)
}
