package integration_test

// Evidence under the contribution contract (bittype#196,
// docs/contribution-contract.md, "Evidence"): gate line 2, rows 3, 7, 8 and
// 9, with the decisions O23, O29, O30, O35 and O40. The tests on the horos
// source use a fixture family that represents a horos list by what was
// selected for its element. The tests of a cycle and of evidence through a
// declaration run on the independent fixture source.

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// gateAnswered is what a rule read about its first dependency through its
// step: how many answers the planner issued, and the answer of the kind
// gateNilFree.
type gateAnswered struct {
	issued      int
	outcome     g.Outcome
	content     string
	established bool
}

// gateObserving makes a contribution for the horos source record what its
// rule is answered about its first dependency, before the rule decides.
func gateObserving(base goContribution, seen *[]gateAnswered) goContribution {
	decide := base.Rule.Decide
	base.Rule.Decide = func(ctx context.Context, step horosDecideStep) (g.Decision, error) {
		issued := len(step.Answers(0))
		outcome, content, established := step.Evidence(0, gateNilFree)
		*seen = append(*seen, gateAnswered{issued, outcome, content, established})
		return decide(ctx, step)
	}
	return base
}

// gateAnswerIs compares an answer that the plan records on a link.
func gateAnswerIs(a g.PlanAnswer, outcome g.Outcome, content string, required bool, reference *g.ContributionReference) bool {
	if a.Kind != gateNilFree || a.Outcome != outcome.String() || a.Content != content || a.Required != required {
		return false
	}
	if reference == nil {
		return a.Reference == nil
	}
	return a.Reference != nil && *a.Reference == *reference
}

// gateEvidenceLink is the one link of the node selected for the body of a
// declaration, with its one answer.
func gateEvidenceLink(t *testing.T, plan g.RepresentationPlan, declaration string) (g.PlanNode, g.PlanLink) {
	t.Helper()
	node := gateBody(t, plan, declaration)
	if len(node.Links) != 1 || len(node.Links[0].Answers) != 1 {
		t.Fatalf("the node selected for the body of %s records the links %+v, want one link with one answer", declaration, node.Links)
	}
	return node, node.Links[0]
}

// Row 3 (O13, O23), gate line 2.
//
// Shows: a fixture family for the horos list asks for a fact of one kind
// about what was selected for its element, and its rule reads the answer
// through its step before it decides. Where the selected text contribution
// declared the fact at registration, the rule is answered "proved" with the
// content that the declaring contribution wrote, and the generated field is
// the direct form. Where the selected text contribution states nothing, the
// rule is answered "unknown" and the field has the form that needs no proof.
// The rule is issued one answer, of the kind it asked for, although the
// dependency declares a fact of another kind first. The answer is a link of
// the frozen record of the asking rule: need "evidence", the kind, the
// outcome, the content, and the whole contribution reference of the
// dependency.
// Does not show: that the fact is true. The kernel carries it unread (AC-05).
// A check that computes evidence: nothing computes, the dependency states.
// Red under: E-facts-not-handed, E-declared-facts-dropped,
// E-answer-ignores-kind, E-answer-content-dropped, E-answer-not-a-link,
// E-answer-without-reference.
func TestGateFamilyDecidesOnEvidenceAboutItsOperand(t *testing.T) {
	from, input := gateSource(t, "module example.evidence;\n\ntype Written = { items: list<text>; };\n")
	const content = "a Go string is no pointer"
	another := g.EvidenceKind{Identity: gatePublisher + "/golang/another-kind", Version: "1"}
	stating := gateNilFreeText(t, from, "golang/1", "1")
	stating.Facts = []g.Fact{another.Refuted("a fact of another kind"), gateNilFree.Proved(content)}
	silent := gateGolang(t, from).find(t, "horos-to-golang/text")
	if len(silent.Facts) != 0 {
		t.Fatalf("the builtin text contribution declares the facts %+v, so it is no dependency that states nothing", silent.Facts)
	}
	for _, c := range []struct {
		name    string
		text    goContribution
		outcome g.Outcome
		content string
		want    string
	}{
		{"the dependency declares the fact", stating, g.Proved, content, "[]string"},
		{"the dependency states nothing", silent, g.Unknown, "", "[][]string"},
	} {
		var seen []gateAnswered
		list := gateObserving(gateEvidenceList(t, from, gateEvidenceListOptions{}), &seen)
		generation, err := gatePlanGo(t, input, gateGolang(t, from).replace(list).replace(c.text).withCodecs(gateDirectCodec(t), gateEmptyCodec(t)), gateRequest{})
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		if want := []gateAnswered{{1, c.outcome, c.content, false}}; !slices.Equal(seen, want) {
			t.Errorf("%s: the rule read %+v through its step, want %+v", c.name, seen, want)
		}
		if got := gateParseGo(t, gateText[gs.Language](t, golangText(), generation)).field("Written", "Items"); got != c.want {
			t.Errorf("%s: the field has the type %s, want %s", c.name, got, c.want)
		}
		plan := generation.Representation()
		record := gateBody(t, plan, "Written")
		if len(record.Links) != 1 {
			t.Errorf("%s: the record has the links %+v", c.name, record.Links)
			continue
		}
		element, _ := plan.Node(record.Links[0].Node)
		if element.Contribution == nil || *element.Contribution != list.Reference || len(element.Links) != 1 {
			t.Errorf("%s: the node of the list is %+v", c.name, element)
			continue
		}
		link := element.Links[0]
		reference := c.text.Reference
		if link.Kind != "child" || link.Need != "evidence" || len(link.Answers) != 1 || !gateAnswerIs(link.Answers[0], c.outcome, c.content, false, &reference) {
			t.Errorf("%s: the link of the list to its element is %+v, want one answer %q with the content %q about %s", c.name, link, c.outcome, c.content, reference)
		}
		if selected, ok := plan.Node(link.Node); !ok || selected.Contribution == nil || *selected.Contribution != reference {
			t.Errorf("%s: the link leads to the node %+v and not to the selection of %s", c.name, selected, reference)
		}
	}
}

// gateCoherenceModule states one record twice: as an application of a generic
// declaration, and written out.
const gateCoherenceModule = `module example.coherence;

type Page<T> = { items: list<T>; };
type TextPage = Page<text>;
type Written = { items: list<text>; };
`

// gateCoherent is the coherence test of the fixture family of the list for
// gateCoherenceModule: mapping after substitution agrees with substituting
// the mapped argument into the mapped generic declaration. It compares the
// field of Page with its parameter replaced by the native type of text with
// the field of the written-out record.
func gateCoherent(t *testing.T, generation g.GenerationPlan) error {
	t.Helper()
	native := gateParseGo(t, gateText[gs.Language](t, golangText(), generation))
	if native.aliases["TextPage"] != "Page[string]" || len(native.parameters["Page"]) != 1 {
		return fmt.Errorf("unexpected declarations: aliases %v, parameters %v", native.aliases, native.parameters)
	}
	applied := strings.ReplaceAll(native.field("Page", "Items"), native.parameters["Page"][0], "string")
	if written := native.field("Written", "Items"); applied != written {
		return fmt.Errorf("Page<text> has the field type %s and the written-out record has %s", applied, written)
	}
	return nil
}

// Rows 7 and 8 (O13, O29), gate line 2.
//
// Shows: with a family that declares its fact as a requirement of the binder,
// Page<text> and the written-out record generate the same native type. The
// rule inside Page<T> is answered "deferred" and "established", with no
// contribution reference, and the open declaration records the requirement as
// (binder, kind). The application Page<text> records the discharge: proved,
// with the reference of what was selected for text. An application whose
// actual states nothing or states "exhausted" is refused as requirement
// unresolved, and one whose actual refutes the fact as requirement
// unsatisfied, each with the requirement and the actual's contribution named.
// An actual that is itself a binder moves the requirement to that binder: the
// outer declaration records it, the inner application records "deferred", and
// the applications of the outer declaration discharge it or are refused. The
// same family without the declared requirement makes a choice for want of
// proof inside the generic declaration only, and the family's coherence test
// (gateCoherent) fails for it.
// Does not show: that the kernel refuses the incoherent rule. It does not and
// cannot (O29): it gives no rule an application to read and does not know
// what a rule does with an answer. Coherence over every substitution: the
// test compares one application with one written-out record.
// Red under: E-binder-answer-not-deferred, E-requirement-not-collected,
// E-established-without-requirement, E-discharge-not-recorded,
// E-unproved-requirement-admitted, E-refuted-requirement-is-unresolved,
// E-requirement-not-moved.
func TestGateGenericAndWrittenOutAgree(t *testing.T) {
	from, input := gateSource(t, gateCoherenceModule)
	text := gateNilFreeText(t, from, "golang/1", "1").Reference
	required, err := gatePlanGo(t, input, gateEvidenceBundle(t, from, gateEvidenceListOptions{require: true}), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if err := gateCoherent(t, required); err != nil {
		t.Errorf("a family that declares the binder's requirement is not coherent: %v", err)
	}
	choosing, err := gatePlanGo(t, input, gateEvidenceBundle(t, from, gateEvidenceListOptions{}), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if err := gateCoherent(t, choosing); err == nil {
		t.Error("a family that decides differently inside the generic declaration passes its coherence test")
	} else if !strings.Contains(err.Error(), "[][]string") || !strings.Contains(err.Error(), "the written-out record has []string") {
		t.Errorf("the coherence test fails for another reason than the two forms: %v", err)
	}

	// What the planner recorded for the family that declares the requirement.
	plan := required.Representation()
	page := gateDeclaration(t, plan, "Page")
	if want := []g.PlanRequirement{{Binder: 0, Kind: gateNilFree}}; !slices.Equal(page.Requirements, want) {
		t.Errorf("the open declaration records the requirements %+v, want %+v", page.Requirements, want)
	}
	record := gateBody(t, plan, "Page")
	if len(record.Links) != 1 {
		t.Fatalf("the record of Page has the links %+v", record.Links)
	}
	list, _ := plan.Node(record.Links[0].Node)
	if len(list.Links) != 1 || len(list.Links[0].Answers) != 1 || !gateAnswerIs(list.Links[0].Answers[0], g.Deferred, "", true, nil) {
		t.Errorf("inside Page<T> the list records the links %+v, want one answer that is deferred and required and names no contribution", list.Links)
	}
	application := gateBody(t, plan, "TextPage")
	if a := application.Application; a == nil || len(a.Discharges) != 1 || a.Discharges[0].Binder != 0 || a.Discharges[0].Kind != gateNilFree || a.Discharges[0].Outcome != g.Proved.String() || a.Discharges[0].Reference == nil || *a.Discharges[0].Reference != text {
		t.Errorf("Page<text> records the application %+v, want one discharge of binder 0 that is proved by %s", a, text)
	}
	// Without the declared requirement the answer is deferred and establishes
	// nothing, and neither the declaration nor the application records anything.
	chosen := choosing.Representation()
	chosenRecord := gateBody(t, chosen, "Page")
	chosenList, _ := chosen.Node(chosenRecord.Links[0].Node)
	if len(chosenList.Links) != 1 || len(chosenList.Links[0].Answers) != 1 || !gateAnswerIs(chosenList.Links[0].Answers[0], g.Deferred, "", false, nil) {
		t.Errorf("without a requirement the list inside Page<T> records the links %+v", chosenList.Links)
	}
	if a := gateBody(t, chosen, "TextPage").Application; len(gateDeclaration(t, chosen, "Page").Requirements) != 0 || a == nil || len(a.Discharges) != 0 {
		t.Errorf("without a requirement the plan records requirements or discharges: %+v", a)
	}

	// Each application discharges the requirement or is refused.
	exhausted := func(from representation.Language) gateBundle[gs.Language] {
		b := gateEvidenceBundle(t, from, gateEvidenceListOptions{require: true})
		boolean := b.find(t, "horos-to-golang/boolean")
		boolean.Reference.Identity = g.ContributionIdentity{Publisher: gatePublisher, Key: "exhausted-boolean"}
		boolean.Facts = []g.Fact{gateNilFree.Exhausted()}
		return b.replace(boolean)
	}
	plain := func(from representation.Language) gateBundle[gs.Language] {
		return gateEvidenceBundle(t, from, gateEvidenceListOptions{require: true})
	}
	refuting := func(from representation.Language) gateBundle[gs.Language] {
		return gateEvidenceBundle(t, from, gateEvidenceListOptions{require: true, states: []g.Fact{gateNilFree.Refuted("a Go slice can be nil")}})
	}
	const generic = "module example.discharge;\n\ntype Page<T> = { items: list<T>; };\n"
	for _, c := range []struct {
		name, application, class, culprit string
		bundle                            func(representation.Language) gateBundle[gs.Language]
	}{
		{"an actual that states nothing", "type Refused = Page<bool>;\n", "requirement unresolved", "horos-to-golang/boolean", plain},
		{"an actual that states exhausted", "type Refused = Page<bool>;\n", "requirement unresolved", "exhausted-boolean", exhausted},
		{"an actual that refutes the fact", "type Refused = Page<list<text>>;\n", "requirement unsatisfied", "evidence-list", refuting},
		{"an actual of a declaration that moved the requirement", "type Outer<U> = Page<U>;\ntype Refused = Outer<bool>;\n", "requirement unresolved", "horos-to-golang/boolean", plain},
	} {
		from, input := gateSource(t, generic+c.application)
		bundle := c.bundle(from)
		generation, err := gatePlanGo(t, input, bundle, gateRequest{})
		if err == nil || generation.Representation().Valid() {
			t.Errorf("%s: the application is planned (error: %v)", c.name, err)
			continue
		}
		refusal := gateExpectRefusal(t, c.name, err, c.class)
		if culprit := bundle.find(t, c.culprit).Reference; refusal.Owner != g.KernelOwner || refusal.Requirement != gateNilFree || refusal.Contribution != culprit || !strings.HasPrefix(refusal.Reason, g.ReasonRequirement) {
			t.Errorf("%s: the refusal is %+v, want the kernel's, about the requirement %+v and the actual's contribution %s", c.name, refusal, gateNilFree, culprit)
		}
	}

	// An actual that is itself a binder moves the requirement to that binder.
	from, input = gateSource(t, generic+"type Outer<U> = Page<U>;\ntype TextOuter = Outer<text>;\n")
	moved, err := gatePlanGo(t, input, plain(from), gateRequest{})
	if err != nil {
		t.Fatal(err)
	}
	plan = moved.Representation()
	if outer, want := gateDeclaration(t, plan, "Outer"), []g.PlanRequirement{{Binder: 0, Kind: gateNilFree}}; !slices.Equal(outer.Requirements, want) {
		t.Errorf("the outer declaration records the requirements %+v, want %+v", outer.Requirements, want)
	}
	if a := gateBody(t, plan, "Outer").Application; a == nil || len(a.Discharges) != 1 || a.Discharges[0].Outcome != g.Deferred.String() || a.Discharges[0].Reference != nil || a.Discharges[0].Kind != gateNilFree {
		t.Errorf("the application Page<U> inside Outer records %+v, want one deferred discharge", a)
	}
	if a := gateBody(t, plan, "TextOuter").Application; a == nil || len(a.Discharges) != 1 || a.Discharges[0].Outcome != g.Proved.String() || a.Discharges[0].Reference == nil || *a.Discharges[0].Reference != text {
		t.Errorf("the application Outer<text> records %+v, want one discharge that is proved by %s", a, text)
	}
}

// Row 8 (O2, O23), gate line 2.
//
// Shows: each of the five outcomes reaches a rule as its own answer, read
// through the rule's step, and is recorded under its own name on the link.
// Proved and exhausted come from facts that a contribution declared at
// registration, refuted from a fact that a rule stated beside its decision,
// with its content. Unknown is the planner's answer about a nominal
// declaration, which states nothing. Deferred is the planner's answer about
// an operand that is bound to a binder, and it establishes nothing where no
// requirement was declared. A rule that states a fact with an outcome that
// nobody states (unknown, deferred) or without a kind fails as a rule
// failure, which matches no refusal value.
// Does not show: a check that reaches a limit. "Exhausted" is stated by the
// dependency, nothing computes evidence (O13).
// Red under: E-binder-answer-not-deferred, E-stated-facts-dropped,
// E-exhausted-folded-into-unknown, E-refuted-folded-into-unknown,
// E-facts-not-handed, E-established-without-requirement,
// E-stated-fact-any-outcome, E-stated-fact-without-kind.
func TestGateFiveOutcomesReachARule(t *testing.T) {
	const refutation = "a Go slice can be nil"
	from, input := gateSource(t, `module example.outcomes;

type Person = { name: text; };
type Proved = list<text>;
type Refuted = list<list<text>>;
type Unknown = list<Person>;
type Exhausted = list<bool>;
type Deferred<T> = list<T>;
`)
	bundle := func(seen *[]gateAnswered, states []g.Fact) gateBundle[gs.Language] {
		b := gateEvidenceBundle(t, from, gateEvidenceListOptions{states: states})
		boolean := b.find(t, "horos-to-golang/boolean")
		boolean.Reference.Identity = g.ContributionIdentity{Publisher: gatePublisher, Key: "exhausted-boolean"}
		boolean.Facts = []g.Fact{gateNilFree.Exhausted()}
		return b.replace(boolean).replace(gateObserving(b.find(t, "evidence-list"), seen))
	}
	for _, c := range []struct {
		root    string
		outcome g.Outcome
		content string
		// inner is what the rule of the inner list of Refuted reads first.
		inner []gateAnswered
	}{
		{"Proved", g.Proved, "", nil},
		{"Refuted", g.Refuted, refutation, []gateAnswered{{1, g.Proved, "", false}}},
		{"Unknown", g.Unknown, "", nil},
		{"Exhausted", g.Exhausted, "", nil},
		{"Deferred", g.Deferred, "", nil},
	} {
		var seen []gateAnswered
		generation, err := gatePlanGo(t, input, bundle(&seen, []g.Fact{gateNilFree.Refuted(refutation)}), gateRequest{roots: []string{c.root}})
		if err != nil {
			t.Errorf("%s: %v", c.root, err)
			continue
		}
		if want := append(slices.Clone(c.inner), gateAnswered{1, c.outcome, c.content, false}); !slices.Equal(seen, want) {
			t.Errorf("%s: the rules read %+v through their steps, want %+v", c.root, seen, want)
		}
		_, link := gateEvidenceLink(t, generation.Representation(), c.root)
		if a := link.Answers[0]; a.Kind != gateNilFree || a.Outcome != strings.ToLower(c.root) || a.Content != c.content || a.Required || (a.Reference == nil) != (c.outcome == g.Deferred) {
			t.Errorf("%s: the link records the answer %+v", c.root, a)
		}
	}
	// What a rule may state beside its decision.
	for _, c := range []struct {
		name string
		fact g.Fact
	}{
		{"the outcome unknown", g.Fact{Kind: gateNilFree, Outcome: g.Unknown}},
		{"the outcome deferred", g.Fact{Kind: gateNilFree, Outcome: g.Deferred}},
		{"no outcome", g.Fact{Kind: gateNilFree}},
		{"no kind", g.Fact{Outcome: g.Proved}},
	} {
		var seen []gateAnswered
		generation, err := gatePlanGo(t, input, bundle(&seen, []g.Fact{c.fact}), gateRequest{roots: []string{"Proved"}})
		var failure *g.RuleFailure
		if err == nil || generation.Representation().Valid() || !errors.As(err, &failure) || !errors.Is(failure.Err, g.ErrInvalid) || failure.Contribution != gateReference("evidence-list", "golang/1") {
			t.Errorf("a rule that states a fact with %s: got %v, want a rule failure of the stating contribution", c.name, err)
		}
		if matched := gateRefusalClasses(err); len(matched) != 0 {
			t.Errorf("a rule that states a fact with %s: the failure matches the refusal values %v", c.name, matched)
		}
	}
}

// Row 7 (O23), gate line 2: "a contribution declares at registration the
// facts that its revision fixes", and a stated fact has one of three outcomes.
//
// Shows: a contribution that declares at registration a fact with an outcome
// that nobody states (unknown, deferred), with no outcome or with no kind is
// refused as an invalid description, by the registry or by the planner, and
// no rule is answered from it. The same contribution with a proved fact
// registers and plans.
// Does not show: that a declared fact is true.
// Red under: no mutant. The test is skipped as PENDING: it is red against the
// kernel as built (finding F-B1), so no mutant can show it red.
func TestGateDeclaredFactIsOneOfThreeOutcomes(t *testing.T) {
	t.Skip("PENDING F-B1: a contribution that declares at registration a fact with the outcome deferred, unknown or none, or without a kind, registers and plans, and a rule is answered from it; only a fact stated beside a decision is checked")
	from, input := gateSource(t, "module example.declared;\n\ntype Names = list<text>;\n")
	plan := func(fact g.Fact) ([]gateAnswered, error) {
		var seen []gateAnswered
		text := gateNilFreeText(t, from, "golang/1", "1")
		text.Facts = []g.Fact{fact}
		list := gateObserving(gateEvidenceList(t, from, gateEvidenceListOptions{}), &seen)
		_, err := gatePlanGo(t, input, gateGolang(t, from).replace(list).replace(text).withCodecs(gateDirectCodec(t), gateEmptyCodec(t)), gateRequest{})
		return seen, err
	}
	if seen, err := plan(gateNilFree.Proved("")); err != nil || !slices.Equal(seen, []gateAnswered{{1, g.Proved, "", false}}) {
		t.Fatalf("the contribution with a proved fact: the rule read %+v (error: %v)", seen, err)
	}
	for _, c := range []struct {
		name string
		fact g.Fact
	}{
		{"the outcome unknown", g.Fact{Kind: gateNilFree, Outcome: g.Unknown}},
		{"the outcome deferred", g.Fact{Kind: gateNilFree, Outcome: g.Deferred}},
		{"no outcome", g.Fact{Kind: gateNilFree}},
		{"no kind", g.Fact{Outcome: g.Proved}},
	} {
		seen, err := plan(c.fact)
		if err == nil {
			t.Errorf("a contribution that declares a fact with %s registers and plans, and the rule that asked read %+v", c.name, seen)
			continue
		}
		if !errors.Is(err, g.ErrInvalid) || len(gateRefusalClasses(err)) != 0 || len(seen) != 0 {
			t.Errorf("a contribution that declares a fact with %s: got %v (the rule read %+v), want an invalid description", c.name, err, seen)
		}
	}
}

// Row 9 (O30, O40, row 9a), gate line 2.
//
// Shows, on the independent fixture source: a cycle on which every open
// request demands evidence is refused as an evidence cycle (Loop, whose body
// asks about a use of Loop). A cycle with one reference on it is refused as a
// recursive reference, although the declaration contribution demands evidence
// about the body on the same cycle (Rec, whose body applies Page to Rec). A
// cycle whose one request names its operand for conversion only counts as a
// reference when bindings are asked for, and is no cycle for types alone,
// where the operand is not demanded (Spin). Both refusals are unsupported,
// match exactly that one value and have the kernel as owner. They differ in
// their reason, and the two reasons are two texts.
// Does not show: anything of recursion beyond its refusal (row 9a). The
// same on the horos source: the horos analysis refuses a recursive
// declaration before the kernel meets it.
// Red under: E-cycle-any-evidence, E-cycle-reasons-equal,
// E-conversion-counts-as-evidence, E-cycle-all-references.
func TestGateEvidenceCycleAndRecursiveReference(t *testing.T) {
	if g.ReasonEvidenceCycle == g.ReasonRecursiveReference {
		t.Errorf("the two refusals of a cycle have one reason, %q, and cannot be told apart", g.ReasonEvidenceCycle)
	}
	texts := map[string]string{}
	for _, c := range []struct {
		name, root, reason string
		bindings           bool
	}{
		{"every open request demands evidence", "Loop", g.ReasonEvidenceCycle, false},
		{"one reference among demands for evidence", "Rec", g.ReasonRecursiveReference, false},
		{"a need for conversion, when bindings are asked for", "Spin", g.ReasonRecursiveReference, true},
	} {
		generation, err := indPlan(t, indOptions{}, c.bindings, c.root)
		if err == nil || generation.Representation().Valid() {
			t.Errorf("%s: %s is planned (error: %v)", c.name, c.root, err)
			continue
		}
		refusal := gateExpectRefusal(t, c.name, err, "unsupported")
		if refusal.Owner != g.KernelOwner || refusal.Reason != c.reason || refusal.Subject != "declaration/"+c.root {
			t.Errorf("%s: the refusal has the owner %q, the reason %q and the subject %q, want the kernel's %q about declaration/%s", c.name, refusal.Owner, refusal.Reason, refusal.Subject, c.reason, c.root)
		}
		texts[c.root] = strings.ReplaceAll(err.Error(), c.root, "X")
	}
	if texts["Loop"] == "" || texts["Loop"] == texts["Rec"] {
		t.Errorf("the refusal of an evidence cycle and the refusal of a recursive reference read alike: %q", texts["Loop"])
	}
	// For types alone the operand of Spin is not demanded, so nothing is open
	// on a cycle and the declaration plans.
	generation, err := indPlan(t, indOptions{}, false, "Spin")
	if err != nil {
		t.Fatalf("for types alone Spin does not plan: %v", err)
	}
	if carry := gateBody(t, generation.Representation(), "Spin"); len(carry.Links) != 1 || carry.Links[0].Need != "conversion, not demanded" || carry.Links[0].Node != "" {
		t.Errorf("for types alone the body of Spin records the links %+v", carry.Links)
	}
}

// O35 and O41 (row 7), gate line 2.
//
// Shows, on the independent fixture source: a rule that asks about an operand
// of a named type is answered by the declaration. For a declaration that its
// contribution decides as an alias (Word), the kernel forwards the facts of
// the body: the answer is the body's fact with its outcome and content, and
// it is bound to the contribution reference of the body. That holds for a
// refuted fact as for a proved one, and for an alias whose body is its own
// binder (Wrap<Word>), which is represented as its actual is. For a
// declaration that becomes a defined declaration (Person) the answer is
// unknown and is bound to the declaration contribution's reference, although
// the contribution selected for its body declares the fact. A binder's
// requirement is discharged through an alias by the body's fact, proved or
// refuted, and is unresolved for a nominal declaration.
// Does not show: which facts a nominal declaration restates. Nothing is
// decided for that case (#209). An alias in the source: "alias" is here the
// target form that the declaration contribution chooses (O41).
// Red under: E-alias-facts-not-forwarded, E-nominal-forwards-body-facts,
// E-alias-answers-with-own-reference, E-alias-of-binder-not-forwarded.
func TestGateEvidenceThroughADeclaration(t *testing.T) {
	word, policy, person := gateReference("independent-word", "golang/1"), gateReference("independent-declarations", "golang/1"), gateReference("independent-person", "golang/1")
	const refutation = "a word can be nil here"
	declares := []g.Fact{gateNilFree.Proved("")}
	refutes := []g.Fact{gateNilFree.Refuted(refutation)}
	for _, c := range []struct {
		name, root string
		options    indOptions
		outcome    g.Outcome
		content    string
		reference  g.ContributionReference
		native     string
		// declaration is the declaration that answers, and transparent what its node records.
		declaration string
		transparent bool
	}{
		{"an alias forwards a proved fact of its body", "ProbeWord", indOptions{}, g.Proved, "", word, "[]Word", "Word", true},
		{"an alias forwards a refuted fact of its body", "ProbeWord", indOptions{wordFacts: refutes}, g.Refuted, refutation, word, "[][]Word", "Word", true},
		{"an alias whose body is its binder answers as its actual", "ProbeWrapped", indOptions{}, g.Proved, "", word, "[]Wrap[Word]", "Wrap", true},
		{"a nominal declaration answers for itself", "ProbePerson", indOptions{personFacts: declares}, g.Unknown, "", policy, "[][]Person", "Person", false},
	} {
		generation, err := indPlan(t, c.options, false, c.root)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
			continue
		}
		plan := generation.Representation()
		_, link := gateEvidenceLink(t, plan, c.root)
		reference := c.reference
		if !gateAnswerIs(link.Answers[0], c.outcome, c.content, false, &reference) {
			t.Errorf("%s: the answer is %+v (reference %v), want %q with the content %q about %s", c.name, link.Answers[0], link.Answers[0].Reference, c.outcome, c.content, reference)
		}
		if got := gateParseGo(t, gateText[gs.Language](t, golangText(), generation)).aliases[c.root]; got != c.native {
			t.Errorf("%s: %s is generated as %q, want %q", c.name, c.root, got, c.native)
		}
		if declaration := gateDeclaration(t, plan, c.declaration); declaration.Transparent != c.transparent {
			t.Errorf("%s: the node of %s records transparent = %v", c.name, c.declaration, declaration.Transparent)
		}
	}
	// The control of the nominal case: the contribution selected for the body
	// of Person did declare the fact that the declaration does not forward.
	generation, err := indPlan(t, indOptions{personFacts: declares}, false, "ProbePerson")
	if err != nil {
		t.Fatal(err)
	}
	if entry, ok := generation.Representation().Contribution(person); !ok || !slices.ContainsFunc(entry.Facts, func(fact g.PlanFact) bool {
		return fact.Kind == gateNilFree && fact.Outcome == g.Proved.String()
	}) {
		t.Errorf("the contribution selected for the body of Person declares the facts %+v, so the nominal case shows nothing", entry.Facts)
	}

	// A binder's requirement through a declaration.
	generation, err = indPlan(t, indOptions{probeRequires: true}, false, "ProbedWord")
	if err != nil {
		t.Fatalf("Probed<Word> does not plan: %v", err)
	}
	if a := gateBody(t, generation.Representation(), "ProbedWord").Application; a == nil || len(a.Discharges) != 1 || a.Discharges[0].Outcome != g.Proved.String() || a.Discharges[0].Reference == nil || *a.Discharges[0].Reference != word {
		t.Errorf("Probed<Word> records the application %+v, want one discharge that is proved by %s", a, word)
	}
	for _, c := range []struct {
		name, root, class string
		options           indOptions
		culprit           g.ContributionReference
	}{
		{"the requirement and an alias whose body refutes the fact", "ProbedWord", "requirement unsatisfied", indOptions{probeRequires: true, wordFacts: refutes}, word},
		{"the requirement and a nominal declaration", "ProbedPerson", "requirement unresolved", indOptions{probeRequires: true, personFacts: declares}, policy},
	} {
		generation, err := indPlan(t, c.options, false, c.root)
		if err == nil || generation.Representation().Valid() {
			t.Errorf("%s: %s is planned (error: %v)", c.name, c.root, err)
			continue
		}
		if refusal := gateExpectRefusal(t, c.name, err, c.class); refusal.Requirement != gateNilFree || refusal.Contribution != c.culprit || refusal.Owner != g.KernelOwner {
			t.Errorf("%s: the refusal is %+v, want the kernel's, about %+v and %s", c.name, refusal, gateNilFree, c.culprit)
		}
	}
}
