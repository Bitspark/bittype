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
	Standard bool
	Imports []string
	GoFiles []string
}

func productionPackages(t *testing.T) map[string]productionPackage {
	t.Helper()
	command := exec.Command("go", "list", "-deps", "-json", "./...")
	command.Dir = "../.."
	output, err := command.Output()
	if err != nil { t.Fatal(err) }
	result := map[string]productionPackage{}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var p productionPackage
		if err := decoder.Decode(&p); err == io.EOF { break } else if err != nil { t.Fatal(err) }
		result[p.ImportPath] = p
	}
	return result
}

// Discover the stage matrix from source directories, not a production language list.
// The union catches an orphan parser as well as an incomplete syntax implementation.
func stageLanguages(t *testing.T) []string {
	t.Helper()
	names := map[string]bool{}
	for _, stage := range nativeStages {
		entries, err := os.ReadDir(filepath.Join("../..", stage))
		if err != nil { t.Fatal(err) }
		for _, entry := range entries {
			if !entry.IsDir() { continue }
			if stage == "check" && (entry.Name() == "ast" || entry.Name() == "cst") { continue }
			names[entry.Name()] = true
		}
	}
	var result []string
	for name := range names { result = append(result, name) }
	slices.Sort(result)
	return result
}

func TestLanguageStageMatrix(t *testing.T) {
	packages := productionPackages(t)
	languages := stageLanguages(t)
	if len(languages) == 0 { t.Fatal("no language implementations") }
	qualified := map[string]bool{}
	for _, fixture := range nativeConformanceCases() {
		if qualified[fixture.directory] { t.Fatalf("duplicate conformance fixture: %s", fixture.directory) }
		qualified[fixture.directory] = true
	}
	for _, name := range languages {
		for _, stage := range nativeStages {
			path := modulePath + "/" + stage + "/" + name + "/go"
			if p, ok := packages[path]; !ok || len(p.GoFiles) == 0 { t.Errorf("missing implemented stage %s", path) }
		}
		if !qualified[name] { t.Errorf("%s is missing shared provider conformance/round-trip/fuzz coverage", name) }
		delete(qualified, name)
	}
	for name := range qualified { t.Errorf("conformance fixture has no language implementation: %s", name) }
}

// Each production package is checked, including private helpers. This prevents
// a forbidden dependency from being hidden behind an otherwise allowed helper.
func TestProductionArchitecture(t *testing.T) {
	packages := productionPackages(t)
	languages := stageLanguages(t)
	rules := map[string][]string{
		"trees/ast/go": {},
		"trees/cst/go": {"trees/ast/go"},
		"contracts/language/go": {},
		"contracts/pipeline/go": {},
		"contracts/text/go": {"trees/ast/go", "trees/cst/go", "contracts/language/go"},
		"check/ast/go": {"trees/ast/go"},
		"check/cst/go": {"trees/ast/go", "trees/cst/go", "check/ast/go"},
		"print/go": {"trees/cst/go", "check/cst/go"},
		"trees/representation/go": {"trees/ast/go", "trees/cst/go", "contracts/language/go", "contracts/text/go", "check/ast/go", "check/cst/go", "internal/tree"},
		"contracts/transform/go": {"trees/ast/go", "trees/representation/go", "contracts/language/go", "contracts/pipeline/go"},
		"engine/go": {"trees/representation/go", "contracts/language/go", "contracts/text/go", "contracts/transform/go"},
		"internal/tree": {"trees/ast/go"},
		"internal/grammar": {"trees/ast/go", "check/ast/go", "contracts/text/go"},
		"internal/lex": {"trees/ast/go", "contracts/text/go"},
		"internal/parseutil": {"trees/ast/go", "check/ast/go", "contracts/text/go", "internal/lex"},
		"internal/concrete": {"trees/ast/go", "trees/cst/go", "check/cst/go", "contracts/text/go", "internal/lex", "internal/tree"},
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
			if stage != "check" { allowed = append(allowed, checker, "internal/concrete") }
			if stage == "parse" { allowed = append(allowed, "internal/lex", "internal/parseutil") }
			rules[stage+"/"+name+"/go"] = allowed
		}
		rules["provider/"+name+"/go"] = append(append([]string{}, base...), checker, "parse/"+name+"/go", "format/"+name+"/go", "print/go")
		// Language-specific lexical helpers cannot depend on stages or another language.
		rules[helper] = []string{"trees/ast/go", "contracts/text/go", "internal/grammar", "internal/lex", "internal/tree"}
	}
	for path, p := range packages {
		if p.Standard || path == modulePath+"/integration/go" { continue }
		if !strings.HasPrefix(path, modulePath+"/") { t.Errorf("external production dependency: %s", path); continue }
		relative := strings.TrimPrefix(path, modulePath+"/")
		allowed, ok := rules[relative]
		if !ok { t.Errorf("unclassified production package: %s", relative); continue }
		for _, dependency := range p.Imports {
			if dep, ok := packages[dependency]; ok && dep.Standard { continue }
			if !strings.HasPrefix(dependency, modulePath+"/") || !slices.Contains(allowed, strings.TrimPrefix(dependency, modulePath+"/")) {
				t.Errorf("%s must not import %s", relative, dependency)
			}
		}
	}
}
