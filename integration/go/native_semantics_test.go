package integration_test

import (
	"context"
	"errors"
	"fmt"
	"go/types"
	"reflect"
	"sync"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/golang/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/engine/go"
	p "github.com/Bitspark/bittype/provider/golang/go"
	"github.com/Bitspark/bittype/semantics/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

func goCapture(t *testing.T, texts ...string) (semantic.Snapshot, []s.AST) {
	t.Helper()
	var trees []s.AST
	for i, text := range texts {
		parsed, err := (p.Provider{}).Parse(context.Background(), parse.Input{SourceID: ast.SourceID(fmt.Sprint(i)), TreeID: ast.TreeID(fmt.Sprint(i)), Text: "package example; " + text})
		if err != nil {
			t.Fatal(err)
		}
		trees = append(trees, parsed.ASTView())
	}
	return goCaptureTrees(t, trees), trees
}
func goCaptureTrees(t *testing.T, trees []s.AST) semantic.Snapshot {
	t.Helper()
	definition, err := representation.NewLanguage[s.Language](s.Language{}, "go-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	checker, err := p.NewSyntaxChecker(definition)
	if err != nil {
		t.Fatal(err)
	}
	var envelopes []representation.AST
	for _, tree := range trees {
		envelope, err := representation.Pack[s.Language](definition, tree)
		if err != nil {
			t.Fatal(err)
		}
		envelopes = append(envelopes, envelope)
	}
	snapshot, err := semantic.Capture(context.Background(), checker, semantic.Input{Entry: ast.NodeRef{Tree: trees[0].ID, Node: trees[0].Root}, Trees: envelopes})
	if err != nil {
		t.Fatal(err)
	}
	return snapshot
}
func goResults(t *testing.T, snapshot semantic.Snapshot) (semantic.Result, semantic.Result) {
	t.Helper()
	ctx := context.Background()
	resolver, err := p.NewFilesResolver(snapshot.Language())
	if err != nil {
		t.Fatal(err)
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	analyzer, err := p.NewFilesAnalyzer(snapshot.Language())
	if err != nil {
		t.Fatal(err)
	}
	result, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}})
	if err != nil {
		t.Fatal(err)
	}
	return resolved, result
}
func goModel(t *testing.T, text string) n.Model {
	t.Helper()
	snapshot, _ := goCapture(t, text)
	_, result := goResults(t, snapshot)
	model, err := n.Open(result)
	if err != nil {
		t.Fatalf("model: %v; diagnostics: %+v", err, result.Diagnostics())
	}
	return model
}
func goType(t *testing.T, m n.Model, name string) n.Type {
	t.Helper()
	ty, err := m.Lookup(name)
	if err != nil {
		t.Fatal(err)
	}
	return ty
}

// Called by the shared semantic conformance matrix, independently of text cases.
func qualifyGolang(t *testing.T) {
	ctx := context.Background()
	snapshot, _ := goCapture(t, "type Text string; type Alias = Text; type Box[T any] struct { Value T }; type Concrete = Box[Text]")
	var registry engine.Engine
	if err := registry.RegisterLanguage(snapshot.Language()); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Resolve(ctx, snapshot, semantic.Options{}); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	resolver, err := p.NewFilesResolver(snapshot.Language())
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterResolver(resolver); err != nil {
		t.Fatal(err)
	}
	resolved, err := registry.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	input := analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}}
	if _, err := registry.Analyze(ctx, input); !errors.Is(err, engine.ErrMissingCapability) {
		t.Fatal(err)
	}
	analyzer, err := p.NewFilesAnalyzer(snapshot.Language())
	if err != nil {
		t.Fatal(err)
	}
	if err := registry.RegisterAnalyzer(analyzer); err != nil {
		t.Fatal(err)
	}
	result, err := registry.Analyze(ctx, input)
	if err != nil {
		t.Fatal(err)
	}
	rg, _ := resolved.Semantics()
	ag, _ := result.Semantics()
	if !rg.Same(ag) {
		t.Fatal("analysis rebuilt the graph")
	}
	m, err := n.Open(result)
	if err != nil {
		t.Fatal(err, result.Diagnostics())
	}
	if m.Graph().Data().Entities[0].Kind != n.Key("package") || len(m.Graph().Data().Bindings) == 0 || len(m.Graph().Data().Expressions) == 0 {
		t.Fatal("missing native facts")
	}
	text, alias, box := goType(t, m, "Text"), goType(t, m, "Alias"), goType(t, m, "Box")
	if !text.IsNamed() || !alias.IsAlias() || !text.Same(alias) || !alias.AliasTarget().Same(text) {
		t.Fatal("alias/definition identity")
	}
	if len(box.Parameters()) != 1 || !box.Parameters()[0].Unconstrained() {
		t.Fatal("missing constraint")
	}
	fields, ok := box.Fields()
	if !ok || len(fields) != 1 || !fields[0].Type.Same(box.Parameters()[0]) || fields[0].Site.Tree == "" {
		t.Fatal("missing generic fields")
	}
	concrete := goType(t, m, "Concrete").AliasTarget()
	if !concrete.Origin().Same(box) || len(concrete.Arguments()) != 1 || !concrete.Arguments()[0].Same(text) {
		t.Fatal("missing application")
	}
	fields, _ = concrete.Fields()
	if !fields[0].Type.Same(text) {
		t.Fatal("instantiated fields were not substituted")
	}
	if bindings, ok := m.Graph().DirectParameterBindings(concrete.ID()); !ok || len(bindings) != 1 {
		t.Fatal("application does not use shared parameter machinery")
	}
	data := m.Graph().Data()
	data.Entities[0].Name = "mutated"
	fields[0].Name = "mutated"
	if m.Graph().Data().Entities[0].Name == "mutated" {
		t.Fatal("mutable graph escaped")
	}
}

// bittype:case target=github.com/Bitspark/bittype/resolve/files/golang/go#Provider.Resolve feature=go-semantic-subset case=rejected-and-unsupported-fixtures expects=no-partial-model
// bittype:case target=github.com/Bitspark/bittype/analyze/files/golang/go#Provider.Analyze feature=go-semantic-subset case=rejected-and-unsupported-fixtures expects=no-partial-model
func TestGoSemanticReportsAndCoverage(t *testing.T) {
	for _, tc := range []struct {
		name, text string
		coverage   semantic.Completion
		validity   semantic.Verdict
	}{
		{"unresolved", "type Box Missing", semantic.Complete, semantic.Rejected},
		{"constraint", "type Box[T ~string] struct { Value T }; type Bad = Box[bool]", semantic.Complete, semantic.Rejected},
		{"duplicate", "type Box struct { Value string; Value bool }", semantic.Complete, semantic.Rejected},
		{"pointer", "type Box *string", semantic.Unsupported, semantic.Undetermined},
		{"map", "type Box map[string]bool", semantic.Unsupported, semantic.Undetermined},
		{"array", "type Box [3]string", semantic.Unsupported, semantic.Undetermined},
		{"method", "type Box interface { Value() string }", semantic.Unsupported, semantic.Undetermined},
		{"embedded", "type Text string; type Box struct { Text }", semantic.Unsupported, semantic.Undetermined},
		{"import", "import x \"example.invalid/x\"; type Box string", semantic.Unsupported, semantic.Undetermined},
	} {
		t.Run(tc.name, func(t *testing.T) {
			snapshot, _ := goCapture(t, tc.text)
			resolved, result := goResults(t, snapshot)
			if resolved.HasSemantics() || result.HasSemantics() {
				t.Fatal("partial model published")
			}
			coverage, _ := result.Report(n.Coverage())
			validity, _ := result.Report(n.Validity())
			if coverage.Status != tc.coverage || validity.Verdict != tc.validity {
				t.Fatal(coverage, validity)
			}
			if ds := result.Diagnostics(); len(ds) == 0 || ds[0].At == nil || ds[0].At.Tree == "" {
				t.Fatal("missing captured diagnostic", ds)
			}
			if model, err := n.Open(result); err == nil || model.Valid() {
				t.Fatal("failed result became a native model")
			}
		})
	}
}
func TestGoSemanticReceiptAndEnvironmentIsolation(t *testing.T) {
	ctx := context.Background()
	snapshot, trees := goCapture(t, "type Box[T any] struct {Value T}")
	foreign := goCaptureTrees(t, trees)
	resolved, result := goResults(t, snapshot)
	analyzer, _ := p.NewFilesAnalyzer(snapshot.Language())
	resolver, _ := p.NewFilesResolver(snapshot.Language())
	for name, input := range map[string]analysis.Input{
		"missing":    {Snapshot: snapshot, Options: semantic.Options{Configuration: n.Configuration()}},
		"foreign":    {Snapshot: foreign, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}},
		"wrong-role": {Snapshot: snapshot, Resolution: &result, Options: semantic.Options{Configuration: n.Configuration()}},
	} {
		t.Run(name, func(t *testing.T) {
			if r, err := analyzer.Analyze(ctx, input); err == nil || r.Valid() {
				t.Fatal("bad receipt accepted")
			}
		})
	}
	if _, err := n.Open(resolved); err == nil {
		t.Fatal("resolution receipt mistaken for analysis")
	}
	for i := range n.Configuration() {
		settings := n.Configuration()
		settings[i].Value = []byte("foreign")
		if r, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: settings}); err == nil || r.Valid() {
			t.Fatal("foreign environment accepted", i)
		}
	}
	if _, err := resolver.Resolve(ctx, snapshot, semantic.Options{}); err == nil {
		t.Fatal("implicit environment")
	}
	key := n.Key("unknown-query")
	r, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration(), Requested: []semantics.Key{key}})
	if err != nil {
		t.Fatal(err)
	}
	report, _ := r.Report(key)
	if report.Status != semantic.Unsupported {
		t.Fatal(report)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := resolver.Resolve(canceled, snapshot, semantic.Options{Configuration: n.Configuration()}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	var texts []string
	for i := 0; i < 17; i++ {
		texts = append(texts, fmt.Sprintf("type X%d string", i))
	}
	large, _ := goCapture(t, texts...)
	limited, err := resolver.Resolve(ctx, large, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		t.Fatal(err)
	}
	coverage, _ := limited.Report(n.Coverage())
	if coverage.Status != semantic.Incomplete || limited.HasSemantics() {
		t.Fatal("native bound lost")
	}
}
func TestGoSemanticMultiUnitReferencesAndRecursion(t *testing.T) {
	snapshot, trees := goCapture(t, "type List []Person; type Recursive struct { Children []Recursive }", "type Person struct {Name string; Problem error}")
	_, result := goResults(t, snapshot)
	m, err := n.Open(result)
	if err != nil {
		t.Fatal(err, result.Diagnostics())
	}
	list := goType(t, m, "List")
	element, ok := list.Underlying().SliceElement()
	if !ok || !element.Same(goType(t, m, "Person")) {
		t.Fatal("cross-file lookup")
	}
	query, err := analyze.NewQueries(context.Background(), m)
	if err != nil {
		t.Fatal("reconstruction", err)
	}
	_ = query
	// A caller-supplied reference may assert only the target resolved by Go.
	for i, node := range trees[0].Nodes {
		if node.Kind == s.KindNamedType {
			trees[0].Nodes[i].References = []ast.Reference{{Role: s.RoleTarget, Target: ast.NodeRef{Tree: trees[0].ID, Node: trees[0].Root}}}
			break
		}
	}
	forged := goCaptureTrees(t, trees)
	resolver, _ := p.NewFilesResolver(forged.Language())
	if _, err := resolver.Resolve(context.Background(), forged, semantic.Options{Configuration: n.Configuration()}); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal("forged reference", err)
	}
}
func TestGoSemanticNativeAdmission(t *testing.T) {
	ctx := context.Background()
	m := goModel(t, "type Text string; type Open[T any] struct{ Value T }; type Strict[T string] struct{ Value T }; type Approx[T ~string] struct{ Value T }")
	queries, err := analyze.NewQueries(ctx, m)
	if err != nil {
		t.Fatal(err)
	}
	text, _ := m.Basic(types.String)
	boolean, _ := m.Basic(types.Bool)
	for _, tc := range []struct {
		name string
		arg  n.Type
		want bool
	}{{"Open", boolean, true}, {"Strict", text, true}, {"Strict", goType(t, m, "Text"), false}, {"Approx", goType(t, m, "Text"), true}, {"Approx", boolean, false}} {
		a, err := queries.Admit(ctx, goType(t, m, tc.name), tc.arg)
		if err != nil || a.Admitted() != tc.want || !a.Model().Same(m) {
			t.Fatal(tc.name, a.Detail(), err)
		}
		args := a.Arguments()
		args[0] = n.Type{}
		if !a.Arguments()[0].Valid() {
			t.Fatal("mutable admission")
		}
	}
	foreign := goModel(t, "type Open[T any] struct{Value T}")
	foreignBool, _ := foreign.Basic(types.Bool)
	if _, err := queries.Admit(ctx, goType(t, m, "Open"), foreignBool); !errors.Is(err, n.ErrForeign) {
		t.Fatal("foreign argument", err)
	}
	if _, err := queries.Admit(ctx, goType(t, foreign, "Open"), boolean); !errors.Is(err, n.ErrForeign) {
		t.Fatal("foreign constructor", err)
	}
	if _, err := queries.Admit(ctx, n.Type{}, boolean); !errors.Is(err, n.ErrForeign) {
		t.Fatal("zero constructor", err)
	}
	if _, err := queries.Admit(ctx, goType(t, m, "Open")); !errors.Is(err, n.ErrUnsupported) {
		t.Fatal("arity", err)
	}
	if _, err := analyze.NewQueries(ctx, n.Model{}); !errors.Is(err, n.ErrForeign) {
		t.Fatal("zero model", err)
	}
	canceled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := queries.Admit(canceled, goType(t, m, "Open"), boolean); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
	if _, err := analyze.NewQueries(canceled, m); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
func TestGoSemanticConcurrentQueries(t *testing.T) {
	m := goModel(t, "type Box[T any] struct{ Value T }")
	q, err := analyze.NewQueries(context.Background(), m)
	if err != nil {
		t.Fatal(err)
	}
	box := goType(t, m, "Box")
	b, _ := m.Basic(types.Bool)
	original := m.Graph().Data()
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 8; j++ {
				a, err := q.Admit(context.Background(), box, b)
				if err != nil || !a.Admitted() {
					t.Error(err)
				}
				_, _ = box.Fields()
				if !reflect.DeepEqual(original, m.Graph().Data()) {
					t.Error("query changed semantics")
				}
			}
		}()
	}
	wg.Wait()
}
