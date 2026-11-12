package horos_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/semantics/go"
	horos "github.com/Bitspark/bittype/semantics/horos/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
)

// This file spells every published value of the vocabulary a second time, on
// purpose. Production code takes each name from the vocabulary, so a writer
// and its readers follow a changed value together and no other test notices.
// A published value is held by exported graphs, receipts and canonical module
// values. Changing one is a decision: make it here as well.

var publishedStrings = []struct {
	name, got, want string
}{
	{"ResolutionProfile", horos.ResolutionProfile, "horos-declarations-4/resolution-1"},
	{"AnalysisProfile", horos.AnalysisProfile, "horos-declarations-4/analysis-1"},
	{"DeclarationTagPrefix", horos.DeclarationTagPrefix, "horos/decl/"},
	{"DescriptorTagPrefix", horos.DescriptorTagPrefix, "horos/"},
}

var publishedKeys = []struct {
	name        string
	got         semantics.Key
	owner, want string
}{
	{"Key", horos.Key("any-name"), "horos", "any-name"},
	{"ReferenceResolution", horos.ReferenceResolution(), "horos", "reference-resolution"},
	{"Validity", horos.Validity(), "horos", "program-validity"},
	{"Indexing", horos.Indexing(), "horos", "declaration-index"},
	{"Bindings", horos.Bindings(), "horos", "occurrence-bindings"},
	{"Applications", horos.Applications(), "horos", "applications"},
	{"Dependencies", horos.Dependencies(), "horos", "reference-dependencies"},
	{"TypeFacts", horos.TypeFacts(), "horos", "type-facts"},
	{"TargetRole", horos.TargetRole(), "horos", "target"},
	{"ExpressionRole", horos.ExpressionRole(), "horos", "expression"},
	{"ReferenceDependency", horos.ReferenceDependency(), "horos", "declaration-reference"},
	{"SortRelation", horos.SortRelation(), "horos", "sort"},
	{"ResultSortRelation", horos.ResultSortRelation(), "horos", "result-sort"},
	{"DeclaredKindRelation", horos.DeclaredKindRelation(), "horos", "declared-kind"},
	{"ResultKindRelation", horos.ResultKindRelation(), "horos", "result-kind"},
	{"SymbolicBodyRelation", horos.SymbolicBodyRelation(), "horos", "symbolic-body"},
	{"MembershipRelation", horos.MembershipRelation(), "horos", "membership"},
	{"IdentityRelation", horos.IdentityRelation(), "horos", "canonical-identity"},
	{"ModuleKind", horos.ModuleKind(), "horos", "module"},
	{"TypeDefinitionKind", horos.TypeDefinitionKind(), "horos", "type-definition"},
	{"ImportRelation", horos.ImportRelation(), "horos", "import"},
	{"ProperSort", horos.ProperSort(), "horos", "sort/proper"},
	{"ModuleIdentity", horos.ModuleIdentity(), "horos", "identity/module"},
	{"DescriptorOperator", horos.DescriptorOperator("record"), "horos", "descriptor/record"},
	{"TermEncoding", horos.TermEncoding(), "ontos", "codec-v1"},
	{"LimitsSetting", horos.LimitsSetting(), "horos", "declaration-limits-json-1"},
}

// The exported names that the tests of this file cover without a row above.
var otherPublished = []string{
	"DescriptorConstructor", "PreludePin",
	"ResolutionContract", "AnalysisContract", "AnalysisContractWithResolutionConfiguration",
}

func TestPublishedNamesAreUnchanged(t *testing.T) {
	for _, c := range publishedStrings {
		if c.got != c.want {
			t.Errorf("%s = %q, published as %q", c.name, c.got, c.want)
		}
	}
	for _, c := range publishedKeys {
		if want := (semantics.Key{Owner: c.owner, Name: c.want}); c.got != want {
			t.Errorf("%s = %#v, published as %#v", c.name, c.got, want)
		}
	}
}

// The prelude pin is the horos library's and changes with the version that
// go.mod selects, so it has no literal here.
func TestPreludePinIsTheLibrarys(t *testing.T) {
	pin := horos.PreludePin()
	if pin != decl.PreludePin().String() {
		t.Fatalf("PreludePin = %q, the library's is %q", pin, decl.PreludePin().String())
	}
	if len(pin) != 64 || strings.Trim(pin, "0123456789abcdef") != "" {
		t.Fatalf("PreludePin = %q, want 64 lowercase hexadecimal digits", pin)
	}
}

// Every exported name of the package has a pin or a test in this file. A new
// one fails here until it gets one.
func TestEveryExportedNameIsCovered(t *testing.T) {
	covered := slices.Clone(otherPublished)
	for _, c := range publishedStrings {
		covered = append(covered, c.name)
	}
	for _, c := range publishedKeys {
		covered = append(covered, c.name)
	}
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	var exported []string
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
			continue
		}
		file, err := parser.ParseFile(token.NewFileSet(), entry.Name(), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range file.Decls {
			switch d := d.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil && d.Name.IsExported() {
					exported = append(exported, d.Name.Name)
				}
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch spec := spec.(type) {
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							if name.IsExported() {
								exported = append(exported, name.Name)
							}
						}
					case *ast.TypeSpec:
						if spec.Name.IsExported() {
							exported = append(exported, spec.Name.Name)
						}
					}
				}
			}
		}
	}
	if len(exported) == 0 {
		t.Fatal("found no exported name; the test must run in the package directory")
	}
	for _, name := range exported {
		if !slices.Contains(covered, name) {
			t.Errorf("exported %s has no pin in this file", name)
		}
	}
	for _, name := range covered {
		if !slices.Contains(exported, name) {
			t.Errorf("this file pins %s, which the package does not export", name)
		}
	}
}

func TestContractsAreUnchanged(t *testing.T) {
	key := func(name string) semantics.Key { return semantics.Key{Owner: "horos", Name: name} }
	type report struct {
		name    string
		verdict semantic.Verdict
	}
	reports := []report{
		{"reference-resolution", semantic.Holds},
		{"program-validity", semantic.Holds},
		{"declaration-index", semantic.Undetermined},
		{"occurrence-bindings", semantic.Undetermined},
		{"applications", semantic.Undetermined},
		{"reference-dependencies", semantic.Undetermined},
		{"type-facts", semantic.Undetermined},
	}
	check := func(what string, got semantic.Contract, profile string, want []report) {
		t.Helper()
		if got.Profile != profile {
			t.Errorf("%s profile = %q, published as %q", what, got.Profile, profile)
		}
		if len(got.Obligations) != len(want) || len(got.SemanticsRequirements) != len(want) {
			t.Fatalf("%s has %d obligations and %d requirements, published with %d", what, len(got.Obligations), len(got.SemanticsRequirements), len(want))
		}
		for i, r := range want {
			if o := (semantic.Obligation{Key: key(r.name), Domain: key("entry-transitive-module-closure")}); got.Obligations[i] != o {
				t.Errorf("%s obligation %d = %#v, published as %#v", what, i, got.Obligations[i], o)
			}
			if q := (semantic.Requirement{Obligation: key(r.name), Verdict: r.verdict}); got.SemanticsRequirements[i] != q {
				t.Errorf("%s requirement %d = %#v, published as %#v", what, i, got.SemanticsRequirements[i], q)
			}
		}
	}

	resolved := horos.ResolutionContract()
	check("resolution", resolved.Semantic, "horos-declarations-4/resolution-1", reports)
	if resolved.Bindings != key("reference-resolution") {
		t.Errorf("resolution bindings = %#v", resolved.Bindings)
	}

	// Analysis repeats every obligation but the bindings, and requires the
	// resolution with all of its reports and the given configuration.
	setting := semantic.Setting{Key: key("declaration-limits-json-1"), Value: []byte("[1,2,3,4]")}
	for _, c := range []struct {
		what     string
		got      analysis.Contract
		settings []semantic.Setting
	}{
		{"analysis", horos.AnalysisContract(), nil},
		{"configured analysis", horos.AnalysisContractWithResolutionConfiguration([]semantic.Setting{setting}), []semantic.Setting{setting}},
	} {
		check(c.what, c.got.Semantic, "horos-declarations-4/analysis-1", reports[1:])
		required := c.got.Resolution
		if required == nil {
			t.Fatalf("%s requires no resolution", c.what)
		}
		if required.Profile != "horos-declarations-4/resolution-1" {
			t.Errorf("%s requires resolution profile %q", c.what, required.Profile)
		}
		if len(required.Configuration) != len(c.settings) {
			t.Fatalf("%s carries %d settings, want %d", c.what, len(required.Configuration), len(c.settings))
		}
		for i, s := range c.settings {
			if got := required.Configuration[i]; got.Key != s.Key || string(got.Value) != string(s.Value) {
				t.Errorf("%s setting %d = %#v, want %#v", c.what, i, got, s)
			}
		}
		if len(required.Reports) != len(reports) {
			t.Fatalf("%s requires %d resolution reports, published with %d", c.what, len(required.Reports), len(reports))
		}
		for i, r := range reports {
			if q := (semantic.Requirement{Obligation: key(r.name), Verdict: r.verdict}); required.Reports[i] != q {
				t.Errorf("%s resolution report %d = %#v, published as %#v", c.what, i, required.Reports[i], q)
			}
		}
	}
}

// The constructors of the descriptors that the horos parser writes.
var constructors = []string{
	"any", "atom", "scalar", "tuple", "list-of", "seq", "set-of", "map-of", "union", "and", "not", "wire",
	"exact", "atom-exact", "atom-in", "atom-len", "text-len", "key-grammar", "text-non-blank", "lowercase-uuid",
	"compound", "rec", "ref", "record",
}

func TestDescriptorOperatorAndConstructorAreInverse(t *testing.T) {
	for _, constructor := range append(slices.Clone(constructors), "", "descriptor/record", "a/b", "über", " ", "RECORD") {
		operator := horos.DescriptorOperator(constructor)
		if want := (semantics.Key{Owner: "horos", Name: "descriptor/" + constructor}); operator != want {
			t.Errorf("DescriptorOperator(%q) = %#v, published as %#v", constructor, operator, want)
		}
		if got, ok := horos.DescriptorConstructor(operator); !ok || got != constructor {
			t.Errorf("DescriptorConstructor(DescriptorOperator(%q)) = %q, %v", constructor, got, ok)
		}
	}
}

func TestDescriptorConstructorReadsDescriptorOperatorsOnly(t *testing.T) {
	for _, c := range []struct {
		what        string
		operator    semantics.Key
		constructor string
		ok          bool
	}{
		{"a constructor", semantics.Key{Owner: "horos", Name: "descriptor/record"}, "record", true},
		{"the bare prefix", semantics.Key{Owner: "horos", Name: "descriptor/"}, "", true},
		{"a nested prefix", semantics.Key{Owner: "horos", Name: "descriptor/descriptor/record"}, "descriptor/record", true},
		{"the prefix without its separator", semantics.Key{Owner: "horos", Name: "descriptor"}, "", false},
		{"a prefix that is not at the start", semantics.Key{Owner: "horos", Name: "x/descriptor/record"}, "", false},
		{"another case", semantics.Key{Owner: "horos", Name: "Descriptor/record"}, "", false},
		{"another owner", semantics.Key{Owner: "ontos", Name: "descriptor/record"}, "", false},
		{"no owner", semantics.Key{Name: "descriptor/record"}, "", false},
		{"an owner that only starts alike", semantics.Key{Owner: "horos/files", Name: "descriptor/record"}, "", false},
		{"a kind operator", horos.Key("kind/node"), "", false},
		{"the proper sort", horos.ProperSort(), "", false},
		{"the module identity", horos.ModuleIdentity(), "", false},
		{"the term encoding", horos.TermEncoding(), "", false},
		{"the zero key", semantics.Key{}, "", false},
	} {
		constructor, ok := horos.DescriptorConstructor(c.operator)
		if ok != c.ok || constructor != c.constructor {
			t.Errorf("%s: DescriptorConstructor(%#v) = %q, %v, want %q, %v", c.what, c.operator, constructor, ok, c.constructor, c.ok)
		}
	}
}
