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
}

func TestKernelScopeIsDirectBinderOwnership(t *testing.T) {
	fixture := independent(t)
	// The kernel classifies the reference itself; it is not left to fail later
	// for want of a native binding.
	foreign := func(err error) bool {
		var located *g.SiteError
		return errors.Is(err, g.ErrScope) && errors.As(err, &located) && located.Message == "reference to a binder of another declaration"
	}
	if err := independentFails(t, fixture, f.Standard, "Stray"); !foreign(err) {
		t.Fatal("foreign binder:", err)
	}
	// One shared term reached through two bodies is two occurrences.
	if err := independentFails(t, fixture, f.Standard, "Borrow"); !foreign(err) {
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
	if err := independentFails(t, fixture, f.StrayCapture, "Page"); !errors.Is(err, g.ErrScope) {
		t.Fatal("non-child capture:", err)
	}
	// The first plan keeps a handle; a second plan over the same capture is a
	// different run, and the kept handle is foreign there.
	if _, err := fixture.PlanRoots(context.Background(), f.ForeignHandle, "Page"); err != nil {
		t.Fatal(err)
	}
	if err := independentFails(t, fixture, f.ForeignHandle, "Page"); !errors.Is(err, g.ErrScope) {
		t.Fatal("foreign handle:", err)
	}
}

func TestKernelRootsAndEvidence(t *testing.T) {
	fixture := independent(t)
	other := independent(t)
	// Evidence published for another snapshot is not admitted.
	input := fixture.Input()
	input.Evidence = other.Input().Evidence
	if plan, err := fixture.PlanInput(context.Background(), input, "Person"); err == nil || plan.Valid() {
		t.Fatal("foreign evidence admitted")
	}
	if plan, err := fixture.PlanInput(context.Background(), transform.Input{Source: fixture.Input().Source}, "Person"); err == nil || plan.Valid() {
		t.Fatal("missing evidence admitted")
	}
}

// A zero-parameter constructor alias applied to an argument, and a curried
// application, are valid horos with no direct structural reading. They are
// unsupported at their site, not malformed input.
func TestHorosNoStructuralViewIsUnsupported(t *testing.T) {
	for name, source := range map[string]string{
		"applied-alias": "module aliases; type ListType<T> = list<T>; type Alias: top(P) -> top(P) = ListType; type Texts = Alias<text>;",
		"curried":       "module aliases; type ListType<T> = list<T>; type Alias: top(P) -> top(P) = ListType; type Later<F: top(P) -> top(P)>: top(P) -> top(P) = F; type Texts = Later<Alias><text>;",
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
			last := func(reg *c.Registry, input transform.Input, to representation.Language) (c.Request, error) {
				q, err := tsRequest(reg, input, to)
				q.Roots = q.Roots[len(q.Roots)-1:]
				return q, err
			}
			output, err := runGeneration(t, input, nativeChecker(t, typescriptText()), last)
			var diagnostic *text.Diagnostic
			if !errors.Is(err, g.ErrUnsupported) || errors.Is(err, semantic.ErrInvalidRequest) || !errors.As(err, &diagnostic) || diagnostic.Code != text.Unsupported || diagnostic.Node == nil || output.Snapshot().Valid() {
				t.Fatal(err)
			}
		})
	}
}
