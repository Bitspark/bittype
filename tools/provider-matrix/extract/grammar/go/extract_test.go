package grammar_test

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"testing"

	"github.com/Bitspark/bittype/contracts/text/go"
	g "github.com/Bitspark/bittype/internal/grammar"
	"github.com/Bitspark/bittype/model/syntax/go"
	c "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
	extract "github.com/Bitspark/bittype/tools/provider-matrix/extract/grammar/go"
	carrier "github.com/Bitspark/bittype/trees/ast/go"
)

const grammarPath = "github.com/Bitspark/bittype/model/syntax/go"
const fixturePath = "example.test/fixture"

// These declarations provide compiler identities; they are never executed.
// What the maintained table type means is separately checked below against
// the interpreter. This stand-in is keyed by strings. The extractor reads
// source syntax and compiler facts, which a language parameter does not
// change; the typed stand-in below mirrors the owning package and has its own
// test.
const grammarSource = `package syntax
type Predicate struct { Owner, Name string }
type Attribute struct { Required bool; Predicate Predicate; OneOf []string }
type Role struct { Name string; Min,Max int; Kinds []string }
type Rule struct { Attributes map[string]Attribute; Roles []Role; Reference bool }
type Rules map[string]Rule
func Bool() Predicate { return Predicate{"bittype","bool"} }
func Nonempty() Predicate { return Predicate{"bittype","nonempty"} }
`

// language marks the maintained table that the tests below build. The owning
// table is typed by language; a fixture's table is source text and never runs.
type language struct{}

func names[K ~string](kinds []K) []string {
	result := make([]string, len(kinds))
	for i, kind := range kinds {
		result[i] = string(kind)
	}
	return result
}

// The typed declarations mirror the owning package: the table is keyed by one
// language's node kinds, and role and attribute names are strings.
const kindSource = `package ast
type Kind[L any] string
`
const kindPath = "example.test/ast"
const typedGrammarSource = `package syntax
import "` + kindPath + `"
type Predicate struct { Owner, Name string }
type Attribute struct { Required bool; Predicate Predicate; OneOf []string }
type Role[L any] struct { Name string; Min,Max int; Kinds []ast.Kind[L] }
type Rule[L any] struct { Attributes map[string]Attribute; Roles []Role[L]; Reference bool }
type Rules[L any] map[ast.Kind[L]]Rule[L]
func Bool() Predicate { return Predicate{"bittype","bool"} }
`
const vocabularyPath = "example.test/vocabulary"
const vocabularySource = `package vocabulary
import "` + kindPath + `"
type Language struct{}
type Kind = ast.Kind[Language]
const (
 KindRecord Kind = "record"
 KindField Kind = "field"
 KindAtom Kind = "atom"
)
const (
 RoleField = "field"
 RoleType = "type"
)
const AttributeName = "name"
`

type fixtureSource struct {
	set          *token.FileSet
	files        map[c.FileRef]c.File
	declarations map[c.Symbol]c.Declaration
	facts        map[c.Span]c.Facts
	packages     map[string]*types.Package
}

func (s *fixtureSource) File(ref c.FileRef) (c.File, error) {
	if file, ok := s.files[ref]; ok {
		return file, nil
	}
	return c.File{}, fmt.Errorf("foreign file")
}
func (s *fixtureSource) Declaration(symbol c.Symbol) (c.Declaration, error) {
	if declaration, ok := s.declarations[symbol]; ok {
		return declaration, nil
	}
	return c.Declaration{}, fmt.Errorf("unavailable declaration %s", symbol)
}
func (s *fixtureSource) Facts(span c.Span) (c.Facts, error) {
	if facts, ok := s.facts[span]; ok {
		return facts, nil
	}
	return c.Facts{}, fmt.Errorf("not an expression: %#v", span)
}
func (s *fixtureSource) Import(path string) (*types.Package, error) {
	if pkg := s.packages[path]; pkg != nil {
		return pkg, nil
	}
	return nil, fmt.Errorf("missing import %s", path)
}

func symbol(obj types.Object) *c.Symbol {
	if obj == nil {
		return nil
	}
	if obj.Pkg() == nil {
		if obj.Parent() == types.Universe {
			return &c.Symbol{Package: "builtin", Name: obj.Name()}
		}
		return nil
	}
	if obj.Parent() != obj.Pkg().Scope() {
		return nil
	}
	return &c.Symbol{Package: obj.Pkg().Path(), Name: obj.Name()}
}

func (s *fixtureSource) add(t *testing.T, path, text string) {
	t.Helper()
	file, err := parser.ParseFile(s.set, path+"/source.go", text, 0)
	if err != nil {
		t.Fatal(err)
	}
	info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{}}
	config := types.Config{Importer: s}
	pkg, err := config.Check(path, s.set, []*ast.File{file}, info)
	if err != nil {
		t.Fatal(err)
	}
	s.packages[path] = pkg
	ref := c.FileRef(path)
	s.files[ref] = c.File{Ref: ref, Path: path + "/source.go", Text: text}
	tokenFile := s.set.File(file.Pos())
	span := func(n ast.Node) c.Span {
		if n == nil {
			return c.Span{}
		}
		return c.Span{File: ref, Start: tokenFile.Offset(n.Pos()), End: tokenFile.Offset(n.End())}
	}
	var object func(ast.Expr) types.Object
	object = func(expr ast.Expr) types.Object {
		switch n := expr.(type) {
		case *ast.Ident:
			return info.Uses[n]
		case *ast.SelectorExpr:
			return info.Uses[n.Sel]
		case *ast.ParenExpr:
			return object(n.X)
		}
		return nil
	}
	for _, declaration := range file.Decls {
		switch n := declaration.(type) {
		case *ast.FuncDecl:
			id := symbol(info.Defs[n.Name])
			s.declarations[*id] = c.Declaration{Symbol: *id, Span: span(n), Type: span(n.Type)}
		case *ast.GenDecl:
			for _, spec := range n.Specs {
				switch n := spec.(type) {
				case *ast.ValueSpec:
					for i, name := range n.Names {
						id := symbol(info.Defs[name])
						var value c.Span
						if len(n.Values) == len(n.Names) {
							value = span(n.Values[i])
						}
						s.declarations[*id] = c.Declaration{Symbol: *id, Span: span(n), Value: value, Type: span(n.Type)}
					}
				case *ast.TypeSpec:
					id := symbol(info.Defs[n.Name])
					s.declarations[*id] = c.Declaration{Symbol: *id, Span: span(n), Type: span(n.Type)}
				}
			}
		}
	}
	ast.Inspect(file, func(node ast.Node) bool {
		expr, ok := node.(ast.Expr)
		if !ok {
			return true
		}
		facts := c.Facts{Symbol: symbol(object(expr))}
		if call, ok := expr.(*ast.CallExpr); ok {
			facts.Callee = symbol(object(call.Fun))
		}
		if typ := info.TypeOf(expr); typ != nil {
			facts.Type.Text = types.TypeString(typ, nil)
			if named, ok := types.Unalias(typ).(*types.Named); ok {
				facts.Type.Origin = symbol(named.Origin().Obj())
			}
		}
		if value := info.Types[expr].Value; value != nil {
			text := value.ExactString()
			if value.Kind() == constant.String {
				text = constant.StringVal(value)
			}
			facts.Constant = &c.Constant{Kind: strings.ToLower(value.Kind().String()), Value: text}
		}
		s.facts[span(expr)] = facts
		return true
	})
}

func request(t *testing.T, body string) c.DefinitionRequest {
	t.Helper()
	source := &fixtureSource{set: token.NewFileSet(), files: map[c.FileRef]c.File{}, declarations: map[c.Symbol]c.Declaration{}, facts: map[c.Span]c.Facts{}, packages: map[string]*types.Package{}}
	source.add(t, grammarPath, grammarSource)
	source.add(t, fixturePath, "package fixture\nimport g \""+grammarPath+"\"\n"+body)
	return c.DefinitionRequest{Root: c.Symbol{Package: fixturePath, Name: "rules"}, Source: source}
}

func inventory(t *testing.T, body string) c.Inventory {
	t.Helper()
	result, err := (extract.Extractor{}).Extract(context.Background(), request(t, body))
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestAliasesCompositionAndLocalKnowledge(t *testing.T) {
	result := inventory(t, `
type Alias = g.Rules
const atom = "atom"
var kinds = []string{atom}
var modes = []string{"one","two"}
func prefix() []g.Role { return []g.Role{{Name:"items",Max:-1,Kinds:kinds}} }
func with(base []g.Role, more ...g.Role) []g.Role { return append(base,more...) }
func custom() g.Predicate { return g.Predicate{"fixture","custom"} }
func named()map[string]g.Attribute{return map[string]g.Attribute{"mode":{OneOf:modes},"free":{},"custom":{Required:true,Predicate:custom()},"flag":{Predicate:g.Bool()}}}
var rules = Alias{
 "record": {Attributes:named(), Roles:with(prefix(),g.Role{Name:"next",Max:1,Kinds:append([]string{},kinds...)}),Reference:true},
 "empty": {},
}`)
	if result.Extent != c.PartlyDecoded || len(result.Features) != 2 {
		t.Fatalf("extent/features: %#v", result)
	}
	record := result.Features[0]
	if record.Key != "record" || record.Roles.State != c.Known || len(record.Roles.Value) != 2 {
		t.Fatalf("roles: %#v", record)
	}
	items := record.Roles.Value[0]
	maintained := syntax.Role[language]{Name: "items", Max: -1, Kinds: []carrier.Kind[language]{"atom"}}
	if items.Name != maintained.Name || items.Min.Value != maintained.Min || !items.Max.Value.Unbounded || maintained.Max >= 0 || strings.Join(items.Kinds.Value, ",") != strings.Join(names(maintained.Kinds), ",") {
		t.Fatalf("table type semantics: %#v", items)
	}
	if record.Roles.Value[1].Max.Value.Value != 1 || record.References.Value[0] != "target" {
		t.Fatalf("optional/reference: %#v", record)
	}
	attrs := record.Attributes.Value
	if record.Attributes.State != c.Known || len(attrs) != 4 {
		t.Fatalf("attributes: %#v", record.Attributes)
	}
	if attrs[0].Required.State != c.Known || attrs[0].Required.Value || strings.Join(attrs[0].Values.Value.Values, ",") != "one,two" {
		t.Fatalf("finite constraint: %#v", attrs[0])
	}
	if !attrs[1].Values.Value.Unrestricted || attrs[2].Values.State != c.Opaque || !attrs[2].Required.Value {
		t.Fatalf("local knowledge: %#v", attrs)
	}
	if strings.Join(attrs[3].Values.Value.Values, ",") != "true,false" || !g.Bool("true") || !g.Bool("false") || g.Bool("other") {
		t.Fatalf("Boolean predicate: %#v", attrs[3])
	}
	empty := result.Features[1]
	if empty.Roles.State != c.Known || len(empty.Roles.Value) != 0 || len(empty.Attributes.Value) != 0 || len(empty.References.Value) != 0 {
		t.Fatalf("known absence: %#v", empty)
	}
}

// The extractor reads a table from source. The interpreter applies the same
// table as a value. Both must say the same of every role and of every
// attribute whose constraint the extractor interprets.
func TestInterpretationMatchesTheInterpreter(t *testing.T) {
	result := inventory(t, `var rules=g.Rules{"x":{
 Roles:[]g.Role{{Name:"one",Min:1,Max:1,Kinds:[]string{"a"}},{Name:"optional",Max:1,Kinds:[]string{"a","b"}},{Name:"many",Min:2,Max:-1,Kinds:[]string{"b"}},{Name:"none"}},
 Attributes:map[string]g.Attribute{"finite":{Required:true,OneOf:[]string{"a","b"}},"empty":{OneOf:[]string{}},"free":{},"zero":{Predicate:g.Predicate{}},"flag":{Predicate:g.Bool()},"both":{Required:true,Predicate:g.Bool(),OneOf:[]string{"true","after"}}},Reference:true,
}}`)
	maintained := syntax.Rule[language]{
		Roles: []syntax.Role[language]{{Name: "one", Min: 1, Max: 1, Kinds: []carrier.Kind[language]{"a"}}, {Name: "optional", Max: 1, Kinds: []carrier.Kind[language]{"a", "b"}}, {Name: "many", Min: 2, Max: -1, Kinds: []carrier.Kind[language]{"b"}}, {Name: "none"}},
		Attributes: map[string]syntax.Attribute{
			"finite": {Required: true, OneOf: []string{"a", "b"}}, "empty": {OneOf: []string{}}, "free": {}, "zero": {Predicate: syntax.Predicate{}},
			"flag": {Predicate: syntax.Bool()}, "both": {Required: true, Predicate: syntax.Bool(), OneOf: []string{"true", "after"}},
		},
		Reference: true,
	}
	feature := result.Features[0]
	if result.Extent != c.Decoded || len(feature.Roles.Value) != len(maintained.Roles) || len(feature.Attributes.Value) != len(maintained.Attributes) {
		t.Fatalf("unexpected interpretation: %#v", result)
	}
	for i, actual := range maintained.Roles {
		model := feature.Roles.Value[i]
		if model.Name != actual.Name || model.Min.Value != actual.Min || model.Max.Value.Unbounded != (actual.Max < 0) || actual.Max >= 0 && model.Max.Value.Value != actual.Max || !slices.Equal(model.Kinds.Value, names(actual.Kinds)) {
			t.Fatalf("role %d disagrees: %#v versus %#v", i, model, actual)
		}
	}
	// The interpreter judges one attribute of a tree of one node.
	check := func(name, sample string) error {
		rules := syntax.Rules[language]{"x": {Attributes: map[string]syntax.Attribute{name: maintained.Attributes[name]}}}
		tree := carrier.AST[language]{ID: "sample", Root: 1, Nodes: []carrier.Node[language]{{Kind: "x", Attributes: []carrier.Attribute{{Name: name, Text: sample}}}}}
		return g.Check(tree, "x", rules, nil)
	}
	accepts := func(name, sample string) bool { return check(name, sample) == nil }
	// A set that is present and empty admits no text. The extractor reports no
	// accepted value, and the interpreter refuses the table that holds the set,
	// as unsupported and for every text.
	var refusal *text.Diagnostic
	if err := check("empty", "a"); !errors.As(err, &refusal) || refusal.Code != text.Unsupported {
		t.Fatalf("the interpreter does not refuse a table with an empty set as unsupported: %v", err)
	}
	for _, model := range feature.Attributes.Value {
		actual, ok := maintained.Attributes[model.Name]
		if !ok || model.Required.Value != actual.Required || model.Values.State != c.Known {
			t.Fatalf("attribute disagrees: %#v", model)
		}
		for _, sample := range []string{"", "a", "b", "c", "true", "false", "after"} {
			represented := model.Values.Value.Unrestricted || slices.Contains(model.Values.Value.Values, sample)
			if accepted := accepts(model.Name, sample); accepted != represented {
				t.Fatalf("%s(%q): interpreter=%v extractor=%v", model.Name, sample, accepted, represented)
			}
		}
	}
	if !slices.Equal(feature.References.Value, []string{"target"}) {
		t.Fatal("reference allowance changed")
	}
}

func TestTableTypedByLanguageReadsVocabularyConstants(t *testing.T) {
	source := &fixtureSource{set: token.NewFileSet(), files: map[c.FileRef]c.File{}, declarations: map[c.Symbol]c.Declaration{}, facts: map[c.Span]c.Facts{}, packages: map[string]*types.Package{}}
	source.add(t, kindPath, kindSource)
	source.add(t, grammarPath, typedGrammarSource)
	source.add(t, vocabularyPath, vocabularySource)
	source.add(t, fixturePath, "package fixture\nimport g \""+grammarPath+"\"\nimport v \""+vocabularyPath+"\"\n"+`
var types = []v.Kind{v.KindAtom, v.KindRecord}
func prefix() []g.Role[v.Language] { return []g.Role[v.Language]{{Name:v.RoleField,Max:-1,Kinds:[]v.Kind{v.KindField}}} }
func with(base []g.Role[v.Language], more ...g.Role[v.Language]) []g.Role[v.Language] { return append(base,more...) }
var rules = g.Rules[v.Language]{
 v.KindRecord: {Attributes:map[string]g.Attribute{v.AttributeName:{Required:true}}, Roles:with(prefix(),g.Role[v.Language]{Name:v.RoleType,Max:1,Kinds:append(types,v.KindField)})},
 v.KindField: {Roles:[]g.Role[v.Language]{{Name:v.RoleType,Min:1,Max:1,Kinds:types}},Reference:true},
 v.KindAtom: {},
}`)
	result, err := (extract.Extractor{}).Extract(context.Background(), c.DefinitionRequest{Root: c.Symbol{Package: fixturePath, Name: "rules"}, Source: source})
	if err != nil {
		t.Fatal(err)
	}
	if result.Extent != c.Decoded || len(result.Features) != 3 || len(result.Limitations) != 0 {
		t.Fatalf("typed table was not decoded: %#v", result)
	}
	record, field, atom := result.Features[0], result.Features[1], result.Features[2]
	if record.Key != "record" || field.Key != "field" || atom.Key != "atom" {
		t.Fatalf("keys are not the vocabulary's kind values: %#v", result.Features)
	}
	if len(record.Attributes.Value) != 1 || record.Attributes.Value[0].Name != "name" || !record.Attributes.Value[0].Required.Value {
		t.Fatalf("attribute name constant: %#v", record.Attributes)
	}
	if len(record.Roles.Value) != 2 {
		t.Fatalf("roles through structural helpers: %#v", record.Roles)
	}
	for i, maintained := range []syntax.Role[language]{{Name: "field", Max: -1, Kinds: []carrier.Kind[language]{"field"}}, {Name: "type", Max: 1, Kinds: []carrier.Kind[language]{"atom", "record", "field"}}} {
		model := record.Roles.Value[i]
		if model.Name != maintained.Name || model.Min.Value != maintained.Min || model.Max.Value.Unbounded != (maintained.Max < 0) || !slices.Equal(model.Kinds.Value, names(maintained.Kinds)) {
			t.Fatalf("role %d through typed helpers: %#v", i, record.Roles)
		}
	}
	one := syntax.Role[language]{Name: "type", Min: 1, Max: 1, Kinds: []carrier.Kind[language]{"atom", "record"}}
	if len(field.Roles.Value) != 1 || field.Roles.Value[0].Name != one.Name || field.Roles.Value[0].Min.Value != one.Min || field.Roles.Value[0].Max.Value.Value != one.Max || !slices.Equal(field.Roles.Value[0].Kinds.Value, names(one.Kinds)) || !slices.Equal(field.References.Value, []string{"target"}) {
		t.Fatalf("typed kind list and reference: %#v", field)
	}
	if atom.Roles.State != c.Known || len(atom.Roles.Value) != 0 || len(atom.Attributes.Value) != 0 || len(atom.References.Value) != 0 {
		t.Fatalf("known absence: %#v", atom)
	}
}

func TestUnresolvedAttributeKeyLeavesTheCollectionPartial(t *testing.T) {
	result := inventory(t, `
func key()string{return "dynamic"}
var rules=g.Rules{"node":{Attributes:map[string]g.Attribute{"same":{OneOf:[]string{"after"}},key():{Required:true}}}}
`)
	attributes := result.Features[0].Attributes
	if result.Extent != c.PartlyDecoded || attributes.State != c.Partial || len(attributes.Value) != 1 {
		t.Fatalf("unknown map portion lost: %#v", result)
	}
	value := attributes.Value[0]
	if value.Name != "same" || value.Required.Value || !slices.Equal(value.Values.Value.Values, []string{"after"}) {
		t.Fatalf("the known entry was not reported: %#v", value)
	}
}

func TestUnknownFieldsPreserveIndependentFacts(t *testing.T) {
	result := inventory(t, `
func dynamic()[]string{values:=[]string{};return values}
func extra()map[string]g.Attribute{m:=map[string]g.Attribute{};return m}
var rules=g.Rules{"node":{Attributes:extra(),Roles:[]g.Role{{Name:"part",Min:2,Max:-1,Kinds:dynamic()}}}}
`)
	feature := result.Features[0]
	if feature.Attributes.State != c.Unknown || feature.Roles.State != c.Known {
		t.Fatalf("independent collections: %#v", feature)
	}
	role := feature.Roles.Value[0]
	if role.Min.State != c.Known || role.Min.Value != 2 || !role.Max.Value.Unbounded || role.Kinds.State == c.Known {
		t.Fatalf("independent role fields: %#v", role)
	}
	// A set that cannot be decoded is not reported as unrestricted.
	result = inventory(t, `
func dynamic()[]string{values:=[]string{};return values}
var rules=g.Rules{"node":{Attributes:map[string]g.Attribute{"v":{OneOf:dynamic()},"w":{Predicate:g.Bool(),OneOf:dynamic()}}}}
`)
	for _, attribute := range result.Features[0].Attributes.Value {
		if attribute.Values.State == c.Known || attribute.Values.Value.Unrestricted {
			t.Fatalf("undecoded set: %#v", attribute)
		}
	}
}

func TestExactNameIdentity(t *testing.T) {
	result := inventory(t, `
func Bool() g.Predicate { return g.Predicate{"fixture","bool"} }
type Predicate struct{ Owner, Name string }
func role()g.Role{r:=g.Role{Name:"fake"};return r}
var flag = g.Bool()
var rules=g.Rules{"node":{Attributes:map[string]g.Attribute{"own":{Predicate:Bool()},"shared":{Predicate:flag},"other":{Predicate:g.Nonempty()},"literal":{Predicate:g.Predicate{"bittype","bool"}}},Roles:[]g.Role{role()}}}
`)
	f := result.Features[0]
	if f.Roles.State != c.Partial || len(f.Roles.Value) != 0 {
		t.Fatalf("a role that is no literal was decoded: %#v", f.Roles)
	}
	attributes := map[string]c.Attribute{}
	for _, attribute := range f.Attributes.Value {
		attributes[attribute.Name] = attribute
	}
	if own := attributes["own"].Values; own.State != c.Opaque || own.Value.Unrestricted || len(own.Value.Values) != 0 {
		t.Fatalf("an unrelated Bool was mistaken for the shared name: %#v", own)
	}
	if shared := attributes["shared"].Values; shared.State != c.Known || !slices.Equal(shared.Value.Values, []string{"true", "false"}) {
		t.Fatalf("the shared name through a variable: %#v", shared)
	}
	// Only the declared name function is interpreted. A literal that spells the
	// same owner and name, and every other name, stays opaque.
	for _, name := range []string{"other", "literal"} {
		if values := attributes[name].Values; values.State != c.Opaque || values.Value.Unrestricted {
			t.Fatalf("%s: %#v", name, values)
		}
	}
}

func TestEmptySetAdmitsNothingAndUnknownKeyIsPartial(t *testing.T) {
	result := inventory(t, `
func key()string{return "dynamic"}
var none = []string{}
var absent []string = nil
var rules=g.Rules{"known":{Attributes:map[string]g.Attribute{"none":{OneOf:none},"literal":{OneOf:[]string{}},"all":{},"nil":{OneOf:nil},"named":{OneOf:absent}}},key():{}}
`)
	if result.Extent != c.PartlyDecoded || len(result.Limitations) != 1 || len(result.Features) != 1 {
		t.Fatalf("unknown key: %#v", result)
	}
	// The table type says of a set: absent, which is nil, it restricts nothing;
	// present and empty it admits no text.
	for _, attribute := range result.Features[0].Attributes.Value {
		unrestricted := attribute.Name == "all" || attribute.Name == "nil" || attribute.Name == "named"
		if attribute.Values.State != c.Known || attribute.Values.Value.Unrestricted != unrestricted || len(attribute.Values.Value.Values) != 0 {
			t.Fatalf("%s: %#v", attribute.Name, attribute)
		}
	}
	if len(result.Features[0].Attributes.Value) != 5 {
		t.Fatalf("attributes: %#v", result.Features[0].Attributes)
	}
}

func TestRepeatedStructuralHelperHasDistinctInvocationBindings(t *testing.T) {
	result := inventory(t, `
func with(base []g.Role, more ...g.Role)[]g.Role{return append(base,more...)}
var rules=g.Rules{"x":{Roles:with(with([]g.Role{},g.Role{Name:"first",Min:1,Max:1,Kinds:[]string{"x"}}),g.Role{Name:"second",Max:1,Kinds:[]string{"x"}})}}
`)
	roles := result.Features[0].Roles.Value
	if result.Extent != c.Decoded || len(roles) != 2 || roles[0].Name != "first" || roles[1].Name != "second" {
		t.Fatalf("lawful nested helper composition was lost: %#v", result)
	}
	missing := inventory(t, `var rules g.Rules`)
	if missing.Extent != c.PartlyDecoded || len(missing.Limitations) != 1 || len(missing.Features) != 0 {
		t.Fatalf("missing initializer was not localized: %#v", missing)
	}
}

func TestDefinedTypeIsNotRulesAndLimitsAreFailures(t *testing.T) {
	_, err := (extract.Extractor{}).Extract(context.Background(), request(t, `type Distinct g.Rules;var rules=Distinct{"x":{}}`))
	if err == nil {
		t.Fatal("new defined type treated as an alias")
	}
	r := request(t, `func loop()[]g.Role{return loop()};var rules=g.Rules{"x":{Roles:loop()}}`)
	if _, err := (extract.Extractor{}).Extract(context.Background(), r); err == nil || !strings.Contains(err.Error(), "cycle") {
		t.Fatalf("cycle: %v", err)
	}
	r = request(t, `var rules=g.Rules{"x":{Roles:[]g.Role{{Name:"part",Min:1,Max:1,Kinds:[]string{"x"}}}}}`)
	r.Limits.MaxExpansions = 1
	if _, err := (extract.Extractor{}).Extract(context.Background(), r); err == nil {
		t.Fatal("traversal exhaustion reported as success")
	}
	r.Limits = c.Limits{MaxRecords: 1}
	if _, err := (extract.Extractor{}).Extract(context.Background(), r); err == nil {
		t.Fatal("record exhaustion reported as success")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := (extract.Extractor{}).Extract(ctx, r); err != context.Canceled {
		t.Fatalf("cancellation: %v", err)
	}
}
