package declsyntax

import (
	"strings"
	"testing"

	structure "github.com/Bitspark/bittype/check/ast/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

type kindLanguage struct{}

// F has a constructor interface, but Missing is deliberately unresolved. The
// syntax checker must validate their representation without deciding meaning.
func kindFixture() ast.AST[kindLanguage] {
	return ast.AST[kindLanguage]{ID: "kinds", Root: 11, Nodes: []ast.Node[kindLanguage]{
		{Kind: KindInterface, Attributes: []ast.Attribute{{Name: AttributeName, Text: "type"}}},
		{Kind: KindInterface, Attributes: []ast.Attribute{{Name: AttributeName, Text: "data"}}},
		{Kind: KindInterface, Attributes: []ast.Attribute{{Name: AttributeName, Text: "constructor"}}, Children: []ast.Child{{Role: RoleInput, Node: 1}, {Role: RoleResult, Node: 2}}},
		{Kind: KindTypeParameter, Attributes: []ast.Attribute{{Name: AttributeName, Text: "F"}, {Name: AttributeIndex, Text: "0"}}, Children: []ast.Child{{Role: RoleInterface, Node: 3}}},
		{Kind: KindInterface, Attributes: []ast.Attribute{{Name: AttributeName, Text: "type"}}},
		{Kind: KindParameter, Attributes: []ast.Attribute{{Name: AttributeIndex, Text: "0"}}},
		{Kind: KindReference, Attributes: []ast.Attribute{{Name: AttributeScope, Text: "local"}, {Name: AttributeName, Text: "Missing"}}},
		{Kind: KindApplication, Children: []ast.Child{{Role: RoleCallee, Node: 6}, {Role: RoleArgument, Node: 7}}},
		{Kind: KindDeclaration, Attributes: []ast.Attribute{{Name: AttributeName, Text: "Apply"}}, Children: []ast.Child{{Role: RoleParameter, Node: 4}, {Role: RoleResult, Node: 5}, {Role: RoleBody, Node: 8}}},
		{Kind: KindModule, Attributes: []ast.Attribute{{Name: AttributeName, Text: "example.kinds"}, {Name: AttributePin, Text: strings.Repeat("0", 64)}}, Children: []ast.Child{{Role: RoleDeclaration, Node: 9}}},
		{Kind: KindProgram, Children: []ast.Child{{Role: RoleModule, Node: 10}}, References: []ast.Reference{{Role: RoleRoot, Target: ast.NodeRef{Tree: "kinds", Node: 10}}}},
	}}
}

func TestKindInterfaceSyntax(t *testing.T) {
	if err := Check(kindFixture()); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		change func(*ast.AST[kindLanguage])
	}{
		{"legacy bound kind", func(tree *ast.AST[kindLanguage]) { tree.Nodes[0].Kind = "bound" }},
		{"legacy bound role", func(tree *ast.AST[kindLanguage]) { tree.Nodes[3].Children[0].Role = "bound" }},
		{"unknown interface", func(tree *ast.AST[kindLanguage]) { tree.Nodes[0].Attributes[0].Text = "Value" }},
		{"proper interface children", func(tree *ast.AST[kindLanguage]) { tree.Nodes[2].Attributes[0].Text = "type" }},
		{"missing constructor result", func(tree *ast.AST[kindLanguage]) { tree.Nodes[2].Children[1].Role = RoleInput }},
		{"constructor result", func(tree *ast.AST[kindLanguage]) {
			tree.Nodes[3].Children[0].Node = 5
			tree.Nodes[8].Children[1].Node = 3
		}},
		{"legacy application target", func(tree *ast.AST[kindLanguage]) {
			tree.Nodes[7].Attributes = []ast.Attribute{{Name: AttributeScope, Text: "local"}, {Name: AttributeName, Text: "Missing"}}
		}},
		{"application binding annotation", func(tree *ast.AST[kindLanguage]) {
			tree.Nodes[7].References = []ast.Reference{{Role: RoleTarget, Target: ast.NodeRef{Tree: "kinds", Node: 9}}}
		}},
		{"missing callee", func(tree *ast.AST[kindLanguage]) { tree.Nodes[7].Children[0].Role = RoleArgument }},
		{"imported reference without alias", func(tree *ast.AST[kindLanguage]) { tree.Nodes[6].Attributes[0].Text = "imported" }},
		{"local reference with alias", func(tree *ast.AST[kindLanguage]) {
			tree.Nodes[6].Attributes = append(tree.Nodes[6].Attributes, ast.Attribute{Name: AttributeAlias, Text: "foreign"})
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			tree := kindFixture()
			test.change(&tree)
			if err := structure.Structure(tree); err != nil {
				t.Fatalf("fixture must reach syntax admission: %v", err)
			}
			if err := Check(tree); err == nil {
				t.Fatal("malformed revision-3 syntax admitted")
			}
		})
	}
}
