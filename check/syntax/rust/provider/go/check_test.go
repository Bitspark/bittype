package check_test

import (
	"errors"
	"testing"

	check "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

// Independently authored syntax: Unresolved<'a, U, Item = <T as Trait>::Item>.
func associatedTree() ast.AST[s.Language] {
	return ast.AST[s.Language]{ID: "associated", Root: 15, Nodes: []ast.Node[s.Language]{
		{Kind: s.KindPathSegment, Attributes: tr.Attrs("name", "T")},
		{Kind: s.KindPathType, Children: []ast.Child{tr.Child("segment", 1)}},
		{Kind: s.KindPathSegment, Attributes: tr.Attrs("name", "Trait")},
		{Kind: s.KindPathType, Children: []ast.Child{tr.Child("segment", 3)}},
		{Kind: s.KindPathSegment, Attributes: tr.Attrs("name", "Item")},
		{Kind: s.KindQualifiedPathType, Children: []ast.Child{tr.Child("subject", 2), tr.Child("trait", 4), tr.Child("segment", 5)}},
		{Kind: s.KindPathSegment, Attributes: tr.Attrs("name", "Item")},
		{Kind: s.KindAssociatedTypeBinding, Children: []ast.Child{tr.Child("path", 7), tr.Child("type", 6)}},
		{Kind: s.KindLifetimeReference, Attributes: tr.Attrs("name", "a")},
		{Kind: s.KindPathSegment, Attributes: tr.Attrs("name", "U")},
		{Kind: s.KindPathType, Children: []ast.Child{tr.Child("segment", 10)}},
		{Kind: s.KindPathSegment, Attributes: tr.Attrs("name", "Unresolved"), Children: []ast.Child{tr.Child("argument", 9), tr.Child("argument", 11), tr.Child("argument", 8)}},
		{Kind: s.KindPathType, Children: []ast.Child{tr.Child("segment", 12)}},
		{Kind: s.KindTrait, Attributes: tr.Attrs("name", "X"), Children: []ast.Child{tr.Child("bound", 13)}},
		{Kind: s.KindFile, Children: []ast.Child{tr.Child("item", 14)}},
	}}
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/rust/provider/go#Check feature=associated-type-projections case=independent-ast-and-references expects=accept
func TestAssociatedTypeAdmission(t *testing.T) {
	tree := associatedTree()
	// The syntax checker accepts resolved references without interpreting them.
	tree.Nodes[4].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "external", Node: 7}}}
	tree.Nodes[6].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "external", Node: 8}}}
	if err := check.Check(tree); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"r#type", "bool"} {
		tree.Nodes[6].Attributes = tr.Attrs("name", name)
		if err := check.Check(tree); err != nil {
			t.Fatal(err)
		}
	}
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/rust/provider/go#Check feature=associated-type-projections case=invalid-shapes-and-ordering expects=located-rejection
func TestInvalidAssociatedTypeShapes(t *testing.T) {
	cases := map[string]func(*ast.AST[s.Language]){
		"missing subject":          func(a *ast.AST[s.Language]) { a.Nodes[5].Children[0].Role = "type" },
		"missing suffix":           func(a *ast.AST[s.Language]) { a.Nodes[5].Children[2].Role = "trait" },
		"trait must be path":       func(a *ast.AST[s.Language]) { a.Nodes[5].Children[1].Node, a.Nodes[5].Children[2].Node = 5, 4 },
		"binding requires value":   func(a *ast.AST[s.Language]) { a.Nodes[7].Children[1].Role = "subject" },
		"binding requires segment": func(a *ast.AST[s.Language]) { a.Nodes[7].Children[0].Node, a.Nodes[7].Children[1].Node = 6, 7 },
		"binding is not a type":    func(a *ast.AST[s.Language]) { a.Nodes[11].Children[2].Role = "type" },
		"type after binding": func(a *ast.AST[s.Language]) {
			a.Nodes[11].Children[1], a.Nodes[11].Children[2] = a.Nodes[11].Children[2], a.Nodes[11].Children[1]
		},
		"lifetime after type": func(a *ast.AST[s.Language]) {
			a.Nodes[11].Children[0], a.Nodes[11].Children[1] = a.Nodes[11].Children[1], a.Nodes[11].Children[0]
		},
		"lifetime after binding": func(a *ast.AST[s.Language]) {
			a.Nodes[11].Children[0], a.Nodes[11].Children[2] = a.Nodes[11].Children[2], a.Nodes[11].Children[0]
		},
	}
	for _, name := range []string{"_", "Self", "self", "crate", "super"} {
		cases["binding name "+name] = func(a *ast.AST[s.Language]) { a.Nodes[6].Attributes = tr.Attrs("name", name) }
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			tree := associatedTree()
			mutate(&tree)
			var diagnostic *text.Diagnostic
			if err := check.Check(tree); !errors.As(err, &diagnostic) || diagnostic.Node == nil {
				t.Fatalf("expected node-located rejection: %v", err)
			}
		})
	}
}
