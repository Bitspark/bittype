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

var nativeStages = []string{"syntax", "parse/syntax", "check/syntax", "format/syntax", "provider"}

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
	for _, stage := range append(slices.Clone(nativeStages), "resolve/files", "analyze/files") {
		entries, err := os.ReadDir(filepath.Join("../..", stage))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() == "files" {
				continue
			}
			if stage == "check/syntax" && (entry.Name() == "ast" || entry.Name() == "cst") {
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
	semanticQualified := map[string]bool{}
	for _, fixture := range semanticConformanceCases() {
		if semanticQualified[fixture.directory] || fixture.qualify == nil {
			t.Fatalf("duplicate or unqualified semantic/text fixture: %s", fixture.directory)
		}
		qualified[fixture.directory] = true
		semanticQualified[fixture.directory] = true
		profiles[fixture.directory] = []string{"syntax", "parse/syntax", "resolve/files", "check/syntax", "analyze/files", "format/syntax", "provider"}
	}
	for _, name := range languages {
		for _, stage := range profiles[name] {
			path := modulePath + "/" + stage + "/" + name + "/go"
			if p, ok := packages[path]; !ok || len(p.GoFiles) == 0 {
				t.Errorf("missing implemented stage %s", path)
			}
		}
		for _, stage := range append(slices.Clone(nativeStages), "resolve/files", "analyze/files") {
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
		"tools/provider-matrix/go":  {}, // development tooling, standard library only
		"tools/horos-module/go":     {"contracts/parse/go", "internal/horosmodel", "internal/lex", "parse/syntax/horos/go", "trees/ast/go"},
		"trees/ast/go":              {},
		"semantics/go":              {"trees/ast/go"},
		"trees/cst/go":              {"trees/ast/go"},
		"contracts/language/go":     {},
		"contracts/pipeline/go":     {},
		"contracts/check/go":        {"trees/representation/go", "trees/ast/go", "contracts/language/go"},
		"contracts/parse/go":        {"trees/representation/go", "trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"contracts/format/go":       {"trees/representation/go", "trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"contracts/print/go":        {"trees/representation/go", "trees/cst/go", "contracts/language/go"},
		"contracts/text/go":         {"trees/ast/go"},
		"check/syntax/ast/go":       {"trees/ast/go", "contracts/text/go"},
		"check/syntax/cst/go":       {"trees/ast/go", "trees/cst/go", "check/syntax/ast/go", "contracts/text/go"},
		"print/syntax/go":           {"trees/cst/go", "check/syntax/cst/go", "contracts/language/go", "contracts/print/go", "trees/representation/go"},
		"trees/representation/go":   {"trees/ast/go", "trees/cst/go", "contracts/language/go", "check/syntax/ast/go", "check/syntax/cst/go", "internal/tree"},
		"contracts/transform/go":    {"trees/ast/go", "trees/representation/go", "contracts/language/go", "contracts/semantic/go", "internal/tree"},
		"contracts/semantic/go":     {"contracts/check/go", "trees/ast/go", "trees/representation/go", "contracts/language/go", "semantics/go", "internal/tree"},
		"contracts/verification/go": {"contracts/semantic/go", "semantics/go", "trees/ast/go", "trees/representation/go"},
		"contracts/resolution/go":   {"trees/representation/go", "contracts/language/go", "contracts/semantic/go", "semantics/go"},
		"contracts/analysis/go":     {"trees/representation/go", "contracts/language/go", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"},
		"engine/go":                 {"trees/ast/go", "trees/representation/go", "contracts/language/go", "contracts/parse/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", "contracts/transform/go", "contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go"},
		"internal/tree":             {"trees/ast/go"},
		"internal/grammar":          {"trees/ast/go", "check/syntax/ast/go", "contracts/text/go"},
		"internal/lex":              {"trees/ast/go", "contracts/parse/go", "contracts/text/go"},
		"internal/parseutil":        {"trees/ast/go", "check/syntax/ast/go", "contracts/parse/go", "contracts/text/go", "internal/lex"},
		"internal/concrete":         {"trees/ast/go", "trees/cst/go", "check/syntax/cst/go", "contracts/parse/go", "contracts/text/go", "internal/lex", "internal/tree"},
	}
	for _, name := range languages {
		syntax := "syntax/" + name + "/go"
		checker := "check/syntax/" + name + "/go"
		helper := "internal/" + name
		rules[syntax] = []string{"trees/ast/go", "trees/cst/go", "contracts/language/go"}
		base := []string{"trees/ast/go", "trees/cst/go", "contracts/language/go", "contracts/text/go", syntax}
		for _, stage := range []string{"check", "parse", "format"} {
			allowed := append([]string{}, base...)
			allowed = append(allowed, "check/syntax/ast/go", "internal/tree", "internal/grammar", helper)
			allowed = append(allowed, "contracts/"+stage+"/go")
			if stage != "check" {
				allowed = append(allowed, checker, "internal/concrete")
			}
			if stage == "parse" {
				allowed = append(allowed, "internal/lex", "internal/parseutil")
			}
			rules[stage+"/syntax/"+name+"/go"] = allowed
		}
		rules["provider/"+name+"/go"] = append(append([]string{}, base...), "contracts/parse/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", checker, "parse/syntax/"+name+"/go", "format/syntax/"+name+"/go", "print/syntax/go", "print/files/go", "trees/representation/go")
		// Language-specific lexical helpers cannot depend on stages or another language.
		rules[helper] = []string{"trees/ast/go", "contracts/text/go", "internal/grammar", "internal/lex", "internal/tree"}
	}
	// The Horos frontend uses private mechanics, not a language registry.
	rules["internal/declsyntax"] = []string{"trees/ast/go", "check/syntax/ast/go", "contracts/text/go", "internal/grammar"}
	rules["internal/horosmodel"] = []string{"trees/ast/go", "check/syntax/ast/go", "contracts/text/go", "internal/declsyntax", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"}
	for _, name := range []string{"horos"} {
		rules["resolve/files/"+name+"/go"] = []string{"trees/ast/go", "syntax/" + name + "/go", "check/syntax/" + name + "/go", "internal/horosmodel", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"}
		rules["check/syntax/"+name+"/go"] = []string{"trees/ast/go", "syntax/" + name + "/go", "internal/declsyntax", "contracts/check/go"}
		rules["analyze/files/"+name+"/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go", "semantics/go", "syntax/" + name + "/go", "resolve/files/" + name + "/go"}
		rules["provider/"+name+"/go"] = append(rules["provider/"+name+"/go"], []string{"trees/ast/go", "contracts/check/go", "contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go", "syntax/" + name + "/go", "check/syntax/" + name + "/go", "analyze/files/" + name + "/go", "resolve/files/" + name + "/go"}...)
		for _, stage := range []string{"parse", "format"} {
			rules[stage+"/syntax/"+name+"/go"] = append(rules[stage+"/syntax/"+name+"/go"], "internal/horosmodel", "internal/lex")
		}
	}
	rules["internal/horosgen"] = []string{"analyze/files/horos/go", "contracts/analysis/go", "contracts/semantic/go", "contracts/text/go", "contracts/transform/go", "resolve/files/horos/go", "semantics/go", "syntax/horos/go", "trees/ast/go"}
	for target, syntax := range map[string]string{"typescript": "typescript", "go": "golang", "rust": "rust", "cpp": "cpp"} {
		rules["transforms/horos-to-"+target+"/go"] = []string{"contracts/transform/go", "contracts/text/go", "internal/horosgen", "internal/" + syntax, "syntax/horos/go", "syntax/" + syntax + "/go", "trees/ast/go", "trees/representation/go"}
	}
	// The Files dimension has explicit stage paths and a separate sealed carrier.
	rules["files/go"] = []string{"trees/ast/go"}
	rules["source/go"] = []string{"contracts/semantic/go", "files/go", "semantics/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"}
	rules["trees/representation/go"] = append(rules["trees/representation/go"], "semantics/go")
	rules["contracts/semantic/go"] = append(rules["contracts/semantic/go"], "files/go")
	for _, stage := range []string{"parse", "format", "print"} {
		key := "contracts/" + stage + "/go"
		rules[key] = append(rules[key], "files/go", "source/go", "semantics/go", "contracts/semantic/go")
	}
	rules["engine/go"] = append(rules["engine/go"], "files/go", "source/go")
	rules["print/files/go"] = []string{"contracts/language/go", "contracts/print/go", "files/go", "source/go", "trees/cst/go", "trees/representation/go"}
	composition := []string{"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/parse/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go", "files/go", "source/go", "semantics/go", "syntax/horos/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"}
	// Files composition uses the public Horos prelude identity from the resolver.
	rules["internal/horosfiles"] = append(append([]string{}, composition...), "internal/tree", "resolve/files/horos/go")
	for _, stage := range []string{"parse", "format"} {
		rules[stage+"/files/horos/go"] = append(append([]string{}, composition...), "internal/horosfiles")
	}
	rules["provider/horos/go"] = append(rules["provider/horos/go"], "parse/files/horos/go", "format/files/horos/go", "print/files/go", "trees/representation/go")
	rules["internal/horosgen"] = append(rules["internal/horosgen"], "trees/representation/go")
	// Verification projections are explicit concrete compositions, not registered
	// language stages. The agnostic association above never imports their models.
	rules["verification/files/source/horos/go"] = []string{"contracts/semantic/go", "contracts/resolution/go", "internal/horosmodel", "resolve/files/horos/go", "syntax/horos/go", "trees/ast/go", "trees/representation/go"}
	rules["semantics/golang/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "semantics/go", "syntax/golang/go", "trees/ast/go", "trees/representation/go"}
	rules["internal/golangsem"] = []string{"contracts/semantic/go", "semantics/go", "semantics/golang/go", "syntax/golang/go", "trees/ast/go"}
	rules["resolve/files/golang/go"] = []string{"contracts/resolution/go", "contracts/semantic/go", "internal/golangsem", "semantics/go", "semantics/golang/go", "syntax/golang/go"}
	rules["analyze/files/golang/go"] = []string{"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "internal/golangsem", "resolve/files/golang/go", "semantics/go", "semantics/golang/go", "syntax/golang/go", "trees/representation/go"}
	rules["provider/golang/go"] = append(rules["provider/golang/go"], "contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "analyze/files/golang/go", "resolve/files/golang/go")
	rules["verification/files/horos-to-go/go"] = []string{"contracts/semantic/go", "contracts/verification/go", "semantics/go", "trees/ast/go", "verification/files/source/horos/go", "semantics/golang/go", "analyze/files/golang/go"}
	for path, p := range packages {
		if p.Standard || path == modulePath+"/integration/go" || path == modulePath+"/verification/integration/go" {
			continue
		}
		if !strings.HasPrefix(path, modulePath+"/") {
			if !slices.Contains(modelDependencies, path) {
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
			if slices.Contains(modelExternalImports[relative], dependency) {
				continue
			}
			if !strings.HasPrefix(dependency, modulePath+"/") || !slices.Contains(allowed, strings.TrimPrefix(dependency, modulePath+"/")) {
				t.Errorf("%s must not import %s", relative, dependency)
			}
		}
	}
}

// Only owning frontends/compositions consume external semantic models. A global exemption for
// external modules would let model dependencies leak into generic infrastructure.
var modelDependencies = []string{
	"github.com/Bitspark/horos/core/go",
	"github.com/Bitspark/horos/core/go/declarations",
	"github.com/Bitspark/horos/core/go/internal/declmodel",
	"github.com/Bitspark/horos/core/go/internal/kindmodel",
	"github.com/bitspark/ontos/core/go",
	"github.com/bitspark/ontos/data/go",
	"github.com/bitspark/ontos/codec/go",
	"github.com/Bitspark/logos-contract/contract/go/typed",
}

var modelExternalImports = map[string][]string{
	"tools/horos-module/go":              {"github.com/bitspark/ontos/codec/go"},
	"internal/horosgen":                  {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"parse/syntax/horos/go":              {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"format/syntax/horos/go":             {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"internal/horosmodel":                {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"resolve/files/horos/go":             {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
	"verification/files/source/horos/go": {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"verification/files/horos-to-go/go":  {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/Bitspark/logos-contract/contract/go/typed"},
}

func TestFilesDimensionStageCoverage(t *testing.T) {
	// Files is a dimension, never a discovered language named "files".
	for _, stage := range []string{"parse", "format"} {
		entries, err := os.ReadDir(filepath.Join("../..", stage, "files"))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if entry.IsDir() && entry.Name() != "horos" {
				t.Fatalf("%s Files semantics lack a documented fixture", entry.Name())
			}
		}
	}
}
