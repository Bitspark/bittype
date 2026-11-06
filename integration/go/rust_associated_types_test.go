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

func rustAssociatedFixture() ast.AST[rust.Language] {
	f := fixture[rust.Language]{}
	path := func(name string, args ...ast.Child) ast.NodeID {
		segment := f.node("path-segment", attr("name", name), args...)
		return f.node("path-type", nil, child("segment", segment))
	}
	parameter := f.node("type-parameter", attr("name", "I"), child("bound", path("Iterator")))
	subject, trait := path("I"), path("Iterator")
	item := f.node("path-segment", attr("name", "Item"))
	projection := f.node("qualified-path-type", nil, child("subject", subject), child("trait", trait), child("segment", item))
	alias := f.node("type-alias", attr("name", "Item", "visibility", "public"), child("generic-parameter", parameter), child("type", projection))
	name := f.node("path-segment", attr("name", "Item"))
	value := f.node("primitive-type", attr("name", "u8"))
	binding := f.node("associated-type-binding", nil, child("path", name), child("type", value))
	bound := path("Iterator", child("argument", binding))
	decl := f.node("trait", attr("name", "Bytes", "visibility", "public"), child("bound", bound))
	f.tree.Root = f.node("file", nil, child("item", alias), child("item", decl))
	f.tree.ID = "independent-associated-types"
	for _, id := range []ast.NodeID{item, name} {
		f.tree.Nodes[id-1].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "external", Node: 1}}}
		f.tree.Nodes[id-1].Origin = ast.Origin{Spans: []ast.Span{{Source: "input", Start: 0, End: 4}}}
	}
	return f.tree
}

// bittype:case target=github.com/Bitspark/bittype/provider/rust/go#Provider.Format feature=associated-type-projections case=independent-ast-formatting expects=ast-preserved-and-golden-text
func TestIndependentRustAssociatedFormatter(t *testing.T) {
	a := rustAssociatedFixture()
	p := rustprovider.Provider{}
	c, err := p.Format(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.Abstract, a) {
		t.Fatal("formatter changed the projection/binding AST or its references/origins")
	}
	got, err := p.Print(context.Background(), c)
	want := "pub type Item<I: Iterator> = <I as Iterator>::Item;\n\npub trait Bytes: Iterator<Item = u8> {\n}\n"
	if err != nil || got != want {
		t.Fatalf("want %q, got %q: %v", want, got, err)
	}
	reparsed, err := p.Parse(context.Background(), parse.Input{SourceID: "formatted", TreeID: "reparsed", Text: got})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(shape(a, a.Root), shape(reparsed.Abstract, reparsed.Abstract.Root)) {
		t.Fatal("formatting and parsing changed independently authored syntax")
	}
}

func TestRustAssociatedFormatting(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"type X=<T> /* projection */ :: Item;", "type X = <T>::Item;\n"},
		{"type X=<<T as A>::B as ::api::Trait::<U>>::Item::<'a,V>::Leaf;", "type X = <<T as A>::B as ::api::Trait<U>>::Item<'a, V>::Leaf;\n"},
		{"type X=dyn Source<'a,T,Item::<'a>=&'a T,Other=<U>::Result,>;", "type X = dyn Source<'a, T, Item<'a> = &'a T, Other = <U>::Result>;\n"},
		{"type X=dyn Source<r#type=<T as Trait>::r#type,bool=u8>;", "type X = dyn Source<r#type = <T as Trait>::r#type, bool = u8>;\n"},
		{"type X=<(T,U) as Trait>::Item;", "type X = <(T, U) as Trait>::Item;\n"},
		{"type X=<&(dyn A+B) as Trait>::Item;", "type X = <&(dyn A + B) as Trait>::Item;\n"},
	} {
		t.Run(tc.source, func(t *testing.T) {
			p := rustprovider.Provider{}
			c, err := p.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: tc.source})
			if err != nil {
				t.Fatal(err)
			}
			if source, err := p.Print(context.Background(), c); err != nil || source != tc.source {
				t.Fatalf("parser lost exact source: %v", err)
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
				t.Fatal("projection/binding formatting is not a stable syntax round trip")
			}
		})
	}
}
