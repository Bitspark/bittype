// Package typescript maps analyzed Horos pagination types to TypeScript syntax.
package typescript

import (
	"context"
	"github.com/Bitspark/bittype/contracts/transform/go"
	g "github.com/Bitspark/bittype/internal/horosgen"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strconv"
)

// RepresentationProfile permits kind-constraint erasure for Horos-checked uses.
const RepresentationProfile = g.RepresentationProfile

// EnforcementProfile requires native kind enforcement and is not implemented.
const EnforcementProfile = g.EnforcementProfile

type Config struct {
	TreeID  ast.TreeID
	Profile string
}

func New(from, to representation.Language, config Config) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	if from.Profile() != h.Profile || config.TreeID == "" || !g.KnownProfile(config.Profile) {
		return nil, transform.ErrInvalidTransformer
	}
	return transform.NewTransformer[h.Language, s.Language](from, to, g.Contract(), func(ctx context.Context, request transform.Request[h.Language]) (ast.Set[s.Language], error) {
		output, err := g.Generate(ctx, request, config.TreeID, config.Profile, dialect{})
		if err != nil {
			return ast.Set[s.Language]{}, err
		}
		return ast.Set[s.Language]{Entry: ast.NodeRef{Tree: output.ID, Node: output.Root}, Trees: []ast.AST[s.Language]{output}}, nil
	})
}

type dialect struct{}

func (dialect) ReservedNames() []string { return nil }

func (dialect) FieldName(name string) (string, error) { return name, nil }
func (dialect) Parameter(a *g.Arena[s.Language], name string, o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindTypeParameter, g.Attr("name", name), o)
}
func (dialect) Reference(a *g.Arena[s.Language], name string, target ast.NodeID, args []ast.NodeID, o ast.NodeRef) ast.NodeID {
	ref := a.Add(s.KindTypeReference, g.Attr("name", name), o)
	a.Target(ref, target)
	if len(args) == 0 {
		return ref
	}
	children := append([]ast.Child{g.Child(s.RoleConstructor, ref)}, g.Children(s.RoleArgument, args)...)
	return a.Add(s.KindTypeApplication, nil, o, children...)
}
func (dialect) Text(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindKeywordType, g.Attr("name", "string"), o)
}
func (dialect) Boolean(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindKeywordType, g.Attr("name", "boolean"), o)
}
func (dialect) List(a *g.Arena[s.Language], element ast.NodeID, o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindArrayType, nil, o, g.Child(s.RoleElement, element))
}
func (dialect) Record(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, fields []g.Field) {
	children := g.Children(s.RoleTypeParameter, params)
	for _, field := range fields {
		member := a.Add(s.KindProperty, g.Attr("name", field.Name, "optional", strconv.FormatBool(field.Optional)), field.Origin, g.Child(s.RoleType, field.Type))
		children = append(children, g.Child(s.RoleMember, member))
	}
	a.Set(id, s.KindInterface, g.Attr("name", name, "exported", "true"), children...)
}
func (dialect) Alias(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, body ast.NodeID) {
	a.Set(id, s.KindTypeAlias, g.Attr("name", name, "exported", "true"), append(g.Children(s.RoleTypeParameter, params), g.Child(s.RoleType, body))...)
}
func (dialect) Finish(a *g.Arena[s.Language], decls []ast.NodeID) ast.NodeID {
	return a.Add(s.KindSourceFile, nil, ast.NodeRef{}, g.Children(s.RoleDeclaration, decls)...)
}
