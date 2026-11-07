package parse_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	tr "github.com/Bitspark/bittype/internal/tree"
	tsparse "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
)

// bittype:case target=github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Parse feature=conditional-mapped-syntax case=structured-binders-and-modifiers expects=structure-preserved
func TestConditionalMappedStructure(t *testing.T) {
	source := "type X<T> = { -readonly [K in keyof T as T[K] extends string ? K : never]+?: T[K] extends infer U extends string ? U : never };"
	c, err := tsparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	v := tr.View[s.Language]{Tree: c.Abstract}
	mapped := v.Child(v.Child(c.Abstract.Root, "declaration"), "type")
	if v.Kind(mapped) != "mapped-type" || v.Attr(mapped, "readonly-modifier") != "remove" || v.Attr(mapped, "optional-modifier") != "add" {
		t.Fatal("mapped modifiers lost")
	}
	binder := v.Child(mapped, "type-parameter")
	if v.Attr(binder, "name") != "K" || v.Kind(v.Child(binder, "constraint")) != "type-operator" {
		t.Fatal("mapped key binder must retain the in type as its constraint")
	}
	remap, value := v.Child(mapped, "name-type"), v.Child(mapped, "type")
	if v.Kind(remap) != "conditional-type" || v.Kind(value) != "conditional-type" || v.Kind(v.Child(value, "check")) != "indexed-access-type" {
		t.Fatal("conditional remapping/value roles lost")
	}
	inferred := v.Child(value, "constraint")
	parameter := v.Child(inferred, "type-parameter")
	if v.Kind(inferred) != "infer-type" || v.Attr(parameter, "name") != "U" || v.Attr(v.Child(parameter, "constraint"), "name") != "string" {
		t.Fatal("infer binder and constraint lost")
	}
	span := v.Node(inferred).Origin.Spans[0]
	if span.Source != "source" || source[span.Start:span.End] != "infer U extends string" {
		t.Fatalf("infer origin: %#v", span)
	}
}

// bittype:case target=github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Parse feature=conditional-mapped-syntax case=unresolved-syntax expects=accept
func TestConditionalMappedWithoutResolution(t *testing.T) {
	for _, source := range []string{
		"type X = Unknown extends Other ? Left : Right;",
		"type X<T> = T extends [infer U extends string] ? U : never;",
		"type X<T> = T extends (infer U)[] ? U : never;",
		"type X<T> = T extends (infer U extends X ? A : B) ? Y : Z;",
		"type X<T> = T extends { x: infer U } ? U : never;",
		"type X = { [K in Missing as Other<K>]: Result<K> };",
		"type X = { [K in 123]: Unknown };",
		"type X = { readonly: string; [key: string]: unknown };",
		"type X<T> = { readonly [key: string]: T };",
		"type X<T> = { [K in keyof T]: { [K in keyof T]: T[K] } };",
		"type X<T> = T extends unknown ? () => string : (x: T) => number;",
	} {
		t.Run(source, func(t *testing.T) {
			if _, err := tsparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// bittype:case target=github.com/Bitspark/bittype/parse/syntax/typescript/provider/go#Parse feature=conditional-mapped-syntax case=malformed-and-bounded-input expects=located-error-without-output
func TestMalformedConditionalMappedTypes(t *testing.T) {
	for _, source := range []string{
		"type X = T extends U;", "type X = T extends U ? A;", "type X = T extends U ? : B;",
		"type X = T extends ? A : B;", "type X = T extends U ? A :;", "type X = T ? A : B;",
		"type X = T\nextends U ? A : B;", "type X = T /*\n*/ extends U ? A : B;",
		"type X = T extends U extends V ? A : B;",
		"type X = infer T;", "type X = T extends U ? infer V : never;",
		"type X = T extends U ? true : infer V;", "type X = T extends infer ? true : false;",
		"type X = T extends infer U = string ? U : never;", "type X = T extends infer U extends ? U : never;",
		"type X = { [K in]: T };", "type X = { [K in keyof T as]: T };",
		"type X = { [K in T, P in U]: V };", "type X = { [K in T = U]: V };",
		"type X = { [K in T]: };", "type X = { [K in T] };",
		"type X = { [K in T]: U; other: string };", "type X = { other: string; [K in T]: U };",
		"type X = { [K in T]: U, };", "type X = { [K in T]: U } garbage;",
		"type X = { ++readonly [K in T]: U };", "type X = { - [K in T]: U };",
		"type X = { [K in T]+: U };", "type X = { [K in T]??: U };",
		"interface X { [K in T]: U }", "type X = { [K in T as `get${K}`]: U };",
		"type X = " + strings.Repeat("T extends U ? A : ", 140) + "B;",
		"type X = " + strings.Repeat("{ [K in T]: ", 140) + "U" + strings.Repeat(" }", 140) + ";",
	} {
		t.Run(source, func(t *testing.T) {
			c, err := tsparse.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: source})
			var d *text.Diagnostic
			if !errors.As(err, &d) || d.Span == nil || len(c.Nodes) != 0 || len(c.Abstract.Nodes) != 0 {
				t.Fatalf("expected located rejection and no partial output: %v", err)
			}
		})
	}
}
