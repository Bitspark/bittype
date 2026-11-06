package report

import (
	"context"
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
func fixture(t *testing.T, additions ...map[string]string) *program {
	t.Helper()
	const module = "example.test/fixture"
	sources := map[string]string{
		"contracts/language/go": `package language; type ID string
// bittype:shared
type Provider interface { ID() ID; Name() string }`,
		"contracts/print/go": `package print; import "example.test/fixture/contracts/language/go"; type Provider[I,O any] interface { language.Provider; Print(I)(O,error) }`,
		"contracts/parse/go": `package parse; import "example.test/fixture/contracts/language/go"; type Provider[I,O any] interface { language.Provider; Parse(I)(O,error) }`,
		"trees/cst/go": `package cst
// bittype:dimension Syntax
type CST[L any] struct {}`,
		"source/go": `package source
// bittype:dimension Files
type Source[L any] struct {}`,
		"files/go": `package files
// bittype:dimension Files
type Files struct {}`,
		"syntax/ember/go": `package ember; import "example.test/fixture/contracts/language/go"; type Language struct{};func(Language)ID()language.ID{return "ember"};func(Language)Name()string{return "Ember"}`,
		"syntax/ice/go":   `package ice; import "example.test/fixture/contracts/language/go"; type Language struct{};func(Language)ID()language.ID{return "ice"};func(Language)Name()string{return "Ice"}`,
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
	for _, extra := range additions {
		for path, source := range extra {
			sources[path] = source
		}
	}
	p := &program{root: t.TempDir(), module: module, fset: token.NewFileSet(), packages: map[string]*sourcePackage{}, functions: map[*types.Func]declaration{}, ctx: context.Background(), limits: DefaultLoadLimits(), sources: map[string]string{}, metadata: map[string]packageInfo{}}
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
	file, err := parser.ParseFile(l.p.fset, filepath.Join(l.p.root, rel, "implementation.go"), source, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	p := &sourcePackage{files: []*ast.File{file}, sources: map[string]string{filepath.Join(l.p.root, rel, "implementation.go"): source}, info: &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{}, Instances: map[*ast.Ident]types.Instance{}}}
	config := types.Config{Importer: l}
	p.types, err = config.Check(path, l.p.fset, p.files, p.info)
	if err != nil {
		return nil, err
	}
	l.p.packages[path] = p
	if l.p.sources == nil {
		l.p.sources = map[string]string{}
	}
	filename := filepath.Join(l.p.root, rel, "implementation.go")
	l.p.sources[filename] = source
	if l.p.metadata == nil {
		l.p.metadata = map[string]packageInfo{}
	}
	l.p.metadata[path] = packageInfo{ImportPath: path, Name: p.types.Name(), Dir: filepath.Dir(filename), GoFiles: []string{"implementation.go"}}
	if err := l.p.readDimensions(p); err != nil {
		return nil, err
	}
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
	if files == nil || files.level != "Files" || len(files.cells["ember"]) == 0 || len(files.cells["ice"]) != 0 {
		t.Fatal("generic printing must require syntax printing, but no semantics")
	}
	if m.rows["language.Provider"].level != "Shared" || m.rows["print.Provider[cst.CST[L], string]"].level != "Syntax" {
		t.Fatal("matrix dimensions must come from source metadata")
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
	if !strings.Contains(providers, "[provider/ember/go.Output.Print](https://example.test/blob/main/provider/ember/go/implementation.go#L4)") {
		t.Fatal("missing public provider method source position", providers)
	}
	if !strings.Contains(transformers, "[transforms/example/go.New](https://example.test/blob/main/transforms/example/go/implementation.go#L3)") {
		t.Fatal("missing transformer factory link", transformers)
	}
}

func TestRendererLinksCapabilityDetails(t *testing.T) {
	p := fixture(t)
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	providers, transformers, err := p.render(m, "https://example.test/blob/main")
	if err != nil {
		t.Fatal(err)
	}
	for name, table := range map[string]string{"providers": providers, "transformers": transformers} {
		if !strings.Contains(table, "[details](provider-capabilities.md#") {
			t.Fatalf("%s matrix has no link to implementation details", name)
		}
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

func TestComposedTransformerFactoriesRetainTypedEndpoints(t *testing.T) {
	p := fixture(t, map[string]string{
		"generation/syntax/go": `package generation
import bridge "example.test/fixture/contracts/transform/go"
func Compose[F,T any]()bridge.ASTTransformer{return bridge.NewTransformer[F,T]()}`,
		"generation/relay/go": `package relay
import (g "example.test/fixture/generation/syntax/go";bridge "example.test/fixture/contracts/transform/go")
func Forward[Out,In any]()bridge.ASTTransformer{return g.Compose[In,Out]()}`,
		"transforms/configured/go": `package configured
import (r "example.test/fixture/generation/relay/go";bridge "example.test/fixture/contracts/transform/go";a "example.test/fixture/syntax/ember/go";b "example.test/fixture/syntax/ice/go")
func New()bridge.ASTTransformer{return r.Forward[b.Language,a.Language]()}`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(m.edges["ember"]["ice"]) != 2 || len(m.edges["ice"]["ember"]) != 0 {
		t.Fatal("delegation lost endpoints or advertised the open helper", m.edges)
	}
	_, table, err := p.render(m, "https://example.test/blob/main")
	if err != nil || !strings.Contains(table, "/transforms/configured/go/implementation.go#L3") {
		t.Fatal("missing concrete factory location", table, err)
	}
}

func TestCyclicTransformerFactoryFails(t *testing.T) {
	p := fixture(t, map[string]string{
		"transforms/cycle/go": `package cyclic
import bridge "example.test/fixture/contracts/transform/go"
func New()bridge.ASTTransformer{return New()}`,
	})
	if _, err := p.discover(); err == nil || !strings.Contains(err.Error(), "cyclic transformer") {
		t.Fatal("cyclic factory accepted", err)
	}
}

func TestTransformerInputsDoNotInventOrHideFixedEndpoints(t *testing.T) {
	p := fixture(t, map[string]string{
		"generation/decorated/go": `package decorated
import (bridge "example.test/fixture/contracts/transform/go";a "example.test/fixture/syntax/ember/go";b "example.test/fixture/syntax/ice/go")
type wrapped struct{bridge.ASTTransformer}
type Collaborator interface{bridge.ASTTransformer}
type Alias = Collaborator
func Wrap(input Alias)bridge.ASTTransformer{return wrapped{input}}
func Again(input bridge.ASTTransformer)bridge.ASTTransformer{return Wrap(input)}
func New()bridge.ASTTransformer{return Again(bridge.NewTransformer[a.Language,b.Language]())}
func WithInput(input bridge.ASTTransformer)bridge.ASTTransformer{_ = input;return bridge.NewTransformer[b.Language,a.Language]()}
`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	if len(m.edges["ember"]["ice"]) != 2 || len(m.edges["ice"]["ember"]) != 1 {
		t.Fatal("transformer inputs invented or hid a fixed endpoint", m.edges)
	}
	for _, targets := range m.edges {
		for _, locations := range targets {
			for _, loc := range locations {
				if name := loc.op.public.Name(); name == "Wrap" || name == "Again" {
					t.Fatal("caller-supplied transformer advertised as a fixed edge", loc.op)
				}
			}
		}
	}
	if m.edges["ice"]["ember"][0].op.public.Name() != "WithInput" {
		t.Fatal("a transformer input hid the function's own typed factory")
	}
}

func TestUnboundTransformerFactoriesStillFail(t *testing.T) {
	for name, body := range map[string]string{
		"unbound":   "return nil",
		"decorated": "return Wrap(nil)",
	} {
		t.Run(name, func(t *testing.T) {
			p := fixture(t, map[string]string{
				"generation/unbound/go": `package unbound
import bridge "example.test/fixture/contracts/transform/go"
func Wrap(input bridge.ASTTransformer)bridge.ASTTransformer{return input}
func New()bridge.ASTTransformer{` + body + `}`,
			})
			if _, err := p.discover(); err == nil || !strings.Contains(err.Error(), "cannot derive transformer endpoints") {
				t.Fatal("factory without a typed endpoint accepted", err)
			}
		})
	}
}

func TestTransformerRetainsRelatedCallsFromEveryTypedPath(t *testing.T) {
	p := fixture(t, map[string]string{
		"generation/paths/go": `package generation
import (bridge "example.test/fixture/contracts/transform/go";a "example.test/fixture/syntax/ember/go";b "example.test/fixture/syntax/ice/go")
func first()bridge.ASTTransformer{return bridge.NewTransformer[a.Language,b.Language]()}
func second()bridge.ASTTransformer{return bridge.NewTransformer[a.Language,b.Language]()}
func New()bridge.ASTTransformer{result:=first();_ = second();return result}
`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	for _, loc := range m.edges["ember"]["ice"] {
		if loc.op.public.Pkg().Path() != p.module+"/generation/paths/go" {
			continue
		}
		seen := map[string]bool{}
		for _, link := range loc.op.chain {
			seen[link.callee.Name()] = true
		}
		if !seen["first"] || !seen["second"] || len(loc.op.chain) != 4 {
			t.Fatal("same-endpoint paths discarded related operation identities", loc.op.chain)
		}
		return
	}
	t.Fatal("missing factory candidate")
}

func TestDiscoveryRetainsPublicMethodAndRelatedDelegateIdentity(t *testing.T) {
	p := fixture(t, map[string]string{
		"provider/ember/go": `package emberprovider
import (s "example.test/fixture/syntax/ember/go"; tree "example.test/fixture/trees/cst/go")
type Output struct{s.Language}
type Other struct{s.Language}
func(Output)Print(v tree.CST[s.Language])(string,error){return emit(rewrite(v))}
func(Other)Print(v tree.CST[s.Language])(string,error){return emit(v)}
func rewrite(v tree.CST[s.Language]) tree.CST[s.Language] { return v }
func emit(tree.CST[s.Language])(string,error){return "",nil}
`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	locations := m.rows["print.Provider[cst.CST[L], string]"].cells["ember"]
	if len(locations) != 2 {
		t.Fatalf("expected two public contributions, got %d", len(locations))
	}
	ids := map[string]bool{}
	for _, loc := range locations {
		op := loc.op
		if op == nil || op.public == nil || op.public.Name() != "Print" || op.entry.Name() != "emit" || len(op.chain) != 1 || op.chain[0].callee != op.entry {
			t.Fatal("public identity collapsed to delegate", op)
		}
		if op.receiver == nil || op.contract == nil || op.language.id != "ember" || op.id == "" || ids[op.id] {
			t.Fatal("incomplete or duplicate contribution identity", op)
		}
		ids[op.id] = true
		if canonicalSymbol(op.public).Receiver == "" || canonicalSymbol(op.entry).Receiver != "" {
			t.Fatal("method and free function identity collapsed")
		}
	}
	// The typed factory retains every single path step and its bound arguments;
	// it remains a source relation, not evidence that the factory was executed.
	op := m.edges["ember"]["ice"][0].op
	if op.kind != "transformer" || op.from.id != "ember" || op.to.id != "ice" || len(op.chain) != 1 || len(op.chain[0].arguments) != 2 {
		t.Fatal("transformer identity or generic endpoints lost", op)
	}
}

func TestDiscoveryRejectsSameDisplayContractFromDifferentOwners(t *testing.T) {
	p := fixture(t, map[string]string{
		"contracts/elsewhere/go": `package print
import "example.test/fixture/contracts/language/go"
type Provider[I,O any] interface { language.Provider; Print(I)(O,error) }`,
	})
	if _, err := p.discover(); err == nil || !strings.Contains(err.Error(), "ambiguous matrix contract label") {
		t.Fatal("unrelated contract identities merged by display spelling", err)
	}
}

func TestDependenciesDoNotUseDisplayContractNames(t *testing.T) {
	p := fixture(t, map[string]string{
		"contracts/elsewhere/go": `package print
import "example.test/fixture/contracts/language/go"
type Provider[I,O any] interface { language.Provider; Render(I)(O,error) }`,
		"print/files/go": `package filesprint
import (id "example.test/fixture/contracts/language/go"; other "example.test/fixture/contracts/elsewhere/go"; tree "example.test/fixture/trees/cst/go"; source "example.test/fixture/source/go"; files "example.test/fixture/files/go")
type Output[L any] struct {identity id.Provider; required other.Provider[tree.CST[L],string]}
func New[L any](identity id.Provider,required other.Provider[tree.CST[L],string])(Output[L],error){return Output[L]{identity,required},nil}
func(p Output[L])ID()id.ID{return p.identity.ID()}
func(p Output[L])Name()string{return p.identity.Name()}
func(p Output[L])Print(source.Source[L])(files.Files,error){return files.Files{},nil}
`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	files := m.rows["print.Provider[source.Source[L], files.Files]"]
	if files == nil || len(files.cells["ember"]) != 0 || len(files.cells["ice"]) != 0 {
		t.Fatal("a different print.Provider's Print operation discharged the required Render contract")
	}
}

func TestCurrentRepositoryProviderDiscovery(t *testing.T) {
	// Exercise the maintained documentation command's source-discovery path in
	// the existing hosted Go test gate, before a main-only publication job runs.
	p, err := load(context.Background(), filepath.Join("..", "..", "..", ".."), DefaultLoadLimits(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.discover(); err != nil {
		t.Fatal(err)
	}
}
