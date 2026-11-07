package check_test

import (
	"errors"
	"testing"

	check "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func syntaxTree() (ast.AST[s.Language], map[string]ast.NodeID) {
	tree := ast.AST[s.Language]{ID: "conditional-mapped"}
	ids := map[string]ast.NodeID{}
	add := func(kind s.Kind, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
		tree.Nodes = append(tree.Nodes, ast.Node[s.Language]{Kind: kind, Attributes: attrs, Children: children})
		return ast.NodeID(len(tree.Nodes))
	}
	ref := func(name string) ast.NodeID { return add(s.KindTypeReference, tr.Attrs("name", name)) }
	ids["parameter"] = add(s.KindTypeParameter, tr.Attrs("name", "K"), tr.Child("constraint", ref("UnresolvedKeys")))
	ids["inferred-parameter"] = add(s.KindTypeParameter, tr.Attrs("name", "U"))
	ids["infer"] = add(s.KindInferType, nil, tr.Child("type-parameter", ids["inferred-parameter"]))
	ids["conditional"] = add(s.KindConditionalType, nil, tr.Child("check", ref("T")), tr.Child("constraint", ids["infer"]), tr.Child("true", ref("U")), tr.Child("false", ref("Fallback")))
	ids["mapped"] = add(s.KindMappedType, tr.Attrs("readonly-modifier", "remove", "optional-modifier", "add"), tr.Child("type-parameter", ids["parameter"]), tr.Child("name-type", ref("K")), tr.Child("type", ids["conditional"]))
	alias := add(s.KindTypeAlias, tr.Attrs("name", "X"), tr.Child("type", ids["mapped"]))
	tree.Root = add(s.KindSourceFile, nil, tr.Child("declaration", alias))
	return tree, ids
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Check feature=conditional-mapped-syntax case=independent-ast expects=accept
func TestConditionalMappedAdmission(t *testing.T) {
	tree, _ := syntaxTree()
	if err := check.Check(tree); err != nil {
		t.Fatal(err)
	}
}

// bittype:case target=github.com/Bitspark/bittype/check/syntax/typescript/provider/go#Check feature=conditional-mapped-syntax case=invalid-shapes-and-contexts expects=located-rejection
func TestInvalidConditionalMappedShapes(t *testing.T) {
	cases := map[string]func(*ast.AST[s.Language], map[string]ast.NodeID){
		"conditional branch role": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["conditional"]-1].Children[3].Role = "result"
		},
		"conditional role order": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			children := a.Nodes[ids["conditional"]-1].Children
			children[2], children[3] = children[3], children[2]
		},
		"infer in check": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			children := a.Nodes[ids["conditional"]-1].Children
			children[0].Node, children[1].Node = children[1].Node, children[0].Node
		},
		"infer in branch": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			children := a.Nodes[ids["conditional"]-1].Children
			children[3].Node, children[1].Node = children[1].Node, children[3].Node
		},
		"mapped requires in type": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			id := ids["parameter"]
			children := a.Nodes[id-1].Children
			// Preserve the key subtree elsewhere, so only mapped admission fails.
			parent := ids["mapped"]
			nameType := a.Nodes[parent-1].Children[1].Node
			a.Nodes[nameType-1].Kind = s.KindTypeOperator
			a.Nodes[nameType-1].Attributes = tr.Attrs("operator", "keyof")
			a.Nodes[nameType-1].Children = []ast.Child{tr.Child("type", children[0].Node)}
			a.Nodes[id-1].Children = nil
		},
		"invalid modifier": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["mapped"]-1].Attributes[0].Text = "false"
		},
		"mapped cannot have members": func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes[ids["mapped"]-1].Children[2].Role = "member"
		},
	}
	for _, key := range []string{"parameter", "inferred-parameter"} {
		cases["no default on "+key] = func(a *ast.AST[s.Language], ids map[string]ast.NodeID) {
			a.Nodes = append(a.Nodes, ast.Node[s.Language]{Kind: s.KindKeywordType, Attributes: tr.Attrs("name", "string")})
			id := ids[key]
			a.Nodes[id-1].Children = append(a.Nodes[id-1].Children, tr.Child("default", ast.NodeID(len(a.Nodes))))
		}
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			tree, ids := syntaxTree()
			mutate(&tree, ids)
			var d *text.Diagnostic
			if err := check.Check(tree); !errors.As(err, &d) || d.Node == nil {
				t.Fatalf("expected node-located syntax rejection: %v", err)
			}
		})
	}
}
