package integration_test

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	resolution "github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func goImportsHold(t *testing.T, result semantic.Result) {
	t.Helper()
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		report, _ := result.Report(key)
		if report.Status != semantic.Complete || report.Verdict != semantic.Holds {
			t.Fatal(key, report, result.Diagnostics())
		}
	}
}

// Captured declarations built directly on supplied structure export only the
// supplied identities, so unsupported supplied internals neither leak nor fail.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=captured-types-on-supplied-structure expects=supplied-structure-not-exported
func TestGoImportsDoNotExportSuppliedStructureThroughCapturedTypes(t *testing.T) {
	box := n.Package{ImportPath: "example.org/box", Files: []n.PackageFile{{Name: "box.go", Text: "package box\n\ntype Box struct {\n\tsecret string `json:\"hidden\"`\n\tnext   *Box\n}\n\ntype C interface {\n\t~string\n\tM() *int\n}\n\ntype Ptr = *Box\n\ntype G[T interface{ ~string | *int }] struct {\n\tv T\n}\n"}}}
	snapshot, _ := goPackageCapture(t, map[string]string{
		"a/a.go": "package a\n\nimport (\n\t\"example.org/box\"\n\t\"example.org/cell\"\n)\n\ntype A = cell.Cell[string]\n\ntype B box.Box\n\ntype D B\n\ntype W[T any] cell.Cell[T]\n\ntype P[T box.C] struct {\n\tV T\n}\n\ntype E = box.Ptr\n\ntype H = box.G[string]\n",
	}, true)
	_, result, err := goImportResults(t, snapshot, goImportSettings(t, box, cellPackage))
	if err != nil {
		t.Fatal(err)
	}
	goImportsHold(t, result)
	model, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	data := model.Graph().Data()
	for _, e := range data.Entities {
		if e.Name == "secret" || e.Name == "next" || e.Name == "value" {
			t.Fatal("a supplied field was exported", e)
		}
	}
	for _, term := range data.Terms {
		if strings.Contains(string(term.Data), "hidden") || strings.Contains(string(term.Data), "secret") {
			t.Fatal("supplied structure reached the graph", string(term.Data))
		}
	}
	for _, name := range []string{"A", "B", "D", "W", "E", "H"} {
		ty, err := model.LookupIn("example.org/app/a", name)
		if err != nil {
			t.Fatal(name, err)
		}
		if facts := model.Graph().Facts(semantics.Subject{Term: ty.ID()}, n.Key("underlying")); len(facts) != 0 {
			t.Fatal(name, "exports a supplied underlying structure")
		}
	}
}

// Analysis and model admission are bound to the environment the graph was computed under.
//
// bittype:case target=github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze feature=go-imported-packages case=resolution-and-analysis-environment-mismatch expects=invalid-request
func TestGoImportsAnalysisRefusesAnotherEnvironment(t *testing.T) {
	ctx := context.Background()
	imports := goImportSettings(t, cellPackage)
	for name, tc := range map[string]struct{ resolve, analyze []semantic.Setting }{
		"base resolution, import analysis": {n.Configuration(), imports},
		"import resolution, base analysis": {imports, n.Configuration()},
	} {
		t.Run(name, func(t *testing.T) {
			snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
			resolver, _ := resolution.NewResolver[s.Language](snapshot.Language(), resolvegolang.Contract(), resolvegolang.Provider{})
			resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: tc.resolve})
			if err != nil {
				t.Fatal(err)
			}
			contract, err := analyze.ContractFor(tc.analyze)
			if err != nil {
				t.Fatal(err)
			}
			// The contract names the analysis environment; the resolution supplied
			// was computed under the other one.
			contract.Resolution.Configuration = tc.resolve
			analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), contract, analyze.Provider{})
			if err != nil {
				t.Fatal(err)
			}
			if _, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: tc.analyze}}); !errors.Is(err, semantic.ErrInvalidRequest) {
				t.Fatal(err)
			}
		})
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=go-module-layout-rules expects=go-tool-classification
func TestGoImportsFollowGoLayoutRules(t *testing.T) {
	internal := n.Package{ImportPath: "example.org/lib/internal/x", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n\ntype X string\n"}}}
	program := n.Package{ImportPath: "example.org/tool", Files: []n.PackageFile{{Name: "main.go", Text: "package main\n\ntype T string\n"}}}
	for _, tc := range []struct {
		name     string
		sources  map[string]string
		packages []n.Package
		coverage semantic.Completion
		validity semantic.Verdict
	}{
		{"test file", map[string]string{"a/a_test.go": "package a\n\ntype A string\n"}, nil, semantic.Unsupported, semantic.Undetermined},
		{"testdata", map[string]string{"testdata/a.go": "package a\n\ntype A string\n"}, nil, semantic.Unsupported, semantic.Undetermined},
		{"underscore directory", map[string]string{"_a/a.go": "package a\n\ntype A string\n"}, nil, semantic.Unsupported, semantic.Undetermined},
		{"invalid package directory", map[string]string{"a b/a.go": "package a\n\ntype A string\n"}, nil, semantic.Unsupported, semantic.Undetermined},
		{"supplied internal from outside", map[string]string{"a/a.go": "package a\n\nimport \"example.org/lib/internal/x\"\n\ntype A = x.X\n"}, []n.Package{internal}, semantic.Complete, semantic.Rejected},
		{"captured internal from outside", map[string]string{"a/a.go": "package a\n\nimport \"example.org/app/b/internal/c\"\n\ntype A = c.C\n", "b/internal/c/c.go": "package c\n\ntype C string\n"}, nil, semantic.Complete, semantic.Rejected},
		{"captured internal from within", map[string]string{"b/b.go": "package b\n\nimport \"example.org/app/b/internal/c\"\n\ntype B = c.C\n", "b/internal/c/c.go": "package c\n\ntype C string\n"}, nil, semantic.Complete, semantic.Holds},
		{"import of a program", map[string]string{"a/a.go": "package a\n\nimport \"example.org/tool\"\n\ntype A = main.T\n"}, []n.Package{program}, semantic.Complete, semantic.Rejected},
		{"self import", map[string]string{"a/a.go": "package a\n\nimport \"example.org/app/a\"\n\ntype A = a.B\n\ntype B string\n"}, nil, semantic.Complete, semantic.Rejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snapshot, _ := goPackageCapture(t, tc.sources, true)
			_, result, err := goImportResults(t, snapshot, goImportSettings(t, tc.packages...))
			if err != nil {
				t.Fatal(err)
			}
			coverage, _ := result.Report(n.Coverage())
			validity, _ := result.Report(n.Validity())
			if coverage.Status != tc.coverage || validity.Verdict != tc.validity {
				t.Fatal(coverage, validity, result.Diagnostics())
			}
			if tc.validity != semantic.Holds {
				if ds := result.Diagnostics(); len(ds) == 0 || ds[0].At == nil || ds[0].At.Tree == "" {
					t.Fatal("missing captured diagnostic", ds)
				}
			}
		})
	}
}

// The environment is validated before captured code is judged, independently of
// how packages and trees sort.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=environment-before-captured-judgment expects=order-independent-classification
func TestGoImportsClassifyIndependentlyOfOrder(t *testing.T) {
	broken := func(path string) n.Package {
		return n.Package{ImportPath: path, Files: []n.PackageFile{{Name: "b.go", Text: "package b\n\ntype B Missing\n"}}}
	}
	for _, path := range []string{"a.org/broken", "z.org/broken"} {
		snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A Undefined\n"}, true)
		if _, _, err := goImportResults(t, snapshot, goImportSettings(t, broken(path))); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(path, err)
		}
	}
	// A broken environment outranks an unsupported captured construct.
	snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\nimport . \"example.org/cell\"\n\ntype A = Cell[string]\n"}, true)
	if _, _, err := goImportResults(t, snapshot, goImportSettings(t, cellPackage, broken("a.org/broken"))); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	// A cycle among supplied packages is the environment's fault.
	cyclic := []n.Package{
		{ImportPath: "example.org/p", Files: []n.PackageFile{{Name: "p.go", Text: "package p\n\nimport \"example.org/q\"\n\ntype P = q.Q\n"}}},
		{ImportPath: "example.org/q", Files: []n.PackageFile{{Name: "q.go", Text: "package q\n\nimport \"example.org/p\"\n\ntype Q = p.P\n"}}},
	}
	snapshot, _ = goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
	if _, _, err := goImportResults(t, snapshot, goImportSettings(t, cyclic...)); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	// The same sources captured in different tree orders give the same graph.
	sources := map[string]string{
		"pages/page.go":    "package pages\n\nimport \"example.org/app/people\"\n\ntype Page struct {\n\tItems []people.Person\n}\n",
		"people/person.go": "package people\n\ntype Person struct {\n\tName string\n}\n",
	}
	var graphs []semantics.Data
	for _, order := range [][]string{{"pages/page.go", "people/person.go"}, {"people/person.go", "pages/page.go"}} {
		snapshot, _ := goPackageCaptureIn(t, order, sources, true)
		_, result, err := goImportResults(t, snapshot, goImportSettings(t))
		if err != nil {
			t.Fatal(err)
		}
		goImportsHold(t, result)
		model, err := n.Open(result)
		if err != nil {
			t.Fatal(err)
		}
		data := model.Graph().Data()
		// Package declarations list roots in capture order; compare the rest exactly.
		for i := range data.Entities {
			if data.Entities[i].Kind == n.Key("package") {
				data.Entities[i].Declarations = nil
			}
		}
		graphs = append(graphs, data)
	}
	if !reflect.DeepEqual(graphs[0].Entities, graphs[1].Entities) || len(graphs[0].Bindings) != len(graphs[1].Bindings) {
		t.Fatal("tree order changed the exported graph")
	}
}

// bittype:case target=github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze feature=go-imported-packages case=captured-lookup-and-base-parity expects=captured-ambiguity-only
func TestGoImportsLookupAndBaseParity(t *testing.T) {
	// Base parity: one package entity 1 named "captured".
	base := goModel(t, "type Text string")
	if first := base.Graph().Data().Entities[0]; first.Kind != n.Key("package") || first.Name != n.PackagePath {
		t.Fatal("base configuration changed its package entity", first)
	}
	// Base parity: a qualified type stays unsupported, naming its node kind.
	qualified, _ := goCapture(t, "type Box x.Y")
	_, refused := goResults(t, qualified)
	if ds := refused.Diagnostics(); len(ds) == 0 || !strings.Contains(ds[0].Message, string(s.KindQualifiedType)) {
		t.Fatal("base qualified-type diagnostic changed", ds)
	}
	// An export-time limit is reported at the declaration being exported. The
	// first declaration starts the deepest alias chain, past the export depth.
	chain := "type A140 string"
	for i := 139; i >= 0; i-- {
		chain = "type A" + strconv.Itoa(i) + " = A" + strconv.Itoa(i+1) + "; " + chain
	}
	deep, _ := goCapture(t, chain)
	_, limited := goResults(t, deep)
	if coverage, _ := limited.Report(n.Coverage()); coverage.Status != semantic.Incomplete {
		t.Fatal("expected an incomplete result at the export limit", coverage, limited.Diagnostics())
	}
	if ds := limited.Diagnostics(); len(ds) == 0 || ds[0].At == nil || ds[0].At.Tree == "" {
		t.Fatal("an export-time limit has no captured site", ds)
	}
	snapshot, _ := goPackageCapture(t, map[string]string{
		"a/a.go": "package a\n\nimport \"example.org/other\"\n\ntype Person struct{}\n\ntype A = other.A\n",
		"b/b.go": "package b\n\ntype Person struct{}\n",
	}, true)
	other := n.Package{ImportPath: "example.org/other", Files: []n.PackageFile{{Name: "o.go", Text: "package other\n\ntype A string\n\ntype Unused string\n"}}}
	_, result, err := goImportResults(t, snapshot, goImportSettings(t, other))
	if err != nil {
		t.Fatal(err)
	}
	model, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := model.Lookup("Person"); err == nil {
		t.Fatal("an ambiguous captured name was resolved")
	}
	for _, path := range []string{"example.org/app/a", "example.org/app/b"} {
		if _, err := model.LookupIn(path, "Person"); err != nil {
			t.Fatal(path, err)
		}
	}
	// A used supplied declaration never makes a captured name ambiguous.
	captured, err := model.Lookup("A")
	if err != nil {
		t.Fatal(err)
	}
	if want, _ := model.LookupIn("example.org/app/a", "A"); !captured.Same(want) {
		t.Fatal("Lookup did not return the captured declaration")
	}
	if _, err := model.LookupIn("example.org/other", "A"); err != nil {
		t.Fatal(err)
	}
	if _, err := model.LookupIn("example.org/other", "Unused"); err == nil {
		t.Fatal("an unused supplied declaration was exported")
	}
}

// Every configuration ImportConfiguration accepts can be carried by analysis
// contracts and model admission.
//
// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=configuration-at-the-evidence-budget expects=accepted-environments-admissible
func TestGoImportConfigurationFitsEvidence(t *testing.T) {
	environment := func(size int) n.ImportEnvironment {
		// JSON escapes '<', so the encoded setting grows six bytes per character.
		text := "package big\n\n// " + strings.Repeat("<", size) + "\n"
		return n.ImportEnvironment{Module: goImportModule, Packages: []n.Package{{ImportPath: "example.org/big", Files: []n.PackageFile{{Name: "big.go", Text: text}}}}}
	}
	lo, hi := 0, 1<<20
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if _, err := n.ImportConfiguration(environment(mid)); err == nil {
			lo = mid
		} else if errors.Is(err, semantic.ErrLimit) {
			hi = mid - 1
		} else {
			t.Fatal(err)
		}
	}
	settings, err := n.ImportConfiguration(environment(lo))
	if err != nil {
		t.Fatal(err)
	}
	contract, err := analyze.ContractFor(settings)
	if err != nil {
		t.Fatal(err)
	}
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "go-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := analysis.NewAnalyzer[s.Language](definition, contract, analyze.Provider{}); err != nil {
		t.Fatal("the largest accepted environment cannot be analyzed:", err)
	}
}
