package integration_test

import (
	"github.com/Bitspark/bittype/trees/ast/go"
	cppprovider "github.com/Bitspark/bittype/provider/cpp/go"
	goprovider "github.com/Bitspark/bittype/provider/golang/go"
	rustprovider "github.com/Bitspark/bittype/provider/rust/go"
	tsprovider "github.com/Bitspark/bittype/provider/typescript/go"
	cpp "github.com/Bitspark/bittype/syntax/cpp/go"
	golang "github.com/Bitspark/bittype/syntax/golang/go"
	rust "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// These ASTs are authored independently of every parser. The small helper only
// allocates addresses; there is no type mapping or source-to-AST fixture producer.
type fixture[L any] struct{ tree ast.AST[L] }

func (f *fixture[L]) node(kind string, attrs []ast.Attribute, children ...ast.Child) ast.NodeID {
	f.tree.Nodes = append(f.tree.Nodes, ast.Node[L]{Kind: ast.Kind[L](kind), Attributes: attrs, Children: children})
	return ast.NodeID(len(f.tree.Nodes))
}
func attr(pairs ...string) []ast.Attribute {
	var a []ast.Attribute
	for i := 0; i < len(pairs); i += 2 {
		a = append(a, ast.Attribute{Name: pairs[i], Text: pairs[i+1]})
	}
	return a
}
func child(role string, id ast.NodeID) ast.Child { return ast.Child{Role: role, Node: id} }
func tsFixture() ast.AST[ts.Language] {
	f := fixture[ts.Language]{}
	str := f.node("keyword-type", attr("name", "string"))
	null := f.node("literal-type", attr("literal-kind", "null", "text", "null"))
	union := f.node("union-type", nil, child("type", str), child("type", null))
	array := f.node("array-type", nil, child("element", union))
	alias := f.node("type-alias", attr("name", "Strings", "exported", "true"), child("type", array))
	f.tree.Root = f.node("source-file", nil, child("declaration", alias))
	f.tree.ID = "independent-ts"
	return f.tree
}
func goFixture() ast.AST[golang.Language] {
	f := fixture[golang.Language]{}
	str := f.node("named-type", attr("name", "string"))
	slice := f.node("slice-type", nil, child("element", str))
	alias := f.node("type-alias", attr("name", "Names"), child("type", slice))
	f.tree.Root = f.node("file", attr("package", "fixture"), child("declaration", alias))
	f.tree.ID = "independent-go"
	return f.tree
}
func rustFixture() ast.AST[rust.Language] {
	f := fixture[rust.Language]{}
	str := f.node("primitive-type", attr("name", "str"))
	lifetime := f.node("lifetime-reference", attr("name", "static"))
	ref := f.node("reference-type", nil, child("element", str), child("lifetime", lifetime))
	alias := f.node("type-alias", attr("name", "Name", "visibility", "public"), child("type", ref))
	f.tree.Root = f.node("file", nil, child("item", alias))
	f.tree.ID = "independent-rust"
	return f.tree
}
func cppFixture() ast.AST[cpp.Language] {
	f := fixture[cpp.Language]{}
	integer := f.node("builtin-type", attr("name", "int"))
	array := f.node("array-type", attr("length", "3"), child("element", integer))
	pointer := f.node("pointer-type", nil, child("element", array))
	alias := f.node("type-alias", attr("name", "MatrixPtr"), child("type", pointer))
	visibleType := f.node("builtin-type", attr("name", "int"))
	visible := f.node("field", attr("name", "visible", "access", "public"), child("type", visibleType))
	hiddenType := f.node("builtin-type", attr("name", "int"))
	hidden := f.node("field", attr("name", "hidden"), child("type", hiddenType))
	record := f.node("record", attr("name", "AccessFixture", "record-key", "class"), child("member", visible), child("member", hidden))
	f.tree.Root = f.node("translation-unit", nil, child("declaration", alias), child("declaration", record))
	f.tree.ID = "independent-cpp"
	return f.tree
}
func formatted[L any](t *testing.T, p text.Provider[L], a ast.AST[L]) string {
	t.Helper()
	c, err := p.Format(a)
	if err != nil {
		t.Fatal(err)
	}
	result, err := p.Print(c)
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func TestIndependentFormatterGoldens(t *testing.T) {
	cases := map[string]struct{ got, want string }{
		"typescript": {formatted(t, tsprovider.Provider{}, tsFixture()), "export type Strings = (string | null)[];\n"},
		"go":         {formatted(t, goprovider.Provider{}, goFixture()), "package fixture\n\ntype Names = []string\n"},
		"rust":       {formatted(t, rustprovider.Provider{}, rustFixture()), "pub type Name = &'static str;\n"},
		"cpp":        {formatted(t, cppprovider.Provider{}, cppFixture()), "using MatrixPtr = int (*)[3];\n\nclass AccessFixture {\n  public:\n  int visible;\n  private:\n  int hidden;\n};\n"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Fatalf("want %q; got %q", tc.want, tc.got)
			}
		})
	}
}
func TestFormatterSnapshotOwnership(t *testing.T) {
	a := tsFixture()
	a.Nodes[0].Origin = ast.Origin{Spans: []ast.Span{{Source: "original", Start: 0, End: 1}}, Nodes: []ast.NodeRef{{Tree: "external", Node: 1}}}
	// An ordinary named use can retain a separately resolved target.
	a.Nodes[0].Kind = ts.KindTypeReference
	a.Nodes[0].Attributes = attr("name", "External")
	a.Nodes[0].References = []ast.Reference{{Role: "target", Target: ast.NodeRef{Tree: "external", Node: 2}}}
	p := tsprovider.Provider{}
	c, err := p.Format(a)
	if err != nil {
		t.Fatal(err)
	}
	before, err := p.Print(c)
	if err != nil {
		t.Fatal(err)
	}
	a.Nodes[0].Attributes[0].Text = "Changed"
	a.Nodes[0].References[0].Target.Node = 99
	a.Nodes[0].Origin.Spans[0].Source = "changed"
	a.Nodes[0].Origin.Nodes[0].Node = 99
	a.Nodes[a.Root-1].Children[0].Node = 99
	view := c.ASTView()
	if view.Nodes[0].Attributes[0].Text != "External" || view.Nodes[0].References[0].Target.Node != 2 || view.Nodes[0].Origin.Spans[0].Source != "original" || view.Nodes[0].Origin.Nodes[0].Node != 1 || view.Nodes[view.Root-1].Children[0].Node == 99 {
		t.Fatal("CST aliases original AST storage")
	}
	after, err := p.Print(c)
	if err != nil || after != before {
		t.Fatal("formatter output changed with input mutation")
	}
}
func invalidAST[L any](t *testing.T, p text.Provider[L], makeTree func() ast.AST[L]) {
	t.Helper()
	mutations := []func(*ast.AST[L]){
		func(a *ast.AST[L]) { a.ID = "" },
		func(a *ast.AST[L]) { a.Root = 9999 },
		func(a *ast.AST[L]) {
			a.Nodes[a.Root-1].Children = append(a.Nodes[a.Root-1].Children, child("unknown", 1))
		},
		func(a *ast.AST[L]) {
			a.Nodes[a.Root-1].Children = append(a.Nodes[a.Root-1].Children, child("declaration", a.Root))
		},
		func(a *ast.AST[L]) {
			a.Nodes[a.Root-1].Attributes = append(a.Nodes[a.Root-1].Attributes, attr("unknown", "x")...)
		},
		func(a *ast.AST[L]) { a.Nodes[0].Attributes = append(a.Nodes[0].Attributes, a.Nodes[0].Attributes[0]) },
		func(a *ast.AST[L]) { a.Nodes[0].Kind = "unowned-kind" },
		func(a *ast.AST[L]) { a.Nodes = append(a.Nodes, ast.Node[L]{Kind: "orphan"}) },
		func(a *ast.AST[L]) { a.Nodes[0].Origin.Spans = []ast.Span{{Source: "s", Start: 2, End: 1}} },
	}
	for i, mutate := range mutations {
		a := makeTree()
		mutate(&a)
		if err := p.Check(a); err == nil {
			t.Errorf("mutation %d accepted by checker", i)
		}
		out, err := p.Format(a)
		if err == nil || len(out.Nodes) > 0 {
			t.Errorf("mutation %d produced output", i)
		}
	}
}
func TestMalformedASTs(t *testing.T) {
	invalidAST(t, tsprovider.Provider{}, tsFixture)
	invalidAST(t, goprovider.Provider{}, goFixture)
	invalidAST(t, rustprovider.Provider{}, rustFixture)
	invalidAST(t, cppprovider.Provider{}, cppFixture)
}

func TestPrimitiveSyntaxPositions(t *testing.T) {
	a := tsFixture()
	a.Nodes[0].Kind = ts.KindTypeReference
	if err := (tsprovider.Provider{}).Check(a); err == nil {
		t.Fatal("keyword spelling accepted as unqualified reference")
	}
	f := fixture[rust.Language]{}
	segment := f.node("path-segment", attr("name", "u8"))
	path := f.node("path-type", nil, child("segment", segment))
	alias := f.node("type-alias", attr("name", "X"), child("type", path))
	f.tree.ID = "primitive-position"
	f.tree.Root = f.node("file", nil, child("item", alias))
	if err := (rustprovider.Provider{}).Check(f.tree); err == nil {
		t.Fatal("primitive spelling accepted as an unqualified type path")
	}
	for _, source := range []string{"type X = NS.string;", "interface X { string: string; }"} {
		if _, err := (tsprovider.Provider{}).Parse(text.Input{SourceID: "native-name", TreeID: "native-name", Text: source}); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := (rustprovider.Provider{}).Parse(text.Input{SourceID: "native-path", TreeID: "native-path", Text: "type X = u8::Associated;"}); err != nil {
		t.Fatal(err)
	}
}
func nativeFiles[L any](t *testing.T, p text.Provider[L], language, ext string, independent ast.AST[L]) string {
	t.Helper()
	directory := t.TempDir()
	paths, err := filepath.Glob("../../testdata/text/" + language + "/valid/*." + ext)
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		bytes, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		c, err := p.Parse(text.Input{SourceID: "native-input", TreeID: "native-tree", Text: string(bytes)})
		if err != nil {
			t.Fatal(err)
		}
		output := formatted(t, p, c.Abstract)
		if err := os.WriteFile(filepath.Join(directory, filepath.Base(path)), []byte(output), 0600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(directory, "independent."+ext), []byte(formatted(t, p, independent)), 0600); err != nil {
		t.Fatal(err)
	}
	return directory
}
func command(t *testing.T, dir, name string, args ...string) {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	t.Logf("%s %s\n%s", name, strings.Join(args, " "), output)
	if err != nil {
		t.Fatal(err)
	}
}
func TestNativeCompilerFixtures(t *testing.T) {
	if os.Getenv("BITTYPE_NATIVE") != "1" {
		t.Skip("native compiler qualification runs in the dedicated remote job")
	}
	t.Run("typescript", func(t *testing.T) {
		dir := nativeFiles(t, tsprovider.Provider{}, "typescript", "ts", tsFixture())
		command(t, dir, "tsc", "--version")
		command(t, dir, "tsc", "--noEmit", "--strict", "--target", "ES2020", "--lib", "ES2020", "--module", "ESNext", "--moduleResolution", "node", "types.ts", "external.ts", "independent.ts")
	})
	t.Run("go", func(t *testing.T) {
		dir := nativeFiles(t, goprovider.Provider{}, "golang", "go", goFixture())
		// The independent example has its own package and is compiled separately.
		independent := t.TempDir()
		b, err := os.ReadFile(filepath.Join(dir, "independent.go"))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(independent, "independent.go"), b, 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Remove(filepath.Join(dir, "independent.go")); err != nil {
			t.Fatal(err)
		}
		for _, d := range []string{dir, independent} {
			if err := os.WriteFile(filepath.Join(d, "go.mod"), []byte("module qualification\n\ngo 1.27.0\n"), 0600); err != nil {
				t.Fatal(err)
			}
			command(t, d, "go", "test", "-vet=off", "./...")
		}
	})
	t.Run("rust", func(t *testing.T) {
		dir := nativeFiles(t, rustprovider.Provider{}, "rust", "rs", rustFixture())
		for _, file := range []string{"types.rs", "independent.rs"} {
			command(t, dir, "docker", "run", "--rm", "--network=none", "-v", dir+":/work", "-w", "/work", "rust:1.90.0", "rustc", "--edition=2024", "--crate-type=lib", "-A", "dead_code", "-A", "unused_imports", file)
		}
	})
	t.Run("cpp", func(t *testing.T) {
		dir := nativeFiles(t, cppprovider.Provider{}, "cpp", "hpp", cppFixture())
		probe := "#include <type_traits>\n#include \"types.hpp\"\n#include \"independent.hpp\"\nstatic_assert(std::is_same_v<MatrixPtr, int (*)[3]>);\nstatic_assert(std::is_same_v<decltype(api::Shapes::array), int *[3]>);\nstatic_assert(std::is_same_v<decltype(api::Shapes::pointer), int (*)[3]>);\nstatic_assert(std::is_same_v<decltype(api::Shapes::callback), int (*)(double) noexcept>);\nstatic_assert(std::is_same_v<decltype(api::Shapes::member), int api::Person::*>);\n"
		if err := os.WriteFile(filepath.Join(dir, "probe.cpp"), []byte(probe), 0600); err != nil {
			t.Fatal(err)
		}
		command(t, dir, "docker", "run", "--rm", "--network=none", "-v", dir+":/work", "-w", "/work", "gcc:14.2.0", "g++", "-std=c++20", "-fsyntax-only", "probe.cpp")
		if err := os.WriteFile(filepath.Join(dir, "access.cpp"), []byte("#include \"independent.hpp\"\nusing Hidden = decltype(AccessFixture::hidden);\n"), 0600); err != nil {
			t.Fatal(err)
		}
		cmd := exec.Command("docker", "run", "--rm", "--network=none", "-v", dir+":/work", "-w", "/work", "gcc:14.2.0", "g++", "-std=c++20", "-fsyntax-only", "access.cpp")
		output, err := cmd.CombinedOutput()
		if err == nil || !strings.Contains(string(output), "private") {
			t.Fatalf("expected native private-access rejection: %v\n%s", err, output)
		}
		t.Log("Native compiler rejected access to the independently authored default-private member.")
	})
}
