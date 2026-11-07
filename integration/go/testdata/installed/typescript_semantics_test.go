package consumer_test

import (
	checktypescript "github.com/Bitspark/bittype/check/syntax/typescript/provider/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	parsetypescript "github.com/Bitspark/bittype/parse/syntax/typescript/provider/go"

	"context"
	"os"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/typescript/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"

	r "github.com/Bitspark/bittype/resolve/files/typescript/provider/go"
	n "github.com/Bitspark/bittype/semantics/typescript/go"
	s "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func TestInstalledTypeScriptSemantics(t *testing.T) {
	ctx := context.Background()
	text := "export interface Page<T> { items:T[] } export type Pages = Page<Page<string>>;"
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "native-types-1")
	if err != nil {
		t.Fatal(err)
	}
	parser, err := parse.NewSyntaxParser[s.Language](definition, parsetypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := parser.Parse(ctx, parse.Input{Text: text, TreeID: "page", SourceID: "page"})
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](definition, checktypescript.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	info, err := parsed.ASTView().Inspect()
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: ast.NodeRef{Tree: info.ID, Node: info.Root}, Trees: []representation.AST{parsed.ASTView()}, Files: &files.Context{Files: files.Files{Entries: []files.File{{ID: "page", Path: "page.ts", Text: text}}}, Links: []files.SyntaxLink{{File: "page", Tree: "page", Source: "page"}}}})
	if err != nil {
		t.Fatal(err)
	}
	compiler := r.Compiler{NodePath: os.Getenv("BITTYPE_TYPESCRIPT_NODE"), TypeScriptPath: os.Getenv("BITTYPE_TYPESCRIPT_MODULE")}
	resolver, err := r.New(definition, compiler)
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](definition, analyze.Contract(), analyze.Provider{})
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
	pages, err := m.Lookup("page", "Pages")
	if err != nil {
		t.Fatal(err)
	}
	fields, ok := pages.Type().Properties()
	if !ok || len(fields) != 1 || fields[0].Type.String() != "Page<string>[]" {
		t.Fatal(fields)
	}
	q, err := analyze.NewQueries(m, compiler)
	if err != nil {
		t.Fatal(err)
	}
	a, err := q.Assignable(ctx, pages.Type(), pages.Type())
	if err != nil || !a.Holds() {
		t.Fatal(err)
	}
}
