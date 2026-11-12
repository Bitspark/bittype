package integration_test

import (
	"context"
	"errors"
	"go/types"
	"reflect"
	"strconv"
	"strings"
	"testing"

	nativequery "github.com/Bitspark/bittype/analyze/files/golang/queries/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

// plainPackage has an unexported tagged field, a pointer alias, a constraint with
// a method and a pointer term, and a generic whose constraint has a pointer term.
var plainPackage = n.Package{ImportPath: "example.org/plain", Files: []n.PackageFile{{Name: "plain.go", Text: "package plain\n\ntype S struct {\n\tsecret string `json:\"hidden\"`\n\tnext   *S\n}\n\ntype Ptr = *S\n\ntype C interface {\n\t~string\n\tM() *int\n}\n\ntype G[T interface{ ~string | *int }] struct {\n\tv T\n}\n"}}}

// goNoSuppliedText checks the structural projection: no supplied member or its
// metadata. Human-readable display text may name supplied types; it is never
// semantic identity.
func goNoSuppliedText(t *testing.T, data semantics.Data) {
	t.Helper()
	for _, term := range data.Terms {
		if term.Operator == n.Key("display") {
			continue
		}
		for _, leak := range []string{"hidden", "secret"} {
			if strings.Contains(string(term.Data), leak) {
				t.Fatalf("supplied structure %q reached the graph: %s", leak, term.Data)
			}
		}
	}
	for _, e := range data.Entities {
		if e.Name == "secret" || e.Name == "next" || e.Name == "v" {
			t.Fatal("a supplied field was exported", e)
		}
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=long-chains-over-supplied-structure expects=supplied-structure-not-exported
func TestGoImportsGuardLongChainsOverSuppliedStructure(t *testing.T) {
	var b strings.Builder
	b.WriteString("package a\n\nimport \"example.org/plain\"\n\ntype T0 plain.S\n\ntype R0 = plain.Ptr\n\ntype C0 plain.C\n")
	for i := 1; i <= 70; i++ {
		k, p := strconv.Itoa(i), strconv.Itoa(i-1)
		b.WriteString("\ntype T" + k + " T" + p + "\n\ntype R" + k + " = R" + p + "\n\ntype C" + k + " C" + p + "\n")
	}
	b.WriteString("\ntype P[T C70] struct {\n\tV T\n}\n\ntype H = plain.G[string]\n")
	snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": b.String()}, true)
	_, result, err := goImportResults(t, snapshot, goImportSettings(t, plainPackage))
	if err != nil {
		t.Fatal(err)
	}
	goImportsHold(t, result)
	model, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	goNoSuppliedText(t, model.Graph().Data())
	// A supplied generic is admitted by the native compiler like any other, and
	// a constructor must be a generic declaration, never an instance or a
	// non-generic type.
	g, err := model.LookupIn("example.org/plain", "G")
	if err != nil {
		t.Fatal(err)
	}
	queries, err := nativequery.NewQueries(context.Background(), model)
	if err != nil {
		t.Fatal(err)
	}
	boolean, err := model.Basic(types.Bool)
	if err != nil {
		t.Fatal(err)
	}
	text, err := model.Basic(types.String)
	if err != nil {
		t.Fatal(err)
	}
	if admission, err := queries.Admit(context.Background(), g, boolean); err != nil || admission.Admitted() {
		t.Fatal("bool admitted for a ~string | *int constraint", err)
	}
	if admission, err := queries.Admit(context.Background(), g, text); err != nil || !admission.Admitted() {
		t.Fatal("string refused for a ~string | *int constraint", err, admission.Detail())
	}
	alias, err := model.LookupIn("example.org/app/a", "H")
	if err != nil {
		t.Fatal(err)
	}
	// H = plain.G[string]: its target is an instance of the supplied generic.
	instance := alias.AliasTarget()
	if !instance.Valid() || len(instance.Arguments()) != 1 || !instance.Origin().Same(g) {
		t.Fatal("expected H's target to instantiate plain.G", instance.String())
	}
	t0, err := model.LookupIn("example.org/app/a", "T0")
	if err != nil {
		t.Fatal(err)
	}
	for name, tc := range map[string]struct {
		constructor n.Type
		arguments   []n.Type
	}{"an instance": {instance, []n.Type{text}}, "a non-generic type": {t0, nil}} {
		if _, err := queries.Admit(context.Background(), tc.constructor, tc.arguments...); !errors.Is(err, n.ErrUnsupported) {
			t.Fatal(name, "accepted as a constructor:", err)
		}
	}
}

// Every environment ImportConfiguration accepts survives a full resolve, analyze and open.
//
// bittype:case target=github.com/Bitspark/bittype/analyze/files/golang/provider/go#Provider.Analyze feature=go-imported-packages case=largest-environment-end-to-end expects=accepted-environments-analyzable
func TestGoImportLargestEnvironmentEndToEnd(t *testing.T) {
	environment := func(size int) n.ImportEnvironment {
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
	snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
	_, result, err := goImportResults(t, snapshot, settings)
	if err != nil {
		t.Fatal("the largest accepted environment cannot be resolved:", err)
	}
	if _, err := n.Open(result); err != nil {
		t.Fatal("the largest accepted environment cannot be opened:", err, result.Diagnostics())
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=deterministic-diagnostics-and-graphs expects=capture-order-independent
func TestGoImportsDiagnosticsAndGraphsIgnoreCaptureOrder(t *testing.T) {
	// The environment is judged before captured limits: 17 captured trees exceed
	// the captured-file bound, but the broken supplied package decides.
	many := map[string]string{}
	for i := 0; i < 17; i++ {
		many["p"+strconv.Itoa(i)+"/x.go"] = "package p" + strconv.Itoa(i) + "\n\ntype X string\n"
	}
	snapshot, _ := goPackageCapture(t, many, true)
	broken := n.Package{ImportPath: "example.org/broken", Files: []n.PackageFile{{Name: "b.go", Text: "package broken\n\ntype B Missing\n"}}}
	if _, _, err := goImportResults(t, snapshot, goImportSettings(t, broken)); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	// Disagreeing package clauses give one diagnostic whatever the capture order.
	sources := map[string]string{
		"t/a.go": "package main\n\ntype A string\n",
		"t/b.go": "package t\n\ntype B string\n",
		"u/u.go": "package u\n\nimport \"example.org/app/t\"\n\ntype U = t.B\n",
	}
	var messages []string
	for _, order := range [][]string{{"t/a.go", "t/b.go", "u/u.go"}, {"u/u.go", "t/b.go", "t/a.go"}} {
		snapshot, trees := goPackageCaptureIn(t, order, sources, true)
		_, result, err := goImportResults(t, snapshot, goImportSettings(t))
		if err != nil {
			t.Fatal(err)
		}
		ds := result.Diagnostics()
		if len(ds) == 0 || ds[0].At == nil || ds[0].At.Tree != trees["t/b.go"].ID {
			t.Fatal("expected the disagreeing clause at t/b.go", ds)
		}
		messages = append(messages, ds[0].Message)
	}
	if messages[0] != messages[1] {
		t.Fatal("capture order changed the diagnostic", messages)
	}
	// The whole exported graph is the same for any capture order.
	held := map[string]string{
		"pages/page.go":    "package pages\n\nimport (\n\t\"example.org/app/people\"\n\t\"example.org/cell\"\n)\n\ntype Page[T any] struct {\n\tItems []T\n}\n\ntype Cells = Page[cell.Cell[people.Person]]\n",
		"people/person.go": "package people\n\ntype Person struct {\n\tName string\n}\n",
	}
	var graphs []semantics.Data
	for _, order := range [][]string{{"pages/page.go", "people/person.go"}, {"people/person.go", "pages/page.go"}} {
		snapshot, _ := goPackageCaptureIn(t, order, held, true)
		_, result, err := goImportResults(t, snapshot, goImportSettings(t, cellPackage))
		if err != nil {
			t.Fatal(err)
		}
		goImportsHold(t, result)
		model, err := n.Open(result)
		if err != nil {
			t.Fatal(err)
		}
		graphs = append(graphs, model.Graph().Data())
	}
	if !reflect.DeepEqual(graphs[0], graphs[1]) {
		t.Fatal("capture order changed the exported graph")
	}
	// Captured package entities come in import-path order, then supplied ones.
	var order []string
	for _, e := range graphs[0].Entities {
		if e.Kind == n.Key("package") {
			order = append(order, e.Name)
		}
	}
	if !reflect.DeepEqual(order, []string{"example.org/app/pages", "example.org/app/people", "example.org/cell"}) {
		t.Fatal("package entity order", order)
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=go-tool-collisions-and-sites expects=go-build-agreement
func TestGoImportsCollisionsAndSites(t *testing.T) {
	for _, tc := range []struct {
		name     string
		sources  map[string]string
		coverage semantic.Completion
		validity semantic.Verdict
	}{
		{"case-insensitive packages", map[string]string{"A/x.go": "package a\n\ntype X string\n", "a/x.go": "package a\n\ntype X string\n"}, semantic.Complete, semantic.Rejected},
		{"underscore file", map[string]string{"a/_x.go": "package a\n\ntype X string\n"}, semantic.Unsupported, semantic.Undetermined},
		{"case-insensitive files", map[string]string{"a/X.go": "package a\n\ntype X string\n", "a/x.go": "package a\n\ntype Y string\n"}, semantic.Complete, semantic.Rejected},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snapshot, _ := goPackageCapture(t, tc.sources, true)
			_, result, err := goImportResults(t, snapshot, goImportSettings(t))
			if err != nil {
				t.Fatal(err)
			}
			coverage, _ := result.Report(n.Coverage())
			validity, _ := result.Report(n.Validity())
			if coverage.Status != tc.coverage || validity.Verdict != tc.validity {
				t.Fatal(coverage, validity, result.Diagnostics())
			}
		})
	}
	// Environments the go tool refuses are invalid configurations.
	for name, packages := range map[string][]n.Package{
		"case-insensitive supplied paths":   {{ImportPath: "example.org/x", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n"}}}, {ImportPath: "example.org/X", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n"}}}},
		"case-insensitive file names":       {{ImportPath: "example.org/x", Files: []n.PackageFile{{Name: "A.go", Text: "package x\n"}, {Name: "B.go", Text: "package x\n"}, {Name: "a.go", Text: "package x\n"}}}},
		"supplied test file":                {{ImportPath: "example.org/x", Files: []n.PackageFile{{Name: "x_test.go", Text: "package x_test\n"}}}},
		"supplied ignored file":             {{ImportPath: "example.org/x", Files: []n.PackageFile{{Name: "_x.go", Text: "package x\n"}}}},
		"supplied under the module, folded": {{ImportPath: "Example.org/App/x", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n"}}}},
	} {
		if _, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, Packages: packages}); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, err)
		}
	}
	// Supplied packages obey internal and main rules, as an invalid request.
	internal := n.Package{ImportPath: "example.org/lib/internal/x", Files: []n.PackageFile{{Name: "x.go", Text: "package x\n\ntype X string\n"}}}
	outside := n.Package{ImportPath: "example.org/other", Files: []n.PackageFile{{Name: "o.go", Text: "package other\n\nimport \"example.org/lib/internal/x\"\n\ntype O = x.X\n"}}}
	program := n.Package{ImportPath: "example.org/tool", Files: []n.PackageFile{{Name: "main.go", Text: "package main\n\ntype T string\n"}}}
	user := n.Package{ImportPath: "example.org/user", Files: []n.PackageFile{{Name: "u.go", Text: "package user\n\nimport \"example.org/tool\"\n\ntype U = main.T\n"}}}
	for name, packages := range map[string][]n.Package{"supplied internal from outside": {internal, outside}, "supplied imports a program": {program, user}} {
		snapshot, _ := goPackageCapture(t, map[string]string{"a/a.go": "package a\n\ntype A string\n"}, true)
		if _, _, err := goImportResults(t, snapshot, goImportSettings(t, packages...)); !errors.Is(err, semantic.ErrInvalidRequest) {
			t.Fatal(name, err)
		}
	}
	// Diagnostics sit at the node that causes them.
	for name, tc := range map[string]struct {
		source string
		kind   s.Kind
		attr   string
		value  string
	}{
		"unused import":  {"package a\n\nimport \"example.org/app/b\"\n\ntype A string\n", s.KindImport, s.AttributePath, "example.org/app/b"},
		"cycle":          {"package a\n\nimport \"example.org/app/b\"\n\ntype A = b.B\n", s.KindImport, s.AttributePath, "example.org/app/b"},
		"unexported use": {"package a\n\nimport \"example.org/app/b\"\n\ntype A = b.hidden\n", s.KindQualifiedType, s.AttributeName, "hidden"},
	} {
		t.Run(name, func(t *testing.T) {
			b := "package b\n\ntype B string\n\ntype hidden string\n"
			if name == "cycle" {
				b = "package b\n\nimport \"example.org/app/a\"\n\ntype B = a.A\n"
			}
			snapshot, trees := goPackageCapture(t, map[string]string{"a/a.go": tc.source, "b/b.go": b}, true)
			_, result, err := goImportResults(t, snapshot, goImportSettings(t))
			if err != nil {
				t.Fatal(err)
			}
			want := nodeOf(t, trees["a/a.go"], tc.kind, tc.attr, tc.value)
			if ds := result.Diagnostics(); len(ds) == 0 || ds[0].At == nil || *ds[0].At != want {
				t.Fatal("diagnostic site", ds, want)
			}
		})
	}
	// LookupIn needs a package; Imports and the analysis contract refuse altered settings.
	model := goModel(t, "type Text string")
	if _, err := model.LookupIn("", "Text"); err == nil {
		t.Fatal("LookupIn accepted an empty import path")
	}
	settings := goImportSettings(t, cellPackage)
	reordered := append([]semantic.Setting{settings[len(settings)-1]}, settings[:len(settings)-1]...)
	extra := append(append([]semantic.Setting(nil), settings...), semantic.Setting{Key: n.Key("extra"), Value: []byte("x")})
	for name, altered := range map[string][]semantic.Setting{"reordered": reordered, "extra": extra} {
		if _, _, err := n.Imports(altered); err == nil {
			t.Fatal(name, "settings accepted by Imports")
		}
		if _, err := n.AnalysisContractFor(altered); err == nil {
			t.Fatal(name, "settings accepted by AnalysisContractFor")
		}
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/provider/go#Provider.Resolve feature=go-imported-packages case=documented-check-order expects=layout-then-syntax-then-go-tool-rejections
func TestGoImportsCheckOrderAndVendor(t *testing.T) {
	// A layout refusal comes before any tree's syntax limit.
	deep := "package b\n\ntype B " + strings.Repeat("[]", 70) + "string\n"
	snapshot, trees := goPackageCapture(t, map[string]string{"a/a_test.go": "package a\n\ntype A string\n", "b/b.go": deep}, true)
	_, result, err := goImportResults(t, snapshot, goImportSettings(t))
	if err != nil {
		t.Fatal(err)
	}
	coverage, _ := result.Report(n.Coverage())
	if ds := result.Diagnostics(); coverage.Status != semantic.Unsupported || len(ds) == 0 || ds[0].At == nil || ds[0].At.Tree != trees["a/a_test.go"].ID {
		t.Fatal("expected the test file refused before the deep tree's limit", coverage, ds)
	}
	// Case collisions come before disagreeing package clauses, whatever the paths.
	snapshot, _ = goPackageCapture(t, map[string]string{
		"B/x.go": "package b\n\ntype X string\n",
		"b/x.go": "package b\n\ntype Y string\n",
		"a/y.go": "package a\n\ntype Y string\n",
		"a/z.go": "package q\n\ntype Z string\n",
	}, true)
	_, result, err = goImportResults(t, snapshot, goImportSettings(t))
	if err != nil {
		t.Fatal(err)
	}
	if ds := result.Diagnostics(); len(ds) == 0 || !strings.Contains(ds[0].Message, "case-insensitive import collision") {
		t.Fatal("expected the case collision first", ds)
	}
	// Vendor directories are outside the profile, and supplied paths may not name them.
	snapshot, _ = goPackageCapture(t, map[string]string{"vendor/v/v.go": "package v\n\ntype V string\n"}, true)
	_, result, err = goImportResults(t, snapshot, goImportSettings(t))
	if err != nil {
		t.Fatal(err)
	}
	if coverage, _ := result.Report(n.Coverage()); coverage.Status != semantic.Unsupported {
		t.Fatal("a vendored captured file was accepted", coverage)
	}
	vendored := n.Package{ImportPath: "example.org/lib/vendor/v", Files: []n.PackageFile{{Name: "v.go", Text: "package v\n"}}}
	if _, err := n.ImportConfiguration(n.ImportEnvironment{Module: goImportModule, Packages: []n.Package{vendored}}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("a vendored supplied path was accepted", err)
	}
}
