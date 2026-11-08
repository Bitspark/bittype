package integration_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	"github.com/Bitspark/bittype/contracts/parse/go"
	"github.com/Bitspark/bittype/contracts/semantic/go"
	"github.com/Bitspark/bittype/contracts/text/go"
	"github.com/Bitspark/bittype/contracts/transform/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	parsehoros "github.com/Bitspark/bittype/parse/syntax/horos/provider/go"
	s "github.com/Bitspark/bittype/syntax/golang/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/ast/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// These qualify GEN-01: the kernel walks the independent source's published
// graph exactly as it walks horos, and enforces scope, whole-use admission,
// child relations and handle identity itself.

func independent(t *testing.T) *f.Fixture {
	t.Helper()
	fixture, err := f.New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return fixture
}

func independentText(t *testing.T, fixture *f.Fixture, roots ...string) string {
	t.Helper()
	plan, err := fixture.PlanRoots(context.Background(), f.Standard, roots...)
	if err != nil {
		t.Fatal(err)
	}
	out, err := plan.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return formatted(t, golangText(), outputTree[s.Language](t, out))
}

func independentFails(t *testing.T, fixture *f.Fixture, mode f.Mode, roots ...string) error {
	t.Helper()
	plan, err := fixture.PlanRoots(context.Background(), mode, roots...)
	if err == nil || plan.Valid() {
		t.Fatalf("%v planned: %v", roots, err)
	}
	return err
}

func TestKernelTraversalIndependentNestedLibraryCell(t *testing.T) {
	text := independentText(t, independent(t), "Cells")
	// Cell's library represents Cell<Person> as []Person; Cell's own body is never opened.
	if !strings.Contains(text, "type Cells = Page[[]Person, [][]Person]") || strings.Contains(text, "type Cell[") {
		t.Fatal(text)
	}
}

func TestKernelAdmitsOmittedActuals(t *testing.T) {
	fixture := independent(t)
	// The convention omits Ignore's argument; an admitted argument is not generated.
	if text := independentText(t, fixture, "Kept"); !strings.Contains(text, "type Kept = Ignore") || strings.Contains(text, "type Person") {
		t.Fatal(text)
	}
	// A refuted omitted argument is a rejection; an uncovered one is missing
	// evidence. Neither is the other.
	if err := independentFails(t, fixture, f.Standard, "Leak"); !errors.Is(err, g.ErrRejected) || errors.Is(err, semantic.ErrEvidenceUnavailable) {
		t.Fatal("refuted actual:", err)
	}
	if err := independentFails(t, fixture, f.Standard, "Blind"); !errors.Is(err, semantic.ErrEvidenceUnavailable) || errors.Is(err, g.ErrRejected) {
		t.Fatal("uncovered actual:", err)
	}
	// Admission covers the whole omitted subtree, and operands no rule captures.
	if err := independentFails(t, fixture, f.Standard, "Deep"); !errors.Is(err, g.ErrRejected) {
		t.Fatal("refuted content nested in an omitted actual:", err)
	}
	if err := independentFails(t, fixture, f.Standard, "Smuggle"); !foreignBinder(err) {
		t.Fatal("foreign binder as an omitted actual:", err)
	}
	if err := independentFails(t, fixture, f.Standard, "SmuggleLoose"); !errors.Is(err, g.ErrUnsupported) || !strings.Contains(err.Error(), "generic declaration used without arguments") {
		t.Fatal("unapplied generic as an omitted actual:", err)
	}
	if err := independentFails(t, fixture, f.Standard, "Hidden"); !foreignBinder(err) {
		t.Fatal("foreign binder in an uncaptured operand:", err)
	}
	if err := independentFails(t, fixture, f.Standard, "Shadow"); !foreignBinder(err) {
		t.Fatal("foreign binder in a body the recipe never represents:", err)
	}
}

// foreignBinder: the kernel classified the reference itself, at a site; it is
// not left to fail later for want of a native binding.
func foreignBinder(err error) bool {
	var located *g.SiteError
	return errors.Is(err, g.ErrScope) && errors.As(err, &located) && located.Message == "reference to a binder of another declaration" && located.Site.Node != 0
}

func TestKernelMapsEveryVerdict(t *testing.T) {
	fixture := independent(t)
	if err := independentFails(t, fixture, f.Standard, "Refused"); !errors.Is(err, g.ErrRejected) {
		t.Fatal("rejected declaration:", err)
	}
	if err := independentFails(t, fixture, f.Standard, "Void"); !errors.Is(err, g.ErrInvalid) || errors.Is(err, g.ErrRejected) || errors.Is(err, semantic.ErrEvidenceUnavailable) {
		t.Fatal("zero verdict:", err)
	}
}

func TestKernelScopeIsDirectBinderOwnership(t *testing.T) {
	fixture := independent(t)
	if err := independentFails(t, fixture, f.Standard, "Stray"); !foreignBinder(err) {
		t.Fatal("foreign binder:", err)
	}
	// One shared term reached through two bodies is two occurrences.
	if err := independentFails(t, fixture, f.Standard, "Borrow"); !foreignBinder(err) {
		t.Fatal("borrowed binder:", err)
	}
	if text := independentText(t, fixture, "Lend"); !strings.Contains(text, "type Lend[T0 interface {}] = Page[T0, []T0]") {
		t.Fatal(text)
	}
	if text := independentText(t, fixture, "SameA", "SameB"); !strings.Contains(text, "type SameA = Page[Person, []Person]") || !strings.Contains(text, "type SameB = Page[Person, []Person]") {
		t.Fatal(text)
	}
}

func TestKernelStructuralLimits(t *testing.T) {
	fixture := independent(t)
	for root, want := range map[string]string{"Later": "application without a direct declaration callee", "Loose": "generic declaration used without arguments"} {
		err := independentFails(t, fixture, f.Standard, root)
		var located *g.SiteError
		if !errors.Is(err, g.ErrUnsupported) || !errors.As(err, &located) || !strings.Contains(located.Message, want) {
			t.Fatal(root, err)
		}
	}
	if err := independentFails(t, fixture, f.Standard, "Rec"); !errors.Is(err, g.ErrUnsupported) || !strings.Contains(err.Error(), "recursive") {
		t.Fatal("recursion:", err)
	}
}

func TestKernelCapturesAreChildrenOfThisRun(t *testing.T) {
	fixture := independent(t)
	notChild := func(err error) bool {
		return errors.Is(err, g.ErrScope) && strings.Contains(err.Error(), "a capture must be a child")
	}
	if err := independentFails(t, fixture, f.StrayCapture, "Page"); !notChild(err) {
		t.Fatal("non-child capture:", err)
	}
	if err := independentFails(t, fixture, f.PatternStray, "Cells"); !notChild(err) {
		t.Fatal("pattern capture of the occurrence itself:", err)
	}
	// A declaration recipe may map only its own body.
	if err := independentFails(t, fixture, f.RecipeStray, "Page", "Outer"); !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), "may map only its body") {
		t.Fatal("recipe mapped a non-body use:", err)
	}
	// The first plan keeps a handle; a second plan over the same capture is a
	// different run, and the kept handle is foreign there.
	keeper := independent(t)
	if _, err := keeper.PlanRoots(context.Background(), f.ForeignHandle, "Page"); err != nil {
		t.Fatal(err)
	}
	if err := independentFails(t, keeper, f.ForeignHandle, "Page"); !notChild(err) {
		t.Fatal("foreign handle:", err)
	}
}

func TestKernelRootsAndEvidence(t *testing.T) {
	fixture := independent(t)
	other := independent(t)
	// Evidence published for another snapshot is not admitted.
	input := fixture.Input()
	input.Evidence = other.Input().Evidence
	if plan, err := fixture.PlanInput(context.Background(), input, "Person"); !errors.Is(err, semantic.ErrInvalidRequest) || plan.Valid() {
		t.Fatal("foreign evidence admitted:", err)
	}
	if plan, err := fixture.PlanInput(context.Background(), transform.Input{Source: fixture.Input().Source}, "Person"); !errors.Is(err, semantic.ErrInvalidRequest) || plan.Valid() {
		t.Fatal("missing evidence admitted:", err)
	}
	// A root must be the site of exactly one entity the source generates.
	for _, site := range []string{"Kindish", ""} {
		if plan, err := fixture.PlanSites(context.Background(), fixture.Site(site)); !errors.Is(err, g.ErrScope) || plan.Valid() {
			t.Fatalf("root %q: %v", site, err)
		}
	}
}

// A zero-parameter constructor alias applied to an argument and a curried
// application are valid horos with no direct structural reading; a
// constructor-sorted parameter is a declaration the horos judgment refuses to
// represent. Each is unsupported at its site, not malformed input.
func TestHorosUnsupportedFormsAtTheirSite(t *testing.T) {
	messages := map[string]string{"constructor-parameter": "constructor-valued types require an explicit representation", "applied-alias": "application without a direct declaration callee", "curried": "application without a direct declaration callee"}
	for name, source := range map[string]string{
		"constructor-parameter": "module aliases; type Wrap<F: top(P) -> top(P)> = text;",
		"applied-alias":         "module aliases; type ListType<T> = list<T>; type Alias: top(P) -> top(P) = ListType; type Texts = Alias<text>;",
		"curried":               "module aliases; type ListType<T> = list<T>; type Alias: top(P) -> top(P) = ListType; type Later<F: top(P) -> top(P)>: top(P) -> top(P) = F; type Texts = Later<Alias><text>;",
	} {
		t.Run(name, func(t *testing.T) {
			parsed, err := parsehoros.Parse(context.Background(), parse.Input{SourceID: "aliases.horos", TreeID: "canonical", Text: source})
			if err != nil {
				t.Fatal(err)
			}
			_, _, input := generationSource(t, parsed.Abstract)
			if !input.Evidence[0].HasSemantics() {
				t.Fatal("valid horos published no model")
			}
			var root ast.NodeRef
			last := func(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
				q, err := tsRequest(reg, input, to)
				q.Roots = q.Roots[len(q.Roots)-1:]
				root = q.Roots[0]
				return q, err
			}
			output, err := runGeneration(t, input, nativeChecker(t, typescriptText()), last)
			var diagnostic *text.Diagnostic
			if !errors.Is(err, g.ErrUnsupported) || errors.Is(err, semantic.ErrInvalidRequest) || !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || diagnostic.Node == nil || output.Snapshot().Valid() {
				t.Fatal(err)
			}
			// The diagnostic is the one located at the form, not the registry's
			// entry-node fallback.
			if *diagnostic.Node == input.Source.Entry() || !strings.Contains(diagnostic.Message, messages[name]) {
				t.Fatal("unlocated diagnostic:", diagnostic.Node, diagnostic.Message)
			}
			// The form is inside the root's body: a term with one recorded site
			// is reported there, not at its declaration.
			if *diagnostic.Node == root {
				t.Fatal("diagnostic at the declaration, not the form:", diagnostic.Node)
			}
		})
	}
}

// The regressions below come from the independent review that rejected the
// first GEN-01 delivery (board claim bittype-source-traversal-20261108, record
// 000003). Each has a correctly composed control.

// TestKernelPlansExactlyTheAdmittedEvidence: the source always interprets the
// evidence the kernel admitted. Two inputs share a snapshot and graph and
// differ only in a second report. Through the registry and through the
// exported planner, each input is judged by its own report, because Plan opens
// the source itself. A source that serves an interpretation opened over the
// other input cannot resolve this run's handles, though the graph is the same.
func TestKernelPlansExactlyTheAdmittedEvidence(t *testing.T) {
	ctx := context.Background()
	fixture := independent(t)
	holds, err := fixture.AdmissionInput(ctx, semantic.Holds)
	if err != nil {
		t.Fatal(err)
	}
	rejected, err := fixture.AdmissionInput(ctx, semantic.Rejected)
	if err != nil {
		t.Fatal(err)
	}
	if holds.Evidence[0] != rejected.Evidence[0] || holds.Evidence[1] == rejected.Evidence[1] {
		t.Fatal("the inputs must share the graph evidence and differ in the report")
	}
	admission := f.SourceVariant{Admission: true}
	for name, plan := range map[string]func(f.SourceVariant, transform.Input) (g.GenerationPlan, error){
		"registry": func(v f.SourceVariant, in transform.Input) (g.GenerationPlan, error) {
			return fixture.PlanSource(ctx, v, in, "Person")
		},
		"direct": func(v f.SourceVariant, in transform.Input) (g.GenerationPlan, error) {
			return fixture.PlanDirect(ctx, v, in, "Person")
		},
	} {
		if p, err := plan(admission, holds); err != nil || !p.Valid() {
			t.Fatal(name, "holding report:", err)
		}
		if p, err := plan(admission, rejected); !errors.Is(err, g.ErrRejected) || p.Valid() {
			t.Fatal(name, "rejecting report:", err)
		}
		reusing := admission.Reusing()
		if p, err := plan(reusing, holds); err != nil || !p.Valid() {
			t.Fatal(name, "first capture:", err)
		}
		if p, err := plan(reusing, rejected); !errors.Is(err, g.ErrScope) || errors.Is(err, g.ErrRejected) || p.Valid() {
			t.Fatal(name, "handles resolved against another capture of the same graph:", err)
		}
	}
}

// TestKernelSelectorsAreQualified: a source selector whose language differs in
// its profile, its identity or its Go marker, with an otherwise identical key,
// is outside the run. It chooses neither an exact binding, an exact declaration
// recipe nor the default recipe.
func TestKernelSelectorsAreQualified(t *testing.T) {
	ctx := context.Background()
	foreign, err := independent(t).ForeignLanguages()
	if err != nil {
		t.Fatal(err)
	}
	for _, example := range []struct{ root, selector string }{
		{"Person", "expression/person-body"}, // exact descriptor binding
		{"Cells", "declaration/Cell"},        // exact application binding
		{"Kept", "declaration/Ignore"},       // exact declaration recipe
		{"Person", "declaration/Person"},     // default declaration recipe
	} {
		for differs, language := range foreign {
			t.Run(example.selector+"/"+differs, func(t *testing.T) {
				fixture := independent(t)
				if p, err := fixture.PlanSource(ctx, f.SourceVariant{}, fixture.Input(), example.root); err != nil || !p.Valid() {
					t.Fatal("control:", err)
				}
				p, err := fixture.PlanSource(ctx, f.SourceVariant{Foreign: example.selector, ForeignIn: language}, fixture.Input(), example.root)
				if !errors.Is(err, g.ErrScope) || p.Valid() || !strings.Contains(err.Error(), "selector outside the source language") {
					t.Fatal("foreign selector:", err)
				}
			})
		}
	}
}

// TestKernelLocatesUseFailuresAtTheUse: when a use reaches a callee the source
// does not generate, or one whose selector is outside the source language, the
// failure belongs to the use and is reported in the declaration it occurs in,
// not at the callee.
func TestKernelLocatesUseFailuresAtTheUse(t *testing.T) {
	ctx := context.Background()
	fixture := independent(t)
	p, err := fixture.PlanSource(ctx, f.SourceVariant{Foreign: "declaration/Page"}, fixture.Input(), "Lend")
	var callee *g.SiteError
	if !errors.Is(err, g.ErrScope) || p.Valid() || !errors.As(err, &callee) || callee.Site != fixture.Site("Lend") {
		t.Fatalf("foreign callee selector at %v, want Lend %v (Page is %v): %v", callee, fixture.Site("Lend"), fixture.Site("Page"), err)
	}
	err = independentFails(t, fixture, f.Standard, "UsesKind")
	var located *g.SiteError
	if !errors.Is(err, g.ErrUnsupported) || !errors.As(err, &located) || located.Site != fixture.Site("UsesKind") || !strings.Contains(err.Error(), "the source does not generate this declaration") {
		t.Fatalf("ungenerated callee at %v, want UsesKind %v (Kindish is %v): %v", located, fixture.Site("UsesKind"), fixture.Site("Kindish"), err)
	}
}

// TestKernelLocatesSharedTermsAtTheirDeclaration: when the graph records a term
// at several syntax sites, the kernel cannot tell which one an occurrence is
// and reports the declaration the occurrence is in.
func TestKernelLocatesSharedTermsAtTheirDeclaration(t *testing.T) {
	fixture, err := f.NewSharedSites(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if p, err := fixture.PlanRoots(context.Background(), f.Standard, "Lend"); err != nil || !p.Valid() {
		t.Fatal("control:", err)
	}
	err = independentFails(t, fixture, f.Standard, "Borrow")
	var located *g.SiteError
	if !foreignBinder(err) || !errors.As(err, &located) || located.Site != fixture.Site("Borrow") {
		t.Fatalf("site %v, want Borrow %v (Lend is %v): %v", located, fixture.Site("Borrow"), fixture.Site("Lend"), err)
	}
}
