package integration_test

import (
	"context"
	"errors"
	"go/types"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	analyze "github.com/Bitspark/bittype/analyze/files/horos/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	gp "github.com/Bitspark/bittype/provider/golang/go"
	hp "github.com/Bitspark/bittype/provider/horos/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	hs "github.com/Bitspark/bittype/syntax/horos/go"
	generate "github.com/Bitspark/bittype/transforms/horos-to-go/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	n "github.com/Bitspark/bittype/verification/files/target/golang/go"
	decl "github.com/Bitspark/horos/core/go/declarations"
	"github.com/Bitspark/logos-contract/contract/go/typed"
	core "github.com/bitspark/ontos/core/go"
	data "github.com/bitspark/ontos/data/go"
)

type syntaxProvider[L any] interface {
	language.Provider
	Check(ast.View[L]) error
	Parse(context.Context, parse.Input) (cst.CST[L], error)
}

func capture[L any](t *testing.T, p syntaxProvider[L], profile, text string) semantic.Snapshot {
	t.Helper()
	ctx := context.Background()
	definition, err := representation.NewLanguage[L](p, profile)
	if err != nil {
		t.Fatal(err)
	}
	c, err := check.NewChecker[L](definition, p)
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := p.Parse(ctx, parse.Input{SourceID: "types", TreeID: "same-address", Text: text})
	if err != nil {
		t.Fatal(err)
	}
	tree := parsed.ASTView()
	envelope, err := representation.Pack[L](definition, tree)
	if err != nil {
		t.Fatal(err)
	}
	s, err := semantic.Capture(ctx, c, semantic.Input{Entry: ast.NodeRef{Tree: tree.ID, Node: tree.Root}, Trees: []representation.AST{envelope}})
	if err != nil {
		t.Fatal(err)
	}
	return s
}
func sourceModel(t *testing.T, text string) h.Model {
	t.Helper()
	text = "module example; import h from horos.kinds pin \"sha256:" + decl.PreludePin().String() + "\"; " + text
	s := capture[hs.Language](t, hp.Provider{}, hs.Profile, text)
	r, err := resolution.NewResolver[hs.Language](s.Language(), resolve.Contract(), hp.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	e, err := r.Resolve(context.Background(), s, semantic.Options{})
	if err != nil {
		t.Fatal(err)
	}
	m, err := h.Inspect(context.Background(), s, e)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
func nativeModel(t *testing.T, text string) n.Model {
	t.Helper()
	s := capture[gs.Language](t, gp.Provider{}, "go-syntax-1", "package example; "+text)
	m, err := n.Inspect(context.Background(), s)
	if err != nil {
		t.Fatal(err)
	}
	return m
}
func declaration(t *testing.T, m h.Model, name string) h.Declaration {
	t.Helper()
	d, err := m.Declaration(decl.DeclarationID{Module: m.Root(), Name: name})
	if err != nil {
		t.Fatal(err)
	}
	return d
}
func native(t *testing.T, m n.Model, name string) n.Type {
	t.Helper()
	d, err := m.Lookup(name)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
func importedKind(t *testing.T, m h.Model, name string) h.Kind {
	t.Helper()
	k, err := m.CompileKind(m.Root(), core.NewTuple(core.NewAtom([]byte("horos/decl/kind-imported")), data.MustEncodeText("h"), data.MustEncodeText(name)))
	if err != nil {
		t.Fatal(err)
	}
	return k
}
func composition(t *testing.T, policy typed.Policy, limits v.Limits) v.Composition {
	t.Helper()
	c, err := v.New(policy, limits)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func receiver(t *testing.T) v.Composition {
	t.Helper()
	limits := v.DefaultLimits()
	limits.Checker.Work = 20000000
	return composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: v.MappingDependencies()}, limits)
}
func checkClaim(t *testing.T, c v.Composition, claim v.Claim, candidate typed.Expr) v.Result {
	t.Helper()
	r, err := c.Prepare(context.Background(), claim)
	if err != nil {
		t.Fatal(err)
	}
	out, err := c.Check(context.Background(), r, candidate)
	if err != nil {
		t.Fatal(err)
	}
	return out
}
func accepted(t *testing.T, r v.Result) {
	t.Helper()
	o := r.Outcome()
	if !r.Valid() || !o.IsAccepted() {
		t.Fatalf("%s: %s %s at %+v, steps=%d", o.Status(), o.Code(), o.Detail(), o.Location(), o.Steps())
	}
}

func box(t *testing.T, constraint string) (h.Model, n.Model, v.Claim, error) {
	t.Helper()
	s := sourceModel(t, "type TextBox<T: h.Text> = { value: T; };")
	g := nativeModel(t, "type TextBox[T "+constraint+"] struct { Value T }")
	d := declaration(t, s, "TextBox")
	kind, err := d.ParameterKind(0)
	if err != nil {
		t.Fatal(err)
	}
	boolean, err := s.Proper(core.NewTuple(core.NewAtom([]byte("horos/scalar")), core.NewAtom([]byte("bool"))))
	if err != nil {
		t.Fatal(err)
	}
	nonmember, err := s.Satisfies(context.Background(), boolean, kind)
	if err != nil {
		t.Fatal(err)
	}
	if nonmember.Verdict() != decl.Refuted {
		t.Fatal("Horos did not refute bool:Text", nonmember.Verdict())
	}
	ng := native(t, g, "TextBox")
	nb, err := g.Basic(types.Bool)
	if err != nil {
		t.Fatal(err)
	}
	admission, err := g.Admit(context.Background(), ng, nb)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := v.Counterexample(d, ng, nonmember, admission)
	return s, g, claim, err
}

func TestTextBoxExplicitCounterexample(t *testing.T) {
	s, g, claim, err := box(t, "any")
	if err != nil {
		t.Fatal(err)
	}
	r := checkClaim(t, receiver(t), claim, v.CounterexampleCandidate())
	accepted(t, r)
	if !r.Request().Association().Participants()[0].Snapshot.Same(s.Snapshot()) || !r.Request().Association().Participants()[1].Snapshot.Same(g.Snapshot()) {
		t.Fatal("lost actual captures")
	}
	a, _ := r.Outcome().Accepted()
	want := []string{v.ModelBasis, v.NativeBasis, v.ScalarLaw, v.SourceBasis}
	slices.Sort(want)
	if !slices.Equal(a.Dependencies(), want) {
		t.Fatal(a.Dependencies())
	}
	if r.Request().Question().Name != "admission-counterexample" {
		t.Fatal("counterexample was mislabeled as preservation")
	}
	_, _, _, err = box(t, "interface{ string }")
	if !errors.Is(err, semantic.ErrEvidenceUnavailable) {
		t.Fatal("bool must cease to be an admission counterexample", err)
	}
	_, _, _, err = box(t, "interface{ ~string }")
	if !errors.Is(err, semantic.ErrEvidenceUnavailable) {
		t.Fatal(err)
	}
}

func TestKindSpellingDoesNotSubstituteItsResolutionModule(t *testing.T) {
	s := sourceModel(t, "kind Text = h.Bool; type Box<T: Text> = {value:T;};")
	g := nativeModel(t, "type Box[T any] struct {Value T}")
	d := declaration(t, s, "Box")
	bound, err := d.ParameterKind(0)
	if err != nil {
		t.Fatal(err)
	}
	boolean, err := s.Proper(core.NewTuple(core.NewAtom([]byte("horos/scalar")), core.NewAtom([]byte("bool"))))
	if err != nil {
		t.Fatal(err)
	}
	actual, err := s.Satisfies(context.Background(), boolean, bound)
	if err != nil || actual.Verdict() != decl.Proven {
		t.Fatal(actual.Verdict(), err)
	}
	foreignScope, err := s.CompileKind(decl.PreludePin(), core.NewTuple(core.NewAtom([]byte("horos/decl/kind-local")), data.MustEncodeText("Text")))
	if err != nil {
		t.Fatal(err)
	}
	refuted, err := s.Satisfies(context.Background(), boolean, foreignScope)
	if err != nil || refuted.Verdict() != decl.Refuted {
		t.Fatal(refuted.Verdict(), err)
	}
	if !bound.Expression().Value().Equal(foreignScope.Expression().Value()) {
		t.Fatal("fixture must use the same local spelling")
	}
	b, _ := g.Basic(types.Bool)
	admission, err := g.Admit(context.Background(), native(t, g, "Box"), b)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := v.Counterexample(d, native(t, g, "Box"), refuted, admission); !errors.Is(err, v.ErrCoverage) {
		t.Fatal("query was rebound by spelling", err)
	}
	if _, err := foreignScope.Arrow(s.Root()); !errors.Is(err, h.ErrForeign) {
		t.Fatal("arrow reinterpreted a bound kind in a foreign module", err)
	}
}

const sourcePage = "type Page<T> = { items: list<T>; next?: text; }; type Person = {name: text; active: bool;};"
const nativePage = "type BittypeOptional[V any] struct { Present bool; Value V }; type Page[N any] struct { Items []N; Next BittypeOptional[string] }; type Person struct { Name string; Active bool }"

type pagination struct {
	source   h.Model
	target   n.Model
	family   v.Family
	kind     h.Kind
	argument v.Argument
}

func page(t *testing.T, sourceText, nativeText string) (pagination, error) {
	t.Helper()
	return bindPage(t, sourceModel(t, sourceText), nativeModel(t, nativeText))
}
func bindPage(t *testing.T, s h.Model, g n.Model) (pagination, error) {
	t.Helper()
	d := declaration(t, s, "Page")
	k := importedKind(t, s, "Data")
	constructor, err := d.Value()
	if err != nil {
		t.Fatal(err)
	}
	arrow, err := k.Arrow(s.Root())
	if err != nil {
		t.Fatal(err)
	}
	guarantee, err := s.Satisfies(context.Background(), constructor, arrow)
	if err != nil {
		t.Fatal(err)
	}
	f, err := v.BindFamily(d, native(t, g, "Page"), []v.FieldBinding{{Source: "items", Target: "Items"}, {Source: "next", Target: "Next"}}, v.Optional{Constructor: native(t, g, "BittypeOptional"), PresenceField: "Present", ValueField: "Value"}, k, guarantee)
	if err != nil {
		return pagination{}, err
	}
	p := declaration(t, s, "Person")
	pv, err := p.Value()
	if err != nil {
		t.Fatal(err)
	}
	member, err := s.Satisfies(context.Background(), pv, k)
	if err != nil {
		t.Fatal(err)
	}
	a, err := f.BindArgument(p, native(t, g, "Person"), []v.FieldBinding{{Source: "name", Target: "Name"}, {Source: "active", Target: "Active"}}, member)
	return pagination{s, g, f, k, a}, err
}

func TestGeneratedPageUsesActualTransformerOutput(t *testing.T) {
	s := sourceModel(t, sourcePage+" type PeoplePage = Page<Person>; type NestedPage = Page<Page<Person>>;")
	ctx := context.Background()
	from := s.Snapshot().Language()
	a, err := analysis.NewAnalyzer[hs.Language](from, analyze.Contract(), hp.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	receipt := s.Receipt()
	evidence, err := a.Analyze(ctx, analysis.Input{Snapshot: s.Snapshot(), Resolution: &receipt})
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[gs.Language](gp.Provider{}, "go-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	mapper, err := generate.New(from, to, generate.Config{TreeID: "generated", Package: "example", Profile: generate.RepresentationProfile})
	if err != nil {
		t.Fatal(err)
	}
	out, err := mapper.Transform(ctx, transform.Input{Source: s.Snapshot(), Evidence: []semantic.Result{evidence}})
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[gs.Language](to, gp.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	snapshot, err := semantic.Capture(ctx, checker, semantic.Input{Entry: out.Entry, Trees: out.Trees})
	if err != nil {
		t.Fatal(err)
	}
	g, err := n.Inspect(ctx, snapshot)
	if err != nil {
		t.Fatal(err)
	}
	if !native(t, g, "NestedPage").IsAlias() {
		t.Fatal("missing actual nested Go instantiation")
	}
	p, err := bindPage(t, s, g)
	if err != nil {
		t.Fatal(err)
	}
	claim, _ := p.family.Theorem()
	c := receiver(t)
	theorem := checkClaim(t, c, claim, p.family.Candidate())
	accepted(t, theorem)
	application, err := p.family.Application(theorem, p.argument, 2)
	if err != nil {
		t.Fatal(err)
	}
	r := checkClaim(t, c, application, p.argument.Candidate(2))
	accepted(t, r)
	if !r.Request().Association().Participants()[1].Snapshot.Same(snapshot) {
		t.Fatal("verification substituted a hand-authored target")
	}
}

func TestOrdinaryRenamedFamilyUsesSameMapping(t *testing.T) {
	s := sourceModel(t, "kind Payload = h.Data; type Envelope<Element> = { entries: list<Element>; cursor?: text; };")
	g := nativeModel(t, "type Maybe[X any] struct { Has bool; Item X }; type Envelope[Value any] struct { Entries []Value; Cursor Maybe[string] }")
	d := declaration(t, s, "Envelope")
	k, err := s.CompileKind(s.Root(), core.NewTuple(core.NewAtom([]byte("horos/decl/kind-local")), data.MustEncodeText("Payload")))
	if err != nil {
		t.Fatal(err)
	}
	arrow, err := k.Arrow(s.Root())
	if err != nil {
		t.Fatal(err)
	}
	value, err := d.Value()
	if err != nil {
		t.Fatal(err)
	}
	guarantee, err := s.Satisfies(context.Background(), value, arrow)
	if err != nil {
		t.Fatal(err)
	}
	f, err := v.BindFamily(d, native(t, g, "Envelope"), []v.FieldBinding{{Source: "entries", Target: "Entries"}, {Source: "cursor", Target: "Cursor"}}, v.Optional{Constructor: native(t, g, "Maybe"), PresenceField: "Has", ValueField: "Item"}, k, guarantee)
	if err != nil {
		t.Fatal(err)
	}
	claim, _ := f.Theorem()
	accepted(t, checkClaim(t, receiver(t), claim, f.Candidate()))
}
func TestPageClosedTheoremAndNestedReuse(t *testing.T) {
	p, err := page(t, sourcePage, nativePage)
	if err != nil {
		t.Fatal(err)
	}
	claim, err := p.family.Theorem()
	if err != nil {
		t.Fatal(err)
	}
	c := receiver(t)
	theorem := checkClaim(t, c, claim, p.family.Candidate())
	accepted(t, theorem)
	if len(theorem.Request().Statement().Context) != 0 || theorem.Request().Statement().Expected.Kind() != typed.PiKind {
		t.Fatal("not a closed quantified theorem")
	}
	acceptance, _ := theorem.Outcome().Accepted()
	want := v.MappingDependencies()
	want = slices.DeleteFunc(want, func(s string) bool { return s == "" })
	slices.Sort(want)
	if !slices.Equal(acceptance.Dependencies(), want) {
		t.Fatalf("dependencies=%v want %v", acceptance.Dependencies(), want)
	}
	for _, depth := range []int{1, 2} {
		application, err := p.family.Application(theorem, p.argument, depth)
		if err != nil {
			t.Fatal(err)
		}
		result := checkClaim(t, c, application, p.argument.Candidate(depth))
		accepted(t, result)
		found := false
		for _, d := range result.Request().Statement().Declarations {
			if d.Name() == "family.theorem" {
				_, defined := d.Definition()
				found = defined && !d.IsAssumption()
			}
		}
		if !found {
			t.Fatal("accepted theorem was not retained as a checked definition")
		}
		wrong := typed.Fst(p.argument.Candidate(depth))
		invalid := checkClaim(t, c, application, wrong)
		if invalid.Outcome().Status() != typed.InvalidStatus {
			t.Fatal("membership alone supplied representation", invalid.Outcome().Status())
		}
	}
	application, err := p.family.Application(theorem, p.argument, 2)
	if err != nil {
		t.Fatal(err)
	}
	dependencies := slices.DeleteFunc(v.MappingDependencies(), func(dependency string) bool { return dependency == v.OptionalLaw })
	limits := v.DefaultLimits()
	limits.Checker.Work = 20000000
	restricted := composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: dependencies}, limits)
	refused := checkClaim(t, restricted, application, p.argument.Candidate(2))
	if refused.Outcome().Status() != typed.UnsupportedStatus || refused.Outcome().Code() != typed.UnadmittedDependency {
		t.Fatal("reused theorem bypassed the current receiver policy", refused.Outcome().Status(), refused.Outcome().Code())
	}
}

func TestChangedHelperAndOriginalConstraintCannotReuseOldBasis(t *testing.T) {
	_, err := page(t, sourcePage, strings.Replace(nativePage, "Present bool", "Present string", 1))
	if !errors.Is(err, v.ErrCoverage) {
		t.Fatal("changed helper passed old mapping", err)
	}
	_, err = page(t, sourcePage, strings.Replace(nativePage, "Items []N", "Items []bool", 1))
	if !errors.Is(err, v.ErrCoverage) {
		t.Fatal("changed field passed old mapping", err)
	}
	_, err = page(t, strings.Replace(sourcePage, "Page<T>", "Page<T: h.Text>", 1), nativePage)
	if !errors.Is(err, semantic.ErrEvidenceUnavailable) {
		t.Fatal("original bound was discarded", err)
	}
}

func TestFreshCapturesAndForeignProjectionHandles(t *testing.T) {
	p, err := page(t, sourcePage, nativePage)
	if err != nil {
		t.Fatal(err)
	}
	q, err := page(t, sourcePage, nativePage)
	if err != nil {
		t.Fatal(err)
	}
	if p.source.Snapshot().Entry() != q.source.Snapshot().Entry() {
		t.Fatal("fixture must exercise equal addresses")
	}
	if _, err := h.Inspect(context.Background(), p.source.Snapshot(), q.source.Receipt()); !errors.Is(err, semantic.ErrInvalidRequest) {
		t.Fatal(err)
	}
	d := declaration(t, p.source, "Page")
	value, err := d.Value()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := p.source.Satisfies(context.Background(), value, q.kind); !errors.Is(err, h.ErrForeign) {
		t.Fatal(err)
	}
	claim, _ := p.family.Theorem()
	theorem := checkClaim(t, receiver(t), claim, p.family.Candidate())
	accepted(t, theorem)
	if _, err := q.family.Application(theorem, q.argument, 1); !errors.Is(err, v.ErrInvalid) {
		t.Fatal("retagged theorem", err)
	}
	fresh, _ := q.family.Theorem()
	result := checkClaim(t, receiver(t), fresh, p.family.Candidate())
	accepted(t, result)
	if result.Request().Association().Same(theorem.Request().Association()) {
		t.Fatal("proof syntax reuse reused authority")
	}
	copy := d.Definition()
	*copy = decl.Declaration{}
	if d.Definition().ID().Name != "Page" {
		t.Fatal("definition view mutated retained handle")
	}
}

func TestJointOutcomesAndImmutability(t *testing.T) {
	_, _, claim, err := box(t, "any")
	if err != nil {
		t.Fatal(err)
	}
	policy := typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: v.MappingDependencies()}
	limits := v.DefaultLimits()
	c := composition(t, policy, limits)
	request, err := c.Prepare(context.Background(), claim)
	if err != nil {
		t.Fatal(err)
	}
	policy.Dependencies[0] = "changed"
	statement := request.Statement()
	statement.Declarations[0] = typed.Assume("broken", typed.Bool())
	statement.Context = append(statement.Context, typed.Bool())
	view := request.Policy()
	view.Dependencies[0] = "changed"
	origins := request.Origins()
	origins[0].Sites[0].Node = ast.NodeRef{}
	result, err := c.Check(context.Background(), request, v.CounterexampleCandidate())
	if err != nil {
		t.Fatal(err)
	}
	accepted(t, result)
	retained := result.Request()
	request = v.Request{}
	resultView := result.Outcome()
	a, _ := resultView.Accepted()
	a = typed.Acceptance{}
	if a.Valid() || !retained.Valid() || !result.Request().Same(retained) {
		t.Fatal("mutable outer handle")
	}
	if _, err := receiver(t).Check(context.Background(), retained, v.CounterexampleCandidate()); !errors.Is(err, v.ErrInvalid) {
		t.Fatal("foreign composition request", err)
	}
	if _, err := c.Check(context.Background(), v.Request{}, v.CounterexampleCandidate()); !errors.Is(err, v.ErrInvalid) {
		t.Fatal(err)
	}
	invalid := checkClaim(t, c, claim, typed.True())
	if invalid.Outcome().Status() != typed.InvalidStatus || invalid.Diagnostics()[0].Location.Phase != typed.CandidatePhase || len(invalid.Diagnostics()[0].Sites) != 0 {
		t.Fatal("invented candidate location")
	}
	unsupported := checkClaim(t, composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker}, limits), claim, v.CounterexampleCandidate())
	if unsupported.Outcome().Status() != typed.UnsupportedStatus {
		t.Fatal(unsupported.Outcome().Status())
	}
	limits.Checker.Work = 1
	incomplete := checkClaim(t, composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: v.MappingDependencies()}, limits), claim, v.CounterexampleCandidate())
	if incomplete.Outcome().Status() != typed.IncompleteStatus {
		t.Fatal(incomplete.Outcome().Status())
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if r, err := c.Prepare(ctx, claim); !errors.Is(err, context.Canceled) || r.Valid() {
		t.Fatal(err)
	}
	if r, err := c.Check(ctx, retained, v.CounterexampleCandidate()); !errors.Is(err, context.Canceled) || r.Valid() {
		t.Fatal(err)
	}
}

func TestJointConcurrentChecks(t *testing.T) {
	_, _, claim, err := box(t, "any")
	if err != nil {
		t.Fatal(err)
	}
	c := receiver(t)
	r, err := c.Prepare(context.Background(), claim)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out, err := c.Check(context.Background(), r, v.CounterexampleCandidate())
			if err != nil || !out.Outcome().IsAccepted() {
				t.Error(err, out.Outcome().Status())
			}
			if !reflect.DeepEqual(out.Request().Statement(), r.Statement()) {
				t.Error("statement mutated")
			}
		}()
	}
	wg.Wait()
}
