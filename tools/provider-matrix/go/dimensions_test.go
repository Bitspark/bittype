package main

import (
	"go/token"
	"go/types"
	"path/filepath"
	"strings"
	"testing"
)

// None of these names or paths indicate a known provider, language or carrier.
// Only source-owned comments classify their operation boundaries.
const dimensionFixtureSource = `package elsewhere
// bittype:dimension Syntax
type TokenGraph struct{}
// bittype:dimension Files
type BoundResult struct { References TokenGraph }
func (BoundResult) Tree() TokenGraph { return TokenGraph{} }
type Hidden struct { References TokenGraph }
type Box[T any] struct{ Value T }
type Alias = TokenGraph
// bittype:dimension Artifacts
type SomethingNew struct{}
// bittype:dimension Files
type Sealed[T any] struct{ Value T }
type NewDefinition TokenGraph
type (
    // bittype:dimension Syntax
    Grouped struct{}
    Unmarked struct{}
)
/* bittype:dimension Files */
type BlockComment struct{}
// bittype:shared
type Identity interface { ID() string }
type Operation[I, O any] interface { Identity; Apply(I) O }
type Original interface { Read() TokenGraph }
type Derived interface { Original }
// bittype:dimension Files
type Contradictory interface { Read() TokenGraph }
type Holder interface { Apply(func(TokenGraph) BoundResult) }
`

func dimensionFixture(t *testing.T) (*program, *sourcePackage) {
	t.Helper()
	p := &program{root: t.TempDir(), module: "example.test/dimensions", fset: token.NewFileSet(), packages: map[string]*sourcePackage{}, functions: map[*types.Func]declaration{}}
	l := &fixtureLoader{p: p, sources: map[string]string{"unrelated/area": dimensionFixtureSource}}
	path := p.module + "/unrelated/area"
	if _, err := l.Import(path); err != nil {
		t.Fatal(err)
	}
	return p, p.packages[path]
}

func TestDimensionsFollowMetadataOnOperationBoundaries(t *testing.T) {
	p, pkg := dimensionFixture(t)
	get := func(name string) types.Type {
		obj := pkg.types.Scope().Lookup(name)
		if obj == nil {
			t.Fatalf("missing fixture type %s", name)
		}
		return obj.Type()
	}
	instantiate := func(name string, args ...types.Type) types.Type {
		instance, err := types.Instantiate(nil, get(name), args, true)
		if err != nil {
			t.Fatal(err)
		}
		return instance
	}
	for _, test := range []struct {
		name          string
		input, output types.Type
		want          string
	}{
		{"syntax boundary", get("TokenGraph"), types.Typ[types.String], "Syntax"},
		{"semantic result with AST storage and helpers", types.Typ[types.String], get("BoundResult"), "Files"},
		{"conflicting boundaries", get("TokenGraph"), get("BoundResult"), "Mixed: Files + Syntax"},
		{"unclassified storage", get("Hidden"), types.Typ[types.String], "Unclassified"},
		{"generic wrapper", types.Typ[types.String], instantiate("Box", get("BoundResult")), "Files"},
		{"classified generic carrier owns its boundary", types.Typ[types.String], instantiate("Sealed", get("TokenGraph")), "Files"},
		{"alias retains its target", get("Alias"), types.Typ[types.String], "Syntax"},
		{"new definition needs its own metadata", get("NewDefinition"), types.Typ[types.String], "Unclassified"},
		{"custom dimension label", get("SomethingNew"), types.Typ[types.String], "Artifacts"},
		{"pointer and slice boundaries", types.NewSlice(types.NewPointer(get("TokenGraph"))), types.Typ[types.String], "Syntax"},
		{"grouped type metadata", get("Grouped"), types.Typ[types.String], "Syntax"},
		{"untagged grouped sibling", get("Unmarked"), types.Typ[types.String], "Unclassified"},
		{"block documentation comment", get("BlockComment"), types.Typ[types.String], "Files"},
		{"shared input does not classify an operation", get("Identity"), types.Typ[types.String], "Unclassified"},
	} {
		t.Run(test.name, func(t *testing.T) {
			contract := instantiate("Operation", test.input, test.output).(*types.Named)
			if got := p.level(contract); got != test.want {
				t.Fatalf("got %q, want %q", got, test.want)
			}
		})
	}
	for name, want := range map[string]string{"Identity": "Shared", "Derived": "Syntax", "Contradictory": "Mixed: Files + Syntax", "Holder": "Mixed: Files + Syntax"} {
		if got := p.level(get(name).(*types.Named)); got != want {
			t.Errorf("%s: got %q, want %q", name, got, want)
		}
	}
}

func TestDimensionDiscoveryDoesNotGuessFromPackageNames(t *testing.T) {
	p, _ := dimensionFixture(t)
	// Recreate the exact path and type name recognized by the former heuristic,
	// but without metadata. The operation must remain unclassified.
	l := &fixtureLoader{p: p, sources: map[string]string{
		"trees/ast/go": "package ast; type AST struct{}; type Provider interface { Read() AST }",
	}}
	pkg, err := l.Import(p.module + "/trees/ast/go")
	if err != nil {
		t.Fatal(err)
	}
	if got := p.level(pkg.Scope().Lookup("Provider").Type().(*types.Named)); got != "Unclassified" {
		t.Fatal(got)
	}
}

func TestDimensionMetadataRejectsAmbiguousDeclarations(t *testing.T) {
	for name, source := range map[string]string{
		"duplicate":             "// bittype:dimension Syntax\n// bittype:dimension Files\ntype T struct{}",
		"shared and dimension":  "// bittype:shared\n// bittype:dimension Files\ntype T struct{}",
		"missing label":         "// bittype:dimension\ntype T struct{}",
		"extra labels":          "// bittype:dimension Files Syntax\ntype T struct{}",
		"reserved label":        "// bittype:dimension Shared\ntype T struct{}",
		"extra shared argument": "// bittype:shared Files\ntype T struct{}",
		"unknown directive":     "// bittype:dimesion Files\ntype T struct{}",
		"group annotation":      "// bittype:dimension Files\ntype (A struct{}; B struct{})",
		"alias annotation":      "type Original struct{}\n// bittype:dimension Files\ntype T = Original",
	} {
		t.Run(name, func(t *testing.T) {
			p := &program{root: t.TempDir(), module: "example.test/invalid", fset: token.NewFileSet(), packages: map[string]*sourcePackage{}, functions: map[*types.Func]declaration{}}
			l := &fixtureLoader{p: p, sources: map[string]string{"somewhere": "package invalid\n" + source}}
			_, err := l.Import(p.module + "/somewhere")
			if err == nil || !strings.Contains(err.Error(), filepath.Join("somewhere", "implementation.go")) {
				t.Fatalf("expected located metadata error, got %v", err)
			}
		})
	}
}
