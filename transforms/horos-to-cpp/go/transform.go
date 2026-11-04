// Package cpp maps analyzed Horos pagination types to C++ syntax.
package cpp

import (
	"context"
	"fmt"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	native "github.com/Bitspark/bittype/internal/cpp"
	g "github.com/Bitspark/bittype/internal/horosgen"
	s "github.com/Bitspark/bittype/syntax/cpp/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
)

type Config struct {
	TreeID    ast.TreeID
	Namespace string
}

func New(from, to representation.Language, config Config) (transform.ASTTransformer, error) {
	if from.ID() != (h.Language{}).ID() || to.ID() != (s.Language{}).ID() {
		return nil, transform.ErrInvalidTransformer
	}
	if from.Profile() != h.Profile || config.TreeID == "" || !g.Identifier(config.Namespace) || !native.Identifier(config.Namespace) || strings.HasPrefix(config.Namespace, "_") || config.Namespace == "std" {
		return nil, transform.ErrInvalidTransformer
	}
	return transform.NewTransformer[h.Language, s.Language](from, to, g.Contract(), func(ctx context.Context, request transform.Request[h.Language]) (ast.Set[s.Language], error) {
		output, err := g.Generate(ctx, request, config.TreeID, dialect{namespace: config.Namespace})
		if err != nil {
			return ast.Set[s.Language]{}, err
		}
		return ast.Set[s.Language]{Entry: ast.NodeRef{Tree: output.ID, Node: output.Root}, Trees: []ast.AST[s.Language]{output}}, nil
	})
}

type dialect struct{ namespace string }

func (dialect) ReservedNames() []string { return nil }

func (dialect) FieldName(name string) (string, error) {
	if !g.Identifier(name) || !native.Identifier(name) || strings.HasPrefix(name, "_") {
		return "", fmt.Errorf("unsupported C++ field name %q", name)
	}
	return name, nil
}
func (dialect) Parameter(a *g.Arena[s.Language], name string, o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindTypeParameter, g.Attr("name", name), o)
}
func path(a *g.Arena[s.Language], names []string, args []ast.NodeID, target ast.NodeID, o ast.NodeRef, global bool) ast.NodeID {
	var segments []ast.Child
	for i, name := range names {
		var children []ast.Child
		if i == len(names)-1 {
			children = g.Children(s.RoleArgument, args)
		}
		segment := a.Add(s.KindNameSegment, g.Attr("name", name), o, children...)
		if i == len(names)-1 && target != 0 {
			a.Target(segment, target)
		}
		segments = append(segments, g.Child(s.RoleSegment, segment))
	}
	var attrs []ast.Attribute
	if global {
		attrs = g.Attr("global", "true")
	}
	return a.Add(s.KindNamedType, attrs, o, segments...)
}
func (d dialect) Reference(a *g.Arena[s.Language], name string, target ast.NodeID, args []ast.NodeID, o ast.NodeRef) ast.NodeID {
	if a.Err != nil {
		return 0
	}
	if a.Tree.Nodes[target-1].Kind == s.KindTypeParameter {
		return path(a, []string{name}, args, target, o, false)
	}
	return path(a, []string{d.namespace, name}, args, target, o, true)
}
func (dialect) Text(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return path(a, []string{"std", "string"}, nil, 0, o, true)
}
func (dialect) Boolean(a *g.Arena[s.Language], o ast.NodeRef) ast.NodeID {
	return a.Add(s.KindBuiltinType, g.Attr("name", "bool"), o)
}
func (dialect) List(a *g.Arena[s.Language], element ast.NodeID, o ast.NodeRef) ast.NodeID {
	return path(a, []string{"std", "vector"}, []ast.NodeID{element}, 0, o, true)
}
func (dialect) Record(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, fields []g.Field) {
	if a.Err != nil {
		return
	}
	children := g.Children(s.RoleTemplateParameter, params)
	for _, field := range fields {
		for _, p := range params {
			for _, attr := range a.Tree.Nodes[p-1].Attributes {
				if attr.Name == "name" && attr.Text == field.Name {
					a.Err = &text.Diagnostic{Code: text.Unsupported, Message: "C++ field shadows template parameter", Node: &field.Origin}
					return
				}
			}
		}
		ty := field.Type
		if field.Optional {
			ty = path(a, []string{"std", "optional"}, []ast.NodeID{ty}, 0, field.Origin, true)
		}
		member := a.Add(s.KindField, g.Attr("name", field.Name), field.Origin, g.Child(s.RoleType, ty))
		children = append(children, g.Child(s.RoleMember, member))
	}
	a.Set(id, s.KindRecord, g.Attr("name", name, "record-key", "struct"), children...)
}
func (dialect) Alias(a *g.Arena[s.Language], id ast.NodeID, name string, params []ast.NodeID, body ast.NodeID) {
	a.Set(id, s.KindTypeAlias, g.Attr("name", name), append(g.Children(s.RoleTemplateParameter, params), g.Child(s.RoleType, body))...)
}
func (d dialect) Finish(a *g.Arena[s.Language], decls []ast.NodeID) ast.NodeID {
	var children []ast.Child
	for _, header := range []string{"optional", "string", "vector"} {
		include := a.Add(s.KindInclude, g.Attr("path", header, "system", "true"), ast.NodeRef{})
		children = append(children, g.Child(s.RoleInclude, include))
	}
	namespace := a.Add(s.KindNamespace, g.Attr("name", d.namespace), ast.NodeRef{}, g.Children(s.RoleDeclaration, decls)...)
	children = append(children, g.Child(s.RoleDeclaration, namespace))
	return a.Add(s.KindTranslationUnit, nil, ast.NodeRef{}, children...)
}
