// Package target constructs native ASTs from a closed native semantic plan.
package target

import (
	"context"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	model "github.com/Bitspark/bittype/transform/syntax/rust/model/go"
	"github.com/Bitspark/bittype/transform/syntax/rust/symbol/go"
	m "github.com/Bitspark/bittype/transform/syntax/shared/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Provider struct{}
type Backend struct {
	*m.Store[s.Language]
	config   Provider
	language representation.Language
}

// CheckExpression admits the forms of the Rust model. An external form is
// admitted when its crate is a dependency that a destination can name.
func (p Provider) CheckExpression(ctx context.Context, e g.Expression[s.Language]) error {
	if err := model.CheckExpression(ctx, e); err != nil {
		return err
	}
	if name, _ := model.Name(e.Form()); name == model.FormExternal {
		_, err := external(e.Data())
		return err
	}
	return nil
}

// external returns the destination that the data of an external form names:
// the crate first, then the path of the item. The model's name of the
// standard library is mapped to its dependency. For every other text the
// symbol package says why no destination can name it.
func external(data []string) (symbol.Symbol, error) {
	crate := symbol.Std()
	if data[0] != model.Std {
		var err error
		if crate, err = symbol.Dependency(data[0]); err != nil {
			return symbol.Symbol{}, err
		}
	}
	return symbol.Item(crate, data[1:]...)
}
func (p Provider) CheckDeclaration(ctx context.Context, d g.TargetDeclaration[s.Language]) error {
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

// Item writes a use of an item of a dependency as an absolute path: the name
// of the crate, then the path of the item. The last segment names the item,
// takes the arguments and records the destination. The segments before it
// qualify the use.
func (b *Backend) Item(destination symbol.Symbol, args []m.Fragment[s.Language], origin ast.NodeRef) (m.Fragment[s.Language], error) {
	crate, path, ok := destination.Item()
	if !ok {
		return m.Fragment[s.Language]{}, g.ErrInvalid
	}
	use, err := symbol.ExternalUse(b.language, ast.NodeRef{}, destination)
	if err != nil {
		return m.Fragment[s.Language]{}, err
	}
	names := append([]string{crate.Name()}, path...)
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
		if i == len(names)-1 {
			if p, err = b.External(p, use); err != nil {
				return p, err
			}
		}
		segments = append(segments, m.Child[s.Language]{Role: s.RoleSegment, Fragment: p})
	}
	return b.Node(m.Node(s.KindPathType, m.Attributes("absolute", "true"), origin), segments...)
}

// primitive writes a use of a primitive type. Its name is no keyword, so the
// use records its destination.
func (b *Backend) primitive(destination symbol.Symbol, origin ast.NodeRef) (m.Fragment[s.Language], error) {
	name, ok := destination.Primitive()
	if !ok {
		return m.Fragment[s.Language]{}, g.ErrInvalid
	}
	node, err := b.Node(m.Node(s.KindPrimitiveType, m.Attributes("name", name), origin))
	if err != nil {
		return node, err
	}
	use, err := symbol.ExternalUse(b.language, ast.NodeRef{}, destination)
	if err != nil {
		return node, err
	}
	return b.External(node, use)
}
func (b *Backend) Declaration(d g.TargetDeclaration[s.Language], params []m.Fragment[s.Language], body m.Fragment[s.Language]) (m.Fragment[s.Language], error) {
	name, origin := d.Name, d.Origin
	form, _ := model.Name(d.Form)
	record := form != model.FormAlias

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
	case model.FormPrimitive:
		if data[0] == model.PrimitiveString {
			return b.Item(symbol.StdString(), nil, origin)
		}
		return b.primitive(symbol.Bool(), origin)
	case model.FormSequence:
		return b.Item(symbol.StdVec(), args, origin)
	case model.FormField:
		return b.Node(m.Node(s.KindField, m.Attributes("name", data[0], "visibility", "public"), origin), m.Child[s.Language]{Role: s.RoleType, Fragment: args[0]})
	case model.FormRecord:
		return b.Node(m.Node(s.KindStruct, nil, origin), edges(s.RoleField, args)...)
	case model.FormExternal:
		destination, err := external(data)
		if err != nil {
			return m.Fragment[s.Language]{}, err
		}
		return b.Item(destination, args, origin)
	}
	return m.Fragment[s.Language]{}, g.ErrUnsupported
}

func (p Provider) CheckProgram(ctx context.Context, program g.Program[s.Language]) error {
	return m.CheckSingleScope(ctx, program)
}

// CheckMember states the member vocabulary of the target (bittype#196, row
// 13): generation.MemberVocabulary.
func (p Provider) CheckMember(ctx context.Context, e g.Expression[s.Language]) error {
	return model.CheckMember(ctx, e)
}
