package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
	"testing"
)

// An unrelated language set demonstrates source discovery without extending a
// production language/capability catalogue or executing provider code.
func fixture(t *testing.T) *program {
	t.Helper()
	const module = "example.test/fixture"
	sources := map[string]string{
		"contracts/language/go": `package language; type ID string; type Provider interface { ID() ID; Name() string }`,
		"contracts/print/go":    `package print; import "example.test/fixture/contracts/language/go"; type Provider[I,O any] interface { language.Provider; Print(I)(O,error) }`,
		"contracts/parse/go":    `package parse; import "example.test/fixture/contracts/language/go"; type Provider[I,O any] interface { language.Provider; Parse(I)(O,error) }`,
		"trees/cst/go":          `package cst; type CST[L any] struct {}`,
		"source/go":             `package source; type Source[L any] struct {}`,
		"files/go":              `package files; type Files struct {}`,
		"syntax/ember/go":       `package ember; import "example.test/fixture/contracts/language/go"; type Language struct{};func(Language)ID()language.ID{return "ember"};func(Language)Name()string{return "Ember"}`,
		"syntax/ice/go":         `package ice; import "example.test/fixture/contracts/language/go"; type Language struct{};func(Language)ID()language.ID{return "ice"};func(Language)Name()string{return "Ice"}`,
		"provider/ember/go": `package emberprovider
import (s "example.test/fixture/syntax/ember/go"; tree "example.test/fixture/trees/cst/go")
type Output struct{s.Language}
func(Output)Print(v tree.CST[s.Language])(string,error){return emit(v)}
func emit(tree.CST[s.Language])(string,error){return "",nil}
`,
		"print/files/go": `package filesprint
import (id "example.test/fixture/contracts/language/go"; printing "example.test/fixture/contracts/print/go"; tree "example.test/fixture/trees/cst/go"; source "example.test/fixture/source/go"; files "example.test/fixture/files/go")
type Output[L any] struct {identity id.Provider;syntax printing.Provider[tree.CST[L],string]}
func New[L any](identity id.Provider,syntax printing.Provider[tree.CST[L],string])(Output[L],error){return Output[L]{identity,syntax},nil}
func(p Output[L])ID()id.ID{return p.identity.ID()}
func(p Output[L])Name()string{return p.identity.Name()}
func(p Output[L])Print(source.Source[L])(files.Files,error){return files.Files{},nil}
`,
		// This adapter implements Parse but must not advertise a parser for Ice.
		"contracts/parse/go/adapter": `package adapter
import s "example.test/fixture/syntax/ice/go"
type Wrapper struct{s.Language}
func(Wrapper)Parse(string)(string,error){return "",nil}
`,
		"contracts/transform/go": `package transform
type ASTTransformer interface{FromID()string;ToID()string;Transform()}
type bound[F,T any] struct{}
func(bound[F,T])FromID()string{return ""};func(bound[F,T])ToID()string{return ""};func(bound[F,T])Transform(){}
func NewTransformer[F,T any]()ASTTransformer{return bound[F,T]{}}
`,
		"transforms/example/go": `package mapping
import(bridge "example.test/fixture/contracts/transform/go";a "example.test/fixture/syntax/ember/go";b "example.test/fixture/syntax/ice/go")
func New()bridge.ASTTransformer{return bridge.NewTransformer[a.Language,b.Language]()}
`,
	}
	p := &program{root: t.TempDir(), module: module, fset: token.NewFileSet(), packages: map[string]*sourcePackage{}, functions: map[*types.Func]declaration{}}
	l := &fixtureLoader{p: p, sources: sources}
	for path := range sources {
		if _, err := l.Import(module + "/" + path); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

type fixtureLoader struct {
	p       *program
	sources map[string]string
}

func (l *fixtureLoader) Import(path string) (*types.Package, error) {
	if p := l.p.packages[path]; p != nil {
		return p.types, nil
	}
	rel := strings.TrimPrefix(path, l.p.module+"/")
	source, ok := l.sources[rel]
	if !ok {
		return nil, fmt.Errorf("missing fixture %s", path)
	}
	file, err := parser.ParseFile(l.p.fset, filepath.Join(l.p.root, rel, "implementation.go"), source, 0)
	if err != nil {
		return nil, err
	}
	p := &sourcePackage{files: []*ast.File{file}, info: &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{}, Instances: map[*ast.Ident]types.Instance{}}}
	config := types.Config{Importer: l}
	p.types, err = config.Check(path, l.p.fset, p.files, p.info)
	if err != nil {
		return nil, err
	}
	l.p.packages[path] = p
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			l.p.functions[p.info.Defs[fn.Name].(*types.Func)] = declaration{p, fn}
		}
	}
	return p.types, nil
}

func TestGenericFilesPrintingUsesConstructorCollaborators(t *testing.T) {
	p := fixture(t)
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(m.languages) != 2 || m.languages[0].id != "ember" || m.languages[1].id != "ice" {
		t.Fatal("language discovery", m.languages)
	}
	files := m.rows["print.Provider[source.Source[L], files.Files]"]
	if files == nil || len(files.cells["ember"]) == 0 || len(files.cells["ice"]) != 0 {
		t.Fatal("generic printing must require syntax printing, but no semantics")
	}
	for key := range m.rows {
		if strings.HasPrefix(key, "parse.Provider") {
			t.Fatal("contract adapter advertised as implementation")
		}
	}
	if len(m.edges["ember"]["ice"]) == 0 || len(m.edges["ice"]["ember"]) != 0 {
		t.Fatal("typed transformer direction")
	}
	providers, transformers, err := p.render(m, "https://example.test/blob/main")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(providers, "[emit](https://example.test/blob/main/provider/ember/go/implementation.go#L5)") {
		t.Fatal("missing direct delegation source position", providers)
	}
	if !strings.Contains(transformers, "[New](https://example.test/blob/main/transforms/example/go/implementation.go#L3)") {
		t.Fatal("missing transformer factory link", transformers)
	}
}

func TestMissingSyntaxPrinterDoesNotCreateFilesPrinting(t *testing.T) {
	p := fixture(t)
	// Remove the concrete syntax printer; the generic printer must not seed itself.
	delete(p.packages, p.module+"/provider/ember/go")
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range m.rows {
		if row.level == "Files" {
			for _, cell := range row.cells {
				if len(cell) != 0 {
					t.Fatal("unsatisfied generic implementation reported")
				}
			}
		}
	}
}

func TestGeneratedRegionsPreserveProseAndRejectAmbiguousMarkers(t *testing.T) {
	before := "authored\n<!-- generated:providers:begin -->\nold\n<!-- generated:providers:end -->\nauthored tail\n"
	after, err := replaceRegion(before, "providers", "new table\n")
	if err != nil {
		t.Fatal(err)
	}
	if after != "authored\n<!-- generated:providers:begin -->\n\nnew table\n\n<!-- generated:providers:end -->\nauthored tail\n" {
		t.Fatal(after)
	}
	if _, err := replaceRegion(before+before, "providers", "new"); err == nil {
		t.Fatal("duplicate markers accepted")
	}
	if _, err := replaceRegion("authored", "providers", "new"); err == nil {
		t.Fatal("missing markers accepted")
	}
}
