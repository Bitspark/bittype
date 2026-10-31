package integration_test

import (
	"fmt"
	goast "go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const modulePath = "github.com/Bitspark/bittype"

type contractImporter map[string]*types.Package

func (packages contractImporter) Import(path string) (*types.Package, error) {
	if result, ok := packages[path]; ok {
		return result, nil
	}
	if !strings.HasPrefix(path, modulePath+"/") { return importer.Default().Import(path) }
	return nil, fmt.Errorf("unexpected dependency in public contract fixture: %s", path)
}

// Check separately authored consumers against the actual public declarations.
// No language grammar or transformation implementation is needed by a consumer.
func loadContracts(t *testing.T) contractImporter {
	t.Helper()
	packages := contractImporter{}
	for _, name := range []string{
		"ast", "cst", "pipeline", "language", "text", "syntax/typescript", "syntax/golang", "syntax/rust", "syntax/cpp",
	} {
		dir := filepath.Join("..", "..", name, "go")
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		files := []*goast.File{}
		positions := token.NewFileSet()
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
				continue
			}
			file, err := parser.ParseFile(positions, filepath.Join(dir, entry.Name()), nil, parser.AllErrors)
			if err != nil {
				t.Fatal(err)
			}
			files = append(files, file)
		}
		path := modulePath + "/" + name + "/go"
		config := types.Config{GoVersion: "go1.27", Importer: packages}
		pkg, err := config.Check(path, positions, files, nil)
		if err != nil {
			t.Fatal(err)
		}
		packages[path] = pkg
	}
	return packages
}

func checkConsumer(t *testing.T, packages contractImporter, name string) error {
	t.Helper()
	if diagnostics := consumerDiagnostics(t, packages, name); len(diagnostics) != 0 {
		return diagnostics[0]
	}
	return nil
}

func consumerDiagnostics(t *testing.T, packages contractImporter, name string) []error {
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
