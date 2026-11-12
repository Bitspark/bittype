package format

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	"github.com/Bitspark/bittype/source/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// The formatter cannot be given these two inputs: the TypeScript bridge refuses
// a namespace import and never publishes a binding without a single target. The
// integration table covers every refusal that a formatter call reaches.

// category requires that a refusal matches want and none of the other values.
func category(t *testing.T, err, want error) {
	t.Helper()
	for _, value := range []error{source.ErrInvalid, source.ErrUnsupported, semantic.ErrEvidenceUnavailable, source.ErrLimit} {
		if errors.Is(err, value) != (value == want) {
			t.Errorf("errors.Is(err, %q) = %v for %q", value, errors.Is(err, value), err)
		}
	}
}

func attributes(pairs ...string) []ast.Attribute {
	var out []ast.Attribute
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, ast.Attribute{Name: pairs[i], Text: pairs[i+1]})
	}
	return out
}

func TestImportFormOutsideTheProfileIsUnsupported(t *testing.T) {
	root := func(role string, kind s.Kind) s.AST {
		return s.AST{ID: "input", Root: 1, Nodes: []s.Node{
			{Kind: s.KindSourceFile, Children: []ast.Child{{Role: role, Node: 2}}},
			{Kind: s.KindImport, Attributes: attributes(s.AttributeModule, "@example/all"), Children: []ast.Child{{Role: s.RoleBinding, Node: 3}}},
			{Kind: kind, Attributes: attributes(s.AttributeName, "all")},
		}}
	}
	if sites, err := declarations(root(s.RoleImport, s.KindImportBinding)); err != nil || len(sites) != 0 {
		t.Fatal("a named type import is inside the profile", err)
	}
	sites, err := declarations(root(s.RoleImport, s.KindNamespaceImport))
	if err == nil || sites != nil {
		t.Fatal("a namespace import was read as a named type import")
	}
	category(t, err, source.ErrUnsupported)
	// A child of the file that is neither an import nor a declaration is no
	// TypeScript file of any profile.
	sites, err = declarations(root("statement", s.KindImportBinding))
	if err == nil || sites != nil {
		t.Fatal("an unknown child of the file was read")
	}
	category(t, err, source.ErrInvalid)
}

type admitAll struct{ s.Language }

func (admitAll) Check(ast.View[s.Language]) error { return nil }

func TestAmbiguousNativeBindingIsInvalid(t *testing.T) {
	ctx := context.Background()
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](definition, admitAll{})
	if err != nil {
		t.Fatal(err)
	}
	tree := s.AST{ID: "input", Root: 1, Nodes: []s.Node{
		{Kind: s.KindSourceFile, Children: []ast.Child{{Role: s.RoleDeclaration, Node: 2}}},
		{Kind: s.KindTypeAlias, Attributes: attributes(s.AttributeName, "Use", s.AttributeExported, "true"), Children: []ast.Child{{Role: s.RoleType, Node: 3}}},
		{Kind: s.KindTypeReference, Attributes: attributes(s.AttributeName, "Twice")},
	}}
	packed, err := representation.Pack(definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: ref(tree.ID, tree.Root), Trees: []representation.AST{packed}, Files: &files.Context{
		Files: files.Files{Entries: []files.File{{ID: "input", Path: "input.ts", Text: "export type Use = Twice;"}}},
		Links: []files.SyntaxLink{{File: "input", Tree: tree.ID, Source: "input"}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	graph := func(target semantics.Entity) semantics.Semantics[s.Language] {
		built, err := semantics.New[s.Language](semantics.Data{
			Entities: []semantics.Entity{{Kind: n.ModuleKind(), Name: "input", Declarations: []ast.NodeRef{ref(tree.ID, 1)}}, target},
			Bindings: []semantics.Binding{{At: semantics.Occurrence{Site: ref(tree.ID, 3), Role: n.NameRole()}, Target: 2}},
		}, semantics.DefaultLimits(), func(site ast.NodeRef) bool { return site.Tree == tree.ID && int(site.Node) <= len(tree.Nodes) })
		if err != nil {
			t.Fatal(err)
		}
		return built
	}
	if p, err := readProgram(ctx, snapshot, graph(semantics.Entity{Kind: n.BuiltinKind(), Name: "Twice"})); err != nil || p.bindings[ref(tree.ID, 3)].builtin != "Twice" {
		t.Fatal("a binding with one target must be read", err)
	}
	// The binding's target has no declaration site and is neither a package type
	// nor a builtin, so the output does not say what the name refers to.
	p, err := readProgram(ctx, snapshot, graph(semantics.Entity{Kind: n.InterfaceKind(), Name: "Twice"}))
	if err == nil || p != nil {
		t.Fatal("a binding without one target was read")
	}
	category(t, err, source.ErrInvalid)
}
