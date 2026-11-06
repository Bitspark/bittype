package integration_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	rustprovider "github.com/Bitspark/bittype/provider/rust/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func rustDeclarationsFixture() ast.AST[rust.Language] {
	f := fixture[rust.Language]{tree: ast.AST[rust.Language]{ID: "independent-rust-declarations"}}
	path := func(name string, args ...ast.Child) ast.NodeID {
		return f.node("path-type", nil, child("segment", f.node("path-segment", attr("name", name), args...)))
	}
	declarationParam := f.node("type-parameter", attr("name", "T"), child("bound", path("Clone")))
	declarationWhere := f.node("where-predicate", nil, child("subject", path("T")), child("bound", path("Default")))
	declaration := f.node("associated-type", attr("name", "Item"), child("generic-parameter", declarationParam), child("bound", path("Clone")), child("where", declarationWhere))
	trait := f.node("trait", attr("name", "Family", "visibility", "public"), child("member", declaration))
	structure := f.node("struct", attr("name", "VectorFamily", "visibility", "public"))
	definitionParam := f.node("type-parameter", attr("name", "T"), child("bound", path("Clone")))
	definitionWhere := f.node("where-predicate", nil, child("subject", path("T")), child("bound", path("Default")))
	use := f.node("path-segment", attr("name", "T"))
	f.tree.Nodes[use-1].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: f.tree.ID, Node: definitionParam}}}
	f.tree.Nodes[use-1].Origin = ast.Origin{Spans: []ast.Span{{Source: "input", Start: 0, End: 1}}}
	value := path("Vec", child("argument", f.node("path-type", nil, child("segment", use))))
	definition := f.node("type-alias", attr("name", "Item"), child("generic-parameter", definitionParam), child("where", definitionWhere), child("type", value))
	impl := f.node("trait-impl", nil, child("trait", path("Family")), child("subject", path("VectorFamily")), child("member", definition))
	argument := f.node("primitive-type", attr("name", "u8"))
	suffix := f.node("path-segment", attr("name", "Item"), child("argument", argument))
	projection := f.node("qualified-path-type", nil, child("subject", path("VectorFamily")), child("trait", path("Family")), child("segment", suffix))
	alias := f.node("type-alias", attr("name", "Bytes", "visibility", "public"), child("type", projection))
	f.tree.Root = f.node("file", nil, child("item", trait), child("item", structure), child("item", impl), child("item", alias))
	return f.tree
}

func TestIndependentRustDeclarationFormatter(t *testing.T) {
	a := rustDeclarationsFixture()
	p := rustprovider.Provider{}
	c, err := p.Format(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.Abstract, a) {
		t.Fatal("formatter changed declarations, binder references or origins")
	}
	got, err := p.Print(context.Background(), c)
	want := "pub trait Family {\n    type Item<T: Clone>: Clone where T: Default;\n}\n\npub struct VectorFamily {}\n\nimpl Family for VectorFamily {\n    type Item<T: Clone> = Vec<T> where T: Default;\n}\n\npub type Bytes = <VectorFamily as Family>::Item<u8>;\n"
	if err != nil || got != want {
		t.Fatalf("want %q, got %q: %v", want, got, err)
	}
	reparsed, err := p.Parse(context.Background(), parse.Input{SourceID: "formatted", TreeID: "reparsed", Text: got})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(shape(a, a.Root), shape(reparsed.Abstract, reparsed.Abstract.Root)) {
		t.Fatal("independent declarations changed syntax on round trip")
	}
}

func TestRustDeclarationFormatting(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"impl Family for V{type Item<'a>=&'a str where Self:'a;}", "impl Family for V {\n    type Item<'a> = &'a str where Self: 'a;\n}\n"},
		{"impl Family for V{/* keep in CST */type Item<T>where T:Clone=Vec<T>;}", "impl Family for V {\n    type Item<T> = Vec<T> where T: Clone;\n}\n"},
		{"trait Family{type Item<'a,T:?Sized> : ?Sized+'a where T:'a;}", "trait Family {\n    type Item<'a, T: ?Sized>: ?Sized + 'a where T: 'a;\n}\n"},
		{"impl<'a,T:?Sized> ::api::Trait<T> for &'a T where T:'a{}", "impl<'a, T: ?Sized> ::api::Trait<T> for &'a T where T: 'a {\n}\n"},
		{"mod nested{impl Trait for (){type r#type<T>=T;}}", "mod nested {\n    impl Trait for () {\n        type r#type<T> = T;\n    }\n}\n"},
	} {
		t.Run(tc.source, func(t *testing.T) {
			p := rustprovider.Provider{}
			c, err := p.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: tc.source})
			if err != nil {
				t.Fatal(err)
			}
			if original, err := p.Print(context.Background(), c); err != nil || original != tc.source {
				t.Fatalf("original CST changed: %v", err)
			}
			got := formatted(t, p, c.Abstract)
			if got != tc.want {
				t.Fatalf("want %q, got %q", tc.want, got)
			}
			reparsed, err := p.Parse(context.Background(), parse.Input{SourceID: "formatted", TreeID: "reparsed", Text: got})
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(shape(c.Abstract, c.Abstract.Root), shape(reparsed.Abstract, reparsed.Abstract.Root)) || formatted(t, p, reparsed.Abstract) != got {
				t.Fatal("declaration syntax or formatting is unstable")
			}
		})
	}
}
