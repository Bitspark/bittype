package integration_test

// Guarantees 1 to 4 of the kernel (bittype#196, docs/contribution-contract.md,
// "What the kernel carries without reading it", row 27). Each is a structural,
// authority, scope or association guarantee. The facet kinds and the product
// of these tests belong to a fixture product of the test: the kernel knows no
// kind.

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	g "github.com/Bitspark/bittype/contracts/generation/go"
	shared "github.com/Bitspark/bittype/integration/go/testdata/composition"
	gs "github.com/Bitspark/bittype/syntax/golang/go"
	h "github.com/Bitspark/bittype/syntax/horos/go"
	source "github.com/Bitspark/bittype/transform/syntax/horos/source/go"
	c "github.com/Bitspark/bittype/transform/syntax/shared/generation/go"
	"github.com/Bitspark/bittype/trees/representation/go"
)

// gateFacetModule selects five contributions: text, list and record for uses,
// the optional contribution for a field position and the declaration
// contribution. It selects no contribution for booleans.
const gateFacetModule = `module example.facets;

type Person = { name: text; nick?: text; };
type Page<T> = { items: list<T>; };
type People = Page<Person>;
`

// gateFacetSelected are the local keys of the Go contributions that
// gateFacetModule selects.
var gateFacetSelected = []string{"horos-to-golang/text", "horos-to-golang/list", "horos-to-golang/record", "horos-to-golang/optional", "horos-to-golang/declarations"}

// gateFacetProduct is a fixture product that asks for bindings. It requires
// one facet kind of a contribution selected for a use, another of one selected
// for a field position and a third of one selected for a declaration.
type gateFacetProduct struct {
	value, field, declared g.FacetKind
	product                g.Product
}

func gateNewFacetProduct(t *testing.T) gateFacetProduct {
	t.Helper()
	kind := func(name string) g.FacetKind {
		return gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/" + name, Version: "1"}))
	}
	p := gateFacetProduct{value: kind("value-adapter"), field: kind("field-state"), declared: kind("declared-adapter")}
	p.product = g.Product{Identity: gatePublisher + "/conversion", Requires: map[g.SubjectKind][]g.FacetKind{g.UseSubject: {p.value}, g.FieldSubject: {p.field}, g.DeclarationSubject: {p.declared}}}
	return p
}

// kind is the facet kind that the product requires of a contribution.
func (p gateFacetProduct) kind(contribution goContribution) g.FacetKind {
	switch contribution.Domain.Kind() {
	case g.RoleDomain:
		return p.field
	case g.DeclarationPolicyDomain:
		return p.declared
	}
	return p.value
}

// compose is the Go bundle with the product and one facet for each selected
// contribution except the ones named, and with the further facets.
func (p gateFacetProduct) compose(t *testing.T, from representation.Language, except []string, more ...g.Facet) gateBundle[gs.Language] {
	t.Helper()
	b := gateGolang(t, from)
	b.products = []g.Product{p.product}
	for _, key := range gateFacetSelected {
		if slices.Contains(except, key) {
			continue
		}
		contribution := b.find(t, key)
		identity := contribution.Reference.Identity
		b = b.withFacets(gateMust[g.Facet](t)(p.kind(contribution).Issue(identity.Publisher, identity, goTarget, "the facet of "+key)))
	}
	return b.withFacets(more...)
}

// Guarantee 1 (rows 26 and 27, O5, O8, O32), gate line 2.
//
// Shows: with a facet for each selected contribution, a request for types
// and bindings plans, and the plan names both products. With the facet of one
// selected contribution removed, the request for bindings is refused as
// "missing facet", which matches exactly that one value, with the kernel as
// owner, the kind, the product and the contribution named, and the request
// for types alone succeeds. That holds for each of the five selected
// contributions: three selected for a use, one for a field position, which
// has no facet of the kind that a use needs, and the declaration
// contribution. The contribution for booleans, which the module does not
// select, has no facet and none is asked of it. What is checked is the set of
// facets that the composition registered for that contribution: the facets of
// the other contributions, of the same kind, do not answer for it. The
// registry admits a facet only for a registered contribution of the facet's
// publisher.
// Does not show: what a facet does. The kernel hands a value to the owner of
// its kind and reads none. The planner's own check of a facet's publisher,
// which no composition reaches through the registry.
// Red under: G1-facet-of-any-contribution, G1-products-not-checked,
// G1-field-subject-unchecked, G1-missing-facet-as-unsupported,
// G1-registry-admits-stranger, G1-registry-admits-unregistered.
func TestGuarantee1RequiredFacetsBelongToTheSelectedContribution(t *testing.T) {
	from, input := gateSource(t, gateFacetModule)
	p := gateNewFacetProduct(t)
	bindings := gateRequest{products: []string{p.product.Identity}}
	generation, err := gatePlanGo(t, input, p.compose(t, from, nil), bindings)
	if err != nil {
		t.Fatalf("the bundle with every facet does not plan for bindings: %v", err)
	}
	plan := generation.Representation()
	if want := []string{g.TypesProduct, p.product.Identity}; !slices.Equal(plan.Products(), want) {
		t.Errorf("the plan names the products %v, want %v", plan.Products(), want)
	}
	var selected []string
	for _, entry := range plan.Contributions() {
		selected = append(selected, entry.Reference.Identity.Key)
	}
	slices.Sort(selected)
	if want := slices.Sorted(slices.Values(gateFacetSelected)); !slices.Equal(selected, want) {
		t.Fatalf("the module selects %v, want %v", selected, want)
	}
	for _, key := range gateFacetSelected {
		without := p.compose(t, from, []string{key})
		contribution := without.find(t, key)
		generation, err := gatePlanGo(t, input, without, bindings)
		if err == nil || generation.Representation().Valid() {
			t.Errorf("without the facet of %s the request for bindings is planned (error: %v)", key, err)
			continue
		}
		refusal := gateExpectRefusal(t, "bindings without the facet of "+key, err, "missing facet")
		kind := p.kind(contribution).Reference().Identity
		if refusal.Owner != g.KernelOwner || refusal.Contribution != contribution.Reference || refusal.Target != goTarget || !strings.HasPrefix(refusal.Reason, g.ReasonFacet) || !strings.Contains(refusal.Reason, kind) || !strings.Contains(refusal.Reason, p.product.Identity) {
			t.Errorf("without the facet of %s the refusal is %+v, want the kernel's, about the kind %s, the product and %s", key, refusal, kind, contribution.Reference)
		}
		if generation, err := gatePlanGo(t, input, without, gateRequest{}); err != nil || !generation.Representation().Valid() {
			t.Errorf("without the facet of %s the request for types alone fails: %v", key, err)
		}
	}
	// A facet is admitted only for a registered contribution of its publisher.
	list := gateGolang(t, from).find(t, "horos-to-golang/list").Reference.Identity
	register := func(facet g.Facet) error {
		var registry c.Registry
		b := gateGolang(t, from)
		to := nativeChecker(t, golangText()).Language()
		if err := (shared.Set[h.Language, gs.Language, source.Reader]{Contributions: b.contributions, Codecs: b.codecs}).Register(&registry, from, to); err != nil {
			t.Fatal(err)
		}
		return registry.RegisterFacet(facet)
	}
	if err := register(gateMust[g.Facet](t)(p.value.Issue(list.Publisher, list, goTarget, "the facet"))); err != nil {
		t.Errorf("the facet of the list's publisher for the list is not admitted: %v", err)
	}
	if err := register(gateMust[g.Facet](t)(p.value.Issue("example.org/a-stranger", list, goTarget, "a stranger's facet"))); !errors.Is(err, g.ErrScope) {
		t.Errorf("a facet of another publisher than its contribution's is admitted: %v", err)
	}
	unregistered := g.ContributionIdentity{Publisher: list.Publisher, Key: "horos-to-golang/unregistered"}
	if err := register(gateMust[g.Facet](t)(p.value.Issue(list.Publisher, unregistered, goTarget, "a facet of nothing"))); !errors.Is(err, g.ErrScope) {
		t.Errorf("a facet that names no registered contribution is admitted: %v", err)
	}
}

// Guarantee 2 (rows 26 and 27), gate line 2.
//
// Shows: a facet under the identity of the kind that the product requires is
// rejected where it does not match that kind. A facet of a kind that another
// owner made under the same identity and version, a facet of the required
// kind for another target, and a facet of a kind with the same identity and
// another contract version are each rejected with the facet error, which is a
// scope error and matches no refusal value: it is no "missing facet". No plan
// results. The matching facet plans, and a request for types alone, which
// requires no facet, plans with each mismatched facet. A kind's validator
// runs when a facet is issued, and a kind hands a facet's value to its own
// owner only.
// Does not show: the comparison of the contract version alone. A kind's
// issuer is made with the kind, so a kind with another version has another
// issuer and is rejected by the comparison of the issuer first. No value that
// the contract can express reaches the comparison of the version with the
// issuers equal. The mutant G2-facet-ignores-version therefore survives.
// Red under: G2-facet-ignores-issuer, G2-facet-ignores-target,
// G2-issue-skips-validator, G2-read-ignores-issuer. Not red under
// G2-facet-ignores-version, for the reason above.
func TestGuarantee2FacetMatchesKindTargetAndVersion(t *testing.T) {
	from, input := gateSource(t, gateFacetModule)
	p := gateNewFacetProduct(t)
	bindings := gateRequest{products: []string{p.product.Identity}}
	const key = "horos-to-golang/list"
	list := gateGolang(t, from).find(t, key).Reference.Identity
	sameName := func(version string) g.FacetKind {
		return gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: p.value.Reference().Identity, Version: version}))
	}
	issue := func(kind g.FacetKind, target string) g.Facet {
		return gateMust[g.Facet](t)(kind.Issue(list.Publisher, list, target, "the facet of the list"))
	}
	if generation, err := gatePlanGo(t, input, p.compose(t, from, []string{key}, issue(p.value, goTarget)), bindings); err != nil || !generation.Representation().Valid() {
		t.Fatalf("the matching facet does not plan: %v", err)
	}
	for _, c := range []struct {
		name  string
		facet g.Facet
	}{
		{"a facet of a kind from another issuer under the same identity and version", issue(sameName("1"), goTarget)},
		{"a facet of the required kind for another target", issue(p.value, tsTarget)},
		{"a facet of a kind with the same identity and another contract version", issue(sameName("2"), goTarget)},
	} {
		composition := p.compose(t, from, []string{key}, c.facet)
		generation, err := gatePlanGo(t, input, composition, bindings)
		if err == nil || generation.Representation().Valid() {
			t.Errorf("%s is accepted (error: %v)", c.name, err)
			continue
		}
		if !errors.Is(err, g.ErrFacet) || !errors.Is(err, g.ErrScope) {
			t.Errorf("%s: got %v, want the facet error", c.name, err)
		}
		if matched := gateRefusalClasses(err); len(matched) != 0 {
			t.Errorf("%s: the error matches the refusal values %v", c.name, matched)
		}
		if generation, err := gatePlanGo(t, input, composition, gateRequest{}); err != nil || !generation.Representation().Valid() {
			t.Errorf("%s: the request for types alone fails: %v", c.name, err)
		}
	}
	// A facet is a validated value, and its kind's owner alone reads it.
	refused := errors.New("the value is not of this kind")
	validated := gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/validated", Version: "1", Validate: func(value any) error {
		if value != "valid" {
			return refused
		}
		return nil
	}}))
	if _, err := validated.Issue(list.Publisher, list, goTarget, "invalid"); !errors.Is(err, g.ErrInvalid) {
		t.Errorf("a kind issues a facet with a value that its validator refuses: %v", err)
	}
	facet, err := validated.Issue(list.Publisher, list, goTarget, "valid")
	if err != nil {
		t.Fatalf("a kind does not issue a facet with a value that its validator admits: %v", err)
	}
	if value, ok := validated.Read(facet); !ok || value != "valid" {
		t.Errorf("the owner of the kind reads %v, %v", value, ok)
	}
	stranger := gateMust[g.FacetKind](t)(g.NewFacetKind(g.FacetKindDescription{Identity: gatePublisher + "/facet/validated", Version: "1"}))
	if value, ok := stranger.Read(facet); ok || value != nil {
		t.Errorf("a kind of another owner under the same identity reads the facet's value: %v", value)
	}
}

// Guarantee 3 (rows 27 and 29), gate line 2.
//
// Shows: a reference is issued by one plan with a category and, for a binder
// and its adapter, the scope of the declaration that binds it. A reference of
// each of the six categories that a plan issues is admitted as its own
// category by its own plan and is rejected as a scope error as each other
// category, by a second plan of the same request, whose content is equal, and
// by the plan parsed from the first one's serialized form. A binder and its
// adapter are rejected in the scope of another declaration and admitted in
// their own. A zero reference is rejected. The plan issues no reference to a
// node it does not hold, to a member a node does not have, or to a binder a
// declaration does not have, and it resolves a member reference to the
// spelled name only after the same check.
// Does not show: a reference to a requirement and its evidence: the plan
// issues none of that category. What a reader does with a reference.
// Red under: G3-reference-ignores-plan, G3-reference-ignores-category,
// G3-reference-ignores-scope, G3-member-unchecked, G3-binder-of-any-index,
// G3-member-of-any-ordinal, G3-reference-to-any-node.
func TestGuarantee3ReferenceHasItsPlanCategoryAndScope(t *testing.T) {
	from, input := gateSource(t, gateFacetModule)
	planned := func() g.RepresentationPlan {
		generation, err := gatePlanGo(t, input, gateGolang(t, from), gateRequest{})
		if err != nil {
			t.Fatal(err)
		}
		return generation.Representation()
	}
	one, other := planned(), planned()
	parsed, err := g.ParsePlan(gateMust[[]byte](t)(one.Serialize()))
	if err != nil {
		t.Fatal(err)
	}
	page, person, record := gateDeclaration(t, one, "Page"), gateDeclaration(t, one, "Person"), gateBody(t, one, "Page")
	if other.Nodes()[0].ID != one.Nodes()[0].ID || len(record.Members) != 1 || page.Binders != 1 {
		t.Fatalf("the two plans differ, or Page is not one record with one member and one binder: %+v", record)
	}
	owned := slices.IndexFunc(one.Symbols(), func(s g.PlanSymbol) bool { return s.Role == "owned" })
	if owned < 0 {
		t.Fatal("the plan has no declaration that a contribution owns")
	}
	must := gateMust[g.PlanReference](t)
	references := map[g.ReferenceCategory]g.PlanReference{
		g.TypeReference:        must(one.TypeOf(record.ID)),
		g.MemberReference:      must(one.MemberOf(record.ID, 0)),
		g.BinderReference:      must(one.BinderOf(page.ID, 0)),
		g.AdapterReference:     must(one.AdapterOf(page.ID, 0)),
		g.DeclarationReference: must(one.DeclarationOf(page.Symbol)),
		g.OwnedReference:       must(one.DeclarationOf(one.Symbols()[owned].ID)),
	}
	scopeOf := func(category g.ReferenceCategory) string {
		if category == g.BinderReference || category == g.AdapterReference {
			return page.Symbol
		}
		return ""
	}
	isScope := func(err error) bool { return errors.Is(err, g.ErrScope) }
	categories := []g.ReferenceCategory{g.TypeReference, g.BinderReference, g.AdapterReference, g.RequirementReference, g.DeclarationReference, g.OwnedReference, g.MemberReference}
	for category, reference := range references {
		if reference.Category() != category {
			t.Errorf("the %s reference has the category %s", category, reference.Category())
		}
		if err := one.Check(reference, category, scopeOf(category)); err != nil {
			t.Errorf("the plan rejects its own %s reference: %v", category, err)
		}
		for name, another := range map[string]g.RepresentationPlan{"a second plan of the same request": other, "the plan parsed from the serialized form": parsed} {
			if err := another.Check(reference, category, scopeOf(category)); !isScope(err) {
				t.Errorf("%s admits the %s reference of the first plan: %v", name, category, err)
			}
		}
		for _, want := range categories {
			if err := one.Check(reference, want, scopeOf(want)); want != category && !isScope(err) {
				t.Errorf("a %s reference is admitted where a %s reference is expected: %v", category, want, err)
			}
		}
	}
	for _, category := range []g.ReferenceCategory{g.BinderReference, g.AdapterReference} {
		for _, scope := range []string{person.Symbol, ""} {
			if err := one.Check(references[category], category, scope); !isScope(err) {
				t.Errorf("the %s of Page is admitted in the scope %q: %v", category, scope, err)
			}
		}
	}
	if err := one.Check(g.PlanReference{}, g.TypeReference, ""); !isScope(err) {
		t.Errorf("a reference that no plan issued is admitted: %v", err)
	}
	// What a plan issues.
	for name, err := range map[string]error{
		"a type reference to a node of no plan":                 second(one.TypeOf("no such node")),
		"a member reference beyond the members of the record":   second(one.MemberOf(record.ID, 1)),
		"a member reference to a negative ordinal":              second(one.MemberOf(record.ID, -1)),
		"a binder reference beyond the binders of Page":         second(one.BinderOf(page.ID, 1)),
		"an adapter reference beyond the binders of Page":       second(one.AdapterOf(page.ID, 1)),
		"a binder reference into a node that is no declaration": second(one.BinderOf(record.ID, 0)),
		"a binder reference into a declaration without binders": second(one.BinderOf(person.ID, 0)),
		"a declaration reference to a symbol of no plan":        second(one.DeclarationOf("source:no such symbol")),
	} {
		if !isScope(err) {
			t.Errorf("the plan issues %s: %v", name, err)
		}
	}
	// A member reference resolves to the spelled name, after the same check.
	if name, err := one.Member(references[g.MemberReference]); err != nil || name != "Items" {
		t.Errorf("the member reference resolves to %q: %v", name, err)
	}
	if name, err := other.Member(references[g.MemberReference]); !isScope(err) || name != "" {
		t.Errorf("another plan resolves the member reference of the first to %q: %v", name, err)
	}
	if name, err := one.Member(references[g.TypeReference]); !isScope(err) || name != "" {
		t.Errorf("a type reference resolves as a member to %q: %v", name, err)
	}
}

// second returns the error of a call that returns a value and an error.
func second[V any](_ V, err error) error { return err }

// gateHooked makes a contribution for the horos source call hook with the
// step of its rule before the rule decides, and gives it the local key key.
func gateHooked(base goContribution, key string, hook func(step horosDecideStep)) goContribution {
	base.Reference.Identity = g.ContributionIdentity{Publisher: gatePublisher, Key: key}
	decide := base.Rule.Decide
	base.Rule.Decide = func(ctx context.Context, step horosDecideStep) (g.Decision, error) {
		hook(step)
		return decide(ctx, step)
	}
	return base
}

// gateAsking makes a record contribution ask for evidence about each of its
// fields, so that its rule has one answer for each of several dependencies.
func gateAsking(base goContribution) goContribution {
	base.Rule.Name = func(_ context.Context, step horosNameStep) ([]horosRequest, error) {
		reader, err := step.Reader()
		if err != nil {
			return nil, err
		}
		fields, err := reader.Fields(step.Subject.Use)
		var requests []horosRequest
		for _, field := range fields {
			requests = append(requests, g.Evidence(field.Use, g.Demand{Kind: gateNilFree}))
		}
		return requests, err
	}
	return base
}

// gatePresented is what a rule got when it presented an answer to its step.
type gatePresented struct {
	called  int
	outcome g.Outcome
	content string
	err     error
}

func (p *gatePresented) present(step horosDecideStep, dependency int, answer g.Answer) {
	p.called++
	p.outcome, p.content, p.err = step.Present(dependency, answer)
}

// Guarantee 4 (rows 7 and 27, O13, O19, O23), gate line 2.
//
// Shows: an answer is read through the step it was issued to, and it is bound
// by value to the asking contribution, the place of the subject, the
// dependency, the whole contribution reference of the dependency and the
// generic environment. A rule that presents an answer that differs from what
// its step expects in exactly one of these is given no outcome and a scope
// error, and the planning of its subject fails with that error although the
// rule drops it. The cases: an answer about another dependency of the same
// subject; an answer issued at another place of the same declaration; an
// answer issued in another generic environment, with two generic
// declarations, once about the binder of the first declaration and once about
// a ground operand inside it; an answer about a dependency whose contribution
// reference differs in one part, for each of its five parts (publisher, local
// key, meaning identity, meaning revision, representation revision); and an
// answer issued to another asking contribution. The step's own answer is read
// with its outcome.
// Does not show: that a rule cannot act on an answer it remembers without
// presenting it: a rule is trusted code of its own package. That an answer
// kept from another run whose five bindings are equal is rejected: the binding
// is by value, and such an answer is read.
// Red under: G4-answer-ignores-dependency, G4-planner-binds-dependency-zero,
// G4-answer-ignores-subject,
// G4-answer-ignores-environment, G4-planner-binds-no-environment,
// G4-answer-ignores-publisher, G4-answer-ignores-key,
// G4-answer-ignores-meaning-identity, G4-answer-ignores-meaning-revision,
// G4-answer-ignores-representation-revision, G4-answer-ignores-asker,
// G4-rejected-answer-does-not-fail.
func TestGuarantee4AnswerBelongsToItsDependency(t *testing.T) {
	const single = "module example.bound;\n\ntype Names = list<text>;\n"
	// plan plans a module with a text contribution that declares the fact and
	// with the two contributions that keep and present answers.
	plan := func(module string, text func(*g.ContributionReference), listKey string, list, record func(horosDecideStep)) error {
		from, input := gateSource(t, module)
		b := gateGolang(t, from)
		stating := gateNilFreeText(t, from, "golang/1", "1")
		if text != nil {
			text(&stating.Reference)
		}
		b = b.replace(stating).replace(gateHooked(gateEvidenceList(t, from, gateEvidenceListOptions{}), listKey, list))
		if record != nil {
			b = b.replace(gateHooked(gateAsking(b.find(t, "horos-to-golang/record")), "asking-record", record))
		}
		generation, err := gatePlanGo(t, input, b.withCodecs(gateDirectCodec(t), gateEmptyCodec(t)), gateRequest{})
		if err == nil && !generation.Representation().Valid() {
			t.Fatal("the planner returned no plan and no error")
		}
		return err
	}
	nothing := func(horosDecideStep) {}
	// keepThenPresent keeps the answer of the first step that runs and
	// presents it to the second.
	keepThenPresent := func(kept *[]g.Answer, presented *gatePresented) func(horosDecideStep) {
		return func(step horosDecideStep) {
			if *kept == nil {
				*kept = step.Answers(0)
				return
			}
			presented.present(step, 0, (*kept)[0])
		}
	}
	// across keeps the answer of one plan and presents it in a second plan in
	// which one thing differs.
	across := func(text func(*g.ContributionReference), listKey string) func(*gatePresented) error {
		return func(presented *gatePresented) error {
			var kept []g.Answer
			if err := plan(single, nil, "evidence-list", func(step horosDecideStep) { kept = step.Answers(0) }, nil); err != nil || len(kept) != 1 {
				t.Fatalf("the first plan kept %d answers: %v", len(kept), err)
			}
			return plan(single, text, listKey, func(step horosDecideStep) { presented.present(step, 0, kept[0]) }, nil)
		}
	}
	within := func(module string) func(*gatePresented) error {
		return func(presented *gatePresented) error {
			var kept []g.Answer
			return plan(module, nil, "evidence-list", keepThenPresent(&kept, presented), nil)
		}
	}
	for _, c := range []struct {
		name string
		run  func(*gatePresented) error
	}{
		{"another dependency of the same subject", func(presented *gatePresented) error {
			return plan("module example.bound;\n\ntype Both = { a: text; b: text; };\n", nil, "evidence-list", nothing, func(step horosDecideStep) {
				presented.present(step, 1, step.Answers(0)[0])
			})
		}},
		{"another place of the same declaration", within("module example.bound;\n\ntype Both = { a: list<text>; b: list<text>; };\n")},
		{"another generic environment, about the binder", within("module example.bound;\n\ntype First<T> = list<T>;\ntype Second<U> = list<U>;\n")},
		{"another generic environment, about a ground operand", within("module example.bound;\n\ntype First<T> = { a: list<text>; b: T; };\ntype Second<U> = { a: list<text>; b: U; };\n")},
		{"another publisher of the dependency's contribution", across(func(r *g.ContributionReference) { r.Identity.Publisher += "/other" }, "evidence-list")},
		{"another local key of the dependency's contribution", across(func(r *g.ContributionReference) { r.Identity.Key += "-other" }, "evidence-list")},
		{"another meaning identity of the dependency's contribution", across(func(r *g.ContributionReference) { r.Meaning.Identity += "/other" }, "evidence-list")},
		{"another meaning revision of the dependency's contribution", across(func(r *g.ContributionReference) { r.Meaning.Revision = "2" }, "evidence-list")},
		{"another representation revision of the dependency's contribution", across(func(r *g.ContributionReference) { r.Representation = "golang/2" }, "evidence-list")},
		{"another asking contribution", across(nil, "another-evidence-list")},
	} {
		var presented gatePresented
		err := c.run(&presented)
		if presented.called != 1 {
			t.Errorf("an answer issued for %s: the rule presented %d answers, so the case shows nothing", c.name, presented.called)
			continue
		}
		if presented.outcome != 0 || presented.content != "" || !errors.Is(presented.err, g.ErrScope) {
			t.Errorf("an answer issued for %s is read by the rule: outcome %q, content %q, error %v", c.name, presented.outcome, presented.content, presented.err)
		}
		if err == nil {
			t.Errorf("an answer issued for %s: the planning succeeds although the answer is none of this step's", c.name)
		} else if !errors.Is(err, g.ErrScope) || !strings.Contains(err.Error(), "was not issued for this dependency") || len(gateRefusalClasses(err)) != 0 {
			t.Errorf("an answer issued for %s: the planning fails with %v, want the scope error of the rejected answer", c.name, err)
		}
	}
	// The control: an answer presented to the step it was issued to is read.
	var own gatePresented
	if err := plan(single, nil, "evidence-list", func(step horosDecideStep) { own.present(step, 0, step.Answers(0)[0]) }, nil); err != nil || own.called != 1 || own.err != nil || own.outcome != g.Proved {
		t.Errorf("the step's own answer is not read: outcome %q, error %v, planning error %v", own.outcome, own.err, err)
	}
}

// Guarantee 4 and row 7, the table of "Evidence": the kernel checks "that an
// answer belongs to the dependency it was asked about: the selected
// representation and its whole contribution reference, the target, profile
// and configuration, the generic environment".
//
// Shows: an answer that a rule kept from a plan under one target
// configuration is rejected when the rule presents it in a plan of the same
// request under another target configuration.
// Does not show: the same for the target and for the profile.
// Red under: no mutant. The test is skipped as PENDING: it is red against the
// kernel as built (finding F-B3), so no mutant can show it red.
func TestGuarantee4AnswerIsBoundToTheTargetConfiguration(t *testing.T) {
	t.Skip("PENDING F-B3: the plan records the target configuration and binds nothing to it; an answer kept from a plan under one target configuration is read in a plan under another")
	from, input := gateSource(t, "module example.bound;\n\ntype Names = list<text>;\n")
	plan := func(value string, hook func(horosDecideStep)) error {
		b := gateGolang(t, from).replace(gateNilFreeText(t, from, "golang/1", "1")).replace(gateHooked(gateEvidenceList(t, from, gateEvidenceListOptions{}), "evidence-list", hook)).withCodecs(gateDirectCodec(t), gateEmptyCodec(t))
		_, err := gatePlanGo(t, input, b, gateRequest{configuration: []g.Setting{{Name: "example.test/gate/setting", Value: value}}})
		return err
	}
	var kept []g.Answer
	if err := plan("1", func(step horosDecideStep) { kept = step.Answers(0) }); err != nil || len(kept) != 1 {
		t.Fatalf("the first plan kept %d answers: %v", len(kept), err)
	}
	var presented gatePresented
	err := plan("2", func(step horosDecideStep) { presented.present(step, 0, kept[0]) })
	if presented.called != 1 {
		t.Fatalf("the rule presented %d answers", presented.called)
	}
	if !errors.Is(presented.err, g.ErrScope) || presented.outcome != 0 {
		t.Errorf("an answer issued under another target configuration is read by the rule: outcome %q, error %v", presented.outcome, presented.err)
	}
	if !errors.Is(err, g.ErrScope) {
		t.Errorf("an answer issued under another target configuration: the planning returns %v, want a scope error", err)
	}
}
