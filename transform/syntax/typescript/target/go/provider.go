// Package target constructs native ASTs from a closed native semantic plan.
package target

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	lit "github.com/Bitspark/bittype/internal/typescript"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	model "github.com/Bitspark/bittype/transform/syntax/typescript/model/go"
	"github.com/Bitspark/bittype/transform/syntax/typescript/symbol/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Provider struct{}
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

	b := &Backend{Store: m.New[s.Language](ctx, in.Limits), config: p, language: in.Language}
	return m.Compile(ctx, in, b.Store, b)
}
func (b *Backend) Parameter(name string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin))
}
func (b *Backend) Reference(symbol m.Symbol[s.Language], args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	node, err := b.Node(m.Node(s.KindTypeReference, m.Attributes("name", symbol.Name()), origin))
	if err != nil {
		return node, err
	}
	node, err = b.Link(node, symbol, s.RoleTarget)
	if err != nil {
		return node, err
	}
	return b.apply(node, args, origin)
}

func (b *Backend) apply(node m.Fragment[s.Language], args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	if len(args) == 0 {
		return node, nil
	}
	children := append([]m.Child[s.Language]{{Role: s.RoleConstructor, Fragment: node}}, edges(s.RoleArgument, args)...)
	return b.Node(m.Node(s.KindTypeApplication, nil, origin), children...)
}
func edges(role string, fragments []m.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}

// PackageType refers to a named export of a selected package.
func (b *Backend) PackageType(module, export string, args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	if !lit.PackageSpecifier(module) || !lit.TypeName(export) {
		return m.Fragment[s.Language]{}, g.ErrInvalid
	}
	return b.externalType(module, export, args, origin)
}

// GlobalType identifies a public type supplied by the target library environment.
func (b *Backend) GlobalType(name string, args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	if !lit.TypeName(name) {
		return m.Fragment[s.Language]{}, g.ErrInvalid
	}
	return b.externalType("", name, args, origin)
}
func (b *Backend) externalType(module, export string, args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	ref, err := b.Node(m.Node(s.KindTypeReference, m.Attributes("name", export), origin))
	if err != nil {
		return ref, err
	}
	external, err := symbol.ExternalUse(b.language, ast.NodeRef{}, symbol.Type{Module: module, Name: export})
	if err != nil {
		return ref, err
	}
	ref, err = b.External(ref, external)
	if err != nil {
		return ref, err
	}
	return b.apply(ref, args, origin)
}
func (b *Backend) Declaration(d g.NativeDeclaration[s.Language], params []m.Fragment[s.Language], body m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	name, origin := d.Name, d.Origin
	form, _ := model.Name(d.Form)
	record := form != "alias"

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
		return b.Node(m.Node(s.KindInterface, m.Attributes("name", name), origin), children...)
	}
	children = append(children, m.Child[s.Language]{Role: s.RoleType, Fragment: body})
	return b.Node(m.Node(s.KindTypeAlias, m.Attributes("name", name), origin), children...)
}
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindSourceFile, nil, ast.NodeRef{}), edges(s.RoleDeclaration, declarations)...)
}
func (b *Backend) Expression(e g.Expression[s.Language], args []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if err := model.CheckExpression(context.Background(), e); err != nil {
		return m.Fragment[s.Language]{}, err
	}
	name, _ := model.Name(e.Form())
	data, origin := e.Data(), e.Origin()
	switch name {
	case "primitive":
		spelling := "boolean"
		if data[0] == "string" {
			spelling = "string"
		}
		return b.Node(m.Node(s.KindKeywordType, m.Attributes("name", spelling), origin))
	case "array":
		return b.Node(m.Node(s.KindArrayType, nil, origin), m.Child[s.Language]{Role: s.RoleElement, Fragment: args[0]})
	case "field":
		return b.Node(m.Node(s.KindProperty, m.Attributes("name", data[0], "optional", data[1]), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case "record":
		return b.Node(m.Node(s.KindObjectType, nil, origin), edges(s.RoleMember, args)...)
	case "external":
		if data[0] == "" {
			return b.GlobalType(data[1], args, origin)
		}
		return b.PackageType(data[0], data[1], args, origin)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

func (p Provider) CheckProgram(ctx context.Context, program g.Program[s.Language]) error {
	return m.CheckSingleScope(ctx, program)
}
