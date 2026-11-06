package check_test

import (
	"errors"
	"testing"

	check "github.com/Bitspark/bittype/check/syntax/rust/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func declarationsTree() (ast.AST[s.Language], map[string]ast.NodeID) {
	tree := ast.AST[s.Language]{ID: "declarations"}
	ids := map[string]ast.NodeID{}
	add := func(kind s.Kind, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
		tree.Nodes = append(tree.Nodes, ast.Node[s.Language]{Kind: kind, Attributes: attrs, Children: children})
		return ast.NodeID(len(tree.Nodes))
	}
	path := func(name string) ast.NodeID {
		segment := add(s.KindPathSegment, tr.Attrs("name", name))
		return add(s.KindPathType, nil, tr.Child("segment", segment))
	}
	member := func(definition bool) ast.NodeID {
		parameter := add(s.KindTypeParameter, tr.Attrs("name", "T"), tr.Child("bound", path("Clone")))
		where := add(s.KindWherePredicate, nil, tr.Child("subject", path("T")), tr.Child("bound", path("Default")))
		if definition {
			ids["def-param"] = parameter
			return add(s.KindTypeAlias, tr.Attrs("name", "Item"), tr.Child("generic-parameter", parameter), tr.Child("where", where), tr.Child("type", path("T")))
		}
		ids["decl-param"] = parameter
		return add(s.KindAssociatedType, tr.Attrs("name", "Item"), tr.Child("generic-parameter", parameter), tr.Child("bound", path("Clone")), tr.Child("where", where))
	}
	ids["declaration"] = member(false)
	trait := add(s.KindTrait, tr.Attrs("name", "Family"), tr.Child("member", ids["declaration"]))
	ids["definition"] = member(true)
	ids["impl-param"] = add(s.KindTypeParameter, tr.Attrs("name", "U"))
	ids["impl"] = add(s.KindTraitImpl, nil, tr.Child("generic-parameter", ids["impl-param"]), tr.Child("trait", path("Family")), tr.Child("subject", path("Unknown")), tr.Child("member", ids["definition"]))
	tree.Root = add(s.KindFile, nil, tr.Child("item", trait), tr.Child("item", ids["impl"]))
	return tree, ids
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/rust/go#Check feature=generic-associated-declarations case=independent-ast expects=accept
func TestGenericAssociatedDeclarationAdmission(t *testing.T) {
	tree, _ := declarationsTree()
	if err := check.Check(tree); err != nil {
		t.Fatal(err)
	}
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/rust/go#Check feature=generic-associated-declarations case=relaxed-bound-context expects=context-sensitive-admission
func TestRelaxedBoundContext(t *testing.T) {
	for _, tc := range []struct {
		owner string
		index int
		valid bool
	}{{"declaration", 1, true}, {"definition", 2, false}} {
		t.Run(tc.owner, func(t *testing.T) {
			tree, ids := declarationsTree()
			owner := ids[tc.owner]
			path := tree.Nodes[owner-1].Children[tc.index].Node
			tree.Nodes = append(tree.Nodes, ast.Node[s.Language]{Kind: s.KindRelaxedBound, Children: []ast.Child{tr.Child("path", path)}})
			tree.Nodes[owner-1].Children[tc.index].Node = ast.NodeID(len(tree.Nodes))
			if err := check.Check(tree); (err == nil) != tc.valid {
				t.Fatalf("bound modifiers belong in bounds, not type expressions: %v", err)
			}
		})
	}
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/rust/go#Check feature=generic-associated-declarations case=invalid-declaration-contexts expects=located-rejection
func TestGenericAssociatedContextRules(t *testing.T) {
	cases := map[string]func(*ast.AST[s.Language], map[string]ast.NodeID){
		"impl has no visibility": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["impl"]-1].Attributes = tr.Attrs("visibility", "public")
		},
		"impl member inherits visibility": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["definition"]-1].Attributes = append(a.Nodes[ids["definition"]-1].Attributes, tr.Attrs("visibility", "public")...)
		},
		"definition requires type": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["definition"]-1].Children[2].Role = "bound"
		},
		"declaration cannot define type": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["declaration"]-1].Children[1].Role = "type"
		},
		"impl cannot contain declarations": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["definition"]-1].Kind = s.KindAssociatedType
		},
		"trait impl requires trait and subject": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["impl"]-1].Children[1].Role = "subject"
		},
		"definition children remain ordered": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			children := a.Nodes[ids["definition"]-1].Children
			children[1], children[2] = children[2], children[1]
		},
	}
	for _, key := range []string{"decl-param", "def-param", "impl-param"} {
		cases["no default on "+key] = func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes = append(a.Nodes, ast.Node[s.Language]{Kind: s.KindPrimitiveType, Attributes: tr.Attrs("name", "u8")})
			id := ids[key]
			a.Nodes[id-1].Children = append(a.Nodes[id-1].Children, tr.Child("default", ast.NodeID(len(a.Nodes))))
		}
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			tree, ids := declarationsTree()
			mutate(&tree, ids)
			var diagnostic *text.Diagnostic
			if err := check.Check(tree); !errors.As(err, &diagnostic) || diagnostic.Node == nil {
				t.Fatalf("expected node-located syntax rejection: %v", err)
			}
		})
	}
}
