package integration_test

import (
	analyzerust "github.com/Bitspark/bittype/analyze/files/rust/provider/go"
	checkrust "github.com/Bitspark/bittype/check/syntax/rust/provider/go"
	check "github.com/Bitspark/bittype/contracts/check/go"
	parserust "github.com/Bitspark/bittype/parse/syntax/rust/provider/go"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"
	"unicode"

	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
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

// rsEnvironment is the environment of these tests. Its crate name is not the
// name of any profile, so a path that starts with it shows that the compiler
// analyzed under the stated name.
func rsEnvironment() n.Environment {
	return n.Environment{Edition: n.Edition2024, Target: n.TargetLinuxX64, Crate: n.Crate{Name: "pages", Kind: n.KindLibrary}}
}
func rsSettingsFor(t *testing.T, environment n.Environment) []semantic.Setting {
	t.Helper()
	settings, err := n.ConfigurationFor(environment)
	if err != nil {
		t.Fatal(err)
	}
	return settings
}
func rsSettings(t *testing.T) []semantic.Setting {
	t.Helper()
	return rsSettingsFor(t, rsEnvironment())
}
func rsContract(t *testing.T, settings []semantic.Setting) analysis.Contract {
	t.Helper()
	contract, err := n.AnalysisContractFor(settings)
	if err != nil {
		t.Fatal(err)
	}
	return contract
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
	return rsResultsFor(t, snapshot, rsSettings(t))
}
func rsResultsFor(t *testing.T, snapshot semantic.Snapshot, settings []semantic.Setting) (semantic.Result, semantic.Result) {
	t.Helper()
	resolver, err := resolve.New(snapshot.Language(), rsCompiler(t))
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: settings})
	if err != nil {
		t.Fatal(err)
	}
	return resolved, rsAnalyze(t, snapshot, resolved, settings)
}
func rsAnalyze(t *testing.T, snapshot semantic.Snapshot, resolved semantic.Result, settings []semantic.Setting) semantic.Result {
	t.Helper()
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), rsContract(t, settings), analyzerust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: settings}})
	if err != nil {
		t.Fatal(err)
	}
	return result
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
	d, err := m.Lookup(rsEnvironment().Crate.Name + "::" + name)
	if err != nil {
		t.Fatal(name, err)
	}
	return d
}

// rsRoot is the crate root, the module that pub(crate) restricts to.
func rsRoot(t *testing.T, m n.Model) n.Declaration {
	t.Helper()
	d, err := m.Lookup(rsEnvironment().Crate.Name)
	if err != nil {
		t.Fatal(err)
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
	resolved, err := e.Resolve(ctx, snapshot, semantic.Options{Configuration: rsSettings(t)})
	if err != nil {
		t.Fatal(err)
	}
	input := analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: rsSettings(t)}}
	if _, err = e.Analyze(ctx, input); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), rsContract(t, rsSettings(t)), analyzerust.Provider{})
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
	fields, ok, err := people.Fields()
	if err != nil || !ok || len(fields) != 2 || fields[0].Name != "items" {
		t.Fatal(fields, err)
	}
	if visibility, err := fields[0].Declaration.Visibility(); err != nil || !visibility.Public {
		t.Fatal(visibility, err)
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
	if _, err = analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Options: semantic.Options{Configuration: rsSettings(t)}}); err == nil {
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
	fields, _, _ := rsType(t, m, "Nested").Fields()
	arguments := fields[0].Type.Term().Arguments
	if len(arguments) == 0 || arguments[0] != rsType(t, m, "People").ID() {
		t.Fatal("nested substitution", arguments)
	}
	other := rsDecl(t, m, "0", "Boxed")
	if page.Parameters()[0].ID() == other.Parameters()[0].ID() {
		t.Fatal("collapsed generic owners")
	}
	pf, _, _ := page.Type().Fields()
	bf, _, _ := other.Type().Fields()
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
	if _, err := m.Lookup(rsEnvironment().Crate.Name + "::Boxed::T"); !errors.Is(err, n.ErrForeign) {
		t.Fatal("ambiguous native namespaces selected by spelling", err)
	}
	box := rsDecl(t, m, "0", "Boxed")
	fields, _, _ := box.Type().Fields()
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
	analyzer, _ := analysis.NewAnalyzer[s.Language](snapshot.Language(), rsContract(t, rsSettings(t)), analyzerust.Provider{})
	other := rsCaptureTrees(t, trees, source)
	if _, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: other, Resolution: &resolved, Options: semantic.Options{Configuration: rsSettings(t)}}); err == nil {
		t.Fatal("foreign receipt")
	}
	for i := range trees[0].Nodes {
		if trees[0].Nodes[i].Kind == s.KindPathSegment {
			trees[0].Nodes[i].Attributes[0].Text = "B"
		}
	}
	forged := rsCaptureTrees(t, trees, source)
	if _, err := resolver.Resolve(ctx, forged, semantic.Options{Configuration: rsSettings(t)}); !errors.Is(err, semantic.ErrInvalidRequest) {
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
	if _, err := resolver.Resolve(ctx, forged, semantic.Options{Configuration: rsSettings(t)}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("false reference admitted", err)
	}
	if _, err := resolver.Resolve(ctx, snapshot, semantic.Options{}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("implicit configuration", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := resolver.Resolve(cancelled, snapshot, semantic.Options{Configuration: rsSettings(t)}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}

func TestRustSemanticConcurrent(t *testing.T) {
	snapshot, _, _ := rsCapture(t, "pub struct Page<T> {pub item:T} pub type Bytes = Page<u8>;")
	resolver, _ := resolve.New(snapshot.Language(), rsCompiler(t))
	settings := rsSettings(t)
	analyzer, _ := analysis.NewAnalyzer[s.Language](snapshot.Language(), rsContract(t, settings), analyzerust.Provider{})
	resolved, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: settings})
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				result, err := analyzer.Analyze(context.Background(), analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: settings}})
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
	} else if fields, ok, err := people.Fields(); err != nil || !ok || len(fields) != 1 {
		t.Error("fields operator: the fields fact carries no fields term", fields, err)
	}
	if len(page.Parameters()) != 1 || !page.Parameters()[0].Default().Valid() {
		t.Error("default relation: a defaulted parameter has no default")
	}
	bounded := rsModel(t, `pub trait Marker {} pub struct Person {} impl Marker for Person {} pub struct Boxed<T: Marker> {pub value: T} pub type Good = Boxed<Person>;`)
	if len(rsDecl(t, bounded, "0", "Boxed").Predicates()) == 0 {
		t.Error("predicate relation: a bounded declaration has no predicate")
	}
	// The visibility names. Each is read from the raw graph, so that a name
	// the vocabulary spells differently from the helper fails by itself.
	visible := rsModel(t, `pub struct Public { pub field: Hidden } struct Hidden {}`)
	graph := visible.Graph()
	fact := func(d n.Declaration, relation semantics.Key) semantics.Term {
		terms := graph.Facts(semantics.Subject{Entity: d.ID()}, relation)
		if len(terms) != 1 {
			return semantics.Term{}
		}
		term, _ := graph.Term(terms[0])
		return term
	}
	public, hidden := rsDecl(t, visible, "0", "Public"), rsDecl(t, visible, "0", "Hidden")
	if len(graph.Facts(semantics.Subject{Entity: public.ID()}, n.DeclaredVisibilityRelation())) != 1 {
		t.Error("declared-visibility relation: an item has no declared visibility")
	}
	if fact(public, n.Key("declared-visibility")).Operator != n.PublicOperator() {
		t.Error("public operator: a public item's declared visibility has another operator")
	}
	restricted := fact(hidden, n.Key("declared-visibility"))
	if restricted.Operator != n.RestrictedOperator() {
		t.Error("restricted operator: a private item's declared visibility has another operator")
	}
	if len(restricted.Operands) != 1 {
		t.Error("restricted operator: no scope operand")
	} else if scope, _ := graph.Term(restricted.Operands[0].Term); rsRoot(t, visible).Entity().Kind != n.ModuleKind() || scope.Target != rsRoot(t, visible).ID() {
		t.Error("module kind: the scope of a private root item is not the crate root module")
	}
	if len(graph.Facts(semantics.Subject{Entity: public.ID()}, n.EffectiveVisibilityRelation())) != 1 {
		t.Error("effective-visibility relation: an item has no effective visibility")
	}
	effective := fact(public, n.Key("effective-visibility"))
	if effective.Operator != n.EffectiveVisibilityOperator() {
		t.Error("effective-visibility operator: the effective fact has another operator")
	}
	var stated struct{ Audience, Path string }
	if err := json.Unmarshal(effective.Data, &stated); err != nil {
		t.Error(err)
	}
	if stated.Audience != n.AudienceAnotherCrate {
		t.Error("audience: the effective fact states another audience", stated.Audience)
	}
	if stated.Path != n.PathDefinition {
		t.Error("path: the effective fact states another path", stated.Path)
	}
	fields, _, _ := public.Type().Fields()
	if len(fields) != 1 || len(graph.Facts(semantics.Subject{Entity: fields[0].Declaration.ID()}, n.PrivateInterfaceRelation())) != 1 {
		t.Error("private-interface relation: a field with a more private type has no observation")
	}
	// A field's visibility is the fact on its entity and nothing else: the
	// members of the fields term are the two that model.go reads.
	var members []map[string]json.RawMessage
	listed, _ := graph.Term(graph.Facts(semantics.Subject{Term: public.Type().ID()}, n.FieldsRelation())[0])
	if err := json.Unmarshal(listed.Data, &members); err != nil || len(members) != 1 || len(members[0]) != 2 || members[0]["Name"] == nil || members[0]["Entity"] == nil {
		t.Error("fields operator: a field is described by other members than its name and its entity", string(listed.Data), err)
	}
}

// The tests below are the Verify section of #278. They need the compiler
// helper and run where it is provisioned.

type rsVisibilityCase struct {
	model    n.Model
	name     string
	public   bool
	scope    string
	nameable bool
}

// Declared and effective visibility are two facts. D is declared public and
// another crate cannot name it, because its module is private.
func TestRustSemanticItemVisibility(t *testing.T) {
	m := rsModel(t, `pub struct A {} struct B {} pub(crate) struct C {} mod m { pub struct D {} } pub mod n { pub struct E {} struct F {} }`)
	root := rsRoot(t, m)
	// A restriction is the one that was written, which inside a module is not
	// the enclosing module: G is restricted to the root. Every module around
	// a declaration counts for the effective fact, not only the nearest: b
	// and X are declared public inside the private a.
	deeper := rsModel(t, `pub mod n { pub(crate) struct G {} } mod a { pub mod b { pub struct X {} } }`)
	for _, c := range []rsVisibilityCase{
		{name: "A", public: true, nameable: true},
		{name: "B", scope: ""},
		{name: "C", scope: ""},
		{name: "m", scope: ""},
		{name: "m::D", public: true},
		{name: "n", public: true, nameable: true},
		{name: "n::E", public: true, nameable: true},
		{name: "n::F", scope: "n"},
		{model: deeper, name: "n", public: true, nameable: true},
		{model: deeper, name: "n::G", scope: ""},
		{model: deeper, name: "a", scope: ""},
		{model: deeper, name: "a::b", public: true},
		{model: deeper, name: "a::b::X", public: true},
	} {
		m, root := m, root
		if c.model.Valid() {
			m, root = c.model, rsRoot(t, c.model)
		}
		d := rsDecl(t, m, "0", c.name)
		declared, err := d.Visibility()
		if err != nil {
			t.Errorf("%s: declared visibility: %v", c.name, err)
			continue
		}
		scope := root
		if c.scope != "" {
			scope = rsDecl(t, m, "0", c.scope)
		}
		if c.public && (!declared.Public || declared.Scope.Valid()) {
			t.Errorf("%s: declared %+v, want public", c.name, declared)
		}
		if !c.public && (declared.Public || declared.Scope.ID() != scope.ID()) {
			t.Errorf("%s: declared %+v, want restricted to %s", c.name, declared, scope.Entity().Name)
		}
		effective, err := d.Effective()
		if err != nil {
			t.Errorf("%s: effective visibility: %v", c.name, err)
			continue
		}
		if effective != (n.Effective{Audience: n.AudienceAnotherCrate, Path: n.PathDefinition, Nameable: c.nameable}) {
			t.Errorf("%s: effective %+v, want nameable %v", c.name, effective, c.nameable)
		}
	}
	// The crate root is a scope and no item: nothing declares it.
	if _, err := root.Visibility(); !errors.Is(err, n.ErrMissingFact) {
		t.Error("the crate root has a declared visibility", err)
	}
	// A public field hands D to other crates, which the compiler calls
	// reachable. They still cannot name it, and that is the effective fact.
	leaked := rsModel(t, `mod m { pub struct D {} } pub struct A { pub d: m::D }`)
	if e, err := rsDecl(t, leaked, "0", "m::D").Effective(); err != nil || e.Nameable {
		t.Error("a reachable declaration that cannot be named", e, err)
	}
	if e, err := rsDecl(t, leaked, "0", "A").Effective(); err != nil || !e.Nameable {
		t.Error("the declaration that leaks it", e, err)
	}
}

// Every captured item has both facts and every captured field the declared
// one, whatever its kind. What Rust gives no visibility has neither.
func TestRustSemanticVisibilityCoverage(t *testing.T) {
	m := rsModel(t, `pub trait Marker {} pub struct Person {pub name: String, age: u8} impl Marker for Person {}
pub struct Pair(pub u8, bool);
pub type People<T = Person> = Vec<T>;
pub enum Choice<T> {Empty, Value(T), Named {value: T}}
mod inner {pub(crate) enum Hidden {One} pub trait Local {} pub type Byte = u8;}`)
	seen := map[string]int{}
	for i, e := range m.Graph().Data().Entities {
		if len(e.Declarations) == 0 {
			continue
		}
		d, err := m.Declaration(semantics.EntityID(i + 1))
		if err != nil {
			t.Fatal(err)
		}
		_, declared := d.Visibility()
		_, effective := d.Effective()
		seen[e.Kind.Name]++
		switch e.Kind.Name {
		case "module", "alias", "struct", "enum", "trait":
			if declared != nil || effective != nil {
				t.Errorf("%s %s: declared %v, effective %v", e.Kind.Name, e.Name, declared, effective)
			}
		case "field":
			if declared != nil || !errors.Is(effective, n.ErrMissingFact) {
				t.Errorf("field %s: declared %v, effective %v", e.Name, declared, effective)
			}
		case "variant", "impl", "parameter":
			if !errors.Is(declared, n.ErrMissingFact) || !errors.Is(effective, n.ErrMissingFact) {
				t.Errorf("%s %s has a visibility fact: declared %v, effective %v", e.Kind.Name, e.Name, declared, effective)
			}
		default:
			t.Errorf("unexpected captured kind %s", e.Kind.Name)
		}
	}
	for _, kind := range []string{"module", "alias", "struct", "enum", "trait", "field", "variant", "impl", "parameter"} {
		if seen[kind] == 0 {
			t.Errorf("the source has no captured %s", kind)
		}
	}
	// A sysroot declaration is outside the capture and carries no fact.
	fields, _, _ := rsType(t, m, "Person").Fields()
	text, err := m.Declaration(fields[0].Type.Term().Target)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := text.Visibility(); !errors.Is(err, n.ErrMissingFact) {
		t.Error("a sysroot declaration has a declared visibility", err)
	}
}

// A field's fact has the form of an item's. A yes or no could not tell
// pub(crate) from private.
func TestRustSemanticFieldVisibility(t *testing.T) {
	m := rsModel(t, `mod m { pub struct S { pub(crate) shared: u8, own: u8, pub open: u8 } } pub enum Choice { Named { value: u8 } } enum Private { Named { value: u8 } }`)
	fields, ok, err := rsType(t, m, "m::S").Fields()
	if err != nil || !ok || len(fields) != 3 {
		t.Fatal(fields, err)
	}
	shared, err := fields[0].Declaration.Visibility()
	if err != nil {
		t.Fatal(err)
	}
	own, err := fields[1].Declaration.Visibility()
	if err != nil {
		t.Fatal(err)
	}
	open, err := fields[2].Declaration.Visibility()
	if err != nil {
		t.Fatal(err)
	}
	if shared.Public || shared.Scope.ID() != rsRoot(t, m).ID() {
		t.Error("pub(crate) field is not restricted to the crate root", shared)
	}
	if own.Public || own.Scope.ID() != rsDecl(t, m, "0", "m").ID() {
		t.Error("private field is not restricted to its module", own)
	}
	if shared.Scope.ID() == own.Scope.ID() {
		t.Error("pub(crate) and private read alike")
	}
	if !open.Public || open.Scope.Valid() {
		t.Error("pub field", open)
	}
	// A variant field takes no modifier. The compiler resolves it to the
	// visibility of its enum.
	if v, err := rsDecl(t, m, "0", "Choice::Named::value").Visibility(); err != nil || !v.Public {
		t.Error("field of a variant of a public enum", v, err)
	}
	if v, err := rsDecl(t, m, "0", "Private::Named::value").Visibility(); err != nil || v.Public || v.Scope.ID() != rsRoot(t, m).ID() {
		t.Error("field of a variant of a private enum", v, err)
	}
}

// rsTampered publishes facts under the Rust resolution contract. It stands for
// a bridge that omits a fact.
type rsTampered struct {
	s.Language
	facts semantics.Data
}

func (p rsTampered) Resolve(_ context.Context, request semantic.Request[s.Language]) (semantic.Result, error) {
	publication := semantic.Publication{Facts: &p.facts}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Complete, Verdict: semantic.Holds})
	}
	return request.Publish(publication)
}

// rsWithoutFact opens a model of the same capture whose graph lacks the one
// fact of relation about the named declaration.
func rsWithoutFact(t *testing.T, snapshot semantic.Snapshot, resolved semantic.Result, name string, relation semantics.Key) n.Model {
	t.Helper()
	envelope, _ := resolved.Semantics()
	graph, err := representation.UnpackSemantics[s.Language](snapshot.Language(), envelope)
	if err != nil {
		t.Fatal(err)
	}
	data := graph.Data()
	subject := semantics.EntityID(0)
	for i, e := range data.Entities {
		if e.Name == rsEnvironment().Crate.Name+"::"+name {
			subject = semantics.EntityID(i + 1)
		}
	}
	kept := data.Facts[:0:0]
	for _, fact := range data.Facts {
		if fact.Subject.Entity != subject || fact.Relation != relation {
			kept = append(kept, fact)
		}
	}
	if subject == 0 || len(kept) != len(data.Facts)-1 {
		t.Fatalf("%s has %d facts of %s", name, len(data.Facts)-len(kept), relation.Name)
	}
	data.Facts = kept
	resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), n.ResolutionContract(), rsTampered{facts: data})
	if err != nil {
		t.Fatal(err)
	}
	tampered, err := resolver.Resolve(context.Background(), snapshot, semantic.Options{Configuration: rsSettings(t)})
	if err != nil {
		t.Fatal(err)
	}
	m, err := n.Open(rsAnalyze(t, snapshot, tampered, rsSettings(t)))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// A graph from which one visibility fact was removed makes its reader fail.
// The declaration is never read as private or as not nameable.
func TestRustSemanticMissingVisibilityFactFails(t *testing.T) {
	snapshot, _, _ := rsCapture(t, `pub struct A { pub value: u8 } pub struct B {}`)
	resolved, result := rsResults(t, snapshot)
	whole, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	if v, err := rsDecl(t, whole, "0", "A").Visibility(); err != nil || !v.Public {
		t.Fatal("the untouched graph", v, err)
	}

	m := rsWithoutFact(t, snapshot, resolved, "A", n.DeclaredVisibilityRelation())
	if v, err := rsDecl(t, m, "0", "A").Visibility(); !errors.Is(err, n.ErrMissingFact) || v != (n.Visibility{}) {
		t.Error("an item without its declared fact was read", v, err)
	}
	if _, err := rsDecl(t, m, "0", "A").Effective(); err != nil {
		t.Error("the effective fact of the same item", err)
	}
	if v, err := rsDecl(t, m, "0", "B").Visibility(); err != nil || !v.Public {
		t.Error("another item", v, err)
	}

	m = rsWithoutFact(t, snapshot, resolved, "A", n.EffectiveVisibilityRelation())
	if e, err := rsDecl(t, m, "0", "A").Effective(); !errors.Is(err, n.ErrMissingFact) || e != (n.Effective{}) {
		t.Error("an item without its effective fact was read", e, err)
	}

	m = rsWithoutFact(t, snapshot, resolved, "A::value", n.DeclaredVisibilityRelation())
	fields, ok, err := rsType(t, m, "A").Fields()
	if err != nil || !ok || len(fields) != 1 {
		t.Fatal(fields, err)
	}
	if v, err := fields[0].Declaration.Visibility(); !errors.Is(err, n.ErrMissingFact) || v != (n.Visibility{}) {
		t.Error("a field without its declared fact was read", v, err)
	}
}

// A public record whose field has a more private type is legal Rust and stays
// analyzed. The compiler's observation is a fact on the field, and the same
// source with a public type has none.
func TestRustSemanticPrivateInterface(t *testing.T) {
	read := func(text string) (n.Model, []n.Declaration) {
		t.Helper()
		snapshot, _, _ := rsCapture(t, text)
		_, result := rsResults(t, snapshot)
		if r, _ := result.Report(n.Validity()); r.Status != semantic.Complete || r.Verdict != semantic.Holds {
			t.Fatal("a private interface is no invalid input", r, result.Diagnostics())
		}
		m, err := n.Open(result)
		if err != nil {
			t.Fatal(err, result.Diagnostics())
		}
		fields, ok, err := rsType(t, m, "Public").Fields()
		if err != nil || !ok || len(fields) != 1 {
			t.Fatal(fields, err)
		}
		observed, err := fields[0].Declaration.PrivateInterface()
		if err != nil {
			t.Fatal(err)
		}
		return m, observed
	}
	m, observed := read(`pub struct Public { pub field: Hidden } struct Hidden {}`)
	if len(observed) != 1 || observed[0].ID() != rsDecl(t, m, "0", "Hidden").ID() {
		t.Error("the private interface of Public::field is not observed", observed)
	}
	if on, err := rsDecl(t, m, "0", "Public").PrivateInterface(); err != nil || len(on) != 0 {
		t.Error("the observation is about the field, as the compiler reports it", on, err)
	}
	if _, observed := read(`pub struct Public { pub field: Hidden } pub struct Hidden {}`); len(observed) != 0 {
		t.Error("a public type is observed as private", observed)
	}
	if _, observed := read(`pub struct Public { field: Hidden } struct Hidden {}`); len(observed) != 0 {
		t.Error("a private field is observed", observed)
	}

	// The compiler reports the same for an alias, a tuple field and a variant
	// field, and once for each more private type.
	m = rsModel(t, `struct Hidden {} struct Other {} pub type Alias = Hidden; pub struct Tuple(pub Hidden, pub Other); pub enum Choice { Value(Hidden) } pub struct Both { pub field: (Hidden, Other) } pub struct Twice { pub field: (Hidden, Hidden) } pub struct Defaulted<T = Hidden> { pub value: T } pub type Chosen<T = Other> = Vec<T>;`)
	hidden, other := rsDecl(t, m, "0", "Hidden").ID(), rsDecl(t, m, "0", "Other").ID()
	for name, want := range map[string][]semantics.EntityID{
		"Alias":            {hidden},
		"Tuple::0":         {hidden},
		"Tuple::1":         {other},
		"Choice::Value::0": {hidden},
		"Both::field":      {hidden, other},
		"Twice::field":     {hidden},
		// For the default of a parameter the compiler reports the item itself.
		"Defaulted":        {hidden},
		"Defaulted::value": nil,
		"Chosen":           {other},
		"Hidden":           nil,
	} {
		observed, err := rsDecl(t, m, "0", name).PrivateInterface()
		if err != nil {
			t.Fatal(name, err)
		}
		var got []semantics.EntityID
		for _, d := range observed {
			got = append(got, d.ID())
		}
		slices.Sort(got)
		slices.Sort(want)
		if !slices.Equal(got, want) {
			t.Errorf("%s: observed %v, want %v", name, got, want)
		}
	}

	// The limit of the fact: it follows the declared visibility of the type.
	// D is declared public, so the field carries no observation, although
	// another crate cannot name D. The effective fact of D says that.
	m = rsModel(t, `mod m { pub struct D {} } pub struct A { pub d: m::D }`)
	if observed, err := rsDecl(t, m, "0", "A::d").PrivateInterface(); err != nil || len(observed) != 0 {
		t.Error("a public type in a private module is observed", observed, err)
	}
	if effective, err := rsDecl(t, m, "0", "m::D").Effective(); err != nil || effective.Nameable {
		t.Error("the type of that field is nameable", effective, err)
	}

	// A record with many such fields resolves, and every field has its
	// observation.
	const many = 600
	var text strings.Builder
	text.WriteString("struct Hidden {} pub struct Wide {")
	for i := 0; i < many; i++ {
		fmt.Fprintf(&text, " pub f%d: Hidden,", i)
	}
	text.WriteString(" }")
	m = rsModel(t, text.String())
	fields, ok, err := rsType(t, m, "Wide").Fields()
	if err != nil || !ok || len(fields) != many {
		t.Fatal("the wide record", len(fields), ok, err)
	}
	for _, field := range fields {
		if observed, err := field.Declaration.PrivateInterface(); err != nil || len(observed) != 1 {
			t.Fatal(field.Name, observed, err)
		}
	}
}

// The compiler analyzes under the environment's crate name, and the receipt
// is bound to the environment.
func TestRustSemanticEnvironment(t *testing.T) {
	ctx := context.Background()
	snapshot, _, _ := rsCapture(t, `pub struct A {}`)
	other := rsEnvironment()
	other.Crate.Name = "other"
	_, result := rsResultsFor(t, snapshot, rsSettingsFor(t, other))
	m, err := n.Open(result)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := m.Lookup("other::A"); err != nil {
		t.Error("the compiler did not analyze under the stated crate name", err)
	}
	if _, err := m.Lookup(rsEnvironment().Crate.Name + "::A"); !errors.Is(err, n.ErrForeign) {
		t.Error("a path of another crate name resolves", err)
	}

	// Analysis claims the environment of its resolution and no other.
	resolved, _ := rsResults(t, snapshot)
	for name, contract := range map[string][]semantic.Setting{"contract of the resolution": rsSettings(t), "contract of the request": rsSettingsFor(t, other)} {
		analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), rsContract(t, contract), analyzerust.Provider{})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: rsSettingsFor(t, other)}}); err == nil {
			t.Errorf("%s: analysis under one environment took the resolution of another", name)
		}
	}

	// A compiler that does not exist shows that a refusal precedes the compiler.
	absent := resolve.Compiler{HelperPath: rsCompiler(t).SysrootPath + "/no-helper", SysrootPath: rsCompiler(t).SysrootPath}
	resolver, err := resolve.New(snapshot.Language(), absent)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: rsSettings(t)}); !errors.Is(err, n.ErrCompiler) {
		t.Fatal("the absent compiler was not reached", err)
	}
	setting := func(value string) []semantic.Setting {
		settings := rsSettings(t)
		for i := range settings {
			if settings[i].Key == n.Key("environment-1") {
				settings[i].Value = []byte(value)
				return settings
			}
		}
		t.Fatal("no environment setting")
		return nil
	}
	analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), rsContract(t, rsSettings(t)), analyzerust.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	for name, value := range map[string]string{
		"edition 2021":   `{"Edition":"2021","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`,
		"another target": `{"Edition":"2024","Target":"aarch64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[]}`,
		"a binary crate": `{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"bin"},"Dependencies":[]}`,
		"one dependency": `{"Edition":"2024","Target":"x86_64-unknown-linux-gnu","Crate":{"Name":"pages","Kind":"lib"},"Dependencies":[{"Name":"serde"}]}`,
	} {
		if _, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: setting(value)}); !errors.Is(err, n.ErrUnsupported) || errors.Is(err, semantic.ErrInvalidRequest) {
			t.Errorf("resolution with %s: %v, want unsupported", name, err)
		}
		if _, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: setting(value)}}); !errors.Is(err, n.ErrUnsupported) || errors.Is(err, semantic.ErrInvalidRequest) {
			t.Errorf("analysis with %s: %v, want unsupported", name, err)
		}
	}
	// The configuration of the earlier helper profile is no configuration of
	// this one.
	earlier := []semantic.Setting{
		{Key: n.Key("compiler"), Value: []byte("1.90.0")},
		{Key: n.Key("environment"), Value: []byte("edition=2024;target=x86_64-unknown-linux-gnu;crate=captured;lib;sysroot-only;no-codegen")},
		{Key: n.Key("helper"), Value: []byte("captured-types-1")},
	}
	if _, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: earlier}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("resolution with the earlier configuration", err)
	}
	if _, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: earlier}}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("analysis with the earlier configuration", err)
	}
}

// rsReusing publishes its resolution's graph as an analysis. It stands for an
// analyzer of another profile.
type rsReusing struct{ s.Language }

func (rsReusing) Analyze(_ context.Context, request analysis.Request[s.Language]) (semantic.Result, error) {
	resolved, _ := request.Resolution()
	publication := semantic.Publication{ReuseSemanticsFrom: &resolved}
	for _, key := range []semantics.Key{n.Coverage(), n.Bindings(), n.Validity()} {
		publication.Reports = append(publication.Reports, semantic.Report{Obligation: key, Status: semantic.Complete, Verdict: semantic.Holds})
	}
	return request.Publish(publication)
}

// The profile's vocabulary changed incompatibly, so the profile has a new
// version. A receipt of the earlier version is refused by its profile alone:
// these receipts carry the configuration of this version. The test needs no
// compiler.
func TestRustSemanticEarlierProfileIsRefused(t *testing.T) {
	ctx := context.Background()
	snapshot, _, _ := rsCapture(t, `pub struct A {}`)
	settings := rsSettings(t)
	const earlierResolution, earlierAnalysis = "rustc-1.90.0/captured-types-1/resolution-1", "rustc-1.90.0/captured-types-1/analysis-1"
	resolve := func(profile string) semantic.Result {
		t.Helper()
		contract := n.ResolutionContract()
		contract.Semantic.Profile = profile
		resolver, err := resolution.NewResolver[s.Language](snapshot.Language(), contract, rsTampered{})
		if err != nil {
			t.Fatal(err)
		}
		resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: settings})
		if err != nil {
			t.Fatal(err)
		}
		return resolved
	}
	analyze := func(profile string, resolved semantic.Result) (semantic.Result, error) {
		t.Helper()
		contract := rsContract(t, settings)
		contract.Semantic.Profile = profile
		analyzer, err := analysis.NewAnalyzer[s.Language](snapshot.Language(), contract, rsReusing{})
		if err != nil {
			t.Fatal(err)
		}
		return analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: settings}})
	}

	// The same graph and configuration under this version are admitted.
	current := resolve(n.ResolutionProfile)
	analyzed, err := analyze(n.AnalysisProfile, current)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := n.Open(analyzed); err != nil {
		t.Fatal("a receipt of this profile", err)
	}

	// Analysis does not take a resolution receipt of the earlier profile.
	if _, err := analyze(n.AnalysisProfile, resolve(earlierResolution)); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("a resolution receipt of the earlier profile was admitted", err)
	}
	// The reader does not open an analysis receipt of the earlier profile.
	earlier, err := analyze(earlierAnalysis, current)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := n.Open(earlier); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Error("an analysis receipt of the earlier profile was opened", err)
	}
}

// rsRustc runs the pinned compiler on a one-line source and reports whether
// it took the arguments. Its answer is the vocabulary that the environment
// follows.
func rsRustc(t *testing.T, arguments ...string) (bool, string) {
	t.Helper()
	directory := t.TempDir()
	source := filepath.Join(directory, "source.rs")
	if err := os.WriteFile(source, []byte("fn main() {}\n"), 0600); err != nil {
		t.Fatal(err)
	}
	command := exec.Command(filepath.Join(rsCompiler(t).SysrootPath, "bin", "rustc"), append(arguments, "-A", "warnings", "--out-dir", directory, source)...)
	output, err := command.CombinedOutput()
	var exit *exec.ExitError
	if err != nil && !errors.As(err, &exit) {
		t.Fatal(err)
	}
	return err == nil, string(output)
}

// The generator of the crate-name table of semantics/rust/go/cratename.go: it
// prints, as ranges, every code point that the compiler's rule for a
// character of a crate name takes.
const rsCrateNameCharacters = `fn main() {
    let mut start: Option<u32> = None;
    for u in 0..=0x110000u32 {
        let ok = char::from_u32(u).map_or(false, |c| c.is_alphanumeric() || c == '_');
        match (ok, start) {
            (true, None) => start = Some(u),
            (false, Some(s)) => { println!("{:X}..{:X}", s, u - 1); start = None; }
            _ => {}
        }
    }
}
`

// The environment's vocabulary is the pinned compiler's. This puts the same
// texts to the compiler and to the environment and compares the answers: a
// text the compiler takes is admitted or unsupported, and a text it refuses
// is invalid. For a crate name the comparison covers every code point.
func TestRustEnvironmentEqualsCompiler(t *testing.T) {
	invalid := func(change func(*n.Environment)) bool {
		t.Helper()
		environment := rsEnvironment()
		change(&environment)
		_, err := n.ConfigurationFor(environment)
		if err != nil && errors.Is(err, semantic.ErrInvalidRequest) == errors.Is(err, n.ErrUnsupported) {
			t.Fatal("a refusal that is not exactly one of invalid and unsupported", err)
		}
		return errors.Is(err, semantic.ErrInvalidRequest)
	}
	base := []string{"--emit", "metadata", "--edition", "2024"}

	// Every character of a crate name.
	directory := t.TempDir()
	generator := filepath.Join(directory, "characters.rs")
	if err := os.WriteFile(generator, []byte(rsCrateNameCharacters), 0600); err != nil {
		t.Fatal(err)
	}
	program := filepath.Join(directory, "characters")
	if output, err := exec.Command(filepath.Join(rsCompiler(t).SysrootPath, "bin", "rustc"), "-O", "-o", program, generator).CombinedOutput(); err != nil {
		t.Fatal(err, string(output))
	}
	printed, err := exec.Command(program).Output()
	if err != nil {
		t.Fatal(err)
	}
	taken := map[rune]bool{}
	for _, line := range strings.Fields(string(printed)) {
		var low, high rune
		if _, err := fmt.Sscanf(line, "%X..%X", &low, &high); err != nil {
			t.Fatal(line, err)
		}
		for r := low; r <= high; r++ {
			taken[r] = true
		}
	}
	if len(taken) < 100000 {
		t.Fatal("the compiler's set of crate name characters", len(taken))
	}
	differ := 0
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		if invalid(func(e *n.Environment) { e.Crate.Name = "a" + string(r) }) == taken[r] {
			if differ++; differ <= 10 {
				t.Errorf("U+%04X in a crate name: the compiler takes it: %v", r, taken[r])
			}
		}
	}
	if differ != 0 {
		t.Errorf("%d characters of a crate name differ from the compiler", differ)
	}

	// Whole crate names, put to the compiler as it is called.
	for _, name := range []string{"pages", "_", "__", "1a", "é", "x²", "ⅷ", "٣", "a٣", "ª", "ⓐ", "日本語", "self", "crate", "Self", "type", "std", "core", "test",
		"", "a-b", "a b", "a.b", "a\u0301", "r#type", "-a", "pages\n", "a/b", "a::b"} {
		takes, output := rsRustc(t, append(base, "--crate-type", "lib", "--crate-name="+name)...)
		if takes == invalid(func(e *n.Environment) { e.Crate.Name = name }) {
			t.Errorf("crate name %q: the compiler takes it: %v\n%s", name, takes, output)
		}
		// A dependency's name is judged by the same rule.
		if takes == invalid(func(e *n.Environment) { e.Dependencies = []n.Dependency{{Name: name}} }) {
			t.Errorf("dependency name %q: the compiler takes it as a crate name: %v", name, takes)
		}
	}

	// Editions and crate kinds.
	for _, edition := range []string{"2015", "2018", "2021", "2024", "future", "2027", "2024 ", "banana", ""} {
		takes, output := rsRustc(t, "--emit", "metadata", "--crate-type", "lib", "--crate-name", "pages", "--edition="+edition)
		if takes == invalid(func(e *n.Environment) { e.Edition = edition }) {
			t.Errorf("edition %q: the compiler takes it: %v\n%s", edition, takes, output)
		}
	}
	for _, kind := range []string{"lib", "rlib", "staticlib", "dylib", "cdylib", "bin", "proc-macro", "sdylib", "Lib", "banana", ""} {
		takes, output := rsRustc(t, append(base, "--crate-name", "pages", "--crate-type="+kind)...)
		if takes == invalid(func(e *n.Environment) { e.Crate.Kind = kind }) {
			t.Errorf("crate kind %q: the compiler takes it: %v\n%s", kind, takes, output)
		}
	}

	// Every target the compiler has is a well-formed target text.
	targets, err := exec.Command(filepath.Join(rsCompiler(t).SysrootPath, "bin", "rustc"), "--print", "target-list").Output()
	if names := strings.Fields(string(targets)); err != nil || len(names) < 100 || !slices.Contains(names, n.TargetLinuxX64) {
		t.Fatal("the compiler's targets", err, len(names))
	} else {
		for _, target := range names {
			if invalid(func(e *n.Environment) { e.Target = target }) {
				t.Errorf("the compiler's target %q is no well-formed target here", target)
			}
		}
	}

	// A captured crate named std, core or alloc is unsupported, and one named
	// like any other crate of the sysroot is admitted.
	crate := regexp.MustCompile(`^lib([A-Za-z0-9_]+)-[0-9a-f]+\.(rlib|rmeta)$`)
	entries, err := os.ReadDir(filepath.Join(rsCompiler(t).SysrootPath, "lib", "rustlib", n.TargetLinuxX64, "lib"))
	if err != nil {
		t.Fatal(err)
	}
	present := map[string]bool{}
	for _, entry := range entries {
		if match := crate.FindStringSubmatch(entry.Name()); match != nil {
			present[match[1]] = true
		}
	}
	if !present["std"] || !present["core"] || !present["alloc"] || !present["libc"] || len(present) < 19 {
		t.Fatal("the crates of the sysroot", len(present))
	}
	for name := range present {
		environment := rsEnvironment()
		environment.Crate.Name = name
		_, err := n.ConfigurationFor(environment)
		if refused := name == "std" || name == "core" || name == "alloc"; refused != errors.Is(err, n.ErrUnsupported) || !refused && err != nil {
			t.Errorf("crate name %s of the sysroot: %v", name, err)
		}
	}
	// The reason for the three: the graph names an entity by its crate and its
	// path. Under the name of another crate that the sysroot loads, the crate
	// root is the one entity of its name and a lookup finds it.
	for _, name := range []string{"addr2line", "adler2", "cfg_if", "compiler_builtins", "gimli", "hashbrown", "libc", "memchr", "miniz_oxide", "object", "panic_unwind", "rustc_demangle", "rustc_std_workspace_alloc", "rustc_std_workspace_core", "std_detect", "unwind"} {
		environment := rsEnvironment()
		environment.Crate.Name = name
		snapshot, _, _ := rsCapture(t, `pub struct A { pub o: Option<u8>, pub s: String, pub m: ::std::collections::HashMap<u8, u8> }`)
		_, result := rsResultsFor(t, snapshot, rsSettingsFor(t, environment))
		m, err := n.Open(result)
		if err != nil {
			t.Fatal(name, err, result.Diagnostics())
		}
		if root, err := m.Lookup(name); err != nil || root.Entity().Kind != n.ModuleKind() {
			t.Errorf("the crate root of a crate named %s: %v", name, err)
		}
		if _, err := m.Lookup(name + "::A"); err != nil {
			t.Errorf("a declaration of a crate named %s: %v", name, err)
		}
	}
}
