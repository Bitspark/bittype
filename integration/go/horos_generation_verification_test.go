package integration_test

import (
	"context"
	"testing"

	analyzegolang "github.com/Bitspark/bittype/analyze/files/golang/provider/go"
	analyze "github.com/Bitspark/bittype/analyze/files/horos/provider/go"
	checkgolang "github.com/Bitspark/bittype/check/syntax/golang/provider/go"
	"github.com/Bitspark/bittype/contracts/analysis/go"
	"github.com/Bitspark/bittype/contracts/check/go"
	"github.com/Bitspark/bittype/contracts/language/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/resolution/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	genfixture "github.com/Bitspark/bittype/integration/go/testdata/horosgeneration"
	resolvegolang "github.com/Bitspark/bittype/resolve/files/golang/provider/go"
	resolve "github.com/Bitspark/bittype/resolve/files/horos/provider/go"
	n "github.com/Bitspark/bittype/semantics/golang/go"
	"github.com/Bitspark/bittype/semantics/horos/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	hs "github.com/Bitspark/bittype/syntax/horos/go"
	target "github.com/Bitspark/bittype/transform/syntax/golang/target/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/cst/go"
	"github.com/Bitspark/bittype/trees/representation/go"
	v "github.com/Bitspark/bittype/verification/files/horos-to-go/go"
	h "github.com/Bitspark/bittype/verification/files/source/horos/go"
	joint "github.com/Bitspark/bittype/verification/logos/go"
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
	s := capture[hs.Language](t, horosText(), hs.Profile, text)
	r, err := resolution.NewResolver[hs.Language](s.Language(), horos.ResolutionContract(), resolve.Provider{})
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
func composition(t *testing.T, policy typed.Policy, limits joint.Limits, authorities ...joint.Authority) joint.Composition {
	t.Helper()
	c, err := joint.New(policy, limits, authorities...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
func receiver(t *testing.T, authorities ...joint.Authority) joint.Composition {
	t.Helper()
	limits := joint.DefaultLimits()
	limits.Checker.Work = 20000000
	return composition(t, typed.Policy{Foundation: typed.Foundation, Checker: typed.Checker, Dependencies: allowedLaws()}, limits, authorities...)
}
func checkClaim(t *testing.T, c joint.Composition, claim joint.Claim, candidate typed.Expr) joint.Result {
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
func accepted(t *testing.T, r joint.Result) {
	t.Helper()
	o := r.Outcome()
	if !r.Valid() || !o.IsAccepted() {
		t.Fatalf("%s: %s %s at %+v, steps=%d", o.Status(), o.Code(), o.Detail(), o.Location(), o.Steps())
	}
}

const sourcePage = "type Page<T> = { items: list<T>; next?: text; }; type Person = {name: text; active: bool;};"

type pagination struct {
	authorities []joint.Authority
	source      h.Model
	target      n.Model
	family      v.Family
	kind        h.Kind
	argument    v.Argument
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
	rules := []v.Rule{textLaw(), booleanLaw(), listLaw(), recordLaw(), optionalLaw(native(t, g, "BittypeOptional"), "Present", "Value")}
	library := v.NewFamilies(selectMappings(t, rules...))
	f, err := library.Bind(d, native(t, g, "Page"), []v.FieldBinding{{Source: "items", Target: "Items"}, {Source: "next", Target: "Next"}}, k, guarantee)
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
	return pagination{library.Authorities(), s, g, f, k, a}, err
}

func TestGeneratedPageUsesActualTransformerOutput(t *testing.T) {
	s := sourceModel(t, sourcePage+" type PeoplePage = Page<Person>; type NestedPage = Page<Page<Person>>;")
	ctx := context.Background()
	from := s.Snapshot().Language()
	a, err := analysis.NewAnalyzer[hs.Language](from, horos.AnalysisContract(), analyze.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	receipt := s.Receipt()
	evidence, err := a.Analyze(ctx, analysis.Input{Snapshot: s.Snapshot(), Resolution: &receipt})
	if err != nil {
		t.Fatal(err)
	}
	to, err := representation.NewLanguage[gs.Language](gs.Language{}, "go-syntax-1")
	if err != nil {
		t.Fatal(err)
	}
	var registry c.Registry
	request, err := genfixture.Compose(&registry, transform.Input{Source: s.Snapshot(), Evidence: []semantic.Result{evidence}}, to, "generated", target.Provider{Package: "example"}, golangBundles()...)
	if err != nil {
		t.Fatal(err)
	}
	checker, err := check.NewChecker[gs.Language](to, checkgolang.Provider{})
	if err != nil {
		t.Fatal(err)
	}
	out, err := registry.Generate(ctx, request, checker)
	if err != nil {
		t.Fatal(err)
	}
	snapshot := out.Snapshot()
	g, err := checkedNative(ctx, snapshot)
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
	c := receiver(t, p.authorities...)
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

// Each semantic role is explicit; opening the typed view compiles nothing.
func checkedNative(ctx context.Context, snapshot semantic.Snapshot) (n.Model, error) {
	resolver, err := resolution.NewResolver[gs.Language](snapshot.Language(), n.ResolutionContract(), resolvegolang.Provider{})
	if err != nil {
		return n.Model{}, err
	}
	resolved, err := resolver.Resolve(ctx, snapshot, semantic.Options{Configuration: n.Configuration()})
	if err != nil {
		return n.Model{}, err
	}
	analyzer, err := analysis.NewAnalyzer[gs.Language](snapshot.Language(), n.AnalysisContract(), analyzegolang.Provider{})
	if err != nil {
		return n.Model{}, err
	}
	analyzed, err := analyzer.Analyze(ctx, analysis.Input{Snapshot: snapshot, Resolution: &resolved, Options: semantic.Options{Configuration: n.Configuration()}})
	if err != nil {
		return n.Model{}, err
	}
	return n.Open(analyzed)
}

// The receiver owns this policy independently of the mapping selection.
func allowedLaws() []string {
	return []string{v.ModelBasis, v.SourceBasis, v.NativeBasis, scalarLawName, listLawName, optionalLawName, recordLawName, v.NamedLaw}
}
func selectMappings(t *testing.T, rules ...v.Rule) v.Mappings {
	t.Helper()
	selected, err := v.Select(rules...)
	if err != nil {
		t.Fatal(err)
	}
	return selected
}
