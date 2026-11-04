package horosparse_test

import (
	"context"
	"encoding/hex"
	"errors"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	provider "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/horos/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

func pins(tree s.AST) []string {
	var result []string
	for _, child := range tree.Nodes[tree.Root-1].Children {
		for _, attr := range tree.Nodes[child.Node-1].Attributes {
			if attr.Name == s.AttributePin {
				result = append(result, attr.Text)
			}
		}
	}
	return result
}

func roundTrip(t *testing.T, source string) s.AST {
	t.Helper()
	p := provider.Provider{}
	parsed, err := p.Parse(context.Background(), parse.Input{SourceID: "input.horos", TreeID: "input", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	if err := p.Check(parsed); err != nil {
		t.Fatal(err)
	}
	printed, err := p.Print(context.Background(), parsed)
	if err != nil || printed != source {
		t.Fatalf("lossless source round trip: %v\n%q", err, printed)
	}
	formatted, err := p.Format(context.Background(), parsed)
	if err != nil {
		t.Fatal(err)
	}
	output, err := p.Print(context.Background(), formatted)
	if err != nil {
		t.Fatal(err)
	}
	reparsed, err := p.Parse(context.Background(), parse.Input{SourceID: "formatted.horos", TreeID: "formatted", Text: output})
	if err != nil {
		t.Fatalf("formatter produced invalid source: %v\n%s", err, output)
	}
	if strings.Join(pins(parsed.Abstract), ",") != strings.Join(pins(reparsed.Abstract), ",") {
		t.Fatalf("canonical document changed:\n%s", output)
	}
	again, err := p.Format(context.Background(), reparsed)
	if err != nil {
		t.Fatal(err)
	}
	stable, err := p.Print(context.Background(), again)
	if err != nil || stable != output {
		t.Fatalf("formatting is not idempotent: %v", err)
	}
	// The formatter owns a complete snapshot rather than aliasing borrowed ASTs.
	parsed.Abstract.Nodes[0].Kind = "mutation"
	if formatted.Abstract.Nodes[0].Kind == "mutation" {
		t.Fatal("formatter aliases input")
	}
	return reparsed.Abstract
}

func TestCompleteSourceVocabulary(t *testing.T) {
	prelude := decl.PreludePin()
	source := `module examples.vocabulary;
import h from horos.kinds pin "sha256:` + hex.EncodeToString(prelude[:]) + `";
kind Type = top(P);
kind Data = h.Data;
type Anything = any;
type Blob = atom;
type Ready = atom-exact(#"ready");
type Bytes = atom-exact(#hex"00ff");
type Small = atom-len(0, 4096);
type Operations = atom-in(set(#"read", #"write"));
type Scalars = tuple<int, text, bool>;
type Empty = tuple<>;
type Labelled<T> = compound(#"item", T, text);
type Collections<D: Data> = tuple<list<D>, seq<D>, set<D>, map<text, D>>;
type Record<T> = { value: T; next?: text; "content-type"?: text; };
type Open<D: Data> = open { value: D; };
type Alternatives = union<> | union<text> | (text | int);
type Intersections = and<> & and<text> & (text & text-non-blank());
type Negated = not<exact(null)>;
type Short = text & text-len(1, 80) & text-non-blank();
type Key = key-grammar();
type UUID = lowercase-uuid();
type Tree<D: Data> = rec(#"tree", { value: D; children: list<ref(#"tree")>; });
type Wire<T: Type> = wire<T>;
type Literals = exact(tuple("hello", true, false, null, -123456789012345678901234567890, tuple(), set(), #hex"fffe"));
type Opaque<T> = exact(tuple(#"horos/decl/param", 0));
type Unicode = exact("\u0000\u00ff\U0001f642\\x01\"\t");
type Order = (text | int) | (bool & atom);
type Explicit = union<union<text, int>, bool>;
type T = text;
type Shadow<T> = tuple<T, self.T>;
type ` + "`list`" + `<` + "`text`" + `> = ` + "`text`" + `;
type Escaped = ` + "`list`" + `<text>;
`
	tree := roundTrip(t, source)
	model, err := resolve.Model(tree)
	if err != nil {
		t.Fatal(err)
	}
	projected, err := resolve.Project(model, "projected-vocabulary")
	if err != nil {
		t.Fatal(err)
	}
	p := provider.Provider{}
	formatted, err := p.Format(context.Background(), projected)
	if err != nil {
		t.Fatal(err)
	}
	output, err := p.Print(context.Background(), formatted)
	if err != nil {
		t.Fatal(err)
	}
	reparsed := roundTrip(t, output)
	recovered, err := resolve.Model(reparsed)
	if err != nil || recovered.Root() != model.Root() {
		t.Fatalf("canonical constructor vocabulary changed: %v", err)
	}
}

func TestSourceTriviaAndSpans(t *testing.T) {
	source := "\ufeff// heading\r\nmodule examples.trivia; /* comment */\r\ntype Page<T> = {\r\n items: list<T>; // item\r\n next?: text;\r\n};\r\n// tail"
	tree := roundTrip(t, source)
	for _, n := range tree.Nodes {
		for _, span := range n.Origin.Spans {
			if span.Source != ast.SourceID("formatted.horos") || span.Start > span.End {
				t.Fatal("invalid source origin")
			}
		}
	}
}

func TestSyntaxDoesNotResolveOrCheckContracts(t *testing.T) {
	for _, body := range []string{
		"type Missing = Unknown<text>;",
		`kind Data = node("scalar", {}, epsilon); type Page<D: Data> = list<D>; type Invalid<T> = Page<T>;`,
		"type X<T: Nope> = T;",
		"type X<T> = T<text>;",
		"type Invalid = map<wire<text>, text>;",
		"type Invalid = scalarTypo;",
		"type Invalid = atom-len(-1, 2);",
	} {
		t.Run(body, func(t *testing.T) {
			tree := roundTrip(t, "module invalid.semantics; "+body)
			if _, err := resolve.Model(tree); err == nil {
				t.Fatal("semantic invalidity was lost")
			}
		})
	}
	roundTrip(t, `module missing.context; import pages from example.pagination pin "sha256:`+strings.Repeat("AB", 32)+`"; type People = pages.Page<text>;`)
}

func TestRejectWholeInput(t *testing.T) {
	cases := []string{
		"", "type X = text;", "module x", "module x; garbage", "module x; type X = text; junk",
		"module x; type X<> = text;", "module x; type X<T,T> = T;",
		"module x; type X = list<>;", "module x; type X = map<text>;", "module x; type X = wire<text,int>;",
		"module x; type X = exact-(1);", "module x; type X = text-len-(1, 2);",
		"module x; type X = { foo: text };", "module x; type X = exact(T);", "module x; type X = exact(1.2);",
		"module x; type X = exact(01);", "module x; type X = exact(-0);", "module x; type X = exact(1e3);",
		`module x; type X = exact("\xFF");`, `module x; type X = exact("\uD800");`, `module x; type X = atom-exact(#hex"f");`,
		"module x; type text = int;", "module x; type X = 'text';",
		"module x; type X = text; type X = int;", "module x; import a from x pin \"sha256:nope\";",
		"module x; type X = text; import a from x pin \"sha256:" + strings.Repeat("0", 64) + "\";",
		"module x; /* unterminated", "module x; type X = exact(" + strings.Repeat("1", 4097) + ");",
		"module x; type X = " + strings.Repeat("list<", 300) + "text" + strings.Repeat(">", 300) + ";",
		"module x;" + strings.Repeat(" ", 4<<20), "module x;\xff",
	}
	p := provider.Provider{}
	for i, source := range cases {
		out, err := p.Parse(context.Background(), parse.Input{SourceID: "bad.horos", TreeID: "bad", Text: source})
		var diagnostic *text.Diagnostic
		if err == nil || !errors.As(err, &diagnostic) || len(out.Nodes) != 0 || len(out.Abstract.Nodes) != 0 {
			t.Fatalf("case %d: expected located failure without partial output: %v", i, err)
		}
		if diagnostic.Span == nil {
			t.Fatalf("case %d has no location", i)
		}
	}
}

func FuzzSourceRoundTrip(f *testing.F) {
	f.Add("module pagination; type Page<T> = { items: list<T>; next?: text; };")
	f.Add("module x; type X = (text | int) | union<text>;")
	f.Add(`module x; type X = exact(tuple(#"horos/decl/param", 0, #hex"ff"));`)
	f.Fuzz(func(t *testing.T, source string) {
		if len(source) > 65536 {
			t.Skip()
		}
		p := provider.Provider{}
		parsed, err := p.Parse(context.Background(), parse.Input{SourceID: "fuzz", TreeID: "fuzz", Text: source})
		if err != nil {
			if len(parsed.Nodes) != 0 || len(parsed.Abstract.Nodes) != 0 {
				t.Fatal("partial parse result")
			}
			return
		}
		roundTrip(t, source)
	})
}
