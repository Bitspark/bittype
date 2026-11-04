package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var nativeStages = []string{"syntax", "parse", "check", "format", "provider"}

type productionPackage struct {
	ImportPath string
	Standard   bool
	Imports    []string
	GoFiles    []string
}

func productionPackages(t *testing.T) map[string]productionPackage {
	t.Helper()
	command := exec.Command("go", "list", "-deps", "-json", "./...")
	command.Dir = "../.."
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	result := map[string]productionPackage{}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var p productionPackage
		if err := decoder.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		result[p.ImportPath] = p
	}
	return result
}

// Discover the stage matrix from source directories, not a production language list.
// The union catches an orphan parser as well as an incomplete syntax implementation.
func stageLanguages(t *testing.T) []string {
	t.Helper()
	names := map[string]bool{}
	for _, stage := range append(slices.Clone(nativeStages), "resolve", "analyze") {
		entries, err := os.ReadDir(filepath.Join("../..", stage))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if stage == "check" && (entry.Name() == "ast" || entry.Name() == "cst") {
				continue
			}
			names[entry.Name()] = true
		}
	}
	var result []string
	for name := range names {
		result = append(result, name)
	}
	slices.Sort(result)
	return result
}

func TestLanguageStageMatrix(t *testing.T) {
	packages := productionPackages(t)
	languages := stageLanguages(t)
	if len(languages) == 0 {
		t.Fatal("no language implementations")
	}
	qualified := map[string]bool{}
	profiles := map[string][]string{}
	for _, fixture := range nativeConformanceCases() {
		if qualified[fixture.directory] {
			t.Fatalf("duplicate conformance fixture: %s", fixture.directory)
		}
		qualified[fixture.directory] = true
		profiles[fixture.directory] = nativeStages
	}
	for _, fixture := range semanticConformanceCases() {
		if qualified[fixture.directory] || fixture.qualify == nil {
			t.Fatalf("duplicate or unqualified semantic/text fixture: %s", fixture.directory)
		}
		qualified[fixture.directory] = true
		profiles[fixture.directory] = []string{"syntax", "parse", "resolve", "check", "analyze", "format", "provider"}
	}
	for _, name := range languages {
		for _, stage := range profiles[name] {
			path := modulePath + "/" + stage + "/" + name + "/go"
			if p, ok := packages[path]; !ok || len(p.GoFiles) == 0 {
				t.Errorf("missing implemented stage %s", path)
			}
		}
		for _, stage := range append(slices.Clone(nativeStages), "resolve", "analyze") {
			if _, ok := packages[modulePath+"/"+stage+"/"+name+"/go"]; ok && !slices.Contains(profiles[name], stage) {
				t.Errorf("%s has an unqualified capability stage: %s", name, stage)
			}
		}
		if !qualified[name] {
			t.Errorf("%s is missing a tested native-text or semantic/text capability profile", name)
		}
		delete(qualified, name)
	}
	for name := range qualified {
		t.Errorf("conformance fixture has no language implementation: %s", name)
	}
}

// Each production package is checked, including private helpers. This prevents
// a forbidden dependency from being hidden behind an otherwise allowed helper.
func TestProductionArchitecture(t *testing.T) {
	packages := productionPackages(t)
	languages := stageLanguages(t)
	rules := map[string][]string{
		"trees/ast/go":            {},
		"semantics/go":            {"trees/ast/go"},
		"trees/cst/go":            {"trees/ast/go"},
		"contracts/language/go":   {},
		"contracts/pipeline/go":   {},
		"contracts/check/go":      {"trees/representation/go", "trees/ast/go", "contracts/language/go"},
		"contracts/parse/go":      {"trees/representation/go", "trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"contracts/format/go":     {"trees/representation/go", "trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"contracts/print/go":      {"trees/representation/go", "trees/cst/go", "contracts/language/go"},
		"contracts/text/go":       {"trees/ast/go"},
		"check/ast/go":            {"trees/ast/go", "contracts/text/go"},
		"check/cst/go":            {"trees/ast/go", "trees/cst/go", "check/ast/go", "contracts/text/go"},
		"print/go":                {"trees/cst/go", "check/cst/go"},
		"trees/representation/go": {"trees/ast/go", "trees/cst/go", "contracts/language/go", "check/ast/go", "check/cst/go", "internal/tree"},
		"contracts/transform/go":  {"trees/ast/go", "trees/representation/go", "contracts/language/go", "contracts/pipeline/go"},
		"contracts/semantic/go":   {"contracts/check/go", "trees/ast/go", "trees/representation/go", "contracts/language/go", "semantics/go", "internal/tree"},
		"contracts/resolution/go": {"trees/representation/go", "contracts/language/go", "contracts/semantic/go", "semantics/go"},
		"contracts/analysis/go":   {"trees/representation/go", "contracts/language/go", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"},
		"engine/go":               {"trees/representation/go", "contracts/language/go", "contracts/parse/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", "contracts/transform/go", "contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go"},
		"internal/tree":           {"trees/ast/go"},
		"internal/grammar":        {"trees/ast/go", "check/ast/go", "contracts/text/go"},
		"internal/lex":            {"trees/ast/go", "contracts/parse/go", "contracts/text/go"},
		"internal/parseutil":      {"trees/ast/go", "check/ast/go", "contracts/parse/go", "contracts/text/go", "internal/lex"},
		"internal/concrete":       {"trees/ast/go", "trees/cst/go", "check/cst/go", "contracts/parse/go", "contracts/text/go", "internal/lex", "internal/tree"},
	}
	for _, name := range languages {
		syntax := "syntax/" + name + "/go"
		checker := "check/" + name + "/go"
		helper := "internal/" + name
		rules[syntax] = []string{"trees/ast/go", "trees/cst/go", "contracts/language/go"}
		base := []string{"trees/ast/go", "trees/cst/go", "contracts/language/go", "contracts/text/go", syntax}
		for _, stage := range []string{"check", "parse", "format"} {
			allowed := append([]string{}, base...)
			allowed = append(allowed, "check/ast/go", "internal/tree", "internal/grammar", helper)
			allowed = append(allowed, "contracts/"+stage+"/go")
			if stage != "check" {
				allowed = append(allowed, checker, "internal/concrete")
			}
			if stage == "parse" {
				allowed = append(allowed, "internal/lex", "internal/parseutil")
			}
			rules[stage+"/"+name+"/go"] = allowed
		}
		rules["provider/"+name+"/go"] = append(append([]string{}, base...), "contracts/parse/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", checker, "parse/"+name+"/go", "format/"+name+"/go", "print/go")
		// Language-specific lexical helpers cannot depend on stages or another language.
		rules[helper] = []string{"trees/ast/go", "contracts/text/go", "internal/grammar", "internal/lex", "internal/tree"}
	}
	// The Horos frontend uses private mechanics, not a language registry.
	rules["internal/declsyntax"] = []string{"trees/ast/go", "check/ast/go", "contracts/text/go", "internal/grammar"}
	rules["internal/horosmodel"] = []string{"trees/ast/go", "check/ast/go", "contracts/text/go", "internal/declsyntax", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"}
	for _, name := range []string{"horos"} {
		rules["resolve/"+name+"/go"] = []string{"trees/ast/go", "syntax/" + name + "/go", "check/" + name + "/go", "internal/horosmodel", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"}
		rules["check/"+name+"/go"] = []string{"trees/ast/go", "syntax/" + name + "/go", "internal/declsyntax"}
		rules["analyze/"+name+"/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go", "semantics/go", "syntax/" + name + "/go", "resolve/" + name + "/go"}
		rules["provider/"+name+"/go"] = append(rules["provider/"+name+"/go"], []string{"trees/ast/go", "contracts/check/go", "contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go", "syntax/" + name + "/go", "check/" + name + "/go", "analyze/" + name + "/go", "resolve/" + name + "/go"}...)
		for _, stage := range []string{"parse", "format"} { rules[stage+"/"+name+"/go"] = append(rules[stage+"/"+name+"/go"], "internal/horosmodel", "internal/lex") }
	}
	for path, p := range packages {
		if p.Standard || path == modulePath+"/integration/go" {
			continue
		}
		if !strings.HasPrefix(path, modulePath+"/") {
			if !slices.Contains(horosModelDependencies, path) {
				t.Errorf("unclassified external production dependency: %s", path)
			}
			continue
		}
		relative := strings.TrimPrefix(path, modulePath+"/")
		allowed, ok := rules[relative]
		if !ok {
			t.Errorf("unclassified production package: %s", relative)
			continue
		}
		for _, dependency := range p.Imports {
			if dep, ok := packages[dependency]; ok && dep.Standard {
				continue
			}
			if slices.Contains(horosExternalImports[relative], dependency) {
				continue
			}
			if !strings.HasPrefix(dependency, modulePath+"/") || !slices.Contains(allowed, strings.TrimPrefix(dependency, modulePath+"/")) {
				t.Errorf("%s must not import %s", relative, dependency)
			}
		}
	}
}

// Only the Horos frontend consumes the semantic model. A global exemption for
// external modules would let model dependencies leak into generic infrastructure.
var horosModelDependencies = []string{
	"github.com/Bitspark/horos/core/go",
	"github.com/Bitspark/horos/core/go/declarations",
	"github.com/Bitspark/horos/core/go/internal/declmodel",
	"github.com/bitspark/ontos/core/go",
	"github.com/bitspark/ontos/data/go",
	"github.com/bitspark/ontos/codec/go",
}

var horosExternalImports = map[string][]string{
	"parse/horos/go": {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"format/horos/go": {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"internal/horosmodel": {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"resolve/horos/go":    {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
}
