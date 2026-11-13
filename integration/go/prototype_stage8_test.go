package integration_test

// PROTOTYPE (bittype#196), stage 8: the fixtures of the vertical cases and
// the identity cases of row 32.
//
//	Case 1  not built (O6): recursion is refused, see stage 2.
//	Case 2  below, at the kernel's level, with a MOCK family.
//	Case 3  below: an exact generic that omits its parameter natively, with a value facet.
//	Case 4  stage 5: Page[Person] lowered and executed.
//	Case 5  below: a fixture generic that requires evidence about its argument.
//	Case 6  stage 9: the contribution of step 3.
//
// The four identity cases of row 32 are in stage 4: repeated uses share a
// node (TwoUsesOneNode), equal native types with different adapters stay
// apart and an erased native argument is kept for conversion
// (OmittedArgumentIsKeptForConversion), one plan under two TypeScript
// placements (RealizationManifest).

import (
	"context"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	f "github.com/Bitspark/bittype/integration/go/testdata/generationfixture"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
)

// Vertical case 2, `Page<Outcome>` with optional nullable fields, at the
// kernel's level, through the independent fixture source, with a MOCK family.
//
// The family "optional nullable fields" is not built (#210, #211). The mock
// stands in for it: the fixture source issues a role for each field of
// Outcome, a mock contribution answers that role with a member fragment, and
// the mock record lays out nothing itself. What the mock does with a value
// means nothing.
//
// Shows: the kernel plans a generic applied to a record whose every field is
// delegated, through a source that is not horos: the role selectors of that
// source, the selection by role, the fragments in the record, one field node
// for two positions with equal decisions, and the application Page<Outcome>.
// Does not show: null against absence, or discriminators. It shows nothing
// about the family the mock stands for, and nothing about conversion.
func TestPrototypeStage8Case2PageOfOutcomeWithAMockFamily(t *testing.T) {
	plan, err := independentPlan(t, independentOptions{}, true, "Outcomes")
	if err != nil {
		t.Fatal(err)
	}
	outcome := bodyNode(t, plan.Plan, "Outcome")
	if len(outcome.Links) != 2 || outcome.Links[0].Kind != "field" || outcome.Links[0].Role != f.OptionalNullableRole || outcome.Links[0].Stated != "Field0" || outcome.Links[1].Stated != "Field1" {
		t.Fatalf("the mock record's links are %+v", outcome.Links)
	}
	for i, link := range outcome.Links {
		field, _ := plan.Plan.Node(link.Node)
		if field.Kind != "field" || field.Contribution.Identity.Key != "mock-optional-nullable" {
			t.Errorf("field %d is answered by %+v", i, field.Contribution)
		}
	}
	if strings.Join(outcome.Members, ",") != "Field0,Field1" {
		t.Errorf("the mock record has the members %v", outcome.Members)
	}
	application := bodyNode(t, plan.Plan, "Outcomes")
	if application.Application == nil || len(application.Application.Adapters) != 1 {
		t.Errorf("Page<Outcome> is %+v", application.Application)
	}
	text := protoText(t, golangText(), plan)
	if !strings.Contains(text, "type Outcomes = Page[Outcome]") || !strings.Contains(text, "Field0 []") {
		t.Errorf("unexpected native text:\n%s", text)
	}
	t.Logf("GREEN vertical case 2 with a MOCK family plans at the kernel's level:\n%s", text)
}

// fixtureValue is a facet kind of a fixture product: "a value facet". Its
// value means nothing; the case is about the facet being required and found.
var fixtureValue = func() g.FacetKind {
	k, err := g.NewFacetKind(g.FacetKindDescription{Identity: "example.org/fixtures/value-facet", Version: "1"})
	if err != nil {
		panic(err)
	}
	return k
}()

// Vertical case 3: "An exact generic with an omitted or an additional native
// parameter and a value facet".
//
// Shows: an exact contribution for Cell<T> whose native type omits T. Asked
// for a product that requires a value facet of every selected use, the
// request plans, the argument is linked for conversion, and the Cell
// contribution's facet is the one the product finds. Asked for types alone,
// the argument is not planned. Without the facet the request for the product
// is "missing facet".
// Does not show: capabilities, or what Cell does with its slot. It does not
// show an additional native parameter; the fixture has the omitted form only.
func TestPrototypeStage8Case3ExactGenericWithAnOmittedParameter(t *testing.T) {
	fixture := protoFixture(t)
	composition := independentComposition(t, fixture, independentOptions{cellConversionOnly: true})
	product := g.Product{Identity: "example.org/fixtures/values", Requires: map[g.SubjectKind][]g.FacetKind{g.UseSubject: {fixtureValue}}}
	composition.Products, composition.Requested = []g.Product{product}, []string{product.Identity}
	var cell g.ContributionIdentity
	for _, c := range composition.Contributions {
		if c.Project == nil || c.Domain.Kind() == g.RoleDomain {
			continue
		}
		if c.Domain.Kind() == g.ExactDomain {
			cell = c.Reference.Identity
		}
		composition.Facets = append(composition.Facets, protoMust[g.Facet](t)(fixtureValue.Issue(fixturePublisher, c.Reference.Identity, golangTarget, "a value facet of "+c.Reference.Identity.Key)))
	}
	plan, err := fixture.PlanContributions(context.Background(), composition, "CellOfPerson")
	if err != nil {
		t.Fatal(err)
	}
	use := bodyNode(t, plan.Plan, "CellOfPerson")
	if use.Kind != "use" || use.Contribution.Identity != cell || use.Links[0].Need != "conversion" || use.Links[0].Node == "" {
		t.Errorf("Cell<Person> is %+v", use)
	}
	if text := protoText[gs.Language](t, golangText(), plan); !strings.Contains(text, "type CellOfPerson = string") {
		t.Errorf("the native type does not omit the parameter:\n%s", text)
	}
	t.Log("GREEN vertical case 3: the exact generic omits its parameter natively, its argument is linked for conversion and its value facet is found")
	// Without the Cell contribution's facet.
	without := composition
	without.Facets = nil
	for _, facet := range composition.Facets {
		if facet.Contribution() != cell {
			without.Facets = append(without.Facets, facet)
		}
	}
	_, err = fixture.PlanContributions(context.Background(), without, "CellOfPerson")
	if refusal := expectRefusal(t, "the exact generic without its value facet", err, "missing facet"); refusal.Contribution.Identity != cell {
		t.Errorf("the refusal names %s", refusal.Contribution)
	}
}

// Vertical case 5: "A fixture generic that requires evidence about its
// argument, accepted for one operand and refused for another".
//
// Probed<T> is a fixture generic whose body's rule declares that T's actuals
// must be proved nil-free. The fixture policy states beside its decision what
// was established about a declaration's body, so a use of Person can answer.
//
// Shows: the requirement is found on the open declaration and discharged at
// each application: Probed<Person> is accepted where Person's body is proved
// nil-free, refused as "requirement unsatisfied" where it is refuted, and
// Probed<Page<Person>> is refused as "requirement unresolved", because nothing
// is stated about a page. The refusal names the requirement.
// Does not show: the map family or the meaning of its keys.
func TestPrototypeStage8Case5GenericRequirementAcceptedAndRefused(t *testing.T) {
	accepted, err := independentPlan(t, independentOptions{probeRequires: true, personFacts: []g.Fact{fixtureNilFree.Proved("")}}, false, "ProbedPerson")
	if err != nil {
		t.Fatalf("the requirement is not discharged by a proved fact: %v", err)
	}
	application := bodyNode(t, accepted.Plan, "ProbedPerson")
	if d := application.Application.Discharges; len(d) != 1 || d[0].Outcome != "proved" || d[0].Kind != fixtureNilFree || d[0].Reference == nil {
		t.Errorf("Probed<Person> records the discharges %+v", d)
	}
	if r := symbolNode(t, accepted.Plan, "Probed").Requirements; len(r) != 1 || r[0].Kind != fixtureNilFree {
		t.Errorf("Probed<T> records the requirements %+v", r)
	}
	t.Logf("GREEN vertical case 5: Probed<Person> is accepted, discharged by %s", application.Application.Discharges[0].Reference)
	_, err = independentPlan(t, independentOptions{probeRequires: true, personFacts: []g.Fact{fixtureNilFree.Refuted("")}}, false, "ProbedPerson")
	if refusal := expectRefusal(t, "Probed<Person> where the fact is refuted", err, "requirement unsatisfied"); refusal.Requirement != fixtureNilFree {
		t.Errorf("the refusal names the requirement %+v", refusal.Requirement)
	}
	_, err = independentPlan(t, independentOptions{probeRequires: true, personFacts: []g.Fact{fixtureNilFree.Proved("")}}, false, "ProbedPage")
	if refusal := expectRefusal(t, "Probed<Page<Person>> where nothing is stated", err, "requirement unresolved"); refusal.Requirement != fixtureNilFree {
		t.Errorf("the refusal names the requirement %+v", refusal.Requirement)
	}
}
