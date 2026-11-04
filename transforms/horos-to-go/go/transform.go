// Package golang maps analyzed Horos pagination types to Go syntax.
package golang

import (
	"context"
	"github.com/Bitspark/bittype/contracts/transform/go"
	g "github.com/Bitspark/bittype/internal/horosgen"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"go/token"
)

type Config struct {
	TreeID  ast.TreeID
	Package string
}

func New(from, to representation.Language, config Config) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	if from.Profile() != h.Profile || config.TreeID == "" || !g.Identifier(config.Package) || token.Lookup(config.Package).IsKeyword() {
		return nil, transform.ErrInvalidTransformer
	}
	return transform.NewTransformer[h.Language, s.Language](from, to, g.Contract(), func(ctx context.Context, request transform.Request[h.Language]) (s.AST, error) {
		return g.Generate(ctx, request, config.TreeID, &dialect{pkg: config.Package})
	})
}

type dialect struct {
	pkg      string
	optional ast.NodeID
}

func (*dialect) ReservedNames() []string { return []string{"BittypeOptional"} }

func (*dialect) FieldName(name string) (string, error) { return g.Exported(name) }
func (*dialect) Parameter(a *g.Arena[s.Language], name string, o ast.NodeRef) ast.NodeID {
	constraint := a.Add(s.KindInterfaceType, nil, o)
	return a.Add(s.KindTypeParameter, g.Attr("name", name), o, g.Child(s.RoleConstraint, constraint))
}
func (*dialect) Reference(a *g.Arena[s.Language], name string, target ast.NodeID, args []ast.NodeID, o ast.NodeRef) ast.NodeID {
	ref := a.Add(s.KindNamedType, g.Attr("name", name), o)
	a.Target(ref, target)
	if len(args) == 0 {
		return ref
	}
	return a.Add(s.KindTypeApplication, nil, o, append([]ast.Child{g.Child(s.RoleConstructor, ref)}, g.Children(s.RoleArgument, args)...)...)
}
func (*dialect) Text(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindNamedType, g.Attr("name", "string"), o)
}
func (*dialect) Boolean(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindNamedType, g.Attr("name", "bool"), o)
}
func (*dialect) List(a *g.Arena[s.Language], element ast.NodeID, o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindSliceType, nil, o, g.Child(s.RoleElement, element))
}
func (d *dialect) presence(a *g.Arena[s.Language]) ast.NodeID {
	if d.optional != 0 {
		return d.optional
	}
	d.optional = a.Add(s.KindTypeDeclaration, nil, ast.NodeRef{})
	p := d.Parameter(a, "T", ast.NodeRef{})
	value := d.Reference(a, "T", p, nil, ast.NodeRef{})
	boolean := d.Boolean(a, ast.NodeRef{})
	d.Record(a, d.optional, "BittypeOptional", []ast.NodeID{p}, []g.Field{{Name: "Present", Type: boolean}, {Name: "Value", Type: value}})
	return d.optional
}
func (d *dialect) Record(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, fields []g.Field) {
	var members []ast.Child
	for _, field := range fields {
		ty := field.Type
		if field.Optional {
			ty = d.Reference(a, "BittypeOptional", d.presence(a), []ast.NodeID{ty}, field.Origin)
		}
		member := a.Add(s.KindField, g.Attr("name", field.Name), field.Origin, g.Child(s.RoleType, ty))
		members = append(members, g.Child(s.RoleField, member))
	}
	body := a.Add(s.KindStructType, nil, a.Origin(id), members...)
	a.Set(id, s.KindTypeDeclaration, g.Attr("name", name), append(g.Children(s.RoleTypeParameter, params), g.Child(s.RoleType, body))...)
}
func (*dialect) Alias(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, body ast.NodeID) {
	a.Set(id, s.KindTypeAlias, g.Attr("name", name), append(g.Children(s.RoleTypeParameter, params), g.Child(s.RoleType, body))...)
}
func (d *dialect) Finish(a *g.Arena[s.Language], decls []ast.NodeID) ast.NodeID {
	if d.optional != 0 {
		decls = append([]ast.NodeID{d.optional}, decls...)
	}
	return a.Add(s.KindFile, g.Attr("package", d.pkg), ast.NodeRef{}, g.Children(s.RoleDeclaration, decls)...)
}
