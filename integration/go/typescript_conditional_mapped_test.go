package integration_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
)

func tsMappedFixture() ast.AST[ts.Language] {
	f := fixture[ts.Language]{tree: ast.AST[ts.Language]{ID: "independent-ts-mapped"}}
	ref := func(name string, target ast.NodeID) ast.NodeID {
		id := f.node("type-reference", attr("name", name))
		f.tree.Nodes[id-1].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: f.tree.ID, Node: target}}}
		return id
	}
	parameter := f.node("type-parameter", attr("name", "T"))
	keys := f.node("type-operator", attr("operator", "keyof"), child("type", ref("T", parameter)))
	key := f.node("type-parameter", attr("name", "K"), child("constraint", keys))
	index := f.node("indexed-access-type", nil, child("object", ref("T", parameter)), child("index", ref("K", key)))
	inferred := f.node("type-parameter", attr("name", "U"))
	inferType := f.node("infer-type", nil, child("type-parameter", inferred))
	whenFalse := f.node("keyword-type", attr("name", "never"))
	conditional := f.node("conditional-type", nil, child("check", index), child("constraint", inferType), child("true", ref("U", inferred)), child("false", whenFalse))
	mapped := f.node("mapped-type", attr("readonly-modifier", "remove", "optional-modifier", "remove"), child("type-parameter", key), child("name-type", ref("K", key)), child("type", conditional))
	f.tree.Nodes[mapped-1].Origin = ast.Origin{Spans: []ast.Span{{Source: "original", Start: 0, End: 1}}}
	alias := f.node("type-alias", attr("name", "Mapped", "exported", "true"), child("type-parameter", parameter), child("type", mapped))
	f.tree.Root = f.node("source-file", nil, child("declaration", alias))
	return f.tree
}

// bittype:case target=github.com/Bitspark/bittype/provider/typescript/go#Provider.Format feature=conditional-mapped-syntax case=independent-ast-formatting expects=ast-preserved-and-golden-text
func TestIndependentTypeScriptMappedFormatter(t *testing.T) {
	a := tsMappedFixture()
	p := tsprovider.Provider{}
	c, err := p.Format(context.Background(), a)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(c.Abstract, a) {
		t.Fatal("formatter changed syntax, binder references or origins")
	}
	got, err := p.Print(context.Background(), c)
	want := "export type Mapped<T> = {\n  -readonly [K in keyof T as K]-?: T[K] extends infer U ? U : never;\n};\n"
	if err != nil || got != want {
		t.Fatalf("want %q, got %q: %v", want, got, err)
	}
	reparsed, err := p.Parse(context.Background(), parse.Input{SourceID: "formatted", TreeID: "reparsed", Text: got})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(shape(a, a.Root), shape(reparsed.Abstract, reparsed.Abstract.Root)) {
		t.Fatal("independently authored syntax changed on round trip")
	}
}

func TestTypeScriptConditionalMappedFormatting(t *testing.T) {
	for _, tc := range []struct{ source, want string }{
		{"type X=A extends B?C:D extends E?F:G;", "type X = A extends B ? C : D extends E ? F : G;\n"},
		{"type X=(A extends B?C:D) extends E?F:G;", "type X = (A extends B ? C : D) extends E ? F : G;\n"},
		{"type X=A extends (B extends C?D:E)?F:G;", "type X = A extends (B extends C ? D : E) ? F : G;\n"},
		{"type X=(() => A) extends (() => B)?C:D;", "type X = (() => A) extends (() => B) ? C : D;\n"},
		{"type X=() => A extends B?C:D;", "type X = () => A extends B ? C : D;\n"},
		{"type X=(A extends B?C:D)[];", "type X = (A extends B ? C : D)[];\n"},
		{"type X=A|(B extends C?D:E)&F;", "type X = A | (B extends C ? D : E) & F;\n"},
		{"type X=keyof (A extends B?C:D);", "type X = keyof (A extends B ? C : D);\n"},
		{"type X=(A extends B?C:D)[K];", "type X = (A extends B ? C : D)[K];\n"},
		{"type X=T extends (infer U)[]?U:T;", "type X = T extends (infer U)[] ? U : T;\n"},
		{"type X=T extends infer U extends string?U:never;", "type X = T extends (infer U extends string) ? U : never;\n"},
		{"type X=T extends (infer U extends string)|number?U:never;", "type X = T extends (infer U extends string) | number ? U : never;\n"},
		{"type X=T extends (infer U extends string)&Tag?U:never;", "type X = T extends (infer U extends string) & Tag ? U : never;\n"},
		{"type X=T extends (infer U extends X ? A:B)?Y:Z;", "type X = T extends (infer U extends X ? A : B) ? Y : Z;\n"},
		{"type X={+readonly [K in keyof T]+?:T[K]};", "type X = {\n  readonly [K in keyof T]?: T[K];\n};\n"},
		{"type X={-readonly [K in keyof T as K extends string?K:never]-?:T[K]};", "type X = {\n  -readonly [K in keyof T as K extends string ? K : never]-?: T[K];\n};\n"},
		{"type X={ [K in keyof T]: { [P in keyof T[K]]?: T[K][P] } };", "type X = {\n  [K in keyof T]: {\n    [P in keyof T[K]]?: T[K][P];\n  };\n};\n"},
	} {
		t.Run(tc.source, func(t *testing.T) {
			p := tsprovider.Provider{}
			c, err := p.Parse(context.Background(), parse.Input{SourceID: "source", TreeID: "tree", Text: tc.source})
			if err != nil {
				t.Fatal(err)
			}
			if original, err := p.Print(context.Background(), c); err != nil || original != tc.source {
				t.Fatalf("original source changed: %v", err)
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
				t.Fatal("conditional/mapped syntax or formatting is unstable")
			}
		})
	}
}
