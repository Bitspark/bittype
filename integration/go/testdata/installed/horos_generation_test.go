package consumer_test

import (
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	formattypescript "github.com/Bitspark/bittype/format/syntax/typescript/provider/go"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	syntaxprint "github.com/Bitspark/bittype/print/syntax/shared/provider/go"

	"context"
	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	generate "github.com/Bitspark/bittype/transform/syntax/horos-to-typescript/mapping/go"
	rules "github.com/Bitspark/bittype/transform/syntax/horos/mapping/go"
	reader "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	target "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"

	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	"strings"
	"testing"
)

func TestInstalledSemanticGeneration(t *testing.T) {
	ctx := context.Background()
	source, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "input.horos", TreeID: "input", Text: "module installed; type Page<T> = { items: list<T>; next?: text; }; type People = Page<text>;"})
	if err != nil {
		t.Fatal(err)
	}
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[ts.Language](ts.Language{}, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	for _, d := range []representation.Language{from, to} {
		if err = e.RegisterLanguage(d); err != nil {
			t.Fatal(err)
		}
	}
	hc, err := check.NewChecker[h.Language](from, checkhoros.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	tc, err := check.NewChecker[ts.Language](to, checktypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []check.Checker{hc, tc} {
		if err = e.RegisterChecker(c); err != nil {
			t.Fatal(err)
		}
	}
	resolver, err := resolution.NewResolver[h.Language](from, resolve.Contract(), resolve.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[h.Language](from, analyze.Contract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	packed, err := representation.Pack(from, source)
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := e.Capture(ctx, semantic.Input{Entry: ast.NodeRef{Tree: source.Abstract.ID, Node: source.Abstract.Root}, Trees: []representation.AST{packed}})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := e.Resolve(ctx, snapshot, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	analyzed, err := e.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved})
	if err != nil {
		t.Fatal(err)
	}
	registry := e.Generation()
	readerRef, err := c.RegisterSource[h.Language, reader.Reader](registry, from, reader.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	targetRef, err := c.RegisterTarget[h.Language, ts.Language, target.Builder[h.Language]](registry, from, to, target.Provider[h.Language]{})
	if err != nil {
		t.Fatal(err)
	}
	var refs []c.MappingRef
	for _, provider := range generate.Providers(rules.RepresentationProfile) {
		ref, err := c.RegisterMapping[h.Language, ts.Language, reader.Reader, target.Builder[h.Language]](registry, from, to, provider)
		if err != nil {
			t.Fatal(err)
		}
		refs = append(refs, ref)
	}
	var roots []ast.NodeRef
	for i, node := range source.Abstract.Nodes {
		if node.Kind == h.KindTypeDefinition {
			roots = append(roots, ast.NodeRef{Tree: source.Abstract.ID, Node: ast.NodeID(i + 1)})
		}
	}
	request := c.Request{Input: transform.Input{Source: snapshot, Evidence: []semantic.Result{analyzed}}, To: to, Roots: roots, Output: "generated-types", Limits: g.DefaultLimits(), Selection: c.Selection{Source: readerRef, Target: targetRef, Mappings: refs}}
	output, err := e.Generate(ctx, request)
	if err != nil {
		t.Fatal(err)
	}
	typed, err := representation.Unpack[ts.Language](to, output.Snapshot().Inputs()[0])
	if err != nil {
		t.Fatal(err)
	}
	concrete, err := formattypescript.Format(context.Background(), typed)
	if err != nil {
		t.Fatal(err)
	}
	text, err := syntaxprint.Print[ts.Language](context.Background(), concrete)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "interface Page<T>") || !strings.Contains(text, "next?: string;") || !strings.Contains(text, "type People = Page<string>;") {
		t.Fatal(text)
	}
}
