package integration_test

import (
	"context"
	"errors"
	"go/types"
	"reflect"
	"slices"
	"sort"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
	nativequery "github.com/Bitspark/bittype/analyze/files/golang/queries/go"
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	resolution "github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/files/go"
	parsegolang "github.com/Bitspark/bittype/parse/syntax/golang/provider/go"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

const goImportModule = "example.org/app"

// cellPackage is a supplied dependency whose internals are outside the first
// profile (a pointer field). Only its declaration identities may be exported.
var cellPackage = n.Package{ImportPath: "example.org/cell", Files: []n.PackageFile{{Name: "cell.go", Text: "package cell\n\ntype Cell[T any] struct {\n\tvalue T\n\tnext  *Cell[T]\n}\n"}}}

func goImportSettings(t *testing.T, packages ...n.Package) []semantic.Setting {
	t.Helper()
	settings, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: packages})
	if err != nil {
		t.Fatal(err)
	}
	return settings
}

// goPackageCapture captures one tree per file; each tree is linked to its path.
func goPackageCapture(t *testing.T, sources map[string]string, withFiles bool) (semantic.Snapshot, map[string]s.AST) {
	t.Helper()
	var paths []string
	for p := range sources {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	return goPackageCaptureIn(t, paths, sources, withFiles)
}

// goTreeID gives each file a stable tree ID that sorts opposite to its path
// wherever two paths first differ in lowercase letters: those are mirrored
// (a↔z, b↔y, ...). For such paths, results that follow file-path order cannot
// be produced by sorting tree IDs.
func goTreeID(path string) ast.TreeID {
	id := []rune("tree:")
	for _, r := range path {
		if r >= 'a' && r <= 'z' {
			r = 'a' + 'z' - r
		}
		id = append(id, r)
	}
	return ast.TreeID(string(id))
}

// goPackageCaptureIn captures the trees in exactly the given order.
func goPackageCaptureIn(t *testing.T, paths []string, sources map[string]string, withFiles bool) (semantic.Snapshot, map[string]s.AST) {
	t.Helper()
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "go-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](definition, checkgolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	trees := map[string]s.AST{}
	var envelopes []representation.AST
	layout := &files.Context{}
	for _, p := range paths {
		id := goTreeID(p)
		parsed, err := parsegolang.Parse(context.Background(), parse.Input{SourceID: ast.SourceID(p), TreeID: id, Text: sources[p]})
		if err != nil {
			t.Fatal(err)
		}
		tree := parsed.ASTView()
		trees[p] = tree
		envelope, err := representation.Pack[s.Language](definition, tree)
		if err != nil {
			t.Fatal(err)
		}
		envelopes = append(envelopes, envelope)
		layout.Files.Entries = append(layout.Files.Entries, files.File{ID: files.ID(p), Path: files.Path(p), Text: sources[p]})
		layout.Links = append(layout.Links, files.SyntaxLink{File: files.ID(p), Tree: id, Source: ast.SourceID(p)})
	}
	input := semantic.Input{Entry: ast.NodeRef{Tree: trees[paths[0]].ID, Node: trees[paths[0]].Root}, Trees: envelopes}
	if withFiles {
		input.Files = layout
	}
	snapshot, err := semantic.Capture(context.Background(), checker, input)
	if err != nil {
		t.Fatal(err)
	}
	return snapshot, trees
}

func goImportResults(t *testing.T, snapshot semantic.Snapshot, settings []semantic.Setting) (semantic.Result, semantic.Result, error) {
	t.Helper()
	ctx := context.Background()
	resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: settings})
	if err != nil {
		return semantic.Result{}, semantic.Result{}, err
	}
	contract, err := n.AnalysisContractFor(settings)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), contract, analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: settings}})
	if err != nil {
		t.Fatal(err)
	}
	return resolved, result, nil
}

// nodeOf finds the first node of kind whose attribute name has value.
func nodeOf(t *testing.T, tree s.AST, kind s.Kind, attribute, value string) ast.NodeRef {
	t.Helper()
	for i, node := range tree.Nodes {
		if node.Kind != kind {
			continue
		}
		for _, a := range node.Attributes {
			if a.Name == attribute && a.Text == value {
				return ast.NodeRef{Tree: tree.ID, Node: ast.NodeID(i + 1)}
			}
		}
	}
	t.Fatalf("no %s with %s=%s", kind, attribute, value)
	return ast.NodeRef{}
}

func bindingAt(data semantics.Data, at ast.NodeRef, role semantics.Key) (semantics.EntityID, bool) {
	for _, b := range data.Bindings {
		if b.At.Site == at && b.At.Role == role {
			return b.Target, true
		}
	}
	return 0, false
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=captured-and-supplied-imports expects=qualified-uses-bind-imported-declarations
func TestGoImportsBindQualifiedUsesToImportedDeclarations(t *testing.T) {
	snapshot, trees := goPackageCapture(t, map[string]string{
		"pages/page.go":    "package pages\n\nimport (\n\t\"example.org/app/people\"\n\tcells \"example.org/cell\"\n)\n\ntype Page[T any] struct {\n\tItems []T\n}\n\ntype Cells = Page[cells.Cell[people.Person]]\n",
		"people/person.go": "package people\n\ntype Person struct {\n\tName string\n}\n",
	}, true)
	settings := goImportSettings(t, cellPackage)
	resolved, result, err := goImportResults(t, snapshot, settings)
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		report, _ := result.Report(key)
		if report.Status != semantic.Complete || report.Verdict != semantic.Holds {
			t.Fatal(key, report, result.Diagnostics())
		}
	}
	if !resolved.HasSemantics() {
		t.Fatal("resolution published no semantics")
	}
	model, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	data := model.Graph().Data()
	packages := map[string]semantics.EntityID{}
	for i, e := range data.Entities {
		if e.Kind == n.Key("package") {
			packages[e.Name] = semantics.EntityID(i + 1)
		}
	}
	for _, want := range []string{"example.org/app/pages", "example.org/app/people", "example.org/cell"} {
		if packages[want] == 0 {
			t.Fatal("missing package entity", want, packages)
		}
	}
	if supplied, _ := model.Graph().Entity(packages["example.org/cell"]); len(supplied.Declarations) != 0 {
		t.Fatal("a supplied package has captured declarations", supplied)
	}
	person, err := model.LookupIn("example.org/app/people", "Person")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := model.LookupIn("example.org/cell", "Cell"); err != nil {
		t.Fatal(err)
	}
	if _, err := model.LookupIn("example.org/app/pages", "Person"); err == nil {
		t.Fatal("Person found in the wrong package")
	}
	// The qualified uses bind to the imported declarations, owned by their packages.
	page := trees["pages/page.go"]
	personUse := nodeOf(t, page, s.KindQualifiedType, s.AttributeName, "Person")
	cellUse := nodeOf(t, page, s.KindQualifiedType, s.AttributeName, "Cell")
	personEntity, ok := bindingAt(data, personUse, n.Key("name"))
	if !ok || data.Entities[personEntity-1].Owner != packages["example.org/app/people"] {
		t.Fatal("people.Person does not bind to the people package's declaration", personEntity)
	}
	if decl := data.Entities[personEntity-1].Declarations; len(decl) != 1 || decl[0].Tree != trees["people/person.go"].ID {
		t.Fatal("Person's declaration site is not in its own file", decl)
	}
	cellEntity, ok := bindingAt(data, cellUse, n.Key("name"))
	if !ok || data.Entities[cellEntity-1].Owner != packages["example.org/cell"] || data.Entities[cellEntity-1].Name != "Cell" {
		t.Fatal("cells.Cell does not bind to the supplied declaration", cellEntity)
	}
	if !person.Valid() {
		t.Fatal("invalid Person handle")
	}
	// Each import binds to the package it names; the alias does not change the target.
	for path, want := range map[string]string{"example.org/app/people": "example.org/app/people", "example.org/cell": "example.org/cell"} {
		at := nodeOf(t, page, s.KindImport, s.AttributePath, path)
		target, ok := bindingAt(data, at, n.Key("import"))
		if !ok || target != packages[want] {
			t.Fatal("import binding", path, target)
		}
	}
	// Supplied internals are not exported: Cell's term has no underlying fact.
	cellBody := data.Entities[cellEntity-1].Body
	if facts := model.Graph().Facts(semantics.Subject{Term: cellBody}, n.Key("underlying")); len(facts) != 0 {
		t.Fatal("supplied structure was exported", facts)
	}
	// The explicit query session reconstructs the same graph from the same environment.
	if _, err := nativequery.NewQueries(context.Background(), model); err != nil {
		t.Fatal(err)
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=rejected-unsupported-and-invalid-environments expects=no-partial-model
func TestGoImportsRejectAndRefuse(t *testing.T) {
	person := "package people\n\ntype Person struct {\n\tName string\n}\n"
	for _, tc := range []struct {
		name     string
		sources  map[string]string
		packages []n.Package
		coverage semantic.Completion
		validity semantic.Verdict
	}{
		{"unused import", map[string]string{"a/a.go": "package a\n\nimport \"example.org/app/people\"\n\ntype A string\n", "people/p.go": person}, nil, semantic.Complete, semantic.Rejected},
		{"unexported use", map[string]string{"a/a.go": "package a\n\nimport \"example.org/app/people\"\n\ntype A = people.person\n", "people/p.go": "package people\n\ntype person struct{}\n"}, nil, semantic.Complete, semantic.Rejected},
		{"cycle", map[string]string{"a/a.go": "package a\n\nimport \"example.org/app/b\"\n\ntype A = b.B\n", "b/b.go": "package b\n\nimport \"example.org/app/a\"\n\ntype B = a.A\n"}, nil, semantic.Complete, semantic.Rejected},
		{"package clause", map[string]string{"a/a.go": "package a\n\ntype A string\n", "a/b.go": "package b\n\ntype B string\n"}, nil, semantic.Complete, semantic.Rejected},
		{"outside environment", map[string]string{"a/a.go": "package a\n\nimport \"example.org/missing\"\n\ntype A = missing.M\n"}, nil, semantic.Unsupported, semantic.Undetermined},
		{"dot import", map[string]string{"a/a.go": "package a\n\nimport . \"example.org/cell\"\n\ntype A = Cell[string]\n"}, []n.Package{cellPackage}, semantic.Unsupported, semantic.Undetermined},
		{"blank import", map[string]string{"a/a.go": "package a\n\nimport _ \"example.org/cell\"\n\ntype A string\n"}, []n.Package{cellPackage}, semantic.Unsupported, semantic.Undetermined},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snapshot, _ := goPackageCapture(t, tc.sources, true)
			resolved, result, err := goImportResults(t, snapshot, goImportSettings(t, tc.packages...))
			if err != nil {
				t.Fatal(err)
			}
			if resolved.HasSemantics() || result.HasSemantics() {
				t.Fatal("partial model published")
			}
			coverage, _ := result.Report(n.Coverage())
			validity, _ := result.Report(n.Validity())
			if coverage.Status != tc.coverage || validity.Verdict != tc.validity {
				t.Fatal(coverage, validity, result.Diagnostics())
			}
			if ds := result.Diagnostics(); len(ds) == 0 || ds[0].At == nil || ds[0].At.Tree == "" {
				t.Fatal("missing captured diagnostic", ds)
			}
		})
	}
	// Environment errors are invalid requests, never a judgment of the captured code.
	valid := map[string]string{"a/a.go": "package a\n\nimport \"example.org/cell\"\n\ntype A = cell.Cell[string]\n"}
	for name, packages := range map[string][]n.Package{
		"broken supplied package":   {{ImportPath: "example.org/cell", Files: []n.PackageFile{{Name: "cell.go", Text: "package cell\n\ntype Cell[T any] Missing\n"}}}},
		"supplied import outside":   {{ImportPath: "example.org/cell", Files: []n.PackageFile{{Name: "cell.go", Text: "package cell\n\nimport \"example.org/other\"\n\ntype Cell[T any] other.O\n"}}}},
		"supplied package unparsed": {{ImportPath: "example.org/cell", Files: []n.PackageFile{{Name: "cell.go", Text: "package cell\n\ntype Cell[T any struct{}\n"}}}},
	} {
		t.Run(name, func(t *testing.T) {
			snapshot, _ := goPackageCapture(t, valid, true)
			if _, _, err := goImportResults(t, snapshot, goImportSettings(t, packages...)); !errors.Is(err, semantic.ErrInvalidRequest) {
				t.Fatal(err)
			}
		})
	}
	t.Run("no file layout", func(t *testing.T) {
		snapshot, _ := goPackageCapture(t, valid, false)
		if _, _, err := goImportResults(t, snapshot, goImportSettings(t, cellPackage)); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(err)
		}
	})
	t.Run("not a Go file path", func(t *testing.T) {
		snapshot, _ := goPackageCapture(t, map[string]string{"a/a.txt": valid["a/a.go"]}, true)
		if _, _, err := goImportResults(t, snapshot, goImportSettings(t, cellPackage)); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(err)
		}
	})
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=normalized-receipt-environment expects=equal-environments-equal-settings
func TestGoImportConfigurationIsNormalizedAndBound(t *testing.T) {
	two := n.Package{ImportPath: "example.org/two", Files: []n.PackageFile{{Name: "b.go", Text: "package two\n"}, {Name: "a.go", Text: "package two\n\ntype T string\n"}}}
	a, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: []n.Package{two, cellPackage}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: []n.Package{cellPackage, {ImportPath: two.ImportPath, Files: []n.PackageFile{two.Files[1], two.Files[0]}}}})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("equal environments produced different settings")
	}
	if !n.MatchesConfiguration(a) || !n.MatchesConfiguration(n.Configuration()) {
		t.Fatal("valid configuration refused")
	}
	environment, enabled, err := n.Imports(a)
	if err != nil || !enabled || environment.Module != goImportModule || len(environment.Packages) != 2 || environment.Packages[0].ImportPath != "example.org/cell" {
		t.Fatal(environment, enabled, err)
	}
	tampered := append([]semantic.Setting(nil), a...)
	tampered[len(tampered)-1].Value = append([]byte(" "), tampered[len(tampered)-1].Value...)
	if n.MatchesConfiguration(tampered) {
		t.Fatal("non-normalized environment accepted")
	}
	for name, environment := range map[string]n.ImportEnvironment{
		"module without dot":    {Module: "app", OutputRoot: n.ModuleRoot},
		"dot-dot element":       {Module: "example.org/../app", OutputRoot: n.ModuleRoot},
		"supplied under module": {Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: []n.Package{{ImportPath: goImportModule + "/x", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n"}}}}},
		"duplicate package":     {Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: []n.Package{cellPackage, cellPackage}},
		"not a go file":         {Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: []n.Package{{ImportPath: "example.org/x", Files: []n.PackageFile{{Name: "x.txt", Text: "package x\n"}}}}},
		"no files":              {Module: goImportModule, OutputRoot: n.ModuleRoot, Packages: []n.Package{{ImportPath: "example.org/x"}}},
	} {
		if _, err := n.ImportConfiguration(environment); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, err)
		}
	}
	many := n.ImportEnvironment{Module: goImportModule, OutputRoot: n.ModuleRoot}
	for i := 0; i < 65; i++ {
		many.Packages = append(many.Packages, n.Package{ImportPath: "example.org/p" + string(rune('a'+i%26)) + string(rune('a'+i/26)), Files: []n.PackageFile{{Name: "p.go", Text: "package p\n"}}})
	}
	if _, err := n.ImportConfiguration(many); !errors.Is(err, semantic.ErrLimit) {
		t.Fatal(err)
	}
	// Analysis is bound to the exact environment its resolution used.
	snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\nimport \"example.org/cell\"\n\ntype A = cell.Cell[string]\n"}, true)
	resolver, _ := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: goImportSettings(t, cellPackage)})
	if err != nil {
		t.Fatal(err)
	}
	other := n.Package{ImportPath: cellPackage.ImportPath, Files: []n.PackageFile{{Name: "cell.go", Text: "package cell\n\ntype Cell[T any] struct{ Value T }\n"}}}
	otherSettings := goImportSettings(t, other)
	contract, err := n.AnalysisContractFor(otherSettings)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, _ := analysis.NewAnalyzer[s.Language](snapshot.Language(), contract, analyze.Provider{})
	if _, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: otherSettings}}); err == nil {
		t.Fatal("analysis accepted a resolution from a different supplied package text")
	}
}

// goImportSettingsAt states where the output root sits in the module.
func goImportSettingsAt(t *testing.T, outputRoot string, packages ...n.Package) []semantic.Setting {
	t.Helper()
	settings, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, OutputRoot: outputRoot, Packages: packages})
	if err != nil {
		t.Fatal(err)
	}
	return settings
}

// goImportModel resolves, analyzes and opens a capture whose reports must hold.
func goImportModel(t *testing.T, sources map[string]string, settings []semantic.Setting) (n.Model, map[string]s.AST) {
	t.Helper()
	snapshot, trees := goPackageCapture(t, sources, true)
	_, result, err := goImportResults(t, snapshot, settings)
	if err != nil {
		t.Fatal(err)
	}
	goImportsHold(t, result)
	model, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	return model, trees
}

// goPackageNames reads the graph itself: per package entity, the data of the
// terms of its package-name facts.
func goPackageNames(model n.Model) map[string][]string {
	names := map[string][]string{}
	for i, e := range model.Graph().Data().Entities {
		if e.Kind != n.PackageKind() {
			continue
		}
		names[e.Name] = []string{}
		for _, id := range model.Graph().Facts(semantics.Subject{Entity: semantics.EntityID(i + 1)}, n.PackageNameRelation()) {
			term, _ := model.Graph().Term(id)
			names[e.Name] = append(names[e.Name], string(term.Data))
		}
	}
	return names
}

// Two captures at one path that differ only in the package clause differ in
// the package-name fact and in nothing else.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=package-name-fact expects=name-published-separately-from-path
func TestGoPackageNameIsAFactSeparateFromThePath(t *testing.T) {
	settings := goImportSettingsAt(t, "generated", cellPackage)
	source := func(name string) map[string]string {
		return map[string]string{"pages/types.go": "package " + name + "\n\nimport \"example.org/cell\"\n\ntype Page struct {\n\tTitle string\n}\n\ntype Cells = cell.Cell[Page]\n"}
	}
	pages, _ := goImportModel(t, source("pages"), settings)
	people, _ := goImportModel(t, source("people"), settings)
	// Every package entity, captured and supplied, carries exactly one name fact.
	const captured, supplied = "example.org/app/generated/pages", "example.org/cell"
	for name, model := range map[string]n.Model{"pages": pages, "people": people} {
		if got := goPackageNames(model); !reflect.DeepEqual(got, map[string][]string{captured: {name}, supplied: {"cell"}}) {
			t.Fatal("package-name facts", name, got)
		}
		// The reader returns the fact, which is not the last element of the path.
		if got, err := model.PackageName(captured); err != nil || got != name {
			t.Fatal("captured package name", got, err)
		}
		if got, err := model.PackageName(supplied); err != nil || got != "cell" {
			t.Fatal("supplied package name", got, err)
		}
		for _, missing := range []string{"example.org/missing", "", name} {
			if _, err := model.PackageName(missing); !errors.Is(err, n.ErrForeign) {
				t.Fatal("a package that is not in the graph has a name", missing, err)
			}
		}
	}
	a, b := pages.Graph().Data(), people.Graph().Data()
	if reflect.DeepEqual(a, b) {
		t.Fatal("two package clauses gave one graph")
	}
	// The package clause is the only difference: with the one term of that fact
	// made equal, the graphs are equal.
	fact := people.Graph().Facts(semantics.Subject{Entity: 1}, n.PackageNameRelation())
	if len(fact) != 1 || string(b.Terms[fact[0]-1].Data) != "people" {
		t.Fatal("the captured package is not entity 1 with one name fact", fact)
	}
	b.Terms[fact[0]-1].Data = []byte("pages")
	if !reflect.DeepEqual(a, b) {
		t.Fatal("the graphs differ beyond the package-name fact")
	}
	// The base profile's one package has the constant path and its clause's name.
	base := goModel(t, "type Text string")
	if got := goPackageNames(base); !reflect.DeepEqual(got, map[string][]string{n.PackagePath: {"example"}}) {
		t.Fatal("base package-name fact", got)
	}
	if got, err := base.PackageName(n.PackagePath); err != nil || got != "example" {
		t.Fatal("base package name", got, err)
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=output-root-in-module expects=captured-path-under-output-root
func TestGoOutputRootIsRelatedToTheModuleRoot(t *testing.T) {
	sources := map[string]string{
		"pages/types.go":   "package pages\n\nimport \"example.org/app/generated/people\"\n\ntype Author = people.Person\n",
		"people/person.go": "package people\n\ntype Person struct {\n\tName string\n}\n",
	}
	model, _ := goImportModel(t, sources, goImportSettingsAt(t, "generated"))
	if got := goPackageNames(model); !reflect.DeepEqual(got, map[string][]string{"example.org/app/generated/pages": {"pages"}, "example.org/app/generated/people": {"people"}}) {
		t.Fatal("captured packages under the output root", got)
	}
	if _, err := model.LookupIn("example.org/app/generated/pages", "Author"); err != nil {
		t.Fatal(err)
	}
	if _, err := model.LookupIn("example.org/app/pages", "Author"); err == nil {
		t.Fatal("a captured package is also at the module root")
	}
	// The query session rebuilds the same graph, so it applies the same mapping.
	if _, err := nativequery.NewQueries(context.Background(), model); err != nil {
		t.Fatal(err)
	}
	// With the output root at the module root, the same files form other
	// packages: the import of the path under "generated" is outside the environment.
	snapshot, _ := goPackageCapture(t, sources, true)
	_, result, err := goImportResults(t, snapshot, goImportSettingsAt(t, n.ModuleRoot))
	if err != nil {
		t.Fatal(err)
	}
	if coverage, _ := result.Report(n.Coverage()); coverage.Status != semantic.Unsupported || result.HasSemantics() {
		t.Fatal("the import resolved although the output root is the module root", coverage)
	}
	atRoot, _ := goImportModel(t, map[string]string{"pages/types.go": "package pages\n\ntype Page string\n", "root.go": "package app\n\ntype Root string\n"}, goImportSettingsAt(t, n.ModuleRoot))
	if got := goPackageNames(atRoot); !reflect.DeepEqual(got, map[string][]string{"example.org/app/pages": {"pages"}, "example.org/app": {"app"}}) {
		t.Fatal("captured packages at the module root", got)
	}
	// The exported derivation gives the identities the bridge grouped by.
	generated := n.ImportEnvironment{Module: goImportModule, OutputRoot: "generated"}
	for _, tc := range []struct {
		environment n.ImportEnvironment
		file        string
		want        string
	}{
		{generated, "pages/types.go", "example.org/app/generated/pages"},
		{generated, "people/person.go", "example.org/app/generated/people"},
		{generated, "types.go", "example.org/app/generated"},
		{generated, "a/b/types.go", "example.org/app/generated/a/b"},
		{n.ImportEnvironment{Module: goImportModule, OutputRoot: "internal/gen"}, "pages/types.go", "example.org/app/internal/gen/pages"},
		{n.ImportEnvironment{Module: goImportModule, OutputRoot: n.ModuleRoot}, "pages/types.go", "example.org/app/pages"},
		{n.ImportEnvironment{Module: goImportModule, OutputRoot: n.ModuleRoot}, "types.go", "example.org/app"},
	} {
		if got, ok := n.CapturedImportPath(tc.environment, tc.file); !ok || got != tc.want {
			t.Fatal("derived import path", tc.environment.OutputRoot, tc.file, got, ok)
		}
	}
	// No identity is derived from an environment that does not state the
	// mapping, or for a path that names no file in a package directory.
	for name, tc := range map[string]struct {
		environment n.ImportEnvironment
		file        string
	}{
		"no output root":      {n.ImportEnvironment{Module: goImportModule}, "pages/types.go"},
		"unclean output root": {n.ImportEnvironment{Module: goImportModule, OutputRoot: "generated/"}, "pages/types.go"},
		"no module":           {n.ImportEnvironment{OutputRoot: "generated"}, "pages/types.go"},
		"invalid module":      {n.ImportEnvironment{Module: "app", OutputRoot: "generated"}, "pages/types.go"},
		"invalid directory":   {generated, "a b/types.go"},
		"parent directory":    {generated, "../types.go"},
		"absolute path":       {generated, "/pages/types.go"},
		"unclean directory":   {generated, "pages//types.go"},
		"backslash":           {generated, "pages\\sub/types.go"},
		"no file name":        {generated, "pages/"},
		"empty path":          {generated, ""},
	} {
		if got, ok := n.CapturedImportPath(tc.environment, tc.file); ok || got != "" {
			t.Fatal("an import path was derived:", name, got)
		}
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=supplied-package-inside-module expects=admitted-outside-output-root-only
func TestGoSuppliedPackageInsideTheModuleAndOutsideTheOutputRoot(t *testing.T) {
	supplied := func(path string) n.Package {
		return n.Package{ImportPath: path, Files: []n.PackageFile{{Name: "model.go", Text: "package model\n\ntype Person struct {\n\tName string\n}\n"}}}
	}
	// A hand-written package of the consumer's own module is supplied and used.
	model, _ := goImportModel(t, map[string]string{"pages/types.go": "package pages\n\nimport \"example.org/app/model\"\n\ntype Author = model.Person\n"}, goImportSettingsAt(t, "generated", supplied("example.org/app/model")))
	if got := goPackageNames(model); !reflect.DeepEqual(got, map[string][]string{"example.org/app/generated/pages": {"pages"}, "example.org/app/model": {"model"}}) {
		t.Fatal("packages", got)
	}
	person, err := model.LookupIn("example.org/app/model", "Person")
	if err != nil {
		t.Fatal(err)
	}
	if author, err := model.LookupIn("example.org/app/generated/pages", "Author"); err != nil || !author.AliasTarget().Same(person) {
		t.Fatal("the captured alias does not name the supplied declaration", err)
	}
	for _, tc := range []struct {
		outputRoot string
		path       string
		admitted   bool
	}{
		{"generated", "example.org/app/model", true},
		{"generated", "example.org/app", true},
		{"generated", "example.org/app/generated2", true},
		{"generated", "example.org/app/internal/generated", true},
		{"generated", "example.org/app/generated", false},
		{"generated", "example.org/app/generated/model", false},
		{"generated", "example.org/app/Generated/model", false},
		{"generated", "Example.org/App/GENERATED", false},
		{"internal/gen", "example.org/app/internal", true},
		{"internal/gen", "example.org/app/internal/gen", false},
		{"internal/gen", "example.org/app/internal/gen/model", false},
		{n.ModuleRoot, "example.org/app/model", false},
		{n.ModuleRoot, "example.org/app", false},
		{n.ModuleRoot, "example.org/application", true},
	} {
		_, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, OutputRoot: tc.outputRoot, Packages: []n.Package{supplied(tc.path)}})
		if tc.admitted && err != nil || !tc.admitted && !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal("supplied package", tc.outputRoot, tc.path, err)
		}
	}
}

// Each declaration reads its own package from the graph's ownership, and that
// identity is the one derived from the environment and the declaration's file.
//
// bittype:case target=github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze feature=go-imported-packages case=declaration-package-reader expects=owner-identity-equals-derived-identity
func TestGoDeclarationPackageIsReadFromOwnership(t *testing.T) {
	environment := n.ImportEnvironment{Module: goImportModule, OutputRoot: "generated", Packages: []n.Package{cellPackage}}
	settings, err := n.ImportConfiguration(environment)
	if err != nil {
		t.Fatal(err)
	}
	model, trees := goImportModel(t, map[string]string{
		"people/person.go": "package people\n\ntype Person struct {\n\tName string\n}\n",
		"other/person.go":  "package people\n\nimport \"example.org/cell\"\n\ntype Person struct {\n\tName string\n}\n\ntype Box[T any] struct {\n\tValue T\n}\n\ntype Boxed = Box[Person]\n\ntype Cells = cell.Cell[Person]\n",
	}, settings)
	data := model.Graph().Data()
	for file, want := range map[string]string{"people/person.go": "example.org/app/generated/people", "other/person.go": "example.org/app/generated/other"} {
		// The declaration is found by its site, not by a package path.
		var declared []n.Type
		for _, e := range data.Entities {
			if e.Name == "Person" && e.Kind == n.DefinedKind() && len(e.Declarations) == 1 && e.Declarations[0].Tree == trees[file].ID {
				ty, err := model.Type(e.Body)
				if err != nil {
					t.Fatal(err)
				}
				declared = append(declared, ty)
			}
		}
		if len(declared) != 1 {
			t.Fatal("declarations of Person in", file, len(declared))
		}
		pkg, err := declared[0].Package()
		derived, ok := n.CapturedImportPath(environment, file)
		if err != nil || !ok || pkg.Path != derived || pkg.Path != want || pkg.Name != "people" {
			t.Fatal("package of Person in", file, pkg, err, derived, ok)
		}
	}
	const other = "example.org/app/generated/other"
	lookup := func(importPath, name string) n.Type {
		ty, err := model.LookupIn(importPath, name)
		if err != nil {
			t.Fatal(importPath, name, err)
		}
		return ty
	}
	// An alias has its own package; an instance has its generic declaration's;
	// a used supplied declaration has the supplied package.
	for name, tc := range map[string]struct {
		ty   n.Type
		want n.DeclaringPackage
	}{
		"alias":    {lookup(other, "Boxed"), n.DeclaringPackage{Path: other, Name: "people"}},
		"instance": {lookup(other, "Boxed").AliasTarget(), n.DeclaringPackage{Path: other, Name: "people"}},
		"supplied": {lookup("example.org/cell", "Cell"), n.DeclaringPackage{Path: "example.org/cell", Name: "cell"}},
	} {
		if got, err := tc.ty.Package(); err != nil || got != tc.want {
			t.Fatal(name, got, err)
		}
	}
	// What is not a package-level declaration has no package.
	text, err := model.Basic(types.String)
	if err != nil {
		t.Fatal(err)
	}
	box := lookup(other, "Box")
	if len(box.Parameters()) != 1 {
		t.Fatal("Box has no type parameter")
	}
	for name, ty := range map[string]n.Type{
		"predeclared type": text,
		"type parameter":   box.Parameters()[0],
		"unnamed type":     lookup(other, "Person").Underlying(),
		"zero handle":      {},
	} {
		if got, err := ty.Package(); !errors.Is(err, n.ErrForeign) || got != (n.DeclaringPackage{}) {
			t.Fatal(name, "has a package:", got, err)
		}
	}
	// In the base profile the identity is the constant path.
	base := goModel(t, "type Text string")
	if got, err := goType(t, base, "Text").Package(); err != nil || got != (n.DeclaringPackage{Path: n.PackagePath, Name: "example"}) {
		t.Fatal("base declaration package", got, err)
	}
}

// The first Go Files slice admits a file name only when it can carry no build
// constraint, and an ordinary package only.
func TestGoFilesAdmitOnlyNamesWithoutBuildConstraints(t *testing.T) {
	judged := func(err error) string {
		unsupported, invalid := errors.Is(err, n.ErrUnsupported), errors.Is(err, semantic.ErrInvalidRequest)
		switch {
		case err == nil:
			return "admitted"
		case unsupported && !invalid:
			return "unsupported"
		case invalid && !unsupported:
			return "invalid"
		}
		return "uncategorized: " + err.Error()
	}
	for want, paths := range map[string][]string{
		// "linux.go" has no underscore, so Go reads no constraint from it either.
		// A directory is not judged: the rule reads the base name.
		"admitted":    {"pages/types.go", "types.go", "pages/Types.go", "a/b/types.gen.go", "pages/linux.go", "my_pages/types.go", ".hidden/types.go"},
		"unsupported": {"pages/types_linux.go", "pages/types_amd64.go", "pages/types_test.go", "pages/_types.go", "pages/.types.go", "pages/types_gen.go", "types_linux_amd64.go", "pages/_.go", "pages/.go", "_test.go"},
		"invalid":     {"pages/types.txt", "pages/types", "pages/types_linux.txt", "pages/types.GO", "pages/types.go/", "pages/", "", "pages\\.types.go", "pages/ty\x00pes.go"},
	} {
		for _, path := range paths {
			if got := judged(n.AdmitFileName(path)); got != want {
				t.Errorf("file name %q: %s, want %s", path, got, want)
			}
		}
	}
	for want, names := range map[string][]string{
		"admitted":    {"pages", "people", "Pages", "p1", "a_b", "test", "tests", "mainline", "testing_"},
		"unsupported": {"main", "pages_test", "main_test", "_test"},
		"invalid":     {"", "_", "1pages", "my-pages", "a.b", "pages people", "type", "package"},
	} {
		for _, name := range names {
			if got := judged(n.AdmitPackageName(name)); got != want {
				t.Errorf("package name %q: %s, want %s", name, got, want)
			}
		}
	}
}

// The environment must state where the output root sits. Nothing is defaulted,
// and a setting written before the statement existed is refused.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=output-root-must-be-stated expects=absent-or-unclean-root-refused
func TestGoImportEnvironmentMustStateTheOutputRoot(t *testing.T) {
	// The bytes of the setting, which every receipt under an import environment carries.
	for outputRoot, want := range map[string]string{
		"generated":  `{"Module":"example.org/app","OutputRoot":"generated","Packages":[]}`,
		n.ModuleRoot: `{"Module":"example.org/app","OutputRoot":".","Packages":[]}`,
	} {
		settings := goImportSettingsAt(t, outputRoot)
		if len(settings) != len(n.Configuration())+1 {
			t.Fatal(settings)
		}
		if last := settings[len(settings)-1]; last.Key != n.Key("imports-2") || string(last.Value) != want {
			t.Fatal("import setting", last.Key, string(last.Value))
		}
		environment, enabled, err := n.Imports(settings)
		if err != nil || !enabled || environment.Module != goImportModule || environment.OutputRoot != outputRoot {
			t.Fatal("decoded environment", environment, enabled, err)
		}
	}
	// The module root is stated as ".", which is not the absent value.
	if n.ModuleRoot != "." {
		t.Fatal("the stated module root changed", n.ModuleRoot)
	}
	admitted := map[string]bool{}
	for _, outputRoot := range []string{n.ModuleRoot, "generated", "Generated", "internal/gen", "a/b/c", "gen-1.x", "v2"} {
		settings := goImportSettingsAt(t, outputRoot)
		admitted[string(settings[len(settings)-1].Value)] = true
	}
	if len(admitted) != 7 {
		t.Fatal("two output roots gave one setting", admitted)
	}
	for name, outputRoot := range map[string]string{
		"absent":               "",
		"parent":               "..",
		"parent first":         "../generated",
		"parent inside":        "generated/../other",
		"parent last":          "generated/..",
		"absolute":             "/generated",
		"drive":                "C:/generated",
		"trailing slash":       "generated/",
		"current first":        "./generated",
		"current last":         "generated/.",
		"module root, slashed": "./",
		"empty element":        "generated//pages",
		"backslash":            "generated\\pages",
		"space":                " generated",
		"hidden directory":     ".generated",
		"trailing dot":         "generated.",
		"underscore directory": "_generated",
		"underscore inside":    "a/_generated",
		"testdata":             "testdata",
		"testdata inside":      "a/testdata/b",
		"vendor":               "vendor",
		"vendor inside":        "a/vendor",
	} {
		if _, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, OutputRoot: outputRoot}); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal("an output root was accepted:", name, err)
		}
	}
	// A setting under the old key, and one whose bytes state no clean output
	// root, is refused wherever a configuration is read.
	snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
	resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	current := goImportSettingsAt(t, "generated")
	if _, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: current}); err != nil {
		t.Fatal(err)
	}
	for name, last := range map[string]semantic.Setting{
		"old key, old bytes":           {Key: n.Key("imports-1"), Value: []byte(`{"Module":"example.org/app","Packages":[]}`)},
		"old key, current bytes":       {Key: n.Key("imports-1"), Value: current[len(current)-1].Value},
		"current key, no output root":  {Key: n.Key("imports-2"), Value: []byte(`{"Module":"example.org/app","Packages":[]}`)},
		"current key, empty root":      {Key: n.Key("imports-2"), Value: []byte(`{"Module":"example.org/app","OutputRoot":"","Packages":[]}`)},
		"current key, unclean root":    {Key: n.Key("imports-2"), Value: []byte(`{"Module":"example.org/app","OutputRoot":"generated/","Packages":[]}`)},
		"current key, root after list": {Key: n.Key("imports-2"), Value: []byte(`{"Module":"example.org/app","Packages":[],"OutputRoot":"generated"}`)},
	} {
		altered := append(slices.Clone(n.Configuration()), last)
		if _, _, err := n.Imports(altered); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, "accepted by Imports:", err)
		}
		if n.MatchesConfiguration(altered) {
			t.Fatal(name, "matches the configuration")
		}
		if _, err := n.AnalysisContractFor(altered); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, "accepted by AnalysisContractFor:", err)
		}
		if _, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: altered}); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, "accepted by the resolver:", err)
		}
	}
}
