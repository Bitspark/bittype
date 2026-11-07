package report

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
	observation "github.com/Bitspark/bittype/tools/provider-matrix/observation/go"
)

func TestOverviewScopesProfilesAndCaseDeduplication(t *testing.T) {
	p := fixture(t)
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(p.root, "docs"), 0700); err != nil {
		t.Fatal(err)
	}
	scope := "# Scope\n<!-- bittype:overview scope=scope part=features -->\nEmber forms.\n<!-- bittype:overview:end -->\n<!-- bittype:overview scope=scope part=exclusions -->\nBodies excluded.\n<!-- bittype:overview:end -->\n## Foreign\n<!-- bittype:overview scope=foreign part=features -->\nForeign-only forms.\n<!-- bittype:overview:end -->\n"
	if err := os.WriteFile(filepath.Join(p.root, "docs/scope.md"), []byte(scope), 0600); err != nil {
		t.Fatal(err)
	}
	target := capability.Symbol{Package: p.module + "/print/syntax/ember/provider/go", Name: "emit"}
	a := &annotations{operations: map[capability.Symbol][]authored{target: {
		{directive: directive{name: "scope"}, document: "docs/scope.md#scope"},
		{directive: directive{name: "scope", profile: profile{present: true, value: "strict"}}, document: "docs/scope.md#scope"},
		{directive: directive{name: "scope"}, document: "docs/scope.md"},
	}}, cases: []caseAssociation{
		{target: target, directive: directive{values: map[string]string{"feature": "emit", "case": "one", "expects": "owner-assertion"}}, test: testDeclaration{id: "test-id", name: "TestEmit", file: "emit_test.go", variant: "internal"}},
	}}
	a.cases = append(a.cases, a.cases[0]) // repeated association, not another case
	second := a.cases[0]
	second.directive.profile = profile{present: true, value: "strict"}
	a.cases = append(a.cases, second)
	if err := associateCases(m, a); err != nil {
		t.Fatal(err)
	}
	views := collectOverview(m, a)
	if len(views) != 2 || len(views[0].tests) != 1 || len(views[0].cases) != 2 || len(views[1].cases) != 0 {
		t.Fatalf("misassociated evidence: %#v", views)
	}
	config := Config{MatrixPath: "docs/providers.md", DetailsPath: "docs/provider-capabilities.md", OverviewPath: "guide/languages.md"}
	render := func() string {
		s, err := p.renderOverview(context.Background(), m, a, config, "https://example.test/source", nil)
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	s := render()
	if s != render() {
		t.Fatal("overview is not deterministic")
	}
	for _, want := range []string{"../docs/providers.md", "../docs/provider-capabilities.md", "Source → target", "| Dimension |", "Bodies excluded.", "profile: strict", "authored profile absent", "1 associated test declarations, 2 authored cases", "no execution results", "Delegate evidence does not qualify wrappers"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q", want)
		}
	}
	if strings.Contains(s, "Foreign-only") || strings.Contains(s, "**pass**") {
		t.Fatal("scope leak or manufactured outcome")
	}
	if strings.Count(s, "Supported features:** Ember forms.") != 2 {
		t.Fatal("file-wide scope acted as a wildcard or profile was erased")
	}
	ice := s[strings.Index(s, "## Ice"):]
	if strings.Contains(ice, "Ember forms") || !strings.Contains(ice, "Supported breadth and exclusions are unreported") {
		t.Fatal("scope transferred to unrelated language")
	}
}

func TestOverviewMatchedOccurrencesKeepFailuresAndApplicability(t *testing.T) {
	ctx := observation.SourceContext{Revision: "rev", ModulePath: "example.test/reports", ModuleRevision: "rev", Harness: "repository", GoVersion: "go1.27.3", GOOS: "linux", GOARCH: "amd64", CGOEnabled: "1", InputState: "clean", Complete: true, Modules: []observation.ModuleSelection{{Path: "example.test/reports", Revision: "rev"}}}
	a := &annotations{tests: []testDeclaration{{id: "test", pkg: "example.test/reports", name: "TestOne", variant: "external", file: "report_test.go"}}}
	views := []languageOverview{{language: language{name: "Ember"}, tests: map[string]bool{"test": true}}, {language: language{name: "Ice"}, tests: map[string]bool{}}}
	failure := 1
	execution := observation.Execution{Invocation: observation.Invocation{Schema: observation.Schema, Context: ctx, RunID: "17", Attempt: "1", Job: "checks", CommandID: "tests", Args: []string{"go", "test", "-json", "-count=2", "./..."}, CaptureComplete: true, ExitCode: &failure}, ProtocolComplete: true, Tests: []observation.TestRun{
		{Package: "example.test/reports", Name: "TestOne", Occurrence: 1, Started: true, Outcome: observation.Failed},
		{Package: "example.test/reports", Name: "TestOne", Occurrence: 2, Started: true, Outcome: observation.Passed},
		{Package: "example.test/reports", Name: "TestOne/child", Occurrence: 1, Started: true, Outcome: observation.Passed},
	}}
	config := Config{Context: &ctx, Executions: []observation.Execution{execution}}
	var b limitedText
	renderOverviewObservations(&b, a, config, views)
	for _, want := range []string{"process exit: **1**", "| Language | pass | fail | skip | incomplete |", "| Ember | 1 | 1 | 0 | 0 |", "| Ice | / | / | / | / |"} {
		if !strings.Contains(b.String(), want) {
			t.Fatalf("missing %q: %s", want, b.String())
		}
	}
	config.Executions[0].Invocation.Context.Harness = "installed-consumer"
	b = limitedText{}
	renderOverviewObservations(&b, a, config, views)
	if !strings.Contains(b.String(), "mismatched") || !strings.Contains(b.String(), "| Ember | / | / | / | / |") {
		t.Fatal("foreign harness credited", b.String())
	}
}

func TestOverviewOutputPathsMustRemainDistinct(t *testing.T) {
	for _, path := range []string{"docs/providers.md", "docs/provider-capabilities.md", "../escape.md"} {
		_, err := Generate(context.Background(), Config{Root: t.TempDir(), OverviewPath: path})
		if err == nil || !strings.Contains(err.Error(), "paths") && !strings.Contains(err.Error(), "escapes") {
			t.Fatalf("path %q was not rejected before loading source: %v", path, err)
		}
	}
}

func TestSharedScopeFoldingPreservesProfilesAndLanguageAssociation(t *testing.T) {
	p := fixture(t)
	if err := os.WriteFile(filepath.Join(p.root, "scope.md"), []byte("# Shared\n<!-- bittype:overview scope=shared part=features -->\nCommon forms.\n<!-- bittype:overview:end -->\n"), 0600); err != nil {
		t.Fatal(err)
	}
	entry := authored{directive: directive{name: "scope"}, document: "scope.md#shared"}
	named := entry
	named.directive.profile = profile{present: true, value: "named"}
	views := []languageOverview{
		{language: language{id: "ember"}, bindings: []overviewBinding{{entry: entry}, {entry: entry}, {entry: named}}},
		{language: language{id: "ice"}, bindings: []overviewBinding{{entry: entry}}},
	}
	scopes, err := p.overviewScopes(views, "https://example.test/source")
	if err != nil || len(scopes) != 2 {
		t.Fatal(scopes, err)
	}
	for _, scope := range scopes {
		if scope.entry.directive.profile.present {
			if len(scope.languages) != 1 || scope.languages[0] != "ember" {
				t.Fatal("named profile broadened", scope)
			}
		} else if len(scope.languages) != 2 {
			t.Fatal("common scope not deduplicated", scope)
		}
		if strings.Count(scope.text, "Common forms.") != 1 || !strings.Contains(scope.text, "No compact exclusion passage") {
			t.Fatal("scope content or missing evidence lost", scope)
		}
	}
}
