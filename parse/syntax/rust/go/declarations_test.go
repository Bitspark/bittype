package parse_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	rustparse "github.com/Bitspark/bittype/parse/syntax/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

func TestGenericAssociatedDeclarationStructure(t *testing.T) {
	source := "trait Family { type Item<'a, T: Clone>: Clone + 'a where T: 'a; } impl<U> Family for Wrap<U> where U: Clone { type Item<'a, T: Clone> = Vec<T> where T: 'a; }"
	c, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	v := tr.View[s.Language]{Tree: c.Abstract}
	items := v.Children(c.Abstract.Root, "item")
	declaration := v.Child(items[0], "member")
	if v.Kind(declaration) != "associated-type" || len(v.Children(declaration, "generic-parameter")) != 2 || len(v.Children(declaration, "bound")) != 2 || len(v.Children(declaration, "where")) != 1 {
		t.Fatal("associated declaration lost its parameters, result bounds or predicates")
	}
	impl := items[1]
	if v.Kind(impl) != "trait-impl" || len(v.Children(impl, "generic-parameter")) != 1 || v.Kind(v.Child(impl, "trait")) != "path-type" || v.Kind(v.Child(impl, "subject")) != "path-type" || len(v.Children(impl, "where")) != 1 {
		t.Fatal("impl head lost its structured children")
	}
	definition := v.Child(impl, "member")
	if v.Kind(definition) != "type-alias" || len(v.Children(definition, "generic-parameter")) != 2 || len(v.Children(definition, "where")) != 1 || v.Child(definition, "type") == 0 || len(v.Children(definition, "bound")) != 0 {
		t.Fatal("impl definition must reuse the type-alias shape")
	}
	if v.Children(declaration, "generic-parameter")[1] == v.Children(definition, "generic-parameter")[1] {
		t.Fatal("same-spelled parameters in declaration and definition must remain separate nodes")
	}
	span := v.Node(definition).Origin.Spans[0]
	if source[span.Start:span.End] != "type Item<'a, T: Clone> = Vec<T> where T: 'a;" {
		t.Fatalf("definition origin: %#v", span)
	}
}

func TestGenericAssociatedSyntaxWithoutResolution(t *testing.T) {
	for _, source := range []string{
		"trait X { type Item<T>; type Other; }",
		"trait X { type Item<'a: 'b, 'b, T: ?Sized + 'a>: ?Sized + 'b where T: Clone, 'a: 'b,; }",
		"trait X { type Item<T>: Missing<T> where <T as Unknown>::Item: Clone; }",
		"impl Missing for Unknown {}",
		"impl<'a, T: ?Sized> Missing<T> for &'a T where T: 'a { type Item = str; }",
		"impl<T> ::api::Family<T> for [T; 2] { type Item<U> = (T, U); }",
		"impl Family for () { type r#type<T> = T; }",
		"impl Family for Unknown { type Item<T> where T: Clone = Vec<T>; }",
		"impl Family for Unknown { type Item<T> = Vec<T> where T: Clone,; }",
	} {
		t.Run(source, func(t *testing.T) {
			if _, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMalformedGenericAssociatedDeclarations(t *testing.T) {
	for _, source := range []string{
		"trait X { type Item<>; }", "trait X { type Item<T> }", "trait X { type Item<T>:; }",
		"trait X { type Item<'a = 'b>; }", "trait X { type Item<T = u8>; }",
		"trait X { type Item<T, 'a>; }", "trait X { pub type Item<T>; }",
		"trait X { type Item<T> where; }", "trait X { type Item<'a> where 'a: Clone; }",
		"trait X { type Item<'a> where 'a: ?Sized; }", "trait X { type Item: ??Sized; }",
		"trait X: ?Sized {}", "type X = dyn ?Sized;",
		"impl Family for {}", "impl Family for V { type Item; }",
		"impl<T = u8> Family for V {}", "impl Family for V { type Item<T = u8> = T; }",
		"impl Family for V { type Item: Clone = u8; }", "impl Family for V { type Item<T>: Clone = T; }",
		"impl Family for V { type Item<T> = ; }", "impl Family for V { type Item<T> = T }",
		"impl Family for V { type Item<T> where T: Clone = T where T: Copy; }",
		"impl Family for V where {}", "pub impl Family for V {}",
		"impl Family for V {} trailing", "impl<> Family for V {}",
	} {
		t.Run(source, func(t *testing.T) {
			c, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source})
			var diagnostic *text.Diagnostic
			if !errors.As(err, &diagnostic) || diagnostic.Span == nil || len(c.Nodes) != 0 || len(c.Abstract.Nodes) != 0 {
				t.Fatalf("expected located rejection without partial output: %v", err)
			}
		})
	}
}

func TestImplScopeRejection(t *testing.T) {
	for _, source := range []string{
		"impl V {}", "impl<T> V<T> { type Item = T; }",
		"impl !Family for V {}", "unsafe impl Family for V {}",
		"impl Family for V { fn method() {} }", "impl Family for V { fn method(); }",
		"impl Family for V { type Item = u8; const VALUE: u8 = 0; }",
		"impl Family for V { type Item = u8; fn method() {} }",
		"impl Family for V { pub type Item = u8; }", "impl Family for V { struct Nested {} }",
		"trait X { type Item<T> = Vec<T>; }", "trait X { type Item = u8; }",
		"impl Family for V { default type Item = u8; }",
	} {
		t.Run(source, func(t *testing.T) {
			c, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source})
			var diagnostic *text.Diagnostic
			if !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || len(c.Nodes) != 0 || len(c.Abstract.Nodes) != 0 {
				t.Fatalf("expected whole-input unsupported rejection: %v", err)
			}
		})
	}
}
