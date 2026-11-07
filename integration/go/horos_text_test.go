package integration_test

import (
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	formathoros "github.com/Bitspark/bittype/format/syntax/horos/provider/go"
	textfixture "github.com/Bitspark/bittype/integration/go/testdata/textfixture"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"

	"context"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"

	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	s "github.com/Bitspark/bittype/syntax/horos/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

// bittype:case target=github.com/Bitspark/bittype/parse/syntax/horos/provider/go#Provider.Parse feature=independent-engine-stages case=pagination-source expects=independently-composable
// bittype:case target=github.com/Bitspark/bittype/check/syntax/horos/provider/go#Provider.Check feature=independent-engine-stages case=pagination-source expects=independently-composable
// bittype:case target=github.com/Bitspark/bittype/format/syntax/horos/provider/go#Provider.Format feature=independent-engine-stages case=pagination-source expects=independently-composable
// bittype:case target=github.com/Bitspark/bittype/print/syntax/shared/provider/go#Provider.Print feature=independent-engine-stages case=pagination-source expects=independently-composable language=github.com/Bitspark/bittype/syntax/horos/go#Language
// bittype:case target=github.com/Bitspark/bittype/resolve/files/horos/provider/go#Provider.Resolve feature=independent-engine-stages case=pagination-source expects=independently-composable
// bittype:case target=github.com/Bitspark/bittype/analyze/files/horos/provider/go#Provider.Analyze feature=independent-engine-stages case=pagination-source expects=independently-composable
func TestHorosSourceThroughIndependentEngineStages(t *testing.T) {
	p := textfixture.New(parsehoros.Provider{}, checkhoros.Provider{}, formathoros.Provider{})
	definition, err := representation.NewLanguage[s.Language](p, s.Profile)
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	registerNativeStages(t, e, definition, p)
	source := "// pagination\r\nmodule example.pagination;\r\nimport h from horos.kinds pin \"sha256:" + decl.PreludePin().String() + "\";\r\ntype Page<T> = { items: list<T>; next?: text; };\r\ntype DataPage<D: h.Data> = Page<D>;\r\ntype People = DataPage<text>;\r\n"
	parsed, err := e.ParseSyntax(context.Background(), p.ID(), parse.Input{SourceID: "page.horos", TreeID: "page", Text: source})
	if err != nil {
		t.Fatal(err)
	}
	printed, err := e.PrintSyntax(context.Background(), parsed)
	if err != nil || printed != source {
		t.Fatalf("engine lost trivia: %v", err)
	}
	rendered, err := e.RenderSyntax(context.Background(), parsed.ASTView())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.ParseSyntax(context.Background(), p.ID(), parse.Input{SourceID: "output.horos", TreeID: "output", Text: rendered}); err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[s.Language](definition, parsed.ASTView())
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := e.Capture(context.Background(), semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{parsed.ASTView()}})
	if err != nil {
		t.Fatal(err)
	}
	resolver, err := resolution.NewResolver[s.Language](definition, resolve.Contract(), resolve.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](definition, analyze.Contract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	if err := e.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	resolved, err := e.Resolve(context.Background(), snapshot, semantic.Options{})
	if err != nil || !resolved.HasSemantics() {
		t.Fatalf("source resolution: %v", err)
	}
	analyzed, err := e.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved})
	if err != nil || !analyzed.HasSemantics() {
		t.Fatalf("source analysis: %v", err)
	}
}
