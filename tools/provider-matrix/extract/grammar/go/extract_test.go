package grammar_test

import (
	"context"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"go/types"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/internal/grammar"
	c "github.com/Bitspark/bittype/tools/provider-matrix/capability/go"
	extract "github.com/Bitspark/bittype/tools/provider-matrix/extract/grammar/go"
)

const grammarPath = "github.com/Bitspark/bittype/internal/grammar"
const fixturePath = "example.test/fixture"

// These declarations provide compiler identities/signatures; they are never
// executed. Maintained helper meaning is separately checked against g below.
const grammarSource = `package grammar
type Attribute struct { Required bool; Accept func(string)bool }
type Role struct { Name string; Min,Max int; Kinds []string }
type Rule struct { Attributes map[string]Attribute; Roles []Role; Reference bool }
type Rules map[string]Rule
func Attr(bool,func(string)bool)Attribute{return Attribute{}}
func Attributes(...any)map[string]Attribute{return nil}
func One(string,...string)Role{return Role{}}
func Opt(string,...string)Role{return Role{}}
func Many(string,int,...string)Role{return Role{}}
func Set(...string)func(string)bool{return nil}
func Bool(string)bool{return false}
`

type fixtureSource struct {
	set *token.FileSet
	files map[c.FileRef]c.File
	declarations map[c.Symbol]c.Declaration
	facts map[c.Span]c.Facts
	packages map[string]*types.Package
}

func (s *fixtureSource) File(ref c.FileRef) (c.File, error) {
	if file, ok := s.files[ref]; ok { return file, nil }
	return c.File{}, fmt.Errorf("foreign file")
}
func (s *fixtureSource) Declaration(symbol c.Symbol) (c.Declaration, error) {
	if declaration, ok := s.declarations[symbol]; ok { return declaration, nil }
	return c.Declaration{}, fmt.Errorf("unavailable declaration %s", symbol)
}
func (s *fixtureSource) Facts(span c.Span) (c.Facts, error) {
	if facts, ok := s.facts[span]; ok { return facts, nil }
	return c.Facts{}, fmt.Errorf("not an expression: %#v", span)
}
func (s *fixtureSource) Import(path string) (*types.Package, error) {
	if pkg := s.packages[path]; pkg != nil { return pkg, nil }
	return nil, fmt.Errorf("missing import %s", path)
}

func symbol(obj types.Object) *c.Symbol {
	if obj == nil { return nil }
	if obj.Pkg() == nil {
		if obj.Parent() == types.Universe { return &c.Symbol{Package: "builtin", Name: obj.Name()} }
		return nil
	}
	if obj.Parent() != obj.Pkg().Scope() { return nil }
	return &c.Symbol{Package: obj.Pkg().Path(), Name: obj.Name()}
}

func (s *fixtureSource) add(t *testing.T, path, text string) {
	t.Helper()
	file, err := parser.ParseFile(s.set, path+"/source.go", text, 0)
	if err != nil { t.Fatal(err) }
	info := &types.Info{Types: map[ast.Expr]types.TypeAndValue{}, Defs: map[*ast.Ident]types.Object{}, Uses: map[*ast.Ident]types.Object{}, Selections: map[*ast.SelectorExpr]*types.Selection{}}
	config := types.Config{Importer: s}
	pkg, err := config.Check(path, s.set, []*ast.File{file}, info)
	if err != nil { t.Fatal(err) }
	s.packages[path] = pkg
	ref := c.FileRef(path)
	s.files[ref] = c.File{Ref: ref, Path: path+"/source.go", Text: text}
	tokenFile := s.set.File(file.Pos())
	span := func(n ast.Node) c.Span {
		if n == nil { return c.Span{} }
		return c.Span{File: ref, Start: tokenFile.Offset(n.Pos()), End: tokenFile.Offset(n.End())}
	}
	var object func(ast.Expr) types.Object
	object = func(expr ast.Expr) types.Object {
		switch n := expr.(type) {
		case *ast.Ident: return info.Uses[n]
		case *ast.SelectorExpr: return info.Uses[n.Sel]
		case *ast.ParenExpr: return object(n.X)
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
						if len(n.Values) == len(n.Names) { value = span(n.Values[i]) }
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
		if !ok { return true }
		facts := c.Facts{Symbol: symbol(object(expr))}
		if call, ok := expr.(*ast.CallExpr); ok { facts.Callee = symbol(object(call.Fun)) }
		if typ := info.TypeOf(expr); typ != nil {
			facts.Type.Text = types.TypeString(typ, nil)
			if named, ok := types.Unalias(typ).(*types.Named); ok { facts.Type.Origin = symbol(named.Origin().Obj()) }
		}
		if value := info.Types[expr].Value; value != nil {
			text := value.ExactString()
			if value.Kind() == constant.String { text = constant.StringVal(value) }
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
	if err != nil { t.Fatal(err) }
	return result
}

func TestAliasesCompositionAndLocalKnowledge(t *testing.T) {
	result := inventory(t, `
type Alias = g.Rules
const atom = "atom"
var kinds = []string{atom}
func prefix() []g.Role { return []g.Role{g.Many("items",0,kinds...)} }
func with(base []g.Role, more ...g.Role) []g.Role { return append(base,more...) }
func custom(string)bool{return true}
func named()map[string]g.Attribute{return g.Attributes("mode",g.Attr(false,g.Set("one","two")),"free",g.Attr(false,nil),"custom",g.Attr(true,custom),"flag",g.Attr(false,g.Bool))}
var rules = Alias{
 "record": {Attributes:named(), Roles:with(prefix(),g.Opt("next",append([]string{},kinds...)...)),Reference:true},
 "empty": {},
}`)
	if result.Extent != c.PartlyDecoded || len(result.Features) != 2 { t.Fatalf("extent/features: %#v", result) }
	record := result.Features[0]
	if record.Key != "record" || record.Roles.State != c.Known || len(record.Roles.Value) != 2 { t.Fatalf("roles: %#v", record) }
	items := record.Roles.Value[0]
	maintained := g.Many("items",0,"atom")
	if items.Name != maintained.Name || items.Min.Value != maintained.Min || !items.Max.Value.Unbounded || maintained.Max >= 0 || strings.Join(items.Kinds.Value,",") != strings.Join(maintained.Kinds,",") { t.Fatalf("owning helper semantics: %#v", items) }
	if record.Roles.Value[1].Max.Value.Value != 1 || record.References.Value[0] != "target" { t.Fatalf("optional/reference: %#v", record) }
	attrs := record.Attributes.Value
	if record.Attributes.State != c.Known || len(attrs) != 4 { t.Fatalf("attributes: %#v", record.Attributes) }
	if attrs[0].Required.State != c.Known || attrs[0].Required.Value || strings.Join(attrs[0].Values.Value.Values,",") != "one,two" { t.Fatalf("finite constraint: %#v", attrs[0]) }
	if !attrs[1].Values.Value.Unrestricted || attrs[2].Values.State != c.Opaque || !attrs[2].Required.Value { t.Fatalf("local knowledge: %#v", attrs) }
	if strings.Join(attrs[3].Values.Value.Values,",") != "true,false" || !g.Bool("true") || !g.Bool("false") || g.Bool("other") { t.Fatalf("Boolean predicate: %#v", attrs[3]) }
	empty := result.Features[1]
	if empty.Roles.State != c.Known || len(empty.Roles.Value) != 0 || len(empty.Attributes.Value) != 0 || len(empty.References.Value) != 0 { t.Fatalf("known absence: %#v", empty) }
}

func TestInterpretationMatchesMaintainedOwningHelpers(t *testing.T) {
	result := inventory(t, `var rules=g.Rules{"x":{
 Roles:[]g.Role{g.One("one","a"),g.Opt("optional","a","b"),g.Many("many",2,"b")},
 Attributes:g.Attributes("finite",g.Attr(true,g.Set("a","b")),"empty",g.Attr(false,g.Set()),"free",g.Attr(false,nil)),Reference:true,
}}`)
	maintained := g.Rule{
		Roles: []g.Role{g.One("one","a"),g.Opt("optional","a","b"),g.Many("many",2,"b")},
		Attributes: g.Attributes("finite",g.Attr(true,g.Set("a","b")),"empty",g.Attr(false,g.Set()),"free",g.Attr(false,nil)),
		Reference: true,
	}
	feature := result.Features[0]
	if result.Extent != c.Decoded || len(feature.Roles.Value) != len(maintained.Roles) { t.Fatalf("unexpected interpretation: %#v",result) }
	for i,actual := range maintained.Roles {
		model := feature.Roles.Value[i]
		if model.Name != actual.Name || model.Min.Value != actual.Min || model.Max.Value.Unbounded != (actual.Max<0) || actual.Max>=0 && model.Max.Value.Value!=actual.Max || !slices.Equal(model.Kinds.Value,actual.Kinds) {
			t.Fatalf("owning role %d disagrees: %#v versus %#v",i,model,actual)
		}
	}
	for _,model := range feature.Attributes.Value {
		actual,ok := maintained.Attributes[model.Name]
		if !ok || model.Required.Value != actual.Required || model.Values.State != c.Known { t.Fatalf("owning attribute disagrees: %#v",model) }
		for _,sample := range []string{"","a","b","c"} {
			accepted := actual.Accept==nil || actual.Accept(sample)
			represented := model.Values.Value.Unrestricted || slices.Contains(model.Values.Value.Values,sample)
			if accepted != represented { t.Fatalf("%s(%q): maintained=%v interpreted=%v",model.Name,sample,accepted,represented) }
		}
	}
	if !slices.Equal(feature.References.Value,[]string{"target"}) { t.Fatal("reference allowance changed") }
}

func TestUnknownFieldsPreserveIndependentFacts(t *testing.T) {
	result := inventory(t, `
func dynamic()[]string{values:=[]string{};return values}
func extra()map[string]g.Attribute{m:=map[string]g.Attribute{};return m}
var rules=g.Rules{"node":{Attributes:extra(),Roles:[]g.Role{g.Many("part",2,dynamic()...)}}}
`)
	feature := result.Features[0]
	if feature.Attributes.State != c.Unknown || feature.Roles.State != c.Known { t.Fatalf("independent collections: %#v", feature) }
	role := feature.Roles.Value[0]
	if role.Min.State != c.Known || role.Min.Value != 2 || !role.Max.Value.Unbounded || role.Kinds.State == c.Known { t.Fatalf("independent role fields: %#v", role) }
}

func TestExactHelperIdentityAndShadowedNil(t *testing.T) {
	result := inventory(t, `
func One(string,...string)g.Role{r:=g.Role{};return r}
var nil=func(string)bool{return true}
var rules=g.Rules{"node":{Attributes:g.Attributes("v",g.Attr(false,nil)),Roles:[]g.Role{One("fake","x")}}}
`)
	f := result.Features[0]
	if f.Roles.State != c.Partial || len(f.Roles.Value) != 0 { t.Fatalf("unrelated One was mistaken for owning helper: %#v", f.Roles) }
	if f.Attributes.Value[0].Values.State != c.Opaque || f.Attributes.Value[0].Values.Value.Unrestricted { t.Fatalf("shadowed nil was treated as absent predicate: %#v", f.Attributes) }
}

func TestFiniteEmptyIsNotUnrestrictedAndUnknownKeyIsPartial(t *testing.T) {
	result := inventory(t, `
func key()string{return "dynamic"}
var rules=g.Rules{"known":{Attributes:g.Attributes("none",g.Attr(false,g.Set()),"all",g.Attr(false,nil))},key():{}}
`)
	if result.Extent != c.PartlyDecoded || len(result.Limitations) != 1 || len(result.Features) != 1 { t.Fatalf("unknown key: %#v", result) }
	attrs := result.Features[0].Attributes.Value
	if attrs[0].Values.State != c.Known || attrs[0].Values.Value.Unrestricted || len(attrs[0].Values.Value.Values) != 0 || !attrs[1].Values.Value.Unrestricted { t.Fatalf("empty versus unrestricted: %#v", attrs) }
}

func TestRepeatedStructuralHelperHasDistinctInvocationBindings(t *testing.T) {
	result := inventory(t, `
func with(base []g.Role, more ...g.Role)[]g.Role{return append(base,more...)}
var rules=g.Rules{"x":{Roles:with(with([]g.Role{},g.One("first","x")),g.Opt("second","x"))}}
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
	if err == nil { t.Fatal("new defined type treated as an alias") }
	r := request(t, `func loop()[]g.Role{return loop()};var rules=g.Rules{"x":{Roles:loop()}}`)
	if _, err := (extract.Extractor{}).Extract(context.Background(), r); err == nil || !strings.Contains(err.Error(),"cycle") { t.Fatalf("cycle: %v",err) }
	r = request(t, `var rules=g.Rules{"x":{Roles:[]g.Role{g.One("part","x")}}}`)
	r.Limits.MaxExpansions = 1
	if _, err := (extract.Extractor{}).Extract(context.Background(), r); err == nil { t.Fatal("traversal exhaustion reported as success") }
	r.Limits = c.Limits{MaxRecords:1}
	if _, err := (extract.Extractor{}).Extract(context.Background(), r); err == nil { t.Fatal("record exhaustion reported as success") }
	ctx,cancel := context.WithCancel(context.Background());cancel()
	if _, err := (extract.Extractor{}).Extract(ctx,r); err != context.Canceled { t.Fatalf("cancellation: %v",err) }
}
