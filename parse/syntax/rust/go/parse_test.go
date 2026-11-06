package parse_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	rustparse "github.com/Bitspark/bittype/parse/syntax/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
)

func TestAssociatedTypeStructure(t *testing.T) {
	source := "type X = <T as ::api::Source<'a, U, Item<'a> = <V>::Output>>::Result<W>::Leaf;"
	c, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	v := tr.View[s.Language]{Tree: c.Abstract}
	projection := v.Child(v.Child(c.Abstract.Root, "item"), "type")
	if v.Kind(projection) != "qualified-path-type" || v.Kind(v.Child(projection, "subject")) != "path-type" {
		t.Fatal("projection must retain its subject as a type child")
	}
	trait := v.Child(projection, "trait")
	segments := v.Children(trait, "segment")
	if !v.Flag(trait, "absolute") || len(segments) != 2 || v.Attr(segments[1], "name") != "Source" {
		t.Fatal("qualified trait path lost")
	}
	args := v.Children(segments[1], "argument")
	if len(args) != 3 || v.Kind(args[0]) != "lifetime-reference" || v.Kind(args[1]) != "path-type" || v.Kind(args[2]) != "associated-type-binding" {
		t.Fatal("generic argument roles lost")
	}
	binding := args[2]
	name := v.Child(binding, "path")
	if v.Kind(name) != "path-segment" || v.Attr(name, "name") != "Item" || len(v.Children(name, "argument")) != 1 {
		t.Fatal("binding must retain the associated name and its own generic arguments")
	}
	value := v.Child(binding, "type")
	if v.Kind(value) != "qualified-path-type" || v.Child(value, "trait") != 0 {
		t.Fatal("binding RHS must retain the projection without an as-clause")
	}
	suffix := v.Children(projection, "segment")
	if len(suffix) != 2 || v.Attr(suffix[0], "name") != "Result" || len(v.Children(suffix[0], "argument")) != 1 || v.Attr(suffix[1], "name") != "Leaf" {
		t.Fatal("projection suffix lost")
	}
	span := v.Node(binding).Origin.Spans[0]
	if span.Source != "source" || source[span.Start:span.End] != "Item<'a> = <V>::Output" {
		t.Fatalf("binding origin: %#v", span)
	}
}

func TestAssociatedTypeSyntaxWithoutResolution(t *testing.T) {
	for _, source := range []string{
		"type X = <Unknown>::Item;",
		"type X = <<T as Outer>::Item as Inner>::Result;",
		"type X = <&'a mut [T] as Source>::Item<'a>;",
		"type X = <(T, U) as Trait>::Item;",
		"type X = <fn(T) -> U as Trait>::Item;",
		"type X = <T as Trait::<U>>::Item::<'a, V>::Leaf;",
		"type X = <T as bool>::bool;",
		"trait X: Source<'a, T, Item<'a> = &'a T, Other = <U as Trait>::Item> {}",
		"trait X: Source<Item::<T> = U> {}",
		"trait X: Source<r#type = <T>::r#type, bool = u8> {}",
		"trait X where <T as Iterator>::Item: Clone {}",
	} {
		t.Run(source, func(t *testing.T) {
			if _, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestMalformedAssociatedTypes(t *testing.T) {
	for _, source := range []string{
		"type X = <T as Trait>;", "type X = <T>::;", "type X = <as Trait>::Item;",
		"type X = <T as >::Item;", "type X = <T Trait>::Item;", "type X = <T>::Item<>;",
		"trait X: Trait<Item = > {}", "trait X: Trait<Item = u8 = u16> {}",
		"trait X: Trait<A::Item = u8> {}", "trait X: Trait<::Item = u8> {}",
		"trait X: Trait<(Item) = u8> {}", "trait X: Trait<(bool) = u8> {}",
		"trait X: Trait<&Item = u8> {}", "trait X: Trait<<T>::Item = u8> {}",
		"trait X: Trait<'a = T> {}", "trait X: Trait<_ = T> {}",
		"trait X: Trait<Self = T> {}", "trait X: Trait<self = T> {}",
		"trait X: Trait<Item = T, U> {}", "trait X: Trait<Item = T, 'a> {}",
		"trait X: Trait<T, 'a, Item = U> {}",
		"type X = " + strings.Repeat("<", 130) + "T" + strings.Repeat(">::Item", 130) + ";",
		"trait X: " + strings.Repeat("Trait<Item = ", 130) + "T" + strings.Repeat(">", 130) + " {}",
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

func TestAssociatedTypeBoundsRemainUnsupported(t *testing.T) {
	c, err := rustparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: "trait X: Iterator<Item: Clone> {}"})
	var diagnostic *text.Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || diagnostic.Span == nil || len(c.Nodes) != 0 {
		t.Fatalf("expected explicit unsupported bound diagnostic: %v", err)
	}
}
