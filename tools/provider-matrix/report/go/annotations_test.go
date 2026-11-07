package report

import (
	"context"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	capability "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
)

func addTestFile(t *testing.T, p *program, relative, name, text string, external bool) {
	t.Helper()
	path := p.module + "/" + relative
	if p.metadata == nil {
		p.metadata = map[string]packageInfo{}
	}
	item := p.metadata[path]
	item.ImportPath = path
	item.Dir = filepath.Join(p.root, relative)
	item.Name = "fixture"
	item.Module = &struct {
		Path, Dir, GoVersion string
		Main                 bool
	}{Path: p.module, Dir: p.root, Main: true}
	if external {
		item.XTestGoFiles = append(item.XTestGoFiles, name)
	} else {
		item.TestGoFiles = append(item.TestGoFiles, name)
	}
	p.metadata[path] = item
	if err := os.MkdirAll(item.Dir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(item.Dir, name), []byte(text), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestTestDeclarationsStayOutsideProductionDiscovery(t *testing.T) {
	p := fixture(t)
	addTestFile(t, p, "print/syntax/ember/provider/go", "internal_test.go", `package emberprovider
import t "testing"
// bittype:case target=emit feature=printing case=internal expects=round-trip
func TestSame(v *t.T) {}
type TestOnly struct{}
func(TestOnly)ID()string{return "ice"}
func(TestOnly)Name()string{return "Ice"}
func(TestOnly)Parse(string)(string,error){return "",nil}
`, false)
	addTestFile(t, p, "print/syntax/ember/provider/go", "external_test.go", `package emberprovider_test
import("testing";ep "example.test/fixture/print/syntax/ember/provider/go")
// bittype:case target=ep.Output.Print feature=printing case=external expects=round-trip
func TestSame(t *testing.T) {_ = ep.Output{}}
`, true)
	a, err := p.readAnnotations()
	if err != nil {
		t.Fatal(err)
	}
	if len(a.tests) != 2 || len(a.cases) != 2 || a.tests[0].id == a.tests[1].id || a.tests[0].pkg != a.tests[1].pkg {
		t.Fatalf("lost test variants: %#v", a.tests)
	}
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	if err := associateCases(m, a); err != nil {
		t.Fatal(err)
	}
	for _, row := range m.rows {
		if strings.HasPrefix(row.contract, "parse.") {
			t.Fatal("test-only provider entered discovery")
		}
	}
}

func TestExternalTestReferencesAndSignatures(t *testing.T) {
	for _, body := range []string{
		`// bittype:case target=emit feature=f case=c expects=ok
func TestCase(t *testing.T){}`,
		`// bittype:case target=example.test/fixture/print/syntax/ember/provider/go#emit feature=f case=c expects=ok
func TestCase(t,u *testing.T){}`,
		`type Alias = testing.T
// bittype:case target=example.test/fixture/print/syntax/ember/provider/go#emit feature=f case=c expects=ok
func TestCase(t *Alias){}`,
	} {
		t.Run(body[:min(35, len(body))], func(t *testing.T) {
			p := fixture(t)
			addTestFile(t, p, "print/syntax/ember/provider/go", "bad_test.go", "package emberprovider_test\nimport \"testing\"\n"+body, true)
			if _, err := p.readAnnotations(); err == nil {
				t.Fatal("unsupported/ambiguous annotation accepted")
			}
		})
	}
	file, err := parser.ParseFile(token.NewFileSet(), "test.go", "package p; import t \"testing\"; func TestValid(v *(t.T)) {}", 0)
	if err != nil {
		t.Fatal(err)
	}
	fn := file.Decls[1].(*ast.FuncDecl)
	if !testFunction(fn, map[string]string{"t": "testing"}) {
		t.Fatal("direct parenthesized testing.T rejected")
	}
}

func TestGenericCaseNeedsActualMarker(t *testing.T) {
	p := fixture(t, map[string]string{"provider/ice/go": `package iceprovider
import(s "example.test/fixture/syntax/ice/go";tree "example.test/fixture/trees/cst/go")
type Output struct{s.Language}
func(Output)Print(tree.CST[s.Language])(string,error){return "",nil}
`})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	target := capability.Symbol{Package: p.module + "/print/files/shared/provider/go", Receiver: "Output", Name: "Print"}
	a := &annotations{cases: []caseAssociation{{target: target, test: testDeclaration{file: "test.go"}}}}
	if err := associateCases(m, a); err == nil {
		t.Fatal("unqualified shared target credited across markers")
	}
	a.cases[0].marker = capability.Symbol{Package: p.module + "/syntax/ember/go", Name: "Language"}
	if err := associateCases(m, a); err != nil {
		t.Fatal(err)
	}
	a.cases[0].marker.Name = "Fake"
	if err := associateCases(m, a); err == nil {
		t.Fatal("unrecognized language qualifier accepted")
	}
}

func TestQueryLeaseAndProfileIdentity(t *testing.T) {
	p := fixture(t)
	view := &queryView{base: newSourceView(p), ctx: context.Background(), remaining: 1}
	symbol := capability.Symbol{Package: p.module + "/print/syntax/ember/provider/go", Name: "emit"}
	if _, err := view.Declaration(symbol); err != nil {
		t.Fatal(err)
	}
	if _, err := view.Declaration(symbol); err == nil {
		t.Fatal("query budget not enforced")
	}
	view.remaining = 1
	view.closed = true
	if _, err := view.Declaration(symbol); err == nil {
		t.Fatal("retained query lease remained active")
	}
	if stableKey("ab", "c") == stableKey("a", "bc") || stableKey("false", "") == stableKey("true", "unspecified") {
		t.Fatal("identity tuple collision")
	}
	if profileLabel(profile{}) == profileLabel(profile{true, "unspecified"}) {
		t.Fatal("absent profile conflated with named profile")
	}
	op := &operation{id: "operation"}
	entry := authored{directive: directive{name: "grammar"}}
	absent := bindingAnchor(op, symbol, entry)
	entry.directive.profile = profile{present: true, value: "unspecified"}
	if absent == bindingAnchor(op, symbol, entry) {
		t.Fatal("absent profile collided with explicit binding anchor")
	}
}

func TestMarkerAliasesRetainActualGoIdentity(t *testing.T) {
	for _, name := range []string{"Alias", "Pointer"} {
		t.Run(name, func(t *testing.T) {
			p := fixture(t, map[string]string{"unrelated/aliases": `package aliases; import s "example.test/fixture/syntax/ember/go"; type Alias = s.Language; type Pointer = *s.Language`})
			text := "package emberprovider_test\nimport \"testing\"\n// bittype:case target=example.test/fixture/print/files/shared/provider/go#Output.Print feature=f case=c expects=ok language=example.test/fixture/unrelated/aliases#" + name + "\nfunc TestAlias(t *testing.T){}\n"
			addTestFile(t, p, "print/syntax/ember/provider/go", "alias_test.go", text, true)
			a, err := p.readAnnotations()
			if name == "Pointer" {
				if err == nil {
					t.Fatal("pointer substituted for marker")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if a.cases[0].marker != (capability.Symbol{Package: p.module + "/syntax/ember/go", Name: "Language"}) {
				t.Fatal("alias spelling replaced Go identity")
			}
		})
	}
}
