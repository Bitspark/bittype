package integration_test

// PROTOTYPE (bittype#196), stage 2: evidence. Rows 3, 5, 7, 8 and 9 of gate
// line 2 and guarantee 4.

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

const protoEvidenceModule = `module example.evidence;

type Page<T> = { items: list<T>; };
type TextPage = Page<text>;
type Written = { items: list<text>; };
`

// evidenceBundle is the Go bundle with the fixture family in the place of the
// list contribution and the fixture text, which declares the fact, in the
// place of the text contribution.
func evidenceBundle(t *testing.T, from representation.Language, o evidenceListOptions) bundle[gs.Language] {
	t.Helper()
	list, listCodec := evidenceList(t, from, o)
	text, textCodec := nilFreeText(t, from, "golang/1", "1")
	return golangBundle(t, from).replace(list).replace(text).withCodecs(listCodec, textCodec)
}

// symbolNode finds the node of the declaration with the given source name.
func symbolNode(t *testing.T, plan g.RepresentationPlan, name string) g.PlanNode {
	t.Helper()
	for _, s := range plan.Symbols() {
		if strings.HasSuffix(s.ID, "/"+name) {
			node, _ := plan.Node(s.Node)
			return node
		}
	}
	t.Fatalf("the plan has no symbol %s", name)
	return g.PlanNode{}
}

// bodyNode is the node selected for the body of a declaration.
func bodyNode(t *testing.T, plan g.RepresentationPlan, name string) g.PlanNode {
	t.Helper()
	declaration := symbolNode(t, plan, name)
	if len(declaration.Links) != 1 {
		t.Fatalf("declaration %s has %d links", name, len(declaration.Links))
	}
	node, _ := plan.Node(declaration.Links[0].Node)
	return node
}

// Row 3: "A fixture family whose decision depends on evidence about its
// operand plans and generates".
//
// Shows: a rule asks for a typed fact about what was selected for its
// operand, gets it as a typed answer before it decides, and its decision and
// native type follow the answer: direct where the fact is proved, the form
// that needs no proof where nothing is stated. The fact is one the dependency
// declared at registration. The answer is a link of the frozen record, bound
// to the dependency's contribution reference.
// Does not show: that the fact is true (AC-05). The kernel carries it unread.
func TestPrototypeStage2FixtureFamilyDecidesOnEvidence(t *testing.T) {
	from, input := protoSource(t, "module example.evidence;\n\ntype Written = { items: list<text>; };\n")
	proved, err := protoPlanGo(t, input, evidenceBundle(t, from, evidenceListOptions{}), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	// The same family over bittype's own text contribution, which states no
	// such fact: the answer is unknown.
	list, listCodec := evidenceList(t, from, evidenceListOptions{})
	unknown, err := protoPlanGo(t, input, golangBundle(t, from).replace(list).withCodecs(listCodec), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name          string
		plan          g.ContributionPlan
		outcome, want string
	}{{"the dependency declares the fact", proved, "proved", "[]string"}, {"the dependency states nothing", unknown, "unknown", "[][]string"}} {
		fields, _, _ := goFieldTypes(t, protoText[gs.Language](t, golangText(), c.plan))
		if got := fields["Written"]["Items"]; got != c.want {
			t.Errorf("%s: the field has the type %s, want %s", c.name, got, c.want)
		}
		record := bodyNode(t, c.plan.Plan, "Written")
		element, _ := c.plan.Plan.Node(record.Links[0].Node)
		if len(element.Links) != 1 || len(element.Links[0].Answers) != 1 {
			t.Fatalf("%s: the list node records no answer: %+v", c.name, element.Links)
		}
		answer := element.Links[0].Answers[0]
		if answer.Kind != fixtureNilFree || answer.Outcome != c.outcome || answer.Reference == nil || element.Links[0].Need != "evidence" {
			t.Errorf("%s: the answer is %+v", c.name, answer)
		}
		t.Logf("GREEN %s: answer %s about %s, native type %s", c.name, answer.Outcome, answer.Reference, c.want)
	}
}

// Rows 7 and 8: "`Page<text>` and the written-out record generate the same
// native type when a rule uses evidence, and a rule that reads one
// application to decide fails."
//
// Shows: a rule that declares the fact as a requirement of the binder plans
// the open declaration once, the application Page<text> discharges the
// requirement with what was selected for text, and the application and the
// written-out record have the same native type. A rule that instead makes a
// choice for want of proof inside the generic declaration gives Page<text>
// another native type than the written-out record, and this check detects it.
// The kernel gives no rule an application to read: inside Page<T> the answer
// about T is "deferred", whatever Page is applied to.
// Does not show: that the kernel refuses the incoherent rule. It cannot: it
// does not know what a rule does with an answer. The check is a qualification
// test of a family, and "fails" is read as "fails this test" (finding F5).
func TestPrototypeStage2GenericAndWrittenOutAgree(t *testing.T) {
	from, input := protoSource(t, protoEvidenceModule)
	coherent := func(o evidenceListOptions) (g.ContributionPlan, error) {
		plan, err := protoPlanGo(t, input, evidenceBundle(t, from, o), protoRequest{})
		if err != nil {
			return plan, err
		}
		fields, parameters, aliases := goFieldTypes(t, protoText[gs.Language](t, golangText(), plan))
		if aliases["TextPage"] != "Page[string]" || len(parameters["Page"]) != 1 {
			return plan, fmt.Errorf("unexpected declarations: %v %v", aliases, parameters)
		}
		applied := strings.ReplaceAll(fields["Page"]["Items"], parameters["Page"][0], "string")
		if written := fields["Written"]["Items"]; applied != written {
			return plan, fmt.Errorf("Page<text> has the field type %s and the written-out record has %s", applied, written)
		}
		return plan, nil
	}
	plan, err := coherent(evidenceListOptions{require: true})
	if err != nil {
		t.Fatal(err)
	}
	t.Log("GREEN a rule that declares the binder's requirement: Page<text> and the written-out record have the same native type")
	if _, err := coherent(evidenceListOptions{}); err == nil {
		t.Error("a rule that decides differently inside the generic declaration is not detected")
	} else {
		t.Logf("RED   a rule that makes a choice for want of proof inside the generic declaration only: %v", err)
	}
	// What the planner recorded. The requirement is stated in the frozen
	// record of the rule that declared it, as the answer "deferred, required",
	// and on the node of the open declaration. The application discharges it.
	page := symbolNode(t, plan.Plan, "Page")
	if len(page.Requirements) != 1 || page.Requirements[0] != (g.PlanRequirement{Binder: 0, Kind: fixtureNilFree}) {
		t.Errorf("the open declaration records the requirements %+v", page.Requirements)
	}
	record, _ := plan.Plan.Node(page.Links[0].Node)
	list, _ := plan.Plan.Node(record.Links[0].Node)
	if a := list.Links[0].Answers[0]; a.Outcome != "deferred" || !a.Required || a.Reference != nil {
		t.Errorf("inside Page<T> the answer about T is %+v", a)
	}
	application := bodyNode(t, plan.Plan, "TextPage")
	if application.Application == nil || len(application.Application.Discharges) != 1 || application.Application.Discharges[0].Outcome != "proved" || application.Application.Discharges[0].Reference == nil {
		t.Errorf("Page<text> records the discharges %+v", application.Application)
	}
}

// Row 8, O2: the five outcomes.
//
// Shows: each outcome reaches a rule as its own typed answer: proved and
// exhausted from facts declared at registration, refuted from a fact stated
// beside a decision, unknown where the selection states nothing, deferred
// where the operand is bound to a binder.
// Does not show: a check that reaches a limit. "Exhausted" is stated by the
// dependency here; nothing in the prototype computes evidence (O13).
func TestPrototypeStage2FiveOutcomes(t *testing.T) {
	from, input := protoSource(t, `module example.outcomes;

type Person = { name: text; };
type Proved = list<text>;
type Refuted = list<list<text>>;
type Unknown = list<Person>;
type Exhausted = list<bool>;
type Deferred<T> = list<T>;
`)
	b := evidenceBundle(t, from, evidenceListOptions{states: []g.Fact{fixtureNilFree.Refuted("a Go slice can be nil")}})
	boolean := b.find(t, "horos-to-golang/boolean")
	boolean.Reference.Identity = g.ContributionIdentity{Publisher: fixturePublisher, Key: "exhausted-boolean"}
	boolean.Facts = []g.Fact{fixtureNilFree.Exhausted()}
	plan, err := protoPlanGo(t, input, b.replace(boolean), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Proved", "Refuted", "Unknown", "Exhausted", "Deferred"} {
		answers := bodyNode(t, plan.Plan, name).Links[0].Answers
		if len(answers) != 1 || answers[0].Outcome != strings.ToLower(name) {
			t.Errorf("%s: the answer is %+v", name, answers)
		} else {
			t.Logf("GREEN the rule of %s was answered %q", name, answers[0].Outcome)
		}
	}
}

// Row 3, the first consumer of the evidence path: "the declaration recipe's
// need to know that a selected body is a record".
//
// Shows: the declaration contribution decides between the nominal form and
// the alias from the answer about its body, and reads no payload and no form
// of another family. Without the fact the same module becomes aliases.
// Does not show: anything about a body that is a record under another
// contribution; that contribution states the fact or not, as its owner decides.
func TestPrototypeStage2DeclarationPolicyAsksAboutItsBody(t *testing.T) {
	from, input := protoSource(t, "module example.body;\n\ntype Person = { name: text; };\n")
	b := golangBundle(t, from)
	with, err := protoPlanGo(t, input, b, protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	silent := b.find(t, "horos-to-golang/record")
	silent.Facts = nil
	without, err := protoPlanGo(t, input, b.replace(silent), protoRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if text := protoText[gs.Language](t, golangText(), with); !strings.Contains(text, "type Person struct") {
		t.Errorf("with the fact the declaration is not nominal:\n%s", text)
	}
	if text := protoText[gs.Language](t, golangText(), without); !strings.Contains(text, "type Person = struct") {
		t.Errorf("without the fact the declaration is not an alias:\n%s", text)
	}
	if a, b := symbolNode(t, with.Plan, "Person").Links[0].Answers[0].Outcome, symbolNode(t, without.Plan, "Person").Links[0].Answers[0].Outcome; a != "proved" || b != "unknown" {
		t.Errorf("the answers were %s and %s", a, b)
	}
	t.Log("GREEN proved: type Person struct; unknown: type Person = struct")
}

// Row 5, O20: "Two contributions that differ only in their representation
// revision give equal payloads and different node identities".
//
// Shows: an invariant of a revision is carried by the reference and needs no
// field: the payloads are byte-identical, and the node is another. A node that
// depends on it is another too, because its answers and links are bound to
// the reference.
// Does not show: what the revision means. That is its publisher's.
func TestPrototypeStage2RevisionIsIdentityAndNoPayload(t *testing.T) {
	from, input := protoSource(t, "module example.revision;\n\ntype Names = list<text>;\n")
	plan := func(revision string) g.RepresentationPlan {
		list, listCodec := evidenceList(t, from, evidenceListOptions{})
		text, textCodec := nilFreeText(t, from, revision, "1")
		p, err := protoPlanGo(t, input, golangBundle(t, from).replace(list).replace(text).withCodecs(listCodec, textCodec), protoRequest{})
		if err != nil {
			t.Fatal(err)
		}
		return p.Plan
	}
	one, two := plan("golang/1"), plan("golang/2")
	listOne, listTwo := bodyNode(t, one, "Names"), bodyNode(t, two, "Names")
	textOne, _ := one.Node(listOne.Links[0].Node)
	textTwo, _ := two.Node(listTwo.Links[0].Node)
	if textOne.Payload != textTwo.Payload || listOne.Payload != listTwo.Payload {
		t.Error("the payloads differ")
	}
	if textOne.ID == textTwo.ID || listOne.ID == listTwo.ID {
		t.Error("the nodes are the same under two representation revisions")
	}
	if textOne.Contribution.Representation == textTwo.Contribution.Representation {
		t.Error("the plan does not name the revision")
	}
	t.Logf("GREEN equal payloads %q, nodes %s and %s", textOne.Payload, textOne.ID, textTwo.ID)
}

// Row 9: "An evidence cycle and a recursive reference are refused with
// different reasons".
//
// Shows: the contract tells a reference from a demand for evidence, and the
// two refusals can be told apart by their reason, in one class.
// Does not show: anything of recursion beyond its refusal (O6, row 9a).
func TestPrototypeStage2CycleAndRecursionAreToldApart(t *testing.T) {
	// The cases are stated in the independent fixture source: the horos
	// analysis establishes no sort for a recursive declaration, so the horos
	// source refuses one before the kernel meets the recursion.
	fixture := protoFixture(t)
	composition := independentComposition(t, fixture, independentOptions{})
	expect := func(name, root, reason string) {
		_, err := fixture.PlanContributions(context.Background(), composition, root)
		var refusal *g.Refusal
		if !errors.As(err, &refusal) || !errors.Is(err, g.ErrUnsupported) || refusal.Reason != reason || refusal.Owner != g.KernelOwner {
			t.Errorf("%s: got %v, want the refusal %q", name, err, reason)
			return
		}
		t.Logf("GREEN %s is refused: %v", name, err)
	}
	// Rec = Page<Rec>. One open request on the cycle is a reference: the
	// application refers to its argument.
	expect("a recursive reference", "Rec", g.ReasonRecursiveReference)
	// Loop = probe(Loop). Every open request on the cycle demands evidence:
	// the declaration asks about its body, and the probe asks about its operand.
	expect("an evidence cycle", "Loop", g.ReasonEvidenceCycle)
}

// stashing wraps the Go record contribution as a fixture that asks for
// evidence about each of its fields and can keep and present answers.
func stashing(t *testing.T, base golangContribution, hook func(step horosDecide)) golangContribution {
	base.Reference.Identity = g.ContributionIdentity{Publisher: fixturePublisher, Key: "stashing-record"}
	decide := base.Rule.Decide
	base.Rule.Name = func(_ context.Context, step horosName) ([]horosRequest, error) {
		fields, err := step.Reader.Fields(step.Subject.Use)
		var requests []horosRequest
		for _, f := range fields {
			requests = append(requests, g.Evidence(f.Use, g.Demand{Kind: fixtureNilFree}))
		}
		return requests, err
	}
	base.Rule.Decide = func(ctx context.Context, step horosDecide) (g.Decision, error) {
		hook(step)
		return decide(ctx, step)
	}
	return base
}

// Guarantee 4: "An answer issued for another dependency, another generic
// environment or another contribution reference is rejected. The reference
// counts whole, both of its parts (O19)."
//
// The entry point the prototype fixes: an answer is a value that a rule reads
// through the step it was issued to, DecideStep.Present(dependency, answer).
// A rule that kept an answer from another step and presents it is rejected,
// and the planning of its subject fails even where the rule drops the error.
//
// Shows: each of the four bindings is compared: the dependency, the generic
// environment, and both parts of the third axis of the dependency's
// reference. Each is red under a kernel that leaves that one comparison out.
// Does not show: that a rule cannot act on an answer it remembers without
// presenting it. A rule is trusted code in its own package; the guarantee is
// about what the kernel lets a rule read as an answer.
func TestPrototypeStage2AnswersAreBoundToTheirDependency(t *testing.T) {
	two := "module example.bound;\n\ntype A = { a: text; b: text; };\n"
	both := two + "type B = { a: text; b: text; };\n"
	plan := func(module string, revision, meaning string, hook func(horosDecide)) error {
		from, input := protoSource(t, module)
		text, textCodec := nilFreeText(t, from, revision, meaning)
		b := golangBundle(t, from)
		_, err := protoPlanGo(t, input, b.replace(stashing(t, b.find(t, "horos-to-golang/record"), hook)).replace(text).withCodecs(textCodec), protoRequest{})
		return err
	}
	rejected := func(err error) error {
		if err == nil {
			return fmt.Errorf("the answer was read")
		}
		if !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), "not issued for this dependency") {
			return fmt.Errorf("unexpected failure: %w", err)
		}
		return nil
	}
	// The rule drops the error of Present on purpose.
	redGreen(t, "an answer issued for another dependency is rejected", func() error {
		return rejected(plan(two, "golang/1", "1", func(step horosDecide) { step.Present(1, step.Answers(0)[0]) }))
	}, "evidence-ignores-dependency")
	redGreen(t, "an answer issued in another generic environment is rejected", func() error {
		var kept []g.Answer
		return rejected(plan(both, "golang/1", "1", func(step horosDecide) {
			if kept == nil {
				kept = step.Answers(0)
				return
			}
			step.Present(0, kept[0])
		}))
	}, "evidence-ignores-environment")
	across := func(revision, meaning string) func() error {
		return func() error {
			var kept []g.Answer
			if err := plan(two, "golang/1", "1", func(step horosDecide) { kept = step.Answers(0) }); err != nil {
				return err
			}
			return rejected(plan(two, revision, meaning, func(step horosDecide) { step.Present(0, kept[0]) }))
		}
	}
	redGreen(t, "an answer issued for another representation revision of the dependency is rejected", across("golang/2", "1"), "evidence-ignores-representation-revision")
	redGreen(t, "an answer issued for another meaning revision of the dependency is rejected", across("golang/1", "2"), "evidence-ignores-meaning-revision")
	// The control: an answer presented to the step it was issued to is read.
	if err := plan(two, "golang/1", "1", func(step horosDecide) {
		if outcome, _, err := step.Present(0, step.Answers(0)[0]); err != nil || outcome != g.Proved {
			t.Errorf("the step's own answer was not read: %v %v", outcome, err)
		}
	}); err != nil {
		t.Error(err)
	}
}
