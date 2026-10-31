package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"testing"
)

// Inspect production dependency closure, excluding integration-test imports.
// Syntax packages may depend on the two agnostic contracts, never each other.
// The AST, pipeline and language identity foundations remain independent.
func TestFoundationDependencies(t *testing.T) {
	command := exec.Command("go", "list", "-deps", "-json",
		"./ast/go", "./pipeline/go", "./language/go", "./syntax/...")
	command.Dir = "../.."
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	allowed := map[string]map[string]bool{
		modulePath + "/ast/go":      {},
		modulePath + "/pipeline/go": {},
		modulePath + "/language/go": {},
	}
	for _, name := range []string{"typescript", "golang", "rust", "cpp"} {
		allowed[modulePath+"/syntax/"+name+"/go"] = map[string]bool{
			modulePath + "/ast/go":      true,
			modulePath + "/language/go": true,
		}
	}
	type packageInfo struct {
		ImportPath string
		Standard   bool
		Imports    []string
	}
	var packages []packageInfo
	standard := map[string]bool{}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var pkg packageInfo
		if err := decoder.Decode(&pkg); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		packages = append(packages, pkg)
		standard[pkg.ImportPath] = pkg.Standard
	}
	seen := map[string]bool{}
	for _, pkg := range packages {
		if pkg.Standard {
			continue
		}
		dependencies, ok := allowed[pkg.ImportPath]
		if !ok {
			t.Fatalf("unexpected production dependency: %s", pkg.ImportPath)
		}
		seen[pkg.ImportPath] = true
		for _, dependency := range pkg.Imports {
			if !standard[dependency] && !dependencies[dependency] {
				t.Errorf("%s must not import %s", pkg.ImportPath, dependency)
			}
		}
	}
	for name := range allowed {
		if !seen[name] {
			t.Errorf("production package missing from dependency evidence: %s", name)
		}
	}
}
