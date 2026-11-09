package integration_test

import (
	"context"
	"errors"
	"reflect"
	"sort"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
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
	settings, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, Packages: packages})
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

// goTreeID gives each file a stable tree ID that sorts opposite to its path:
// lowercase letters are mirrored (a↔z, b↔y, ...). Results that follow file-path
// order therefore cannot be produced by sorting tree IDs.
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
	resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), resolvegolang.Contract(), resolvegolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: settings})
	if err != nil {
		return semantic.Result{}, semantic.Result{}, err
	}
	contract, err := analyze.ContractFor(settings)
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
	if _, err := analyze.NewQueries(context.Background(), model); err != nil {
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
	a, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, Packages: []n.Package{two, cellPackage}})
	if err != nil {
		t.Fatal(err)
	}
	b, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, Packages: []n.Package{cellPackage, {ImportPath: two.ImportPath, Files: []n.PackageFile{two.Files[1], two.Files[0]}}}})
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
		"module without dot":    {Module: "app"},
		"dot-dot element":       {Module: "example.org/../app"},
		"supplied under module": {Module: goImportModule, Packages: []n.Package{{ImportPath: goImportModule + "/x", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n"}}}}},
		"duplicate package":     {Module: goImportModule, Packages: []n.Package{cellPackage, cellPackage}},
		"not a go file":         {Module: goImportModule, Packages: []n.Package{{ImportPath: "example.org/x", Files: []n.PackageFile{{Name: "x.txt", Text: "package x\n"}}}}},
		"no files":              {Module: goImportModule, Packages: []n.Package{{ImportPath: "example.org/x"}}},
	} {
		if _, err := n.ImportConfiguration(environment); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, err)
		}
	}
	many := n.ImportEnvironment{Module: goImportModule}
	for i := 0; i < 65; i++ {
		many.Packages = append(many.Packages, n.Package{ImportPath: "example.org/p" + string(rune('a'+i%26)) + string(rune('a'+i/26)), Files: []n.PackageFile{{Name: "p.go", Text: "package p\n"}}})
	}
	if _, err := n.ImportConfiguration(many); !errors.Is(err, semantic.ErrLimit) {
		t.Fatal(err)
	}
	// Analysis is bound to the exact environment its resolution used.
	snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\nimport \"example.org/cell\"\n\ntype A = cell.Cell[string]\n"}, true)
	resolver, _ := resolution.NewResolver[s.Language](snapshot.Language(), resolvegolang.Contract(), resolvegolang.Provider{})
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: goImportSettings(t, cellPackage)})
	if err != nil {
		t.Fatal(err)
	}
	other := n.Package{ImportPath: cellPackage.ImportPath, Files: []n.PackageFile{{Name: "cell.go", Text: "package cell\n\ntype Cell[T any] struct{ Value T }\n"}}}
	otherSettings := goImportSettings(t, other)
	contract, err := analyze.ContractFor(otherSettings)
	if err != nil {
		t.Fatal(err)
	}
	analyzer, _ := analysis.NewAnalyzer[s.Language](snapshot.Language(), contract, analyze.Provider{})
	if _, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: otherSettings}}); err == nil {
		t.Fatal("analysis accepted a resolution from a different supplied package text")
	}
}
