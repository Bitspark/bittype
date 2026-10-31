package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"testing"
)

// Inspect production dependency closure, excluding integration-test imports.
// Neither foundation is allowed to acquire the other or a language/runtime
// dependency through an indirect import.
func TestFoundationDependencies(t *testing.T) {
	command := exec.Command("go", "list", "-deps", "-json", "./ast/go", "./pipeline/go")
	command.Dir = "../.."
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]bool{
		modulePath + "/ast/go":      false,
		modulePath + "/pipeline/go": false,
	}
	decoder := json.NewDecoder(bytes.NewReader(output))
	for {
		var pkg struct {
			ImportPath string
			Standard   bool
			Imports    []string
		}
		if err := decoder.Decode(&pkg); err == io.EOF {
			break
		} else if err != nil {
			t.Fatal(err)
		}
		if pkg.Standard {
			continue
		}
		if _, ok := expected[pkg.ImportPath]; !ok {
			t.Fatalf("unexpected production dependency: %s", pkg.ImportPath)
		}
		expected[pkg.ImportPath] = true
		for _, dependency := range pkg.Imports {
			if _, foundation := expected[dependency]; foundation {
				t.Fatalf("%s must not import %s", pkg.ImportPath, dependency)
			}
		}
	}
	for name, seen := range expected {
		if !seen {
			t.Errorf("production package missing from dependency evidence: %s", name)
		}
	}
}
