package integration_test

import (
	goast "go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/Bitspark/bittype"

// Check separately authored consumers against actual public declarations.
func loadContracts(t *testing.T) *engineImporter {
	t.Helper()
	return &engineImporter{standard: importer.Default(), packages: make(map[string]*types.Package)}
}

func checkConsumer(t *testing.T, packages types.Importer, name string) error {
	t.Helper()
	if diagnostics := consumerDiagnostics(t, packages, name); len(diagnostics) != 0 {
		return diagnostics[0]
	}
	return nil
}

func consumerDiagnostics(t *testing.T, packages types.Importer, name string) []error {
	t.Helper()
	positions := token.NewFileSet()
	file, err := parser.ParseFile(positions, filepath.Join("testdata", name+".go"), nil, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	var diagnostics []error
	config := types.Config{
		GoVersion: "go1.27",
		Importer:  packages,
		Error:     func(err error) { diagnostics = append(diagnostics, err) },
	}
	_, err = config.Check("example.invalid/consumer", positions, []*goast.File{file}, nil)
	if err != nil && len(diagnostics) == 0 {
		diagnostics = append(diagnostics, err)
	}
	return diagnostics
}

func TestIndependentLanguageConsumer(t *testing.T) {
	if err := checkConsumer(t, loadContracts(t), "consumer"); err != nil {
		t.Fatal(err)
	}
}

func TestLanguageAndStageSeparation(t *testing.T) {
	packages := loadContracts(t)
	cases := []struct {
		name string
		want []string
	}{
		{"mixed_trees", []string{"cannot use", "AST[right]", "AST[left]"}},
		{"mixed_nodes", []string{"cannot use", "Node[right]", "Node[left]"}},
		{"mixed_kinds", []string{"cannot use", "Kind[right]", "Kind[left]"}},
		{"stage_input", []string{"cannot use", "AST[left]", "AST[right]"}},
		{"stage_output", []string{"cannot use", "AST[right]", "string"}},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			err := checkConsumer(t, packages, test.name)
			if err == nil {
				t.Fatal("incompatible public contract unexpectedly type-checked")
			}
			for _, part := range test.want {
				if !strings.Contains(err.Error(), part) {
					t.Fatalf("expected diagnostic containing %q, got: %v", part, err)
				}
			}
		})
	}
}
