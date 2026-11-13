package integration_test

import (
	analyzerust "github.com/Bitspark/bittype/analyze/files/rust/provider/go"
	checkrust "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	parserust "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"

	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	"github.com/Bitspark/bittype/files/go"

	resolve "github.com/Bitspark/bittype/resolve/files/rust/provider/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/rust/go"
	s "github.com/Bitspark/bittype/syntax/rust/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func rsCompiler(t *testing.T) resolve.Compiler {
	t.Helper()
	c := resolve.Compiler{HelperPath: os.Getenv("BITTYPE_RUST_HELPER"), SysrootPath: os.Getenv("BITTYPE_RUST_SYSROOT")}
	if err := c.Validate(); err != nil {
		t.Fatal("Provision the pinned Rust compiler and set BITTYPE_RUST_HELPER and BITTYPE_RUST_SYSROOT:", err)
	}
	return c
}
func rsCapture(t *testing.T, texts ...string) (semantic.Snapshot, []s.AST, *files.Context) {
	t.Helper()
	var trees []s.AST
	source := &files.Context{}
	for i, text := range texts {
		id := fmt.Sprint(i)
		parsed, err := parserust.Parse(context.Background(), parse.Input{Text: text, SourceID: ast.SourceID(id), TreeID: ast.TreeID(id)})
		if err != nil {
			t.Fatal(err)
		}
		trees = append(trees, parsed.ASTView())
		source.Files.Entries = append(source.Files.Entries, files.File{ID: files.ID(id), Path: files.Path(id + ".rs"), Text: text})
		source.Links = append(source.Links, files.SyntaxLink{File: files.ID(id), Tree: ast.TreeID(id), Source: ast.SourceID(id)})
	}
	return rsCaptureTrees(t, trees, source), trees, source
}
func rsCaptureTrees(t *testing.T, trees []s.AST, source *files.Context) semantic.Snapshot {
	t.Helper()
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "rust-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[s.Language](definition, checkrust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	var envelopes []representation.AST
	for _, tree := range trees {
		env, err := representation.Pack[s.Language](definition, tree)
		if err != nil {
			t.Fatal(err)
		}
		envelopes = append(envelopes, env)
	}
	snapshot, err := semantic.Capture(context.Background(), checker, semantic.Input{Entry: ast.NodeRef{Tree: trees[0].ID, Node: trees[0].Root}, Trees: envelopes, Files: source})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}
func rsResults(t *testing.T, snapshot semantic.Snapshot) (semantic.Result, semantic.Result) {
	t.Helper()
	resolver, err := resolve.New(snapshot.Language(), rsCompiler(t))
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), n.AnalysisContract(), analyzerust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}})
	if err != nil {
		t.Fatal(err)
	}
	return resolved, result
}
func rsModel(t *testing.T, texts ...string) n.Model {
	t.Helper()
	snapshot, _, _ := rsCapture(t, texts...)
	_, result := rsResults(t, snapshot)
	model, err := n.Open(result)
	if err != nil {
		t.Fatalf("model: %v; diagnostics: %+v", err, result.Diagnostics())
	}
	return model
}
func rsDecl(t *testing.T, m n.Model, tree ast.TreeID, name string) n.Declaration {
	t.Helper()
	d, err := m.Lookup("captured::" + name)
	if err != nil {
		t.Fatal(name, err)
	}
	return d
}
func rsType(t *testing.T, m n.Model, name string) n.Type {
	t.Helper()
	return rsDecl(t, m, "0", name).Type()
}

func qualifyRust(t *testing.T) {
	ctx := context.Background()
	snapshot, _, _ := rsCapture(t, "pub struct Person { pub name: String } pub struct Page<T> { pub items: Vec<T>, pub next: Option<String> } pub type People = Page<Person>;")
	var e engine.Engine
	if err := e.RegisterLanguage(snapshot.Language()); err != nil {
		t.Fatal(err)
	}
	if _, err := e.Resolve(ctx, snapshot, semantic.Options{}); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	resolver, err := resolve.New(snapshot.Language(), rsCompiler(t))
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	resolved, err := e.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	input := analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}}
	if _, err = e.Analyze(ctx, input); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), n.AnalysisContract(), analyzerust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	if err = e.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	result, err := e.Analyze(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	rg, _ := resolved.Semantics()
	ag, _ := result.Semantics()
	if !rg.Same(ag) {
		t.Fatal("analysis rebuilt compiler facts", result.Diagnostics())
	}
	m, err := n.Open(result)
	if err != nil {
		t.Fatal(err, result.Diagnostics())
	}
	if _, err = n.Open(resolved); err == nil {
		t.Fatal("resolution admitted as analysis")
	}
	if len(m.Graph().Data().Bindings) == 0 || len(m.Graph().Data().Expressions) == 0 {
		t.Fatal("missing native occurrences")
	}
	page := rsDecl(t, m, "0", "Page")
	if len(page.Parameters()) != 1 || page.Parameters()[0].Entity().Owner != page.ID() {
		t.Fatal("binder ownership")
	}
	people := rsType(t, m, "People")
	fields, ok := people.Fields()
	if !ok || len(fields) != 2 || fields[0].Name != "items" || !fields[0].Public {
		t.Fatal(fields)
	}
	if !strings.Contains(fields[0].Type.String(), "Person") {
		t.Fatal("missing field substitution", fields[0].Type.String())
	}
	pairs, ok := m.Graph().DirectParameterBindings(people.ID())
	if !ok || len(pairs) != 1 || pairs[0].Parameter != page.Parameters()[0].ID() || pairs[0].Argument != rsType(t, m, "Person").ID() {
		t.Fatal("effective application", pairs)
	}
	data := m.Graph().Data()
	data.Entities[0].Name = "mutated"
	if reflect.DeepEqual(data, m.Graph().Data()) {
		t.Fatal("mutable graph")
	}
	if _, err = analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Options: semantic.Options{Configuration: n.Configuration()}}); err == nil {
		t.Fatal("hidden resolver")
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/rust/provider/go#Provider.Resolve feature=rust-semantic-subset case=native-types-and-generics expects=native-identity-and-substitution-preserved
// bittype:case target=github.com/Bitspark/bittype/analyze/files/rust/provider/go#Provider.Analyze feature=rust-semantic-subset case=native-types-and-generics expects=native-identity-and-substitution-preserved
func TestRustSemanticNativeTypes(t *testing.T) {
	m := rsModel(t, `pub struct Person {pub name: String}
pub struct Page<T=Person> {pub items: Vec<T>}
pub type DefaultPage = Page;
pub type People = Page<Person>;
pub type Nested = Page<Page<Person>>;
pub struct Boxed<T> {pub value: T}
pub struct Recursive {pub next: *const Recursive}
pub type Pair = (u8, bool);
pub type Bytes = [u8; 8];
pub type Text = &'static str;
pub enum Choice<T> {Empty, Value(T), Named {value: T}}
pub mod models {pub struct Item {pub value: u32}}
pub type Item = models::Item;`)
	if !rsType(t, m, "People").Same(rsType(t, m, "DefaultPage")) {
		t.Fatal("native alias identity")
	}
	page := rsDecl(t, m, "0", "Page")
	if !page.Parameters()[0].Default().Same(rsType(t, m, "Person")) {
		t.Fatal("generic default")
	}
	fields, _ := rsType(t, m, "Nested").Fields()
	arguments := fields[0].Type.Term().Arguments
	if len(arguments) == 0 || arguments[0] != rsType(t, m, "People").ID() {
		t.Fatal("nested substitution", arguments)
	}
	other := rsDecl(t, m, "0", "Boxed")
	if page.Parameters()[0].ID() == other.Parameters()[0].ID() {
		t.Fatal("collapsed generic owners")
	}
	pf, _ := page.Type().Fields()
	bf, _ := other.Type().Fields()
	args := pf[0].Type.Term().Arguments
	if len(args) == 0 || args[0] == bf[0].Type.ID() {
		t.Fatal("same-spelled parameters collapsed")
	}
	if !rsType(t, m, "Item").Same(rsType(t, m, "models::Item")) {
		t.Fatal("module binding")
	}
	if len(m.Graph().Facts(semantics.Subject{Term: rsType(t, m, "Choice").ID()}, n.Key("variants"))) != 1 {
		t.Fatal("enum variants")
	}
}

func TestRustSemanticNamespaceAndCoordinates(t *testing.T) {
	m := rsModel(t, "// café\r\npub struct Boxed<T> {pub T: T}\r\npub type Small = Boxed<u8>;\r\n")
	if _, err := m.Lookup("captured::Boxed::T"); !errors.Is(err, n.ErrForeign) {
		t.Fatal("ambiguous native namespaces selected by spelling", err)
	}
	box := rsDecl(t, m, "0", "Boxed")
	fields, _ := box.Type().Fields()
	if len(fields) != 1 || fields[0].Type.Term().Target != box.Parameters()[0].ID() {
		t.Fatal("namespace or CRLF/UTF-8 source association", fields)
	}
	zero := n.Model{}
	if zero.Valid() || zero.Same(m) || zero.Receipt().Valid() {
		t.Fatal("zero model")
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/rust/provider/go#Provider.Resolve feature=rust-semantic-subset case=bounds-and-profile-outcomes expects=qualified-reports-without-partial-model
// bittype:case target=github.com/Bitspark/bittype/analyze/files/rust/provider/go#Provider.Analyze feature=rust-semantic-subset case=bounds-and-profile-outcomes expects=qualified-reports-without-partial-model
// The probe of row 26 of research 0016. The helper binds every segment of an
// absolute path into the standard library. It names each entity by the crate
// and the path of its definition, which is not the path that was written, and
// such an entity has no declaration site. A primitive type has no binding.
// docs/destinations.md states what this means for a recorded destination.
func TestRustSemanticAbsoluteStandardLibraryPaths(t *testing.T) {
	snapshot, trees, _ := rsCapture(t, `pub struct Person { pub name: ::std::string::String, pub tags: ::std::vec::Vec<::std::string::String>, pub nick: ::std::option::Option<::std::string::String>, pub ok: bool }`)
	_, result := rsResults(t, snapshot)
	for _, obligation := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		if report, ok := result.Report(obligation); !ok || report.Status != semantic.Complete || report.Verdict != semantic.Holds {
			t.Fatalf("%s does not hold: %+v; %+v", obligation.Name, report, result.Diagnostics())
		}
	}
	m, err := n.Open(result)
	if err != nil {
		t.Fatalf("model: %v; diagnostics: %+v", err, result.Diagnostics())
	}
	tree, data := trees[0], m.Graph().Data()
	type entity struct{ kind, name string }
	bound := map[string]entity{}
	sites := map[ast.NodeID]bool{}
	for _, b := range data.Bindings {
		node := tree.Nodes[b.At.Site.Node-1]
		if b.At.Site.Tree != tree.ID || node.Kind != s.KindPathSegment || b.At.Role != n.Key(s.RoleTarget) || sites[b.At.Site.Node] {
			t.Errorf("unexpected binding at %v (%s) with the role %v", b.At.Site, node.Kind, b.At.Role)
			continue
		}
		sites[b.At.Site.Node] = true
		target := data.Entities[b.Target-1]
		if len(target.Declarations) != 0 {
			t.Errorf("the entity %s of the standard library has a declaration site", target.Name)
		}
		written, to := nodeAttribute(node, s.AttributeName), entity{target.Kind.Name, target.Name}
		if previous, seen := bound[written]; seen && previous != to {
			t.Errorf("the segment %s binds to %v and to %v", written, previous, to)
		}
		bound[written] = to
	}
	want := map[string]entity{
		"std":    {"module", "std"},
		"string": {"module", "alloc::string"}, "String": {"struct", "alloc::string::String"},
		"vec": {"module", "alloc::vec"}, "Vec": {"struct", "alloc::vec::Vec"},
		"option": {"module", "core::option"}, "Option": {"enum", "core::option::Option"},
	}
	if !reflect.DeepEqual(bound, want) {
		t.Errorf("the segments bind to %v, want %v", bound, want)
	}
	// Five paths of three segments each.
	if len(data.Bindings) != 15 {
		t.Errorf("%d bindings, want 15", len(data.Bindings))
	}
	primitives := 0
	for i, node := range tree.Nodes {
		id := ast.NodeID(i + 1)
		switch node.Kind {
		case s.KindPathSegment:
			if !sites[id] {
				t.Errorf("the segment %s at node %d has no binding", nodeAttribute(node, s.AttributeName), id)
			}
		case s.KindPrimitiveType:
			primitives++
			if sites[id] {
				t.Errorf("the primitive type at node %d has a binding", id)
			}
		}
	}
	if primitives != 1 {
		t.Errorf("%d primitive types, want the one bool", primitives)
	}
}
func TestRustSemanticBoundsAndOutcomes(t *testing.T) {
	m := rsModel(t, `pub trait Marker {} pub struct Person {} impl Marker for Person {} pub struct Boxed<T: Marker> {pub value: T} pub type Good = Boxed<Person>;`)
	b := rsDecl(t, m, "0", "Boxed")
	marker := rsDecl(t, m, "0", "Marker")
	if len(b.Predicates()) == 0 || len(marker.Parameters()) != 1 || marker.Parameters()[0].ID() == marker.ID() {
		t.Fatal("trait predicates/Self binder")
	}
	for _, text := range []string{
		`pub trait Marker {} pub struct Person {} pub struct Boxed<T: Marker> {pub value: T} pub struct Bad {pub value: Boxed<Person>}`,
		`pub type Missing = Unknown;`,
		`pub struct Boxed<T> {value:T} pub type Missing = Boxed;`,
		`pub trait Marker {} pub struct Person {} impl Marker for Person {} impl Marker for Person {}`,
	} {
		t.Run(text, func(t *testing.T) {
			snapshot, _, _ := rsCapture(t, text)
			_, result := rsResults(t, snapshot)
			r, _ := result.Report(n.Validity())
			if r.Status != semantic.Complete || r.Verdict != semantic.Rejected || result.HasSemantics() || len(result.Diagnostics()) == 0 {
				t.Fatal(r, result.Diagnostics())
			}
		})
	}
	for _, text := range []string{`pub trait Call {fn call(&self);}`, `pub type Call = fn(u8)->u8;`, `pub struct Borrowed<'a> {value: &'a str}`, `pub trait Family {type Item;}`, `use std::string::String; pub type Text = String;`} {
		t.Run(text, func(t *testing.T) {
			snapshot, _, _ := rsCapture(t, text)
			_, result := rsResults(t, snapshot)
			r, _ := result.Report(n.Coverage())
			if r.Status != semantic.Unsupported || result.HasSemantics() {
				t.Fatal(r, result.Diagnostics())
			}
		})
	}
	snapshot, trees, _ := rsCapture(t, "pub type Byte = u8;")
	snapshot = rsCaptureTrees(t, trees, nil)
	_, result := rsResults(t, snapshot)
	r, _ := result.Report(n.Coverage())
	if r.Status != semantic.Incomplete || result.HasSemantics() {
		t.Fatal(r)
	}
	snapshot, _, _ = rsCapture(t, "pub type A = u8;", "pub type B = u8;")
	_, result = rsResults(t, snapshot)
	r, _ = result.Report(n.Coverage())
	if r.Status != semantic.Unsupported {
		t.Fatal(r)
	}
}

func TestRustSemanticCaptureAndReceipts(t *testing.T) {
	ctx := context.Background()
	snapshot, trees, source := rsCapture(t, "pub struct A {} pub struct B {} pub type Alias = A;")
	resolved, _ := rsResults(t, snapshot)
	resolver, _ := resolve.New(snapshot.Language(), rsCompiler(t))
	analyzer, _ := analysis.NewAnalyzer[s.Language](snapshot.Language(), n.AnalysisContract(), analyzerust.Provider{})
	other := rsCaptureTrees(t, trees, source)
	if _, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: other, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}}); err == nil {
		t.Fatal("foreign receipt")
	}
	for i := range trees[0].Nodes {
		if trees[0].Nodes[i].Kind == s.KindPathSegment {
			trees[0].Nodes[i].Attributes[0].Text = "B"
		}
	}
	forged := rsCaptureTrees(t, trees, source)
	if _, err := resolver.Resolve(ctx, forged, semantic.Options{Configuration: n.Configuration()}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("forged source admitted", err)
	}
	_, trees, source = rsCapture(t, "pub struct A {} pub struct B {} pub type Alias = A;")
	var b ast.NodeID
	for i, node := range trees[0].Nodes {
		if node.Kind == s.KindStruct {
			for _, a := range node.Attributes {
				if a.Name == "name" && a.Text == "B" {
					b = ast.NodeID(i + 1)
				}
			}
		}
	}
	for i := range trees[0].Nodes {
		if trees[0].Nodes[i].Kind == s.KindPathSegment {
			trees[0].Nodes[i].References = []ast.Reference{{Role: s.RoleTarget, Target: ast.NodeRef{Tree: trees[0].ID, Node: b}}}
		}
	}
	forged = rsCaptureTrees(t, trees, source)
	if _, err := resolver.Resolve(ctx, forged, semantic.Options{Configuration: n.Configuration()}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("false reference admitted", err)
	}
	if _, err := resolver.Resolve(ctx, snapshot, semantic.Options{}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("implicit configuration", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := resolver.Resolve(cancelled, snapshot, semantic.Options{Configuration: n.Configuration()}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestRustSemanticConcurrent(t *testing.T) {
	snapshot, _, _ := rsCapture(t, "pub struct Page<T> {pub item:T} pub type Bytes = Page<u8>;")
	resolver, _ := resolve.New(snapshot.Language(), rsCompiler(t))
	analyzer, _ := analysis.NewAnalyzer[s.Language](snapshot.Language(), n.AnalysisContract(), analyzerust.Provider{})
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				result, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}})
				if err != nil {
					t.Error(err)
					return
				}
				m, err := n.Open(result)
				if err != nil {
					t.Error(err, result.Diagnostics())
					return
				}
				data := m.Graph().Data()
				data.Entities[0].Name = "local"
			}
		}()
	}
	wg.Wait()
}

// The compiler helper is Rust and spells the graph names itself. This reads the
// real helper's graph through every name that the vocabulary publishes for
// model.go, one failure per name. No other test opens a native type by id,
// which reads the native-type relation.
func TestRustHelperNamesMatchVocabulary(t *testing.T) {
	m := rsModel(t, `pub struct Person {pub name: String}
pub struct Page<T=Person> {pub items: Vec<T>}
pub type People = Page<Person>;`)
	page, people := rsDecl(t, m, "0", "Page"), rsType(t, m, "People")
	if !people.Valid() {
		t.Error("declared-type relation: an alias has no declared type")
	}
	if opened, err := m.Type(people.ID()); err != nil || !opened.Same(people) {
		t.Error("native-type relation: a native type is not opened by id", err)
	}
	if people.String() == "" {
		t.Error("display relation: a native type has no display")
	}
	if len(m.Graph().Facts(semantics.Subject{Term: people.ID()}, n.FieldsRelation())) != 1 {
		t.Error("fields relation: a struct type has no fields fact")
	} else if fields, ok := people.Fields(); !ok || len(fields) != 1 {
		t.Error("fields operator: the fields fact carries no fields term", fields)
	}
	if len(page.Parameters()) != 1 || !page.Parameters()[0].Default().Valid() {
		t.Error("default relation: a defaulted parameter has no default")
	}
	bounded := rsModel(t, `pub trait Marker {} pub struct Person {} impl Marker for Person {} pub struct Boxed<T: Marker> {pub value: T} pub type Good = Boxed<Person>;`)
	if len(rsDecl(t, bounded, "0", "Boxed").Predicates()) == 0 {
		t.Error("predicate relation: a bounded declaration has no predicate")
	}
}
