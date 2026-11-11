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

var nativeStages = []string{"syntax", "parse/syntax", "check/syntax", "format/syntax"}

type productionPackage struct {
	ImportPath string
	Name       string
	Dir        string
	Standard   bool
	Imports    []string
	GoFiles    []string
	CgoFiles   []string
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
			if !entry.IsDir() || entry.Name() == "shared" {
				continue
			}
			files, err := filepath.Glob(filepath.Join("../..", strings.TrimPrefix(stagePath(stage, entry.Name()), modulePath+"/"), "*.go"))
			if err != nil {
				t.Fatal(err)
			}
			if len(files) > 0 {
				names[entry.Name()] = true
			}
		}
	}
	var result []string
	for name := range names {
		result = append(result, name)
	}
	slices.Sort(result)
	return result
}

func stagePath(stage, name string) string {
	if stage == "syntax" {
		return modulePath + "/syntax/" + name + "/go"
	}
	return modulePath + "/" + stage + "/" + name + "/provider/go"
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
		profiles[fixture.directory] = []string{"syntax", "parse/syntax", "resolve/files", "check/syntax", "analyze/files", "format/syntax"}
	}
	for _, name := range languages {
		for _, stage := range profiles[name] {
			path := stagePath(stage, name)
			if p, ok := packages[path]; !ok || len(p.GoFiles) == 0 {
				t.Errorf("missing implemented stage %s", path)
			}
		}
		for _, stage := range append(slices.Clone(nativeStages), "resolve/files", "analyze/files") {
			if _, ok := packages[stagePath(stage, name)]; ok && !slices.Contains(profiles[name], stage) {
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
//
// A package is classified by its entry in rules or, without one, by the role
// its path gives (packageRole). Every allowance in rules must be used. An
// import of a stage provider package also needs an entry in providerExceptions
// that names the symbols used. The rules themselves are in checkArchitecture.
func TestProductionArchitecture(t *testing.T) {
	packages := productionPackages(t)
	languages := stageLanguages(t)
	rules := map[string][]string{
		"tools/provider-matrix/go":                 {"tools/provider-matrix/report/go", "tools/provider-matrix/observation/go", "tools/provider-matrix/extract/grammar/go"},
		"tools/provider-matrix/report/go":          {"tools/provider-matrix/capability/go", "tools/provider-matrix/observation/go"},
		"tools/provider-matrix/capability/go":      {},
		"tools/provider-matrix/observation/go":     {},
		"tools/provider-matrix/extract/grammar/go": {"tools/provider-matrix/capability/go"},
		"tools/horos-module/go":                    {"contracts/parse/go", "internal/horosmodel", "internal/lex", "parse/syntax/horos/provider/go", "trees/ast/go"},
		"trees/ast/go":                             {},
		"semantics/go":                             {"trees/ast/go"},
		"trees/cst/go":                             {"trees/ast/go"},
		"contracts/language/go":                    {},
		"contracts/pipeline/go":                    {},
		"contracts/check/go":                       {"trees/representation/go", "trees/ast/go", "contracts/language/go"},
		"contracts/parse/go":                       {"trees/representation/go", "trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"contracts/format/go":                      {"trees/representation/go", "trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"contracts/print/go":                       {"trees/representation/go", "trees/cst/go", "contracts/language/go"},
		"contracts/text/go":                        {"trees/ast/go"},
		"check/syntax/shared/ast/go":               {"trees/ast/go", "contracts/text/go"},
		"check/syntax/shared/cst/go":               {"trees/cst/go", "check/syntax/shared/ast/go", "contracts/text/go"},
		"print/syntax/shared/provider/go":          {"trees/cst/go", "check/syntax/shared/cst/go", "contracts/language/go", "contracts/print/go", "trees/representation/go"},
		"trees/representation/go":                  {"trees/ast/go", "trees/cst/go", "contracts/language/go", "check/syntax/shared/ast/go", "check/syntax/shared/cst/go", "internal/tree"},
		"contracts/transform/go":                   {"trees/ast/go", "trees/representation/go", "contracts/language/go", "contracts/semantic/go", "internal/tree"},
		"contracts/semantic/go":                    {"contracts/check/go", "trees/ast/go", "trees/representation/go", "contracts/language/go", "semantics/go", "internal/tree"},
		"contracts/verification/go":                {"contracts/semantic/go", "semantics/go", "trees/ast/go", "trees/representation/go"},
		"contracts/generation/go":                  {"contracts/transform/go", "trees/ast/go", "contracts/semantic/go", "semantics/go", "trees/representation/go"},
		"contracts/resolution/go":                  {"trees/representation/go", "contracts/language/go", "contracts/semantic/go", "semantics/go"},
		"contracts/analysis/go":                    {"trees/representation/go", "contracts/language/go", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"},
		"engine/go":                                {"trees/representation/go", "contracts/language/go", "contracts/parse/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", "contracts/transform/go", "contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go"},
		"internal/tree":                            {"trees/ast/go"},
		"internal/grammar":                         {"trees/ast/go", "check/syntax/shared/ast/go", "contracts/text/go"},
		"internal/lex":                             {"trees/ast/go", "contracts/parse/go", "contracts/text/go"},
		"internal/parseutil":                       {"trees/ast/go", "check/syntax/shared/ast/go", "contracts/parse/go", "contracts/text/go", "internal/lex"},
		"internal/concrete":                        {"trees/ast/go", "trees/cst/go", "check/syntax/shared/cst/go", "contracts/parse/go", "contracts/text/go", "internal/lex", "internal/tree"},
	}
	for _, name := range languages {
		syntax := "syntax/" + name + "/go"
		checker := "check/syntax/" + name + "/provider/go"
		rules[syntax] = []string{"trees/ast/go", "trees/cst/go", "contracts/language/go"}
		base := []string{"trees/ast/go", syntax, "internal/tree"}
		for _, stage := range []string{"check", "parse", "format"} {
			allowed := append([]string{}, base...)
			allowed = append(allowed, "contracts/"+stage+"/go")
			if stage == "check" {
				allowed = append(allowed, "internal/grammar")
			} else {
				allowed = append(allowed, "trees/cst/go", checker, "internal/concrete")
			}
			if stage != "format" {
				allowed = append(allowed, "contracts/text/go")
			}
			if stage == "parse" {
				allowed = append(allowed, "internal/lex", "internal/parseutil")
			}
			rules[stage+"/syntax/"+name+"/provider/go"] = allowed
		}
	}
	// Language-specific lexical helpers cannot depend on stages or another language.
	// Go has none. Each stage is allowed its helper only where it uses it.
	for _, name := range []string{"cpp", "rust", "typescript"} {
		helper := "internal/" + name
		rules[helper] = []string{"internal/grammar"}
		for _, stage := range []string{"check", "parse"} {
			rules[stage+"/syntax/"+name+"/provider/go"] = append(rules[stage+"/syntax/"+name+"/provider/go"], helper)
		}
	}
	rules["format/syntax/typescript/provider/go"] = append(rules["format/syntax/typescript/provider/go"], "internal/grammar", "internal/typescript")
	rules["internal/horos"] = []string{}
	// The Horos frontend uses private mechanics, not a language registry.
	rules["internal/declsyntax"] = []string{"trees/ast/go", "check/syntax/shared/ast/go", "contracts/text/go", "internal/grammar"}
	rules["internal/horosmodel"] = []string{"trees/ast/go", "check/syntax/shared/ast/go", "contracts/text/go", "internal/declsyntax", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"}
	for _, name := range []string{"horos"} {
		rules["resolve/files/"+name+"/provider/go"] = []string{"trees/ast/go", "syntax/" + name + "/go", "check/syntax/" + name + "/provider/go", "internal/horosmodel", "contracts/semantic/go", "contracts/resolution/go", "semantics/go"}
		rules["check/syntax/"+name+"/provider/go"] = []string{"trees/ast/go", "syntax/" + name + "/go", "internal/declsyntax", "contracts/check/go"}
		rules["analyze/files/"+name+"/provider/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "contracts/resolution/go", "semantics/go", "syntax/" + name + "/go", "resolve/files/" + name + "/provider/go"}
		for _, stage := range []string{"parse", "format"} {
			rules[stage+"/syntax/"+name+"/provider/go"] = append(rules[stage+"/syntax/"+name+"/provider/go"], "internal/"+name, "internal/horosmodel", "internal/lex")
		}
		rules["parse/syntax/"+name+"/provider/go"] = append(rules["parse/syntax/"+name+"/provider/go"], "check/syntax/shared/ast/go")
		rules["format/syntax/"+name+"/provider/go"] = append(rules["format/syntax/"+name+"/provider/go"], "contracts/text/go")
	}
	rules["transform/syntax/shared/generation/go"] = []string{"contracts/text/go", "contracts/generation/go", "contracts/semantic/go", "contracts/transform/go", "trees/ast/go", "trees/representation/go"}
	rules["transform/syntax/horos/source/go"] = []string{"analyze/files/horos/provider/go", "contracts/analysis/go", "contracts/generation/go", "contracts/semantic/go", "contracts/text/go", "contracts/transform/go", "resolve/files/horos/provider/go", "semantics/go", "syntax/horos/go", "trees/ast/go", "trees/representation/go"}
	rules["transform/syntax/shared/target/go"] = []string{"contracts/generation/go", "trees/ast/go", "trees/representation/go"}
	rules["transform/syntax/horos/mapping/go"] = []string{"contracts/generation/go", "transform/syntax/horos/source/go", "syntax/horos/go", "trees/representation/go"}
	for _, syntax := range []string{"typescript", "golang", "rust", "cpp"} {
		rules["transform/syntax/"+syntax+"/target/go"] = []string{"contracts/generation/go", "transform/syntax/shared/target/go", "syntax/" + syntax + "/go", "trees/ast/go"}
		// These are supplied Horos integrations. The shared composer, source
		// reader, target backends and thin transformer factories cannot import them.
		rules["transform/syntax/horos-to-"+syntax+"/mapping/go"] = []string{"contracts/generation/go", "syntax/horos/go", "syntax/" + syntax + "/go", "transform/syntax/horos/source/go", "transform/syntax/horos/mapping/go", "transform/syntax/" + syntax + "/model/go", "trees/representation/go"}
	}
	// The Files dimension has explicit stage paths and a separate sealed carrier.
	rules["files/go"] = []string{"trees/ast/go"}
	rules["source/go"] = []string{"contracts/semantic/go", "files/go", "semantics/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"}
	rules["trees/representation/go"] = append(rules["trees/representation/go"], "semantics/go")
	rules["contracts/semantic/go"] = append(rules["contracts/semantic/go"], "files/go")
	rules["contracts/parse/go"] = append(rules["contracts/parse/go"], "files/go", "source/go", "contracts/semantic/go")
	rules["contracts/format/go"] = append(rules["contracts/format/go"], "source/go", "semantics/go", "contracts/semantic/go")
	rules["contracts/print/go"] = append(rules["contracts/print/go"], "files/go", "source/go")
	rules["engine/go"] = append(rules["engine/go"], "contracts/generation/go", "transform/syntax/shared/generation/go", "files/go", "source/go")
	rules["print/files/shared/provider/go"] = []string{"contracts/language/go", "contracts/print/go", "files/go", "source/go", "trees/cst/go", "trees/representation/go"}
	composition := []string{"contracts/analysis/go", "contracts/check/go", "contracts/parse/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go", "source/go", "semantics/go", "syntax/horos/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"}
	// Files composition uses the public Horos prelude identity from the resolver.
	rules["internal/horosfiles"] = append(append([]string{}, composition...), "contracts/format/go", "files/go", "internal/tree", "resolve/files/horos/provider/go")
	rules["parse/files/horos/provider/go"] = append(append([]string{}, composition...), "internal/horosfiles")
	rules["format/files/horos/provider/go"] = append(append([]string{}, composition...), "contracts/format/go", "internal/horosfiles")
	// Verification projections are explicit concrete compositions, not registered
	// language stages. The agnostic association above never imports their models.
	rules["verification/files/source/horos/go"] = []string{"contracts/semantic/go", "contracts/resolution/go", "internal/horosmodel", "resolve/files/horos/provider/go", "syntax/horos/go", "trees/ast/go", "trees/representation/go"}
	// semantics/<language>/go is classified by its role. The golang entry states
	// the same list and stays until #204, which edits the next line, has merged.
	rules["semantics/golang/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "semantics/go", "syntax/golang/go", "trees/ast/go", "trees/representation/go"}
	rules["internal/golangsem"] = []string{"contracts/semantic/go", "semantics/go", "semantics/golang/go", "syntax/golang/go", "trees/ast/go"}
	rules["resolve/files/golang/provider/go"] = []string{"contracts/resolution/go", "contracts/semantic/go", "internal/golangsem", "semantics/go", "semantics/golang/go", "syntax/golang/go"}
	rules["analyze/files/golang/provider/go"] = []string{"contracts/analysis/go", "contracts/resolution/go", "contracts/semantic/go", "internal/golangsem", "resolve/files/golang/provider/go", "semantics/go", "semantics/golang/go", "syntax/golang/go", "trees/representation/go"}
	rules["internal/typescriptimports"] = []string{"files/go", "internal/typescript", "syntax/typescript/go", "trees/ast/go"}
	for _, syntax := range []string{"typescript", "golang", "rust", "cpp"} {
		rules["transform/syntax/"+syntax+"/model/go"] = []string{"contracts/generation/go", "syntax/" + syntax + "/go"}
		rules["transform/syntax/"+syntax+"/target/go"] = append(rules["transform/syntax/"+syntax+"/target/go"], "transform/syntax/"+syntax+"/model/go")
	}

	rules["transform/syntax/cpp/target/go"] = append(rules["transform/syntax/cpp/target/go"], "internal/cpp")
	rules["transform/syntax/typescript/target/go"] = append(rules["transform/syntax/typescript/target/go"], "internal/typescript", "trees/representation/go", "transform/syntax/typescript/symbol/go")
	rules["format/files/typescript/provider/go"] = []string{"contracts/analysis/go", "contracts/check/go", "contracts/format/go", "contracts/print/go", "contracts/resolution/go", "contracts/semantic/go", "files/go", "internal/typescriptimports", "semantics/go", "semantics/typescript/go", "source/go", "syntax/typescript/go", "trees/ast/go", "trees/cst/go", "trees/representation/go"}
	rules["contracts/generation/go"] = append(rules["contracts/generation/go"], "contracts/check/go", "semantics/go")
	rules["contracts/format/go"] = append(rules["contracts/format/go"], "contracts/generation/go")
	rules["transform/syntax/shared/generation/go"] = append(rules["transform/syntax/shared/generation/go"], "contracts/check/go")
	rules["format/files/typescript/provider/go"] = append(rules["format/files/typescript/provider/go"], "contracts/generation/go", "transform/syntax/typescript/symbol/go")
	rules["transform/syntax/typescript/symbol/go"] = []string{"contracts/generation/go", "internal/typescript", "semantics/go", "syntax/typescript/go", "trees/ast/go", "trees/representation/go"}
	rules["internal/typescriptsem"] = []string{"contracts/parse/go", "contracts/semantic/go", "files/go", "parse/syntax/typescript/provider/go", "semantics/go", "semantics/typescript/go", "syntax/typescript/go", "trees/ast/go"}
	rules["resolve/files/typescript/provider/go"] = []string{"trees/representation/go", "contracts/resolution/go", "contracts/semantic/go", "internal/typescriptsem", "semantics/go", "semantics/typescript/go", "syntax/typescript/go"}
	rules["analyze/files/typescript/provider/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "semantics/go", "semantics/typescript/go", "syntax/typescript/go"}
	rules["internal/rustsem"] = []string{"contracts/parse/go", "contracts/semantic/go", "files/go", "parse/syntax/rust/provider/go", "semantics/go", "semantics/rust/go", "syntax/rust/go", "trees/ast/go"}
	rules["resolve/files/rust/provider/go"] = []string{"trees/representation/go", "contracts/resolution/go", "contracts/semantic/go", "internal/rustsem", "semantics/go", "semantics/rust/go", "syntax/rust/go"}
	rules["analyze/files/rust/provider/go"] = []string{"contracts/analysis/go", "contracts/semantic/go", "semantics/go", "semantics/rust/go", "syntax/rust/go"}
	rules["verification/files/horos-to-go/go"] = []string{"contracts/semantic/go", "contracts/verification/go", "semantics/go", "trees/ast/go", "verification/files/source/horos/go", "semantics/golang/go", "analyze/files/golang/provider/go", "verification/logos/go"}
	rules["verification/logos/go"] = []string{"contracts/verification/go", "semantics/go"}
	rules["verification/files/horos-to-go/mappings/go"] = []string{"verification/files/horos-to-go/go", "semantics/golang/go", "semantics/go"}
	spec := architectureSpec{Module: modulePath, Rules: rules, Exceptions: providerExceptions, ExternalModels: modelDependencies, ExternalImports: modelExternalImports}
	for _, failure := range checkArchitecture(spec, architecturePackages(t, packages)) {
		t.Error(failure)
	}
}

// architecturePackages adds to the `go list` data the symbols that each
// production package of this module selects on its module imports.
func architecturePackages(t *testing.T, packages map[string]productionPackage) map[string]architecturePackage {
	t.Helper()
	name := func(path string) (string, bool) {
		p, ok := packages[path]
		return p.Name, ok && strings.HasPrefix(path, modulePath+"/")
	}
	result := map[string]architecturePackage{}
	for path, p := range packages {
		files := append(slices.Clone(p.GoFiles), p.CgoFiles...)
		data := architecturePackage{Standard: p.Standard, Production: len(files) > 0, Imports: p.Imports}
		if !p.Standard && strings.HasPrefix(path, modulePath+"/") {
			for i, file := range files {
				files[i] = filepath.Join(p.Dir, file)
			}
			var err error
			if data.Uses, data.DotImports, err = parseSymbolUse(files, name); err != nil {
				t.Fatal(err)
			}
		}
		result[path] = data
	}
	return result
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
	"verification/logos/go":                      {"github.com/Bitspark/logos-contract/contract/go/typed"},
	"verification/files/horos-to-go/mappings/go": {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
	"tools/horos-module/go":                      {"github.com/bitspark/ontos/codec/go"},
	"transform/syntax/horos/source/go":           {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"parse/syntax/horos/provider/go":             {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"format/syntax/horos/provider/go":            {"github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"internal/horosmodel":                        {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/bitspark/ontos/codec/go"},
	"resolve/files/horos/provider/go":            {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go"},
	"verification/files/source/horos/go":         {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go"},
	"verification/files/horos-to-go/go":          {"github.com/Bitspark/horos/core/go/declarations", "github.com/bitspark/ontos/core/go", "github.com/bitspark/ontos/data/go", "github.com/Bitspark/logos-contract/contract/go/typed"},
}

func TestFilesDimensionStageCoverage(t *testing.T) {
	fixtures := filesStageFixtures()
	for _, stage := range []string{"parse", "format"} {
		entries, err := os.ReadDir(filepath.Join("../..", stage, "files"))
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			pkg := stagePath(stage+"/files", entry.Name())
			if fixtures[pkg] == nil {
				t.Errorf("Files stage %s lacks a conformance fixture", pkg)
			}
			delete(fixtures, pkg)
		}
	}
	for pkg := range fixtures {
		t.Errorf("Files fixture names absent stage %s", pkg)
	}
}
