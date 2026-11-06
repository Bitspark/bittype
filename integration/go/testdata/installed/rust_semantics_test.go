package consumer_test

import (
	"context"
	"os"
	"testing"
	"strings"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	p "github.com/Bitspark/bittype/provider/rust/go"
	r "github.com/Bitspark/bittype/resolve/files/rust/go"
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
	parser, err := p.NewSyntaxParser(definition)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, parse.Input{Text: text, TreeID: "page", SourceID: "page"})
	if err != nil {
		t.Fatal(err)
	}
	checker, err := p.NewSyntaxChecker(definition)
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
	resolver, err := p.NewFilesResolver(definition, compiler)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := p.NewFilesAnalyzer(definition)
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}})
	if err != nil {
		t.Fatal(err)
	}
	m, err := n.Open(result)
	if err != nil {
		t.Fatal(err, result.Diagnostics())
	}
	pages, err := m.Lookup("captured::Pages")
	if err != nil {
		t.Fatal(err)
	}
	fields, ok := pages.Type().Fields()
	if !ok || len(fields) != 1 || !strings.Contains(fields[0].Type.String(), "Page<u8>") {
		t.Fatal(fields)
	}
}
