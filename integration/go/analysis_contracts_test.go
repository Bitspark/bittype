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

func TestAnalysisIndependentPublicContracts(t *testing.T) { testSemanticPublicContracts(t, "analysis") }
func TestResolutionIndependentPublicContracts(t *testing.T) {
	testSemanticPublicContracts(t, "resolution")
}

func testSemanticPublicContracts(t *testing.T, stage string) {
	t.Helper()
	imports := &engineImporter{standard: importer.Default(), packages: make(map[string]*types.Package)}
	cases := []struct{ name, want string }{
		{"consumer", ""},
		{"wrong_marker", "cannot use"},
		{"wrong_input", "cannot use"},
		{"missing_identity", "ID"},
		{"private_result", "unexported"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			positions := token.NewFileSet()
			file, err := parser.ParseFile(positions, filepath.Join("testdata", stage, tc.name+".go"), nil, parser.AllErrors)
			if err != nil {
				t.Fatal(err)
			}
			config := types.Config{GoVersion: "go1.27", Importer: imports}
			_, err = config.Check("example.invalid/"+stage+"-consumer/"+tc.name, positions, []*goast.File{file}, nil)
			if tc.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("want %q: %v", tc.want, err)
			}
		})
	}
}
