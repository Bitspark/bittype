// Package target constructs native ASTs from a closed native semantic plan.
package target

import (
	"context"
	"strings"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type Provider struct{}
type Backend struct {
	*m.Store[s.Language]
	config Provider
}

func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	return model.CheckExpression(ctx, e)
}
func (p Provider) CheckDeclaration(ctx context.Context, d g.NativeDeclaration[s.Language]) error {
	return model.CheckDeclaration(ctx, d)
}
func (p Provider) Build(ctx context.Context, in g.BuildInput[s.Language]) (g.Assembly, error) {

	b := &Backend{Store: m.New[s.Language](ctx, in.Limits), config: p}
	return m.Compile(ctx, in, b.Store, b)
}
func (b *Backend) Parameter(name string, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindTypeParameter, m.Attributes("name", name), origin))
}
func (b *Backend) Reference(symbol m.Symbol[s.Language], args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
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
func edges(role string, fragments []m.Fragment[s.Language]) []m.Child[s.Language] {
	out := make([]m.Child[s.Language], len(fragments))
	for i, f := range fragments {
		out[i] = m.Child[s.Language]{Role: role, Fragment: f}
	}
	return out
}

// Path is native syntax; selecting a library artifact and its meaning is the caller's responsibility.
func (b *Backend) Path(names []string, args []m.Fragment[s.Language], origin ast.NodeRef, absolute bool) (m.Fragment[s.Language], error) {
	if len(names) == 0 {
		return m.Fragment[s.Language]{}, g.ErrInvalid
	}
	var segments []m.Child[s.Language]
	for i, name := range names {
		var children []m.Child[s.Language]
		if i == len(names)-1 {
			children = edges(s.RoleArgument, args)
		}
		p, err := b.Node(m.Node(s.KindPathSegment, m.Attributes("name", name), origin), children...)
		if err != nil {
			return p, err
		}
		segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: p})
	}
	var attrs []ast.Attribute
	if absolute {
		attrs = m.Attributes("absolute", "true")
	}
	return b.Node(m.Node(s.KindPathType, attrs, origin), segments...)
}
func (b *Backend) Declaration(d g.NativeDeclaration[s.Language], params []m.Fragment[s.Language], body m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	name, origin := d.Name, d.Origin
	form, _ := model.Name(d.Form)
	record := form != "alias"

	children := edges(s.RoleGenericParameter, params)
	if record {
		n, fields, err := b.Inspect(body)
		if err != nil {
			return body, err
		}
		if n.Kind != s.KindStruct {
			return body, g.ErrInvalid
		}
		children = append(children, fields...)
		return b.Node(m.Node(s.KindStruct, m.Attributes("name", name, "visibility", "public"), origin), children...)
	}
	children = append(children, m.Child[s.Language]{Role: s.RoleType, Fragment: body})
	return b.Node(m.Node(s.KindTypeAlias, m.Attributes("name", name, "visibility", "public"), origin), children...)
}
func (b *Backend) Finish(declarations []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	return b.Node(m.Node(s.KindFile, nil, ast.NodeRef{}), edges(s.RoleItem, declarations)...)
}
func (b *Backend) Expression(e g.Expression[s.Language], args []m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	if err := model.CheckExpression(context.Background(), e); err != nil {
		return m.Fragment[s.Language]{}, err
	}
	name, _ := model.Name(e.Form())
	data, origin := e.Data(), e.Origin()
	switch name {
	case "primitive":
		if data[0] == "string" {
			return b.Path([]string{"std", "string", "String"}, nil, origin, true)
		}
		return b.Node(m.Node(s.KindPrimitiveType, m.Attributes("name", "bool"), origin))
	case "sequence":
		return b.Path([]string{"std", "vec", "Vec"}, args, origin, true)
	case "field":
		return b.Node(m.Node(s.KindField, m.Attributes("name", data[0], "visibility", "public"), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case "record":
		return b.Node(m.Node(s.KindStruct, nil, origin), edges(s.RoleField, args)...)
	case "external":
		return b.Path(append(strings.Split(data[0], "::"), data[1]), args, origin, true)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

func (p Provider) CheckProgram(ctx context.Context, program g.Program[s.Language]) error {
	return m.CheckSingleScope(ctx, program)
}
