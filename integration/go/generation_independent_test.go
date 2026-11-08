package integration_test

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	s "github.com/Bitspark/bittype/syntax/golang/go"
)

func TestDeclarativeIndependentSourceAndExtraParameter(t *testing.T) {
	fixture, err := f.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := fixture.Plan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	tree := outputTree[s.Language](t, out)
	text := formatted(t, golangText(), tree)
	for _, want := range []string{"type Page[T0 interface {}, Storage interface {}] struct", "type Outer[T0 interface {}] = Page[T0, []T0]", "type Nested[T0 interface {}] = Page[Page[T0, []T0], []Page[T0, []T0]]", "type People = Page[Person, []Person]", "Cursor Cursor"} {
		if !strings.Contains(text, want) {
			t.Fatalf("missing %q in\n%s", want, text)
		}
	}
	if strings.Count(text, "type Cursor ") != 1 {
		t.Fatal("helper duplicated", text)
	}
	fs := token.NewFileSet()
	file, err := parser.ParseFile(fs, "types.go", text, parser.AllErrors)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = (&types.Config{}).Check("generated", fs, []*ast.File{file}, nil); err != nil {
		t.Fatal(err)
	}
}

func TestDeclarativeIndependentHostedGoCompiler(t *testing.T) {
	if os.Getenv("BITTYPE_NATIVE") != "1" {
		t.Skip("hosted Go compiler qualification")
	}
	fixture, err := f.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	plan, err := fixture.Plan(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	text := formatted(t, golangText(), outputTree[s.Language](t, out))
	dir := t.TempDir()
	if err = os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module independent

go 1.27.0
`), 0600); err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(dir, "types.go"), []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
	command(t, dir, "go", "test", "./...")
}
