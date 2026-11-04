// Package rust maps analyzed Horos pagination types to Rust syntax.
package rust

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/transform/go"
	g "github.com/Bitspark/bittype/internal/horosgen"
	native "github.com/Bitspark/bittype/internal/rust"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

type Config struct{ TreeID ast.TreeID }

func New(from, to representation.Language, config Config) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	if from.Profile() != h.Profile || config.TreeID == "" {
		return nil, transform.ErrInvalidTransformer
	}
	return transform.NewTransformer[h.Language, s.Language](from, to, g.Contract(), func(ctx context.Context, request transform.Request[h.Language]) (s.AST, error) {
		return g.Generate(ctx, request, config.TreeID, dialect{})
	})
}

type dialect struct{}

func (dialect) ReservedNames() []string { return nil }

func (dialect) FieldName(name string) (string, error) {
	if !g.Identifier(name) || !native.Declaration(name) {
		return "", fmt.Errorf("unsupported Rust field name %q", name)
	}
	return name, nil
}
func (dialect) Parameter(a *g.Arena[s.Language], name string, o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindTypeParameter, g.Attr("name", name), o)
}
func path(a *g.Arena[s.Language], names []string, args []ast.NodeID, target ast.NodeID, o ast.NodeRef, absolute bool) ast.NodeID {
	var segments []ast.Child
	for i, name := range names {
		var children []ast.Child
		if i == len(names)-1 {
			children = g.Children(s.RoleArgument, args)
		}
		segment := a.Add(s.KindPathSegment, g.Attr("name", name), o, children...)
		if i == len(names)-1 && target != 0 {
			a.Target(segment, target)
		}
		segments = append(segments, g.Child(s.RoleSegment, segment))
	}
	var attrs []ast.Attribute
	if absolute {
		attrs = g.Attr("absolute", "true")
	}
	return a.Add(s.KindPathType, attrs, o, segments...)
}
func (dialect) Reference(a *g.Arena[s.Language], name string, target ast.NodeID, args []ast.NodeID, o ast.NodeRef) ast.NodeID {
	return path(a, []string{name}, args, target, o, false)
}
func (dialect) Text(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return path(a, []string{"std", "string", "String"}, nil, 0, o, true)
}
func (dialect) Boolean(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindPrimitiveType, g.Attr("name", "bool"), o)
}
func (dialect) List(a *g.Arena[s.Language], element ast.NodeID, o ast.NodeRef) ast.NodeID {
	return path(a, []string{"std", "vec", "Vec"}, []ast.NodeID{element}, 0, o, true)
}
func (dialect) Record(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, fields []g.Field) {
	children := g.Children(s.RoleGenericParameter, params)
	for _, field := range fields {
		ty := field.Type
		if field.Optional {
			ty = path(a, []string{"std", "option", "Option"}, []ast.NodeID{ty}, 0, field.Origin, true)
		}
		member := a.Add(s.KindField, g.Attr("name", field.Name, "visibility", "public"), field.Origin, g.Child(s.RoleType, ty))
		children = append(children, g.Child(s.RoleField, member))
	}
	a.Set(id, s.KindStruct, g.Attr("name", name, "visibility", "public"), children...)
}
func (dialect) Alias(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, body ast.NodeID) {
	a.Set(id, s.KindTypeAlias, g.Attr("name", name, "visibility", "public"), append(g.Children(s.RoleGenericParameter, params), g.Child(s.RoleType, body))...)
}
func (dialect) Finish(a *g.Arena[s.Language], decls []ast.NodeID) ast.NodeID {
	return a.Add(s.KindFile, nil, ast.NodeRef{}, g.Children(s.RoleItem, decls)...)
}
