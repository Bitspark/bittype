package integration_test

import (
	"context"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	checkhoros "github.com/Bitspark/bittype/check/syntax/horos/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	"github.com/Bitspark/bittype/engine/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	"github.com/Bitspark/bittype/semantics/horos/go"
	cs "github.com/Bitspark/bittype/syntax/cpp/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	rs "github.com/Bitspark/bittype/syntax/rust/go"
	ts "github.com/Bitspark/bittype/syntax/typescript/go"
	cgtarget "github.com/Bitspark/bittype/transform/syntax/cpp/target/go"
	ggtarget "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	rgtarget "github.com/Bitspark/bittype/transform/syntax/rust/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	tgtarget "github.com/Bitspark/bittype/transform/syntax/typescript/target/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

// Canonical fixtures are independent of the native generators and text parsers.
func genAtom(s string) core.Value { return core.NewAtom([]byte(s)) }
func genTuple(tag string, args ...core.Value) core.Value {
	return core.NewTuple(append([]core.Value{genAtom(tag)}, args...)...)
}
func genInt(i int64) core.Value { return data.EncodeInt(big.NewInt(i)) }
func genMap(entries ...[2]core.Value) core.Value {
	v, err := data.EncodeMap(entries)
	if err != nil {
		panic(err)
	}
	return v
}
func genEntry(name string, value core.Value) [2]core.Value {
	return [2]core.Value{data.MustEncodeText(name), value}
}
func genModule(name string, imports core.Value, entries ...[2]core.Value) core.Value {
	return genTuple("horos/decl/module", genInt(4), data.MustEncodeText(name), imports, genMap(entries...))
}
func genTop() core.Value { return genTuple("horos/decl/top", genTuple("horos/decl/sort-proper")) }
func genType(body core.Value, params ...string) core.Value {
	var p []core.Value
	for _, name := range params {
		p = append(p, genTuple("horos/decl/parameter", data.MustEncodeText(name), genTop()))
	}
	return genTuple("horos/decl/type", data.EncodeList(p), genTop(), body)
}
func genParam(i int64) core.Value { return genTuple("horos/decl/param", genInt(i)) }
func genLocal(name string, args ...core.Value) core.Value {
	if len(args) == 0 {
		return genTuple("horos/decl/local", data.MustEncodeText(name))
	}
	return genTuple("horos/decl/apply", genTuple("horos/decl/local", data.MustEncodeText(name)), data.EncodeList(args))
}
func genImported(alias, name string, args ...core.Value) core.Value {
	if len(args) == 0 {
		return genTuple("horos/decl/imported", data.MustEncodeText(alias), data.MustEncodeText(name))
	}
	return genTuple("horos/decl/apply", genTuple("horos/decl/imported", data.MustEncodeText(alias), data.MustEncodeText(name)), data.EncodeList(args))
}
func genScalar(name string) core.Value { return genTuple("horos/scalar", genAtom(name)) }
func genField(name string, optional bool, ty core.Value) core.Value {
	flag := "req"
	if optional {
		flag = "opt"
	}
	return genTuple("horos/field", data.MustEncodeText(name), genAtom(flag), ty)
}
func genRecord(fields ...core.Value) core.Value {
	return genTuple("horos/record", append([]core.Value{genAtom("closed")}, fields...)...)
}
func genPage() core.Value {
	return genRecord(genField("items", false, genTuple("horos/list-of", genParam(0))), genField("next", true, genScalar("utf8-text")))
}
func genPerson() core.Value {
	return genRecord(genField("name", false, genScalar("utf8-text")), genField("active", false, genScalar("bool")))
}
func paginationModule() core.Value {
	return genModule("example.pagination", genMap(genEntry("h", decl.PreludeImport())),
		genEntry("Page", genType(genPage(), "T")), genEntry("Person", genType(genPerson())),
		genEntry("PeoplePage", genType(genLocal("Page", genLocal("Person")))),
		genEntry("NestedPage", genType(genLocal("Page", genLocal("Page", genLocal("Person"))))),
		genEntry("Other", genType(genTuple("horos/list-of", genParam(0)), "T")),
		genEntry("DataPage", genTuple("horos/decl/type", data.EncodeList([]core.Value{genTuple("horos/decl/parameter", data.MustEncodeText("D"), genTuple("horos/decl/kind-imported", data.MustEncodeText("h"), data.MustEncodeText("Data")))}), genTop(), genLocal("Page", genParam(0)))),
	)
}

func generationInput(t *testing.T, module core.Value, deps ...core.Value) (*engine.Engine, representation.Language, transform.Input) {
	t.Helper()
	source, err := resolve.Resolve(module, deps, "canonical")
	if err != nil {
		t.Fatal(err)
	}
	return generationSource(t, source)
}
func generationSource(t *testing.T, source h.AST) (*engine.Engine, representation.Language, transform.Input) {
	t.Helper()
	ctx := context.Background()
	from, err := representation.NewLanguage[h.Language](h.Language{}, h.Profile)
	if err != nil {
		t.Fatal(err)
	}
	envelope, err := representation.Pack(from, source)
	if err != nil {
		t.Fatal(err)
	}
	e := engine.New()
	if err = e.RegisterLanguage(from); err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[h.Language](from, checkhoros.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	resolver, err := resolution.NewResolver[h.Language](from, horos.ResolutionContract(), resolve.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[h.Language](from, horos.AnalysisContract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	snapshot, err := e.Capture(ctx, semantic.Input{Entry: ast.NodeRef{Tree: source.ID, Node: source.Root}, Trees: []representation.AST{envelope}})
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := e.Resolve(ctx, snapshot, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	analyzed, err := e.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved})
	if err != nil {
		t.Fatalf("analysis: %v, %+v", err, analyzed.Diagnostics())
	}
	return e, from, transform.Input{Source: snapshot, Evidence: []semantic.Result{analyzed}}
}

func generatedOutput[L any](t *testing.T, e *engine.Engine, input transform.Input, p nativeProvider[L], configure requestBuilder) g.Output {
	t.Helper()
	to, err := representation.NewLanguage[L](p, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterLanguage(to); err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[L](to, p)
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterChecker(checker); err != nil {
		t.Fatal(err)
	}
	request, err := configure(e.Generation(), input, to)
	if err != nil {
		t.Fatal(err)
	}
	out, err := e.Generate(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func generated[L any](t *testing.T, e *engine.Engine, input transform.Input, p nativeProvider[L], configure requestBuilder) ast.AST[L] {
	t.Helper()
	out := generatedOutput(t, e, input, p, configure)
	return outputTree[L](t, out)
}
func outputTree[L any](t *testing.T, out g.Output) ast.AST[L] {
	t.Helper()
	entry, err := out.Snapshot().EntryTree()
	if err != nil {
		t.Fatal(err)
	}
	tree, err := representation.Unpack[L](out.Snapshot().Language(), entry)
	if err != nil {
		t.Fatal(err)
	}
	return tree
}

type requestBuilder func(*c.Registry, transform.Input, representation.Language) (c.Request, error)

func tsRequest(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
	return genfixture.Compose(reg, input, to, "generated-ts", tgtarget.Provider{}, typescriptBundles()...)
}
func goRequest(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
	return genfixture.Compose(reg, input, to, "generated-go", ggtarget.Provider{Package: "pagination"}, golangBundles()...)
}
func rustRequest(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
	return genfixture.Compose(reg, input, to, "generated-rust", rgtarget.Provider{}, rustBundles()...)
}
func cppRequest(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
	return genfixture.Compose(reg, input, to, "generated-cpp", cgtarget.Provider{Namespace: "pagination"}, cppBundles()...)
}

func nativeChecker[L any](t *testing.T, p nativeProvider[L]) check.Checker {
	t.Helper()
	checker, err := check.NewChecker[L](mustLanguage[L](t, p), p)
	if err != nil {
		t.Fatal(err)
	}
	return checker
}
func runGeneration(t *testing.T, input transform.Input, checker check.Checker, configure requestBuilder) (g.Output, error) {
	t.Helper()
	var registry c.Registry
	request, err := configure(&registry, input, checker.Language())
	if err != nil {
		t.Fatal(err)
	}
	return registry.Generate(context.Background(), request, checker)
}
func mustLanguage[L any](t *testing.T, p language.Provider) representation.Language {
	t.Helper()
	d, err := representation.NewLanguage[L](p, "native-types-v1")
	if err != nil {
		t.Fatal(err)
	}
	return d
}

type paginationOutput struct {
	ts     ts.AST
	goTree gs.AST
	rust   rs.AST
	cpp    cs.AST
}

func paginationOutputs(t *testing.T, module core.Value, deps ...core.Value) paginationOutput {
	e, _, input := generationInput(t, module, deps...)
	model, _ := input.Evidence[0].Semantics()
	before := model.Data()
	out := paginationOutput{placedTS(t, generatedOutput(t, e, input, typescriptText(), tsRequest)), generated(t, e, input, golangText(), goRequest), generated(t, e, input, rustText(), rustRequest), generated(t, e, input, cppText(), cppRequest)}
	after, _ := input.Evidence[0].Semantics()
	if model != after || !reflect.DeepEqual(before, after.Data()) {
		t.Fatal("mapping changed established semantic evidence")
	}
	return out
}
func generationNodeAttribute[L any](node ast.Node[L], name string) string {
	for _, a := range node.Attributes {
		if a.Name == name {
			return a.Text
		}
	}
	return ""
}
func checkGenerationLinks[L any](t *testing.T, tree ast.AST[L]) {
	t.Helper()
	uses := map[ast.NodeID]int{}
	parameters := map[ast.NodeID]bool{}
	for i, node := range tree.Nodes {
		if string(node.Kind) == "type-parameter" && generationNodeAttribute(node, "name") == "T" && len(node.Origin.Nodes) > 0 {
			parameters[ast.NodeID(i+1)] = true
		}
		for _, ref := range node.References {
			if ref.Role == "target" {
				if ref.Target.Tree != tree.ID || ref.Target.Node == 0 || int(ref.Target.Node) > len(tree.Nodes) {
					t.Fatal("dangling native target")
				}
				uses[ref.Target.Node]++
			}
		}
	}
	if len(parameters) != 2 {
		t.Fatalf("lost distinct Page/Other T binders: %v", parameters)
	}
	for p := range parameters {
		if uses[p] == 0 {
			t.Fatal("parameter references rebound by spelling")
		}
	}
	for _, node := range tree.Nodes {
		for _, origin := range node.Origin.Nodes {
			if origin.Tree != "canonical" || origin.Node == 0 {
				t.Fatal("lost canonical provenance")
			}
		}
	}
}
func TestSemanticPaginationMappings(t *testing.T) {
	out := paginationOutputs(t, paginationModule())
	checkGenerationLinks(t, out.ts)
	checkGenerationLinks(t, out.goTree)
	checkGenerationLinks(t, out.rust)
	checkGenerationLinks(t, out.cpp)
	for name, fixture := range map[string]struct {
		source   string
		required []string
	}{
		"typescript": {formatted(t, typescriptText(), out.ts), []string{"export interface Page<T>", "items: T[]", "next?: string", "Page<Page<Person>>", "type DataPage<D> = Page<D>"}},
		"go":         {formatted(t, golangText(), out.goTree), []string{"type Page[T interface {}] struct", "Items []T", "Next BittypeOptional[string]", "Page[Page[Person]]", "type DataPage[D interface {}] = Page[D]"}},
		"rust":       {formatted(t, rustText(), out.rust), []string{"pub struct Page<T>", "::std::vec::Vec<T>", "::std::option::Option<::std::string::String>", "Page<Page<Person>>"}},
		"cpp":        {formatted(t, cppText(), out.cpp), []string{"template<typename T>", "struct Page", "::std::vector<T>", "::std::optional<::std::string>", "::pagination::Page<::pagination::Page<::pagination::Person>>"}},
	} {
		t.Run(name, func(t *testing.T) {
			for _, want := range fixture.required {
				if !strings.Contains(fixture.source, want) {
					t.Fatalf("missing %q:\n%s", want, fixture.source)
				}
			}
			t.Log(fixture.source)
		})
	}
}
func TestPaginationImportedIdentity(t *testing.T) {
	first := genModule("first", genMap(), genEntry("Page", genType(genPage(), "T")))
	second := genModule("second", genMap(), genEntry("Page", genType(genTuple("horos/list-of", genParam(0)), "T")))
	imp := func(name string, module core.Value) core.Value {
		pin, err := decl.Pin(module)
		if err != nil {
			t.Fatal(err)
		}
		return genTuple("horos/decl/import", data.MustEncodeText(name), core.NewAtom(pin[:]))
	}
	root := genModule("root", genMap(genEntry("one", imp("first", first)), genEntry("two", imp("second", second))), genEntry("Person", genType(genPerson())), genEntry("First", genType(genImported("one", "Page", genLocal("Person")))), genEntry("Second", genType(genImported("two", "Page", genLocal("Person")))))
	out := paginationOutputs(t, root, second, first)
	// The same-spelled declarations produce two distinct targets and deterministic output.
	text := formatted(t, typescriptText(), out.ts)
	if strings.Count(text, "_Page<T>") != 2 {
		t.Fatal(text)
	}
	again := paginationOutputs(t, root, first, second)
	if text != formatted(t, typescriptText(), again.ts) {
		t.Fatal("dependency order changed generated names")
	}
}
func TestPaginationRejectsUnsupportedTypes(t *testing.T) {
	cases := map[string]core.Value{
		"wire": genTuple("horos/wire", genScalar("utf8-text")), "integer": genScalar("int"),
		"exact":         genTuple("horos/exact", data.MustEncodeText("only")),
		"open":          genTuple("horos/record", genAtom("open"), genField("name", false, genScalar("utf8-text"))),
		"inline-record": genRecord(genField("inner", false, genPerson())),
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			_, _, input := generationInput(t, genModule("unsupported", genMap(), genEntry("Bad", genType(body))))
			for _, tc := range []struct {
				checker check.Checker
				factory requestBuilder
			}{
				{nativeChecker(t, typescriptText()), tsRequest}, {nativeChecker(t, golangText()), goRequest},
				{nativeChecker(t, rustText()), rustRequest}, {nativeChecker(t, cppText()), cppRequest},
			} {
				output, err := runGeneration(t, input, tc.checker, tc.factory)
				var diagnostic *text.Diagnostic
				if !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || diagnostic.Node == nil || diagnostic.Node.Node == 0 || output.Snapshot().Valid() {
					t.Fatalf("unsupported output: %+v %v", output, err)
				}
			}
		})
	}
}

func TestGenerationRejectsInvalidCanonicalApplications(t *testing.T) {
	for name, source := range map[string]string{
		"data-bound": "module invalid; import h from horos.kinds pin \"sha256:" + decl.PreludePin().String() + "\"; type Page<T> = {items: list<T>;}; type DataPage<D: h.Data> = Page<D>; type Bad = DataPage<wire<text>>;",
		"arity":      "module invalid; type Page<T> = {items: list<T>;}; type Bad = Page<text, bool>;",
		"unbound":    "module invalid; type Bad = Missing<text>;",
	} {
		t.Run(name, func(t *testing.T) {
			parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "invalid.horos", TreeID: "canonical", Text: source})
			if err != nil {
				t.Fatal(err)
			}
			_, _, input := generationSource(t, parsed.Abstract)
			if input.Evidence[0].HasSemantics() {
				t.Fatal("invalid canonical application published a model")
			}
			output, err := runGeneration(t, input, nativeChecker(t, typescriptText()), tsRequest)
			if !errors.Is(err, semantic.ErrEvidenceUnavailable) || output.Snapshot().Valid() {
				t.Fatalf("invalid semantic input mapped: %v", err)
			}
		})
	}
}

func TestGenerationNameCollisionsAndPhantomSlots(t *testing.T) {
	// Two declarations of one name, and a type parameter named like a
	// declaration, also like its own, are no rule of generation for TypeScript: its Files stage
	// decides them per file. Go, Rust and C++ still state one flat scope. A
	// declaration with two parameters of one name, and a parameter that no
	// representation uses, are refused by every target.
	targets := []struct {
		name    string
		checker check.Checker
		request requestBuilder
		// perFile says that names across declarations are decided by Files.
		perFile bool
	}{
		{"typescript", nativeChecker(t, typescriptText()), tsRequest, true},
		{"go", nativeChecker(t, golangText()), goRequest, false},
		{"rust", nativeChecker(t, rustText()), rustRequest, false},
		{"cpp", nativeChecker(t, cppText()), cppRequest, false},
	}
	for _, c := range []struct {
		name, source string
		// perFile is the number of declarations that a target generates which
		// leaves names across declarations to its Files stage. 0 says that every
		// target refuses the source.
		perFile int
	}{
		{"declarations", "module collision; type page = text; type Page = bool;", 2},
		{"binders", "module collision; type T = text; type Page<T> = list<T>;", 2},
		{"own", "module collision; type box<box> = list<box>;", 1},
		{"parameters", "module collision; type Pair<t, T> = { first: t; second: T; };", 0},
		{"unused", "module collision; type Phantom<T> = text;", 0},
	} {
		for _, target := range targets {
			t.Run(c.name+"/"+target.name, func(t *testing.T) {
				parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "names.horos", TreeID: "canonical", Text: c.source})
				if err != nil {
					t.Fatal(err)
				}
				_, _, input := generationSource(t, parsed.Abstract)
				if !input.Evidence[0].HasSemantics() {
					t.Fatal("fixture must be semantically valid")
				}
				output, err := runGeneration(t, input, target.checker, target.request)
				if c.perFile != 0 && target.perFile {
					if err != nil || len(output.Inventory()) != c.perFile {
						t.Fatalf("names across declarations are decided per file, and generation refused them: %v", err)
					}
					return
				}
				var diagnostic *text.Diagnostic
				if !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || output.Snapshot().Valid() {
					t.Fatalf("unsupported names mapped: %v", err)
				}
			})
		}
	}
	_, _, input := generationInput(t, genModule("fields", genMap(), genEntry("Record", genType(genRecord(genField("name", false, genScalar("utf8-text")), genField("Name", false, genScalar("bool")))))))
	if output, err := runGeneration(t, input, nativeChecker(t, golangText()), goRequest); err == nil || output.Snapshot().Valid() {
		t.Fatal("Go field capitalization collision accepted")
	}
}

func TestGenerationTraversalLimit(t *testing.T) {
	nested := genScalar("utf8-text")
	for i := 0; i < 140; i++ {
		nested = genTuple("horos/list-of", nested)
	}
	_, _, input := generationInput(t, genModule("deep", genMap(), genEntry("Deep", genType(nested))))
	output, err := runGeneration(t, input, nativeChecker(t, typescriptText()), tsRequest)
	var diagnostic *text.Diagnostic
	if !errors.As(err, &diagnostic) || diagnostic.Code != text.LimitExceeded || output.Snapshot().Valid() {
		t.Fatalf("deep mapping escaped limit: %v", err)
	}
}

func TestGeneratedPaginationCompilers(t *testing.T) {
	if os.Getenv("BITTYPE_NATIVE") != "1" {
		t.Skip("native generation qualification runs in the maintained hosted job")
	}
	out := paginationOutputs(t, paginationModule())
	write := func(dir, name, source string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(dir, name), []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
	}
	t.Run("typescript", func(t *testing.T) {
		dir := t.TempDir()
		write(dir, "types.ts", formatted(t, typescriptText(), out.ts))
		write(dir, "consumer.ts", `import type {PeoplePage, NestedPage, DataPage, Person} from './types';
const page: PeoplePage = {items: [{name: 'Ada', active: true}]};
const nested: NestedPage = {items: [page]};
const data: DataPage<Person> = page;
// @ts-expect-error exact presence does not admit present undefined
const bad: PeoplePage = {items: [], next: undefined};
// @ts-expect-error argument substitution retains Person
const wrong: PeoplePage = {items: ['Ada']};
`)
		command(t, dir, "tsc", "--noEmit", "--strict", "--exactOptionalPropertyTypes", "--target", "ES2020", "--module", "ESNext", "types.ts", "consumer.ts")
	})
	t.Run("go", func(t *testing.T) {
		dir := t.TempDir()
		write(dir, "types.go", formatted(t, golangText(), out.goTree))
		write(dir, "go.mod", "module pagination\n\ngo 1.27.0\n")
		write(dir, "consumer_test.go", `package pagination
import "testing"
func TestGenerated(t *testing.T) {
 p := PeoplePage{Items: []Person{{Name: "Ada", Active: true}}, Next: BittypeOptional[string]{Present: true, Value: ""}}
 var data DataPage[Person] = p
 nested := NestedPage{Items: []Page[Person]{data}}
 if !nested.Items[0].Next.Present || nested.Items[0].Next.Value != "" {t.Fatal("presence lost")}
}
`)
		command(t, dir, "go", "test", "./...")
	})
	t.Run("rust", func(t *testing.T) {
		dir := t.TempDir()
		write(dir, "types.rs", formatted(t, rustText(), out.rust))
		write(dir, "consumer.rs", `include!("types.rs");
pub fn qualify() -> NestedPage {
 let page: DataPage<Person> = Page { items: vec![Person {name: String::from("Ada"), active: true}], next: Some(String::new()) };
 Page {items: vec![page], next: None}
}
`)
		command(t, dir, "docker", "run", "--rm", "--network=none", "-v", dir+":/work", "-w", "/work", "rust:1.90.0", "rustc", "--edition=2024", "--crate-type=lib", "consumer.rs")
	})
	t.Run("cpp", func(t *testing.T) {
		dir := t.TempDir()
		write(dir, "types.hpp", formatted(t, cppText(), out.cpp))
		e, _, input := generationInput(t, genModule("shadow", genMap(), genEntry("Person", genType(genPerson())), genEntry("Container", genType(genRecord(genField("Person", false, genLocal("Person")))))))
		shadow := generated(t, e, input, cppText(), cppRequest)
		shadowDir := t.TempDir()
		write(shadowDir, "shadow.hpp", formatted(t, cppText(), shadow))
		write(shadowDir, "shadow.cpp", "#include \"shadow.hpp\"\n#include <type_traits>\nstatic_assert(std::is_same_v<decltype(pagination::Container::Person),pagination::Person>);\n")
		command(t, shadowDir, "docker", "run", "--rm", "--network=none", "-v", shadowDir+":/work", "-w", "/work", "gcc:14.2.0", "g++", "-std=c++20", "-fsyntax-only", "shadow.cpp")
		write(dir, "consumer.cpp", `#include "types.hpp"
#include <type_traits>
static_assert(std::is_same_v<pagination::PeoplePage, pagination::Page<pagination::Person>>);
static_assert(std::is_same_v<pagination::DataPage<pagination::Person>, pagination::PeoplePage>);
static_assert(std::is_same_v<decltype(pagination::Page<int>::next), std::optional<std::string>>);
static_assert(std::is_same_v<decltype(pagination::NestedPage::items), std::vector<pagination::Page<pagination::Person>>>);
`)
		command(t, dir, "docker", "run", "--rm", "--network=none", "-v", dir+":/work", "-w", "/work", "gcc:14.2.0", "g++", "-std=c++20", "-fsyntax-only", "consumer.cpp")
	})
}
