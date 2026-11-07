package integration_test

import (
	"context"
	"errors"
	goast "go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	generationfixture "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"

	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func generationSPIFixture(t *testing.T, useFirst bool) (s.AST, string) {
	t.Helper()
	ctx := context.Background()
	f, err := generationfixture.New(ctx)
	if err != nil {
		t.Fatal(err)
	}
	session, err := g.Begin[generationfixture.SourceLanguage, s.Language](ctx, f.Input(), f.Config(), g.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if useFirst {
		for _, name := range []string{"NestedPage", "Inline", "External"} {
			product, err := session.Use(f.Use(name))
			if err != nil || !product.Valid() {
				t.Fatalf("use %s: %v", name, err)
			}
		}
	}
	for _, root := range f.Roots() {
		publication, err := session.Publish(root)
		if err != nil || publication.Root() != root {
			t.Fatalf("publish: %v", err)
		}
	}
	for _, root := range f.Roots() {
		callee, err := session.Acquire(root)
		if err != nil {
			t.Fatal(err)
		}
		again, err := session.Acquire(root)
		if err != nil || callee.Signature() != again.Signature() || callee.Layout() != again.Layout() || callee.Frame() != again.Frame() {
			t.Fatalf("callee was reintroduced: %v", err)
		}
		if _, err = session.Define(callee); err != nil {
			t.Fatalf("define %s: %v", callee.Key(), err)
		}
	}
	helper, err := session.AcquireHelper("cursor-helper")
	if err != nil {
		t.Fatal(err)
	}
	if helper.Open().Declaration.Valid() || helper.Open().Environment.Valid() || len(helper.Open().Formals) != 0 {
		t.Fatal("native helper inherited source declarations or binders")
	}
	if _, err = session.Define(helper); err != nil {
		t.Fatal(err)
	}
	output, err := session.Finish()
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[s.Language](f.Target(), output.EntryAST())
	if err != nil {
		t.Fatal(err)
	}
	return tree, formatted(t, golangText(), tree)
}

func nodeAttribute[L any](node ast.Node[L], name string) string {
	for _, attribute := range node.Attributes {
		if attribute.Name == name {
			return attribute.Text
		}
	}
	return ""
}

func TestGenerationSPIIndependentGoComposition(t *testing.T) {
	tree, source := generationSPIFixture(t, true)
	for _, want := range []string{
		"type Person struct", "Name string", "type Cursor string",
		"type Page[T interface {}, ST interface {}] struct", "Items []T", "Support ST", "Cursor Cursor",
		"type Outer[U interface {}, SU interface {}] = Page[U, SU]",
		"type Nested[V interface {}, SV interface {}] = Page[Page[V, SV], []V]",
		"type People = Page[Person, string]", "type InlinePeople = []Person",
		"type ExternalPeople = ExternalBox[Person, string]",
	} {
		if !strings.Contains(source, want) {
			t.Fatalf("missing %q:\n%s", want, source)
		}
	}
	if strings.Contains(source, "type Inline[") || strings.Contains(source, "type ExternalBox[") {
		t.Fatal("inline/external mapping fabricated a local definition")
	}
	declarations := make(map[string]ast.NodeID)
	parameters := make(map[string]map[string]ast.NodeID)
	for index, node := range tree.Nodes {
		if node.Kind != s.KindTypeDeclaration && node.Kind != s.KindTypeAlias {
			continue
		}
		name := nodeAttribute(node, s.AttributeName)
		declarations[name] = ast.NodeID(index + 1)
		parameters[name] = make(map[string]ast.NodeID)
		for _, child := range node.Children {
			if child.Role == s.RoleTypeParameter {
				parameters[name][nodeAttribute(tree.Nodes[child.Node-1], s.AttributeName)] = child.Node
			}
		}
	}
	if len(parameters["Page"]) != 2 || len(parameters["Person"]) != 0 || len(parameters["Cursor"]) != 0 {
		t.Fatal("complete native header was not retained")
	}
	if parameters["Page"]["T"] == parameters["Outer"]["U"] {
		t.Fatal("callee and caller formal identity collapsed")
	}
	parents := make(map[ast.NodeID]int)
	externalUses := 0
	for _, node := range tree.Nodes {
		for _, child := range node.Children {
			parents[child.Node]++
		}
		name := nodeAttribute(node, s.AttributeName)
		if node.Kind == s.KindNamedType && name == "ExternalBox" {
			externalUses++
			if len(node.References) != 0 {
				t.Fatal("external export received a fabricated NodeRef")
			}
		}
		if node.Kind == s.KindNamedType && (name == "U" || name == "SU") {
			if len(node.References) != 1 || node.References[0].Target.Node != parameters["Outer"][name] {
				t.Fatal("Outer use rebound to the callee's formal")
			}
		}
	}
	for node, count := range parents {
		if count != 1 {
			t.Fatalf("fragment occurrence %d has %d containment parents", node, count)
		}
	}
	if externalUses != 1 || len(declarations) != 8 {
		t.Fatalf("unexpected definition surface: %d external uses, %v", externalUses, declarations)
	}
	_, acquiredFirst := generationSPIFixture(t, false)
	if source != acquiredFirst {
		t.Fatalf("first closed use changed the open header:\n%s\n%s", source, acquiredFirst)
	}
	checkGenerationSPIGo(t, source)
}

func TestGenerationSPIRejectsArgumentOnlyResultSupport(t *testing.T) {
	ctx := context.Background()
	f, err := generationfixture.NewArgumentOnlySupport(ctx)
	if err != nil {
		t.Fatal(err)
	}
	session, err := g.Begin[generationfixture.SourceLanguage, s.Language](ctx, f.Input(), f.Config(), g.DefaultLimits())
	if err != nil {
		t.Fatal(err)
	}
	if product, err := session.Use(f.Use("NestedPage")); !errors.Is(err, g.ErrScope) || product.Valid() {
		t.Fatalf("argument-only support was accepted as Page result support: %v", err)
	}
	if output, err := session.Finish(); !errors.Is(err, g.ErrScope) || output.LanguageID() != "" {
		t.Fatalf("failed relation produced usable output: %v", err)
	}
}

func TestGenerationSPINativeFragmentCapture(t *testing.T) {
	if err := generationfixture.CheckCaptureIsolation(); err != nil {
		t.Fatal(err)
	}
}

const generationSPIConsumer = `package fixture
var _ People = Page[Person, string]{Items: []Person{{Name: "Ada"}}, Support: "person"}
var _ Outer[Person, string] = People{}
var _ Nested[Person, string] = Page[Page[Person, string], []Person]{Items: []Page[Person, string]{{}}}
var _ InlinePeople = []Person{{Name: "Ada"}}
var _ ExternalPeople = ExternalBox[Person, string]{Value: Person{Name: "Ada"}}
`

func checkGenerationSPIGo(t *testing.T, source string) {
	t.Helper()
	fset := token.NewFileSet()
	var files []*goast.File
	for _, input := range []struct{ name, source string }{{"generated.go", source}, {"external.go", generationfixture.ExternalSource}, {"consumer.go", generationSPIConsumer}} {
		file, err := parser.ParseFile(fset, input.name, input.source, parser.AllErrors)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, file)
	}
	if _, err := (&types.Config{GoVersion: "go1.27"}).Check("fixture", fset, files, nil); err != nil {
		t.Fatalf("generated public-SPI output did not typecheck: %v\n%s", err, source)
	}
}

func TestGenerationSPIHostedGoCompiler(t *testing.T) {
	if os.Getenv("BITTYPE_NATIVE") != "1" {
		t.Skip("native compiler qualification runs in the maintained hosted job")
	}
	_, source := generationSPIFixture(t, true)
	dir := t.TempDir()
	for name, content := range map[string]string{"go.mod": "module fixture\n\ngo 1.27.0\n", "generated.go": source, "external.go": generationfixture.ExternalSource, "consumer.go": generationSPIConsumer} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	command(t, dir, "go", "test", "./...")
}
