package report

import (
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestReviewAnonymousArgumentOwnerIsNotDisplayIdentity(t *testing.T) {
	p := fixture(t, map[string]string{
		"provider/ember/go": `package emberprovider
import (s "example.test/fixture/syntax/ember/go"; tree "example.test/fixture/trees/cst/go")
type Output struct{s.Language}
func(Output)Print(tree.CST[s.Language])(struct{value int},error){return struct{value int}{},nil}
`,
		"print/files/go": `package filesprint
import (id "example.test/fixture/contracts/language/go"; printing "example.test/fixture/contracts/print/go"; tree "example.test/fixture/trees/cst/go"; source "example.test/fixture/source/go"; files "example.test/fixture/files/go")
type Output[L any] struct {identity id.Provider;syntax printing.Provider[tree.CST[L],struct{value int}]}
func New[L any](identity id.Provider,syntax printing.Provider[tree.CST[L],struct{value int}])(Output[L],error){return Output[L]{identity,syntax},nil}
func(p Output[L])ID()id.ID{return p.identity.ID()}
func(p Output[L])Name()string{return p.identity.Name()}
func(p Output[L])Print(source.Source[L])(files.Files,error){return files.Files{},nil}
`,
	})
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	files := m.rows["print.Provider[source.Source[L], files.Files]"]
	if files == nil || len(files.cells["ember"]) != 0 {
		t.Fatal("a provider returning another package's unexported struct field incorrectly satisfies the Files constructor")
	}
}

func TestDependencyArgumentsKeepGoStructuralIdentity(t *testing.T) {
	for _, test := range []struct {
		name, output string
		available    bool
	}{
		{"exported struct", "struct{Value int}", true},
		{"private interface", "interface{value() int}", false},
		{"exported interface", "interface{Value() int}", true},
		{"nested private struct", "map[string]func() []struct{value int}", false},
		{"shared private alias", "shared.Payload", true},
	} {
		t.Run(test.name, func(t *testing.T) {
			extra := ""
			if test.output == "shared.Payload" {
				extra = `; shared "example.test/fixture/shared"`
			}
			p := fixture(t, map[string]string{
				"shared": `package shared; type Payload = struct{value int}`,
				"provider/ember/go": `package emberprovider
import(s "example.test/fixture/syntax/ember/go"; tree "example.test/fixture/trees/cst/go"` + extra + `)
type Output struct{s.Language}
func(Output)Print(tree.CST[s.Language])(` + test.output + `,error){panic("source only")}`,
				"print/files/go": `package filesprint
import(id "example.test/fixture/contracts/language/go"; printing "example.test/fixture/contracts/print/go"; tree "example.test/fixture/trees/cst/go"; source "example.test/fixture/source/go"; files "example.test/fixture/files/go"` + extra + `)
type Output[L any] struct{identity id.Provider}
func New[L any](identity id.Provider, syntax printing.Provider[tree.CST[L],` + test.output + `])(Output[L],error){return Output[L]{identity},nil}
func(p Output[L])ID()id.ID{return p.identity.ID()}
func(p Output[L])Name()string{return p.identity.Name()}
func(p Output[L])Print(source.Source[L])(files.Files,error){return files.Files{},nil}`,
			})
			m, err := p.discover()
			if err != nil {
				t.Fatal(err)
			}
			row := m.rows["print.Provider[source.Source[L], files.Files]"]
			if row == nil || (len(row.cells["ember"]) > 0) != test.available {
				t.Fatalf("dependency availability does not follow Go identity: %#v", row)
			}
		})
	}
}

func TestRowsKeepUnnamedMemberOwnership(t *testing.T) {
	for _, test := range []struct {
		shape      string
		compatible bool
	}{
		{"struct{value []LANG.Language}", false},
		{"struct{Value []LANG.Language}", true},
		{"interface{value(LANG.Language)}", false},
		{"interface{Value(LANG.Language)}", true},
		{"map[string]func(...LANG.Language) []struct{value int}", false},
		{"map[string]func(...LANG.Language) []struct{Value int}", true},
	} {
		t.Run(test.shape, func(t *testing.T) {
			sources := map[string]string{}
			for _, name := range []string{"ember", "ice"} {
				sources["provider/"+name+"/go"] = `package output
import(s "example.test/fixture/syntax/` + name + `/go"; tree "example.test/fixture/trees/cst/go")
type Output struct{s.Language}
func(Output)Print(tree.CST[s.Language])(` + strings.ReplaceAll(test.shape, "LANG", "s") + `,error){panic("source only")}`
			}
			p := fixture(t, sources)
			m, err := p.discover()
			if !test.compatible {
				if err == nil || !strings.Contains(err.Error(), "ambiguous matrix contract label") {
					t.Fatal("private owners merged into one language-generalized row", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			joined := false
			for _, row := range m.rows {
				joined = joined || strings.HasPrefix(row.contract, "print.Provider[cst.") && len(row.cells["ember"]) == 1 && len(row.cells["ice"]) == 1
			}
			if !joined {
				t.Fatal("identical exported structures failed to share a row")
			}
		})
	}
}

func TestRowAbstractionPreservesNonMarkerIdentity(t *testing.T) {
	decl := fixture(t).packages["example.test/fixture/syntax/ember/go"].types.Scope().Lookup("Language").Type().(*types.Named)
	other := types.NewNamed(types.NewTypeName(token.NoPos, decl.Obj().Pkg(), "LanguageExtra", nil), types.NewStruct(nil, nil), nil)
	parameter := rowParameter()
	for _, value := range []types.Type{other, types.NewPointer(other), types.NewSlice(other), types.NewArray(other, 2), types.NewMap(types.Typ[types.String], other), types.NewChan(types.SendOnly, other)} {
		normalized, err := rowType(value, decl, parameter)
		if err != nil || !types.Identical(value, normalized) {
			t.Fatalf("changed non-marker type: %v %v", value, err)
		}
	}
	pkg := types.NewPackage("example.test/owner", "owner")
	input := types.NewStruct([]*types.Var{types.NewField(token.NoPos, pkg, "value", decl, false)}, []string{`json:"value"`})
	got, err := rowType(input, decl, parameter)
	if err != nil {
		t.Fatal(err)
	}
	expected := types.NewStruct([]*types.Var{types.NewField(token.NoPos, pkg, "value", parameter, false)}, []string{`json:"value"`})
	if !types.Identical(got, expected) {
		t.Fatal("row substitution lost private owner or field tag")
	}
	wrongTag := types.NewStruct([]*types.Var{types.NewField(token.NoPos, pkg, "value", parameter, false)}, []string{`json:"other"`})
	if types.Identical(got, wrongTag) {
		t.Fatal("field tag disappeared")
	}
}

func TestReportAssociationsDoNotUseTypeLabels(t *testing.T) {
	p := fixture(t)
	m, err := p.discover()
	if err != nil {
		t.Fatal(err)
	}
	base := m.rows["print.Provider[cst.CST[L], string]"].cells["ember"][0].op
	private := func(path string) types.Type {
		return types.NewStruct([]*types.Var{types.NewField(token.NoPos, types.NewPackage(path, "owner"), "value", types.Typ[types.Int], false)}, nil)
	}
	a, b := private("example.test/a"), private("example.test/b")
	if types.Identical(a, b) || typeLabel(a) != typeLabel(b) {
		t.Fatal("fixture must distinguish Go identity from type spelling")
	}
	makeOperation := func(argument types.Type) *operation {
		result := *base
		bound, err := types.Instantiate(nil, base.contract.Origin(), []types.Type{base.contract.TypeArgs().At(0), argument}, true)
		if err != nil {
			t.Fatal(err)
		}
		result.contract = bound.(*types.Named)
		result.id = operationID(result.kind, result.public, result.contract, result.receiver, result.language)
		return &result
	}
	first, second := makeOperation(a), makeOperation(b)
	if first.id != second.id {
		t.Fatal("fixture must exercise a rendered anchor collision")
	}
	m = &matrix{rows: map[string]*row{"test": {cells: map[string][]location{"ember": {{op: first}, {op: second}}}}}}
	if len(operations(m)) != 2 {
		t.Fatal("unequal operations merged by anchor spelling")
	}
	if err := validateOperationAnchors(m); err == nil {
		t.Fatal("ambiguous anchor accepted")
	}
	target := canonicalSymbol(first.public)
	annotations := &annotations{cases: []caseAssociation{{target: target, test: testDeclaration{file: "test.go"}}}}
	if err := associateCases(m, annotations); err == nil || !strings.Contains(err.Error(), "ambiguous case target") {
		t.Fatal("unequal contract contexts merged", err)
	}
	second.contract = first.contract
	if err := associateCases(m, annotations); err != nil {
		t.Fatal("identical context rejected", err)
	}
	if len(operations(m)) != 1 {
		t.Fatal("identical operations did not deduplicate")
	}
	if !sameCallLink(callLink{arguments: []types.Type{a}}, callLink{arguments: []types.Type{a}}) || sameCallLink(callLink{arguments: []types.Type{a}}, callLink{arguments: []types.Type{b}}) {
		t.Fatal("call argument association lost Go identity")
	}
}
