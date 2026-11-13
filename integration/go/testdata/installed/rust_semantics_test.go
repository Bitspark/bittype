package consumer_test

import (
	analyzerust "github.com/Bitspark/bittype/analyze/files/rust/provider/go"
	checkrust "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	parserust "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"

	"context"
	"os"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"

	r "github.com/Bitspark/bittype/resolve/files/rust/provider/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestInstalledRustSemantics(t *testing.T) {
	ctx := context.Background()
	text := "pub struct Page<T> {pub items: Vec<T>} pub type Pages = Page<Page<u8>>;"
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-1")
	if err != nil {
		t.Fatal(err)
	}
	parser, err := parse.NewSyntaxParser[s.Language](definition, parserust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, parse.Input{Text: text, TreeID: "page", SourceID: "page"})
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](definition, checkrust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	info, err := parsed.ASTView().Inspect()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: ast.NodeRef{Tree: info.ID, Node: info.Root}, Trees: []representation.AST{parsed.ASTView()}, Files: &files.Context{Files: files.Files{Entries: []files.File{{ID: "page", Path: "page.rs", Text: text}}}, Links: []files.SyntaxLink{{File: "page", Tree: "page", Source: "page"}}}})
	if err != nil {
		t.Fatal(err)
	}
	compiler := r.Compiler{HelperPath: os.Getenv("BITTYPE_RUST_HELPER"), SysrootPath: os.Getenv("BITTYPE_RUST_SYSROOT")}
	resolver, err := r.New(definition, compiler)
	if err != nil {
		t.Fatal(err)
	}
	settings, err := n.ConfigurationFor(n.Environment{Edition: n.Edition2024, Target: n.TargetLinuxX64, Crate: n.Crate{Name: "consumer", Kind: n.KindLibrary}})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: settings})
	if err != nil {
		t.Fatal(err)
	}
	contract, err := n.AnalysisContractFor(settings)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](definition, contract, analyzerust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: settings}})
	if err != nil {
		t.Fatal(err)
	}
	m, err := n.Open(result)
	if err != nil {
		t.Fatal(err, result.Diagnostics())
	}
	pages, err := m.Lookup("consumer::Pages")
	if err != nil {
		t.Fatal(err)
	}
	fields, ok := pages.Type().Fields()
	if !ok || len(fields) != 1 || !strings.Contains(fields[0].Type.String(), "Page<u8>") {
		t.Fatal(fields)
	}
	if visibility, err := fields[0].Declaration.Visibility(); err != nil || !visibility.Public {
		t.Fatal(visibility, err)
	}
	if effective, err := pages.Effective(); err != nil || !effective.Nameable {
		t.Fatal(effective, err)
	}
}
