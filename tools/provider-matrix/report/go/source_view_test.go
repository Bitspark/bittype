package report

import (
	"go/ast"
	"path/filepath"
	"strings"
	"testing"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

func TestSourceViewKeepsPhysicalSourceAndCanonicalSymbols(t *testing.T) {
	const source = `package facts
type Box[T any] struct { Value T }
type Alias = Box[string]
type Defined Box[string]
const Label = "actual"
var Rules = Alias{Value: Label}
var Other = Defined{Value: Label}
func Helper(v string) string { return v }
//line not-the-physical-file.go:800
func (Box[T]) Apply(v string) string { return Helper(v) }
func Local() string { Label := "local"; return Label }
func Append(v []string) []string { return append(v, Label) }
func Shadow(v []string) []string { append := func(v []string, s string) []string { return v }; return append(v, Label) }
var Empty []string = nil
`
	p := fixture(t, map[string]string{"unrelated/facts": source})
	view := newSourceView(p)
	pkgPath := p.module + "/unrelated/facts"
	decl, err := view.Declaration(capability.Symbol{Package: pkgPath, Name: "Rules"})
	if err != nil {
		t.Fatal(err)
	}
	file, err := view.File(decl.Span.File)
	if err != nil || file.Text != source || !strings.HasSuffix(file.Path, filepath.Join("unrelated", "facts", "implementation.go")) {
		t.Fatal("source capture lost", file, err)
	}
	if !strings.HasPrefix(source[decl.Span.Start:decl.Span.End], "var Rules =") {
		t.Fatal("full declaration omitted its Go keyword")
	}
	typeDecl, err := view.Declaration(capability.Symbol{Package: pkgPath, Name: "Alias"})
	if err != nil || !strings.HasPrefix(source[typeDecl.Span.Start:typeDecl.Span.End], "type Alias =") {
		t.Fatal("type declaration omitted its Go keyword", err)
	}
	facts, err := view.Facts(decl.Value)
	if err != nil || facts.Type.Origin == nil || facts.Type.Origin.Name != "Box" || len(facts.Type.Arguments) != 1 || facts.Type.Arguments[0].Text != "string" {
		t.Fatal("alias must expose its real generic origin", facts, err)
	}
	facts.Type.Origin.Name = "corrupted"
	facts.Type.Arguments[0].Text = "corrupted"
	again, _ := view.Facts(decl.Value)
	if again.Type.Origin.Name != "Box" || again.Type.Arguments[0].Text != "string" {
		t.Fatal("mutable facts escaped the source view")
	}
	other, _ := view.Declaration(capability.Symbol{Package: pkgPath, Name: "Other"})
	otherFacts, _ := view.Facts(other.Value)
	if otherFacts.Type.Origin == nil || otherFacts.Type.Origin.Name != "Defined" {
		t.Fatal("new definition collapsed to underlying Box")
	}
	methodDecl, err := view.Declaration(capability.Symbol{Package: pkgPath, Receiver: "Box", Name: "Apply"})
	if err != nil || !strings.HasPrefix(source[methodDecl.Span.Start:methodDecl.Span.End], "func (Box[T]) Apply") {
		t.Fatal("physical method declaration changed by line directive", methodDecl, err)
	}
	if _, err := view.Facts(capability.Span{File: decl.Value.File, Start: decl.Value.Start, End: decl.Value.End - 1}); err == nil {
		t.Fatal("non-expression subspan accepted")
	}
	if _, err := newSourceView(p).Facts(decl.Value); err == nil {
		t.Fatal("foreign view span accepted")
	}
	if _, err := view.Declaration(capability.Symbol{Package: pkgPath, Name: "v"}); err == nil {
		t.Fatal("local parameter was exposed as package declaration")
	}
	var sawBuiltin, sawShadow, sawNil, sawLocal, sawCall bool
	pkg := p.packages[pkgPath]
	ast.Inspect(pkg.files[0], func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.CallExpr:
			facts, err := view.Facts(view.span(node))
			if err != nil {
				t.Fatal(err)
			}
			if id, ok := node.Fun.(*ast.Ident); ok {
				switch id.Name {
				case "append":
					if facts.Callee == nil {
						sawShadow = true
					} else if *facts.Callee == (capability.Symbol{Package: "builtin", Name: "append"}) {
						sawBuiltin = true
					} else {
						t.Fatal("wrong append identity", facts)
					}
				case "Helper":
					sawCall = facts.Callee != nil && facts.Callee.Package == pkgPath && facts.Callee.Name == "Helper"
				}
			}
		case *ast.Ident:
			facts, err := view.Facts(view.span(node))
			if err != nil {
				t.Fatal(err)
			}
			if node.Name == "nil" {
				sawNil = facts.Symbol != nil && *facts.Symbol == (capability.Symbol{Package: "builtin", Name: "nil"})
			}
			if node.Name == "Label" && pkg.info.Uses[node] != nil && pkg.info.Uses[node].Parent() != pkg.types.Scope() {
				sawLocal = facts.Symbol == nil
			}
		}
		return true
	})
	if !sawBuiltin || !sawShadow || !sawNil || !sawLocal || !sawCall {
		t.Fatalf("missing facts: builtin=%v shadow=%v nil=%v local=%v call=%v", sawBuiltin, sawShadow, sawNil, sawLocal, sawCall)
	}
}

func TestSourceViewResolvesImportedAliasesWithoutDisplayNames(t *testing.T) {
	p := fixture(t, map[string]string{
		"unrelated/one": "package same; const Value = 1; func Apply[T any](v T) T { return v }",
		"unrelated/two": "package same; const Value = 2",
		"unrelated/use": `package use
import (a "example.test/fixture/unrelated/one"; b "example.test/fixture/unrelated/two")
var A = a.Apply(a.Value)
var B = b.Value
`,
	})
	view := newSourceView(p)
	for _, test := range []struct { name, origin, value string }{{"A", "one", ""}, {"B", "two", "2"}} {
		decl, err := view.Declaration(capability.Symbol{Package: p.module + "/unrelated/use", Name: test.name})
		if err != nil {
			t.Fatal(err)
		}
		facts, err := view.Facts(decl.Value)
		if err != nil {
			t.Fatal(err)
		}
		identity := facts.Symbol
		if test.name == "A" {
			identity = facts.Callee
		}
		if identity == nil || identity.Package != p.module+"/unrelated/"+test.origin {
			t.Fatal("import alias or same package name collapsed identity", facts)
		}
		if test.value != "" && (facts.Constant == nil || facts.Constant.Value != test.value) {
			t.Fatal("constant fact missing", facts)
		}
	}
}
